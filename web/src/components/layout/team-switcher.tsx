import { IconAlphabetHebrew } from '@tabler/icons-react'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '~/components/ui/sidebar'
import pkg from '../../../package.json'
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar'

export interface TeamSwitcherProps {
  title?: string
  subTitle?: string
  logo?: string
}
export function TeamSwitcher({
  title = pkg.name,
  subTitle = pkg.description,
  logo,
}: TeamSwitcherProps) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size='lg'
          className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground pointer-events-none'
        >
          <div className='flex aspect-square size-8 items-center justify-center rounded-lg'>
            {logo ? (
              <Avatar className='size-full rounded-md'>
                <AvatarImage src={logo} alt='@shadcn' />
                <AvatarFallback>SW</AvatarFallback>
              </Avatar>
            ) : (
              <IconAlphabetHebrew className='bg-sidebar-primary text-sidebar-primary-foreground size-full rounded-md' />
            )}
          </div>
          <div className='grid flex-1 text-left text-sm leading-tight'>
            <span className='truncate font-semibold'>{title}</span>
            <span className='truncate text-xs'>{subTitle}</span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
