import { SysJob } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { TaskActionDialog } from './task-action-dialog'
import { TaskDeleteDialog } from './task-delete-dialog'
import { TaskViewInfoDialog } from './task-view-info'

interface Props {
  refetch(): void
}
export function TaskDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog<SysJob>()

  return (
    <>
      <TaskActionDialog
        key='task-add'
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
          <TaskActionDialog
            key={`task-edit-${currentRow.id}`}
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
          <TaskDeleteDialog
            key={`task-delete-${currentRow.id}`}
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
          <TaskViewInfoDialog
            key='task-view-info'
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
