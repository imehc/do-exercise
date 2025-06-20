import React from 'react'
import { ColumnDef, CellContext, HeaderContext } from '@tanstack/react-table'
import { i18n } from '@lingui/core'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { DataTableColumnHeader } from './column-header'

type BaseColumnConfig<TData> = {
  id?: string
  title?: () => string
  header?: (props: HeaderContext<TData, unknown>) => React.ReactNode
  cell?: (props: CellContext<TData, unknown>) => React.ReactNode
  options?: Partial<ColumnDef<TData>>
}

type AccessorColumnConfig<TData> = BaseColumnConfig<TData> & {
  key: keyof TData
}

type IdColumnConfig<TData> = BaseColumnConfig<TData> & {
  id: string
}

export type ColumnConfig<TData> =
  | AccessorColumnConfig<TData>
  | IdColumnConfig<TData>

export const createColumn = <TData,>(
  config: ColumnConfig<TData>
): ColumnDef<TData> => {
  const baseConfig: Partial<ColumnDef<TData>> = {
    enableSorting: false,
    ...config.options,
  }

  if ('key' in config) {
    return {
      ...baseConfig,
      accessorKey: config.key,
      header:
        config.header ??
        (config.title
          ? ({ column }) => (
              <DataTableColumnHeader column={column} title={config.title!()} />
            )
          : undefined),
      cell:
        config.cell ??
        (({ row }) => <div>{String(row.original[config.key])}</div>),
    } as ColumnDef<TData>
  }

  return {
    ...baseConfig,
    id: config.id,
    header:
      config.header ??
      (config.title
        ? ({ column }) => (
            <DataTableColumnHeader column={column} title={config.title!()} />
          )
        : undefined),
    cell: config.cell,
  } as ColumnDef<TData>
}

export const createActionColumn = <TData,>(
  cell: (props: CellContext<TData, unknown>) => React.ReactNode,
  options: Partial<ColumnDef<TData>> = {}
): ColumnDef<TData> => ({
  id: 'actions',
  cell,
  ...options,
})

export const createDateColumn = <TData,>(
  key: keyof TData,
  title: () => string,
  options: Partial<ColumnDef<TData>> = {}
): ColumnDef<TData> => ({
  accessorKey: key,
  header: ({ column }) => (
    <DataTableColumnHeader column={column} title={title()} />
  ),
  cell: ({ row }) => (
    <div>
      {i18n.date(row.original[key] as Date, {
        dateStyle: 'short',
        timeStyle: 'medium',
      })}
    </div>
  ),
  enableSorting: false,
  ...options,
})

export const createBadgeColumn = <TData,>(
  key: keyof TData,
  title: () => string,
  getBadgeColor: (value: unknown) => string | undefined,
  getBadgeText: (value: unknown) => React.ReactNode,
  options: Partial<ColumnDef<TData>> = {}
): ColumnDef<TData> => ({
  accessorKey: key,
  header: ({ column }) => (
    <DataTableColumnHeader column={column} title={title()} />
  ),
  cell: ({ row }) => {
    const value = row.original[key]
    const badgeColor = getBadgeColor(value)
    return (
      <div className='flex space-x-2'>
        <Badge variant='outline' className={cn('capitalize', badgeColor)}>
          {getBadgeText(value)}
        </Badge>
      </div>
    )
  },
  enableSorting: false,
  ...options,
})
