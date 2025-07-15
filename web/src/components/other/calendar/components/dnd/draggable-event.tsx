import React, { ReactNode } from 'react'
import { motion } from 'framer-motion'
import { EventDetailsDialog } from '../../components/dialogs/event-details-dialog'
import { useCalendar } from '../../contexts/calendar-context'
import { useDragDrop } from '../../contexts/drag-drop-context'
import { IEvent } from '../../interfaces'

interface DraggableEventProps {
  event: IEvent
  children: ReactNode
  className?: string
}

export function DraggableEvent({
  event,
  children,
  className,
}: DraggableEventProps) {
  const { readonly } = useCalendar()
  const { startDrag, endDrag, isDragging, draggedEvent } = useDragDrop()

  const isCurrentlyDragged = isDragging && draggedEvent?.id === event.id

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    e.stopPropagation()
  }

  return (
    <EventDetailsDialog event={event}>
      <motion.div
        className={`${className || ''} ${
          readonly
            ? 'cursor-pointer'
            : isCurrentlyDragged
              ? 'cursor-grabbing opacity-50'
              : 'cursor-grab'
        }`}
        draggable={!readonly}
        onClick={(e: React.MouseEvent<HTMLDivElement>) => handleClick(e)}
        onDragStart={(e) => {
          if (readonly) {
            e.preventDefault()
            return
          }
          ;(e as DragEvent).dataTransfer!.setData(
            'text/plain',
            event.id.toString()
          )
          startDrag(event)
        }}
        onDragEnd={() => {
          if (!readonly) {
            endDrag()
          }
        }}
      >
        {children}
      </motion.div>
    </EventDetailsDialog>
  )
}
