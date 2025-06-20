import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemApiApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { Route } from '~/routes/_authenticated/api'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import { getColumnTitle, useColumns } from './components/api-columns'
import { ApiDialogs } from './components/api-dialogs'

export default function Api() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const sysApi = useApi(SystemApiApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findApis', pagination],
    queryFn: () => sysApi.findApis(pagination),
  })

  const columns = useColumns()

  return (
    <FormDialogProvider>
      <MainHeader />
      <MainBody
        title={<Trans>接口列表</Trans>}
        subTitle={<Trans>用于查看接口详情和更新接口信息。</Trans>}
        actionElemnt={<ApiDialogs refetch={refetch} />}
      >
        <DataTable
          isLoading={isLoading}
          meta={data?.meta}
          data={data?.data || []}
          columns={columns}
          navigate={navigate}
          getColumnTitle={getColumnTitle}
        />
      </MainBody>
    </FormDialogProvider>
  )
}
