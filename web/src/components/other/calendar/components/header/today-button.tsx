import { motion } from 'framer-motion'
import { useDateFormat } from '~/hooks/use-date-locale'
import { buttonHover, transition } from '../../animations'
import { useCalendar } from '../../contexts/calendar-context'

export function TodayButton() {
  const { setSelectedDate } = useCalendar()
  const { formatMonth } = useDateFormat()

  const today = new Date()
  const handleClick = () => setSelectedDate(today)

  return (
    <motion.button
      className='focus-visible:ring-ring flex size-14 flex-col items-start overflow-hidden rounded-lg border focus-visible:ring-1 focus-visible:outline-none'
      onClick={handleClick}
      variants={buttonHover}
      whileHover='hover'
      whileTap='tap'
      transition={transition}
    >
      <motion.p
        className='bg-primary text-primary-foreground flex h-6 w-full items-center justify-center text-center text-xs font-semibold'
        initial={{ y: -10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.1, ...transition }}
      >
        {formatMonth(today)}
      </motion.p>
      <motion.p
        className='flex w-full items-center justify-center text-lg font-bold'
        initial={{ y: 10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.2, ...transition }}
      >
        {today.getDate()}
      </motion.p>
    </motion.button>
  )
}
