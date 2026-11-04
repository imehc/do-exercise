import { IconTablePlus } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useAtomValue } from 'jotai'
import { originTokenAtom } from '~/atoms'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'

export function RolePrimaryButtons() {
  const { setOpen } = useFormDialog()
  // 平台超管没有租户上下文，插件回填不出 tenant_id，建出来的角色是谁都看不见的孤儿，
  // 服务端直接拒绝（platformTenantRequired）。要为某租户建角色请先切到该租户。
  const isSuperAdmin = !!useAtomValue(originTokenAtom).isSuperAdmin
  if (isSuperAdmin) return null
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
