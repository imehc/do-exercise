import { MenuType, type SysMenuApiBrief, type SysMenuTree } from '~/do-exercise-api'

/**
 * 授权页需要的菜单摘要。
 *
 * 角色授权本质是「勾一批菜单 ID」，但授权人真正要判断的是这批 ID 放开了哪些
 * 权限标识、打通了哪些接口、和现状相比多给/收回了什么。这些都能从菜单树推出来，
 * 所以这里只做纯函数，交由组件渲染，方便单测覆盖。
 */
export interface MenuNodeSummary {
  id: number
  /** 已翻译的展示名，由调用方注入（纯函数不碰 i18n 运行时） */
  label: string
  parentId?: number
  type: MenuType
  permission?: string
  apis: SysMenuApiBrief[]
}

export type MenuIndex = Map<number, MenuNodeSummary>

/** 把菜单树压平成 id -> 摘要，后续所有推导都基于这张索引。 */
export function indexMenuTree(
  nodes: SysMenuTree[],
  getLabel: (node: SysMenuTree) => string
): MenuIndex {
  const index: MenuIndex = new Map()
  const walk = (list: SysMenuTree[]) => {
    for (const node of list) {
      index.set(node.id, {
        id: node.id,
        label: getLabel(node),
        // 后端根节点用 0 表示无父级，这里统一成 undefined
        parentId: node.parentId ? node.parentId : undefined,
        type: node.type,
        permission: node.permission,
        apis: node.apis ?? [],
      })
      if (node.children?.length) walk(node.children)
    }
  }
  walk(nodes)
  return index
}

/**
 * 父子一致性：勾了按钮/页面却漏勾祖先节点。
 *
 * 后端按菜单绑定的 API 下策略，父级页面没授权时按钮所在的页面根本进不去，
 * 保存后表现为「权限给了但界面看不到」，是最容易踩的坑，因此在提交前显式列出。
 * 返回的是缺失的祖先节点，按索引顺序去重。
 */
export function findMissingAncestors(
  selectedIds: readonly number[],
  index: MenuIndex
): MenuNodeSummary[] {
  const selected = new Set(selectedIds)
  const missing = new Map<number, MenuNodeSummary>()
  for (const id of selected) {
    let parentId = index.get(id)?.parentId
    // 沿链上溯，遇到已勾选的祖先也要继续走：可能是「勾了父、漏了祖父」
    while (parentId !== undefined) {
      const parent = index.get(parentId)
      if (!parent) break
      if (!selected.has(parent.id)) missing.set(parent.id, parent)
      parentId = parent.parentId
    }
  }
  return [...missing.values()]
}

/**
 * 父子一致性的另一半：勾了页面却没勾该页面下的任何操作按钮。
 *
 * 这种角色打开页面后什么都干不了（连查询按钮的接口都没放开），属于「配错了但
 * 不会报错」的一类，只有用起来才发现。本来就没有按钮的纯展示页不算问题，
 * 因此只在「有操作可选、却一个都没选」时提示。
 */
export function findPagesWithoutActions(
  selectedIds: readonly number[],
  index: MenuIndex
): MenuNodeSummary[] {
  const selected = new Set(selectedIds)
  const buttonsByParent = new Map<number, MenuNodeSummary[]>()
  for (const node of index.values()) {
    if (node.type !== MenuType.button || node.parentId === undefined) continue
    const siblings = buttonsByParent.get(node.parentId)
    if (siblings) siblings.push(node)
    else buttonsByParent.set(node.parentId, [node])
  }

  const flagged = new Map<number, MenuNodeSummary>()
  for (const id of selected) {
    const node = index.get(id)
    if (!node || node.type !== MenuType.menu) continue
    const buttons = buttonsByParent.get(node.id)
    if (!buttons?.length) continue
    if (!buttons.some((button) => selected.has(button.id))) {
      flagged.set(node.id, node)
    }
  }
  return [...flagged.values()]
}

/** 勾选范围内带权限标识的按钮，用于「这次授权放开了哪些权限」预览。 */
export function collectSelectedPermissions(
  selectedIds: readonly number[],
  index: MenuIndex
): MenuNodeSummary[] {
  return selectedIds
    .map((id) => index.get(id))
    .filter(
      (node): node is MenuNodeSummary =>
        !!node && node.type === MenuType.button && !!node.permission
    )
    .sort((a, b) => (a.permission ?? '').localeCompare(b.permission ?? ''))
}

/**
 * 勾选范围实际放开的接口集合（按 method + path 去重）。
 * 多个按钮共用同一个接口很常见，去重后才看得出真实的攻击面大小。
 */
export function collectSelectedApis(
  selectedIds: readonly number[],
  index: MenuIndex
): SysMenuApiBrief[] {
  const seen = new Map<string, SysMenuApiBrief>()
  for (const id of selectedIds) {
    for (const api of index.get(id)?.apis ?? []) {
      seen.set(`${api.method} ${api.path}`, api)
    }
  }
  return [...seen.values()].sort((a, b) =>
    `${a.path} ${a.method}`.localeCompare(`${b.path} ${b.method}`)
  )
}

export interface MenuSelectionDiff {
  added: MenuNodeSummary[]
  removed: MenuNodeSummary[]
}

/**
 * 变更影响面：与角色当前已授权的菜单相比新增/移除了什么。
 * 移除项可能已经不在菜单树里（菜单被删或改成平台专属），此时用兜底摘要保证
 * 「你正在收回这一条」这件事仍然可见，而不是静默消失。
 */
export function diffMenuSelection(
  before: readonly number[],
  after: readonly number[],
  index: MenuIndex,
  fallbackLabel: (id: number) => string
): MenuSelectionDiff {
  const beforeSet = new Set(before)
  const afterSet = new Set(after)
  const resolve = (id: number): MenuNodeSummary =>
    index.get(id) ?? {
      id,
      label: fallbackLabel(id),
      type: MenuType.menu,
      apis: [],
    }
  return {
    added: after.filter((id) => !beforeSet.has(id)).map(resolve),
    removed: before.filter((id) => !afterSet.has(id)).map(resolve),
  }
}
