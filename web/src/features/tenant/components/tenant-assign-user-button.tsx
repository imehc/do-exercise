import { t } from '@lingui/core/macro'
import { IconUsers } from '@tabler/icons-react'
import { Tenant } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'

interface Props {
  tenant: Tenant
}

// 该按钮渲染在 DataTable（FormDialogProvider 内部）中，
// 因此这里可以安全地调用 useFormDialog 来打开「分配用户」对话框。
export function TenantAssignUserButton({ tenant }: Props) {
  const { setOpen, setCurrentRow } = useFormDialog<Tenant>()

  return (
    <Button
      variant='outline'
      size='icon'
      title={t`分配用户`}
      onClick={() => {
        setCurrentRow(tenant)
        setOpen('assign-users')
      }}
    >
      <IconUsers size={16} />
    </Button>
  )
}
