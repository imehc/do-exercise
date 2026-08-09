import { Trans } from '@lingui/react/macro'
import { IconLockExclamation } from '@tabler/icons-react'
import { useAtomValue } from 'jotai'
import { originTokenAtom } from '~/atoms'
import { Alert, AlertDescription, AlertTitle } from '~/components/ui/alert'
import ContentSection from '../components/content-section'
import PasswordForm from './email-form'

export default function SettingsPassword() {
  const mustChangePassword = useAtomValue(originTokenAtom).mustChangePassword

  return (
    <ContentSection
      title={<Trans>密码设置</Trans>}
      desc={<Trans>在此管理您的密码设置。</Trans>}
    >
      <>
        {mustChangePassword && (
          <Alert variant='destructive' className='mb-4'>
            <IconLockExclamation />
            <AlertTitle>
              <Trans>首次登录必须修改密码</Trans>
            </AlertTitle>
            <AlertDescription>
              <Trans>
                当前使用的是系统默认密码，请立即设置新密码后才能继续使用其他功能。
              </Trans>
            </AlertDescription>
          </Alert>
        )}
        <div className='space-y-6'>
          <PasswordForm />
        </div>
      </>
    </ContentSection>
  )
}
