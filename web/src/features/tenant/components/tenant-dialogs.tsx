import { Tenant } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { TenantActionDialog } from './tenant-action-dialog'
import { TenantAssignUsersDialog } from './tenant-assign-users-dialog'
import { TenantDeleteDialog } from './tenant-delete-dialog'
import { TenantMembersDialog } from './tenant-members-dialog'

interface Props {
  refetch(): void
}
export function TenantDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog<Tenant>()

  return (
    <>
      <TenantActionDialog
        key='tenant-add'
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
          <TenantActionDialog
            key={`tenant-edit-${currentRow.tenantId}`}
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
          <TenantAssignUsersDialog
            key={`tenant-assign-users-${currentRow.tenantId}`}
            open={open === 'assign-users'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('assign-users')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <TenantMembersDialog
            key={`tenant-members-${currentRow.tenantId}`}
            open={open === 'tenant-members'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('tenant-members')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
          <TenantDeleteDialog
            key={`tenant-delete-${currentRow.tenantId}`}
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
        </>
      )}
    </>
  )
}
