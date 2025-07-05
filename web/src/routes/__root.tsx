import { QueryClient } from '@tanstack/react-query'
import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { type Token } from '~/do-exercise-api'
import { Toaster } from '~/components/ui/sonner'
import { NavigationProgress } from '~/components/navigation-progress'
import GeneralError from '~/features/errors/general-error'
import NotFoundError from '~/features/errors/not-found-error'

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
  token?: Token
}>()({
  component: () => {
    return (
      <>
        <NavigationProgress />
        <Outlet />
        <Toaster duration={5000} richColors />
        {import.meta.env.MODE === 'development' && (
          <>
            <ReactQueryDevtools buttonPosition='bottom-left' />
            <TanStackRouterDevtools position='bottom-right' />
          </>
        )}
      </>
    )
  },
  notFoundComponent: NotFoundError,
  errorComponent: GeneralError,
})
