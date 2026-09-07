import { IconUsersGroup } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Tenant } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'

interface Props {
  tenant: Tenant
}

// 该按钮渲染在 DataTable（FormDialogProvider 内部），
// 因此可以安全地调用 useFormDialog 打开「成员管理」对话框。
export function TenantMembersButton({ tenant }: Props) {
  const { setOpen, setCurrentRow } = useFormDialog<Tenant>()

  return (
    <Button
      variant='outline'
      size='icon'
      title={t`成员管理`}
      onClick={() => {
        setCurrentRow(tenant)
        setOpen('tenant-members')
      }}
    >
      <IconUsersGroup size={16} />
    </Button>
  )
}
