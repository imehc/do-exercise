import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import User from '~/features/user'

export const Route = createFileRoute('/_authenticated/user/')({
  component: User,
  validateSearch: paginationSchema,
})
