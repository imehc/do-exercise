import { IconName } from 'lucide-react/dynamic'
import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

export interface Tab {
  label: string
  routePath: string
  id: string
  pathname: string
  icon?: IconName
}

type SidebarStore = {
  /** 是否已完成水合作用（从持久化存储中加载完成） */
  hasHydrated: boolean
  /** 当前打开的标签页列表 */
  tabs: Tab[]
  /** 设置标签页列表 */
  setTabs: (tabs: SidebarStore['tabs']) => void
  /** 添加一个新标签页 */
  appendTab: (tab: Tab) => void
  /** 清空所有标签页 */
  clearTabs: () => void
  /** 内部方法：设置水合状态 */
  _setHasHydrated: (isLoading: SidebarStore['hasHydrated']) => void
}

export const useTabsStore = create(
  persist<SidebarStore>(
    set => ({
      hasHydrated: false,
      tabs: [],
      setTabs: (tabs: SidebarStore['tabs']) => set(() => ({ tabs })),
      appendTab: (tab: Tab) => set(state => ({ tabs: [...state.tabs, tab] })),
      clearTabs: () => set(() => ({ tabs: [] })),
      _setHasHydrated: (hasHydrated: SidebarStore['hasHydrated']) => set(() => ({ hasHydrated }))
    }),
    {
      name: 'tabs-storage',
      storage: createJSONStorage(() => sessionStorage),
      onRehydrateStorage: () => state => {
        state?._setHasHydrated(true)
      }
    }
  )
)
