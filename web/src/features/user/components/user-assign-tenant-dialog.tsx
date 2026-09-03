import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { IconLoader3, IconBuilding, IconArrowsExchange } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { SystemUserApi, SysUser } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { useApi } from '~/hooks/use-api'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'

interface Props {
  currentRow: SysUser
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function UserAssignTenantDialog({
  currentRow,
  open,
  onOpenChange,
}: Props) {
  const sysUserApi = useApi(SystemUserApi)
  const [tenantId, setTenantId] = useState('')

  const { data: tenants = [], refetch } = useQuery({
    queryKey: ['listAssignableTenants', currentRow.id],
    queryFn: () => sysUserApi.listAssignableTenants({ id: currentRow.id ?? '' }),
    enabled: open,
  })

  const { isPending, mutate } = useMutation({
    mutationFn: (id: string) =>
      sysUserApi.assignUserTenant({
        id: currentRow.id ?? '',
        assignUserTenant: { tenantId: id },
      }),
    onSuccess: () => {
      toast.success(t`分配租户成功`)
      setTenantId('')
      refetch()
      onOpenChange(false, true)
    },
    onError: () => {
      toast.error(t`分配租户失败`)
    },
  })

  const handleAssign = () => {
    if (!tenantId) {
      toast.error(t`请选择目标租户`)
      return
    }
    mutate(tenantId)
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        if (!state) {
          setTenantId('')
        }
        onOpenChange(state, false)
      }}
    >
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader className='text-left'>
          <DialogTitle className='flex items-center gap-2'>
            <IconBuilding size={18} />
            <span>
              <Trans>分配用户</Trans>
            </span>
            <span className='text-muted-foreground text-sm font-normal italic'>
              {currentRow?.username}
            </span>
            <span>
              <Trans>到租户</Trans>
            </span>
          </DialogTitle>
          <DialogDescription>
            <Trans>
              将用户复制到目标租户下，原租户记录保留。目标租户下的角色由该租户管理员后续在用户管理中分配。
            </Trans>
          </DialogDescription>
        </DialogHeader>
        <div className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1 p-0.5'>
          <span className='text-muted-foreground col-span-2 text-right text-sm'>
            <Trans>目标租户</Trans>
          </span>
          <Select value={tenantId || undefined} onValueChange={setTenantId}>
            <SelectTrigger className='col-span-8 w-full'>
              <SelectValue placeholder={t`请选择目标租户`} />
            </SelectTrigger>
            <SelectContent>
              {tenants.map((tenant) => (
                <SelectItem
                  key={tenant.tenantId}
                  value={tenant.tenantId ?? ''}
                >
                  {tenant.name}
                  {tenant.code ? `（${tenant.code}）` : ''}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <DialogFooter>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false, false)}
          >
            <Trans>取消</Trans>
          </Button>
          <Button
            type='button'
            onClick={handleAssign}
            disabled={isPending || tenants.length === 0}
          >
            {isPending ? (
              <>
                <IconLoader3 className='animate-spin' />
                <span>
                  <Trans>分配中</Trans>...
                </span>
              </>
            ) : (
              <span>
                <Trans>分配</Trans>
              </span>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export function AssignTenant({ currentRow }: Pick<Props, 'currentRow'>) {
  const { setOpen, setCurrentRow } = useFormDialog<SysUser>()
  return (
    <Button
      variant='outline'
      size='icon'
      title={t`分配租户`}
      onClick={() => {
        setOpen('assign-tenant')
        setCurrentRow(currentRow)
      }}
    >
      <IconArrowsExchange />
    </Button>
  )
}
