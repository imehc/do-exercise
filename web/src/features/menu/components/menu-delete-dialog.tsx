'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconAlertTriangle } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  DeleteMenuRequest,
  SysMenuTree,
  SystemMenuApi,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { ConfirmDialog } from '~/components/confirm-dialog'
import { ConfirmTip } from '~/components/other/confirm-tip'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: SysMenuTree
}

export function MenuDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const sysMenuApi = useApi(SystemMenuApi)
  const { isPending, mutate: handleDel } = useMutation({
    mutationFn: (values: DeleteMenuRequest) => sysMenuApi.deleteMenu(values),
    onSuccess: () => {
      toast.success(t`删除成功`)
      onOpenChange(false, true)
    },
  })

  const [value, setValue] = useState('')

  const handleDelete = () => {
    if (value.trim() !== currentRow.name) return
    handleDel({ id: currentRow.id as number })
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
          <Trans>删除菜单</Trans>
        </span>
      }
      desc={
        <ConfirmTip
          text={<Trans>菜单</Trans>}
          title={currentRow.name}
          value={value}
          onChange={setValue}
        />
      }
      confirmText={<Trans>删除</Trans>}
      destructive
    />
  )
}
