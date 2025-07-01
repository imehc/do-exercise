import { createFileRoute } from '@tanstack/react-router'
import { paginationSchema } from '~/schemas/pagination'
import Task from '~/features/task'

export const Route = createFileRoute('/_authenticated/task/')({
  component: Task,
  validateSearch: paginationSchema,
})
