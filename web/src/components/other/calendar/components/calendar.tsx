import { CalendarBody } from '../components/calendar-body'
import { EventUpdateHandler } from '../components/event-update-handler'
import { CalendarHeader } from '../components/header/calendar-header'
import { CalendarProvider } from '../contexts/calendar-context'
import { DragDropProvider } from '../contexts/drag-drop-context'
import { IEvent, IUser } from '../interfaces'

interface Props {
  events: IEvent[]
  users: IUser[]
  readonly?: boolean
}
/**
 * 日历组件
 *
 * 来源 https://github.com/shadcn-ui/ui/discussions/3214
 *
 * 基于 https://github.com/yassir-jeraidi/full-calendar  修改，感谢 yassir-jeraidi
 *
 * 参考 https://github.com/lramos33/big-calendar 感谢 lramos33
 *
 * react-day-picker兼容 https://date-picker.luca-felix.com/
 *
 * 待完善参考 https://demo.fulleventcalendar.com/
 */
export function Calendar(props: Props) {
  return (
    <DragDropProvider>
      <CalendarProvider {...props} view='month'>
        <div className='w-full rounded-xl border'>
          <EventUpdateHandler />
          <CalendarHeader />
          <CalendarBody />
        </div>
      </CalendarProvider>
    </DragDropProvider>
  )
}
