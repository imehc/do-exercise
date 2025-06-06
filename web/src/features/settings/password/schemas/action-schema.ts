import { z } from 'zod'
import { passwordRule } from '~/features/user/schemas/action-schema'

export const passwordSchema = z
  .object({
    oldPassword: passwordRule,
    password: passwordRule,
    confirmPassword: passwordRule,
  })
  .superRefine(({ oldPassword, password, confirmPassword }, ctx) => {
    if (oldPassword.trim() === password.trim()) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['password'],
        message: 'The new password cannot be the same as the old password.',
      })
    }
    if (password.trim() !== confirmPassword.trim()) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['confirmPassword'],
        message: 'The passwords entered twice are inconsistent.',
      })
    }
  })

export type PasswordSchemaFormValues = z.infer<typeof passwordSchema>
