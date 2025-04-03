import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type SidebarStore = {
  collapsed: boolean
  toggleCollapsed: () => void
  setToggled: (collapsed: boolean) => void
}

export const useSidebarStore = create(
  persist<SidebarStore>(
    set => ({
      collapsed: false,
      toggleCollapsed: () => set(state => ({ collapsed: !state.collapsed })),
      setToggled: (collapsed: boolean) => set(() => ({ collapsed }))
    }),
    {
      name: 'sidebar-storage'
    }
  )
)
