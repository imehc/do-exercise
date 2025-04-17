package system

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
)

type SysUserService struct{}

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
	err := global.DB.Create(user).Error
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
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}, nil
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
func (s *SysUserService) Update(req request.UpdateSysUserReq) error {
	// 先检查用户是否存在
	existUser := &system.SysUser{}
	err := global.DB.Where("id = ?", req.Id).
		First(existUser).
		Error
	if err != nil {
		return err
	}
	existUser.Avatar = req.Avatar
	existUser.Email = req.Email
	existUser.Nickname = req.Nickname

	return global.DB.Model(existUser).
		Select("Avatar", "Email", "Nickname").
		Where("id = ?", req.Id).
		Updates(existUser).
		Error
}

// Get 查询单个用户
func (s *SysUserService) Get(id int64) (*response.SysUserResp, error) {
	user := &system.SysUser{}
	err := global.DB.Where("id =?", id).First(user).Error
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
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}, nil
}

// GetList 查询用户列表
func (s *SysUserService) GetList(req common.Pagination) (common.PageResult[response.SysUserResp], error) {
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
			CreatedBy: user.CreatedBy,
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
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
