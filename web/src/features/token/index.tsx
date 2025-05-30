import { useQuery } from '@tanstack/react-query'
import { SystemTokenApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/token-columns'
import { TokenDialogs } from './components/token-dialogs'

export default function Token() {
  const sysTokenApi = useApi(SystemTokenApi)
  const {
    data = [],
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['findAllToken'],
    queryFn: () => sysTokenApi.findAllToken(),
  })

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
            <h2 className='text-2xl font-bold tracking-tight'>令牌 列表</h2>
            <p className='text-muted-foreground'>用于查看并操作令牌</p>
          </div>
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <DataTable
            isLoading={isLoading}
            data={data}
            columns={columns}
            enableClientPagination
            serchOptions={[
              { key: 'username', placeholder: '用户名' },
              { key: 'accessToken', placeholder: '令牌' },
            ]}
          />
        </div>
      </Main>
      <TokenDialogs refetch={refetch} />
    </FormDialogProvider>
  )
}
