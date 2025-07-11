import { atomWithStorage } from 'jotai/utils'

// 标签页缓存的类型定义
export type CachedTabItem = {
  id: string
  closable: boolean
}

export type RouteTabCache = {
  tabs: CachedTabItem[]
  activeTab: string
}

// 默认缓存数据
const defaultCache: RouteTabCache = {
  tabs: [
    {
      id: 'dashboard',
      closable: false,
    },
  ],
  activeTab: 'dashboard',
}

// 标签页缓存 atom
export const routeTabCacheAtom = atomWithStorage<RouteTabCache>(
  'routeTabCache',
  defaultCache
)
