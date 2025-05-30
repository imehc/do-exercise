import { useQuery } from '@tanstack/react-query'
import { FormDialogProvider } from '~/provider'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/menu-columns'
import { MenuDialogs } from './components/menu-dialogs'
import { MenuPrimaryButtons } from './components/menu-primary-buttons'
import { findMenuTree } from './data/api'

export default function Menu() {
  const { data = [], isLoading, refetch } = useQuery(findMenuTree())

  return (
    <FormDialogProvider>
      <Header fixed>
        <Search />
        <div className='ml-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='mb-2 flex flex-wrap items-center justify-between space-y-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>菜单 列表</h2>
            <p className='text-muted-foreground'>用于查看、创建、编辑菜单</p>
          </div>
          <MenuPrimaryButtons />
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <DataTable isLoading={isLoading} data={data} columns={columns} />
        </div>
      </Main>

      <MenuDialogs treeData={data} refetch={refetch} />
    </FormDialogProvider>
  )
}
