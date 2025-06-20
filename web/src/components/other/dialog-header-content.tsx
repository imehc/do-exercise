import { Trans } from '@lingui/react/macro'
import { DialogDescription, DialogTitle } from '../ui/dialog'

interface Props {
  isEdit: boolean
  text: React.ReactNode
}

export function DialogHeaderContent({ isEdit, text }: Props) {
  return (
    <>
      <DialogTitle>
        {isEdit ? <Trans>修改{text}</Trans> : <Trans>创建{text}</Trans>}
      </DialogTitle>
      <DialogDescription>
        {isEdit ? (
          <Trans>更新{text}相关信息。</Trans>
        ) : (
          <Trans>创建{text}相关信息。</Trans>
        )}
        <Trans>完成后点击保存。</Trans>
      </DialogDescription>
    </>
  )
}
