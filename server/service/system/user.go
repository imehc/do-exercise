package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/pkg/utils/scope"
	"github.com/imehc/do-exercise/server/utils"
	"gorm.io/gorm"
)

type UserService struct{}

// 创建用户
func (u *UserService) CreateUser(request request.UserRequest, createBy uint) (err error) {
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
		ControlWrapper: model.ControlWrapper{
			CreateBy: createBy,
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
func (u *UserService) DeleteUser(param request.UserParam, deleteBy uint) (err error) {
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

	if !user.DeleteAt.Time.IsZero() {
		return errors.New("用户已删除")
	}

	// 检查是否有权限删除
	if user.CreateBy != deleteBy {
		// 检查是否是超级管理员
		var role system.Role
		if err := db.First(&role, user.RoleId).Error; err != nil {
			return err
		}
		if !role.IsAdmin {
			return errors.New("无权删除其他用户创建的数据")
		}
	}

	db.
		Model(system.User{}).
		Where("user_id = ?", param.UserId).
		Update("delete_by", deleteBy).
		Delete(&user)
	return nil
}

// 更新用户
func (u *UserService) UpdateUser(param request.UserParam, request request.UserRequest, updateBy uint) (err error) {
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

	// 检查是否有权限更新
	if user.CreateBy != updateBy {
		// 检查是否是超级管理员
		var role system.Role
		if err := db.First(&role, user.RoleId).Error; err != nil {
			return err
		}
		if !role.IsAdmin {
			return errors.New("无权更新其他用户创建的数据")
		}
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
		UpdateBy: updateBy,
	}

	db.
		Model(system.User{}).
		Where("user_id = ?", param.UserId).
		Updates(&user).
		Omit("user_id", "created_at")

	return nil
}

// 查询单个用户
func (u *UserService) GetUser(param request.UserParam) (response.UserItem, error) {
	db := global.DB

	var userRep response.UserItem
	var user system.User
	result := db.
		Preload("Dept").
		Preload("Post").
		First(&user, param.UserId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return userRep, errors.New("用户不存在")
		}
		return userRep, result.Error
	}

	userRep = response.UserItem{
		IDWrapper: common.IDWrapper{
			ID: user.UserId,
		},
		ControlWrapper: user.ControlWrapper,
		UserItem: request.UserItem{
			Avatar:   user.Avatar,
			Nickname: user.Nickname,
			Phone:    user.Phone,
			Email:    user.Email,
			Username: user.Username,
			Sex:      user.Sex,
			RemarkWrapper: model.RemarkWrapper{
				Remark: user.Remark,
			},
			StatusWrapper: model.StatusWrapper{
				Status: user.Status,
			},
		},
	}
	hasDeptId := user.Dept.DeptId != 0
	hasPostId := user.Post.PostId != 0
	// hasRoleId := user.RoleId != 0
	if hasDeptId {
		userRep.DeptId = user.DeptId
		userRep.Dept = response.DeptItem{
			IDWrapper: common.IDWrapper{
				ID: user.Dept.DeptId,
			},
			ControlWrapper: user.Dept.ControlWrapper,
			DeptRequest: request.DeptRequest{
				ParentId:    user.Dept.ParentId,
				Name:        user.Dept.Name,
				SortWrapper: user.Dept.SortWrapper,
				StatusWrapper: model.StatusWrapper{
					Status: user.Dept.Status,
				},
			},
		}
	}

	if hasPostId {
		userRep.PostId = user.PostId
		userRep.Post = response.PostItem{
			IDWrapper: common.IDWrapper{
				ID: user.Post.PostId,
			},
			ControlWrapper: user.Post.ControlWrapper,
			PostRequest: request.PostRequest{
				Name:          user.Post.Name,
				Code:          user.Post.Code,
				SortWrapper:   user.Post.SortWrapper,
				StatusWrapper: user.Post.StatusWrapper,
				RemarkWrapper: user.Post.RemarkWrapper,
			},
		}
	}

	// if hasRoleId {
	// 	// TODO: 对应角色信息
	// }

	return userRep, nil
}

// 查询用户
func (u *UserService) GetUserList(query request.UserQueryParams, s common.ScopeData) (response.UserResponse, error) {
	db := global.DB
	// 应用数据权限过滤
	db = scope.GetDataScope(db, &s, "sys_user")

	var total int64
	var originUsers []system.User
	users := response.UserResponse{}
	err := db.
		Model(&system.User{}).
		Order("user_id ASC").
		Preload("Dept").
		Preload("Post").
		Where(fmt.Sprintf("username LIKE '%%%s%%'", query.Name)).
		Count(&total).
		Scopes(utils.Paginate(query.PageSize, query.Page)).
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
		users.Data[i].ID = user.UserId
		users.Data[i].ControlWrapper = user.ControlWrapper
		users.Data[i].Nickname = user.Nickname
		users.Data[i].Phone = user.Phone
		users.Data[i].Email = user.Email
		users.Data[i].Username = user.Username
		users.Data[i].Sex = user.Sex
		users.Data[i].Status = user.Status
		users.Data[i].Remark = user.Remark

		hasDeptId := user.Dept.DeptId != 0
		hasPostId := user.Post.PostId != 0
		// hasRoleId := user.RoleId != 0
		if hasDeptId {
			users.Data[i].DeptId = user.DeptId
			users.Data[i].Dept = response.DeptItem{
				IDWrapper: common.IDWrapper{
					ID: user.Dept.DeptId,
				},
				ControlWrapper: user.Dept.ControlWrapper,
				DeptRequest: request.DeptRequest{
					ParentId:    user.Dept.ParentId,
					Name:        user.Dept.Name,
					SortWrapper: user.Dept.SortWrapper,
					StatusWrapper: model.StatusWrapper{
						Status: user.Dept.Status,
					},
				},
			}
		}
		if hasPostId {
			users.Data[i].PostId = user.PostId
			users.Data[i].Post = response.PostItem{
				IDWrapper: common.IDWrapper{
					ID: user.Post.PostId,
				},
				ControlWrapper: user.Post.ControlWrapper,
				PostRequest: request.PostRequest{
					Name:          user.Post.Name,
					Code:          user.Post.Code,
					SortWrapper:   user.Post.SortWrapper,
					StatusWrapper: user.Post.StatusWrapper,
					RemarkWrapper: user.Post.RemarkWrapper,
				},
			}
		}
		// if hasRoleId {
		// 	// TODO: 对应角色信息
		// }
	}

	return users, nil
}
