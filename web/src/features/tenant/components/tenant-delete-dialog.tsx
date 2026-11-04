import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconAlertTriangle } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { SystemTenantApi, Tenant } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { ConfirmDialog } from '~/components/confirm-dialog'
import { ConfirmTip } from '~/components/other/confirm-tip'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: Tenant
}

export function TenantDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const systemTenantApi = useApi(SystemTenantApi)
  const { isPending, mutate: handleDel } = useMutation({
    mutationFn: (values: { id: string }) =>
      systemTenantApi.deleteTenant(values),
    onSuccess: () => {
      toast.success(t`删除成功`)
      onOpenChange(false, true)
    },
  })

  const [value, setValue] = useState('')

  const handleDelete = () => {
    if (value.trim() !== (currentRow.name ?? '')) return
    handleDel({ id: currentRow.tenantId ?? '' })
  }

  return (
    <ConfirmDialog
      isLoading={isPending}
      open={open}
      onOpenChange={(state) => onOpenChange(state, false)}
      handleConfirm={handleDelete}
      disabled={value.trim() !== currentRow.name}
      title={
        <span className='text-destructive'>
          <IconAlertTriangle
            className='stroke-destructive mr-1 inline-block'
            size={18}
          />
          <Trans>删除租户</Trans>
        </span>
      }
      desc={
        <ConfirmTip
          text={<Trans>租户</Trans>}
          title={currentRow.name ?? ''}
          value={value}
          onChange={setValue}
        />
      }
      confirmText='Delete'
      destructive
    />
  )
}