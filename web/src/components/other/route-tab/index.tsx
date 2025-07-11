import { useState, useEffect } from 'react'
import { useRouter, useLocation } from '@tanstack/react-router'
import { Trans } from '@lingui/react/macro'
import { useAtom } from 'jotai'
import AutoSizer from 'react-virtualized-auto-sizer'
import { routeTabCacheAtom } from '~/atoms'
import { Menu } from '~/do-exercise-api'
import { useMenus } from '~/provider'
import { cn } from '~/lib/utils'
import { RightMenu } from './right-menu'
import { Tab } from './tab'

export type TabItem = {
  id: string
  name: React.ReactNode
  route: string
  closable: boolean
  icon: string
}

// 固定页面配置
const FIXED_PAGES: Record<'dashboard' | 'settings', TabItem> = {
  dashboard: {
    id: 'dashboard',
    name: <Trans>首页</Trans>,
    route: '/',
    closable: false,
    icon: 'Dashboard',
  },
  settings: {
    id: 'settings',
    name: <Trans>设置</Trans>,
    route: '/settings',
    closable: true,
    icon: 'Settings',
  },
}

export function RouteTab() {
  const router = useRouter()
  const location = useLocation()
  const [tabs, setTabs] = useState<TabItem[]>([FIXED_PAGES.dashboard])
  const [activeTab, setActiveTab] = useState(FIXED_PAGES.dashboard.id)
  const [cacheTab, setCacheTab] = useAtom(routeTabCacheAtom)
  const originMenus = useMenus()

  useEffect(() => {
    if (!cacheTab.tabs.length) {
      setTabs([FIXED_PAGES.dashboard])
      return
    }
    const tempCacheTabs = cacheTab.tabs
      .filter(
        (item) =>
          originMenus.some((menu) => menu.id.toString() === item.id) ||
          [FIXED_PAGES.dashboard.id, FIXED_PAGES.settings.id].includes(item.id)
      )
      .map((item) => {
        if (item.id === FIXED_PAGES.dashboard.id) {
          return {
            ...item,
            name: FIXED_PAGES.dashboard.name,
            route: FIXED_PAGES.dashboard.route,
            icon: FIXED_PAGES.dashboard.icon,
          } as TabItem
        }
        if (item.id === FIXED_PAGES.settings.id) {
          return {
            ...item,
            name: FIXED_PAGES.settings.name,
            route: FIXED_PAGES.settings.route,
            icon: FIXED_PAGES.settings.icon,
          } as TabItem
        }
        const menu = originMenus.find(
          (i) => i.id.toString() === item.id
        ) as Menu
        return {
          ...item,
          name: menu.name,
          icon: menu.icon,
          route: menu.route,
        } as TabItem
      })
    setTabs(tempCacheTabs)

    return () => {
      setTabs([])
    }
  }, [cacheTab.tabs, cacheTab.tabs.length, originMenus])

  useEffect(() => {
    const handler = () => {
      setCacheTab({ activeTab, tabs })
    }

    window.addEventListener('beforeunload', handler)

    return () => {
      window.removeEventListener('beforeunload', handler)
    }
  }, [activeTab, setCacheTab, tabs])

  // 监听路由变化，自动添加标签页
  useEffect(() => {
    const currentPath = location.pathname

    // 检查是否是固定页面
    const fixedPage = Object.values(FIXED_PAGES).find(
      (page) => page.route === currentPath
    )
    if (fixedPage) {
      if (!tabs.find((tab) => tab.id.toString() === fixedPage.id)) {
        setTabs((prev) => {
          const filtered = prev.filter((tab) => tab.id !== fixedPage.id)

          // Dashboard 始终在第一位
          if (fixedPage.id === 'dashboard') {
            return [fixedPage, ...filtered]
          }

          // 其他固定页面（如 Settings）添加到最后
          return [...filtered, fixedPage]
        })
      }
      setActiveTab(fixedPage.id)
      return
    }

    // 检查是否是菜单页面
    const menuTab = originMenus.find((tab) => tab.route === currentPath)
    if (menuTab) {
      setTabs((prev) => {
        const exists = prev.find((tab) => tab.id === menuTab.id.toString())
        if (!exists) {
          return [
            ...prev,
            {
              id: menuTab.id.toString(),
              name: menuTab.name!,
              route: menuTab.route!,
              closable: true,
              icon: menuTab.icon!,
            } satisfies TabItem,
          ]
        }
        return prev
      })
      setActiveTab(menuTab.id.toString())
    }
  }, [location.pathname, originMenus, tabs])

  const handleTabClose = (id: string) => {
    // 不能删除固定页面
    const tab = tabs.find((t) => t.id === id)
    if (!tab?.closable) return

    const currentIndex = tabs.findIndex((tab) => tab.id === id)
    const isCurrentTab = activeTab === id

    // 如果关闭的是当前选中的tab，需要选中下一个或上一个tab
    if (isCurrentTab) {
      const nextTab = tabs[currentIndex + 1] || tabs[currentIndex - 1]
      if (nextTab) {
        setActiveTab(nextTab.id)
        // 导航到新的tab路由
        router.navigate({ to: nextTab.route })
      }
    }

    // 删除tab
    setTabs((prev) => prev.filter((item) => item.id !== id))
  }

  const handleTabRefresh = (id: string) => {
    const tab = tabs.find((t) => t.id === id)
    if (tab) {
      // 刷新当前路由
      router.invalidate()
    }
  }

  const handleTabClick = (id: string) => {
    const tab = tabs.find((t) => t.id === id)
    if (tab) {
      setActiveTab(id)
      router.navigate({ to: tab.route })
    }
  }

  const handleCloseLeftTabs = (id: string) => {
    const currentIndex = tabs.findIndex((tab) => tab.id === id)

    // 获取要关闭的左侧tabs，但不包括首页
    const leftTabs = tabs.slice(0, currentIndex)
    const tabsToClose = leftTabs.filter((tab) => tab.id !== 'dashboard')

    if (tabsToClose.length === 0) return

    // 如果当前选中的tab在要关闭的tabs中，选中当前tab
    const closedTabIds = tabsToClose.map((tab) => tab.id)
    if (closedTabIds.includes(activeTab)) {
      setActiveTab(id)
    }

    setTabs((prev) => prev.filter((tab) => !closedTabIds.includes(tab.id)))
  }

  const handleCloseRightTabs = (id: string) => {
    const currentIndex = tabs.findIndex((tab) => tab.id === id)
    const rightTabs = tabs.slice(currentIndex + 1)
    const tabsToClose = rightTabs.filter((tab) => tab.id !== 'dashboard')

    if (tabsToClose.length === 0) return

    // 如果当前选中的tab在要关闭的tabs中，选中当前tab
    const closedTabIds = tabsToClose.map((tab) => tab.id)
    if (closedTabIds.includes(activeTab)) {
      setActiveTab(id)
      const currentTab = tabs.find((t) => t.id === id)
      if (currentTab) {
        router.navigate({ to: currentTab.route })
      }
    }

    setTabs((prev) => prev.filter((tab) => !closedTabIds.includes(tab.id)))
  }

  const handleCloseOtherTabs = (id: string) => {
    // 保留首页和当前tab
    const tabsToKeep = tabs.filter(
      (tab) => tab.id === 'dashboard' || tab.id === id
    )

    if (tabsToKeep.length === tabs.length) return

    // 如果当前选中的tab不是要保留的tab，选中当前tab
    if (activeTab !== id && activeTab !== 'dashboard') {
      setActiveTab(id)
      const currentTab = tabs.find((t) => t.id === id)
      if (currentTab) {
        router.navigate({ to: currentTab.route })
      }
    }

    setTabs(tabsToKeep)
  }

  const handleCloseAllTabs = () => {
    // 只保留首页
    const fixedTabs = tabs.filter((tab) => tab.id === 'dashboard')

    if (fixedTabs.length === tabs.length) return

    setTabs(fixedTabs)

    // 如果当前选中的tab被关闭了，选中首页
    if (activeTab !== 'dashboard') {
      const firstFixedTab = fixedTabs[0]
      if (firstFixedTab) {
        setActiveTab(firstFixedTab.id)
        router.navigate({ to: firstFixedTab.route })
      }
    }
  }

  return (
    <AutoSizer>
      {({ width }) => (
        <div
          className={cn(
            'border-border bg-background/80 grid w-full min-w-0 grid-cols-[1fr_auto]'
          )}
          style={{ width }}
        >
          {/* Tab 区域 */}
          <Tab
            tabs={tabs}
            activeTab={activeTab}
            onSort={setTabs}
            onClick={handleTabClick}
            onClose={handleTabClose}
            onRefresh={handleTabRefresh}
            onCloseLeft={handleCloseLeftTabs}
            onCloseRight={handleCloseRightTabs}
            onCloseOther={handleCloseOtherTabs}
          />
          {/* 右侧操作区 */}
          <RightMenu
            activeTab={activeTab}
            totalTabs={tabs.length}
            onRefresh={handleTabRefresh}
            onClose={handleTabClose}
            onCloseAll={handleCloseAllTabs}
          />
        </div>
      )}
    </AutoSizer>
  )
}
