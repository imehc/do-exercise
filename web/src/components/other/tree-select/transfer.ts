import { TreeNodeData } from './types'

export function transformData<T extends { children?: T[] }>(
  data: T[],
  getName: (item: T) => string,
  getValue: (item: T) => string | number
): TreeNodeData[] {
  return data.map((item) => {
    const name = getName(item)
    const value = getValue(item)

    const children = item.children
      ? transformData(item.children, getName, getValue)
      : undefined
    return { name, value, children }
  })
}
