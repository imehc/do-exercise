import { SysMenuTree } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { MenuActionDialog } from './menu-action-dialog'
import { MenuDeleteDialog } from './menu-delete-dialog'
import { MenuViewInfoDialog } from './menu-view-info'

interface Props {
  treeData: SysMenuTree[]
  refetch(): void
}
export function MenuDialogs({ treeData, refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } =
    useFormDialog<SysMenuTree>()

  return (
    <>
      <MenuActionDialog
        key='menu-add'
        treeData={treeData}
        open={open === 'add'}
        onOpenChange={() => setOpen('add')}
      />
      {currentRow && (
        <>
          <MenuActionDialog
            key={`menu-edit-${currentRow.id}`}
            treeData={treeData}
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
          <MenuViewInfoDialog
            key='menu-view-info'
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
