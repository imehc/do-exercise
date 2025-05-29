import { TreeSelectComponentProps } from '.'
import { TreeSelectContext } from './context'
import { TreeNode } from './tree-node'
import {
  getValuesFromState,
  handleNodeCheck,
  makeTreeNodeDataState,
} from './tree-utils'
import type { TreeNodeDataState, TreeSelectProps } from './types'

type TreeViewProps = TreeSelectProps &
  Pick<TreeSelectComponentProps, 'readonly'> & {
    searchValue?: string
  }

export const TreeView = ({
  value,
  onChange,
  data,
  searchValue,
  multiple,
  valueMode = 'all',
  readonly,
}: TreeViewProps) => {
  const dataState = makeTreeNodeDataState(data, value, searchValue)
  const handleCheck = (node: TreeNodeDataState) => {
    const updatedState = handleNodeCheck(dataState, node)
    if (!multiple) {
      // 在单选模式下，根据当前点击的动作（选中/取消选中）来决定返回值
      if (!node.checked) {
        // 如果当前是未选中状态，说明这次操作是选中操作
        onChange?.([node.value])
      } else {
        // 如果当前是选中状态，说明这次操作是取消选中操作
        onChange?.([])
      }
      return
    }
    const values = getValuesFromState(updatedState, valueMode)
    onChange?.(values)
  }

  return (
    <TreeSelectContext.Provider
      value={{
        onCheck: (node) => {
          if (readonly) return
          handleCheck(node)
        },
      }}
    >
      <div className='flex flex-col gap-2'>
        {dataState.map((item) => {
          if (!item.visible) return null

          return <TreeNode key={item.value} data={item} readonly={readonly} />
        })}
      </div>
    </TreeSelectContext.Provider>
  )
}
