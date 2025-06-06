import ContentSection from '../components/content-section'
import EmailForm from './email-form'

export default function SettingsEmail() {
  return (
    <ContentSection title='Mail' desc='Manage your email settings here.'>
      <EmailForm />
    </ContentSection>
  )
}
