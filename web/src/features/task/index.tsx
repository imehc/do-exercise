import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { SystemJobApi } from '~/do-exercise-api'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { Route } from '~/routes/_authenticated/task'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { useColumns } from './components/task-columns'
import { TaskDialogs } from './components/task-dialogs'
import { TaskPrimaryButtons } from './components/task-primary-buttons'

export default function Task() {
  const navigate = Route.useNavigate()
  const pagination = Route.useSearch()
  const enabled = useHasPermission('query')
  const systemJobApi = useApi(SystemJobApi)
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['findSysJobs', pagination],
    queryFn: () => systemJobApi.findSysJobs(pagination),
    enabled,
  })

  const columns = useColumns(refetch)

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainHeader />
        <MainBody
          title={<Trans>定时任务</Trans>}
          subTitle={<Trans>用于查看并操作定时任务。</Trans>}
          element={<TaskPrimaryButtons />}
          actionElemnt={<TaskDialogs refetch={refetch} />}
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
