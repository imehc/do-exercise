'use client'

import { useMemo } from 'react'
import { isToday, startOfDay } from 'date-fns'
import { cva } from 'class-variance-authority'
import { motion } from 'framer-motion'
import { cn } from '~/lib/utils'
import { staggerContainer, transition } from '../../animations'
import { EventListDialog } from '../../components/dialogs/events-list-dialog'
import { DroppableArea } from '../../components/dnd/droppable-area'
import { EventBullet } from '../../components/month-view/event-bullet'
import { MonthEventBadge } from '../../components/month-view/month-event-badge'
import { getMonthCellEvents } from '../../helpers'
import type { ICalendarCell, IEvent } from '../../interfaces'

interface IProps {
  cell: ICalendarCell
  events: IEvent[]
  eventPositions: Record<string, number>
}

const MAX_VISIBLE_EVENTS = 3

export const dayCellVariants = cva('text-white', {
  variants: {
    color: {
      blue: 'bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-400 ',
      green:
        'bg-green-600 dark:bg-green-500 hover:bg-green-700 dark:hover:bg-green-400',
      red: 'bg-red-600 dark:bg-red-500 hover:bg-red-700 dark:hover:bg-red-400',
      yellow:
        'bg-yellow-600 dark:bg-yellow-500 hover:bg-yellow-700 dark:hover:bg-yellow-400',
      purple:
        'bg-purple-600 dark:bg-purple-500 hover:bg-purple-700 dark:hover:bg-purple-400',
      orange:
        'bg-orange-600 dark:bg-orange-500 hover:bg-orange-700 dark:hover:bg-orange-400',
      gray: 'bg-gray-600 dark:bg-gray-500 hover:bg-gray-700 dark:hover:bg-gray-400',
    },
  },
  defaultVariants: {
    color: 'blue',
  },
})

export function DayCell({ cell, events, eventPositions }: IProps) {
  const { day, currentMonth, date } = cell

  const cellEvents = useMemo(
    () => getMonthCellEvents(date, events, eventPositions),
    [date, events, eventPositions]
  )
  const isSunday = date.getDay() === 0

  return (
    <motion.div
      className={cn(
        'flex flex-col gap-1 border-t border-l py-1.5',
        isSunday && 'border-l-0'
      )}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={transition}
    >
      <DroppableArea date={date}>
        <motion.span
          className={cn(
            'mb-1 flex h-6 w-6 translate-x-1 items-center justify-center rounded-full px-1 text-xs font-semibold lg:px-2', // mb-0.5 here, always applied
            !currentMonth && 'text-muted-foreground',
            isToday(date) && 'bg-primary text-primary-foreground'
          )}
          whileHover={{ scale: 1.1 }}
          transition={transition}
        >
          {day}
        </motion.span>

        <motion.div
          className={cn(
            'flex h-6 gap-1 px-2 lg:min-h-[94px] lg:flex-col lg:gap-2 lg:px-0',
            !currentMonth && 'opacity-50'
          )}
          variants={staggerContainer}
        >
          {[0, 1, 2].map((position) => {
            const event = cellEvents.find((e) => e.position === position)
            const eventKey = event
              ? `event-${event.id}-${position}`
              : `empty-${position}`

            return (
              <motion.div
                key={eventKey}
                className='lg:flex-1'
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: position * 0.1, ...transition }}
              >
                {event && (
                  <>
                    <EventBullet className='lg:hidden' color={event.color} />
                    <MonthEventBadge
                      className='hidden lg:flex'
                      event={event}
                      cellDate={startOfDay(date)}
                    />
                  </>
                )}
              </motion.div>
            )
          })}
        </motion.div>

        {cellEvents.length > MAX_VISIBLE_EVENTS && (
          <motion.div
            className={cn(
              'text-muted-foreground my-2 h-4.5 px-1.5 text-end text-xs font-semibold',
              !currentMonth && 'opacity-50'
            )}
            initial={{ opacity: 0, y: 5 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3, ...transition }}
          >
            <EventListDialog date={date} events={cellEvents} />
          </motion.div>
        )}
      </DroppableArea>
    </motion.div>
  )
}
