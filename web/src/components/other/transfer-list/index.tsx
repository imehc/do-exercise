import React, { useCallback, useEffect, useMemo, useState } from 'react'
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import { Checkbox } from '~/components/ui/checkbox'
import { Input } from '~/components/ui/input'
import { AriaInvalidProps } from '..'

type Invalid = React.AriaAttributes['aria-invalid']

export type Item = {
  key: number
  label: string
  selected?: boolean
  disabled?: boolean
}

interface Props<T> extends AriaInvalidProps {
  data: T[]
  className?: string
  value?: number[]
  onChange: (value: number[]) => void
  renderLabel?: (item: T) => React.ReactNode
}

interface TransferPanelProps<T> {
  filteredItems: T[]
  search: string
  onSearchChange: (v: string) => void
  onSelectAll: (checked: boolean) => void
  onToggleItem: (key: number) => void
  onMoveClick?: () => void
  moveIcon?: React.ReactNode
  inputFirst?: boolean
  invalid?: Invalid
  checkboxId: string
  renderLabel?: (item: T) => React.ReactNode
}

function getCheckboxState<T extends Item>(list: T[]) {
  const selectable = list.filter((item) => !item.disabled)
  const selectedCount = selectable.filter((item) => item.selected).length
  if (selectedCount === 0) return false
  if (selectedCount === selectable.length) return true
  return 'indeterminate'
}

function countSelected<T extends Item>(list: T[]) {
  return list.filter((item) => item.selected && !item.disabled).length
}

function TransferPanel<T extends Item>({
  filteredItems,
  search,
  onSearchChange,
  onSelectAll,
  onToggleItem,
  onMoveClick,
  moveIcon,
  inputFirst = true,
  invalid,
  checkboxId,
  renderLabel,
}: TransferPanelProps<T>) {
  return (
    <div
      className={cn(
        'w-1/2 rounded-md border shadow',
        invalid &&
          'border-error-500 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive'
      )}
      aria-invalid={invalid}
    >
      <div className='flex items-center justify-between p-2'>
        {inputFirst ? (
          <>
            <Input
              placeholder={t`搜索`}
              className='w-full'
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
            />
            {onMoveClick && (
              <Button
                type='button'
                variant='outline'
                size='icon'
                onClick={onMoveClick}
                className='hover:bg-muted ml-2'
              >
                {moveIcon}
              </Button>
            )}
          </>
        ) : (
          <>
            {onMoveClick && (
              <Button
                type='button'
                variant='outline'
                size='icon'
                onClick={onMoveClick}
                className='hover:bg-muted mr-2'
              >
                {moveIcon}
              </Button>
            )}
            <Input
              placeholder={t`搜索`}
              className='w-full'
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
            />
          </>
        )}
      </div>
      <div className='flex items-center gap-2 border-t px-3 py-2'>
        <Checkbox
          id={checkboxId}
          checked={getCheckboxState(filteredItems)}
          disabled={
            filteredItems.length === 0 ||
            filteredItems.every((item) => item.disabled)
          }
          onCheckedChange={(checked) => onSelectAll(!!checked)}
        />
        <label
          htmlFor={checkboxId}
          className='text-muted-foreground text-sm font-medium'
        >
          <Trans>{countSelected(filteredItems)} 项已选</Trans>
        </label>
      </div>
      <ul className='h-50 overflow-y-auto border-t p-1'>
        {filteredItems.map((item) => (
          <li
            key={item.key}
            className={cn(
              'group flex items-center gap-2 rounded-sm py-1.5 pr-3 pl-2',
              item.disabled ? 'cursor-not-allowed opacity-50' : 'hover:bg-muted'
            )}
          >
            <Checkbox
              id={`transfer-list-checkbox-${item.key}`}
              checked={item.selected}
              disabled={item.disabled}
              onCheckedChange={() => !item.disabled && onToggleItem(item.key)}
            />
            <label
              htmlFor={`transfer-list-checkbox-${item.key}`}
              className={cn(
                'select-none',
                item.disabled ? 'cursor-not-allowed' : 'cursor-pointer'
              )}
            >
              {renderLabel?.(item) ?? (
                <span className='text-sm'>{item.label}</span>
              )}
            </label>
          </li>
        ))}
      </ul>
    </div>
  )
}

export function TransferList<T extends Item>({
  data,
  value = [],
  onChange,
  renderLabel,
  className,
  'aria-invalid': invalid,
}: Props<T>) {
  const [leftList, setLeftList] = useState<T[]>(() =>
    data.filter((item) => !value.includes(item.key))
  )
  const [rightList, setRightList] = useState<T[]>(() =>
    data.filter((item) => value.includes(item.key))
  )
  const [leftSearch, setLeftSearch] = useState('')
  const [rightSearch, setRightSearch] = useState('')

  useEffect(() => {
    onChange(rightList.map((item) => item.key))
  }, [rightList, onChange])

  const toggleSelection = useCallback(
    (
      list: T[],
      setList: React.Dispatch<React.SetStateAction<T[]>>,
      key: number
    ) => {
      setList(
        list.map((item) =>
          item.key === key && !item.disabled
            ? { ...item, selected: !item.selected }
            : item
        )
      )
    },
    []
  )

  const selectAll = useCallback(
    (
      list: T[],
      setList: React.Dispatch<React.SetStateAction<T[]>>,
      selected: boolean
    ) => {
      setList(
        list.map((item) => (item.disabled ? item : { ...item, selected }))
      )
    },
    []
  )

  const moveToRight = useCallback(() => {
    const selected = leftList.filter((item) => item.selected && !item.disabled)
    setRightList([...rightList, ...selected])
    setLeftList(leftList.filter((item) => !item.selected || item.disabled))
  }, [leftList, rightList])

  const moveToLeft = useCallback(() => {
    const selected = rightList.filter((item) => item.selected && !item.disabled)
    setLeftList([...leftList, ...selected])
    setRightList(rightList.filter((item) => !item.selected || item.disabled))
  }, [leftList, rightList])

  const filteredLeft = useMemo(
    () =>
      leftList.filter((item) =>
        item.label.toLowerCase().includes(leftSearch.toLowerCase())
      ),
    [leftList, leftSearch]
  )

  const filteredRight = useMemo(
    () =>
      rightList.filter((item) =>
        item.label.toLowerCase().includes(rightSearch.toLowerCase())
      ),
    [rightList, rightSearch]
  )

  return (
    <div className={cn('flex gap-2', className)}>
      <TransferPanel
        filteredItems={filteredLeft}
        search={leftSearch}
        onSearchChange={setLeftSearch}
        onSelectAll={(checked) => selectAll(leftList, setLeftList, checked)}
        onToggleItem={(key) => toggleSelection(leftList, setLeftList, key)}
        onMoveClick={moveToRight}
        moveIcon={<IconChevronRight className='h-4 w-4' />}
        inputFirst
        invalid={invalid}
        checkboxId='select-all-left'
        renderLabel={renderLabel}
      />
      <TransferPanel
        filteredItems={filteredRight}
        search={rightSearch}
        onSearchChange={setRightSearch}
        onSelectAll={(checked) => selectAll(rightList, setRightList, checked)}
        onToggleItem={(key) => toggleSelection(rightList, setRightList, key)}
        onMoveClick={moveToLeft}
        moveIcon={<IconChevronLeft className='h-4 w-4' />}
        inputFirst={false}
        invalid={invalid}
        checkboxId='select-all-right'
        renderLabel={renderLabel}
      />
    </div>
  )
}
