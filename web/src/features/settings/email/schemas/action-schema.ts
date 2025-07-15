import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getEmailSchema = () =>
  z.object({
    email: z.email({
      error: (issue) =>
        issue.input === undefined ? t`请输入您的邮箱` : t`参数不合法`,
    }),
    code: z
      .string({
        error: t`请输入邮箱验证码`,
      })
      .min(1, { error: t`请输入邮箱验证码` }),
  })

export type EmailSchemaFormValues = z.infer<ReturnType<typeof getEmailSchema>>
