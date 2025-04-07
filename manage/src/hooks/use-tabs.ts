import { useCallback, useEffect, useMemo, useState } from 'react'
import { useMatchRoute } from '.'
import { IconName } from 'lucide-react/dynamic'
import { useRouter } from '.'
import { type Tab, useTabsStore } from '~/store'
import { useRouterMenus } from '~/provider'

interface KeepAliveTab extends Tab {
  children?: React.ReactNode
}

function getKey() {
  return new Date().getTime().toString()
}

export function useTabs() {
  const { hasHydrated, tabs, appendTab, setTabs } = useTabsStore()
  const { flatMenus } = useRouterMenus()
  // 存放页面记录
  const [keepAliveTabs, setKeepAliveTabs] = useState<KeepAliveTab[]>([])
  // 当前激活的tab
  const [activeTabRoutePath, setActiveTabRoutePath] = useState<string>('')
  // 标记是否正在关闭tab
  const [isClosingTab, setIsClosingTab] = useState(false)

  const matchRoute = useMatchRoute()
  const router = useRouter()

  // 过滤掉不在菜单中的tab
  const filterTabs = useMemo(() => {
    if (!hasHydrated) return
    console.log('tabs', tabs)
    console.log('flatMenus', flatMenus)
    const temp = tabs.filter(tab => flatMenus.some(menu => menu.path === tab.routePath))
    console.log('temp', temp)
    return temp
  }, [hasHydrated, tabs, flatMenus])

  useEffect(() => {
    if (!hasHydrated) return
    if (!filterTabs) return
    setKeepAliveTabs(filterTabs)
  }, [filterTabs, hasHydrated])

  useEffect(() => {
    if (!matchRoute || isClosingTab || !hasHydrated) return
    // 检查当前路由是否在菜单中
    const isRouteInMenu = flatMenus.some(menu => menu.path === matchRoute.routePath)
    if (!isRouteInMenu) return

    const existKeepAliveTab = keepAliveTabs.find(o => o.routePath === matchRoute?.routePath)

    // 如果不存在则需要插入
    if (!existKeepAliveTab) {
      const tab: Tab = {
        label: matchRoute.label,
        id: getKey(),
        routePath: matchRoute.routePath,
        pathname: matchRoute.pathname
      }
      setKeepAliveTabs(prev => [
        ...prev,
        {
          ...tab,
          icon: matchRoute.icon as IconName
        }
      ])
      appendTab(tab)
    }
    setActiveTabRoutePath(matchRoute.routePath)
  }, [hasHydrated, matchRoute, keepAliveTabs, appendTab, isClosingTab])

  const closeTab = useCallback(
    (routePath = activeTabRoutePath) => {
      setIsClosingTab(true)
      const index = keepAliveTabs.findIndex(o => o.routePath === routePath)
      if (index === -1) return

      const newTabs = [...keepAliveTabs]
      newTabs.splice(index, 1)

      if (keepAliveTabs[index]?.routePath === activeTabRoutePath) {
        let nextTab: KeepAliveTab | undefined
        if (newTabs.length > 0) {
          nextTab = index > 0 ? newTabs[index - 1] : newTabs[0]
        }

        if (nextTab) {
          setActiveTabRoutePath(nextTab.routePath)
          router.push(nextTab.pathname)
        } else {
          // 如果没有其他标签页，跳转到默认路由
          router.push('/')
        }
      }

      setKeepAliveTabs(newTabs)
      setTabs(newTabs)
      // 在路由跳转完成后重置标志位
      setTimeout(() => setIsClosingTab(false), 0)
    },
    [router, keepAliveTabs, activeTabRoutePath, setTabs]
  )

  const closeOtherTab = useCallback(
    (routePath = activeTabRoutePath) => {
      const newTabs = keepAliveTabs.filter(o => o.routePath === routePath)
      if (newTabs.length === 0) {
        router.push('/')
      } else if (!newTabs.some(tab => tab.routePath === activeTabRoutePath)) {
        setActiveTabRoutePath(newTabs[0].routePath)
        router.push(newTabs[0].pathname)
      }
      setKeepAliveTabs(newTabs)
      setTabs(newTabs)
    },
    [router, keepAliveTabs, activeTabRoutePath, setTabs]
  )

  const refreshTab = useCallback((routePath = activeTabRoutePath) => {
    setKeepAliveTabs(prev => {
      const index = prev.findIndex(o => o.routePath === routePath)
      if (index >= 0) {
        prev[index].id = getKey()
      }
      return [...prev]
    })
  }, [])

  const handleSetTabs = useCallback((tabs: KeepAliveTab[]) => {
    setKeepAliveTabs(tabs)
    setTabs(tabs.map(({ children, ...props }) => props))
  }, [])

  return {
    tabs: keepAliveTabs,
    setTabs: handleSetTabs,
    activeTabRoutePath,
    closeTab,
    closeOtherTab,
    refreshTab
  }
}
