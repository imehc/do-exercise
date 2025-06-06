import { z } from 'zod'

export const emailSchema = z.object({
  email: z
    .string({ required_error: 'Please enter your email' })
    .min(1, { message: 'Please enter your email' })
    .email({ message: 'Invalid email address' }),
  code: z
    .string({ required_error: 'Please enter your verification code' })
    .min(1, { message: 'Please enter your verification code' }),
})

export type EmailSchemaFormValues = z.infer<typeof emailSchema>
