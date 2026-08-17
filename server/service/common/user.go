package common

import (
	"errors"
	"slices"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/common/response"
	sysReq "github.com/imehc/do-exercise/server/model/system/request"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/imehc/do-exercise/server/model/system"
	sysService "github.com/imehc/do-exercise/server/service/system"
)

type UserService struct{}

// 通过邮箱查询用户
func (u *UserService) FindUserByEmail(db *gorm.DB, email string) (*system.SysUser, error) {
	var user *system.SysUser
	query := db
	// 邮箱登录/找回密码等公共端点无租户上下文；多租户模式限定默认租户，
	// 避免同一邮箱跨租户命中歧义（阶段一限制：仅默认租户支持邮箱流程）
	if global.Config.Tenant.IsMulti() {
		query = query.Where("tenant_id = ?", global.Config.Tenant.DefaultTenantId)
	}
	error := query.
		Where("email = ?", email).
		First(&user).
		Error
	if error != nil {
		return nil, errors.New("userNotFound")
	}

	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("userDeleted")
	}
	return user, nil
}

// BindEmail 绑定邮箱
func (u *UserService) BindEmail(db *gorm.DB, id string, email string) error {
	if err := db.
		Model(&system.SysUser{}).
		Where("id = ?", id).
		Update("email", email).
		Error; err != nil {
		global.Log.Error("绑定邮箱失败", zap.String("id", id), zap.String("email", email), zap.Error(err))
		return errors.New("bindEmailFailed")
	}
	return nil
}

// UpdatePassword 修改密码
func (u *UserService) UpdatePassword(db *gorm.DB, id string, oldPassword, password, accessToken string) error {
	var sysUserService sysService.SysUserService
	return sysUserService.ResetPassword(
		db,
		id,
		sysReq.UpdateSysUserPasswordReq{
			Password: password,
		},
		&oldPassword,
		accessToken,
	)
}

// UpdateProfile 修改用户基本信息
func (u *UserService) UpdateProfile(db *gorm.DB, id string, req request.UserModifyProfileReq) error {
	var user *system.SysUser
	error := db.
		Where("id = ?", id).
		First(&user).
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
func (u *UserService) GetProfile(db *gorm.DB, id string) (*response.UserProfile, error) {
	var user *system.SysUser
	error := db.
		Preload("Roles").
		Where("id = ?", id).
		First(&user).
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

func (u *UserService) GetMenu(db *gorm.DB, id string) (*[]response.UserMenu, error) {
	var user *system.SysUser

	// 查询用户及角色菜单
	err := db.
		Preload("Roles.Menus").
		Where("id = ?", id).
		First(&user).
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

	// 一次加载全量菜单建内存索引，供父级回填用。
	// 菜单表是小型参考表，这里 1 次 Find 替代原逐层 N 次 First。
	allMenus := make(map[uint]system.SysMenu)
	var all []system.SysMenu
	if err := db.Find(&all).Error; err != nil {
		return nil, errors.New("getMenuListFailed")
	}
	for _, m := range all {
		allMenus[m.Id] = m
	}

	// 构建菜单map，避免重复查询
	menuMap := make(map[uint]system.SysMenu)
	if user.IsSuperAdmin {
		// 平台超级管理员：由 is_super_admin 标识直接授予全部菜单（含平台独有租户管理子树），
		// 不依赖角色绑定，账号改名/换账号不影响权限；Casbin 中间件对平台超管同样直接放行。
		for id, m := range allMenus {
			menuMap[id] = m
		}
	} else {
		for _, menu := range menus {
			menuMap[menu.Id] = menu
		}
	}

	// 递归向上查找父菜单：直接从内存索引取，不再逐级打数据库
	toCheck := menus
	for len(toCheck) > 0 {
		var nextToCheck []system.SysMenu
		for _, menu := range toCheck {
			if menu.ParentId != nil && *menu.ParentId != 0 {
				if _, exists := menuMap[*menu.ParentId]; !exists {
					if parentMenu, ok := allMenus[*menu.ParentId]; ok {
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
			Sort:       menu.Sort,
		})
	}

	// 按照 Sort 升序排序，Sort 相同再按 ID 升序，与菜单树排序规则保持一致
	slices.SortFunc(userMenus, func(a, b response.UserMenu) int {
		if a.Sort != b.Sort {
			if a.Sort < b.Sort {
				return -1
			}
			return 1
		}
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})

	return &userMenus, nil
}
