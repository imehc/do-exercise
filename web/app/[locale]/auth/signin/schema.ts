import { z } from "zod";

export const signinSchema = z.object({
  username: z
    .string({ required_error: 'accountFail.required' })
    .trim()
    .min(1, 'accountFail.required')
    .min(4, 'accountFail.min')
    .max(8, 'accountFail.max')
    .regex(/^[a-zA-Z0-9]+$/, 'accountFail.pattern')
    .regex(/^[a-zA-Z]/, 'accountFail.prefix'),
  password: z
    .string({ required_error: 'passwordFail.required' })
    .trim()
    .min(1, 'passwordFail.required')
    .min(6, 'passwordFail.min')
    .max(16, 'passwordFail.max'),
  captcha: z
    .string({ required_error: 'captchaFail.required' })
    .trim()
    .min(1, 'captchaFail.required'),
  captchaId: z
    .string({ required_error: 'captchaIdFail.required' })
    .min(1, 'captchaIdFail.required'),
})

export type SigninSchema = z.infer<typeof signinSchema>