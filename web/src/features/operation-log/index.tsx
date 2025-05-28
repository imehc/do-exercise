import { useQuery } from '@tanstack/react-query'
import { SystemLogApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { Route } from '~/routes/_authenticated/operation-log'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { DataTable, StatusRenderer } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { columns } from './components/operation-log-columns'
import { OperationLogDialogs } from './components/operation-log-dialogs'

export default function OperationLog() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const sysLogApi = useApi(SystemLogApi)
  const { data, isLoading } = useQuery({
    queryKey: ['findLogs', pagination],
    queryFn: () => sysLogApi.findLogs(pagination),
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
              <h2 className='text-2xl font-bold tracking-tight'>
                操作日志列表
              </h2>
              <p className='text-muted-foreground'>
                查看和管理所有操作日志，包括请求信息、响应结果和操作人等。
              </p>
            </div>
          </div>
          <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
            <DataTable
              navigate={navigate}
              meta={data?.meta}
              data={data?.data || []}
              columns={columns}
            />
          </div>
        </Main>
        <OperationLogDialogs />
      </StatusRenderer>
    </FormDialogProvider>
  )
}
