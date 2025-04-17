package system

import (
	"errors"
	"slices"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"gorm.io/gorm"
)

type SysMenuService struct{}

func (s *SysMenuService) Create(req request.CreateSysMenuReq) (*response.SysMenuResp, error) {
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
	err := global.DB.Create(menu).Error
	if err != nil {
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
	}, nil
}

func (s *SysMenuService) Delete(id uint) error {
	var menu system.SysMenu
	result := global.DB.
		Unscoped().
		First(&menu, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
	}

	if !menu.DeletedAt.Time.IsZero() {
		return errors.New("菜单已删除")
	}

	return global.DB.
		Model(&system.SysMenu{}).
		Where("id = ?", id).
		Delete(&menu).
		Error
}

func (s *SysMenuService) Update(req request.UpdateSysMenuReq) error {
	var menu system.SysMenu
	result := global.DB.
		First(&menu, req.Id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
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
	return global.DB.
		Model(system.SysMenu{}).
		Where("id = ?", req.Id).
		Updates(&menu).
		Omit("id", "created_at", "created_by").
		Error
}

func (s *SysMenuService) Get(id uint) (*response.SysMenuResp, error) {
	var menu system.SysMenu
	result := global.DB.
		First(&menu, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &response.SysMenuResp{}, errors.New("菜单不存在")
		}
		return &response.SysMenuResp{}, result.Error
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
