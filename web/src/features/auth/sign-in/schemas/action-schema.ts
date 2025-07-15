import { z } from 'zod'
import { t } from '@lingui/core/macro'
import {
  getPasswordRule,
  getUsernameRule,
} from '~/features/user/schemas/action-schema'

export const getSignInActionSchema = () =>
  z.object({
    username: getUsernameRule(),
    password: getPasswordRule(),
    captchaId: z.string({
      error: t`验证码ID不能为空`,
    }),
    captcha: z
      .string({
        error: t`请输入验证码`,
      })
      .min(1, { error: t`请输入验证码` }),
    publicKey: z.string({
      error: t`公钥不能为空`,
    }),
  })

export type SignInActionFormValues = z.infer<
  ReturnType<typeof getSignInActionSchema>
>
