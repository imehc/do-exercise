import { describe, expect, it } from 'vitest'
import { MenuType, type SysMenuTree } from '~/do-exercise-api'
import {
  collectSelectedApis,
  collectSelectedPermissions,
  diffMenuSelection,
  findMissingAncestors,
  findPagesWithoutActions,
  indexMenuTree,
} from './permission-diff'

const tree: SysMenuTree[] = [
  {
    id: 1,
    name: '系统管理',
    parentId: 0,
    type: MenuType.directory,
    children: [
      {
        id: 5,
        name: '用户管理',
        parentId: 1,
        type: MenuType.menu,
        route: '/user',
        apis: [{ id: 90, method: 'GET', path: '/user' }],
        children: [
          {
            id: 131,
            name: '查询',
            parentId: 5,
            type: MenuType.button,
            permission: 'user:query',
            apis: [
              { id: 90, method: 'GET', path: '/user' },
              { id: 91, method: 'GET', path: '/user/{id}' },
            ],
            children: [],
          },
          {
            id: 134,
            name: '删除',
            parentId: 5,
            type: MenuType.button,
            permission: 'user:delete',
            apis: [{ id: 92, method: 'DELETE', path: '/user/{id}' }],
            children: [],
          },
        ],
      },
    ],
  },
] as unknown as SysMenuTree[]

const index = indexMenuTree(tree, (node) => node.name)

describe('indexMenuTree', () => {
  it('压平全部层级并把根节点的 parentId 0 归一成 undefined', () => {
    expect([...index.keys()].sort((a, b) => a - b)).toEqual([1, 5, 131, 134])
    expect(index.get(1)?.parentId).toBeUndefined()
    expect(index.get(131)?.parentId).toBe(5)
  })
})

describe('findMissingAncestors', () => {
  it('只勾按钮时列出整条祖先链', () => {
    expect(findMissingAncestors([131], index).map((item) => item.id)).toEqual([
      5, 1,
    ])
  })

  it('勾了父级但漏了祖父级也要报出来', () => {
    expect(findMissingAncestors([131, 5], index).map((item) => item.id)).toEqual(
      [1]
    )
  })

  it('整条链都勾上时没有提示', () => {
    expect(findMissingAncestors([1, 5, 131], index)).toEqual([])
  })
})

describe('findPagesWithoutActions', () => {
  it('勾了页面但一个按钮都没勾时提示', () => {
    expect(findPagesWithoutActions([1, 5], index).map((item) => item.id)).toEqual(
      [5]
    )
  })

  it('勾了任意一个按钮就不再提示', () => {
    expect(findPagesWithoutActions([1, 5, 131], index)).toEqual([])
  })

  it('没勾页面本身时不提示（那是缺祖先，由 findMissingAncestors 负责）', () => {
    expect(findPagesWithoutActions([1], index)).toEqual([])
  })

  it('没有按钮子节点的展示页不算问题', () => {
    const displayOnly = indexMenuTree(
      [
        {
          id: 7,
          name: '系统信息',
          parentId: 0,
          type: MenuType.menu,
          children: [],
        },
      ] as unknown as SysMenuTree[],
      (node) => node.name
    )
    expect(findPagesWithoutActions([7], displayOnly)).toEqual([])
  })
})

describe('collectSelectedPermissions', () => {
  it('只取按钮且按权限标识排序', () => {
    expect(
      collectSelectedPermissions([134, 5, 131, 1], index).map(
        (item) => item.permission
      )
    ).toEqual(['user:delete', 'user:query'])
  })
})

describe('collectSelectedApis', () => {
  it('按 method + path 去重', () => {
    const apis = collectSelectedApis([5, 131], index)
    expect(apis.map((api) => `${api.method} ${api.path}`)).toEqual([
      'GET /user',
      'GET /user/{id}',
    ])
  })
})

describe('diffMenuSelection', () => {
  it('区分新增与移除', () => {
    const diff = diffMenuSelection([1, 5, 131], [1, 5, 134], index, () => '?')
    expect(diff.added.map((item) => item.id)).toEqual([134])
    expect(diff.removed.map((item) => item.id)).toEqual([131])
  })

  it('菜单树里已经不存在的旧授权也要显示为移除项', () => {
    const diff = diffMenuSelection([999], [], index, (id) => `#${id}`)
    expect(diff.removed).toHaveLength(1)
    expect(diff.removed[0]?.label).toBe('#999')
  })
})
