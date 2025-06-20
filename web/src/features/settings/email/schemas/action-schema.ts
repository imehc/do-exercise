import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getEmailSchema = () =>
  z.object({
    email: z
      .string({ required_error: t`请输入您的邮箱` })
      .min(1, { message: t`请输入您的邮箱` })
      .email({ message: t`邮箱无效` }),
    code: z
      .string({ required_error: t`请输入邮箱验证码` })
      .min(1, { message: t`请输入邮箱验证码` }),
  })

export type EmailSchemaFormValues = z.infer<ReturnType<typeof getEmailSchema>>
