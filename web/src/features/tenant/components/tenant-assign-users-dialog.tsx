import { useMemo, useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { IconLoader3, IconUsers } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { SystemTenantApi, Tenant } from '~/do-exercise-api'
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
import { TransferList } from '~/components/other/transfer-list'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: Tenant
}

type Item = {
  key: number
  label: string
  selected?: boolean
  disabled?: boolean
}

export function TenantAssignUsersDialog({
  open,
  onOpenChange,
  currentRow,
}: Props) {
  const systemTenantApi = useApi(SystemTenantApi)
  const [selectedKeys, setSelectedKeys] = useState<number[]>([])

  const { data: assignable = [], refetch } = useQuery({
    queryKey: ['listAssignableUsers', currentRow.tenantId],
    queryFn: () =>
      systemTenantApi.listAssignableUsers({ id: currentRow.tenantId ?? '' }),
    enabled: open,
  })

  // index -> userId 映射，TransferList 的 key 为 number，用户 ID 为字符串
  const idByKey = useMemo(() => {
    const map = new Map<number, string>()
    assignable.forEach((user, index) => map.set(index, user.id ?? ''))
    return map
  }, [assignable])

  const items = useMemo<Item[]>(
    () =>
      assignable.map((user, index) => ({
        key: index,
        label: `${user.username ?? ''}${
          user.nickname ? `（${user.nickname}）` : ''
        }`,
        selected: false,
      })),
    [assignable]
  )

  const { isPending, mutate } = useMutation({
    mutationFn: (userIds: string[]) =>
      systemTenantApi.assignTenantUsers({
        id: currentRow.tenantId ?? '',
        assignTenantUsers: { userIds },
      }),
    onSuccess: () => {
      toast.success(t`分配成功`)
      setSelectedKeys([])
      refetch()
      onOpenChange(false, true)
    },
    onError: () => {
      toast.error(t`分配失败`)
    },
  })

  const handleAssign = () => {
    const userIds = selectedKeys
      .map((key) => idByKey.get(key))
      .filter((id): id is string => !!id)
    if (userIds.length === 0) {
      toast.error(t`请先选择要分配的用户`)
      return
    }
    mutate(userIds)
  }

  const tenantName = currentRow.name ?? ''

  const handleDialogChange = (state: boolean) => {
    if (!state) {
      setSelectedKeys([])
    }
    onOpenChange(state, false)
  }

  return (
    <Dialog open={open} onOpenChange={handleDialogChange}>
      <DialogContent className='sm:max-w-3xl'>
        <DialogHeader className='text-left'>
          <DialogTitle className='flex items-center gap-2'>
            <IconUsers size={18} />
            <Trans>分配用户</Trans>
          </DialogTitle>
          <DialogDescription>
            <Trans>
              为租户「{tenantName}」分配现有用户，平台超级管理员与已归属该租户的用户不可选择。
            </Trans>
          </DialogDescription>
        </DialogHeader>
        {assignable.length === 0 ? (
          <div className='text-muted-foreground flex h-40 items-center justify-center text-sm'>
            <Trans>暂无可分配的用户</Trans>
          </div>
        ) : (
          <TransferList
            data={items}
            value={[]}
            onChange={setSelectedKeys}
            className='h-80'
          />
        )}
        <DialogFooter>
          <Button
            type='button'
            variant='outline'
            onClick={() => handleDialogChange(false)}
          >
            <Trans>取消</Trans>
          </Button>
          <Button
            type='button'
            onClick={handleAssign}
            disabled={isPending || assignable.length === 0}
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
