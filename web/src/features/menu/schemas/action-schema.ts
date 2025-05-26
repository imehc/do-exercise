import { z } from 'zod'
import { MenuType } from '~/do-exercise-api'

const basicSchema = z.object({
  parentId: z
    .number({
      required_error: '请选择父级菜单',
      invalid_type_error: '请选择父级菜单',
    })
    .nonnegative({ message: '请选择父级菜单' }),
  name: z.string({ required_error: '请输入菜单名称' }).min(1, '请输入菜单名称'),
})

export const ActionSysMenuWithDirectory = basicSchema.extend({
  type: z.literal(MenuType.directory),
  sort: z.coerce.number().min(0, { message: '排序值必须大于等于0' }).optional(),
})

export const ActionSysMenuWithMenu = basicSchema.extend({
  type: z.literal(MenuType.menu),
  sort: z.coerce.number().min(0, { message: '排序值必须大于等于0' }).optional(),
  icon: z.string({ required_error: '请选择图标' }).min(1, '请选择图标'),
  route: z
    .string({ required_error: '请输入路由地址' })
    .min(1, '请输入路由地址'),
  component: z
    .string({ required_error: '请选择组件路径' })
    .min(1, '请选择组件路径'),
  visible: z.boolean().optional(),
})

export const ActionSysMenuWithButton = basicSchema.extend({
  type: z.literal(MenuType.button),
  permission: z
    .string({ required_error: '请输入权限标识' })
    .min(1, '请输入权限标识'),
  apiIds: z
    .array(z.number(), {
      required_error: '请选择关联的API',
    })
    .min(1, '请选择关联的API'),
})

export const actionSysMenuSchema = z.discriminatedUnion('type', [
  ActionSysMenuWithDirectory,
  ActionSysMenuWithMenu,
  ActionSysMenuWithButton,
])

export type ActionSysMenuFormValues = z.infer<typeof actionSysMenuSchema>
