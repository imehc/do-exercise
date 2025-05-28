import { SysRole } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { RoleActionDialog } from './role-action-dialog'
import { RoleDeleteDialog } from './role-delete-dialog'
import { RoleViewInfoDialog } from './role-view-info'

interface Props {
  refetch(): void
}
export function RoleDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog<SysRole>()

  return (
    <>
      <RoleActionDialog
        key='role-add'
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
          <RoleActionDialog
            key={`role-edit-${currentRow.id}`}
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
          <RoleDeleteDialog
            key={`role-delete-${currentRow.id}`}
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
          <RoleViewInfoDialog
            key='role-view-info'
            open={open === 'view-info'}
            onOpenChange={() => {
              setCurrentRow(null)
              setOpen('view-info')
            }}
          />
        </>
      )}
    </>
  )
}
