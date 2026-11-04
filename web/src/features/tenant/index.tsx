import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemTenantApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/tenant'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { useColumns } from './components/tenant-columns'
import { TenantDialogs } from './components/tenant-dialogs'
import { TenantPrimaryButtons } from './components/tenant-primary-buttons'

export default function Tenant() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const enabled = useHasPermission('query')
  const systemTenantApi = useApi(SystemTenantApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findTenants', pagination],
    queryFn: () => systemTenantApi.findTenants(pagination),
    enabled,
  })

  const columns = useColumns(data?.meta?.total)

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainBody
          title={<Trans>租户管理</Trans>}
          subTitle={<Trans>用于查看并操作用户租户信息。</Trans>}
          element={<TenantPrimaryButtons />}
          actionElemnt={<TenantDialogs refetch={refetch} />}
        >
          <DataTable
            isLoading={isLoading}
            meta={data?.meta}
            data={data?.data || []}
            columns={columns}
            navigate={navigate}
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