import { useQuery } from '@tanstack/react-query'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { SystemTokenApi } from '~/do-exercise-api'
import { FormDialogProvider } from '~/provider'
import { useApi } from '~/hooks/use-api'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import { getColumnTitle, useColumns } from './components/token-columns'
import { TokenDialogs } from './components/token-dialogs'

export default function Token() {
  const sysTokenApi = useApi(SystemTokenApi)
  const {
    data = [],
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['findAllToken'],
    queryFn: () => sysTokenApi.findAllToken(),
  })

  const columns = useColumns()

  return (
    <FormDialogProvider>
      <MainHeader />
      <MainBody
        title={<Trans>令牌列表</Trans>}
        subTitle={<Trans>用于查看并操作令牌。</Trans>}
        actionElemnt={<TokenDialogs refetch={refetch} />}
      >
        <DataTable
          isLoading={isLoading}
          data={data}
          columns={columns}
          enableClientPagination
          getColumnTitle={getColumnTitle}
          serchOptions={[
            { key: 'username', placeholder: t`用户名` },
            { key: 'accessToken', placeholder: t`令牌` },
          ]}
        />
      </MainBody>
    </FormDialogProvider>
  )
}
