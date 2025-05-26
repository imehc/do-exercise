'use client'

import { useState } from 'react'
import { IconAlertTriangle } from '@tabler/icons-react'
import { SysMenuTree } from '~/do-exercise-api'
import { showSubmittedData } from '~/utils/show-submitted-data'
import { Alert, AlertDescription, AlertTitle } from '~/components/ui/alert'
import { Input } from '~/components/ui/input'
import { Label } from '~/components/ui/label'
import { ConfirmDialog } from '~/components/confirm-dialog'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: SysMenuTree
}

export function MenuDeleteDialog({ open, onOpenChange, currentRow }: Props) {
  const [value, setValue] = useState('')

  const handleDelete = () => {
    if (value.trim() !== currentRow.name) return

    onOpenChange(false)
    showSubmittedData(currentRow, 'The following menu has been deleted:')
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={handleDelete}
      disabled={value.trim() !== currentRow.name}
      title={
        <span className='text-destructive'>
          <IconAlertTriangle
            className='stroke-destructive mr-1 inline-block'
            size={18}
          />{' '}
          Delete Menu
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to delete the menu{' '}
            <span className='font-bold'>{currentRow.name}</span>?
            <br />
            This action will permanently remove the menu and its related items
            from the system. This cannot be undone.
          </p>

          <Label className='my-2'>
            Menu name:
            <Input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder='Enter menu name to confirm deletion.'
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
