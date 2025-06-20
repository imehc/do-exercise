import { useQuery } from '@tanstack/react-query'
import { Trans } from '@lingui/react/macro'
import { FormDialogProvider } from '~/provider'
import { DataTable, MainBody, MainHeader } from '~/components/other'
import { getColumnTitle, useColumns } from './components/menu-columns'
import { MenuDialogs } from './components/menu-dialogs'
import { MenuPrimaryButtons } from './components/menu-primary-buttons'
import { findMenuTree } from './data/api'

export default function Menu() {
  const { data = [], isLoading, refetch } = useQuery(findMenuTree())
  const columns = useColumns()

  return (
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
        />
      </MainBody>
    </FormDialogProvider>
  )
}
