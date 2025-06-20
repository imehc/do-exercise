import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { Alert, AlertDescription, AlertTitle } from '../ui/alert'
import { Input } from '../ui/input'
import { Label } from '../ui/label'

interface Props {
  text: React.ReactNode
  title: string
  value: string
  onChange(v: string): void
}

export function ConfirmTip({ text, title, value, onChange }: Props) {
  return (
    <div className='space-y-4'>
      <p className='mb-2'>
        <Trans>您确定要删除这个{text}吗?</Trans>
        <span className='font-bold'>{title}</span>
        <br />
        <Trans>
          此操作将永久删除{text}及其相关的系统中的项目，这是无法撤消的。
        </Trans>
      </p>

      <Label className='my-2'>
        <span className='whitespace-nowrap'>
          <Trans>{text}名称:</Trans>
        </span>
        <Input
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={t`输入名称以确认删除`}
        />
      </Label>

      <Alert variant='destructive'>
        <AlertTitle>
          <Trans>警告！</Trans>
        </AlertTitle>
        <AlertDescription>
          <Trans>请小心，此操作无法撤消。</Trans>
        </AlertDescription>
      </Alert>
    </div>
  )
}
