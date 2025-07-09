import { FC } from 'react'
import { format, parseISO } from 'date-fns'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { cn } from '~/lib/utils'
import { useDateFormat } from '~/hooks/use-date-locale'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '~/components/ui/command'
import { EventDetailsDialog } from '../../components/dialogs/event-details-dialog'
import { EventBullet } from '../../components/month-view/event-bullet'
import { useCalendar } from '../../contexts/calendar-context'
import {
  formatTime,
  getBgColor,
  getColorClass,
  getFirstLetters,
  toCapitalize,
  useGetEventsByMode,
  groupBy,
} from '../../helpers'

export const AgendaEvents: FC = () => {
  useAtomValue(languageAtom)
  const { events, use24HourFormat, badgeVariant, agendaModeGroupBy } =
    useCalendar()
  const { formatDateWithWeek, formatDateTime } = useDateFormat()

  const eventsByMode = groupBy(useGetEventsByMode(events), (event) => {
    return agendaModeGroupBy === 'date'
      ? format(parseISO(event.startDate), 'yyyy-MM-dd')
      : event.color
  })

  const groupedAndSortedEvents = Object.entries(eventsByMode).sort(
    (a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime()
  )

  return (
    <Command className='h-[80vh] bg-transparent py-4'>
      <div className='mx-4 mb-4'>
        <CommandInput placeholder={t`键入命令或搜索...`} />
      </div>
      <CommandList className='max-h-max border-t px-3'>
        {groupedAndSortedEvents.map(([date, groupedEvents]) => (
          <CommandGroup
            key={date}
            heading={
              agendaModeGroupBy === 'date'
                ? formatDateWithWeek(parseISO(date))
                : toCapitalize(groupedEvents[0]?.color || '')
            }
          >
            {groupedEvents?.map((event) => (
              <CommandItem
                key={event.id}
                className={cn(
                  'data-[selected=true]:bg-bg data-[selected=true]:text-none mb-2 rounded-md border p-4 transition-all hover:cursor-pointer',
                  {
                    [getColorClass(event.color)]: badgeVariant === 'colored',
                    'hover:bg-zinc-200 dark:hover:bg-gray-900':
                      badgeVariant === 'dot',
                    'hover:opacity-60': badgeVariant === 'colored',
                  }
                )}
              >
                <EventDetailsDialog event={event}>
                  <div className='flex w-full items-center justify-between gap-4'>
                    <div className='flex items-center gap-2'>
                      {badgeVariant === 'dot' ? (
                        <EventBullet color={event.color} />
                      ) : (
                        <Avatar>
                          <AvatarImage src='' alt='@shadcn' />
                          <AvatarFallback className={getBgColor(event.color)}>
                            {getFirstLetters(event.title)}
                          </AvatarFallback>
                        </Avatar>
                      )}
                      <div className='flex flex-col'>
                        <p
                          className={cn({
                            'font-medium': badgeVariant === 'dot',
                            'text-foreground': badgeVariant === 'dot',
                          })}
                        >
                          {event.title}
                        </p>
                        <p className='text-muted-foreground line-clamp-1 w-1/3 text-sm text-ellipsis md:text-clip'>
                          {event.description}
                        </p>
                      </div>
                    </div>
                    <div className='flex w-40 items-center justify-center gap-1'>
                      {agendaModeGroupBy === 'date' ? (
                        <>
                          <p className='text-sm'>
                            {formatTime(event.startDate, use24HourFormat)}
                          </p>
                          <span className='text-muted-foreground'>-</span>
                          <p className='text-sm'>
                            {formatTime(event.endDate, use24HourFormat)}
                          </p>
                        </>
                      ) : (
                        formatDateTime(parseISO(event.startDate))
                      )}
                    </div>
                  </div>
                </EventDetailsDialog>
              </CommandItem>
            ))}
          </CommandGroup>
        ))}
        <CommandEmpty>
          <Trans>未找到结果。</Trans>
        </CommandEmpty>
      </CommandList>
    </Command>
  )
}
