// Package migration 提供版本化的数据库迁移。
//
// 为什么不能只靠 AutoMigrate：它能加列、能建新索引，但改不了既有唯一约束、
// 迁不了数据，也没有「这个库跑到哪一步」的记录。在此之前这些补丁以裸 db.Exec 的
// 形式内联在 migrate 命令里，靠 IF EXISTS / IF NOT EXISTS 保证可重复执行——
// 正确性完全依赖作者每次都记得这件事，而且 Exec 的 error 没人接，失败是静默的。
//
// 这里的口径：
//   - 每次变更是一个显式版本号 + 名字，执行后登记进 schema_migrations 表；
//   - 每个版本在自己的事务里执行，变更与登记同生共死，不存在「跑了一半没记上」；
//   - 整轮迁移由 advisory lock 串行化，两个进程同时迁移不会互相踩；
//   - 迁移里的错误一律向上抛。存量脏数据（例如同租户重复的角色编码）会让迁移
//     大声失败，而不是静默跳过——脚本无法替运维决定该删哪一行。
//
// 迁移只能追加，不能改已发布版本的内容：已经跑过的库不会再执行它，改了只会让
// 新库和存量库长得不一样。需要修正就再加一个版本。
package migration

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Migration 一次不可再分的 schema / 数据变更。
//
// Up 收到的是所在事务的 *gorm.DB，已带租户旁路上下文，可以直接写裸 SQL 或
// 走模型方法。返回 error 即回滚该版本并中断整轮迁移。
type Migration struct {
	// Version 单调递增的版本号，登记进 schema_migrations 后不可再改。
	Version uint64
	// Name 供人阅读的变更名，写进记录表，用于 `migrate --status` 输出。
	Name string
	// Up 变更本身。应当尽量幂等：事务能回滚 DDL 和 DML，但回滚不了序列操作。
	Up func(db *gorm.DB) error
}

// Validate 校验迁移注册表本身的合法性：版本号非零且严格递增、名字非空且不重复、
// Up 非 nil。这些都是编码错误，宁可在执行任何 SQL 之前就失败，
// 也不要用一个顺序错乱的清单去动生产库。注册表的单测直接锁这个函数。
func Validate(migrations []Migration) error {
	var prev uint64
	names := make(map[string]uint64, len(migrations))
	for i, m := range migrations {
		if m.Version == 0 {
			return fmt.Errorf("第 %d 项迁移缺少版本号", i+1)
		}
		if m.Version <= prev {
			return fmt.Errorf("迁移 %d(%s) 的版本号未严格递增，前一项是 %d", m.Version, m.Name, prev)
		}
		name := strings.TrimSpace(m.Name)
		if name == "" {
			return fmt.Errorf("迁移 %d 缺少名字", m.Version)
		}
		if dup, ok := names[name]; ok {
			return fmt.Errorf("迁移 %d 与 %d 同名（%s）", m.Version, dup, name)
		}
		if m.Up == nil {
			return fmt.Errorf("迁移 %d(%s) 没有实现 Up", m.Version, name)
		}
		names[name] = m.Version
		prev = m.Version
	}
	return nil
}
