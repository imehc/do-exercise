import { createFileRoute } from '@tanstack/react-router'
import EmailSignIn from '~/features/auth/email-sign-in'

export const Route = createFileRoute('/(auth)/email-sign-in')({
  component: EmailSignIn,
})
