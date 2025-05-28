import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import Role from '~/features/role'

export const Route = createFileRoute('/_authenticated/role/')({
  component: Role,
  validateSearch: paginationSchema,
})
