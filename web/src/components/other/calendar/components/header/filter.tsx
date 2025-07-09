import { Trans } from '@lingui/react/macro'
import { CheckIcon, Filter, RefreshCcw } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { Toggle } from '~/components/ui/toggle'
import { useCalendar } from '../../contexts/calendar-context'
import type { TEventColor } from '../../types'

export default function FilterEvents() {
  const { selectedColors, filterEventsBySelectedColors, clearFilter } =
    useCalendar()

  type Color = {
    label: React.ReactNode
    value: TEventColor
  }
  const colors: Color[] = [
    { label: <Trans>蓝色</Trans>, value: 'blue' },
    { label: <Trans>绿色</Trans>, value: 'green' },
    { label: <Trans>红色</Trans>, value: 'red' },
    { label: <Trans>黄色</Trans>, value: 'yellow' },
    { label: <Trans>紫色</Trans>, value: 'purple' },
    { label: <Trans>橙色</Trans>, value: 'orange' },
  ]

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Toggle variant='outline' className='w-fit cursor-pointer'>
          <Filter className='h-4 w-4' />
        </Toggle>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end' className='w-[150px]'>
        {colors.map((color, index) => (
          <DropdownMenuItem
            key={index}
            className='flex cursor-pointer items-center gap-2'
            onClick={(e) => {
              e.preventDefault()
              filterEventsBySelectedColors(color.value)
            }}
          >
            <div
              className={`size-3.5 rounded-full bg-${color.value}-600 dark:bg-${color.value}-700`}
            />
            <span className='flex items-center justify-center gap-2 capitalize'>
              {color.label}
              <span>
                {selectedColors.includes(color.value) && (
                  <span className='text-blue-500'>
                    <CheckIcon className='size-4' />
                  </span>
                )}
              </span>
            </span>
          </DropdownMenuItem>
        ))}
        <DropdownMenuItem
          disabled={selectedColors.length === 0}
          className='flex cursor-pointer gap-2'
          onClick={(e) => {
            e.preventDefault()
            clearFilter()
          }}
        >
          <RefreshCcw className='size-3.5' />
          <Trans>清除筛选</Trans>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
