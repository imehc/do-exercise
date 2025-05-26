import { SysMenuTree } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { MenuActionDialog } from './menu-action-dialog'
import { MenuDeleteDialog } from './menu-delete-dialog'

interface Props {
  refetch(): void
}
export function MenuDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } =
    useFormDialog<SysMenuTree>()

  return (
    <>
      <MenuActionDialog
        key='menu-add'
        open={open === 'add'}
        onOpenChange={() => setOpen('add')}
      />
      {currentRow && (
        <>
          <MenuActionDialog
            key={`menu-edit-${currentRow.id}`}
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
          <MenuDeleteDialog
            key={`menu-delete-${currentRow.id}`}
            open={open === 'delete'}
            onOpenChange={() => {
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
