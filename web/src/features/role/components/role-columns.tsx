import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom, originTokenAtom } from '~/atoms'
import { SysRole } from '~/do-exercise-api'
import { basicMoreOptions, usePermissions } from '~/provider'
import { TENANT_ADMIN_ROLE_CODE } from '~/utils/tenant'
import { DataTableRowActions } from '~/components/other'
import {
  createActionColumn,
  createColumn,
  createDateColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'

const columnTitleMap = {
  id: (): string => t`序号`,
  name: (): string => t`角色名称`,
  code: (): string => t`角色编码`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<DataTableFeatures, SysRole>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  // 租户管理员角色由平台统一维护：租户侧只能查看（详情仍开放），
  // 改名、改菜单、删除都只有平台超级管理员可以做，服务端同样拦截
  // （tenantAdminRoleReadonly），这里只是不摆出注定失败的入口。
  const isSuperAdmin = !!useAtomValue(originTokenAtom).isSuperAdmin
  const hasMore = basicMoreOptions.some((p) => permissions.includes(p))

  return [
    {
      accessorKey: 'id',
      header: () => columnTitleMap.id(),
      cell: ({ row, table }) => {
        const pagination = table.atoms.pagination.get()
        return (
          (pagination.pageIndex ?? 0) * (pagination.pageSize ?? 0) +
          row.index +
          1
        )
      },
    },
    createColumn<SysRole>({
      key: 'name',
      title: columnTitleMap.name,
      cell: ({ row }: CellContext<DataTableFeatures, SysRole, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.name}</div>
      ),
    }),
    createColumn<SysRole>({
      key: 'code',
      title: columnTitleMap.code,
      cell: ({ row }: CellContext<DataTableFeatures, SysRole, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.code}</div>
      ),
    }),
    createDateColumn<SysRole>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysRole>('updatedAt', columnTitleMap.updatedAt),
    ...[
      hasMore
        ? createActionColumn<SysRole>(({ row }) => {
            const readonly =
              !isSuperAdmin && row.original.code === TENANT_ADMIN_ROLE_CODE
            return (
              <DataTableRowActions
                row={row}
                showEdit={!readonly && permissions.some((p) => p === 'update')}
                showDelete={
                  !readonly && permissions.some((p) => p === 'delete')
                }
                showInfo={permissions.some((p) => p === 'info')}
              />
            )
          })
        : [],
    ].flat(),
  ]
}
