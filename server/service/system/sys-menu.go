package system

import (
	"errors"
	"slices"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysMenuService struct{}

// routePermissionKey 把父级菜单路由归一成权限标识前缀，与前端 routePermissionKey 同一算法。
func routePermissionKey(route string) string {
	route = strings.Trim(route, "/")
	return strings.ReplaceAll(route, "/", "_")
}

// validatePermission 校验按钮菜单的权限标识：
// 必须形如 <父级路由归一化>:<允许动作>，动作取自 global.MenuPermissionActions（唯一来源）。
// 权限标识是 Casbin 策略与前端按钮鉴权的共同键，放开自由输入会长出
// user:query / users:query / user:list 这类同义不同名的标识，权限表迅速失控。
func (s *SysMenuService) validatePermission(db *gorm.DB, parentID *uint, menuType uint8, permission *string) error {
	if menuType != 3 {
		return nil
	}
	if parentID == nil || *parentID == 0 || permission == nil {
		return errors.New("invalidPermission")
	}
	parent, err := s.checkMenuExist(db, *parentID, false)
	if err != nil {
		return err
	}
	parts := strings.SplitN(*permission, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[0] != routePermissionKey(parent.Route) {
		return errors.New("invalidPermission")
	}
	if !global.IsMenuPermissionAction(parts[1]) {
		return errors.New("invalidPermission")
	}
	return nil
}

// assignApis 分配API。
// 空 apiIds 的语义是「解除该菜单的全部 API 绑定」，不能提前 return，
// 否则解绑请求返回成功但关联表原样保留。
func (s *SysMenuService) assignApis(tx *gorm.DB, menu *system.SysMenu, apiIds []uint) ([]system.SysApi, error) {
	var apis []system.SysApi
	if len(apiIds) > 0 {
		// 检查API是否存在
		if err := tx.Where("id IN ?", apiIds).Find(&apis).Error; err != nil {
			return nil, errors.New("allApisNotFound")
		}
		if len(apis) != len(apiIds) {
			return nil, errors.New("apiNotFound")
		}
	}
	// 建立/清空api菜单关联
	if len(apis) == 0 {
		if err := tx.Model(menu).Association("Apis").Clear(); err != nil {
			return nil, errors.New("apiAssignFailed")
		}
		return apis, nil
	}
	if err := tx.Model(menu).Association("Apis").Replace(apis); err != nil {
		return nil, errors.New("apiAssignFailed")
	}
	return apis, nil
}

// checkMenuExist 检查菜单是否存在
func (s *SysMenuService) checkMenuExist(db *gorm.DB, menuId uint, isParent bool) (*system.SysMenu, error) {
	var menu system.SysMenu
	result := db.
		Unscoped().
		First(&menu, menuId)
	if result.Error != nil {
		if isParent {
			return nil, errors.New("parentMenuNotFound")
		}
		return nil, errors.New("allMenusNotFound")
	}
	if !menu.DeletedAt.Time.IsZero() {
		return nil, errors.New("menuDeleted")
	}

	return &menu, nil
}

// Create 创建菜单
func (s *SysMenuService) Create(db *gorm.DB, req request.CreateSysMenuReq) (*response.SysMenuResp, error) {
	log := global.Log
	// sys_menu 是全租户共享的定义表，租户改一行会影响所有租户，因此多租户模式下
	// 写操作收归平台超级管理员；单租户模式没有共享问题，行为保持不变。
	if tenantRestricted(db) {
		return nil, errors.New("menuReadonlyForTenant")
	}
	if *req.ParentId != 0 {
		_, err := s.checkMenuExist(db, *req.ParentId, false)
		if err != nil {
			return nil, err
		}
	}
	if err := s.validatePermission(db, req.ParentId, req.Type, req.Permission); err != nil {
		return nil, err
	}
	scope := req.Scope
	if scope == "" {
		scope = global.MenuScopeBoth
	}
	if scope == global.MenuScopePlatform {
		// 与 Update 同一理由：新增平台专属菜单同样是授权边界变更，留一条可追溯的审计
		zap.L().Warn("新增平台专属菜单",
			zap.String("menuName", req.Name),
			zap.String("operator", model.CurrentUserID(db)),
		)
	}

	menu := &system.SysMenu{
		Name:       req.Name,
		I18nKey:    req.I18nKey,
		ParentId:   req.ParentId,
		Permission: req.Permission,
		Icon:       req.Icon,
		Type:       req.Type,
		Route:      req.Route,
		Component:  req.Component,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Scope:      scope,
		IsSystem:   req.IsSystem,
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("Failed to create menu", zap.Any("panic", r))
			if tx != nil {
				_ = tx.Rollback().Error
			}
		}
	}()

	err := tx.Create(menu).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("createMenuFailed")
	}

	apis, err := s.assignApis(tx, menu, req.ApiIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createMenuFailed")
	}

	return &response.SysMenuResp{
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
		Scope:      menu.Scope,
		IsSystem:   menu.IsSystem,
		CreatedAt:  menu.CreatedAt,
		CreatedBy:  menu.CreatedBy,
		UpdatedAt:  menu.UpdatedAt,
		UpdatedBy:  menu.UpdatedBy,
		Apis: lo.Map(apis, func(item system.SysApi, index int) response.SysApiResp {
			return response.SysApiResp{
				Method:      item.Method,
				Description: item.Description,
				Group:       item.Group,
				Disabled:    item.Disabled,
				Sort:        item.Sort,
			}
		}),
	}, nil
}

