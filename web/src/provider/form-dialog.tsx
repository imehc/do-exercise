import React, { useState } from 'react'
import useDialogState from '~/hooks/use-dialog-state'

export type DialogType =
  | 'invite'
  | 'add'
  | 'edit'
  | 'delete'
  | 'add-child'
  | 'view-msg'
  | 'view-params'
  | 'view-body'
  | 'view-result'
  | 'view-info'
  | 'reset'
  | 'assign-users'
  | 'assign-tenant'
  | 'tenant-members'

type Object = Record<string, unknown>

interface FormDialogContextType<T = unknown> {
  open: DialogType | null
  setOpen: (str: DialogType | null) => void
  currentRow: T | null
  setCurrentRow: React.Dispatch<React.SetStateAction<T | null>>
}

const FormDialogContext =
  React.createContext<FormDialogContextType<Object> | null>(null)

interface Props {
  children: React.ReactNode
}

export function FormDialogProvider({ children }: Props) {
  const [open, setOpen] = useDialogState<DialogType>(null)
  const [currentRow, setCurrentRow] = useState<Object | null>(null)

  return (
    <FormDialogContext
      value={{
        open,
        setOpen,
        currentRow,
        setCurrentRow,
      }}
    >
      {children}
    </FormDialogContext>
  )
}

export function useFormDialog<T>() {
  const usersContext = React.useContext(FormDialogContext)

  if (!usersContext) {
    throw new Error('useUsers has to be used within <UsersContext>')
  }

  return usersContext as FormDialogContextType<T>
}
