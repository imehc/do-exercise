package migration

import "gorm.io/gorm"

func executionStorageSchema(db *gorm.DB) error {
	return execAll(db, []string{
		`INSERT INTO sys_api (path,description,"group",method,disabled,created_at,updated_at) SELECT '/system/jobs/:id/runs','任务执行记录','SYSTEM','GET',FALSE,now(),now() WHERE NOT EXISTS (SELECT 1 FROM sys_api WHERE path = '/system/jobs/:id/runs' AND method = 'GET' AND deleted_at IS NULL)`,
		`INSERT INTO sys_menu_apis (sys_menu_id,sys_api_id) SELECT DISTINCT link.sys_menu_id, target.id FROM sys_menu_apis link JOIN sys_api source ON source.id = link.sys_api_id JOIN sys_api target ON target.path = '/system/jobs/:id/runs' AND target.method = 'GET' AND target.deleted_at IS NULL WHERE source.path = '/system/jobs/:id' AND source.method = 'GET' AND source.deleted_at IS NULL ON CONFLICT DO NOTHING`,
		`INSERT INTO sys_storage_usage (tenant_id,charged_bytes) SELECT tenant_id, COALESCE(SUM(max_size),0) FROM sys_file WHERE state <> 'deleted' GROUP BY tenant_id ON CONFLICT (tenant_id) DO NOTHING`,
		`ALTER TABLE sys_storage_usage ADD CONSTRAINT ck_storage_nonnegative CHECK (charged_bytes >= 0)`,
		`ALTER TABLE sys_tenant ADD CONSTRAINT ck_tenant_storage_nonnegative CHECK (max_storage_bytes >= 0)`,
		`ALTER TABLE sys_file ADD CONSTRAINT ck_file_size CHECK (max_size > 0 AND size >= 0 AND size <= max_size)`,
		`ALTER TABLE sys_file ADD CONSTRAINT ck_file_state CHECK (state IN ('legacy','pending','ready','deleting','deleted'))`,
		`ALTER TABLE sys_job_run ADD CONSTRAINT ck_job_run_state CHECK (status IN ('queued','retrying','succeeded','failed','skipped'))`,
	})
}
