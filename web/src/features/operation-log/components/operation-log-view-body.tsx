import { IconMessage } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
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

export function OperationLogViewBodyDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysOperationLog>()

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
    >
      <DialogContent>
        <DialogHeader className='text-left'>
          <DialogTitle className='flex items-center gap-2'>
            <IconMessage /> <Trans>请求体</Trans>
          </DialogTitle>
          <DialogDescription>
            <Trans>
              下面是随操作一起发送的详细请求正文。你可以检查中包含的参数、有效负载和其他数据请求。
            </Trans>
          </DialogDescription>
        </DialogHeader>

        <CodeBlock language='json' content={currentRow?.body ?? ''} />

        <DialogFooter className='gap-y-2'>
          <DialogClose asChild>
            <Button variant='outline'>
              <Trans>关闭</Trans>
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
