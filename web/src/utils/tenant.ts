import { t } from '@lingui/core/macro'

/** 平台保留租户ID，与后端 global.PlatformTenantID 保持一致 */
export const PLATFORM_TENANT_ID = 'platform'

/**
 * 租户管理员角色编码，与后端 systemService.TenantAdminRoleCode 保持一致。
 * 该角色由平台在创建租户时供应，只有平台超级管理员可以修改或删除。
 */
export const TENANT_ADMIN_ROLE_CODE = 'tenant_admin'

/**
 * 租户展示文案：平台保留租户显示为「平台」，其余优先展示租户名称。
 * 名称缺失（租户已删除或历史脏数据）时回退展示租户ID，便于排查；两者都为空时展示 '-'。
 */
export function formatTenant(tenantId?: string, tenantName?: string): string {
  if (tenantId === PLATFORM_TENANT_ID) return t`平台`
  return tenantName || tenantId || '-'
}
