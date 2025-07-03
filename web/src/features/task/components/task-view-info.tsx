import { useQuery } from '@tanstack/react-query'
import { IconMessage } from '@tabler/icons-react'
import { i18n } from '@lingui/core'
import { Trans } from '@lingui/react/macro'
import { JobStatus, SysJob, SystemJobApi } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
import { useApi } from '~/hooks/use-api'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from '~/components/ui/drawer'
import { StatusRenderer } from '~/components/other'
import { callVisibleTypes } from '~/features/menu/data/data'
import { scheduleStatusTypes } from '../data/data'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function TaskViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysJob>()
  const systemJobApi = useApi(SystemJobApi)
  const { data, isLoading: isLoadingUser } = useQuery({
    queryKey: ['findSysJob', currentRow?.id],
    queryFn: () => systemJobApi.findSysJob({ id: currentRow?.id as number }),
    enabled: !!currentRow?.id,
  })

  return (
    <Drawer
      shouldScaleBackground
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
      direction='bottom'
    >
      <DrawerContent className='max-h-[85vh]'>
        <DrawerHeader className='text-left'>
          <DrawerTitle className='flex items-center gap-2'>
            <IconMessage /> <Trans>定时任务详情</Trans>
          </DrawerTitle>
          <DrawerDescription>
            <Trans>查看定时任务详细信息。</Trans>
          </DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoadingUser} data={data}>
          {(job) => (
            <div className='grid gap-6 overflow-y-auto px-4'>
              <div className='space-y-3'>
                <h4 className='text-lg font-medium'>
                  <Trans>基本信息</Trans>
                </h4>
                <div className='grid gap-3 text-sm'>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>任务名称</Trans>
                    </div>
                    <div className='col-span-2'>{job.name}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>任务分组</Trans>
                    </div>
                    <div className='col-span-2'>{job.jobGroup}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>cron表达式</Trans>
                    </div>
                    <div className='col-span-2'>{job.cronExpression}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>执行命令</Trans>
                    </div>
                    <div className='col-span-2'>{job.command}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>状态</Trans>
                    </div>
                    <div className='col-span-2'>
                      <Badge
                        variant='outline'
                        className={cn(
                          'capitalize',
                          scheduleStatusTypes.get(job.status)
                        )}
                      >
                        {job.status === JobStatus.normal ? (
                          <Trans>正常</Trans>
                        ) : (
                          <Trans>暂停</Trans>
                        )}
                      </Badge>
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>是否并发执行</Trans>
                    </div>
                    <div className='col-span-2'>
                      <Badge
                        variant='outline'
                        className={cn(
                          'capitalize',
                          callVisibleTypes.get(job.concurrent ?? false)
                        )}
                      >
                        {job.concurrent ? <Trans>是</Trans> : <Trans>否</Trans>}
                      </Badge>
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>任务描述</Trans>
                    </div>
                    <div className='col-span-2'>{job.description}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>上次执行时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {job.lastTime
                        ? i18n.date(job.lastTime, {
                          dateStyle: 'short',
                          timeStyle: 'medium',
                        })
                        : '-'}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>下次执行时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {job.nextTime
                        ? i18n.date(job.nextTime, {
                          dateStyle: 'short',
                          timeStyle: 'medium',
                        })
                        : '-'}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>执行次数</Trans>
                    </div>
                    <div className='col-span-2'>{job.times || 0}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>重试次数</Trans>
                    </div>
                    <div className='col-span-2'>{job.retryTimes || 0}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>重试间隔(秒)</Trans>
                    </div>
                    <div className='col-span-2'>{job.retryInterval || 0}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>执行超时时间(秒)</Trans>
                    </div>
                    <div className='col-span-2'>{job.timeout || 0}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>创建时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {i18n.date(job.createdAt, {
                        dateStyle: 'short',
                        timeStyle: 'medium',
                      })}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>更新时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {i18n.date(job.updatedAt, {
                        dateStyle: 'short',
                        timeStyle: 'medium',
                      })}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </StatusRenderer>

        <DrawerFooter className='gap-y-2'>
          <DrawerClose asChild>
            <Button variant='outline'>
              <Trans>关闭</Trans>
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
