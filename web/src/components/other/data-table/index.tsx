import { useEffect, useState } from 'react'
import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getExpandedRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { Trans } from '@lingui/react/macro'
import { Pagination } from '~/do-exercise-api'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '~/components/ui/table'
import { LoadingSpinner } from '..'
import { DataTablePagination } from './pagination'
import { DataTableToolbar, DataTableToolbarProps } from './toolbar'

export { DataTableColumnHeader } from './column-header'
export { DataTableRowActions } from './row-actions'

interface DataTableProps<T>
  extends Pick<
    DataTableToolbarProps<T>,
    'serchOptions' | 'search' | 'enableClientPagination' | 'navigate'
  > {
  columns: ColumnDef<T>[]
  data: T[]
  meta?: Pagination
  isLoading?: boolean
  getColumnTitle?: (columnId: string) => string
}

export function DataTable<T>({
  columns,
  data,
  meta,
  navigate,
  enableClientPagination = false,
  search,
  serchOptions,
  isLoading = false,
  getColumnTitle,
}: DataTableProps<T>) {
  const [rowSelection, setRowSelection] = useState({})
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [sorting, setSorting] = useState<SortingState>([])
  const [expanded, setExpanded] = useState({})
  const [clientPagination, setClientPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  // 默认展开展开一级菜单
  useEffect(() => {
    const defaultExpanded: Record<string, boolean> = {}

    data.forEach((_, index) => {
      const rowId = String(index)
      defaultExpanded[rowId] = true // 展开一级菜单
    })

    setExpanded(defaultExpanded)
  }, [data])

  const page = meta?.page || 1
  const pageSize = meta?.pageSize || 10
  const total = meta?.total || 0

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnVisibility,
      rowSelection,
      columnFilters,
      pagination: {
        pageIndex: enableClientPagination
          ? clientPagination.pageIndex
          : page - 1,
        pageSize: enableClientPagination ? clientPagination.pageSize : pageSize,
      },
      expanded,
    },
    enableExpanding: true,
    onExpandedChange: setExpanded,
    getSubRows: (row) => (row as { children?: T[] })?.children,
    manualPagination: !enableClientPagination,
    pageCount: enableClientPagination
      ? Math.ceil(data.length / clientPagination.pageSize)
      : Math.ceil(total / pageSize),
    enableRowSelection: true,
    onRowSelectionChange: setRowSelection,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
    getExpandedRowModel: getExpandedRowModel(),
    onPaginationChange: enableClientPagination
      ? (updater) =>
          setClientPagination((old) =>
            updater instanceof Function ? updater(old) : updater
          )
      : undefined,
  })

  return (
    <div className='space-y-4'>
      <DataTableToolbar
        table={table}
        serchOptions={serchOptions}
        navigate={navigate}
        search={search}
        enableClientPagination={enableClientPagination}
        getColumnTitle={getColumnTitle}
      />
      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className='group/row'>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead
                      key={header.id}
                      colSpan={header.colSpan}
                      className={
                        (header.column.columnDef.meta as { className?: string })
                          ?.className ?? ''
                      }
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='text-muted-foreground pointer-events-none h-64 text-center'
                >
                  <LoadingSpinner />
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                  className='group/row'
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className={
                        (cell.column.columnDef.meta as { className?: string })
                          ?.className ?? ''
                      }
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  <Trans>没有内容.</Trans>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {((!!meta && !!navigate) || enableClientPagination) && (
        <DataTablePagination
          table={table}
          meta={meta}
          navigate={navigate}
          enableClientPagination={enableClientPagination}
        />
      )}
    </div>
  )
}
