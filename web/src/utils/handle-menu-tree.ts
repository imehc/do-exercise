import { Menu, MenuType, SysMenuTree } from '~/do-exercise-api'

type MenuTypeKey = 'directory' | 'menu' | 'button'

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

export function handleToMenuTree(menus: Menu[]) {
  const menuMap = new Map<number, SysMenuTree>()
  const result: SysMenuTree[] = []

  // 初始化 menuMap
  menus.forEach((menu) => {
    menuMap.set(menu.id, { ...menu, children: [] })
  })

  // 构建树结构
  menus.forEach((menu) => {
    const node = menuMap.get(menu.id)!
    if (menu.parentId === 0 || !menuMap.has(menu.parentId)) {
      result.push(node)
    } else {
      const parent = menuMap.get(menu.parentId)
      if (parent) {
        if (!parent.children) {
          parent.children = []
        }
        parent.children.push(node)
      }
    }
  })

  // 按照 type 分组
  const grouped = result.reduce<Record<number, SysMenuTree[]>>((acc, cur) => {
    if (!acc[cur.type]) {
      acc[cur.type] = []
    }
    acc[cur.type].push(cur)
    return acc
  }, {})

  // 映射为 'directory' | 'menu' | 'button'
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
