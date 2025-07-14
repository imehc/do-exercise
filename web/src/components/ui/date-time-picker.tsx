import type { ControllerRenderProps, UseFormReturn } from 'react-hook-form'
import { IconCalendar } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { cn } from '~/lib/utils'
import { useDateFormat } from '~/hooks/use-date-locale'
import { Button } from '~/components/ui/button'
import { Calendar } from '~/components/ui/calendar'
import {
  FormControl,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '~/components/ui/popover'
import { ScrollArea, ScrollBar } from '~/components/ui/scroll-area'
import { useCalendar } from '../other/calendar/contexts/calendar-context'
import { TEventFormData } from '../other/calendar/schemas'

interface DatePickerProps {
  form: UseFormReturn<TEventFormData>
  field: ControllerRenderProps<TEventFormData, 'endDate' | 'startDate'>
}

export function DateTimePicker({ form, field }: DatePickerProps) {
  const { use24HourFormat } = useCalendar()
  const { formatDateTime } = useDateFormat()

  function handleDateSelect(date: Date | undefined) {
    if (date) {
      form.setValue(field.name, date)
    }
  }

  function handleTimeChange(type: 'hour' | 'minute' | 'ampm', value: string) {
    const currentDate = form.getValues(field.name) || new Date()
    const newDate = new Date(currentDate)

    if (type === 'hour') {
      newDate.setHours(parseInt(value, 10))
    } else if (type === 'minute') {
      newDate.setMinutes(parseInt(value, 10))
    } else if (type === 'ampm') {
      const hours = newDate.getHours()
      if (value === 'AM' && hours >= 12) {
        newDate.setHours(hours - 12)
      } else if (value === 'PM' && hours < 12) {
        newDate.setHours(hours + 12)
      }
    }

    form.setValue(field.name, newDate)
  }

  return (
    <FormItem className='flex flex-col'>
      <FormLabel>
        {field.name === 'startDate' ? (
          <Trans>开始日期</Trans>
        ) : (
          <Trans>结束日期</Trans>
        )}
      </FormLabel>
      <Popover modal={true}>
        <PopoverTrigger asChild>
          <FormControl>
            <Button
              variant={'outline'}
              className={cn(
                'w-full pl-3 text-left font-normal',
                !field.value && 'text-muted-foreground'
              )}
            >
              {field.value ? (
                formatDateTime(field.value)
              ) : (
                // format(
                //   field.value,
                //   use24HourFormat ? 'MM/dd/yyyy HH:mm' : 'MM/dd/yyyy hh:mm aa'
                // )
                <span>MM/DD/YYYY hh:mm aa</span>
              )}
              <IconCalendar className='ml-auto h-4 w-4 opacity-50' />
            </Button>
          </FormControl>
        </PopoverTrigger>
        <PopoverContent className='w-auto p-0'>
          <div className='sm:flex'>
            <Calendar
              mode='single'
              selected={field.value}
              onSelect={handleDateSelect}
              autoFocus
            />
            <div className='flex flex-col divide-y sm:h-[300px] sm:flex-row sm:divide-x sm:divide-y-0'>
              <ScrollArea className='w-64 sm:w-auto'>
                <div className='flex p-2 sm:flex-col'>
                  {Array.from(
                    { length: use24HourFormat ? 24 : 12 },
                    (_, i) => i
                  ).map((hour) => (
                    <Button
                      key={hour}
                      size='icon'
                      variant={
                        field.value &&
                        field.value.getHours() % (use24HourFormat ? 24 : 12) ===
                          hour % (use24HourFormat ? 24 : 12)
                          ? 'default'
                          : 'ghost'
                      }
                      className='aspect-square shrink-0 sm:w-full'
                      onClick={() => handleTimeChange('hour', hour.toString())}
                    >
                      {hour.toString().padStart(2, '0')}
                    </Button>
                  ))}
                </div>
                <ScrollBar orientation='horizontal' className='sm:hidden' />
              </ScrollArea>
              <ScrollArea className='w-64 sm:w-auto'>
                <div className='flex p-2 sm:flex-col'>
                  {Array.from({ length: 12 }, (_, i) => i * 5).map((minute) => (
                    <Button
                      key={minute}
                      size='icon'
                      variant={
                        field.value && field.value.getMinutes() === minute
                          ? 'default'
                          : 'ghost'
                      }
                      className='aspect-square shrink-0 sm:w-full'
                      onClick={() =>
                        handleTimeChange('minute', minute.toString())
                      }
                    >
                      {minute.toString().padStart(2, '0')}
                    </Button>
                  ))}
                </div>
                <ScrollBar orientation='horizontal' className='sm:hidden' />
              </ScrollArea>
            </div>
          </div>
        </PopoverContent>
      </Popover>
      <FormMessage />
    </FormItem>
  )
}
