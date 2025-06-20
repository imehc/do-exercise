import { Trans } from '@lingui/react/macro'
import ContentSection from '../components/content-section'
import ProfileForm from './profile-form'

export default function SettingsProfile() {
  return (
    <ContentSection
      title={<Trans>个人资料</Trans>}
      desc={<Trans>管理您的个人资料</Trans>}
    >
      <ProfileForm />
    </ContentSection>
  )
}
