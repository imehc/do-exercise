package system

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
)

type SysUserService struct{}

// Create 创建用户
func (s *SysUserService) Create(req system.SysUser) (*system.SysUser, error) {
	user := &system.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Password: req.Password,
	}
	user.Id = util.NextID()
	err := global.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Delete 删除用户
func (s *SysUserService) Delete(id int64) error {
	// 先检查用户是否存在
	existUser := &system.SysUser{}
	err := global.DB.Where("id = ?", id).
		First(existUser).
		Error
	if err != nil {
		return err
	}
	return global.DB.
		Delete(&system.SysUser{}, id).
		Error
}

// Update 更新用户
func (s *SysUserService) Update(req system.SysUser, id int64) error {
	// 先检查用户是否存在
	existUser := &system.SysUser{}
	err := global.DB.Where("id = ?", id).
		First(existUser).
		Error
	if err != nil {
		return err
	}
	return global.DB.Model(&req).
		Select("Avatar", "Email", "Nickname").
		Where("id = ?", id).
		Updates(req).
		Error
}

// Get 查询单个用户
func (s *SysUserService) Get(id int64) (*system.SysUser, error) {
	user := &system.SysUser{}
	err := global.DB.Where("id =?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetList 查询用户列表
func (s *SysUserService) GetList(req common.Pagination) (common.PageResult[system.SysUser], error) {
	var users []system.SysUser
	var total int64
	db := global.DB.Model(&system.SysUser{})
	db.Count(&total)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	db = db.Scopes(util.Paginate(req.PageSize, req.Page))
	err := db.Find(&users).Error
	if err != nil {
		return common.PageResult[system.SysUser]{}, err
	}
	result := common.PageResult[system.SysUser]{
		Data: users,
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}
	return result, nil
}
