import { UseNavigateResult } from '@tanstack/react-router'
import { ReactTable, RowData } from '@tanstack/react-table'
import {
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
} from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { Pagination } from '~/do-exercise-api'
import { Button } from '~/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { DataTableFeatures } from './features'

export interface DataTablePaginationProps<TData extends RowData> {
  navigate?: UseNavigateResult<'/api/'>
  table: ReactTable<DataTableFeatures, TData>
  enableClientPagination?: boolean
  meta?: Pagination
}
export function DataTablePagination<TData extends RowData>({
  table,
  enableClientPagination = false,
  navigate,
  meta,
}: DataTablePaginationProps<TData>) {
  const pageIndex = table.state.pagination.pageIndex
  const pageSize = table.state.pagination.pageSize
  const pageCount = table.getPageCount()
  const total = enableClientPagination
    ? table.getPrePaginatedRowModel().rows.length
    : (meta?.total ?? 0)

  const selectRow = table.getFilteredSelectedRowModel().rows.length
  const selectTotalRow = table.getFilteredRowModel().rows.length
  const currentPage = table.state.pagination.pageIndex + 1
  const totalPage = table.getPageCount()

  return (
    <div
      className='flex items-center justify-between overflow-clip px-2'
      style={{ overflowClipMargin: 1 }}
    >
      <div className='text-muted-foreground hidden flex-1 text-sm sm:block'>
        <Trans>
          已选{selectRow}条,共{selectTotalRow}条
        </Trans>
        <span className='ml-2'>
          <Trans>总共{total}条</Trans>
        </span>
      </div>
      <div className='flex items-center sm:space-x-6 lg:space-x-8'>
        <div className='text-muted-foreground text-sm'></div>
        <div className='flex items-center space-x-2'>
          <p className='hidden text-sm font-medium sm:block'>
            <Trans>每页条数</Trans>
          </p>
          <Select
            value={pageSize + ''}
            onValueChange={(value) => {
              const newSize = Number(value)
              table.setPageSize(newSize)
              table.setPageIndex(0)
              if (!enableClientPagination) {
                navigate?.({
                  search: {
                    page: 1,
                    pageSize: newSize,
                  },
                })
              }
            }}
          >
            <SelectTrigger className='h-8 w-[70px]'>
              <SelectValue placeholder={table.state.pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side='top'>
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className='inline-flex items-center justify-center text-sm font-medium'>
          <Trans>
            第{currentPage}页，共{totalPage}页
          </Trans>
        </div>
        <div className='flex items-center space-x-2'>
          <Button
            variant='outline'
            className='hidden h-8 w-8 p-0 lg:flex'
            onClick={() => {
              table.setPageIndex(0)
              if (!enableClientPagination) {
                navigate?.({ search: (s) => ({ ...s, page: 1 }) })
              }
            }}
            disabled={!table.getCanPreviousPage()}
          >
            <span className='sr-only'>
              <Trans>第一页</Trans>
            </span>
            <IconChevronsLeft className='h-4 w-4' />
          </Button>
          <Button
            variant='outline'
            className='h-8 w-8 p-0'
            onClick={() => {
              table.previousPage()
              if (!enableClientPagination) {
                navigate?.({ search: (s) => ({ ...s, page: s.page - 1 }) })
              }
            }}
            disabled={!table.getCanPreviousPage()}
          >
            <span className='sr-only'>
              <Trans>上一页</Trans>
            </span>
            <IconChevronLeft className='h-4 w-4' />
          </Button>
          <Button
            variant='outline'
            className='h-8 w-8 p-0'
            onClick={() => {
              table.nextPage()
              if (!enableClientPagination) {
                navigate?.({
                  search: (s) => ({
                    ...s,
                    page: pageIndex + 2,
                  }),
                })
              }
            }}
            disabled={!table.getCanNextPage()}
          >
            <span className='sr-only'>
              <Trans>下一页</Trans>
            </span>
            <IconChevronRight className='h-4 w-4' />
          </Button>
          <Button
            variant='outline'
            className='hidden h-8 w-8 p-0 lg:flex'
            onClick={() => {
              table.setPageIndex(pageCount - 1)
              if (!enableClientPagination) {
                navigate?.({
                  search: (s) => ({
                    ...s,
                    page: pageCount,
                  }),
                })
              }
            }}
            disabled={!table.getCanNextPage()}
          >
            <span className='sr-only'>
              <Trans>最后一页</Trans>
            </span>
            <IconChevronsRight className='h-4 w-4' />
          </Button>
        </div>
      </div>
    </div>
  )
}
