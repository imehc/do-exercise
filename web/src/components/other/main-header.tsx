import { LanguageSwitch } from '../language-switch'
import { Header } from '../layout/header'
import { ProfileDropdown } from '../profile-dropdown'
import { Search } from '../search'
import { ThemeSwitch } from '../theme-switch'
import { useSidebar } from '../ui/sidebar'
import { RouteTab } from './route-tab'

export function MainHeader() {
  const { isMobile } = useSidebar()

  return (
    // todo: 可配置
    <Header fixed showTab={!isMobile}>
      <div className='flex h-full flex-1 flex-col gap-y-1'>
        <div className='flex'>
          <Search />
          <div className='ml-auto flex items-center space-x-4'>
            <LanguageSwitch />
            <ThemeSwitch />
            <ProfileDropdown />
          </div>
        </div>
        {!isMobile && <RouteTab />}
      </div>
    </Header>
  )
}
