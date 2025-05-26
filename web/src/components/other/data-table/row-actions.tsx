import { DotsHorizontalIcon } from '@radix-ui/react-icons'
import { Row } from '@tanstack/react-table'
import { IconEdit, IconInfoHexagon, IconTrash } from '@tabler/icons-react'
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
}: DataTableRowActionsProps<T>) {
  const { setOpen, setCurrentRow } = useFormDialog<T>()

  if (!showEdit && !showDelete && !showInfo) {
    return null
  }

  return (
    <>
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button
            variant='ghost'
            className='data-[state=open]:bg-muted flex h-8 w-8 p-0'
          >
            <DotsHorizontalIcon className='h-4 w-4' />
            <span className='sr-only'>Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='end' className='w-[160px]'>
          {showInfo && (
            <DropdownMenuItem
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
          )}
          {showEdit && (
            <DropdownMenuItem
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
          )}
          {showDelete && (
            <>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                disabled={deleteOptions.disable}
                onClick={() => {
                  setCurrentRow(row.original)
                  setOpen('delete')
                }}
                className='text-red-500!'
              >
                {editOptions.title}
                <DropdownMenuShortcut>
                  <IconTrash size={16} className='text-red-500!' />
                </DropdownMenuShortcut>
              </DropdownMenuItem>
            </>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  )
}
