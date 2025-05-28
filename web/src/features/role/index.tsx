import { useQuery } from '@tanstack/react-query'
import { SystemRoleApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { Route } from '~/routes/_authenticated/role'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable, StatusRenderer } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/role-columns'
import { RoleDialogs } from './components/role-dialogs'
import { RolePrimaryButtons } from './components/role-primary-buttons'

export default function Role() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()
  const sysRoleApi = useApi(SystemRoleApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findRoles', pagination],
    queryFn: () => sysRoleApi.findRoles(pagination),
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
              <h2 className='text-2xl font-bold tracking-tight'>角色 列表</h2>
              <p className='text-muted-foreground'>用于查看并编辑角色</p>
            </div>
            <RolePrimaryButtons />
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
        <RoleDialogs refetch={refetch} />
      </StatusRenderer>
    </FormDialogProvider>
  )
}
