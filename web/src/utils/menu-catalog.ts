import { msg } from '@lingui/core/macro'
import type { MessageDescriptor } from '@lingui/core'

/**
 * 内置菜单翻译键的声明处。
 *
 * 菜单是数据库里的行，`i18n_key` 由 `server/migration/menu_seed.go`（迁移版本 7）
 * 回填，运行期通过 `i18n._({ id: key })` 动态取值（见 [getMenuLabel]）。动态 id 不会被
 * `lingui extract` 扫到，因此必须在这里用 `msg()` 显式声明一次，各语言的
 * catalog 才会生成对应条目；否则界面永远回落到数据库里的中文 name。
 *
 * 键的取值规则与 Go 侧一致：
 *   - 页面与目录用 `menu.<路由名>`，与路由一一对应，改中文名不影响翻译；
 *   - 按钮统一用 `menu.action.<动作>`，46 个按钮共用 9 条，动作词表由后端下发。
 *
 * message 是 zh-CN 源文案，与 init.sql 播种的菜单名保持一致。
 */
export const menuMessageDescriptors: MessageDescriptor[] = [
  msg({ id: 'menu.system', message: '系统管理' }),
  msg({ id: 'menu.api', message: '接口管理' }),
  msg({ id: 'menu.menu', message: '菜单管理' }),
  msg({ id: 'menu.role', message: '角色管理' }),
  msg({ id: 'menu.user', message: '用户管理' }),
  msg({ id: 'menu.operation-log', message: '操作日志' }),
  msg({ id: 'menu.token', message: '令牌管理' }),
  msg({ id: 'menu.system-info', message: '系统信息' }),
  msg({ id: 'menu.task', message: '任务管理' }),
  msg({ id: 'menu.tenant', message: '租户管理' }),

  msg({ id: 'menu.action.query', message: '查询' }),
  msg({ id: 'menu.action.info', message: '详情' }),
  msg({ id: 'menu.action.create', message: '创建' }),
  msg({ id: 'menu.action.update', message: '更新' }),
  msg({ id: 'menu.action.delete', message: '删除' }),
  msg({ id: 'menu.action.start', message: '启动' }),
  msg({ id: 'menu.action.stop', message: '停止' }),
  msg({ id: 'menu.action.execute', message: '立即执行' }),
  msg({ id: 'menu.action.reset', message: '重置密码' }),
]

/** 已声明 catalog 条目的菜单翻译键，用于与后端种子对齐的漂移校验。 */
export const knownMenuI18nKeys: ReadonlySet<string> = new Set(
  menuMessageDescriptors.map((item) => item.id).filter((id): id is string =>
    Boolean(id)
  )
)
