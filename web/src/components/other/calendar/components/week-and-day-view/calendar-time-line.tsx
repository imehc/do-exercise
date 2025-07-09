import { useEffect, useState } from 'react'
import { useCalendar } from '../../contexts/calendar-context'
import { formatTime } from '../../helpers'

export function CalendarTimeline() {
  const { use24HourFormat } = useCalendar()
  const [currentTime, setCurrentTime] = useState(new Date())

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), 60 * 1000)
    return () => clearInterval(timer)
  }, [])

  const getCurrentTimePosition = () => {
    const minutes = currentTime.getHours() * 60 + currentTime.getMinutes()
    return (minutes / 1440) * 100
  }

  const formatCurrentTime = () => {
    return formatTime(currentTime, use24HourFormat)
  }

  return (
    <div
      className='border-primary pointer-events-none absolute inset-x-0 z-50 border-t'
      style={{ top: `${getCurrentTimePosition()}%` }}
    >
      <div className='bg-primary absolute -top-1.5 -left-1.5 size-3 rounded-full'></div>

      <div className='bg-background text-primary absolute -left-18 flex w-16 -translate-y-1/2 justify-end pr-1 text-xs font-medium'>
        {formatCurrentTime()}
      </div>
    </div>
  )
}
