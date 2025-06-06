import ContentSection from '../components/content-section'
import PasswordForm from './email-form'

export default function SettingsPassword() {
  return (
    <ContentSection title='Password' desc='Manage your password settings here.'>
      <PasswordForm />
    </ContentSection>
  )
}
