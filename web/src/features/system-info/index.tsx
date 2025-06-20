import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemInfoApi } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { MainBody, MainHeader, StatusRenderer } from '~/components/other'
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
      <MainHeader />
      <MainBody
        title={<Trans>系统信息</Trans>}
        subTitle={<Trans>查看系统、磁盘、CPU、内存等信息</Trans>}
      >
        <StatusRenderer isLoading={isLoading} data={data}>
          {(info) => <SysInfoView data={info} />}
        </StatusRenderer>
      </MainBody>
    </>
  )
}
