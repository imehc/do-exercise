import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import OperationLog from '~/features/operation-log'

export const Route = createFileRoute('/_authenticated/operation-log/')({
  component: OperationLog,
  validateSearch: paginationSchema,
})
