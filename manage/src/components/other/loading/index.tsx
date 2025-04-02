import clsx from 'clsx'

interface LoadingProps {
  global?: boolean
}

export function Loading({ global = false }: LoadingProps) {
  return (
    <div
      className={clsx(
        'flex flex-col justify-center items-center gap-6 bg-background/80 backdrop-blur-sm',
        {
          'fixed left-0 top-0 h-screen w-screen': global
        }
      )}
    >
      <div className="relative">
        <div className="h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent shadow-lg" />
        <div className="absolute -top-2 -left-2 h-16 w-16 animate-ping rounded-full border-2 border-primary opacity-20" />
      </div>
      <div className="relative">
        <p className="text-lg font-medium text-primary animate-pulse">加载中...</p>
        <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-3/4 h-1 bg-gradient-to-r from-transparent via-primary/30 to-transparent rounded-full animate-pulse" />
      </div>
    </div>
  )
}
