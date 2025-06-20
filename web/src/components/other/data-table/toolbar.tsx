import { useEffect, useState } from 'react'
import { IconX } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useDebounce } from '~/hooks/use-debounce'
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import { DataTablePaginationProps } from './pagination'
import { DataTableViewOptions } from './view-options'

export interface DataTableToolbarProps<TData>
  extends Pick<
    DataTablePaginationProps<TData>,
    'enableClientPagination' | 'navigate' | 'table'
  > {
  search?: Record<string, unknown>
  serchOptions?: {
    key: keyof TData
    placeholder?: string
  }[]
  getColumnTitle?: (columnId: string) => string
}

export function DataTableToolbar<TData>({
  table,
  enableClientPagination,
  search,
  navigate,
  serchOptions,
  getColumnTitle,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0
  // 🔍 保存搜索字段值（服务端模式）由于需要防抖，所以必须在这里临时存放，而不是直接使用search
  const [rawSearchParams, setRawSearchParams] = useState<
    Record<string, string>
  >(() => (search as Record<string, string>) || {})

  const debouncedSearchParams = useDebounce(rawSearchParams)

  useEffect(() => {
    navigate?.({
      search: (s) => ({
        ...s,
        ...debouncedSearchParams,
        page: 1,
      }),
    })
  }, [debouncedSearchParams, navigate])

  const handleInputChange = (key: string, value: string) => {
    if (enableClientPagination) {
      table.getColumn(key)?.setFilterValue(value)
    } else {
      setRawSearchParams((old) => ({ ...old, [key]: value }))
    }
  }

  const handleReset = () => {
    if (enableClientPagination) {
      table.resetColumnFilters()
    } else {
      const resetParams = Object.fromEntries(
        Object.keys(rawSearchParams).map((k) => [k, ''])
      )
      setRawSearchParams(resetParams)
      // navigate?.({
      //   search: (s) => ({
      //     page: 1,
      //     pageSize: s.pageSize,
      //   }),
      // })
    }
  }

  return (
    <div className='flex items-center justify-between'>
      <div className='flex flex-1 flex-col-reverse items-start gap-y-2 sm:flex-row sm:items-center sm:space-x-2'>
        {serchOptions?.map((option) => {
          const key = option.key as string
          const value = enableClientPagination
            ? ((table.getColumn(key)?.getFilterValue() as string) ?? '')
            : rawSearchParams[key] || ''

          return (
            <Input
              key={option.key as string}
              placeholder={option.placeholder || t`筛选 ${key}...`}
              value={value}
              onChange={(e) => handleInputChange(key, e.target.value)}
              className='h-8 w-[150px] lg:w-[250px]'
            />
          )
        })}
        {(isFiltered || Object.values(rawSearchParams).some(Boolean)) && (
          <Button
            variant='ghost'
            onClick={handleReset}
            className='h-8 px-2 lg:px-3'
          >
            <Trans>重置</Trans>
            <IconX className='ml-2 h-4 w-4' />
          </Button>
        )}
      </div>
      <DataTableViewOptions table={table} getColumnTitle={getColumnTitle} />
    </div>
  )
}
