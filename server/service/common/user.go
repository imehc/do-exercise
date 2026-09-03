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

// FindUsersByEmail 按邮箱跨租户查出全部可用的业务账号（含角色，供签发 token 用）。
//
// 邮箱在多租户下不唯一：同一个人在 N 个租户里就是 N 行 sys_user（唯一索引是
// email+tenant_id）。邮箱登录与找回密码都是匿名公共端点，请求里没有租户上下文，
// 所以这里不能自己挑一个账号返回——早期实现把范围硬编码在「默认租户」上，
// 其他租户的用户根本用不了邮箱流程。改为返回候选集合，由调用方结合租户选择收敛。
//
// 三类账号被排除：
//   - 软删账号：GORM 默认作用域已处理；
//   - 已停用租户的账号：验证码只证明「请求者拥有这个邮箱」，不该成为绕过停用的入口；
//   - 平台租户账号：平台超管是全系统最高权限，若邮箱可以免口令签发平台 token 或
//     重置平台口令，则邮箱被盗即等于整站被接管。平台账号的口令走后台用户管理重置。
func (u *UserService) FindUsersByEmail(db *gorm.DB, email string) ([]system.SysUser, error) {
	var users []system.SysUser
	if err := db.
		Preload("Roles").
		Where("email = ?", email).
		Where("tenant_id <> ?", global.PlatformTenantID).
		Order("created_at ASC").
		Find(&users).Error; err != nil {
		return nil, errors.New("userNotFound")
	}
	if len(users) == 0 {
		return nil, errors.New("userNotFound")
	}

	tenantIds := lo.Uniq(lo.Map(users, func(user system.SysUser, _ int) string { return user.TenantId }))
	var tenants []system.SysTenant
	if err := db.Where("tenant_id IN ? AND status = ?", tenantIds, true).
		Find(&tenants).Error; err != nil {
		return nil, errors.New("userNotFound")
	}
	enabled := lo.SliceToMap(tenants, func(t system.SysTenant) (string, struct{}) {
		return t.TenantId, struct{}{}
	})

	users = lo.Filter(users, func(user system.SysUser, _ int) bool {
		_, ok := enabled[user.TenantId]
		return ok
	})
	if len(users) == 0 {
		return nil, errors.New("userNotFound")
	}
	return users, nil
}

// UserIdsOf 提取候选账号的 ID 集合，用于绑定/校验邮箱验证码。
func (u *UserService) UserIdsOf(users []system.SysUser) []string {
	return lo.Map(users, func(user system.SysUser, _ int) string { return user.UserId })
}

// TenantIdsOf 提取候选账号所属的租户集合，用于生成租户选择列表。
func (u *UserService) TenantIdsOf(users []system.SysUser) []string {
	return lo.Uniq(lo.Map(users, func(user system.SysUser, _ int) string { return user.TenantId }))
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
	menus = sysService.FilterTenantVisibleMenus(db, menus)

	// 一次加载全量菜单建内存索引，供父级回填用。
	// 菜单表是小型参考表，这里 1 次 Find 替代原逐层 N 次 First。
	allMenus := make(map[uint]system.SysMenu)
	var all []system.SysMenu
	if err := db.Find(&all).Error; err != nil {
		return nil, errors.New("getMenuListFailed")
	}
	all = sysService.FilterTenantVisibleMenus(db, all)
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
			I18nKey:    menu.I18nKey,
			ParentId:   menu.ParentId,
			Permission: menu.Permission,
			Icon:       menu.Icon,
			Type:       menu.Type,
			Route:      menu.Route,
			Component:  menu.Component,
			Sort:       menu.Sort,
			Visible:    menu.Visible,
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
