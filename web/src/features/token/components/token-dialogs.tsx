import { TokenInfo } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { TokenDeleteDialog } from './token-delete-dialog'
import { TokenViewInfoDialog } from './token-view-info'

interface Props {
  refetch(): void
}
export function TokenDialogs({ refetch }: Props) {
  const { open, setOpen, currentRow, setCurrentRow } =
    useFormDialog<TokenInfo>()

  return (
    <>
      {currentRow && (
        <>
          <TokenDeleteDialog
            key={`token-delete-${currentRow.userId}`}
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
          <TokenViewInfoDialog
            key='token-view-info'
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
