import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemApiApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/api'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { getColumnTitle, useColumns } from './components/api-columns'
import { ApiDialogs } from './components/api-dialogs'

export default function Api() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const enabled = useHasPermission('query')
  const sysApi = useApi(SystemApiApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findApis', pagination],
    queryFn: () => sysApi.findApis(pagination),
    enabled,
  })

  const columns = useColumns()

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
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
            stickyColumns={{
              left: ['id'],
              right: ['actions'],
            }}
          />
        </MainBody>
      </FormDialogProvider>
    </WithPermission>
  )
}
