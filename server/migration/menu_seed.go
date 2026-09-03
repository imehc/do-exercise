package migration

import (
	"fmt"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// 内置菜单的元数据种子。
//
// init.sql 只负责插入菜单本身，i18n_key / scope / is_system 三列统一在这里回填，
// 原因是这三列既要作用于新库（执行完 init.sql 之后），也要作用于存量库
// （升级时 init.sql 不会重跑）。放在 Go 里就只有一份口径，不必在 SQL 和代码之间对齐。
//
// 这份种子随迁移一起走（migrations.go 的版本 7），不再由 migrate 命令直接调用：
// 它改的是数据而不是结构，属于「哪个库跑过、哪个库没跑过」需要被记录的那一类变更。

// menuI18nKeys 内置菜单的稳定翻译键。
//
// 键的取值规则：
//   - 页面与目录用 menu.<路由名>，与路由一一对应，改中文名不影响翻译；
//   - 按钮统一用 menu.action.<动作>：「查询」「创建」在每个模块里是同一个词，
//     没必要为 46 个按钮各造一个键，前端 catalog 也只需维护 9 条。
//
// 只回填 i18n_key IS NULL 的行，运维手工改过的键不覆盖。
var menuI18nKeys = map[uint]string{
	1:  "menu.system",
	2:  "menu.api",
	3:  "menu.menu",
	4:  "menu.role",
	5:  "menu.user",
	6:  "menu.operation-log",
	7:  "menu.token",
	8:  "menu.system-info",
	9:  "menu.task",
	10: "menu.tenant",

	101: "menu.action.query", 102: "menu.action.update",
	111: "menu.action.query", 112: "menu.action.create", 113: "menu.action.update",
	114: "menu.action.delete", 115: "menu.action.info",
	121: "menu.action.query", 122: "menu.action.create", 123: "menu.action.update",
	124: "menu.action.delete", 125: "menu.action.info",
	131: "menu.action.query", 132: "menu.action.create", 133: "menu.action.update",
	134: "menu.action.delete", 135: "menu.action.info", 136: "menu.action.reset",
	141: "menu.action.query", 142: "menu.action.info",
	151: "menu.action.query", 152: "menu.action.delete", 153: "menu.action.update",
	154: "menu.action.info",
	161: "menu.action.query",
	171: "menu.action.query", 172: "menu.action.create", 173: "menu.action.update",
	174: "menu.action.delete", 175: "menu.action.start", 176: "menu.action.stop",
	177: "menu.action.execute", 178: "menu.action.info",
	181: "menu.action.query", 182: "menu.action.create", 183: "menu.action.update",
	184: "menu.action.delete", 185: "menu.action.info",
}

// backfillMenuMetadata 幂等回填内置菜单的 scope / is_system / i18n_key。
func backfillMenuMetadata(db *gorm.DB) error {
	// 平台专属菜单：把旧种子的固定 ID 名单一次性翻译成显式 scope。
	// 之后运行期只认 scope，不再按 ID 猜测（见 service/system/scope.go）。
	if len(global.LegacyPlatformOnlyMenuIDs) > 0 {
		if err := db.Exec(
			"UPDATE sys_menu SET scope = ? WHERE id IN ? AND (scope IS NULL OR scope <> ?)",
			global.MenuScopePlatform, global.LegacyPlatformOnlyMenuIDs, global.MenuScopePlatform,
		).Error; err != nil {
			return err
		}
	}
	// 没标注过的行落到默认可见范围，避免 NULL 参与后续判定
	if err := db.Exec(
		"UPDATE sys_menu SET scope = ? WHERE scope IS NULL OR scope = ''",
		global.MenuScopeBoth,
	).Error; err != nil {
		return err
	}

	// 内置菜单的路由、权限、类型是平台契约，标记 is_system 后只允许改显示属性。
	// 这里覆盖全部内置菜单（不只是平台专属那 8 条），否则「系统内置」这个保护
	// 实际只护住了租户管理子树，改坏 /user 的路由照样能把前端打穿。
	builtinIds := make([]uint, 0, len(menuI18nKeys))
	for id := range menuI18nKeys {
		builtinIds = append(builtinIds, id)
	}
	if err := db.Exec(
		"UPDATE sys_menu SET is_system = TRUE WHERE id IN ? AND is_system = FALSE",
		builtinIds,
	).Error; err != nil {
		return err
	}

	// 翻译键：按 key 分组批量更新，46 行菜单只需 19 条 UPDATE
	idsByKey := make(map[string][]uint, len(menuI18nKeys))
	for id, key := range menuI18nKeys {
		idsByKey[key] = append(idsByKey[key], id)
	}
	var filled int64
	for key, ids := range idsByKey {
		result := db.Exec(
			"UPDATE sys_menu SET i18n_key = ? WHERE id IN ? AND i18n_key IS NULL",
			key, ids,
		)
		if result.Error != nil {
			return result.Error
		}
		filled += result.RowsAffected
	}
	if filled > 0 {
		fmt.Printf("已回填 %d 条菜单翻译键\n", filled)
	}
	return nil
}

// MissingMenuI18nKeys 列出仍未标注翻译键的菜单，提示运维补齐 catalog。
// 自建菜单不在 menuI18nKeys 里，只能靠界面填写；这里只提醒，不阻断，
// 因此它不是一条迁移，而是 migrate 命令收尾时的一次只读检查。
func MissingMenuI18nKeys(db *gorm.DB) string {
	var names []string
	if err := db.Raw(
		"SELECT name FROM sys_menu WHERE deleted_at IS NULL AND (i18n_key IS NULL OR i18n_key = '') ORDER BY id",
	).Scan(&names).Error; err != nil {
		return ""
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, ", ")
}
