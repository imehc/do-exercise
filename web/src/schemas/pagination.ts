import { z } from 'zod'

export const paginationSchema = z.object({
  page: z.number().min(1).default(1).catch(1),
  pageSize: z.number().min(1).default(10).catch(10),
})
