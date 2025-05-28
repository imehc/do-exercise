import React from 'react'
import { AlertCircle } from 'lucide-react'
import { LoadingSpinner } from './loading-spinner'

interface Props<T> {
  isLoading: boolean
  className?: string
  isScreen?: boolean
  data?: T
  children: ((data: NonNullable<T>) => React.ReactNode) | React.ReactNode
}

export function StatusRenderer<T>({
  isLoading,
  isScreen,
  data,
  children,
  className,
}: Props<T>): React.ReactNode {
  if (isLoading) {
    return <LoadingSpinner isScreen={isScreen} className={className} />
  }

  if (typeof children === 'function') {
    if (!data) {
      return (
        <div className='text-muted-foreground flex h-24 items-center justify-center gap-2 text-sm'>
          <AlertCircle className='h-4 w-4' />
          暂无数据
        </div>
      )
    }
    return children(data)
  }

  return children
}
