import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getApiActionSchema = () =>
  z.object({
    path: z.string(),
    method: z.string(),
    description: z.string().min(1, { message: t`请输入菜单` }),
    group: z.string().optional(),
    disabled: z.boolean().optional(),
    sort: z.coerce.number().optional(),
  })

export type ApiActionFormValues = z.infer<ReturnType<typeof getApiActionSchema>>
