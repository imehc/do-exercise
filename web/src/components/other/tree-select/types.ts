export interface TreeNodeData {
  name: string
  value: string | number
  children?: Array<TreeNodeData>
}

export interface TreeNodeDataState extends TreeNodeData {
  parent?: TreeNodeDataState
  checked?: boolean
  hasSelectedChildren?: boolean
  visible?: boolean
  children?: Array<TreeNodeDataState>
}

export type TreeSelectProps = {
  value: Array<string | number>
  onChange?: (value: Array<string | number>) => void
  data: Array<TreeNodeData>
  multiple?: boolean
  placeholder?: string
  /** @default all */
  valueMode?: 'all' | 'parent-only'
}
