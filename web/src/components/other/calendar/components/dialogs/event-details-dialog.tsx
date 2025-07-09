'use client'

import { ReactNode } from 'react'
import { parseISO } from 'date-fns'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { Calendar, Clock, Text, User } from 'lucide-react'
import { toast } from 'sonner'
import { useDateFormat } from '~/hooks/use-date-locale'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from '~/components/ui/dialog'
import { ScrollArea } from '~/components/ui/scroll-area'
import { AddEditEventDialog } from '../../components/dialogs/add-edit-event-dialog'
import { useCalendar } from '../../contexts/calendar-context'
import type { IEvent } from '../../interfaces'

interface IProps {
  event: IEvent
  children: ReactNode
}

export function EventDetailsDialog({ event, children }: IProps) {
  const startDate = parseISO(event.startDate)
  const endDate = parseISO(event.endDate)
  const { removeEvent, readonly } = useCalendar()
  const { formatDateTime } = useDateFormat()

  const deleteEvent = (eventId: number) => {
    try {
      removeEvent(eventId)
      toast.success(t`事件删除成功`)
    } catch {
      toast.error(t`删除事件错误`)
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{event.title}</DialogTitle>
        </DialogHeader>

        <ScrollArea className='max-h-[80vh]'>
          <div className='space-y-4 p-4'>
            <div className='flex items-start gap-2'>
              <User className='text-muted-foreground mt-1 size-4 shrink-0' />
              <div>
                <p className='text-sm font-medium'>
                  <Trans>负责人</Trans>
                </p>
                <p className='text-muted-foreground text-sm'>
                  {event.user.name}
                </p>
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <Calendar className='text-muted-foreground mt-1 size-4 shrink-0' />
              <div>
                <p className='text-sm font-medium'>
                  <Trans>开始日期</Trans>
                </p>
                <p className='text-muted-foreground text-sm'>
                  {formatDateTime(startDate)}
                </p>
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <Clock className='text-muted-foreground mt-1 size-4 shrink-0' />
              <div>
                <p className='text-sm font-medium'>
                  <Trans>结束日期</Trans>
                </p>
                <p className='text-muted-foreground text-sm'>
                  {formatDateTime(endDate)}
                </p>
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <Text className='text-muted-foreground mt-1 size-4 shrink-0' />
              <div>
                <p className='text-sm font-medium'>
                  <Trans>描述</Trans>
                </p>
                <p className='text-muted-foreground text-sm'>
                  {event.description}
                </p>
              </div>
            </div>
          </div>
        </ScrollArea>
        <div className='flex justify-end gap-2'>
          {!readonly && (
            <>
              <AddEditEventDialog event={event}>
                <Button variant='outline'>
                  <Trans>编辑</Trans>
                </Button>
              </AddEditEventDialog>
              <Button
                variant='destructive'
                onClick={() => {
                  deleteEvent(event.id)
                }}
              >
                <Trans>删除</Trans>
              </Button>
            </>
          )}
        </div>
        <DialogClose />
      </DialogContent>
    </Dialog>
  )
}
