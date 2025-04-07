import { z } from "zod"

// TODO: 更多的验证
export const loginSchema = z.object({
    username: z.string(),
    password: z.string(),
    captchaId: z.string().optional(),
    captcha: z.string().optional(),
    publicKey: z.string().optional(),
})

export type LoginSchemaType = z.infer<typeof loginSchema>