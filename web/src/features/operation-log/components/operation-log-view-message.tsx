import { IconMessage } from '@tabler/icons-react'
import { SysOperationLog } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'
import { CodeBlock } from '~/components/other'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OperationLogViewMessageDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysOperationLog>()

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
    >
      <DialogContent className='sm:max-w-md'>
        <DialogHeader className='text-left'>
          <DialogTitle className='text-destructive flex items-center gap-2'>
            <IconMessage /> 错误信息
          </DialogTitle>
          <DialogDescription>
            请求返回了一个错误。响应详细信息如下所示。
          </DialogDescription>
        </DialogHeader>

        <CodeBlock content={currentRow?.message ?? ''} />

        <DialogFooter className='gap-y-2'>
          <DialogClose asChild>
            <Button variant='outline'>Close</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
