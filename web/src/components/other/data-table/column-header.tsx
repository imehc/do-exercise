import { CellData, Column, RowData } from '@tanstack/react-table'
import {
  IconArrowDown,
  IconArrowUp,
  IconSelector,
  IconEyeOff,
} from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { DataTableFeatures } from './features'

interface DataTableColumnHeaderProps<
  TData extends RowData,
  TValue,
> extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<DataTableFeatures, TData, TValue>
  title: string
}

export function DataTableColumnHeader<
  TData extends RowData,
  TValue extends CellData,
>({ column, title, className }: DataTableColumnHeaderProps<TData, TValue>) {
  if (!column.getCanSort()) {
    return <div className={cn(className)}>{title}</div>
  }

  return (
    <div className={cn('flex items-center space-x-2', className)}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant='ghost'
            size='sm'
            className='data-[state=open]:bg-accent -ml-3 h-8'
          >
            <span>{title}</span>
            {column.getIsSorted() === 'desc' ? (
              <IconArrowDown className='ml-2 h-4 w-4' />
            ) : column.getIsSorted() === 'asc' ? (
              <IconArrowUp className='ml-2 h-4 w-4' />
            ) : (
              <IconSelector className='ml-2 h-4 w-4' />
            )}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='start'>
          <DropdownMenuItem onClick={() => column.toggleSorting(false)}>
            <IconArrowUp className='text-muted-foreground/70 mr-2 h-3.5 w-3.5' />
            <Trans>升序</Trans>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => column.toggleSorting(true)}>
            <IconArrowDown className='text-muted-foreground/70 mr-2 h-3.5 w-3.5' />
            <Trans>降序</Trans>
          </DropdownMenuItem>
          {column.getCanHide() && (
            <>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => column.toggleVisibility(false)}>
                <IconEyeOff className='text-muted-foreground/70 mr-2 h-3.5 w-3.5' />
                <Trans>隐藏</Trans>
              </DropdownMenuItem>
            </>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
