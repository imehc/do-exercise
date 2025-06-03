package common

import (
	"errors"
	"slices"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/common/response"
	sysReq "github.com/imehc/do-exercise/server/model/system/request"
	"github.com/samber/lo"

	"github.com/imehc/do-exercise/server/model/system"
	sysService "github.com/imehc/do-exercise/server/service/system"
)

type UserService struct{}

// 通过邮箱查询用户
func (u *UserService) FindUserByEmail(email string) (*system.SysUser, error) {
	db := global.DB
	var user *system.SysUser
	error := db.
		Where("email = ?", email).
		First(&user).
		Error
	if error != nil {
		return nil, errors.New("userNotFound")
	}
	return user, nil
}

// BindEmail 绑定邮箱
func (u *UserService) BindEmail(req request.BindEmailReq) error {
	db := global.DB

	if err := db.
		Model(system.SysUser{}).
		Where("id = ?", req.Id).
		Update("Email", req.Email).
		Error; err != nil {
		return errors.New("bindEmailFailed")
	}
	return nil
}

// UpdatePassword 修改密码
func (u *UserService) UpdatePassword(req request.UserModifyPasswordReq) error {
	var sysUserService sysService.SysUserService
	return sysUserService.ResetPassword(
		sysReq.UpdateSysUserPasswordReq{
			Id:       req.Id,
			Password: req.Password,
		},
		&req.OldPassword,
	)
}

// UpdateProfile 修改用户基本信息
func (u *UserService) UpdateProfile(req request.UserModifyProfileReq) error {
	db := global.DB
	var user *system.SysUser
	error := db.
		First(&user, req.Id).
		Error
	if error != nil {
		return errors.New("userNotFound")
	}

	user.Avatar = req.Avatar
	user.Nickname = req.Nickname

	if err := db.
		Model(user).
		Select("Avatar", "Nickname").
		Updates(user).
		Error; err != nil {
		return errors.New("updateFailed")
	}
	return nil
}

// GetProfile 获取用户基本信息
func (u *UserService) GetProfile(id string) (*response.UserProfile, error) {
	db := global.DB
	var user *system.SysUser
	error := db.
		Preload("Roles").
		First(&user, id).
		Error
	if error != nil {
		return nil, errors.New("userNotFound")
	}

	return &response.UserProfile{
		Id:       user.UserId,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Roles: lo.Map(user.Roles, func(item system.SysRole, index int) response.Role {
			return response.Role{
				Id:   item.Id,
				Name: item.Name,
			}
		}),
	}, nil
}

// GetMenu 获取用户菜单

func (u *UserService) GetMenu(id string) (*[]response.UserMenu, error) {
	db := global.DB
	var user *system.SysUser

	// 查询用户及角色菜单
	err := db.
		Preload("Roles.Menus").
		First(&user, id).
		Error
	if err != nil {
		return nil, errors.New("userNotFound")
	}

	// 初始化菜单列表
	var menus []system.SysMenu
	for _, role := range user.Roles {
		menus = append(menus, role.Menus...)
	}
	menus = lo.UniqBy(menus, func(menu system.SysMenu) uint {
		return menu.Id
	})

	// 构建菜单map，避免重复查询
	menuMap := make(map[uint]system.SysMenu)
	for _, menu := range menus {
		menuMap[menu.Id] = menu
	}

	// 递归向上查找父菜单
	toCheck := menus
	for len(toCheck) > 0 {
		var nextToCheck []system.SysMenu
		for _, menu := range toCheck {
			if *menu.ParentId != 0 {
				if _, exists := menuMap[*menu.ParentId]; !exists {
					var parentMenu system.SysMenu
					if err := db.First(&parentMenu, menu.ParentId).Error; err == nil {
						menuMap[parentMenu.Id] = parentMenu
						nextToCheck = append(nextToCheck, parentMenu)
					}
				}
			}
		}
		toCheck = nextToCheck
	}

	// 将 map 转为切片
	finalMenus := make([]system.SysMenu, 0, len(menuMap))
	for _, menu := range menuMap {
		finalMenus = append(finalMenus, menu)
	}

	// 转为响应结构
	userMenus := make([]response.UserMenu, 0, len(finalMenus))
	for _, menu := range finalMenus {
		userMenus = append(userMenus, response.UserMenu{
			Id:         menu.Id,
			Name:       menu.Name,
			ParentId:   menu.ParentId,
			Permission: menu.Permission,
			Icon:       menu.Icon,
			Type:       menu.Type,
			Route:      menu.Route,
			Component:  menu.Component,
		})
	}

	// 按照 ID 升序排序
	slices.SortFunc(userMenus, func(a, b response.UserMenu) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})

	return &userMenus, nil
}
