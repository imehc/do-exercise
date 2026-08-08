import { Row, RowData } from '@tanstack/react-table'
import {
  IconEdit,
  IconInfoHexagon,
  IconTrash,
  IconTablePlus,
  IconDots,
} from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { DialogType, useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { DataTableFeatures } from './features'

type Action = {
  type: DialogType | 'split-line'
  label: React.ReactNode
  icon: React.ReactNode
}

const actions = [
  {
    type: 'view-info',
    label: <Trans>详情</Trans>,
    icon: <IconInfoHexagon size={16} />,
  },
  {
    type: 'edit',
    label: <Trans>编辑</Trans>,
    icon: <IconEdit size={16} />,
  },
  {
    type: 'split-line',
    label: null,
    icon: null,
  },
  {
    type: 'delete',
    label: <Trans>删除</Trans>,
    icon: <IconTrash size={16} className='text-red-500!' />,
  },
] satisfies Action[]

type MenuOption = {
  title?: string
  disable?: boolean
}

interface DataTableRowActionsProps<T extends RowData> {
  row: Row<DataTableFeatures, T>
  showEdit?: boolean
  editOptions?: MenuOption
  showDelete?: boolean
  deleteOptions?: MenuOption
  showInfo?: boolean
  infoOptions?: MenuOption
  showAdd?: boolean
  addOptions?: MenuOption
  children?: ((data: T) => React.ReactNode) | React.ReactNode
}

export function DataTableRowActions<T extends RowData>({
  row,
  showEdit = false,
  editOptions,
  showDelete = false,
  deleteOptions,
  showInfo = false,
  infoOptions,
  showAdd = false,
  addOptions,
  children,
}: DataTableRowActionsProps<T>) {
  const { setOpen, setCurrentRow } = useFormDialog<T>()

  const handleHiddenState = (action: Action) => {
    switch (action.type) {
      case 'view-info': {
        return {
          show: showInfo,
          disabled: !!infoOptions?.disable,
        }
      }
      case 'edit':
        return {
          show: showEdit,
          disabled: !!editOptions?.disable,
        }
      case 'delete':
        return {
          show: showDelete,
          disabled: !!deleteOptions?.disable,
        }
      default:
        return {
          show: false,
          disabled: false,
        }
    }
  }

  if (!showEdit && !showDelete && !showInfo) {
    return null
  }

  return (
    <div className='flex items-center gap-x-2'>
      {typeof children === 'function' ? children?.(row.original) : children}
      {showAdd && (
        <Button
          variant='outline'
          size='icon'
          hidden={!showAdd}
          disabled={addOptions?.disable || false}
          onClick={() => {
            setCurrentRow(row.original)
            setOpen('add-child')
          }}
        >
          <IconTablePlus />
          <span className='sr-only'>
            <Trans>添加子元素</Trans>
          </span>
        </Button>
      )}
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button
            variant='ghost'
            className='data-[state=open]:bg-muted flex h-8 w-8 p-0'
          >
            <IconDots className='h-4 w-4' />
            <span className='sr-only'>
              <Trans>打开菜单</Trans>
            </span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='end' className='w-[160px]'>
          {actions.map((action) => {
            if (action.type === 'split-line') {
              return (
                <DropdownMenuSeparator key={action.type} hidden={!showDelete} />
              )
            }
            return (
              <DropdownMenuItem
                key={action.type}
                hidden={!handleHiddenState(action).show}
                disabled={handleHiddenState(action).disabled}
                onClick={() => {
                  setCurrentRow(row.original)
                  setOpen(action.type as DialogType)
                }}
              >
                <span
                  className={cn(action.type === 'delete' && 'text-red-500')}
                >
                  {action.label}
                </span>
                <DropdownMenuShortcut>{action.icon}</DropdownMenuShortcut>
              </DropdownMenuItem>
            )
          })}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
