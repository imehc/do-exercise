import { IconCheck, IconCopy } from '@tabler/icons-react'
import { cn } from '~/lib/utils'
import { useCopyToClipboard } from '~/hooks/use-copy-to-clipboard'
import { Button } from '~/components/ui/button'

interface InlineCopyProps extends Pick<
  React.HTMLAttributes<HTMLDivElement>,
  'className'
> {
  text: string
}

export function InlineCopy({ text, className }: InlineCopyProps) {
  const { copied, copy } = useCopyToClipboard()

  return (
    <div className={cn('inline-flex items-center gap-2', className)}>
      <span className='text-sm'>{text}</span>
      <Button
        variant='ghost'
        size='icon'
        className='h-6 w-6'
        onClick={() => copy(text)}
      >
        {copied ? (
          <IconCheck className='h-4 w-4 text-green-500' />
        ) : (
          <IconCopy className='h-4 w-4' />
        )}
      </Button>
    </div>
  )
}
