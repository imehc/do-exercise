import { z } from 'zod'
import { t } from '@lingui/core/macro'
import { isValidCron } from 'cron-validator'
import { JobStatus } from '~/do-exercise-api'

export const getSchema = () =>
  z.object({
    name: z
      .string({
        error: t`请输入任务名称`,
      })
      .min(1, t`请输入任务名称`),
    jobGroup: z
      .string({
        error: t`请输入任务分组`,
      })
      .min(1, t`请输入任务分组`),
    cronExpression: z
      .string({
        error: t`请输入cron表达式`,
      })
      .min(1, t`请输入cron表达式`)
      .refine((val) => isValidCron(val, { seconds: true }), {
        error: t`请输入有效的cron表达式`,
      }),
    command: z
      .string({
        error: t`请输入执行命令`,
      })
      .min(1, t`请输入执行命令`),
    status: z.enum(JobStatus, {
      error: t`请选择状态`,
    }),
    description: z.string().optional(),
    concurrent: z.boolean().optional(),
    retryTimes: z.number().int().min(0).max(5).optional(),
    retryInterval: z.number().int().min(0).max(300).optional(),
    timeout: z.number().int().min(0).max(3600).optional(),
    isEdit: z.boolean().default(false).optional(),
  })

export type ActionSysJobFormValues = z.infer<ReturnType<typeof getSchema>>
