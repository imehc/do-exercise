import { useQuery } from '@tanstack/react-query'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { SystemRoleApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/role'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { getColumnTitle, useColumns } from './components/role-columns'
import { RoleDialogs } from './components/role-dialogs'
import { RolePrimaryButtons } from './components/role-primary-buttons'

export default function Role() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const enabled = useHasPermission('query')
  const sysRoleApi = useApi(SystemRoleApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findRoles', pagination],
    queryFn: () => sysRoleApi.findRoles(pagination),
    enabled,
  })

  const columns = useColumns()

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainBody
          title={<Trans>角色列表</Trans>}
          subTitle={<Trans>用于创建、更新、以及查看角色信息。</Trans>}
          element={<RolePrimaryButtons />}
          actionElemnt={<RoleDialogs refetch={refetch} />}
        >
          <DataTable
            isLoading={isLoading}
            meta={data?.meta}
            data={data?.data || []}
            columns={columns}
            search={pagination}
            navigate={navigate}
            serchOptions={[
              {
                key: 'name',
                placeholder: t`角色名称`,
              },
              {
                key: 'code',
                placeholder: t`角色编码`,
              },
            ]}
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
