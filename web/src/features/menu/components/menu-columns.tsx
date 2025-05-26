import { ChevronDownIcon, ChevronRightIcon } from '@radix-ui/react-icons'
import { ColumnDef } from '@tanstack/react-table'
import { SysMenuTree } from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTableColumnHeader, DataTableRowActions } from '~/components/other'
import { callMenuMapping, callMenuTypes, callVisibleTypes } from '../data/data'

export const columns: ColumnDef<SysMenuTree>[] = [
  {
    id: 'expander',
    header: () => null,
    minSize: 20,
    size: 30,
    maxSize: 40,
    cell: ({ row }) => {
      return row.getCanExpand() ? (
        <Button
          variant='ghost'
          className='h-6 w-6 p-0 hover:bg-transparent'
          onClick={row.getToggleExpandedHandler()}
        >
          {row.getIsExpanded() ? (
            <ChevronDownIcon className='h-4 w-4' />
          ) : (
            <ChevronRightIcon className='h-4 w-4' />
          )}
        </Button>
      ) : null
    },
  },
  {
    accessorKey: 'name',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='菜单名称' />
    ),
    cell: ({ row }) => (
      <div
        className='w-fit text-nowrap'
        style={{
          paddingLeft: `${row.depth * 12}px`,
        }}
      >
        {row.original.name}
      </div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'permission',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='权限标识' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.permission || '-'}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'icon',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='图标' />
    ),
    cell: ({ row }) => (
      <div className='w-fit text-nowrap'>{row.original.icon || '-'}</div>
    ),
    enableSorting: false,
  },
  {
    accessorKey: 'type',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='菜单类型' />
    ),
    cell: ({ row }) => {
      const { type } = row.original
      const badgeColor = callMenuTypes.get(type)
      return (
        <div className='flex space-x-2'>
          <Badge variant='outline' className={cn('capitalize', badgeColor)}>
            {callMenuMapping.get(row.getValue('type')) ?? '-'}
          </Badge>
        </div>
      )
    },
    enableSorting: false,
  },
  {
    accessorKey: 'route',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='路由' />
    ),
    cell: ({ row }) => <div>{row.original.route || '-'}</div>,
    enableSorting: false,
  },
  {
    accessorKey: 'component',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='组件' />
    ),
    cell: ({ row }) => <div>{row.original.component || '-'}</div>,
    enableSorting: false,
  },
  {
    accessorKey: 'sort',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='序号' />
    ),
    cell: ({ row }) => <div>{row.original.sort ?? '-'}</div>,
    enableSorting: false,
  },
  {
    accessorKey: 'visible',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='是否可见' />
    ),
    cell: ({ row }) => {
      const { visible } = row.original
      const badgeColor = callVisibleTypes.get(visible || false)
      return (
        <div className='flex space-x-2'>
          <Badge variant='outline' className={cn('capitalize', badgeColor)}>
            {row.getValue('visible') ? '是' : '否'}
          </Badge>
        </div>
      )
    },
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => (
      <DataTableRowActions row={row} showEdit showDelete showInfo />
    ),
  },
]
