import { z } from 'zod'

export const profileSchema = z.object({
  nickname: z
    .string()
    .max(10, { message: '昵称长度不能超过10个字符' })
    .optional(),
  avatar: z.string().optional(),
})

export type ProfileSchemaFormValues = z.infer<typeof profileSchema>
