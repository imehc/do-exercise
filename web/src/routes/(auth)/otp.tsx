import { z } from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import Otp from '~/features/auth/otp'

export const Route = createFileRoute('/(auth)/otp')({
  component: Otp,
  validateSearch: z.object({
    email: z
      .string()
      .email({ message: 'Invalid email address' })
      .catch('')
      .optional(),
  }),
})
