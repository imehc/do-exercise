import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemUserApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/user'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { useColumns } from './components/user-columns'
import { UserDialogs } from './components/user-dialogs'
import { UserPrimaryButtons } from './components/user-primary-buttons'

export default function Role() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()

  const enabled = useHasPermission('query')
  const systemUserApi = useApi(SystemUserApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findUsers', pagination],
    queryFn: () => systemUserApi.findUsers(pagination),
    enabled,
  })

  const columns = useColumns()

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainHeader />
        <MainBody
          title={<Trans>用户列表</Trans>}
          subTitle={<Trans>用于查看并操作用户信息。</Trans>}
          element={<UserPrimaryButtons />}
          actionElemnt={<UserDialogs refetch={refetch} />}
        >
          <DataTable
            isLoading={isLoading}
            meta={data?.meta}
            data={data?.data || []}
            columns={columns}
            navigate={navigate}
          />
        </MainBody>
      </FormDialogProvider>
    </WithPermission>
  )
}
