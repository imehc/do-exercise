import { format } from 'date-fns'
import { ColumnDef } from '@tanstack/react-table'
import { SysOperationLog } from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import {
  DataTableColumnHeader,
  DataTableRowActions,
  Status,
} from '~/components/other'
import { callMethodTypes } from '~/features/api/data/data'
import { callCodeTypes } from '../data/data'
import { OperationLogViewDialog } from './operation-log-view-dialog'

export const columns: ColumnDef<SysOperationLog>[] = [
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
    accessorKey: 'username',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求用户' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.username || '-'}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'ip',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='IP地址' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.ip}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'address',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='IP属地' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>
        {row.original.address
          ?.split('|')
          .filter((part) => part && !/^\d+$/.test(part))
          .slice(0, -1)
          .join(' ') || '-'}
      </div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'os',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='系统' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.os}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'browser',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='浏览器' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.browser}</div>
    ),
    enableSorting: false,
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
    enableSorting: true,
  },
  {
    accessorKey: 'success',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='是否成功' />
    ),
    cell: ({ row }) => {
      const { success } = row.original
      return (
        <div className='flex space-x-2'>
          <Status
            color={success ? 'success' : 'error'}
            label={success ? '成功' : '失败'}
          />
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
    accessorKey: 'code',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='状态码' />
    ),
    cell: ({ row }) => {
      const { code } = row.original
      const badgeColor = callCodeTypes.get(code)
      return (
        <div className='flex space-x-2'>
          <Badge variant='outline' className={cn('capitalize', badgeColor)}>
            {row.getValue('code')}
          </Badge>
        </div>
      )
    },
    enableSorting: false,
  },
  {
    accessorKey: 'message',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='错误信息' />
    ),
    cell: ({ row }) => (
      <OperationLogViewDialog data={row.original} type='view-msg' />
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'params',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求参数' />
    ),
    cell: ({ row }) => (
      <OperationLogViewDialog data={row.original} type='view-params' />
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'body',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求体' />
    ),
    cell: ({ row }) => (
      <OperationLogViewDialog data={row.original} type='view-body' />
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'result',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='响应结果' />
    ),
    cell: ({ row }) => (
      <OperationLogViewDialog data={row.original} type='view-result' />
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'startTime',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求开始时间' />
    ),
    cell: ({ row }) => (
      <div>{format(row.original.startTime, 'yyyy-MM-dd HH:mm:ss')}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'endTime',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='请求结束时间' />
    ),
    cell: ({ row }) => (
      <div>{format(row.original.endTime, 'yyyy-MM-dd HH:mm:ss')}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'latency',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='耗时' />
    ),
    cell: ({ row }) => <div>{row.original.latency}ms</div>,
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => <DataTableRowActions row={row} showInfo />,
  },
]
