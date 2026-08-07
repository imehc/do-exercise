import { ReactNode } from 'react'
import { format, addMinutes, set } from 'date-fns'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { Button } from '~/components/ui/button'
import { DateTimePicker } from '~/components/ui/date-time-picker'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogTrigger,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '~/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { Textarea } from '~/components/ui/textarea'
import { COLORS } from '../../constants'
import { useCalendar } from '../../contexts/calendar-context'
import { useDisclosure } from '../../hooks'
import { IEvent } from '../../interfaces'
import { getEventSchema, TEventFormData } from '../../schemas'

interface IProps {
  children: ReactNode
  startDate?: Date
  startTime?: { hour: number; minute: number }
  event?: IEvent
}

function generateId() {
  return Math.floor(Math.random() * 1000000)
}

export function AddEditEventDialog({
  children,
  startDate,
  startTime,
  event,
}: IProps) {
  const { isOpen, onClose, onToggle } = useDisclosure()
  const { addEvent, updateEvent } = useCalendar()
  const isEditing = !!event

  const getInitialDates = () => {
    if (!startDate)
      return { startDate: new Date(), endDate: addMinutes(new Date(), 30) }
    const start = startTime
      ? set(new Date(startDate), {
          hours: startTime.hour,
          minutes: startTime.minute,
          seconds: 0,
        })
      : new Date(startDate)
    const end = addMinutes(start, 30)
    return { startDate: start, endDate: end }
  }

  const initialDates = getInitialDates()

  const parseEventDates = () => {
    if (!event) return null

    return {
      startDate: new Date(event.startDate),
      endDate: new Date(event.endDate),
    }
  }

  const eventDates = parseEventDates()

  const form = useForm<TEventFormData>({
    resolver: zodResolver(getEventSchema()),
    defaultValues: isEditing
      ? {
          title: event.title,
          description: event.description,
          startDate: eventDates?.startDate,
          endDate: eventDates?.endDate,
          color: event.color,
        }
      : {
          title: '',
          description: '',
          startDate: initialDates.startDate,
          endDate: initialDates.endDate,
          color: 'blue' as const,
        },
  })

  const handleDialogChange = (open: boolean) => {
    if (!open) {
      // 清除表单验证错误
      form.clearErrors()
      // 重置表单到默认值
      form.reset()
    }
    onToggle()
  }

  const onSubmit = (values: TEventFormData) => {
    try {
      // Format event data for API
      const formattedEvent: IEvent = {
        ...values,
        startDate: format(values.startDate, "yyyy-MM-dd'T'HH:mm:ss"),
        endDate: format(values.endDate, "yyyy-MM-dd'T'HH:mm:ss"),
        id: isEditing ? event.id : generateId(),
        user: isEditing
          ? event.user
          : {
              id: generateId().toString(),
              name: 'Jeraidi Yassir',
              picturePath: null,
            },
        color: values.color,
      }

      if (isEditing) {
        updateEvent(formattedEvent)
        toast.success(t`事件更新成功`)
      } else {
        addEvent(formattedEvent)
        toast.success(t`事件创建成功`)
      }

      onClose()
      form.reset()
    } catch (error) {
      console.error(`Error ${isEditing ? 'editing' : 'adding'} event:`, error)
      toast.error(isEditing ? t`编辑事件失败` : t`添加事件失败`)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={handleDialogChange} modal={false}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {isEditing ? <Trans>编辑事件</Trans> : <Trans>添加新事件</Trans>}
          </DialogTitle>
          <DialogDescription>
            {isEditing ? (
              <Trans>修改您的现有事件。</Trans>
            ) : (
              <Trans>为您的日历创建新事件。</Trans>
            )}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form
            id='event-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='grid gap-4 py-4'
          >
            <FormField
              control={form.control}
              name='title'
              render={({ field, fieldState }) => (
                <FormItem>
                  <FormLabel htmlFor='title' className='required'>
                    <Trans>标题</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      id='title'
                      placeholder={t`请输入标题`}
                      {...field}
                      className={fieldState.invalid ? 'border-red-500' : ''}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='startDate'
              render={({ field }) => (
                <DateTimePicker form={form} field={field} />
              )}
            />
            <FormField
              control={form.control}
              name='endDate'
              render={({ field }) => (
                <DateTimePicker form={form} field={field} />
              )}
            />
            <FormField
              control={form.control}
              name='color'
              render={({ field, fieldState }) => (
                <FormItem>
                  <FormLabel className='required'>
                    <Trans>颜色</Trans>
                  </FormLabel>
                  <FormControl>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <SelectTrigger
                        className={`w-full ${fieldState.invalid ? 'border-red-500' : ''}`}
                      >
                        <SelectValue placeholder={t`选择一个颜色`} />
                      </SelectTrigger>
                      <SelectContent>
                        {COLORS.map((color) => (
                          <SelectItem value={color} key={color}>
                            <div className='flex items-center gap-2'>
                              <div
                                className={`size-3.5 rounded-full bg-${color}-600 dark:bg-${color}-700`}
                              />
                              {color}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='description'
              render={({ field, fieldState }) => (
                <FormItem>
                  <FormLabel className='required'>
                    <Trans>描述</Trans>
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      {...field}
                      placeholder={t`请输入描述`}
                      className={fieldState.invalid ? 'border-red-500' : ''}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <DialogFooter>
          <DialogClose asChild>
            <Button type='button' variant='outline'>
              <Trans>取消</Trans>
            </Button>
          </DialogClose>
          <Button form='event-form' type='submit'>
            {isEditing ? <Trans>保存更改</Trans> : <Trans>创建事件</Trans>}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
