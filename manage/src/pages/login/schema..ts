import { z } from "zod"

// TODO: 更多的验证
export const loginSchema = z.object({
    username: z.string(),
    password: z.string(),
    captchaId: z.string(),
    captcha: z.string(),
    publicKey: z.string(),
})

export type LoginSchemaType = z.infer<typeof loginSchema>