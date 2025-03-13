import { z } from 'zod';

export const signinSchema = z.object({
  username: z
    .string({ required_error: 'accountFail.required' })
    .trim()
    .nonempty('accountFail.required')
    .min(4, 'accountFail.min')
    .max(8, 'accountFail.max')
    .regex(/^[a-zA-Z0-9]+$/, 'accountFail.pattern')
    .regex(/^[a-zA-Z]/, 'accountFail.prefix'),
  password: z
    .string({ required_error: 'passwordFail.required' })
    .trim()
    .nonempty('passwordFail.required')
    .min(6, 'passwordFail.min')
    .max(16, 'passwordFail.max'),
  captcha: z
    .string({ required_error: 'captchaFail.required' })
    .trim()
    .nonempty('captchaFail.required'),
  captchaId: z
    .string({ required_error: 'captchaIdFail.required' })
    .nonempty('captchaIdFail.required'),
});

export type SigninSchema = z.infer<typeof signinSchema>;
