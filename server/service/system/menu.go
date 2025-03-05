package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/pkg/utils/scope"
	"github.com/imehc/do-exercise/server/utils"
	"gorm.io/gorm"
)

type MenuService struct{}

// 创建菜单
func (m *MenuService) CreateMenu(request request.MenuRequest, createdBy uint) (err error) {
	db := global.DB

	err = m.checkParentMenuExist(*request.ParentId)
	if err != nil {
		return
	}

	menu := system.Menu{
		ParentId:   *request.ParentId,
		Name:       request.Name,
		Icon:       request.Icon,
		Type:       request.Type,
		Action:     request.Action,
		IsFrame:    request.IsFrame,
		Visible:    request.Visible,
		Title:      request.Title,
		Component:  request.Component,
		Path:       request.Path,
		Permission: request.Permission,
		ControlWrapper: model.ControlWrapper{
			CreatedBy: createdBy,
		},
	}

	if request.Sort != 0 {
		menu.Sort = request.Sort
	}

	tx := db.Begin()
	// 创建菜单
	err = tx.Create(&menu).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if menu.ParentId == 0 {
		menu.Path = utils.FormatFullpath(uint(menu.ParentId), menu.ID, "")
	} else {
		var parentMenu system.Menu
		err = tx.First(&parentMenu, *request.ParentId).Error
		if err != nil {
			return err
		}
		menu.Path = utils.FormatFullpath(uint(menu.ParentId), menu.ID, parentMenu.Path)
	}

	// 更新菜单的Path字段
	if err = tx.Model(&menu).Update("path", menu.Path).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果有关联的API，创建关联关系
	if len(request.ApiIds) > 0 {
		var apis []system.Api
		if err = tx.Where("id IN ?", request.ApiIds).Find(&apis).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Model(&menu).Association("Apis").Replace(apis); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

// 删除菜单
func (m *MenuService) DeleteMenu(param request.MenuParam, deletedBy uint) (err error) {
	db := global.DB

	var menu system.Menu
	result := db.
		Unscoped().
		First(&menu, param.MenuId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
	}

	if !menu.DeletedAt.Time.IsZero() {
		return errors.New("菜单已删除")
	}

	// 检查是否有子菜单
	var childMenu system.Menu
	result = db.
		Model(system.Menu{}).
		Where("parent_id = ?", param.MenuId).
		First(&childMenu)
	if result.Error == nil {
		return errors.New("该菜单下存在子菜单，无法删除")
	}

	// 删除菜单与API的关联关系
	if err = db.Model(&menu).Association("Apis").Clear(); err != nil {
		return err
	}

	db.
		Model(system.Menu{}).
		Where("id = ?", param.MenuId).
		Update("deleted_by", deletedBy).
		Delete(&menu)
	return nil
}

// 更新菜单
func (m *MenuService) UpdateMenu(param request.MenuParam, request request.MenuRequest, updatedBy uint) (err error) {
	db := global.DB

	var menu system.Menu
	result := db.
		First(&menu, param.MenuId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
	}

	err = m.checkParentMenuExist(*request.ParentId)
	if err != nil {
		return err
	}

	menu.ParentId = *request.ParentId
	menu.Name = request.Name
	menu.Icon = request.Icon
	menu.Type = request.Type
	menu.Action = request.Action
	menu.IsFrame = request.IsFrame
	menu.Visible = request.Visible
	menu.Title = request.Title
	menu.Component = request.Component
	menu.Path = request.Path
	menu.Permission = request.Permission
	menu.Sort = request.Sort
	menu.ControlWrapper = model.ControlWrapper{
		UpdatedBy: updatedBy,
	}

	if menu.ParentId == 0 {
		menu.Path = utils.FormatFullpath(uint(menu.ParentId), menu.ID, "")
	} else {
		var parentMenu system.Menu
		err = db.First(&parentMenu, *request.ParentId).Error
		if err != nil {
			return err
		}
		menu.Path = utils.FormatFullpath(uint(menu.ParentId), menu.ID, parentMenu.Path)
	}

	// 更新菜单基本信息
	db.
		Model(system.Menu{}).
		Where("id = ?", param.MenuId).
		Updates(&menu).
		Omit("id", "created_at")

	// 更新菜单与API的关联关系
	if len(request.ApiIds) > 0 {
		var apis []system.Api
		if err = db.Where("id IN ?", request.ApiIds).Find(&apis).Error; err != nil {
			return err
		}
		if err = db.Model(&menu).Association("Apis").Replace(apis); err != nil {
			return err
		}
	} else {
		// 如果没有关联的API，清除所有关联关系
		if err = db.Model(&menu).Association("Apis").Clear(); err != nil {
			return err
		}
	}

	return nil
}

// 获取菜单
func (m *MenuService) GetMenu(param request.MenuParam) (response sysRes.MenuItem, err error) {
	if param.MenuId == 0 {
		return response, errors.New("菜单ID不能为0")
	}

	db := global.DB

	// 查询菜单基本信息
	var menu system.Menu
	result := db.Preload("Apis").First(&menu, param.MenuId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("菜单不存在")
		}
		return response, result.Error
	}

	// 转换为响应结构
	response = sysRes.MenuItem{
		IDWrapper:      menu.IDWrapper,
		ControlWrapper: menu.ControlWrapper,
		MenuRequest: request.MenuRequest{
			ParentId:   &menu.ParentId,
			Name:       menu.Name,
			Icon:       menu.Icon,
			Type:       menu.Type,
			Action:     menu.Action,
			IsFrame:    menu.IsFrame,
			Visible:    menu.Visible,
			Title:      menu.Title,
			Component:  menu.Component,
			Path:       menu.Path,
			Permission: menu.Permission,
			SortWrapper: model.SortWrapper{
				Sort: menu.Sort,
			},
		},
		NoCache: menu.NoCache,
		Params:  menu.Params,
		Route:   menu.Route,
		Apis:    menu.Apis,
	}

	// 递归查询子菜单
	response.Children, err = m.getChildrenMenus(menu.ID)
	if err != nil {
		return
	}

	return
}

// 获取菜单列表
func (m *MenuService) GetMenuTreeList(s common.ScopeData) (menus []sysRes.MenuItem, err error) {
	// TODO: 支持模糊查询

	db := global.DB
	// 应用数据权限过滤
	db = scope.GetDataScope(db, &s, "sys_menu")

	// 查询所有根菜单（ParentId为0的菜单）
	var rootMenus []system.Menu
	if err = db.Preload("Apis").Where("parent_id = ?", 0).Find(&rootMenus).Error; err != nil {
		return nil, err
	}

	// 初始化空数组
	menus = make([]sysRes.MenuItem, 0)

	// 遍历根菜单，构建菜单树
	for _, menu := range rootMenus {
		item := sysRes.MenuItem{
			IDWrapper:      menu.IDWrapper,
			ControlWrapper: menu.ControlWrapper,
			MenuRequest: request.MenuRequest{
				ParentId:   &menu.ParentId,
				Name:       menu.Name,
				Icon:       menu.Icon,
				Type:       menu.Type,
				Action:     menu.Action,
				IsFrame:    menu.IsFrame,
				Visible:    menu.Visible,
				Title:      menu.Title,
				Component:  menu.Component,
				Path:       menu.Path,
				Permission: menu.Permission,
				SortWrapper: model.SortWrapper{
					Sort: menu.Sort,
				},
			},
			NoCache: menu.NoCache,
			Params:  menu.Params,
			Route:   menu.Route,
			Apis:    menu.Apis,
		}

		// 使用getChildrenMenus获取子菜单
		item.Children, err = m.getChildrenMenus(menu.ID)
		if err != nil {
			return nil, err
		}

		menus = append(menus, item)
	}

	return menus, nil
}

// 获取菜单树
func (m *MenuService) GetMenuTree() (response []sysRes.MenuTree, err error) {
	// 初始化空数组
	response = make([]sysRes.MenuTree, 0)

	// 查询顶级菜单
	db := global.DB
	var menus []system.Menu
	if err = db.Where("parent_id = ?", 0).Order("sort DESC").Order("id ASC").Find(&menus).Error; err != nil {
		return
	}

	// 转换为树形结构
	for _, menu := range menus {
		treeNode := sysRes.MenuTree{
			ID:    int(menu.ID),
			Label: menu.Title,
		}

		// 递归获取子菜单
		treeNode.Children, err = m.getMenuTreeChildren(menu.ID)
		if err != nil {
			return nil, err
		}

		response = append(response, treeNode)
	}

	return
}

// 递归获取子菜单树
func (m *MenuService) getMenuTreeChildren(parentId uint) (children []sysRes.MenuTree, err error) {
	// 初始化空数组
	children = make([]sysRes.MenuTree, 0)

	// 查询子菜单
	db := global.DB
	var menus []system.Menu
	if err = db.Where("parent_id = ?", parentId).Order("sort DESC").Order("id ASC").Find(&menus).Error; err != nil {
		return
	}

	// 转换为树形结构
	for _, menu := range menus {
		treeNode := sysRes.MenuTree{
			ID:    int(menu.ID),
			Label: menu.Title,
		}

		// 递归获取子菜单
		treeNode.Children, err = m.getMenuTreeChildren(menu.ID)
		if err != nil {
			return nil, err
		}

		children = append(children, treeNode)
	}

	return
}

// 递归获取子菜单
func (m *MenuService) getChildrenMenus(parentId uint) (children []sysRes.MenuItem, err error) {
	db := global.DB

	// 查询子菜单，使用预加载和排序
	var menus []system.Menu
	if err = db.Preload("Apis").
		Where("parent_id = ?", parentId).
		Order("sort DESC").
		Order("id ASC").
		Find(&menus).Error; err != nil {
		// 确保返回空数组而不是nil
		children = make([]sysRes.MenuItem, 0)
		return
	}

	// 初始化空数组
	children = make([]sysRes.MenuItem, 0)

	// 转换子菜单
	for _, menu := range menus {
		childItem := sysRes.MenuItem{
			IDWrapper:      menu.IDWrapper,
			ControlWrapper: menu.ControlWrapper,
			MenuRequest: request.MenuRequest{
				ParentId:   &menu.ParentId,
				Name:       menu.Name,
				Icon:       menu.Icon,
				Type:       menu.Type,
				Action:     menu.Action,
				IsFrame:    menu.IsFrame,
				Visible:    menu.Visible,
				Title:      menu.Title,
				Component:  menu.Component,
				Path:       menu.Path,
				Permission: menu.Permission,
				SortWrapper: model.SortWrapper{
					Sort: menu.Sort,
				},
			},
			NoCache: menu.NoCache,
			Params:  menu.Params,
			Route:   menu.Route,
			Apis:    menu.Apis,
			// 初始化空的子菜单数组
			Children: make([]sysRes.MenuItem, 0),
		}

		// 递归查询子菜单
		childItem.Children, err = m.getChildrenMenus(menu.ID)
		if err != nil {
			return children, err
		}

		children = append(children, childItem)
	}

	return children, nil
}

// 检查上级菜单是否存在
func (m *MenuService) checkParentMenuExist(parentId uint) (err error) {
	if parentId == 0 {
		return nil
	}
	db := global.DB
	var menu system.Menu
	result := db.
		First(&menu, parentId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("上级菜单不存在")
		}
		return result.Error
	}
	return nil
}
