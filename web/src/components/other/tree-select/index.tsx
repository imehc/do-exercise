import { useRef, useState, useDeferredValue, useMemo } from 'react'
import { IconChevronDown, IconSearch, IconX } from '@tabler/icons-react'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { Button, buttonVariants } from '~/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '~/components/ui/popover'
import { LoadingSpinner } from '..'
import { getTreeValueLabelMap } from './tree-utils'
import { TreeView } from './tree-view'
import { TreeSelectProps } from './types'

export { transformData } from './transfer'

export type AriaInvalidProps = Pick<
  React.HTMLAttributes<HTMLInputElement>,
  'aria-invalid'
>

export type TreeSelectComponentProps = TreeSelectProps &
  AriaInvalidProps & {
    className?: string
    loading?: boolean
    mode?: 'popup' | 'inline' | 'view'
    readonly?: boolean
  }

export const TreeSelect = ({
  value,
  onChange,
  data,
  className,
  loading,
  mode = 'popup',
  multiple = false,
  placeholder = '',
  'aria-invalid': invalid,
  readonly = false,
}: TreeSelectComponentProps) => {
  const ref = useRef<HTMLButtonElement>(null)
  const [search, setSearch] = useState<string | undefined>('')
  const deferredSearch = useDeferredValue(search)

  const valueLabelMap = useMemo(() => {
    return getTreeValueLabelMap(data)
  }, [data])

  if (mode === 'inline' || mode === 'view') {
    return (
      <div
        className={cn(
          'border-input w-[var(--radix-popover-trigger-width)] rounded-md border p-0 shadow-xs transition-[color,box-shadow]',
          className
        )}
      >
        <TreeSelectPanel
          loading={loading}
          data={data}
          value={value}
          onChange={onChange}
          multiple={multiple}
          search={search}
          deferredSearch={deferredSearch}
          setSearch={setSearch}
          readonly={readonly}
          mode={mode}
        />
      </div>
    )
  }

  return (
    <Popover modal>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          className={cn(
            'hover:bg-background h-fit min-h-10 items-center justify-end py-1.5 pr-0 pl-1.5',
            invalid && 'border-error-500',
            className,
            value.length > 1 && 'h-auto'
          )}
          ref={ref}
          aria-invalid={invalid}
        >
          <div className='relative flex grow flex-wrap items-center gap-[6px] overflow-hidden'>
            {value.length > 0 ? (
              value.map((v) => (
                <Badge
                  key={v}
                  variant='secondary'
                  className='hover:text-primary gap-1.5 rounded-sm px-1.5 py-0.5 text-left font-semibold text-wrap hover:bg-indigo-50'
                >
                  {valueLabelMap.get(v) ?? v}
                  <div
                    onClick={(e) => {
                      e.preventDefault()
                      onChange?.(value.filter((value) => value !== v))
                    }}
                    onKeyDown={(e) => {
                      if (e.key === ' ' || e.key === 'Enter') {
                        onChange?.(value.filter((value) => value !== v))
                        ref.current?.focus()
                      }
                    }}
                    role='button'
                    tabIndex={0}
                  >
                    <IconX size={14} />
                  </div>
                </Badge>
              ))
            ) : (
              <span className='text-muted-foreground ml-1.5 text-base font-normal md:text-sm'>
                {placeholder}
              </span>
            )}
          </div>

          {value.length > 0 && (
            <div
              className={cn(
                buttonVariants({ size: 'sm', variant: 'ghost' }),
                'hover:text-primary flex h-auto rounded-sm border-none px-2 py-0 text-gray-500 transition-colors hover:bg-transparent'
              )}
              onClick={(e) => {
                e.preventDefault()
                onChange?.([])
              }}
              onKeyDown={(e) => {
                if (e.key === ' ' || e.key === 'Enter') {
                  onChange?.([])
                  ref.current?.focus()
                }
              }}
              role='button'
              tabIndex={0}
            >
              <IconX className='size-4' />
            </div>
          )}
          <span className='bg-border w-px self-stretch' />
          <div className='text-muted-foreground mr-1.25 flex items-center px-2 opacity-50'>
            <IconChevronDown size={16} />
          </div>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className='w-[var(--radix-popover-trigger-width)] p-0'
        align='start'
      >
        <TreeSelectPanel
          loading={loading}
          data={data}
          value={value}
          onChange={onChange}
          multiple={multiple}
          search={search}
          deferredSearch={deferredSearch}
          setSearch={setSearch}
          readonly={readonly}
          mode={mode}
        />
      </PopoverContent>
    </Popover>
  )
}

interface TreeSelectPanelProps extends Pick<
  TreeSelectComponentProps,
  'loading' | 'data' | 'onChange' | 'multiple' | 'value' | 'readonly' | 'mode'
> {
  search?: string
  setSearch: React.Dispatch<React.SetStateAction<string | undefined>>
  deferredSearch?: string
}

function TreeSelectPanel({
  search,
  setSearch,
  loading,
  data,
  value,
  onChange,
  multiple,
  deferredSearch,
  readonly,
  mode,
}: TreeSelectPanelProps) {
  return (
    <>
      {mode !== 'view' && (
        <div className='flex items-center border-b px-3' cmdk-input-wrapper=''>
          <IconSearch className='mr-2 size-4 shrink-0 opacity-50' />
          <input
            value={search ?? ''}
            onChange={(e) => {
              setSearch(e.target.value)
            }}
            readOnly={readonly}
            className='placeholder:text-muted-foreground flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50'
          />
        </div>
      )}
      <div
        className={cn(
          'max-h-56 overflow-y-auto p-2',
          mode === 'view' && 'max-h-full'
        )}
      >
        {loading && (
          <div className='p-8 text-center text-sm text-gray-400'>
            <LoadingSpinner />
          </div>
        )}
        {!loading && data.length === 0 && (
          <div className='p-8 text-center text-sm'>No results</div>
        )}
        {!loading && data.length > 0 && (
          <TreeView
            value={value}
            onChange={onChange}
            data={data}
            searchValue={deferredSearch}
            multiple={multiple}
            readonly={readonly}
          />
        )}
      </div>
    </>
  )
}
