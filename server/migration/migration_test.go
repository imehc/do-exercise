package migration

import (
	"strings"
	"testing"

	"gorm.io/gorm"
)

// noopUp 占位实现，用于只关心注册表结构的用例
func noopUp(*gorm.DB) error { return nil }

func TestValidateAcceptsStrictlyIncreasingVersions(t *testing.T) {
	migrations := []Migration{
		{Version: 1, Name: "first", Up: noopUp},
		{Version: 2, Name: "second", Up: noopUp},
		{Version: 10, Name: "third", Up: noopUp},
	}
	if err := Validate(migrations); err != nil {
		t.Fatalf("合法的迁移清单被拒绝: %v", err)
	}
}

func TestValidateAcceptsEmptyList(t *testing.T) {
	if err := Validate(nil); err != nil {
		t.Fatalf("空清单应当合法: %v", err)
	}
}

// 顺序错乱是最危险的一类编码错误：清单顺序决定执行顺序，
// 而已执行的版本永远不会重放，写错一次就让新库和存量库长得不一样。
func TestValidateRejectsBadRegistries(t *testing.T) {
	cases := []struct {
		name       string
		migrations []Migration
		wantIn     string
	}{
		{
			name:       "缺少版本号",
			migrations: []Migration{{Name: "no-version", Up: noopUp}},
			wantIn:     "缺少版本号",
		},
		{
			name: "版本号回退",
			migrations: []Migration{
				{Version: 2, Name: "second", Up: noopUp},
				{Version: 1, Name: "first", Up: noopUp},
			},
			wantIn: "未严格递增",
		},
		{
			name: "版本号重复",
			migrations: []Migration{
				{Version: 1, Name: "first", Up: noopUp},
				{Version: 1, Name: "duplicate", Up: noopUp},
			},
			wantIn: "未严格递增",
		},
		{
			name:       "名字为空",
			migrations: []Migration{{Version: 1, Name: "   ", Up: noopUp}},
			wantIn:     "缺少名字",
		},
		{
			name: "名字重复",
			migrations: []Migration{
				{Version: 1, Name: "same", Up: noopUp},
				{Version: 2, Name: "same", Up: noopUp},
			},
			wantIn: "同名",
		},
		{
			name:       "没有实现 Up",
			migrations: []Migration{{Version: 1, Name: "empty"}},
			wantIn:     "没有实现 Up",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.migrations)
			if err == nil {
				t.Fatal("期望校验失败，实际通过")
			}
			if !strings.Contains(err.Error(), tc.wantIn) {
				t.Fatalf("错误信息 %q 未包含 %q", err.Error(), tc.wantIn)
			}
		})
	}
}

// All 是真正会作用到生产库的清单，它自己必须先过校验。
// 这条用例的作用是：新增迁移时写错版本号或复制粘贴忘改名字，在这里就红，
// 而不是等到某个环境上跑 migrate 才发现。
func TestAllIsValid(t *testing.T) {
	if err := Validate(All("-- seed")); err != nil {
		t.Fatalf("内置迁移清单不合法: %v", err)
	}
}

// 版本号一旦发布就等于写进了各个环境的 schema_migrations，改动会让存量库永久跳过
// 新内容。这里把当前的版本号与名字钉死，任何改名/改号都必须是有意识的一次修改。
func TestAllVersionsAreStable(t *testing.T) {
	want := []struct {
		version uint64
		name    string
	}{
		{1, "seed_initial_data"},
		{2, "force_default_admin_password_rotation"},
		{3, "sys_user_soft_delete_aware_unique_indexes"},
		{4, "sys_role_code_partial_unique_index"},
		{5, "sys_menu_permission_partial_unique_index"},
		{6, "casbin_rule_sequence_setval"},
		{7, "menu_metadata_backfill"},
		{8, "sys_menu_scope_not_null"},
		{9, "sys_user_tenant_membership"},
	}

	got := All("")
	if len(got) != len(want) {
		t.Fatalf("迁移数量变化：期望 %d 个，实际 %d 个（新增迁移时请同步本用例）", len(want), len(got))
	}
	for i, w := range want {
		if got[i].Version != w.version || got[i].Name != w.name {
			t.Errorf("第 %d 项迁移变成了 %d_%s，期望 %d_%s", i+1, got[i].Version, got[i].Name, w.version, w.name)
		}
	}
}

// 种子播种是唯一依赖外部文件的迁移，其余版本都只依赖库内状态。
// 拿不到种子内容时它必须报错而不是静默放过——否则后续版本会在空表上跑，
// 最终得到一个没有超管账号、看起来却「迁移成功」的库。
// 该分支需要连库（要先查 sys_user 是否已有数据），由集成测试（E-2）覆盖，
// 这里只锁住「版本 1 就是种子播种」这个前提。
func TestSeedMigrationIsFirst(t *testing.T) {
	seed := All("")[0]
	if seed.Version != 1 || seed.Name != "seed_initial_data" {
		t.Fatalf("版本 1 应当是种子播种，实际是 %d_%s", seed.Version, seed.Name)
	}
}
