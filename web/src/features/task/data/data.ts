import { JobStatus } from '~/do-exercise-api'

export const scheduleStatusTypes = new Map<JobStatus, string>([
  [
    JobStatus.normal, // 正常运行
    'bg-emerald-100/30 text-emerald-900 dark:text-emerald-200 border-emerald-200 dark:bg-emerald-900/30 dark:border-emerald-800',
  ],
  [
    JobStatus.paused, // 已暂停
    'bg-amber-100/30 text-amber-900 dark:text-amber-200 border-amber-200 dark:bg-amber-900/30 dark:border-amber-800',
  ],
])
