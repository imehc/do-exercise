import { useQuery } from '@tanstack/react-query'
import { SystemUserApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { Route } from '~/routes/_authenticated/user'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable, StatusRenderer } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/user-columns'
import { UserDialogs } from './components/user-dialogs'
import { UserPrimaryButtons } from './components/user-primary-buttons'

export default function Role() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()
  const systemUserApi = useApi(SystemUserApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findUsers', pagination],
    queryFn: () => systemUserApi.findUsers(pagination),
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
      <StatusRenderer isLoading={isLoading}>
        <Main>
          <div className='mb-2 flex flex-wrap items-center justify-between space-y-2'>
            <div>
              <h2 className='text-2xl font-bold tracking-tight'>用户 列表</h2>
              <p className='text-muted-foreground'>用于查看并编辑用户</p>
            </div>
            <UserPrimaryButtons />
          </div>
          <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
            <DataTable
              meta={data?.meta}
              data={data?.data || []}
              columns={columns}
              navigate={navigate}
            />
          </div>
        </Main>
        <UserDialogs refetch={refetch} />
      </StatusRenderer>
    </FormDialogProvider>
  )
}
