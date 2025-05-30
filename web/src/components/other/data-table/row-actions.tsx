import { Row } from '@tanstack/react-table'
import {
  IconEdit,
  IconInfoHexagon,
  IconTrash,
  IconTablePlus,
  IconDots,
} from '@tabler/icons-react'
import { useFormDialog } from '~/provider'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

type MenuOption = {
  title?: string
  disable?: boolean
}

interface DataTableRowActionsProps<T> {
  row: Row<T>
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

export function DataTableRowActions<T>({
  row,
  showEdit = false,
  editOptions = {
    title: 'Edit',
    disable: false,
  },
  showDelete = false,
  deleteOptions = {
    title: 'Delete',
    disable: false,
  },
  showInfo = false,
  infoOptions = {
    title: 'Info',
    disable: false,
  },
  showAdd = false,
  addOptions = {
    title: 'Add',
    disable: false,
  },
  children,
}: DataTableRowActionsProps<T>) {
  const { setOpen, setCurrentRow } = useFormDialog<T>()

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
          disabled={addOptions.disable}
          onClick={() => {
            setCurrentRow(row.original)
            setOpen('add-child')
          }}
        >
          <IconTablePlus />
          <span className='sr-only'>Add child</span>
        </Button>
      )}
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button
            variant='ghost'
            className='data-[state=open]:bg-muted flex h-8 w-8 p-0'
          >
            <IconDots className='h-4 w-4' />
            <span className='sr-only'>Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='end' className='w-[160px]'>
          <DropdownMenuItem
            hidden={!showInfo}
            disabled={infoOptions.disable}
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('view-info')
            }}
          >
            {infoOptions.title}
            <DropdownMenuShortcut>
              <IconInfoHexagon size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem
            hidden={!showEdit}
            disabled={editOptions.disable}
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('edit')
            }}
          >
            {editOptions.title}
            <DropdownMenuShortcut>
              <IconEdit size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuSeparator hidden={!showDelete} />
          <DropdownMenuItem
            hidden={!showDelete}
            disabled={deleteOptions.disable}
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('delete')
            }}
            className='text-red-500!'
          >
            {deleteOptions.title}
            <DropdownMenuShortcut>
              <IconTrash size={16} className='text-red-500!' />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
