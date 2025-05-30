import { useQuery } from '@tanstack/react-query'
import { SystemApiApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { Route } from '~/routes/_authenticated/api'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/api-columns'
import { ApiDialogs } from './components/api-dialogs'

export default function Api() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const sysApi = useApi(SystemApiApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findApis', pagination],
    queryFn: () => sysApi.findApis(pagination),
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
            <h2 className='text-2xl font-bold tracking-tight'>API 列表</h2>
            <p className='text-muted-foreground'>用于查看并编辑API</p>
          </div>
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <DataTable
            isLoading={isLoading}
            meta={data?.meta}
            data={data?.data || []}
            columns={columns}
            navigate={navigate}
          />
        </div>
      </Main>

      <ApiDialogs refetch={refetch} />
    </FormDialogProvider>
  )
}
