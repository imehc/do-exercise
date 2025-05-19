import { cn } from '~/lib/utils'

type StatusProps = {
  label?: string
  color?: 'success' | 'error' | 'warning' | 'info' | 'neutral'
  className?: string
}

const colorMap = {
  success: 'bg-green-300',
  warning: 'bg-orange-300',
  info: 'bg-blue-300',
  error: 'bg-destructive/50',
  neutral: 'bg-muted',
}

export function Status({ label, color = 'neutral', className }: StatusProps) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <span className={cn('h-2.5 w-2.5 rounded-full', colorMap[color])}></span>
      <span className='text-muted-foreground text-sm'>{label}</span>
    </div>
  )
}
