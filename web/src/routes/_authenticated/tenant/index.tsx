import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import Tenant from '~/features/tenant'

export const Route = createFileRoute('/_authenticated/tenant/')({
  component: Tenant,
  validateSearch: paginationSchema,
})