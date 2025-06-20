import { Trans } from '@lingui/react/macro'
import ContentSection from '../components/content-section'
import PasswordForm from './email-form'

export default function SettingsPassword() {
  return (
    <ContentSection
      title={<Trans>密码设置</Trans>}
      desc={<Trans>在此管理您的密码设置。</Trans>}
    >
      <PasswordForm />
    </ContentSection>
  )
}
