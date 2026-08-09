package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type SysUserService struct{}

// assignRoles 分配角色。
// 空 roleIds 的语义是「撤销该用户的全部角色」，而不是「不做处理」——
// 早期实现在此提前 return，导致 role_ids:[] 的解除权限请求返回成功，
// 但 Casbin 中的 g 规则原样保留，用户实际权限不变。
func (s *SysUserService) assignRoles(tx *gorm.DB, user *system.SysUser, roleIds []uint) ([]system.SysRole, error) {
	var roles []system.SysRole
	if len(roleIds) > 0 {
		// 检查角色是否存在
		if err := tx.Where("id IN ?", roleIds).Find(&roles).Error; err != nil {
			return nil, errors.New("allRolesNotFound")
		}
		if len(roles) != len(roleIds) {
			return nil, errors.New("roleNotFound")
		}
	}

	// 建立/清空用户角色关联
	if len(roles) == 0 {
		if err := tx.Model(user).Association("Roles").Clear(); err != nil {
			return nil, errors.New("roleAssignFailed")
		}
	} else if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	// 同步 Casbin：无条件先清空该用户的全部角色，再按新集合重建
	enforcer := global.Enforcer
	if _, err := enforcer.DeleteRolesForUser(user.UserId); err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	if len(roles) > 0 {
		if _, err := enforcer.AddRolesForUser(user.UserId, lo.Map(roles, func(item system.SysRole, index int) string {
			return item.Code
		})); err != nil {
			return nil, errors.New("roleAssignFailed")
		}
	}

	if err := util.UpdateUserRoleInCache(user.UserId, lo.Map(roles, func(item system.SysRole, index int) uint {
		return item.Id
	})); err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	return roles, nil
}

// checkUserExist 检查用户是否存在
func (s *SysUserService) checkUserExist(db *gorm.DB, userId string) (*system.SysUser, error) {
	if userId == "" {
		return nil, errors.New("idCannotBeEmpty")
	}
	var user *system.SysUser
	result := db.
		Unscoped().
		Where("id = ?", userId).
		First(&user)
	if result.Error != nil {
		return nil, errors.New(util.TranslateDBError(result.Error, "userNotFound"))
	}

	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("userDeleted")
	}

	return user, nil
}

// checkUserNameDuplication 检查用户名是否重复
func (s *SysUserService) checkUserNameDuplication(db *gorm.DB, username string) error {
	var count int64
	if err := db.Model(&system.SysUser{}).Where("username = ?", username).Count(&count).Error; err != nil || count > 0 {
		return errors.New("usernameExists")
	}
	return nil
}

// checkEmailDuplication 检查邮箱是否重复
func (s *SysUserService) checkEmailDuplication(db *gorm.DB, email string) error {
	var count int64
	if err := db.Model(&system.SysUser{}).Where("email =?", email).Count(&count).Error; err != nil || count > 0 {
		return errors.New("emailExists")
	}
	return nil
}

// Create 创建用户
func (s *SysUserService) Create(db *gorm.DB, req request.CreateSysUserReq) (*response.SysUserResp, error) {
	err := s.checkUserNameDuplication(db, req.Username)
	if err != nil {
		return nil, err
	}
	err = s.checkEmailDuplication(db, req.Email)
	if err != nil {
		return nil, err
	}

	user := &system.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Password: req.Password,
	}
	user.UserId = util.NextID()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err = tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createUserFailed")
	}

	roles, err := s.assignRoles(tx, user, req.RoleIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createUserFailed")
	}
	return &response.SysUserResp{
		Id:        user.UserId,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Roles: lo.Map(roles, func(item system.SysRole, index int) response.SysRoleResp {
			return response.SysRoleResp{
				Id:   item.Id,
				Name: item.Name,
				Code: item.Code,
			}
		}),
	}, nil
}

// Delete 删除用户。
// 除软删除记录外，还必须解除 Casbin 角色绑定并吊销全部会话——
// 否则被删用户的 access token 在 TTL 内仍然有效，refresh token 更能持续换发新令牌，
// 且 casbin_rule 中的 g 规则会永久残留（用户 ID 复用时会泄漏给新用户）。
func (s *SysUserService) Delete(db *gorm.DB, id string) error {
	// 先检查用户是否存在
	user, err := s.checkUserExist(db, id)
	if err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id = ?", id).Delete(user).Error; err != nil {
		tx.Rollback()
		return errors.New("deleteUserFailed")
	}

	// 解除用户与角色的关联
	if err := tx.Model(user).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return errors.New("deleteUserFailed")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("deleteUserFailed")
	}

	// 以下两步作用于 Casbin 与 Redis，不在上面的事务范围内，
	// 因此放在提交成功之后执行，避免事务回滚后权限已被误撤销。
	if _, err := global.Enforcer.DeleteRolesForUser(user.UserId); err != nil {
		return errors.New("deleteUserFailed")
	}

	if err := util.RevokeAllUserTokens(user.UserId); err != nil {
		return errors.New("deleteUserFailed")
	}

	return nil
}

