import { MenuType, SysMenuTree } from '~/do-exercise-api'
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
        onOpenChange={(_, hasRefresh) => {
          if (hasRefresh) {
            refetch()
          }
          setOpen('add')
        }}
      />
      {currentRow && (
        <>
          <MenuActionDialog
            key={`menu-add-child-${currentRow.id}`}
            treeData={treeData}
            open={open === 'add-child'}
            onOpenChange={(_, hasRefresh) => {
              if (hasRefresh) {
                refetch()
              }
              setOpen('add-child')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={
              {
                parentId: currentRow.id,
                type:
                  currentRow.type === MenuType.button // 默认为下一级菜单
                    ? currentRow.type
                    : currentRow.type + 1,
              } as SysMenuTree
            }
          />
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
        </>
      )}
      <MenuViewInfoDialog
        key='menu-view-info'
        open={open === 'view-info'}
        onOpenChange={() => {
          setOpen('view-info')
          setTimeout(() => {
            setCurrentRow(null)
          }, 500)
        }}
      />
    </>
  )
}
