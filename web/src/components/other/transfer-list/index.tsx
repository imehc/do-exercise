import React, { useEffect } from 'react'
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { Checkbox } from '~/components/ui/checkbox'
import { Input } from '~/components/ui/input'
import { AriaInvalidProps } from '..'

type Item = {
  key: number
  label: string
  selected?: boolean
}

interface Props extends AriaInvalidProps {
  data: Item[]
  className?: string
  value?: number[]
  onChange: (value: number[]) => void
}

export function TransferList({
  data,
  value = [],
  onChange,
  className,
  'aria-invalid': invalid,
}: Props) {
  const [leftList, setLeftList] = React.useState<Item[]>(() =>
    data.filter((item) => !value.includes(item.key))
  )
  const [rightList, setRightList] = React.useState<Item[]>(() =>
    data.filter((item) => value.includes(item.key))
  )
  const [leftSearch, setLeftSearch] = React.useState('')
  const [rightSearch, setRightSearch] = React.useState('')

  useEffect(() => {
    onChange(rightList.map((item) => item.key))
  }, [onChange, rightList])

  const toggleSelection = (
    list: Item[],
    setList: React.Dispatch<React.SetStateAction<Item[]>>,
    key: number
  ) => {
    const updatedList = list.map((item) =>
      item.key === key ? { ...item, selected: !item.selected } : item
    )
    setList(updatedList)
  }

  const selectAll = (
    list: Item[],
    setList: React.Dispatch<React.SetStateAction<Item[]>>,
    selected: boolean
  ) => {
    const updatedList = list.map((item) => ({ ...item, selected }))
    setList(updatedList)
  }

  const moveToRight = () => {
    const selected = leftList.filter((item) => item.selected)
    setRightList([...rightList, ...selected])
    setLeftList(leftList.filter((item) => !item.selected))
  }

  const moveToLeft = () => {
    const selected = rightList.filter((item) => item.selected)
    setLeftList([...leftList, ...selected])
    setRightList(rightList.filter((item) => !item.selected))
  }

  const filteredLeft = leftList.filter((item) =>
    item.label.toLowerCase().includes(leftSearch.toLowerCase())
  )
  const filteredRight = rightList.filter((item) =>
    item.label.toLowerCase().includes(rightSearch.toLowerCase())
  )

  const countSelected = (list: Item[]) =>
    list.filter((item) => item.selected).length

  const getCheckboxState = (list: Item[]) => {
    const selectedCount = list.filter((item) => item.selected).length
    if (selectedCount === 0) return false
    if (selectedCount === list.length) return true
    return 'indeterminate'
  }

  return (
    <div className={cn('flex gap-2', className)}>
      <div
        className={cn(
          'w-1/2 rounded-md border shadow',
          invalid &&
            'border-error-500 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive'
        )}
        aria-invalid={invalid}
      >
        <div className='flex items-center justify-between p-2'>
          <Input
            placeholder='搜索'
            className='w-full'
            value={leftSearch}
            onChange={(e) => setLeftSearch(e.target.value)}
          />
          <Button
            type='button'
            variant='outline'
            size='icon'
            onClick={moveToRight}
            className='hover:bg-muted ml-2'
          >
            <ChevronRightIcon className='h-4 w-4' />
          </Button>
        </div>
        <div className='flex items-center gap-2 border-t px-4 py-2'>
          <Checkbox
            id='select-all-left'
            checked={getCheckboxState(filteredLeft)}
            disabled={filteredLeft.length === 0}
            onCheckedChange={(checked) =>
              selectAll(leftList, setLeftList, !!checked)
            }
          />
          <label
            htmlFor='select-all-left'
            className='text-muted-foreground text-sm font-medium'
          >
            {countSelected(filteredLeft)} 项已选
          </label>
        </div>
        <ul className='h-50 overflow-y-auto border-t p-2'>
          {filteredLeft.map((item) => (
            <li
              key={item.key}
              className='hover:bg-muted group flex items-center gap-2 rounded-sm py-1.5 pr-3 pl-2'
            >
              <Checkbox
                checked={item.selected}
                onCheckedChange={() =>
                  toggleSelection(leftList, setLeftList, item.key)
                }
              />
              <span className='text-sm'>{item.label}</span>
            </li>
          ))}
        </ul>
      </div>

      <div
        className={cn(
          'w-1/2 rounded-md border shadow',
          invalid &&
            'border-error-500 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive'
        )}
        aria-invalid={invalid}
      >
        <div className='flex items-center justify-between p-2'>
          <Button
            type='button'
            variant='outline'
            size='icon'
            onClick={moveToLeft}
            className='hover:bg-muted mr-2'
          >
            <ChevronLeftIcon className='h-4 w-4' />
          </Button>
          <Input
            placeholder='搜索'
            className='w-full'
            value={rightSearch}
            onChange={(e) => setRightSearch(e.target.value)}
          />
        </div>
        <div className='flex items-center gap-2 border-t px-4 py-2'>
          <Checkbox
            id='select-all-right'
            checked={getCheckboxState(filteredRight)}
            disabled={filteredRight.length === 0}
            onCheckedChange={(checked) =>
              selectAll(rightList, setRightList, !!checked)
            }
          />
          <label
            htmlFor='select-all-right'
            className='text-muted-foreground text-sm font-medium'
          >
            {countSelected(filteredRight)} 项已选
          </label>
        </div>
        <ul className='h-50 overflow-y-auto border-t p-2'>
          {filteredRight.map((item) => (
            <li
              key={item.key}
              className='hover:bg-muted group flex items-center gap-2 rounded-sm py-1.5 pr-3 pl-2'
            >
              <Checkbox
                checked={item.selected}
                onCheckedChange={() =>
                  toggleSelection(rightList, setRightList, item.key)
                }
              />
              <span className='text-sm'>{item.label}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
