import * as React from 'react'
import { cn } from '~/lib/utils'
import { Avatar, AvatarFallback } from '~/components/ui/avatar'

type AvatarProps = React.ComponentProps<typeof Avatar>

interface AvatarGroupProps extends React.ComponentProps<'div'> {
  children: React.ReactElement<AvatarProps>[]
  max?: number
}

const AvatarGroup = ({
  children,
  max,
  className,
  ...props
}: AvatarGroupProps) => {
  const totalAvatars = React.Children.count(children)
  const displayedAvatars = React.Children.toArray(children)
    .slice(0, max)
    .reverse()
  const remainingAvatars = max ? Math.max(totalAvatars - max, 1) : 0

  return (
    <div
      className={cn('flex flex-row-reverse items-center', className)}
      {...props}
    >
      {remainingAvatars > 0 && (
        <Avatar className='ring-background relative -ml-2 ring-2 hover:z-10'>
          <AvatarFallback className='bg-muted-foreground text-white'>
            +{remainingAvatars}
          </AvatarFallback>
        </Avatar>
      )}
      {displayedAvatars.map((avatar, index) => {
        if (!React.isValidElement(avatar)) return null

        return (
          <div key={index} className='relative -ml-2 hover:z-10'>
            {React.cloneElement(avatar as React.ReactElement<AvatarProps>, {
              className: 'ring-2 ring-background',
            })}
          </div>
        )
      })}
    </div>
  )
}

export { AvatarGroup }
