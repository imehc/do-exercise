import { useQuery } from '@tanstack/react-query'
import { SystemInfoApi } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Header } from '~/components/layout/header'
import { Main } from '~/components/layout/main'
import { StatusRenderer } from '~/components/other'
import { ProfileDropdown } from '~/components/profile-dropdown'
import { Search } from '~/components/search'
import { ThemeSwitch } from '~/components/theme-switch'
import { SysInfoView } from './components/system-info-view'

export default function SystemInfo() {
  const systemInfoApi = useApi(SystemInfoApi)
  const { isLoading, data } = useQuery({
    queryKey: ['getSystemInfo'],
    queryFn: () => systemInfoApi.getSystemInfo(),
    refetchInterval: 1000 * 10,
  })

  return (
    <>
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
            <h2 className='text-2xl font-bold tracking-tight'>系统状态</h2>
            <p className='text-muted-foreground'>
              查看系统、磁盘、CPU、内存等信息
            </p>
          </div>
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <StatusRenderer isLoading={isLoading} data={data}>
            {(info) => <SysInfoView data={info} />}
          </StatusRenderer>
        </div>
      </Main>
    </>
  )
}