// Update 更新用户
func (s *SysUserService) Update(db *gorm.DB, req request.UpdateSysUserReq) error {
	var existUser *system.SysUser
	existUser, err := s.checkUserExist(db, req.Id)
	if err != nil {
		return err
	}
	existUser.Avatar = req.Avatar
	existUser.Nickname = req.Nickname

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户信息
	if err := tx.
		Model(existUser).
		Select("Avatar", "Nickname").
		Updates(existUser).
		Error; err != nil {
		tx.Rollback()
		return errors.New("updateUserFailed")
	}

	if _, err := s.assignRoles(tx, existUser, req.RoleIds); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("updateUserFailed")
	}
	return nil
}

// Get 查询单个用户
func (s *SysUserService) Get(db *gorm.DB, id string) (*response.SysUserResp, error) {
	// 先检查用户是否存在
	_, err := s.checkUserExist(db, id)
	if err != nil {
		return nil, err
	}

	var user system.SysUser
	err = db.
		Preload("Roles").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, errors.New("getUserFailed")
	}

	return &response.SysUserResp{
		Id:        user.UserId,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Roles: lo.Map(user.Roles, func(item system.SysRole, index int) response.SysRoleResp {
			return response.SysRoleResp{
				Id:   item.Id,
				Name: item.Name,
				Code: item.Code,
			}
		}),
	}, nil
}

// GetList 查询用户列表
func (s *SysUserService) GetList(db *gorm.DB, req common.Pagination) (common.PageResult[response.SysUserResp], error) {
	var users []system.SysUser
	var total int64
	db = db.
		Model(&system.SysUser{}).
		Count(&total)
	req.Normalize()
	db = db.
		Preload("Roles").
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id ASC")
	err := db.
		Find(&users).
		Error
	if err != nil {
		return common.PageResult[response.SysUserResp]{}, errors.New("getUserListFailed")
	}
	data := make([]response.SysUserResp, len(users))
	for i, user := range users {
		data[i] = response.SysUserResp{
			Id:        user.UserId,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Roles: lo.Map(user.Roles, func(item system.SysRole, index int) response.SysRoleResp {
				return response.SysRoleResp{
					Id:   item.Id,
					Name: item.Name,
					Code: item.Code,
				}
			}),
		}
	}
	result := common.PageResult[response.SysUserResp]{
		Data: data,
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}
	return result, nil
}

// ResetPassword 重置密码
func (s *SysUserService) ResetPassword(db *gorm.DB, req request.UpdateSysUserPasswordReq, oldPassword *string, accessToken string) error {
	var existUser *system.SysUser
	existUser, err := s.checkUserExist(db, req.Id)
	if err != nil {
		return err
	}

	if oldPassword != nil {
		hash := util.Hash{Value: existUser.Password}
		if !hash.Compare(*oldPassword) {
			return errors.New("passwordError")
		}
	}

	hash := util.Hash{Value: req.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}
	existUser.Password = password

	// 管理员代改（不知道原密码）时置 must_change_password，强制用户下次登录改密；
	// 用户自己验证原密码改密时清除该标记。
	if oldPassword == nil {
		existUser.MustChangePassword = true
	} else {
		existUser.MustChangePassword = false
	}

	// 更新密码
	if err := db.
		Model(existUser).
		Select("Password", "MustChangePassword").
		Updates(existUser).
		Error; err != nil {
		return errors.New("resetPasswordFailed")
	}

	// 改密后处理既有会话：
	// - 用户自己验证原密码改密：保留当前会话（免重新登录），吊销其他会话，并清除
	//   保留会话上的强制改密标记，避免 AuthMiddleware 继续拦截。
	// - 管理员代改（不知道原密码）：无法证明是本人操作，吊销全部会话。
	if oldPassword == nil {
		if err := util.RevokeAllUserTokens(existUser.UserId); err != nil {
			return errors.New("resetPasswordFailed")
		}
	} else {
		if err := util.RevokeAllUserTokensExcept(existUser.UserId, accessToken); err != nil {
			return errors.New("resetPasswordFailed")
		}
	}

	return nil
}
