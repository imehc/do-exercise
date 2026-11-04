package system

import (
	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysMenu struct {
	model.IdWrapper

	Name    string  `json:"name" gorm:"not null;comment:菜单名称"`
	I18nKey *string `json:"i18n_key" gorm:"size:128;default:null;comment:菜单国际化键"`
	// Permission 权限标识全局唯一（菜单由平台统一维护，不分租户），但只在未软删的行之间唯一：
	// 普通唯一索引会让软删掉的按钮永久占用标识，重建同标识只能得到一个通用的写库错误。
	Permission *string `json:"permission" gorm:"size:64;default:null;uniqueIndex:idx_sys_menu_permission,where:deleted_at IS NULL;comment:权限标识"`
	ParentId   *uint   `json:"parent_id" gorm:"default:null;comment:父菜单ID"`
	Icon       string  `json:"icon" gorm:"size:64;default:null;comment:菜单图标"`
	Type       uint8   `json:"type" gorm:"not null;comment:菜单类型(1:目录,2:菜单,3:按钮)"`
	Route      string  `json:"route" gorm:"size:128;default:null;comment:菜单路由"`
	// Component 保留字段：前端是文件式路由，运行期不消费该值。
	// 存量数据里写的是 `/xxx/page.tsx` 这类历史值，界面已不再展示与编辑。
	Component string   `json:"component" gorm:"size:128;default:null;comment:组件地址(保留字段，前端不消费)"`
	Sort      uint     `json:"sort" gorm:"default:0;comment:显示顺序"`
	Visible   bool     `json:"visible" gorm:"default:false;comment:是否显示"`
	Scope     string   `json:"scope" gorm:"size:16;not null;default:both;comment:菜单范围(platform/tenant/both)"`
	IsSystem  bool     `json:"is_system" gorm:"default:false;comment:是否系统内置菜单"`
	Apis      []SysApi `json:"apis" gorm:"many2many:sys_menu_apis;"`

	model.ControlWrapper
}

func (m *SysMenu) BeforeCreate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	m.CreatedBy = userId
	m.UpdatedBy = userId

	return nil
}

func (m *SysMenu) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	if userId == "" {
		return nil
	}

	if m.UpdatedBy != userId && m.Id != 0 {
		m.UpdatedBy = userId
		err = tx.
			Model(m).
			Select("UpdatedBy").
			Updates(m).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SysMenu) BeforeDelete(tx *gorm.DB) (err error) {
	if m.Id != 0 {
		userId := model.CurrentUserID(tx)
		m.DeletedBy = userId
		err = tx.
			Model(m).
			Select("DeletedBy").
			Updates(m).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
