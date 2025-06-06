import { createFileRoute } from '@tanstack/react-router'
import SettingsPassword from '~/features/settings/password'

export const Route = createFileRoute('/_authenticated/settings/password')({
  component: SettingsPassword,
})
