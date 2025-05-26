import { PropsWithChildren } from 'react'
import { LoadingSpinner } from './loading-spinner'

interface Props extends PropsWithChildren {
  isLoading: boolean
  className?: string
  isScreen?: boolean
}

export function StatusRenderer({
  isLoading,
  isScreen,
  children,
  className,
}: Props) {
  if (isLoading) {
    return <LoadingSpinner isScreen={isScreen} className={className} />
  }

  return children
}
