import { useState } from 'react'
import AutoSizer from 'react-virtualized-auto-sizer'
import { cn } from '~/lib/utils'
import { menus } from './mock'
import { RightMenu } from './right-menu'
import { Tab } from './tab'

export function RouteTab() {
  const [tabs, setTabs] = useState(menus)
  const [activeTab, setActiveTab] = useState(menus[0]?.id)
  
  const handleTabClose = (id: string) => {
    // 不能删除第一个tab
    if (id === menus[0]?.id) return
    
    const currentIndex = tabs.findIndex((tab) => tab.id === id)
    const isCurrentTab = activeTab === id
    
    // 如果关闭的是当前选中的tab，需要选中下一个或上一个tab
    if (isCurrentTab) {
      const nextTab = tabs[currentIndex + 1] || tabs[currentIndex - 1]
      if (nextTab) {
        setActiveTab(nextTab.id)
      }
    }
    
    // 删除tab
    setTabs((prev) => prev.filter((item) => item.id !== id))
  }
  
  const handleTabRefresh = (id: string) => {
    console.log('refresh tab:', id)
    // TODO: 实现刷新逻辑
  }
  
  const handleCloseLeftTabs = (id: string) => {
    const currentIndex = tabs.findIndex((tab) => tab.id === id)
    const firstTabId = menus[0]?.id
    
    // 获取要关闭的左侧tabs，但不包括第一个tab
    const leftTabs = tabs.slice(0, currentIndex)
    const tabsToClose = leftTabs.filter((tab) => tab.id !== firstTabId)
    
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
    
    if (rightTabs.length === 0) return
    
    // 如果当前选中的tab在要关闭的tabs中，选中当前tab
    const closedTabIds = rightTabs.map((tab) => tab.id)
    if (closedTabIds.includes(activeTab)) {
      setActiveTab(id)
    }
    
    setTabs((prev) => prev.filter((tab) => !closedTabIds.includes(tab.id)))
  }
  
  const handleCloseOtherTabs = (id: string) => {
    const firstTabId = menus[0]?.id
    
    // 保留第一个tab和当前tab
    const tabsToKeep = tabs.filter((tab) => tab.id === firstTabId || tab.id === id)
    
    if (tabsToKeep.length === tabs.length) return
    
    // 如果当前选中的tab不是要保留的tab，选中当前tab
    if (activeTab !== id && activeTab !== firstTabId) {
      setActiveTab(id)
    }
    
    setTabs(tabsToKeep)
  }
  
  const handleCloseAllTabs = () => {
    const firstTabId = menus[0]?.id
    const firstTab = tabs.find((tab) => tab.id === firstTabId)
    
    if (!firstTab) return
    
    // 只保留第一个tab
    setTabs([firstTab])
    setActiveTab(firstTabId)
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
            onClick={setActiveTab}
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
