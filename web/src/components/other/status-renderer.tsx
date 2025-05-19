import { PropsWithChildren } from 'react'
import { LoadingSpinner } from './loading-spinner'

interface Props extends PropsWithChildren {
  isLoading: boolean
}

export function StatusRenderer({ isLoading, children }: Props) {
  if (isLoading) {
    return <LoadingSpinner isScreen />
  }

  return children
}
