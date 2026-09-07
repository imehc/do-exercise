import { IconTablePlus } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useFormDialog } from '~/provider'
import { usePlatformResourceReadonly } from '~/hooks/use-tenant'
import { Button } from '~/components/ui/button'

export function MenuPrimaryButtons() {
  const { setOpen } = useFormDialog()
  // sys_menu 由平台统一维护，多租户模式下租户侧只读，不提供新增入口
  const readonly = usePlatformResourceReadonly()
  if (readonly) {
    return null
  }
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
