import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import {
  FormDialogProvider,
  useHasPermission,
  WithPermission,
} from '~/provider'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import ForbiddenError from '../errors/forbidden'
import { getColumnTitle, useColumns } from './components/menu-columns'
import { MenuDialogs } from './components/menu-dialogs'
import { MenuPrimaryButtons } from './components/menu-primary-buttons'
import { findMenuTree } from './data/api'

export default function Menu() {
  const enabled = useHasPermission('query')
  const {
    data = [],
    isLoading,
    refetch,
  } = useQuery({ ...findMenuTree(), enabled })
  const columns = useColumns()

  return (
    <WithPermission permission='query' fallback={<ForbiddenError />}>
      <FormDialogProvider>
        <MainHeader />
        <MainBody
          title={<Trans>菜单列表</Trans>}
          subTitle={<Trans>用于创建、更新、以及查看菜单信息。</Trans>}
          element={<MenuPrimaryButtons />}
          actionElemnt={<MenuDialogs treeData={data} refetch={refetch} />}
        >
          <DataTable
            isLoading={isLoading}
            data={data}
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
