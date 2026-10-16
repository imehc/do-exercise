import { SysUser } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { UserActionDialog } from './user-action-dialog'
import { UserAssignTenantDialog } from './user-assign-tenant-dialog'
import { UserDeleteDialog } from './user-delete-dialog'
import { UserResetPasswordDialog } from './user-reset-password-dialog'
import { UserViewInfoDialog } from './user-view-info'

interface Props {
  refetch(): void
}
export function UserDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog<SysUser>()

  return (
    <>
      <UserActionDialog
        key='user-add'
        open={open === 'add'}
        onOpenChange={(_, hasRefresh) => {
          if (hasRefresh) {
            refetch()
          }
          setOpen('add')
        }}
      />

      {currentRow && (
        <>
          <UserActionDialog
            key={`user-edit-${currentRow.id}`}
            open={open === 'edit'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('edit')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <UserResetPasswordDialog
            key='user-reset-password'
            open={open === 'reset'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('reset')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <UserAssignTenantDialog
            key={`user-assign-tenant-${currentRow.id}`}
            open={open === 'assign-tenant'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('assign-tenant')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <UserDeleteDialog
            key={`user-delete-${currentRow.id}`}
            open={open === 'delete'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('delete')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <UserViewInfoDialog
            key='user-view-info'
            open={open === 'view-info'}
            onOpenChange={() => {
              setOpen('view-info')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
          />
        </>
      )}
    </>
  )
}
