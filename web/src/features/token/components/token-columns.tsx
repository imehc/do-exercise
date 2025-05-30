import { format } from 'date-fns'
import { ColumnDef } from '@tanstack/react-table'
import { TokenInfo } from '~/do-exercise-api'
import {
  DataTableColumnHeader,
  DataTableRowActions,
  InlineCopy,
} from '~/components/other'
import { ToggleDisabledSwitch } from './token-toggle-disabled-switch'

export const columns: ColumnDef<TokenInfo>[] = [
  // 列选择
  // {
  //   id: 'select',
  //   header: ({ table }) => {
  //     const isAllPageRowsSelected = table.getIsAllPageRowsSelected()
  //     const isSomePageRowsSelected = table.getIsSomePageRowsSelected()

  //     return (
  //       <Checkbox
  //         checked={
  //           isSomePageRowsSelected && !isAllPageRowsSelected
  //             ? 'indeterminate'
  //             : isAllPageRowsSelected
  //         }
  //         onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
  //         aria-label='Select all'
  //         className='translate-y-[2px]'
  //       />
  //     )
  //   },
  //   cell: ({ row }) => (
  //     <Checkbox
  //       checked={row.getIsSelected()}
  //       onCheckedChange={(value) => row.toggleSelected(!!value)}
  //       aria-label='Select row'
  //       className='translate-y-[2px]'
  //     />
  //   ),
  //   enableSorting: false,
  //   enableHiding: false,
  //   meta: {
  //     className: 'w-10 text-center', // 控制宽度和居中
  //   },
  // },
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
    accessorKey: 'userId',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='用户ID' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.userId}</div>
    ),
    enableSorting: false,
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
    accessorKey: 'accessToken',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='令牌' />
    ),
    cell: ({ row }) => (
      <InlineCopy
        className='w-fit text-nowrap'
        text={row.original.accessToken}
      />
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'disabled',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='禁用状态' />
    ),
    cell: ({ row }) => {
      const token = row.original
      return (
        <ToggleDisabledSwitch
          accessToken={token.accessToken}
          disabled={token.disabled}
        />
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
    accessorKey: 'expiredAt',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='到期时间' />
    ),
    cell: ({ row }) => (
      <div>{format(row.original.expiredAt, 'yyyy-MM-dd HH:mm:ss')}</div>
    ),
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => <DataTableRowActions row={row} showDelete showInfo />,
  },
]
