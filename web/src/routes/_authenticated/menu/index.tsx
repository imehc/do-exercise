import { createFileRoute } from '@tanstack/react-router'
import Menu from '~/features/menu'

export const Route = createFileRoute('/_authenticated/menu/')({
  component: Menu,
})
