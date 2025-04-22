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
		return nil, err
	}
	if len(roles) != len(roleIds) {
		return nil, errors.New("部分角色不存在")
	}
	// 建立用户角色关联
	if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
		return nil, err
	}

	// 添加用户角色到Casbin
	enforcer := global.Enforcer
	// 先清除用户现有的所有角色权限
	_, err := enforcer.DeleteRolesForUser(cast.ToString(user.Id))
	if err != nil {
		return nil, err
	}

	_, err = enforcer.AddRolesForUser(cast.ToString(user.Id), lo.Map(roles, func(item system.SysRole, index int) string {
		return item.Code
	}))
	if err != nil {
		return nil, err
	}

	err = util.UpdateUserRoleInCache(user.Id, lo.Map(roles, func(item system.SysRole, index int) uint {
		return item.Id
	}))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Create 创建用户
func (s *SysUserService) Create(req request.CreateSysUserReq) (*response.SysUserResp, error) {
	user := &system.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Password: req.Password,
	}
	user.Id = util.NextID()

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	roles, err := s.assignRoles(tx, user, req.RoleIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
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
	existUser := &system.SysUser{}
	err := db.
		Where("id = ?", id).
		First(existUser).
		Error
	if err != nil {
		return err
	}
	return db.
		Delete(&system.SysUser{}, id).
		Error
}

// Update 更新用户
func (s *SysUserService) Update(req request.UpdateSysUserReq) error {
	db := global.DB
	// 先检查用户是否存在
	existUser := &system.SysUser{}
	err := db.
		Where("id = ?", req.Id).
		First(existUser).
		Error
	if err != nil {
		return err
	}
	existUser.Avatar = req.Avatar
	existUser.Email = req.Email
	existUser.Nickname = req.Nickname

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户信息
	if err := tx.Model(existUser).
		Select("Avatar", "Email", "Nickname").
		Where("id = ?", req.Id).
		Updates(existUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	if _, err := s.assignRoles(tx, existUser, req.RoleIds); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Get 查询单个用户
func (s *SysUserService) Get(id int64) (*response.SysUserResp, error) {
	db := global.DB
	user := &system.SysUser{}
	err := db.
		Preload("Roles").
		Where("id = ?", id).
		First(user).
		Error
	if err != nil {
		return nil, err
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
		return common.PageResult[response.SysUserResp]{}, err
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
