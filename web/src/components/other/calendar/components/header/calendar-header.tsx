'use client'

import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { motion } from 'framer-motion'
import {
  CalendarRange,
  Columns,
  Grid2X2,
  Grid3X3,
  LayoutList,
  List,
  Plus,
} from 'lucide-react'
import { Button } from '~/components/ui/button'
import { ButtonGroup } from '~/components/ui/button-group'
import { Toggle } from '~/components/ui/toggle'
import {
  buttonHover,
  slideFromLeft,
  slideFromRight,
  transition,
} from '../../animations'
import { AddEditEventDialog } from '../../components/dialogs/add-edit-event-dialog'
import { DateNavigator } from '../../components/header/date-navigator'
import FilterEvents from '../../components/header/filter'
import { TodayButton } from '../../components/header/today-button'
import { UserSelect } from '../../components/header/user-select'
import { Settings } from '../../components/settings/settings'
import { useCalendar } from '../../contexts/calendar-context'
import { useFilteredEvents } from '../../hooks'

export const MotionButton = motion.create(Button)

export function CalendarHeader() {
  const { view, setView, readonly } = useCalendar()

  const events = useFilteredEvents()

  return (
    <div className='flex flex-col gap-4 border-b p-4 lg:flex-row lg:items-center lg:justify-between'>
      <motion.div
        className='flex items-center gap-3'
        variants={slideFromLeft}
        initial='initial'
        animate='animate'
        transition={transition}
      >
        <TodayButton />
        <DateNavigator view={view} events={events} />
      </motion.div>

      <motion.div
        className='flex flex-col gap-4 lg:flex-row lg:items-center lg:gap-1.5'
        variants={slideFromRight}
        initial='initial'
        animate='animate'
        transition={transition}
      >
        <div className='options flex flex-wrap items-center gap-4 md:gap-2'>
          <FilterEvents />
          <MotionButton
            variant='outline'
            onClick={() => setView('agenda')}
            asChild
            variants={buttonHover}
            whileHover='hover'
            whileTap='tap'
          >
            <Toggle className='relative'>
              {view === 'agenda' ? (
                <>
                  <CalendarRange />
                  <span className='absolute -top-1 -right-1 size-3 rounded-full bg-green-400'></span>
                </>
              ) : (
                <LayoutList />
              )}
            </Toggle>
          </MotionButton>
          <ButtonGroup className='flex'>
            <MotionButton
              variant={view === 'day' ? 'default' : 'outline'}
              aria-label={t`按日查看`}
              onClick={() => {
                setView('day')
              }}
              variants={buttonHover}
              whileHover='hover'
              whileTap='tap'
            >
              <List className='h-4 w-4' />
            </MotionButton>

            <MotionButton
              variant={view === 'week' ? 'default' : 'outline'}
              aria-label={t`按周查看`}
              onClick={() => setView('week')}
              variants={buttonHover}
              whileHover='hover'
              whileTap='tap'
            >
              <Columns className='h-4 w-4' />
            </MotionButton>

            <MotionButton
              variant={view === 'month' ? 'default' : 'outline'}
              aria-label={t`按月查看`}
              onClick={() => setView('month')}
              variants={buttonHover}
              whileHover='hover'
              whileTap='tap'
            >
              <Grid3X3 className='h-4 w-4' />
            </MotionButton>
            <MotionButton
              variant={view === 'year' ? 'default' : 'outline'}
              aria-label={t`按年查看`}
              onClick={() => setView('year')}
              variants={buttonHover}
              whileHover='hover'
              whileTap='tap'
            >
              <Grid2X2 className='h-4 w-4' />
            </MotionButton>
          </ButtonGroup>
        </div>

        <div className='flex flex-col gap-4 lg:flex-row lg:items-center lg:gap-1.5'>
          {!readonly && <UserSelect />}

          {!readonly && (
            <AddEditEventDialog>
              <MotionButton
                variants={buttonHover}
                whileHover='hover'
                whileTap='tap'
              >
                <Plus className='h-4 w-4' />
                <Trans>添加事件</Trans>
              </MotionButton>
            </AddEditEventDialog>
          )}
        </div>
        <Settings />
      </motion.div>
    </div>
  )
}
