import { useRef, useState, useDeferredValue, useMemo, useEffect } from 'react'
import { IconChevronDown, IconSearch, IconX } from '@tabler/icons-react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { cn } from '~/lib/utils'
import { getIconComponentList, toIconComponent } from '~/utils/icon'
import { useGridColumnCount } from '~/hooks/use-grid-column-count'
import { Button } from '~/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '~/components/ui/popover'
import { Skeleton } from '~/components/ui/skeleton'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '~/components/ui/tooltip'

export type IconSelectProps = {
  value?: string
  onChange: (value?: string) => void
  className?: string
  loading?: boolean
  placeholder?: string
}

export const IconSelect = ({
  value,
  onChange,
  className,
  loading,
  placeholder = '',
}: IconSelectProps) => {
  const ref = useRef<HTMLButtonElement>(null)

  const SelectedIcon = toIconComponent(value)

  return (
    <Popover modal>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          className={cn(
            'hover:bg-background h-fit min-h-10 items-center justify-end py-1.5 pr-0 pl-1.5',
            className
          )}
          ref={ref}
        >
          <div className='relative flex grow flex-wrap items-center gap-[6px] overflow-hidden'>
            {value && SelectedIcon ? (
              <div className='ml-1.5 flex items-center justify-start gap-x-3'>
                <SelectedIcon className='scale-150' />
                <p>{value}</p>
              </div>
            ) : (
              <span className='text-muted-foreground ml-1.5 text-base font-normal md:text-sm'>
                {placeholder}
              </span>
            )}
          </div>
          {!!value && (
            <div
              className='hover:text-primary flex h-auto rounded-sm border-none px-2 py-0 text-gray-500 transition-colors hover:bg-transparent'
              onClick={(e) => {
                e.preventDefault()
                onChange(undefined)
              }}
              onKeyDown={(e) => {
                if (e.key === ' ' || e.key === 'Enter') {
                  onChange(undefined)
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
        <IconSelectContent
          value={value}
          onChange={onChange}
          loading={loading}
        />
      </PopoverContent>
    </Popover>
  )
}

const estimateSize = 48

interface IconSelectContentProps {
  value?: string
  onChange: (value?: string) => void
  loading?: boolean
}

function IconSelectContent({
  loading,
  value,
  onChange,
}: IconSelectContentProps) {
  const contentRef = useRef<HTMLDivElement>(null)
  const [search, setSearch] = useState<string>(() => value ?? '')
  const deferredSearch = useDeferredValue(search)
  const filteredIcons = useMemo(() => {
    return getIconComponentList().filter(({ label }) =>
      label.toLowerCase().includes(deferredSearch.toLowerCase())
    )
  }, [deferredSearch])

  const count = useGridColumnCount(contentRef, estimateSize)
  const { getVirtualItems, getTotalSize, measure } = useVirtualizer({
    count: filteredIcons.length,
    getScrollElement: () => contentRef.current,
    estimateSize: () => estimateSize,
    enabled: count > 0 && !!contentRef.current && !loading,
    debug: import.meta.env.DEV,
    overscan: 5,
    lanes: count,
    gap: 8,
  })

  useEffect(() => {
    if (count > 0) {
      measure()
    }
    measure()
  }, [count, measure])

  return (
    <>
      <div className='flex items-center border-b px-3'>
        <IconSearch className='mr-2 size-4 shrink-0 opacity-50' />
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className='placeholder:text-muted-foreground flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50'
          placeholder='Search icons...'
        />
      </div>
      <div className='h-56 overflow-y-auto p-2' ref={contentRef}>
        {loading ? (
          <div
            className='relative h-full'
            style={{ height: `${estimateSize * 4}px` }}
          >
            <IconSkeletonGrid count={count || 4} />
          </div>
        ) : (
          <div
            className='relative h-full'
            style={{
              height: `${getTotalSize()}px`,
            }}
          >
            <TooltipProvider>
              {getVirtualItems().map((virtualRow) => {
                const { label, icon: Icon } = filteredIcons[virtualRow.index]

                return (
                  <div
                    key={virtualRow.index}
                    className='flex items-center justify-center'
                    style={{
                      position: 'absolute',
                      top: 0,
                      left: `${virtualRow.lane * (100 / count)}%`,
                      width: `${100 / count}%`,
                      height: `${virtualRow.size}px`,
                      transform: `translateY(${virtualRow.start}px)`,
                    }}
                  >
                    <Tooltip key={label}>
                      <TooltipTrigger asChild>
                        <button
                          onClick={() => onChange(label)}
                          className={cn(
                            'hover:bg-muted flex aspect-square items-center justify-center rounded-md border p-1',
                            value === label &&
                              'border-primary bg-primary/10 text-primary'
                          )}
                        >
                          <Icon size={30} />
                        </button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{label}</p>
                      </TooltipContent>
                    </Tooltip>
                  </div>
                )
              })}
            </TooltipProvider>
          </div>
        )}
      </div>
    </>
  )
}

function IconSkeletonGrid({ count }: { count: number }) {
  return (
    <div className='relative'>
      {Array.from({ length: count * 4 }).map((_, i) => (
        <div
          key={i}
          className='flex items-center justify-center'
          style={{
            position: 'absolute',
            top: 0,
            left: `${(i % count) * (100 / count)}%`,
            width: `${100 / count}%`,
            height: `${estimateSize}px`,
            transform: `translateY(${Math.floor(i / count) * estimateSize}px)`,
          }}
        >
          <Skeleton className='aspect-square w-10 rounded-md' />
        </div>
      ))}
    </div>
  )
}
