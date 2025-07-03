import { useCallback } from 'react'
import { useMutation } from '@tanstack/react-query'
import {
  IconPlayerPlay,
  IconPlayerStop,
  IconPlaystationCircle,
} from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  ExecuteSysJobRequest,
  JobStatus,
  StartSysJobRequest,
  StopSysJobRequest,
  SysJob,
  SystemJobApi,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Button } from '~/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '~/components/ui/tooltip'

interface Props {
  row: SysJob
  type: 'start' | 'stop' | 'execute'
  refresh?: () => void
}
export function TaskActionButton({ row, type, refresh }: Props) {
  const systemJobApi = useApi(SystemJobApi)

  const { isPending: isPendingStart, mutate: start } = useMutation({
    mutationFn: (values: StartSysJobRequest) =>
      systemJobApi.startSysJob(values),
    onSuccess: () => {
      toast.success(t`启动成功`)
      refresh?.()
    },
  })

  const { isPending: isPendingStop, mutate: stop } = useMutation({
    mutationFn: (values: StopSysJobRequest) => systemJobApi.stopSysJob(values),
    onSuccess: () => {
      toast.success(t`停止成功`)
      refresh?.()
    },
  })

  const { isPending: isPendingExecute, mutate: execute } = useMutation({
    mutationFn: (values: ExecuteSysJobRequest) =>
      systemJobApi.executeSysJob(values),
    onSuccess: () => {
      toast.success(t`执行成功`)
      refresh?.()
    },
  })

  const handle = useCallback(() => {
    switch (type) {
      case 'start':
        start({ id: row.id })
        break
      case 'stop':
        stop({ id: row.id })
        break
      default:
        execute({ id: row.id })
        break
    }
  }, [execute, row.id, start, stop, type])

  const isPending = isPendingStart || isPendingStop || isPendingExecute

  return (
    <div className='w-fit text-nowrap'>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant='outline'
              size='icon'
              disabled={
                isPending ||
                (row.status === JobStatus.normal && type === 'start') ||
                (row.status === JobStatus.paused && type === 'stop')
              }
              onClick={handle}
            >
              {type === 'start' ? (
                <IconPlayerPlay />
              ) : type === 'stop' ? (
                <IconPlayerStop />
              ) : (
                <IconPlaystationCircle />
              )}
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            {type === 'start' ? (
              <Trans>开始任务</Trans>
            ) : type === 'stop' ? (
              <Trans>停止任务</Trans>
            ) : (
              <Trans>执行任务</Trans>
            )}
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  )
}
