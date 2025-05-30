import * as React from 'react'
import { IconChevronDown, IconChevronRight } from '@tabler/icons-react'
import clsx from 'clsx'
import { cn } from '~/lib/utils'
import { Checkbox } from '~/components/ui/checkbox'
import { TreeSelectComponentProps } from '.'
import { TreeSelectContext } from './context'
import { TreeNodeDataState } from './types'

interface TreeNodeProps extends Pick<TreeSelectComponentProps, 'readonly'> {
  data: TreeNodeDataState
}

export const TreeNode = ({ data, readonly }: TreeNodeProps) => {
  const hasSelectedChildren = data.hasSelectedChildren === true
  let checked: boolean | 'indeterminate' = data.checked === true
  if (!checked && hasSelectedChildren) {
    checked = 'indeterminate'
  }

  const [isOpen, setIsOpen] = React.useState(true)
  const { onCheck } = React.useContext(TreeSelectContext)

  if (!data.visible) return null

  const isRoot = data.parent === undefined
  const expandable = data.children !== undefined

  const Icon = isOpen ? IconChevronDown : IconChevronRight
  const hasChildren = data.children && data.children.length > 0

  return (
    <div
      data-expandable={expandable}
      className={cn(
        'flex flex-col gap-2 text-left',
        !isRoot && expandable && '-ml-6'
      )}
    >
      <div className='flex items-center gap-2'>
        {data.children && (
          <button
            aria-hidden={!hasChildren}
            className={clsx(!hasChildren && 'pointer-events-none opacity-0')}
            onClick={hasChildren ? () => setIsOpen(!isOpen) : undefined}
            disabled={readonly}
          >
            <Icon className='expand-icon' size={16} />
          </button>
        )}
        <Checkbox
          id={data.value.toString()}
          checked={checked}
          disabled={readonly}
          onCheckedChange={() => {
            onCheck(data)
          }}
        />
        <label
          className='cursor-pointer text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70'
          htmlFor={data.value.toString()}
        >
          {data.name}
        </label>
      </div>
      {isOpen && hasChildren && (
        <div className={cn('flex flex-col gap-2 pl-12')}>
          {data.children?.map((child) => (
            <TreeNode key={child.value} data={child} readonly={readonly} />
          ))}
        </div>
      )}
    </div>
  )
}
