package system

import (
	"errors"
	"slices"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type SysMenuService struct{}

// assignApis 分配API
func (s *SysMenuService) assignApis(tx *gorm.DB, menu *system.SysMenu, apiIds []uint) ([]system.SysApi, error) {
	if len(apiIds) == 0 {
		return []system.SysApi{}, nil
	}
	var apis []system.SysApi
	// 检查菜单是否存在
	if err := tx.Where("id IN ?", apiIds).Find(&apis).Error; err != nil {
		return nil, err
	}
	if len(apis) != len(apiIds) {
		return nil, errors.New("部分Api不存在")
	}
	// 建立角色菜单关联
	if err := tx.Model(menu).Association("Apis").Replace(apis); err != nil {
		return nil, err
	}
	return apis, nil
}

// 检查菜单是否存在
func (s *SysMenuService) checkMenuExist(db *gorm.DB, menuId uint, isParent bool) (*system.SysMenu, error) {
	var menu system.SysMenu
	result := db.
		Unscoped().
		First(&menu, menuId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if isParent {
				return nil, errors.New("父菜单不存在")
			}
			return nil, errors.New("菜单不存在")
		}
		return nil, result.Error
	}
	if !menu.DeletedAt.Time.IsZero() {
		return nil, errors.New("菜单已删除")
	}

	return &menu, nil
}

func (s *SysMenuService) Create(req request.CreateSysMenuReq) (*response.SysMenuResp, error) {
	db := global.DB
	_, err := s.checkMenuExist(db, *req.ParentId, false)
	if err != nil {
		return nil, err
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
			tx.Rollback()
		}
	}()

	err = tx.Create(menu).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	apis, err := s.assignApis(tx, menu, req.ApiIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
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

func (s *SysMenuService) Delete(id uint) error {
	db := global.DB
	var menu system.SysMenu
	_, err := s.checkMenuExist(db, id, false)
	if err != nil {
		return err
	}

	return db.
		Model(&system.SysMenu{}).
		Where("id = ?", id).
		Delete(&menu).
		Error
}

func (s *SysMenuService) Update(req request.UpdateSysMenuReq) error {
	db := global.DB

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
		Model(system.SysMenu{}).
		Where("id = ?", req.Id).
		Updates(&menu).
		Omit("id", "created_at", "created_by").
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if _, err := s.assignApis(tx, menu, req.ApiIds); err != nil {
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

func (s *SysMenuService) Get(id uint) (*response.SysMenuResp, error) {
	db := global.DB
	var menu *system.SysMenu
	menu, err := s.checkMenuExist(db, id, false)
	if err != nil {
		return nil, err
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
	}, nil
}

func (s *SysMenuService) GetTree() ([]response.SysMenuTreeResp, error) {
	var menus []system.SysMenu
	if err := global.DB.Find(&menus).Error; err != nil {
		return nil, err
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
	for _, m := range menus {
		node := menuMap[m.Id]
		if *m.ParentId == 0 {
			// parentId为0的是根节点
			rootMenus = append(rootMenus, *node)
		} else if parent, ok := menuMap[*m.ParentId]; ok {
			// 将子节点添加到父节点的children中
			parent.Children = append(parent.Children, *node)
		}
	}

	// 对根菜单按Sort排序
	sortMenus(rootMenus)

	// 对每个节点的子菜单进行排序
	for _, menu := range menuMap {
		sortMenus(menu.Children)
	}

	return rootMenus, nil
}

// 按Sort字段排序菜单
func sortMenus(menus []response.SysMenuTreeResp) {
	slices.SortFunc(menus, func(a, b response.SysMenuTreeResp) int {
		if a.Sort < b.Sort {
			return -1
		}
		if a.Sort > b.Sort {
			return 1
		}
		return 0
	})
}
