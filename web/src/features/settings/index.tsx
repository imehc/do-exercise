import { Outlet } from '@tanstack/react-router'
import { IconLock, IconMail, IconPalette, IconUser } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { Separator } from '~/components/ui/separator'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import SidebarNav from './components/sidebar-nav'

export default function Settings() {
  return (
    <>
      {/* ===== Top Heading ===== */}
      <Header>
        <Search />
        <div className='ml-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main fixed>
        <div className='space-y-0.5'>
          <h1 className='text-2xl font-bold tracking-tight md:text-3xl'>
            <Trans>设置</Trans>
          </h1>
          <p className='text-muted-foreground'>
            <Trans>管理您的帐户设置。</Trans>
          </p>
        </div>
        <Separator className='my-4 lg:my-6' />
        <div className='flex flex-1 flex-col space-y-2 overflow-hidden md:space-y-2 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <aside className='top-0 lg:sticky lg:w-1/5'>
            <SidebarNav items={sidebarNavItems} />
          </aside>
          <div className='flex w-full overflow-y-hidden p-1'>
            <Outlet />
          </div>
        </div>
      </Main>
    </>
  )
}

const sidebarNavItems = [
  {
    title: <Trans>个人资料</Trans>,
    icon: <IconUser size={18} />,
    href: '/settings',
  },
  {
    title: <Trans>邮箱</Trans>,
    icon: <IconMail size={18} />,
    href: '/settings/email',
  },
  {
    title: <Trans>密码</Trans>,
    icon: <IconLock size={18} />,
    href: '/settings/password',
  },
  // {
  //   title: 'Account',
  //   icon: <IconTool size={18} />,
  //   href: '/settings/account',
  // },
  {
    title: <Trans>外观</Trans>,
    icon: <IconPalette size={18} />,
    href: '/settings/appearance',
  },
  // {
  //   title: 'Notifications',
  //   icon: <IconNotification size={18} />,
  //   href: '/settings/notifications',
  // },
  // {
  //   title: 'Display',
  //   icon: <IconBrowserCheck size={18} />,
  //   href: '/settings/display',
  // },
]
