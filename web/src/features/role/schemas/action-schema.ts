import { z } from 'zod'

export const schema = z.object({
  name: z.string({ required_error: '请输入角色名称' }).min(1, '请输入角色名称'),
  code: z.string({ required_error: '请输入角色编码' }).min(1, '请输入角色编码'),
  menuIds: z
    .array(z.number(), {
      required_error: '请选择关联的菜单',
    })
    .min(1, '请选择关联的菜单'),
})

export type ActionSysRoleFormValues = z.infer<typeof schema>
