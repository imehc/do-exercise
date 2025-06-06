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
          setOpen('view-msg')
          setTimeout(() => {
            setCurrentRow(null)
          }, 500)
        }}
      />

      <OperationLogViewParamsDialog
        key='operation-log-view-params'
        open={open === 'view-params'}
        onOpenChange={() => {
          setOpen('view-params')
          setTimeout(() => {
            setCurrentRow(null)
          }, 500)
        }}
      />

      <OperationLogViewBodyDialog
        key='operation-log-view-body'
        open={open === 'view-body'}
        onOpenChange={() => {
          setOpen('view-body')
          setTimeout(() => {
            setCurrentRow(null)
          }, 500)
        }}
      />

      <OperationLogViewResultDialog
        key='operation-log-view-result'
        open={open === 'view-result'}
        onOpenChange={() => {
          setOpen('view-result')
          setTimeout(() => {
            setCurrentRow(null)
          }, 500)
        }}
      />

      <OperationLogViewInfoDialog
        key='operation-log-view-info'
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
