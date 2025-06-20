import { IconTablePlus } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'

export function RolePrimaryButtons() {
  const { setOpen } = useFormDialog()
  return (
    <div className='flex gap-2'>
      <Button className='space-x-1' onClick={() => setOpen('add')}>
        <span>
          <Trans>添加</Trans>
        </span>
        <IconTablePlus size={18} />
      </Button>
    </div>
  )
}
