import { z } from 'zod'
import { t } from '@lingui/core/macro'
import { MenuType } from '~/do-exercise-api'

const getBasicSchema = () =>
  z.object({
    parentId: z
      .number({
        required_error: t`请选择父级菜单`,
        invalid_type_error: t`请选择父级菜单`,
      })
      .nonnegative({ message: t`请选择父级菜单` }),
    name: z
      .string({ required_error: t`请输入菜单名称` })
      .min(1, t`请输入菜单名称`),
    visible: z.boolean().optional(),
  })

export const GetActionSysMenuWithDirectory = () =>
  getBasicSchema().extend({
    type: z.literal(MenuType.directory),
    sort: z.coerce
      .number()
      .min(0, { message: t`排序值必须大于等于0` })
      .optional(),
  })

export const GetActionSysMenuWithMenu = () =>
  getBasicSchema().extend({
    type: z.literal(MenuType.menu),
    sort: z.coerce
      .number()
      .min(0, { message: t`排序值必须大于等于0` })
      .optional(),
    icon: z.string({ required_error: t`请选择图标` }).min(1, t`请选择图标`),
    route: z
      .string({ required_error: t`请输入路由地址` })
      .min(1, t`请输入路由地址`),
    component: z
      .string({ required_error: t`请选择组件路径` })
      .min(1, t`请选择组件路径`),
  })

export const GetActionSysMenuWithButton = () =>
  getBasicSchema().extend({
    type: z.literal(MenuType.button),
    permission: z
      .string({ required_error: t`请输入权限标识` })
      .min(1, t`请输入权限标识`),
    apiIds: z
      .array(z.number(), {
        required_error: t`请选择关联的API`,
      })
      .min(1, t`请选择关联的API`),
  })

export const getActionSysMenuSchema = () =>
  z.discriminatedUnion('type', [
    GetActionSysMenuWithDirectory(),
    GetActionSysMenuWithMenu(),
    GetActionSysMenuWithButton(),
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
