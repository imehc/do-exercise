import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import Api from '~/features/api'

export const Route = createFileRoute('/_authenticated/api/')({
  component: Api,
  validateSearch: paginationSchema,
})
