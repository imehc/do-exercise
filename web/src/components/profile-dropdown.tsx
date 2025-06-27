import { Link } from '@tanstack/react-router'
import { Trans } from '@lingui/react/macro'
import { ensureHttpPrefix } from '~/utils/url'
import { useLogout, useUserProfile } from '~/hooks/use-user'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

export function ProfileDropdown() {
  const { data: userProfile, isLoading: userProfileIsLoading } =
    useUserProfile()
  const { mutate: logout, isPending: logoutIsPending } = useLogout()

  const isPending = userProfileIsLoading || logoutIsPending

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild disabled={userProfileIsLoading}>
        <Button variant='ghost' className='relative h-8 w-8 rounded-full'>
          <Avatar className='h-8 w-8'>
            <AvatarImage
              src={ensureHttpPrefix(userProfile?.avatar)}
              alt={userProfile?.username}
            />
            <AvatarFallback>ST</AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className='w-56' align='end' forceMount>
        <DropdownMenuLabel className='font-normal'>
          <div className='flex flex-col space-y-1'>
            <p className='text-sm leading-none font-medium'>
              {userProfile?.username}
            </p>
            <p className='text-muted-foreground text-xs leading-none'>
              {userProfile?.email}
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem asChild>
            <Link to='/settings'>
              <Trans>设置</Trans>
              {/* <DropdownMenuShortcut>⌘S</DropdownMenuShortcut> */}
            </Link>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          className='cursor-pointer'
          onClick={() => logout()}
          disabled={isPending}
        >
          <Trans>退出登录</Trans>
          {/* <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut> */}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
