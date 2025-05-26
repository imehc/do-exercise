import { SysApi } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { ApiActionDialog } from './api-action'

interface Props {
  refetch(): void
}
export function ApiDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog<SysApi>()
  return (
    <>
      {currentRow && (
        <ApiActionDialog
          key={`api-edit-${currentRow.id}`}
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
      )}
    </>
  )
}
