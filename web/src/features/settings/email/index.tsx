import { Trans } from '@lingui/react/macro'
import ContentSection from '../components/content-section'
import EmailForm from './email-form'

export default function SettingsEmail() {
  return (
    <ContentSection
      title={<Trans>邮箱设置</Trans>}
      desc={<Trans>在此管理您的电子邮件设置。</Trans>}
    >
      <EmailForm />
    </ContentSection>
  )
}
