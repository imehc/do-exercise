import { createFileRoute } from '@tanstack/react-router'
import Token from '~/features/token'

export const Route = createFileRoute('/_authenticated/token/')({
  component: Token,
})
