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
    .check(({ issues, value }) => {
      if (value.oldPassword.trim() === value.password.trim()) {
        issues.push({
          code: 'custom',
          path: ['password'],
          message: t`新密码不能与旧密码相同`,
          input: value,
        })
      }
      if (value.password.trim() !== value.confirmPassword.trim()) {
        issues.push({
          code: 'custom',
          path: ['confirmPassword'],
          message: t`两次输入的密码不一致`,
          input: value,
        })
      }
    })

export type PasswordSchemaFormValues = z.infer<
  ReturnType<typeof getPasswordSchema>
>
