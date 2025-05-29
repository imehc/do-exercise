import { format } from 'date-fns'
import { ColumnDef } from '@tanstack/react-table'
import { SysUser } from '~/do-exercise-api'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { DataTableColumnHeader, DataTableRowActions } from '~/components/other'

export const columns: ColumnDef<SysUser>[] = [
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
      <DataTableColumnHeader column={column} title='用户名' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.username}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'nickname',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='昵称' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.nickname || '-'}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'email',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='邮箱' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.email || '-'}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'avatar',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='头像' />
    ),
    cell: ({ row }) => {
      const avatar = row.original.avatar
      if (!avatar) {
        return '-'
      }
      // TODO: 根据前缀判断是否需要拼接完整的图片地址
      return (
        <div className='w-fit text-nowrap'>
          <Avatar>
            <AvatarImage src={avatar} alt={row.original.username} />
            <AvatarFallback>{row.original.username.slice(0, 2)}</AvatarFallback>
          </Avatar>
        </div>
      )
    },
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
