import { ReactNode } from 'react'
import { useCalendar } from '../../contexts/calendar-context'
import { useDragDrop } from '../../contexts/drag-drop-context'

interface DroppableAreaProps {
  date: Date
  hour?: number
  minute?: number
  children: ReactNode
  className?: string
}

export function DroppableArea({
  date,
  hour,
  minute,
  children,
  className,
}: DroppableAreaProps) {
  const { readonly } = useCalendar()
  const { handleEventDrop, isDragging } = useDragDrop()

  return (
    <div
      className={`${className || ''} ${isDragging && !readonly ? 'drop-target' : ''}`}
      onDragOver={(e) => {
        if (readonly) return
        // Prevent default to allow drop
        e.preventDefault()
        e.currentTarget.classList.add('bg-primary/10')
      }}
      onDragLeave={(e) => {
        if (readonly) return
        e.currentTarget.classList.remove('bg-primary/10')
      }}
      onDrop={(e) => {
        if (readonly) return
        e.preventDefault()
        e.currentTarget.classList.remove('bg-primary/10')
        handleEventDrop(date, hour, minute)
      }}
    >
      {children}
    </div>
  )
}