// Delete 删除菜单
func (s *SysMenuService) Delete(db *gorm.DB, id uint) error {
	// 共享定义表，写操作理由同 Create
	if tenantRestricted(db) {
		return errors.New("menuReadonlyForTenant")
	}
	var menu *system.SysMenu
	menu, err := s.checkMenuExist(db, id, false)
	if err != nil {
		return err
	}

	err = db.
		Delete(menu, id).
		Error
	if err != nil {
		return errors.New("deleteMenuFailed")
	}
	return nil
}

// Update 更新菜单
func (s *SysMenuService) Update(db *gorm.DB, id uint, req request.UpdateSysMenuReq) error {
	// 共享定义表，写操作理由同 Create
	if tenantRestricted(db) {
		return errors.New("menuReadonlyForTenant")
	}
	var menu *system.SysMenu
	menu, err := s.checkMenuExist(db, id, false)
	if err != nil {
		return err
	}
	if menu.IsSystem {
		// 系统菜单的路由、权限和类型是平台契约，只允许改显示属性。
		// 这一步必须排在校验之前：被强制回填的字段不该再拿请求里的值去校验，
		// 否则编辑一个内置按钮时只改排序也会因为没提交 permission 而被判成非法权限标识。
		req.Name = menu.Name
		// 未提交翻译键时保留已有值；显式提供新键可用于补齐系统菜单 catalog。
		if req.I18nKey == nil {
			req.I18nKey = menu.I18nKey
		}
		req.ParentId = menu.ParentId
		req.Permission = menu.Permission
		req.Type = menu.Type
		req.Route = menu.Route
		req.Component = menu.Component
		req.Scope = menu.Scope
	}
	if req.ParentId != nil && *req.ParentId != 0 {
		_, err = s.checkMenuExist(db, *req.ParentId, true)
		if err != nil {
			return err
		}
	}
	if err := s.validatePermission(db, req.ParentId, req.Type, req.Permission); err != nil {
		return err
	}

	menu.Name = req.Name
	menu.I18nKey = req.I18nKey
	menu.ParentId = req.ParentId
	menu.Permission = req.Permission
	menu.Icon = req.Icon
	menu.Type = req.Type
	menu.Route = req.Route
	menu.Component = req.Component
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	if req.Scope != "" && req.Scope != menu.Scope {
		// scope 决定这条菜单能不能落到业务租户手里，是一条授权边界的变更。
		// 操作记录只留下请求体，这里额外打一条带前后值的审计日志，便于事后追溯
		// 「某个平台菜单是什么时候被放开给租户的」。
		zap.L().Warn("菜单可见范围变更",
			zap.Uint("menuId", menu.Id),
			zap.String("menuName", menu.Name),
			zap.String("from", menu.Scope),
			zap.String("to", req.Scope),
			zap.String("operator", model.CurrentUserID(db)),
		)
		menu.Scope = req.Scope
	}
	menu.IsSystem = menu.IsSystem || req.IsSystem

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Omit("id", "created_at", "created_by").
		Updates(&menu).
		Error; err != nil {
		tx.Rollback()
		return errors.New("updateMenuFailed")
	}

	if _, err := s.assignApis(tx, menu, req.ApiIds); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("updateMenuFailed")
	}
	return nil
}

