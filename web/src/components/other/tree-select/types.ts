export interface TreeNodeData {
  name: string
  value: string
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
  value: Array<string>
  onChange: (value: Array<string>) => void
  data: Array<TreeNodeData>
  multiple?: boolean
  placeholder?: string
}
