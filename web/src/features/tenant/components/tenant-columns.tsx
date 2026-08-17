import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { Tenant } from '~/do-exercise-api'
import { basicMoreOptions, usePermissions } from '~/provider'
import { DataTableRowActions } from '~/components/other'
import { Badge } from '~/components/ui/badge'
import { TenantToggleStatusSwitch } from './tenant-toggle-status-switch'
import { TenantAssignUserButton } from './tenant-assign-user-button'
import {
  createColumn,
  createDateColumn,
  createActionColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'

const columnTitleMap = {
  id: (): string => t`序号`,
  name: (): string => t`租户名称`,
  code: (): string => t`租户编码`,
  status: (): string => t`状态`,
  remark: (): string => t`备注`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (
  total?: number
): ColumnDef<DataTableFeatures, Tenant>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasMore = [...basicMoreOptions].some((p) => permissions.includes(p))
  // 仅剩一个租户时不允许删除，避免平台无可用租户
  const canDelete = permissions.some((p) => p === 'delete') && (total ?? 1) > 1
  const canAssign = permissions.some((p) => p === 'update')

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
    createColumn<Tenant>({
      key: 'name',
      title: columnTitleMap.name,
      cell: ({ row }: CellContext<DataTableFeatures, Tenant, unknown>) => (
        <div className='w-fit text-nowrap font-medium'>{row.original.name}</div>
      ),
    }),
    createColumn<Tenant>({
      key: 'code',
      title: columnTitleMap.code,
      cell: ({ row }: CellContext<DataTableFeatures, Tenant, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.code}</div>
      ),
    }),
    createColumn<Tenant>({
      key: 'status',
      title: columnTitleMap.status,
      cell: ({ row }: CellContext<DataTableFeatures, Tenant, unknown>) => {
        const { tenantId, name, remark, status } = row.original
        return permissions.some((p) => p === 'update') ? (
          <TenantToggleStatusSwitch
            tenantId={tenantId ?? ''}
            name={name ?? ''}
            remark={remark}
            status={!!status}
          />
        ) : status ? (
          <Badge>
            <Trans>启用</Trans>
          </Badge>
        ) : (
          <Badge variant='destructive'>
            <Trans>停用</Trans>
          </Badge>
        )
      },
    }),
    createColumn<Tenant>({
      key: 'remark',
      title: columnTitleMap.remark,
      cell: ({ row }: CellContext<DataTableFeatures, Tenant, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.remark || '-'}</div>
      ),
    }),
    createDateColumn<Tenant>('createdAt', columnTitleMap.createdAt),
    createDateColumn<Tenant>('updatedAt', columnTitleMap.updatedAt),
    ...[
      hasMore
        ? createActionColumn<Tenant>(({ row }) => (
            <DataTableRowActions
              row={row}
              showEdit={permissions.some((p) => p === 'update')}
              showDelete={canDelete}
            >
              {canAssign && <TenantAssignUserButton tenant={row.original} />}
            </DataTableRowActions>
          ))
        : [],
    ].flat(),
  ]
}