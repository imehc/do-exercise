import { MenuType, SysMenuTree } from '~/do-exercise-api'

export function handleMenuTree(menuTree: SysMenuTree[]) {
  const result: SysMenuTree[] = []
  function flatten(node: SysMenuTree) {
    result.push(node)
    if (node.children && node.children.length > 0) {
      node.children.forEach(flatten)
    }
  }
  menuTree.forEach(flatten)

  const grouped = result.reduce<Record<number, SysMenuTree[]>>((acc, cur) => {
    if (!acc[cur.type]) {
      acc[cur.type] = []
    }
    acc[cur.type].push(cur)
    return acc
  }, {})

  type MenuTypeKey = 'directory' | 'menu' | 'button'
  const remapped = Object.entries(grouped).reduce<
    Record<MenuTypeKey, SysMenuTree[]>
  >(
    (acc, [key, value]) => {
      const label = Object.keys(MenuType).find(
        (k) => MenuType[k as keyof typeof MenuType].toString() === key
      )
      if (label) {
        acc[label as MenuTypeKey] = value
      }
      return acc
    },
    {} as Record<MenuTypeKey, SysMenuTree[]>
  )

  return remapped
}
