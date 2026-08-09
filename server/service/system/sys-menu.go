package system

import (
	"errors"
	"slices"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysMenuService struct{}

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
	if *req.ParentId != 0 {
		_, err := s.checkMenuExist(db, *req.ParentId, false)
		if err != nil {
			return nil, err
		}
	}

	menu := &system.SysMenu{
		Name:       req.Name,
		ParentId:   req.ParentId,
		Permission: req.Permission,
		Icon:       req.Icon,
		Type:       req.Type,
		Route:      req.Route,
		Component:  req.Component,
		Sort:       req.Sort,
		Visible:    req.Visible,
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
		ParentId:   menu.ParentId,
		Permission: menu.Permission,
		Icon:       menu.Icon,
		Type:       menu.Type,
		Route:      menu.Route,
		Component:  menu.Component,
		Sort:       menu.Sort,
		Visible:    menu.Visible,
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
func (s *SysMenuService) Update(db *gorm.DB, req request.UpdateSysMenuReq) error {
	var menu *system.SysMenu
	menu, err := s.checkMenuExist(db, req.Id, false)
	if err != nil {
		return err
	}
	if *req.ParentId != 0 {
		_, err = s.checkMenuExist(db, *req.ParentId, true)
		if err != nil {
			return err
		}
	}

	menu.Name = req.Name
	menu.ParentId = req.ParentId
	menu.Permission = req.Permission
	menu.Icon = req.Icon
	menu.Type = req.Type
	menu.Route = req.Route
	menu.Component = req.Component
	menu.Sort = req.Sort
	menu.Visible = req.Visible

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
	_, err := s.checkMenuExist(db, id, false)
	if err != nil {
		return nil, err
	}

	var menu *system.SysMenu
	if err := db.
		Preload("Apis").
		First(&menu, id).
		Error; err != nil {
		return nil, errors.New("getMenuFailed")
	}

	return &response.SysMenuResp{
		Id:         menu.Id,
		Name:       menu.Name,
		ParentId:   menu.ParentId,
		Permission: menu.Permission,
		Icon:       menu.Icon,
		Type:       menu.Type,
		Route:      menu.Route,
		Component:  menu.Component,
		Sort:       menu.Sort,
		Visible:    menu.Visible,
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

// GetTree 获取菜单树
func (s *SysMenuService) GetTree(db *gorm.DB) ([]response.SysMenuTreeResp, error) {
	var menus []system.SysMenu
	if err := db.Find(&menus).Error; err != nil {
		return nil, errors.New("getMenuListFailed")
	}

	// 构建ID到菜单的映射
	menuMap := make(map[uint]*response.SysMenuTreeResp)
	for _, m := range menus {
		menuMap[m.Id] = &response.SysMenuTreeResp{
			Id:         m.Id,
			Name:       m.Name,
			ParentId:   m.ParentId,
			Permission: m.Permission,
			Icon:       m.Icon,
			Type:       m.Type,
			Route:      m.Route,
			Component:  m.Component,
			Sort:       m.Sort,
			Visible:    m.Visible,
			Children:   []response.SysMenuTreeResp{}, // 初始化为空数组而不是nil
		}
	}

	// 构建树结构
	var rootMenus []response.SysMenuTreeResp
	// 使用map记录已处理的节点，避免重复处理
	processed := make(map[uint]bool)

	// 递归构建子树的函数
	var buildSubTree func(uint) []response.SysMenuTreeResp
	buildSubTree = func(parentId uint) []response.SysMenuTreeResp {
		children := make([]response.SysMenuTreeResp, 0)
		for _, m := range menus {
			if m.ParentId != nil && *m.ParentId == parentId && !processed[m.Id] {
				node := *menuMap[m.Id]
				// 递归获取子节点
				node.Children = buildSubTree(m.Id) // 递归调用buildSubTree获取子节点
				// 对子节点排序
				sortMenus(node.Children)
				children = append(children, node)
				processed[m.Id] = true
			}
		}
		// 对当前层级排序
		sortMenus(children)
		return children
	}

	// 从根节点开始构建树
	rootMenus = buildSubTree(0)

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
