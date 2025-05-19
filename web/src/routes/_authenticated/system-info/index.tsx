import { createFileRoute } from '@tanstack/react-router'
import SystemInfo from '~/features/system-info'

export const Route = createFileRoute('/_authenticated/system-info/')({
  component: SystemInfo,
})
