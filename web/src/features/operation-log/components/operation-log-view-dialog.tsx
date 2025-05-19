import { Eye } from 'lucide-react'
import { SysOperationLog } from '~/do-exercise-api'
import { DialogType, useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'

interface Props {
  type: DialogType
  data: SysOperationLog
}

export function OperationLogViewDialog({ type, data }: Props) {
  const { setOpen, setCurrentRow } = useFormDialog()

  let content = null

  switch (type) {
    case 'view-msg':
      content = data.message ? data : null
      break
    case 'view-params':
      content = data.params ? data : null
      break
    case 'view-body':
      content = data.body ? data : null
      break
    case 'view-result':
      content = data.result ? data : null
      break
    default:
      break
  }

  if (!content) return '-'

  return (
    <div className='w-fit text-nowrap'>
      <Button
        variant='outline'
        size='icon'
        onClick={() => {
          setCurrentRow(data)
          setOpen(type)
        }}
      >
        <Eye />
      </Button>
    </div>
  )
}
