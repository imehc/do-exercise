'use client'

import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconAlertTriangle } from '@tabler/icons-react'
import { toast } from 'sonner'
import { DeleteMenuRequest, SysRole, SystemRoleApi } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Alert, AlertDescription, AlertTitle } from '~/components/ui/alert'
import { Input } from '~/components/ui/input'
import { Label } from '~/components/ui/label'
import { ConfirmDialog } from '~/components/confirm-dialog'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: SysRole
}

export function RoleDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const sysRoleApi = useApi(SystemRoleApi)
  const { isPending, mutate: handleDel } = useMutation({
    mutationFn: (values: DeleteMenuRequest) => sysRoleApi.deleteRole(values),
    onSuccess: () => {
      toast.success('删除成功')
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
          Delete Role
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to delete the role
            <span className='font-bold'>{currentRow.name}</span>?
            <br />
            This action will permanently remove the menu and its related items
            from the system. This cannot be undone.
          </p>

          <Label className='my-2'>
            Role name:
            <Input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder='Enter name name to confirm deletion.'
            />
          </Label>

          <Alert variant='destructive'>
            <AlertTitle>Warning!</AlertTitle>
            <AlertDescription>
              Please be careful, this operation cannot be undone.
            </AlertDescription>
          </Alert>
        </div>
      }
      confirmText='Delete'
      destructive
    />
  )
}
