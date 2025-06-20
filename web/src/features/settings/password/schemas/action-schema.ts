import { z } from 'zod'
import { t } from '@lingui/core/macro'
import { getPasswordRule } from '~/features/user/schemas/action-schema'

export const getPasswordSchema = () =>
  z
    .object({
      oldPassword: getPasswordRule(),
      password: getPasswordRule(),
      confirmPassword: getPasswordRule(),
    })
    .superRefine(({ oldPassword, password, confirmPassword }, ctx) => {
      if (oldPassword.trim() === password.trim()) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['password'],
          message: t`新密码不能与旧密码相同`,
        })
      }
      if (password.trim() !== confirmPassword.trim()) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['confirmPassword'],
          message: t`两次输入的密码不一致`,
        })
      }
    })

export type PasswordSchemaFormValues = z.infer<
  ReturnType<typeof getPasswordSchema>
>
