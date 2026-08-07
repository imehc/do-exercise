import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { SysApi } from '~/do-exercise-api'
import { usePermissions } from '~/provider'
import { DataTableRowActions } from '~/components/other'
import {
  createColumn,
  createActionColumn,
  createDateColumn,
  createBadgeColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'
import { callDisabledTypes, callMethodTypes } from '../data/data'

export type Row = { getValue(id: string): string }

const columnTitleMap = {
  id: (): string => t`序号`,
  path: (): string => t`请求路径`,
  description: (): string => t`描述`,
  method: (): string => t`请求方法`,
  group: (): string => t`分组`,
  disabled: (): string => t`禁用状态`,
  sort: (): string => t`排序`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
} as const

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<DataTableFeatures, SysApi>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasMore = permissions.some((p) => p === 'update')

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
    createColumn<SysApi>({
      key: 'path',
      title: columnTitleMap.path,
      cell: ({ row }: CellContext<DataTableFeatures, SysApi, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.path}</div>
      ),
    }),
    createColumn<SysApi>({
      key: 'description',
      title: columnTitleMap.description,
      cell: ({ row }: CellContext<DataTableFeatures, SysApi, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.description}</div>
      ),
    }),
    createBadgeColumn<SysApi>(
      'method',
      columnTitleMap.method,
      (value: unknown) => callMethodTypes.get(value as string),
      (value: unknown) => value as string,
      {
        filterFn: (row: Row, id: string, value: Array<string>) =>
          value.includes(row.getValue(id)),
        enableHiding: false,
        enableSorting: true,
      }
    ),
    createColumn<SysApi>({
      key: 'group',
      title: columnTitleMap.group,
    }),
    createBadgeColumn<SysApi>(
      'disabled',
      columnTitleMap.disabled,
      (value: unknown) => callDisabledTypes.get(value as boolean),
      (value: unknown) => (!value ? <Trans>正常</Trans> : <Trans>禁用</Trans>),
      {
        filterFn: (row: Row, id: string, value: Array<string>) =>
          value.includes(row.getValue(id)),
        enableHiding: false,
        enableSorting: false,
      }
    ),
    createColumn<SysApi>({
      key: 'sort',
      title: columnTitleMap.sort,
      cell: ({ row }: CellContext<DataTableFeatures, SysApi, unknown>) => (
        <div>{row.original.sort ?? 0}</div>
      ),
    }),
    createDateColumn<SysApi>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysApi>('updatedAt', columnTitleMap.updatedAt),
    ...[
      hasMore
        ? createActionColumn<SysApi>(
            ({ row }: CellContext<DataTableFeatures, SysApi, unknown>) => (
              <DataTableRowActions
                row={row}
                showEdit={permissions.some((p) => p === 'update')}
              />
            )
          )
        : [],
    ].flat(),
  ]
}
