import { areIntervalsOverlapping, parseISO } from 'date-fns'
import { EventBlock } from '../../components/week-and-day-view/event-block'
import { getEventBlockStyle } from '../../helpers'
import { IEvent } from '../../interfaces'

interface RenderGroupedEventsProps {
  groupedEvents: IEvent[][]
  day: Date
}

export function RenderGroupedEvents({
  groupedEvents,
  day,
}: RenderGroupedEventsProps) {
  return groupedEvents.map((group, groupIndex) =>
    group.map((event) => {
      let style = getEventBlockStyle(
        event,
        day,
        groupIndex,
        groupedEvents.length
      )
      const hasOverlap = groupedEvents.some(
        (otherGroup, otherIndex) =>
          otherIndex !== groupIndex &&
          otherGroup.some((otherEvent) =>
            areIntervalsOverlapping(
              {
                start: parseISO(event.startDate),
                end: parseISO(event.endDate),
              },
              {
                start: parseISO(otherEvent.startDate),
                end: parseISO(otherEvent.endDate),
              }
            )
          )
      )

      if (!hasOverlap) style = { ...style, width: '100%', left: '0%' }

      return (
        <div key={event.id} className='absolute p-1' style={style}>
          <EventBlock event={event} />
        </div>
      )
    })
  )
}
