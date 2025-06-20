'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconAlertTriangle } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  DeleteTokenRequest,
  SystemTokenApi,
  TokenInfo,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { ConfirmDialog } from '~/components/confirm-dialog'
import { ConfirmTip } from '~/components/other/confirm-tip'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: TokenInfo
}

export function TokenDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const systemTokenApi = useApi(SystemTokenApi)
  const { isPending, mutate: handleDel } = useMutation({
    mutationFn: (values: DeleteTokenRequest) =>
      systemTokenApi.deleteToken(values),
    onSuccess: () => {
      toast.success(t`删除成功`)
      onOpenChange(false, true)
    },
  })

  const [value, setValue] = useState('')

  const handleDelete = () => {
    if (value.trim() !== currentRow.accessToken) return
    handleDel({
      deleteToken: {
        accessToken: currentRow.accessToken,
      },
    })
  }

  return (
    <ConfirmDialog
      isLoading={isPending}
      open={open}
      onOpenChange={(state) => onOpenChange(state, false)}
      handleConfirm={handleDelete}
      disabled={value.trim() !== currentRow.accessToken}
      title={
        <span className='text-destructive'>
          <IconAlertTriangle
            className='stroke-destructive mr-1 inline-block'
            size={18}
          />
          <Trans>删除令牌</Trans>
        </span>
      }
      desc={
        <ConfirmTip
          text={<Trans>令牌</Trans>}
          title={currentRow.accessToken}
          value={value}
          onChange={setValue}
        />
      }
      confirmText={<Trans>删除</Trans>}
      destructive
    />
  )
}
