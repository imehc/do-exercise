import { Button } from '~/components/ui/button'
import { Tabs, TabsContent, TabsList } from '~/components/ui/tabs'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Suspense, useCallback, useMemo, useRef } from 'react'
import { useRouter, useTabs } from '~/hooks'
import { KeepAliveTabContext, type KeepAliveTabContextType } from '~/provider'
import { useLocation, useOutlet } from 'react-router'
import { KeepAlive, useKeepAliveRef } from 'keepalive-for-react'
import { Loading } from '..'
import {
  DndContext,
  DragEndEvent,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors
} from '@dnd-kit/core'
import { SortableContext, horizontalListSortingStrategy } from '@dnd-kit/sortable'
import { SortableTab } from './sortable-tab'

/**
 * 站点标签页组件，提供多标签页管理功能
 * 功能包括：
 * - 标签页切换
 * - 右键菜单操作（刷新、关闭、关闭其他）
 * - KeepAlive缓存管理
 * - 标签页滚动控制
 */
export function SiteTabs() {
  // 获取当前标签页状态和操作函数
  const { activeTabRoutePath, tabs, ...handler } = useTabs()
  // 路由实例
  const router = useRouter()

  const location = useLocation()
  const aliveRef = useKeepAliveRef()
  const outlet = useOutlet()

  /**
   * 关闭标签页
   * @param path 可选，指定要关闭的标签页路径，默认为当前激活标签页
   */
  const closeTab = useCallback(
    (path?: string) => {
      handler.closeTab(path)
      aliveRef.current?.destroy() // 销毁对应的KeepAlive缓存
    },
    [handler, aliveRef]
  )
  /**
   * 关闭其他标签页
   * @param path 可选，指定要保留的标签页路径，默认为当前激活标签页
   */
  const closeOtherTab = useCallback(
    (path?: string) => {
      handler.closeOtherTab(path)
      aliveRef.current?.destroyOther() // 销毁其他KeepAlive缓存
    },
    [handler, aliveRef]
  )
  /**
   * 刷新标签页
   * @param path 可选，指定要刷新的标签页路径，默认为当前激活标签页
   */
  const refreshTab = useCallback(
    (path?: string) => {
      handler.refreshTab(path)
      aliveRef.current?.refresh() // 刷新KeepAlive缓存
    },
    [handler, aliveRef]
  )

  const keepAliveContextValue = useMemo<KeepAliveTabContextType>(
    () => ({
      closeTab,
      closeOtherTab,
      refreshTab
    }),
    [closeTab, closeOtherTab, refreshTab]
  )

  const currentCacheKey = useMemo(() => {
    return location.pathname + location.search
  }, [location.pathname, location.search])

  const tabsListRef = useRef<HTMLDivElement>(null)

  // 配置拖拽传感器
  const mouseSensor = useSensor(MouseSensor, {
    activationConstraint: {
      distance: 10 // 最小拖拽距离，防止误触
    }
  })
  const touchSensor = useSensor(TouchSensor, {
    activationConstraint: {
      delay: 250, // 触摸延迟，防止误触
      tolerance: 5 // 触摸容差
    }
  })
  const sensors = useSensors(mouseSensor, touchSensor)

  /**
   * 处理标签页拖拽结束
   */
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return

    const oldIndex = tabs.findIndex(tab => tab.id === active.id)
    const newIndex = tabs.findIndex(tab => tab.id === over.id)

    const newTabs = [...tabs]
    const [draggedItem] = newTabs.splice(oldIndex, 1)
    newTabs.splice(newIndex, 0, draggedItem)
    handler.setTabs(newTabs)
  }

  /**
   * 处理标签页滚动
   * @param direction 滚动方向 'left'或'right'
   */
  const handleScroll = (direction: 'left' | 'right') => {
    if (!tabsListRef.current) return
    const tabWidth = 100 // 估算每个标签的宽度
    const currentScroll = tabsListRef.current.scrollLeft
    tabsListRef.current.scrollTo({
      left: direction === 'left' ? currentScroll - tabWidth : currentScroll + tabWidth,
      behavior: 'smooth'
    })
  }

  /**
   * 处理标签页双击滚动
   * @param direction 滚动方向 'left'或'right'
   * - 向左滚动到最左端
   * - 向右滚动到最右端
   */
  const handleDoubleClick = (direction: 'left' | 'right') => {
    if (!tabsListRef.current) return
    tabsListRef.current.scrollTo({
      left: direction === 'left' ? 0 : tabsListRef.current.scrollWidth,
      behavior: 'smooth'
    })
  }

  return (
    <KeepAliveTabContext.Provider value={keepAliveContextValue}>
      <Tabs value={activeTabRoutePath} className="w-full h-full my-2 px-4">
        <div className="relative flex items-center">
          <Button
            variant="ghost"
            size="icon"
            className="absolute left-0 z-10 h-8 w-4 rounded-full bg-background/80 backdrop-blur-sm"
            onClick={() => handleScroll('left')}
            onDoubleClick={() => handleDoubleClick('left')}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
            <TabsList
              ref={tabsListRef}
              className="flex justify-start w-full gap-x-2 p-1 overflow-hidden flex-nowrap mx-4"
            >
              <SortableContext items={tabs} strategy={horizontalListSortingStrategy}>
                {tabs.map(tab => (
                  <SortableTab
                    key={tab.id}
                    tab={tab}
                    isActive={tab.routePath === activeTabRoutePath}
                    onClick={() => {
                      if (tab.routePath === activeTabRoutePath) return
                      router.push(tab.routePath)
                    }}
                    onContextMenu={e => {
                      if (tab.routePath !== activeTabRoutePath) {
                        e.preventDefault()
                        e.stopPropagation()
                      }
                    }}
                    onClose={() => closeTab(tab.routePath)}
                    onRefresh={() => refreshTab(tab.routePath)}
                    onCloseOther={() => closeOtherTab(tab.routePath)}
                    showCloseOptions={tabs.length > 1}
                  />
                ))}
              </SortableContext>
            </TabsList>
          </DndContext>
          <Button
            variant="ghost"
            size="icon"
            className="absolute right-0 z-10 h-8 w-4 rounded-full bg-background/80 backdrop-blur-sm"
            onClick={() => handleScroll('right')}
            onDoubleClick={() => handleDoubleClick('right')}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
        <TabsContent value={activeTabRoutePath} className="h-full">
          <KeepAlive transition aliveRef={aliveRef} activeCacheKey={currentCacheKey} max={10}>
            <Suspense fallback={<Loading />}>{outlet}</Suspense>
          </KeepAlive>
        </TabsContent>
      </Tabs>
    </KeepAliveTabContext.Provider>
  )
}
