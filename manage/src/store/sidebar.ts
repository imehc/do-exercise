import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type SidebarStore = {
  collapsed: boolean
  toggleCollapsed: () => void
}

export const useSidebarStore = create(
  persist<SidebarStore>(
    set => ({
      collapsed: false,
      toggleCollapsed: () => set(state => ({ collapsed: !state.collapsed }))
    }),
    {
      name: 'sidebar-storage'
    }
  )
)
