import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getSchema = () =>
  z.object({
    name: z
      .string({
        error: t`请输入角色名称`,
      })
      .min(1, t`请输入角色名称`),
    code: z
      .string({
        error: t`请输入角色编码`,
      })
      .min(1, t`请输入角色编码`),
    menuIds: z
      .array(z.number(), {
        error: t`请选择关联的菜单`,
      })
      .min(1, t`请选择关联的菜单`),
  })

export type ActionSysRoleFormValues = z.infer<ReturnType<typeof getSchema>>
