package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"gorm.io/gorm"
)

type UserService struct{}

// 创建用户
func (u *UserService) CreateUser(request request.UserRequest, createdBy uint) (err error) {
	request.EncryptPassword()

	db := global.DB

	// TODO: 检验部门、角色、岗位是否存在
	var dept system.Dept
	err = db.
		Model(system.Dept{}).
		First(&dept, request.DeptId).
		Error
	if err != nil {
		return errors.New("该部门不存在")
	}

	user := system.User{
		Username: request.Username,
		Password: request.Password,
		DeptId:   request.DeptId,
		PostId:   request.PostId,
		RoleId:   request.RoleId,
		DeptIds:  []int{request.DeptId},
		PostIds:  []int{request.PostId},
		RoleIds:  []int{request.RoleId},
		ControlWrapper: model.ControlWrapper{
			CreatedBy: createdBy,
		},
	}

	if request.Nickname != "" {
		user.Nickname = request.Nickname
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Sex != 0 {
		user.Sex = request.Sex
	}
	if request.Status != 0 {
		user.Status = request.Status
	}
	if request.Remark != "" {
		user.Remark = request.Remark
	}

	err = db.
		Create(&user).
		Error
	return
}

// 删除用户
func (u *UserService) DeleteUser(param request.UserParam, deletedBy uint) (err error) {
	db := global.DB

	var user system.User
	result := db.
		Unscoped().
		First(&user, param.UserId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("用户不存在")
		}
		return result.Error
	}

	if !user.DeletedAt.Time.IsZero() {
		return errors.New("用户已删除")
	}

	db.
		Model(system.User{}).
		Where("id = ?", param.UserId).
		Update("deleted_by", deletedBy).
		Delete(&user)
	return nil
}

// 更新用户
func (u *UserService) UpdateUser(param request.UserParam, request request.UserRequest, updatedBy uint) (err error) {
	db := global.DB

	var user system.User
	result := db.
		First(&user, param.UserId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("用户不存在")
		}
		return result.Error
	}

	// TODO: 检验部门、角色、岗位是否存在
	var dept system.Dept
	err = db.
		Model(system.Dept{}).
		First(&dept, request.DeptId).
		Error
	if err != nil {
		return errors.New("该部门不存在")
	}

	user.Nickname = request.Nickname
	user.DeptId = request.DeptId
	user.Phone = request.Phone
	user.Email = request.Email
	user.Username = request.Username
	user.Sex = request.Sex
	user.Status = request.Status
	user.PostId = request.PostId
	user.RoleId = request.RoleId
	user.Remark = request.Remark
	user.Avatar = request.Avatar
	user.ControlWrapper = model.ControlWrapper{
		UpdatedBy: updatedBy,
	}

	db.
		Model(system.User{}).
		Where("id = ?", param.UserId).
		Updates(&user).
		Omit("id", "created_at")

	return nil
}

// 查询单个用户
func (u *UserService) GetUser(param request.UserParam) (response.UserItem, error) {
	db := global.DB

	var userRep response.UserItem
	var user system.User
	result := db.
		Preload("Dept").
		First(&user, param.UserId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return userRep, errors.New("用户不存在")
		}
		return userRep, result.Error
	}

	return response.UserItem{
		IDWrapper:      user.IDWrapper,
		ControlWrapper: user.ControlWrapper,
		UserItem: request.UserItem{
			Avatar:   user.Avatar,
			Nickname: user.Nickname,
			DeptId:   user.DeptId,
			Phone:    user.Phone,
			Email:    user.Email,
			Username: user.Username,
			Sex:      user.Sex,
			PostId:   user.PostId,
			RoleId:   user.RoleId,
			RemarkWrapper: model.RemarkWrapper{
				Remark: user.Remark,
			},
			StatusWrapper: model.StatusWrapper{
				Status: user.Status,
			},
		},
		DeptIds: user.DeptIds,
		Dept: response.DeptItem{
			IDWrapper:      user.Dept.IDWrapper,
			ControlWrapper: user.Dept.ControlWrapper,
			DeptRequest: request.DeptRequest{
				ParentId:    user.Dept.ParentId,
				Name:        user.Dept.Name,
				Leader:      user.Dept.Leader,
				Phone:       user.Dept.Phone,
				Email:       user.Dept.Email,
				SortWrapper: user.Dept.SortWrapper,
				StatusWrapper: model.StatusWrapper{
					Status: user.Dept.Status,
				},
			},
		},
		PostIds: user.PostIds,
		RoleIds: user.RoleIds,
	}, nil
}

// 查询用户
func (u *UserService) GetUserList(query request.UserQueryParams) (response.UserResponse, error) {
	db := global.DB

	var total int64
	var originUsers []system.User
	offset := (query.Page - 1) * query.PageSize
	users := response.UserResponse{}
	err := db.
		Model(&system.User{}).
		Order("id ASC").Preload("Dept").
		Where(fmt.Sprintf("username LIKE '%%%s%%'", query.Name)).
		Count(&total).Offset(offset).
		Limit(query.PageSize).
		Find(&originUsers).
		Error
	if err != nil {
		return users, err
	}
	users.Meta.Page = query.Page
	users.Meta.PageSize = query.PageSize
	users.Meta.Total = total

	users.Data = make([]response.UserItem, len(originUsers))
	for i, user := range originUsers {
		users.Data[i].ID = user.ID
		users.Data[i].ControlWrapper = user.ControlWrapper
		users.Data[i].Nickname = user.Nickname
		users.Data[i].Phone = user.Phone
		users.Data[i].Email = user.Email
		users.Data[i].Username = user.Username
		users.Data[i].Sex = user.Sex
		users.Data[i].Status = user.Status
		users.Data[i].Remark = user.Remark
		users.Data[i].DeptId = user.DeptId
		users.Data[i].PostId = user.PostId
		users.Data[i].RoleId = user.RoleId
		users.Data[i].DeptIds = user.DeptIds
		users.Data[i].PostIds = user.PostIds
		users.Data[i].RoleIds = user.RoleIds
		users.Data[i].Dept = response.DeptItem{
			IDWrapper:      user.Dept.IDWrapper,
			ControlWrapper: user.Dept.ControlWrapper,
			DeptRequest: request.DeptRequest{
				ParentId:    user.Dept.ParentId,
				Name:        user.Dept.Name,
				Leader:      user.Dept.Leader,
				Phone:       user.Dept.Phone,
				Email:       user.Dept.Email,
				SortWrapper: user.Dept.SortWrapper,
				StatusWrapper: model.StatusWrapper{
					Status: user.Dept.Status,
				},
			},
		}
	}

	return users, nil
}
