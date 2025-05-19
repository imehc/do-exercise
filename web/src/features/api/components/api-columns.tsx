import { format } from 'date-fns'
import { ColumnDef } from '@tanstack/react-table'
import { SysApi } from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { DataTableColumnHeader, DataTableRowActions } from '~/components/other'
import { callDisabledTypes, callMethodTypes } from '../data/data'

export const columns: ColumnDef<SysApi>[] = [
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
    accessorKey: 'path',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求路径' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.path}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'description',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='描述' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.description}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'method',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求方法' />
    ),
    cell: ({ row }) => {
      const { method } = row.original
      const badgeColor = callMethodTypes.get(method)
      return (
        <div className='flex space-x-2'>
          <Badge variant='outline' className={cn('capitalize', badgeColor)}>
            {row.getValue('method')}
          </Badge>
        </div>
      )
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
    enableHiding: false,
    enableSorting: true,
  },
  {
    accessorKey: 'group',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='分组' />
    ),
    cell: ({ row }) => <div>{row.original.group}</div>,
    enableSorting: false,
  },
  {
    accessorKey: 'disabled',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='禁用状态' />
    ),
    cell: ({ row }) => {
      const { disabled } = row.original
      const badgeColor = callDisabledTypes.get(disabled)
      return (
        <div className='flex space-x-2'>
          <Badge variant='outline' className={cn('capitalize', badgeColor)}>
            {!row.getValue('disabled') ? '正常' : '已禁用'}
          </Badge>
        </div>
      )
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
    enableHiding: false,
    enableSorting: false,
  },
  {
    accessorKey: 'sort',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='序号' />
    ),
    cell: ({ row }) => <div>{row.original.sort}</div>,
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
    cell: ({ row }) => <DataTableRowActions row={row} showEdit />,
  },
]
