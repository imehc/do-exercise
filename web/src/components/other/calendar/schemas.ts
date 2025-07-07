import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getEventSchema = () =>
  z.object({
    title: z.string().min(1, t`标题是必需的`),
    description: z.string().min(1, t`描述是必需的`),
    startDate: z.date({
      required_error: t`开始日期是必需的`,
    }),
    endDate: z.date({
      required_error: t`结束日期是必需的`,
    }),
    color: z.enum(['blue', 'green', 'red', 'yellow', 'purple', 'orange'], {
      required_error: t`颜色是必需的`,
    }),
  })

export type TEventFormData = z.infer<ReturnType<typeof getEventSchema>>
