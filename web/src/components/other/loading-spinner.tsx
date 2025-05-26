import { LoaderIcon, LucideProps } from 'lucide-react'
import { cn } from '~/lib/utils'

interface Props extends LucideProps {
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
      <LoaderIcon className='animate-spin' {...props} />
    </div>
  )
}
