import { z } from 'zod'
import {
  passwordRule,
  usernameRule,
} from '~/features/user/schemas/action-schema'

export const signInActionSchema = z.object({
  username: usernameRule,
  password: passwordRule,
  captchaId: z.string({ required_error: '验证码ID不能为空' }),
  captcha: z
    .string({ required_error: '请输入验证码' })
    .min(1, { message: '请输入验证码' }),
  publicKey: z.string({ required_error: '公钥不能为空' }),
})

export type SignInActionFormValues = z.infer<typeof signInActionSchema>
