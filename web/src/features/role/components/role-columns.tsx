import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { SysRole } from '~/do-exercise-api'
import { DataTableRowActions } from '~/components/other'
import {
  createActionColumn,
  createColumn,
  createDateColumn,
} from '~/components/other/data-table/column-utils'

const columnTitleMap = {
  id: (): string => t`序号`,
  name: (): string => t`角色名称`,
  code: (): string => t`角色编码`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<SysRole>[] => {
  useAtomValue(languageAtom)
  return [
    {
      accessorKey: 'id',
      header: () => columnTitleMap.id(),
      cell: ({ row, table }) => {
        const pagination = table.getState().pagination
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
      cell: ({ row }: CellContext<SysRole, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.name}</div>
      ),
    }),
    createColumn<SysRole>({
      key: 'code',
      title: columnTitleMap.code,
      cell: ({ row }: CellContext<SysRole, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.code}</div>
      ),
    }),
    createDateColumn<SysRole>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysRole>('updatedAt', columnTitleMap.updatedAt),
    createActionColumn<SysRole>(({ row }) => (
      <DataTableRowActions row={row} showEdit showDelete showInfo />
    )),
  ]
}
