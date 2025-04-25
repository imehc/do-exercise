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
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type SysUserService struct{}

// assignRoles 分配角色
func (s *SysUserService) assignRoles(tx *gorm.DB, user *system.SysUser, roleIds []uint) ([]system.SysRole, error) {
	if len(roleIds) == 0 {
		return []system.SysRole{}, nil
	}
	var roles []system.SysRole
	// 检查角色是否存在
	if err := tx.Where("id IN ?", roleIds).Find(&roles).Error; err != nil {
		return nil, errors.New("allRolesNotFound")
	}
	if len(roles) != len(roleIds) {
		return nil, errors.New("roleNotFound")
	}
	// 建立用户角色关联
	if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	// 添加用户角色到Casbin
	enforcer := global.Enforcer
	// 先清除用户现有的所有角色权限
	_, err := enforcer.DeleteRolesForUser(cast.ToString(user.Id))
	if err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	_, err = enforcer.AddRolesForUser(cast.ToString(user.Id), lo.Map(roles, func(item system.SysRole, index int) string {
		return item.Code
	}))
	if err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	err = util.UpdateUserRoleInCache(user.Id, lo.Map(roles, func(item system.SysRole, index int) uint {
		return item.Id
	}))
	if err != nil {
		return nil, errors.New("roleAssignFailed")
	}

	return roles, nil
}

// checkUserExist 检查用户是否存在
func (s *SysUserService) checkUserExist(db *gorm.DB, userId int64) (*system.SysUser, error) {
	var user *system.SysUser
	result := db.
		Unscoped().
		First(&user, userId)
	if result.Error != nil {
		return nil, errors.New("userNotFound")
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
func (s *SysUserService) Create(req request.CreateSysUserReq) (*response.SysUserResp, error) {
	db := global.DB
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
	user.Id = util.NextID()

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
		Id:        user.Id,
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

// Delete 删除用户
func (s *SysUserService) Delete(id int64) error {
	db := global.DB
	// 先检查用户是否存在
	user, err := s.checkUserExist(db, id)
	if err != nil {
		return err
	}
	err = db.
		Delete(user, id).
		Error
	if err != nil {
		return errors.New("deleteUserFailed")
	}
	return nil
}

// Update 更新用户
func (s *SysUserService) Update(req request.UpdateSysUserReq) error {
	db := global.DB
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
func (s *SysUserService) Get(id int64) (*response.SysUserResp, error) {
	db := global.DB
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
		Id:        user.Id,
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
func (s *SysUserService) GetList(req common.Pagination) (common.PageResult[response.SysUserResp], error) {
	var users []system.SysUser
	var total int64
	db := global.DB.
		Model(&system.SysUser{})
	db.
		Count(&total)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	db = db.
		Preload("Roles").
		Scopes(util.Paginate(req.PageSize, req.Page))
	err := db.
		Find(&users).
		Error
	if err != nil {
		return common.PageResult[response.SysUserResp]{}, errors.New("getUserListFailed")
	}
	data := make([]response.SysUserResp, len(users))
	for i, user := range users {
		data[i] = response.SysUserResp{
			Id:        user.Id,
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
func (s *SysUserService) ResetPassword(req request.UpdateSysUserPasswordReq, oldPassword *string) error {
	db := global.DB
	var existUser *system.SysUser
	existUser, err := s.checkUserExist(db, req.Id)
	if err != nil {
		return err
	}

	if *oldPassword != "" {
		hash := util.Hash{Value: existUser.Password}
		if !hash.Compare(req.Password) {
			return errors.New("passwordError")
		}
	}

	hash := util.Hash{Value: req.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}
	existUser.Password = password

	// 更新密码
	if err := db.
		Model(existUser).
		Select("Password").
		Updates(existUser).
		Error; err != nil {
		return errors.New("resetPasswordFailed")
	}

	return nil
}
