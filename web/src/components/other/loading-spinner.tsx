import { IconLoader3, IconProps } from '@tabler/icons-react'
import { cn } from '~/lib/utils'

interface Props extends IconProps {
  isScreen?: boolean
  className?: string
}

export const LoadingSpinner = ({
  isScreen = false,
  className,
  ...props
}: Props) => {
  return (
    <div
      className={cn(
        'flex min-h-10 items-center justify-center',
        isScreen ? 'h-screen w-screen' : 'h-full w-full',
        className
      )}
    >
      <IconLoader3 className='animate-spin' {...props} />
    </div>
  )
}
