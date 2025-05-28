import { format } from 'date-fns'
import { ColumnDef } from '@tanstack/react-table'
import { SysRole } from '~/do-exercise-api'
import { DataTableColumnHeader, DataTableRowActions } from '~/components/other'

export const columns: ColumnDef<SysRole>[] = [
  {
    accessorKey: 'id',
    header: '序号',
    cell: ({ row, table }) => {
      const pagination = table.getState().pagination
      return (
        (pagination.pageIndex ?? 0) * (pagination.pageSize ?? 0) + row.index + 1
      )
    },
  },
  {
    accessorKey: 'name',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='角色名称' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.name}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'code',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='角色编码' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.code}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='创建时间' />
    ),
    cell: ({ row }) => (
      <div>{format(row.original.createdAt, 'yyyy-MM-dd HH:mm:ss')}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'updatedAt',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='更新时间' />
    ),
    cell: ({ row }) => (
      <div>{format(row.original.updatedAt, 'yyyy-MM-dd HH:mm:ss')}</div>
    ),
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => (
      <DataTableRowActions row={row} showEdit showDelete showInfo />
    ),
  },
]
