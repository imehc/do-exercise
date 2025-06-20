import { LanguageSwitch } from '../language-switch'
import { Header } from '../layout/header'
import { ProfileDropdown } from '../profile-dropdown'
import { Search } from '../search'
import { ThemeSwitch } from '../theme-switch'

export function MainHeader() {
  return (
    <Header fixed>
      <Search />
      <div className='ml-auto flex items-center space-x-4'>
        <LanguageSwitch />
        <ThemeSwitch />
        <ProfileDropdown />
      </div>
    </Header>
  )
}
