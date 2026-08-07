import React from 'react'
import {
  CellContext,
  ColumnDef,
  HeaderContext,
  RowData,
} from '@tanstack/react-table'
import { i18n } from '@lingui/core'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { DataTableFeatures } from './features'
import { DataTableColumnHeader } from './column-header'

type BaseColumnConfig<TData extends RowData> = {
  id?: string
  title?: () => string
  header?: (
    props: HeaderContext<DataTableFeatures, TData, unknown>
  ) => React.ReactNode
  cell?: (
    props: CellContext<DataTableFeatures, TData, unknown>
  ) => React.ReactNode
  options?: Partial<ColumnDef<DataTableFeatures, TData>>
}

type AccessorColumnConfig<TData extends RowData> = BaseColumnConfig<TData> & {
  key: keyof TData
}

type IdColumnConfig<TData extends RowData> = BaseColumnConfig<TData> & {
  id: string
}

export type ColumnConfig<TData extends RowData> =
  | AccessorColumnConfig<TData>
  | IdColumnConfig<TData>

export const createColumn = <TData extends RowData,>(
  config: ColumnConfig<TData>
): ColumnDef<DataTableFeatures, TData> => {
  const baseConfig: Partial<ColumnDef<DataTableFeatures, TData>> = {
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
    } as ColumnDef<DataTableFeatures, TData>
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
  } as ColumnDef<DataTableFeatures, TData>
}

export const createActionColumn = <TData extends RowData,>(
  cell: (props: CellContext<DataTableFeatures, TData, unknown>) => React.ReactNode,
  options: Partial<ColumnDef<DataTableFeatures, TData>> = {}
): ColumnDef<DataTableFeatures, TData> => ({
  id: 'actions',
  cell,
  ...options,
})

export const createDateColumn = <TData extends RowData,>(
  key: keyof TData,
  title: () => string,
  options: Partial<ColumnDef<DataTableFeatures, TData>> = {}
): ColumnDef<DataTableFeatures, TData> => ({
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

export const createBadgeColumn = <TData extends RowData,>(
  key: keyof TData,
  title: () => string,
  getBadgeColor: (value: unknown) => string | undefined,
  getBadgeText: (value: unknown) => React.ReactNode,
  options: Partial<ColumnDef<DataTableFeatures, TData>> = {}
): ColumnDef<DataTableFeatures, TData> => ({
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
