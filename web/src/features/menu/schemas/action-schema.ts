import { z } from 'zod'
import { t } from '@lingui/core/macro'
import { MenuType } from '~/do-exercise-api'

/**
 * 权限动作词表的唯一来源是后端（`global.MenuPermissionActions`，经
 * `GET /auth/tenants` 的 `permission_actions` 下发）。这里保留一份同名常量，
 * 只作为词表尚未到达时的兜底：表单在首屏就要能渲染出下拉项，而 tenants
 * 查询是异步的。两侧不一致时后端校验（invalidPermission）仍是最终裁判。
 */
export const fallbackMenuPermissionActions = [
  'query',
  'info',
  'create',
  'update',
  'delete',
  'start',
  'stop',
  'execute',
  'reset',
] as const

export type MenuPermissionAction = string

/**
 * 动作的中文/英文展示名。后端新增动作时这里没有对应文案，直接回落到原始
 * 动作串，保证下拉项仍然可选，不会因为缺一条翻译就挡住配置。
 */
export function getMenuPermissionActionLabel(action: string): string {
  switch (action) {
    case 'query':
      return t`查询`
    case 'info':
      return t`详情`
    case 'create':
      return t`创建`
    case 'update':
      return t`更新`
    case 'delete':
      return t`删除`
    case 'start':
      return t`启动`
    case 'stop':
      return t`停止`
    case 'execute':
      return t`立即执行`
    case 'reset':
      return t`重置密码`
    default:
      return action
  }
}

export const menuScopes = ['both', 'tenant', 'platform'] as const

export type MenuScope = (typeof menuScopes)[number]

export function getMenuScopeLabel(scope: MenuScope): string {
  switch (scope) {
    case 'platform':
      return t`仅平台`
    case 'tenant':
      return t`仅租户`
    case 'both':
      return t`平台与租户`
  }
}

const getBasicSchema = () =>
  z.object({
    parentId: z
      .number({
        error: t`请选择父级菜单`,
      })
      .nonnegative({ error: t`请选择父级菜单` }),
    name: z
      .string({
        error: t`请输入菜单名称`,
      })
      .min(1, t`请输入菜单名称`),
    i18nKey: z.string().trim().max(128).optional(),
    scope: z.enum(menuScopes).optional(),
    visible: z.boolean().optional(),
  })

export const GetActionSysMenuWithDirectory = () =>
  getBasicSchema().extend({
    type: z.literal(MenuType.directory),
    sort: z
      .number()
      .min(0, { error: t`排序值必须大于等于0` })
      .optional(),
  })

export const GetActionSysMenuWithMenu = () =>
  getBasicSchema().extend({
    type: z.literal(MenuType.menu),
    sort: z
      .number()
      .min(0, { error: t`排序值必须大于等于0` })
      .optional(),
    icon: z
      .string({
        error: t`请选择图标`,
      })
      .min(1, t`请选择图标`),
    route: z
      .string({
        error: t`请输入路由地址`,
      })
      .min(1, t`请输入路由地址`),
  })

export const GetActionSysMenuWithButton = (
  actions: readonly string[] = fallbackMenuPermissionActions
) => {
  // 词表为空只可能是 tenants 查询还没回来，此时用兜底表，别让下拉变成空的
  const allowed = actions.length ? actions : fallbackMenuPermissionActions
  return getBasicSchema().extend({
    type: z.literal(MenuType.button),
    // 不用 z.enum：动作词表来自后端，运行期才知道取值，这里只校验「在词表内」
    permissionAction: z
      .string({ error: t`请选择权限动作` })
      .refine((value) => allowed.includes(value), {
        error: t`请选择权限动作`,
      }),
    apiIds: z
      .array(z.number(), {
        error: t`请选择关联的API`,
      })
      .min(1, t`请选择关联的API`),
  })
}

export const getActionSysMenuSchema = (actions?: readonly string[]) =>
  z.discriminatedUnion('type', [
    GetActionSysMenuWithDirectory(),
    GetActionSysMenuWithMenu(),
    GetActionSysMenuWithButton(actions),
  ])

export type ActionSysMenuWithDirectoryFormValues = z.infer<
  ReturnType<typeof GetActionSysMenuWithDirectory>
>
export type ActionSysMenuWithMenuFormValues = z.infer<
  ReturnType<typeof GetActionSysMenuWithMenu>
>
export type ActionSysMenuWithButtonFormValues = z.infer<
  ReturnType<typeof GetActionSysMenuWithButton>
>

export type ActionSysMenuFormValues = z.infer<
  ReturnType<typeof getActionSysMenuSchema>
>
