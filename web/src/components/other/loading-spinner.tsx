import clsx from 'clsx'
import { LoaderIcon, LucideProps } from 'lucide-react'

interface Props extends LucideProps {
  isScreen?: boolean
}

export const LoadingSpinner = ({ isScreen = false, ...props }: Props) => {
  return (
    <div
      className={clsx('flex items-center justify-center', [
        isScreen ? 'h-screen w-screen' : 'h-full w-full',
      ])}
    >
      <LoaderIcon className='animate-spin' {...props} />
    </div>
  )
}
