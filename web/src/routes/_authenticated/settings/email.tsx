import { createFileRoute } from '@tanstack/react-router'
import SettingsEmail from '~/features/settings/email'

export const Route = createFileRoute('/_authenticated/settings/email')({
  component: SettingsEmail,
})
