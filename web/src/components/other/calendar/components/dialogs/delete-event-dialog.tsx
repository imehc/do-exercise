import { IconTrash } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '~/components/ui/alert-dialog'
import { Button } from '~/components/ui/button'
import { useCalendar } from '../../contexts/calendar-context'

interface DeleteEventDialogProps {
  eventId: number
}

export default function DeleteEventDialog({ eventId }: DeleteEventDialogProps) {
  const { removeEvent } = useCalendar()

  const deleteEvent = () => {
    try {
      removeEvent(eventId)
      toast.success(t`事件删除成功`)
    } catch {
      toast.error(t`删除事件错误`)
    }
  }

  if (!eventId) {
    return null
  }

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant='destructive'>
          <IconTrash />
          <Trans>删除</Trans>
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            <Trans>您确定要删除吗？</Trans>
          </AlertDialogTitle>
          <AlertDialogDescription>
            <Trans>
              此操作不可撤销。这将永久删除您的事件，并从我们的服务器中移除事件数据。
            </Trans>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>
            <Trans>取消</Trans>
          </AlertDialogCancel>
          <AlertDialogAction onClick={deleteEvent}>
            <Trans>继续</Trans>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
