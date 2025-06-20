'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconAlertTriangle } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { DeleteUserRequest, SystemUserApi, SysUser } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { ConfirmDialog } from '~/components/confirm-dialog'
import { ConfirmTip } from '~/components/other/confirm-tip'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: SysUser
}

export function UserDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const systemUserApi = useApi(SystemUserApi)
  const { isPending, mutate: handleDel } = useMutation({
    mutationFn: (values: DeleteUserRequest) => systemUserApi.deleteUser(values),
    onSuccess: () => {
      toast.success(t`删除成功`)
      onOpenChange(false, true)
    },
  })

  const [value, setValue] = useState('')

  const handleDelete = () => {
    if (value.trim() !== currentRow.username) return
    handleDel({ id: currentRow.id })
  }

  return (
    <ConfirmDialog
      isLoading={isPending}
      open={open}
      onOpenChange={(state) => onOpenChange(state, false)}
      handleConfirm={handleDelete}
      disabled={value.trim() !== currentRow.username}
      title={
        <span className='text-destructive'>
          <IconAlertTriangle
            className='stroke-destructive mr-1 inline-block'
            size={18}
          />
          <Trans>删除用户</Trans>
        </span>
      }
      desc={
        <ConfirmTip
          text={<Trans>用户</Trans>}
          title={currentRow.username}
          value={value}
          onChange={setValue}
        />
      }
      confirmText='Delete'
      destructive
    />
  )
}
