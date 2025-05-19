import { useFormDialog } from '~/provider'
import { OperationLogViewBodyDialog } from './operation-log-view-body'
import { OperationLogViewInfoDialog } from './operation-log-view-info'
import { OperationLogViewMessageDialog } from './operation-log-view-message'
import { OperationLogViewParamsDialog } from './operation-log-view-params'
import { OperationLogViewResultDialog } from './operation-log-view-result'

export function OperationLogDialogs() {
  const { open, setOpen, currentRow, setCurrentRow } = useFormDialog()
  if (!currentRow) return null
  return (
    <>
      <OperationLogViewMessageDialog
        key='operation-log-view-msg'
        open={open === 'view-msg'}
        onOpenChange={() => {
          setCurrentRow(null)
          setOpen('view-msg')
        }}
      />

      <OperationLogViewParamsDialog
        key='operation-log-view-params'
        open={open === 'view-params'}
        onOpenChange={() => {
          setCurrentRow(null)
          setOpen('view-params')
        }}
      />

      <OperationLogViewBodyDialog
        key='operation-log-view-body'
        open={open === 'view-body'}
        onOpenChange={() => {
          setCurrentRow(null)
          setOpen('view-body')
        }}
      />

      <OperationLogViewResultDialog
        key='operation-log-view-result'
        open={open === 'view-result'}
        onOpenChange={() => {
          setCurrentRow(null)
          setOpen('view-result')
        }}
      />

      <OperationLogViewInfoDialog
        key='operation-log-view-info'
        open={open === 'view-info'}
        onOpenChange={() => {
          setCurrentRow(null)
          setOpen('view-info')
        }}
      />
    </>
  )
}
