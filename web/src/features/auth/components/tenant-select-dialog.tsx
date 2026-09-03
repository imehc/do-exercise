import type { ReactNode } from 'react'
import { IconLoader3 } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
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
  /** 说明文案。默认是登录场景的措辞，找回密码等流程可覆盖 */
  description?: ReactNode
  /** 底部提示，传 null 可去掉 */
  hint?: ReactNode
  /** 处理中的状态文案 */
  pendingText?: ReactNode
}

/**
 * 租户选择框。三条流程共用：口令登录、邮箱登录、找回密码。
 * 它们的共同点是「凭据已验证，但目标租户仍不确定」——同一用户名/邮箱在多个租户下
 * 各有一个账号，服务端不擅自替用户挑一个，交由这里显式选择。
 */
export function TenantSelectDialog({
  open,
  tenants,
  isPending = false,
  onSelect,
  description,
  hint,
  pendingText,
}: TenantSelectDialogProps) {
  return (
    <Dialog open={open}>
      <DialogContent className='sm:max-w-md'>
        <DialogHeader>
          <DialogTitle>
            <Trans>选择租户</Trans>
          </DialogTitle>
          <DialogDescription>
            {description ?? (
              <Trans>该账号属于多个租户，请选择要进入的租户</Trans>
            )}
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
          <div className='text-muted-foreground flex items-center justify-center text-xs'>
            <IconLoader3 className='mr-1 h-3.5 w-3.5 animate-spin' />
            {pendingText ?? <Trans>进入中...</Trans>}
          </div>
        )}
        {hint !== null && (
          <p className='text-muted-foreground text-center text-xs'>
            {hint ?? <Trans>提示：登录后可随时在右上角菜单中切换租户</Trans>}
          </p>
        )}
      </DialogContent>
    </Dialog>
  )
}
