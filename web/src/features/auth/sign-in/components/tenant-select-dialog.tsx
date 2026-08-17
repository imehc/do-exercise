import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { IconLoader3 } from '@tabler/icons-react'
import type { TenantOption } from '~/do-exercise-api'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'

interface TenantSelectDialogProps {
  open: boolean
  tenants: TenantOption[]
  isPending?: boolean
  onSelect: (tenant: TenantOption) => void
}

export function TenantSelectDialog({
  open,
  tenants,
  isPending = false,
  onSelect,
}: TenantSelectDialogProps) {
  return (
    <Dialog open={open}>
      <DialogContent className='sm:max-w-md'>
        <DialogHeader>
          <DialogTitle>
            <Trans>选择租户</Trans>
          </DialogTitle>
          <DialogDescription>
            <Trans>该账号属于多个租户，请选择要进入的租户</Trans>
          </DialogDescription>
        </DialogHeader>
        <div className='grid gap-2'>
          {tenants.map((tenant) => (
            <Button
              key={tenant.tenantId}
              variant='outline'
              className='flex h-auto items-center justify-between px-4 py-3'
              disabled={isPending}
              onClick={() => onSelect(tenant)}
            >
              <span>{tenant.name}</span>
              <span className='text-muted-foreground text-xs'>
                {tenant.code}
              </span>
            </Button>
          ))}
        </div>
        {isPending && (
          <div className='flex items-center justify-center text-muted-foreground text-xs'>
            <IconLoader3 className='mr-1 h-3.5 w-3.5 animate-spin' />
            <Trans>进入中...</Trans>
          </div>
        )}
        <p className='text-muted-foreground text-center text-xs'>
          {t`提示：登录后可随时在右上角菜单中切换租户`}
        </p>
      </DialogContent>
    </Dialog>
  )
}
