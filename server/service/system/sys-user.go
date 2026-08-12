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
// 本函数只负责在事务内落库用户-角色关联表；Casbin 与 Redis 的同步
// 放在事务提交之后由 syncRolePolicy 统一处理。
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

	return roles, nil
}

// syncRolePolicy 在事务提交成功后同步 Casbin g 规则与 Redis 里的角色缓存。
// global.Enforcer 走独立 adapter（EnableAutoSave 直写 global.DB），无法加入业务事务，
// 因此必须等 Commit() 成功后执行——否则 Commit 失败时用户行回滚、权限却已生效，
// 造成「从未创建成功的用户持有真实授权」的幽灵授权。
// 此处后续若再失败，状态是「无权限 / 角色已过期」（fail closed）而非「越权」，
// 并由补偿步骤将失效的半成品规则清掉，只把错误抛给上层提示人工介入。
func (s *SysUserService) syncRolePolicy(user *system.SysUser, roles []system.SysRole) error {
	roleCodes := make([]string, 0, len(roles))
	roleIds := make([]uint, 0, len(roles))
	for _, item := range roles {
		roleCodes = append(roleCodes, item.Code)
		roleIds = append(roleIds, item.Id)
	}

	enforcer := global.Enforcer
	// 无条件先清空该用户的全部角色，再按新集合重建
	if _, err := enforcer.DeleteRolesForUser(user.UserId); err != nil {
		return errors.New("roleAssignFailed")
	}

	if len(roleCodes) > 0 {
		if _, err := enforcer.AddRolesForUser(user.UserId, roleCodes); err != nil {
			// 补偿：清掉可能已加入的部分规则，回退到无权限态
			_, _ = enforcer.DeleteRolesForUser(user.UserId)
			return errors.New("roleAssignFailed")
		}
	}

	if err := util.UpdateUserRoleInCache(user.UserId, roleIds); err != nil {
		// 补偿：撤销 Casbin 侧刚写入的规则，保持一致的无权限态
		_, _ = enforcer.DeleteRolesForUser(user.UserId)
		return errors.New("roleAssignFailed")
	}

	return nil
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

	// Casbin 与 Redis 在事务提交成功之后同步
	if err := s.syncRolePolicy(user, roles); err != nil {
		return nil, err
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

	// 账号已被删除，其全部在线会话随 token 一并吊销，推送强制下线
	notifySessionRevoked(user.UserId, "您的账号已被删除")

	return nil
}

// Update 更新用户
func (s *SysUserService) Update(db *gorm.DB, id string, req request.UpdateSysUserReq) error {
	var existUser *system.SysUser
	existUser, err := s.checkUserExist(db, id)
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

	roles, err := s.assignRoles(tx, existUser, req.RoleIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("updateUserFailed")
	}

	// 事务提交成功后同步 Casbin 与 Redis 角色
	if err := s.syncRolePolicy(existUser, roles); err != nil {
		return err
	}

	return nil
}

// Get 查询单个用户
func (s *SysUserService) Get(db *gorm.DB, id string) (*response.SysUserResp, error) {
	// 一次查询带出角色，不再先查存在性再重查同一行
	var user system.SysUser
	result := db.
		Unscoped().
		Preload("Roles").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		return nil, errors.New(util.TranslateDBError(result.Error, "userNotFound"))
	}
	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("userDeleted")
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

	// Count 用独立 builder，避免污染后续 Find 的状态
	countDB := db.Model(&system.SysUser{})
	if err := countDB.Count(&total).Error; err != nil {
		return common.PageResult[response.SysUserResp]{}, errors.New("getUserListFailed")
	}
	req.Normalize()
	err := db.
		Model(&system.SysUser{}).
		Preload("Roles").
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id ASC").
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
func (s *SysUserService) ResetPassword(db *gorm.DB, id string, req request.UpdateSysUserPasswordReq, oldPassword *string, accessToken string) error {
	var existUser *system.SysUser
	existUser, err := s.checkUserExist(db, id)
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
		// 管理员代改密码无法证明是本人操作，吊销全部会话并推送强制下线
		notifySessionRevoked(existUser.UserId, "您的密码已被管理员重置，请重新登录")
	} else {
		if err := util.RevokeAllUserTokensExcept(existUser.UserId, accessToken); err != nil {
			return errors.New("resetPasswordFailed")
		}
	}

	return nil
}