// Get 查询单个菜单
func (s *SysMenuService) Get(db *gorm.DB, id uint) (*response.SysMenuResp, error) {
	// 一次查询带出 Apis，不再先查存在性再重查同一行
	var menu system.SysMenu
	result := scopeTenantVisibleMenus(db).
		Unscoped().
		Preload("Apis").
		First(&menu, id)
	if result.Error != nil {
		return nil, errors.New("allMenusNotFound")
	}
	if !menu.DeletedAt.Time.IsZero() {
		return nil, errors.New("menuDeleted")
	}

	return &response.SysMenuResp{
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
		Scope:      menu.Scope,
		IsSystem:   menu.IsSystem,
		CreatedAt:  menu.CreatedAt,
		CreatedBy:  menu.CreatedBy,
		UpdatedAt:  menu.UpdatedAt,
		UpdatedBy:  menu.UpdatedBy,
		Apis: lo.Map(menu.Apis, func(item system.SysApi, index int) response.SysApiResp {
			return response.SysApiResp{
				Id:          item.Id,
				Path:        item.Path,
				Method:      item.Method,
				Description: item.Description,
				Group:       item.Group,
				Disabled:    item.Disabled,
				Sort:        item.Sort,
			}
		}),
	}, nil
}

// toApiBriefs 把菜单绑定的接口压成授权页预览需要的最小字段集
func toApiBriefs(apis []system.SysApi) []response.SysMenuApiBrief {
	return lo.Map(apis, func(item system.SysApi, _ int) response.SysMenuApiBrief {
		return response.SysMenuApiBrief{
			Id:          item.Id,
			Method:      item.Method,
			Path:        item.Path,
			Description: item.Description,
		}
	})
}

// GetTree 获取菜单树（按调用者的租户可见范围裁剪）
func (s *SysMenuService) GetTree(db *gorm.DB) ([]response.SysMenuTreeResp, error) {
	var menus []system.SysMenu
	// 预加载绑定的 API：角色授权页需要在勾选权限时就地预览这条权限实际放开哪些接口，
	// 否则只能逐个菜单再查一次详情（还要求授权人额外具备 menu:info）。
	if err := scopeTenantVisibleMenus(db).Preload("Apis").Find(&menus).Error; err != nil {
		return nil, errors.New("getMenuListFailed")
	}

	// 构建ID到菜单的映射
	menuMap := make(map[uint]*response.SysMenuTreeResp)
	for _, m := range menus {
		menuMap[m.Id] = &response.SysMenuTreeResp{
			Id:         m.Id,
			Name:       m.Name,
			I18nKey:    m.I18nKey,
			ParentId:   m.ParentId,
			Permission: m.Permission,
			Icon:       m.Icon,
			Type:       m.Type,
			Route:      m.Route,
			Component:  m.Component,
			Sort:       m.Sort,
			Visible:    m.Visible,
			Scope:      m.Scope,
			IsSystem:   m.IsSystem,
			Apis:       toApiBriefs(m.Apis),
			Children:   []response.SysMenuTreeResp{}, // 初始化为空数组而不是nil
		}
	}

	// 构建树结构：先按父节点分组，避免每层重扫整个切片（O(N)）
	childrenOf := make(map[uint][]response.SysMenuTreeResp)
	for _, m := range menus {
		if m.ParentId != nil && *m.ParentId != 0 {
			childrenOf[*m.ParentId] = append(childrenOf[*m.ParentId], *menuMap[m.Id])
		} else {
			childrenOf[0] = append(childrenOf[0], *menuMap[m.Id])
		}
	}
	for parentId := range childrenOf {
		sortMenus(childrenOf[parentId])
	}

	// 递归构建子树
	var buildSubTree func(uint) []response.SysMenuTreeResp
	buildSubTree = func(parentId uint) []response.SysMenuTreeResp {
		children := childrenOf[parentId]
		// 叶子节点在 map 里取不到值，直接返回会是 nil，序列化成 null；
		// 前端生成的解析器对 children 递归 map，null 会直接抛错，所以补成空数组。
		if children == nil {
			return []response.SysMenuTreeResp{}
		}
		for i := range children {
			children[i].Children = buildSubTree(children[i].Id)
		}
		return children
	}

	rootMenus := buildSubTree(0)
	return rootMenus, nil
}

// sortMenus 按Sort字段排序菜单
func sortMenus(menus []response.SysMenuTreeResp) {
	slices.SortFunc(menus, func(a, b response.SysMenuTreeResp) int {
		if a.Sort < b.Sort {
			return -1
		}
		if a.Sort > b.Sort {
			return 1
		}

		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}

		return 0
	})
}
