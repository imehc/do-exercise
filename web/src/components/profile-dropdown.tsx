import { Link } from '@tanstack/react-router'
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
            {/* TODO: 根据前缀判断是否需要拼接完整的图片地址 */}
            <AvatarImage
              src={userProfile?.avatar}
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
              Settings
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
          Log out
          {/* <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut> */}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
