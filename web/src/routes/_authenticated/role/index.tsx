import { createFileRoute } from '@tanstack/react-router'
import Role from '~/features/role'

export const Route = createFileRoute('/_authenticated/role/')({
  component: Role,
})
