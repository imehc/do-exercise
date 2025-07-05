import { isBefore } from 'date-fns'
import { createFileRoute, redirect } from '@tanstack/react-router'
import Dashboard from '~/features/dashboard'

export const Route = createFileRoute('/_authenticated/')({
  beforeLoad: ({ context: { token } }) => {
    if (
      !token ||
      !token.accessToken ||
      isBefore(token.expireTime, new Date())
    ) {
      throw redirect({
        to: '/sign-in',
        search: {
          redirect: window.location.pathname,
        },
      })
    }
  },
  component: Dashboard,
})
