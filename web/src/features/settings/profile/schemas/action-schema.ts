import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getProfileSchema = () =>
  z.object({
    nickname: z
      .string()
      .max(10, { error: t`昵称长度不能超过10个字符` })
      .optional(),
    avatar: z.string().optional(),
  })

export type ProfileSchemaFormValues = z.infer<
  ReturnType<typeof getProfileSchema>
>
