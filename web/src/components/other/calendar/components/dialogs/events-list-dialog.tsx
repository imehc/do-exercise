import { ReactNode } from 'react'
import { format } from 'date-fns'
import { Trans } from '@lingui/react/macro'
import { cn } from '~/lib/utils'
import {
  Dialog,
  DialogContent,
  DialogTrigger,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'
import { dayCellVariants } from '../../components/month-view/day-cell'
import { EventBullet } from '../../components/month-view/event-bullet'
import { useCalendar } from '../../contexts/calendar-context'
import { formatTime } from '../../helpers'
import { IEvent } from '../../interfaces'

interface EventListDialogProps {
  date: Date
  events: IEvent[]
  maxVisibleEvents?: number
  children?: ReactNode
}

export function EventListDialog({
  date,
  events,
  maxVisibleEvents = 3,
  children,
}: EventListDialogProps) {
  const cellEvents = events
  const hiddenEventsCount = Math.max(cellEvents.length - maxVisibleEvents, 0)
  const { badgeVariant, use24HourFormat } = useCalendar()

  const defaultTrigger = (
    <span className='cursor-pointer'>
      <span className='sm:hidden'>+{hiddenEventsCount}</span>
      <span className='my-1 hidden rounded-xl border px-2 py-0.5 sm:inline'>
        {hiddenEventsCount}
        <span className='mx-1'>
          <Trans>更多...</Trans>
        </span>
      </span>
    </span>
  )

  return (
    <Dialog>
      <DialogTrigger asChild>{children || defaultTrigger}</DialogTrigger>
      <DialogContent className='sm:max-w-[425px]'>
        <DialogHeader>
          <DialogTitle>
            <div className='flex items-center gap-2'>
              <EventBullet color={cellEvents[0]?.color} className='' />
              <p className='text-sm font-medium'>
                <Trans>{format(date, 'EEEE, MMMM d, yyyy')} 的事件</Trans>
              </p>
            </div>
          </DialogTitle>
        </DialogHeader>
        <div className='max-h-[60vh] space-y-2 overflow-y-auto'>
          {cellEvents.map((event) => (
            <div
              key={event.id}
              className={cn(
                'hover:bg-muted flex items-center gap-2 rounded-md border p-2',
                {
                  [dayCellVariants({ color: event.color })]:
                    badgeVariant === 'colored',
                }
              )}
            >
              <EventBullet color={event.color} className='' />
              <div className='flex-1'>
                <p className='text-sm font-medium'>{event.title}</p>
                <p
                  className={cn('text-xs', {
                    'text-muted': badgeVariant === 'colored',
                    'text-muted-foreground': badgeVariant === 'dot',
                  })}
                >
                  {formatTime(event.startDate, use24HourFormat)}
                </p>
              </div>
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  )
}
