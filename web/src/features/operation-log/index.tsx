import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemLogApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/operation-log'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { getColumnTitle, useColumns } from './components/operation-log-columns'
import { OperationLogDialogs } from './components/operation-log-dialogs'

export default function OperationLog() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const enabled = useHasPermission('query')
  const sysLogApi = useApi(SystemLogApi)
  const { data, isLoading } = useQuery({
    queryKey: ['findLogs', pagination],
    queryFn: () => sysLogApi.findLogs(pagination),
    enabled,
  })

  const columns = useColumns()

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainHeader />
        <MainBody
          title={<Trans>操作日志列表</Trans>}
          subTitle={
            <Trans>
              用于查看所有操作日志，包括请求信息、响应结果和操作人等。
            </Trans>
          }
          actionElemnt={<OperationLogDialogs />}
        >
          <DataTable
            isLoading={isLoading}
            navigate={navigate}
            meta={data?.meta}
            data={data?.data || []}
            columns={columns}
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
