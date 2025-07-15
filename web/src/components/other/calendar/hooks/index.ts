import { useState } from 'react'
import { useCalendar } from '~/components/other/calendar/contexts/calendar-context'

export function useDisclosure({
  defaultIsOpen = false,
}: { defaultIsOpen?: boolean } = {}) {
  const [isOpen, setIsOpen] = useState(defaultIsOpen)

  const onOpen = () => setIsOpen(true)
  const onClose = () => setIsOpen(false)
  const onToggle = () => setIsOpen((currentValue) => !currentValue)

  return { onOpen, onClose, isOpen, onToggle }
}

export const useFilteredEvents = () => {
  const { events, selectedDate, selectedUserId } = useCalendar()

  return events.filter((event) => {
    const itemStartDate = new Date(event.startDate)
    const itemEndDate = new Date(event.endDate)

    const monthStart = new Date(
      selectedDate.getFullYear(),
      selectedDate.getMonth(),
      1
    )
    const monthEnd = new Date(
      selectedDate.getFullYear(),
      selectedDate.getMonth() + 1,
      0
    )

    const isInSelectedMonth =
      itemStartDate <= monthEnd && itemEndDate >= monthStart
    const isUserMatch =
      selectedUserId === 'all' || event.user.id === selectedUserId
    return isInSelectedMonth && isUserMatch
  })
}

export const useLocalStorage = <T>(
  key: string,
  initialValue: T
): [T, (value: T) => void] => {
  const readValue = (): T => {
    if (typeof window === 'undefined') {
      return initialValue
    }

    try {
      const item = window.localStorage.getItem(key)
      return item ? (JSON.parse(item) as T) : initialValue
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error)
      return initialValue
    }
  }

  const [storedValue, setStoredValue] = useState<T>(readValue)

  const setValue = (value: T) => {
    try {
      const valueToStore =
        value instanceof Function ? value(storedValue) : value
      setStoredValue(valueToStore)
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore))
      }
    } catch (error) {
      console.warn(`Error setting localStorage key "${key}":`, error)
    }
  }

  return [storedValue, setValue]
}
