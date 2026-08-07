import { useEffect, useState } from 'react'
import {
  ColumnDef,
  ColumnFiltersState,
  ColumnVisibilityState,
  RowData,
  SortingState,
  useTable,
} from '@tanstack/react-table'
import { Trans } from '@lingui/react/macro'
import { Pagination } from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '~/components/ui/table'
import { LoadingSpinner } from '..'
import { DataTableFeatures, dataTableFeatures } from './features'
import { DataTablePagination } from './pagination'
import { DataTableToolbar, DataTableToolbarProps } from './toolbar'

export { DataTableColumnHeader } from './column-header'
export { DataTableRowActions } from './row-actions'

// 固定列配置接口
export interface StickyColumnConfig {
  left?: string[] // 左侧固定列的 accessorKey
  right?: string[] // 右侧固定列的 accessorKey
}

interface DataTableProps<T extends RowData>
  extends Pick<
    DataTableToolbarProps<T>,
    'serchOptions' | 'search' | 'enableClientPagination' | 'navigate'
  > {
  columns: ColumnDef<DataTableFeatures, T>[]
  data: T[]
  meta?: Pagination
  isLoading?: boolean
  getColumnTitle?: (columnId: string) => string
  stickyColumns?: StickyColumnConfig // 新的固定列配置
}

export function DataTable<T extends RowData>({
  columns,
  data,
  meta,
  navigate,
  enableClientPagination = false,
  search,
  serchOptions,
  isLoading = false,
  getColumnTitle,
  stickyColumns,
}: DataTableProps<T>) {
  const [rowSelection, setRowSelection] = useState({})
  const [columnVisibility, setColumnVisibility] =
    useState<ColumnVisibilityState>({})
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

  const table = useTable({
    features: dataTableFeatures,
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
    onPaginationChange: enableClientPagination
      ? (updater) =>
          setClientPagination((old) =>
            updater instanceof Function ? updater(old) : updater
          )
      : undefined,
  })

  // 获取固定列配置
  const getStickyConfig = () => {
    if (stickyColumns) {
      return stickyColumns
    }
    return { left: [], right: [] }
  }

  const stickyConfig = getStickyConfig()

  // 获取列的固定样式
  const getStickyClass = (columnId: string, isHeader: boolean = false) => {
    const leftSticky = stickyConfig.left?.includes(columnId)
    const rightSticky = stickyConfig.right?.includes(columnId)

    if (!leftSticky && !rightSticky) return ''

    const baseClass = 'sticky z-50 bg-background'
    const hoverClass = isHeader ? '' : 'group-hover/row:bg-muted/50'

    if (leftSticky) {
      return cn(baseClass, 'left-0 shadow-sm', hoverClass)
    }

    if (rightSticky) {
      return cn(baseClass, 'right-0 shadow-sm', hoverClass)
    }

    return ''
  }

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
        <div className='overflow-x-auto'>
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className='group/row'>
                  {headerGroup.headers.map((header) => {
                    const stickyClass = getStickyClass(header.column.id, true)
                    const isSticky = stickyClass.includes('sticky')
                    const isLeftSticky = stickyConfig.left?.includes(
                      header.column.id
                    )
                    const isRightSticky = stickyConfig.right?.includes(
                      header.column.id
                    )

                    return (
                      <TableHead
                        key={header.id}
                        colSpan={header.colSpan}
                        className={cn(
                          (
                            header.column.columnDef.meta as {
                              className?: string
                            }
                          )?.className ?? '',
                          stickyClass
                        )}
                        style={
                          isSticky
                            ? {
                                backgroundColor: 'var(--background)',
                                position: 'sticky',
                                zIndex: 50,
                                left: isLeftSticky ? 0 : undefined,
                                right: isRightSticky ? 0 : undefined,
                                boxShadow: isLeftSticky
                                  ? '2px 0 4px rgba(0,0,0,0.1)'
                                  : isRightSticky
                                    ? '-2px 0 4px rgba(0,0,0,0.1)'
                                    : undefined,
                              }
                            : undefined
                        }
                      >
                        {header.isPlaceholder ? null : (
                          <table.FlexRender header={header} />
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
                    {row.getVisibleCells().map((cell) => {
                      const stickyClass = getStickyClass(cell.column.id)
                      const isSticky = stickyClass.includes('sticky')
                      const isLeftSticky = stickyConfig.left?.includes(
                        cell.column.id
                      )
                      const isRightSticky = stickyConfig.right?.includes(
                        cell.column.id
                      )

                      return (
                        <TableCell
                          key={cell.id}
                          className={cn(
                            (
                              cell.column.columnDef.meta as {
                                className?: string
                              }
                            )?.className ?? '',
                            stickyClass
                          )}
                          style={
                            isSticky
                              ? {
                                  backgroundColor: 'var(--background)',
                                  position: 'sticky',
                                  zIndex: 50,
                                  left: isLeftSticky ? 0 : undefined,
                                  right: isRightSticky ? 0 : undefined,
                                  boxShadow: isLeftSticky
                                    ? '2px 0 4px rgba(0,0,0,0.1)'
                                    : isRightSticky
                                      ? '-2px 0 4px rgba(0,0,0,0.1)'
                                      : undefined,
                                }
                              : undefined
                          }
                        >
                          <table.FlexRender cell={cell} />
                        </TableCell>
                      )
                    })}
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
