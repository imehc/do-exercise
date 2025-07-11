import React from 'react'
import { cn } from '~/lib/utils'

interface MainProps extends React.HTMLAttributes<HTMLElement> {
  fixed?: boolean
  ref?: React.Ref<HTMLElement>
}

export const Main = ({ fixed, ...props }: MainProps) => {
  return (
    <main
      className={cn(
        'overflow-y-auto peer-[.header-fixed]/header:mt-16 peer-[.header-with-tabs]/header:mt-24',
        'px-4 py-0',
        fixed && 'fixed-main flex grow flex-col overflow-hidden'
      )}
      {...props}
    />
  )
}

Main.displayName = 'Main'
