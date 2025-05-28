import { z } from 'zod'

export const apiActionSchema = z.object({
  path: z.string(),
  method: z.string(),
  description: z.string().min(1, { message: '请输入菜单' }),
  group: z.string().optional(),
  disabled: z.boolean().optional(),
  sort: z.coerce.number().optional(),
})

export type ApiActionFormValues = z.infer<typeof apiActionSchema>
