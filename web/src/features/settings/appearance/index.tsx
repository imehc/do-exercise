import { Trans } from '@lingui/react/macro'
import ContentSection from '../components/content-section'
import { AppearanceForm } from './appearance-form'

export default function SettingsAppearance() {
  return (
    <ContentSection
      title={<Trans>通用设置</Trans>}
      desc={<Trans>自定义应用程序的外观。自动切换日期 和夜晚主题</Trans>}
    >
      <AppearanceForm />
    </ContentSection>
  )
}
