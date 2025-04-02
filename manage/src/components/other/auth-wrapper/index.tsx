import { PropsWithChildren, useEffect } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { useRouter } from '~/hooks'
import NotFoundPage from '~/pages/not-found'

export function AuthWrapper({ children }: PropsWithChildren) {
  const router = useRouter()
  let login = true

  useEffect(() => {
    if (!login) {
      router.replace('/login')
    }
  }, [login, router])

  return <ErrorBoundary FallbackComponent={NotFoundPage}>{children}</ErrorBoundary>
}
