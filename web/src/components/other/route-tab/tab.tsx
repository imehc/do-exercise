import { useRef } from 'react'
import React from 'react'
import {
  IconRefresh,
  IconArrowBarLeft,
  IconArrowBarRight,
  IconX,
  IconMinus,
} from '@tabler/icons-react'
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
} from '@dnd-kit/core'
import { restrictToHorizontalAxis } from '@dnd-kit/modifiers'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  horizontalListSortingStrategy,
  defaultAnimateLayoutChanges,
} from '@dnd-kit/sortable'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Button } from '~/components/ui/button'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '~/components/ui/context-menu'
import { Tabs, TabsList, TabsTrigger } from '../tabs'

export interface TabItem {
  id: string
  icon: React.ReactElement
  title: string
}

interface TabProps {
  tabs: TabItem[]
  activeTab?: string
  onClose?: (id: string) => void
  onClick?: (id: string) => void
  onSort?: (tabs: TabItem[]) => void
  onRefresh?: (id: string) => void
  onCloseLeft?: (id: string) => void
  onCloseRight?: (id: string) => void
  onCloseOther?: (id: string) => void
}

export function Tab({
  tabs,
  activeTab,
  onClose,
  onClick,
  onSort,
  onRefresh,
  onCloseLeft,
  onCloseRight,
  onCloseOther,
}: TabProps) {
  const scrollContainerRef = useRef<HTMLDivElement>(null)
  const [activeId, setActiveId] = React.useState<string | null>(null)
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        delay: 250,
        tolerance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )

  // 使用 useWheel 处理滚轮事件，实现横向滚动
  React.useEffect(() => {
    const container = scrollContainerRef.current
    if (container) {
      // 检查是否有可滚动的内容
      const hasScrollableContent = container.scrollWidth > container.clientWidth

      if (hasScrollableContent) {
        // 阻止默认的垂直滚动行为和事件传播
        const wheelHandler = (event: WheelEvent) => {
          event.preventDefault()
          event.stopPropagation()

          // 将垂直滚动转换为横向滚动
          container.scrollLeft += event.deltaY
        }
        container.addEventListener('wheel', wheelHandler, { passive: false })
        return () => container.removeEventListener('wheel', wheelHandler)
      }
    }
  }, [])

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id.toString())
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (active.id !== over?.id && over) {
      const oldIndex = tabs.findIndex((tab) => tab.id === active.id)
      const newIndex = tabs.findIndex((tab) => tab.id === over.id)

      // 防止第一个tab被拖拽或其他tab被拖拽到第一个位置
      if (oldIndex === 0 || newIndex === 0) {
        setActiveId(null)
        return
      }

      const newTabs = arrayMove(tabs, oldIndex, newIndex)
      onSort?.(newTabs)
    }
    setActiveId(null)
  }

  const activeDragTab = tabs.find((tab) => tab.id === activeId)

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
      modifiers={[restrictToHorizontalAxis]}
    >
      <Tabs
        value={activeTab}
        className='no-scrollbar bg-muted flex overflow-x-auto rounded-lg'
        ref={scrollContainerRef}
      >
        <SortableContext
          items={tabs.map((tab) => tab.id)}
          strategy={horizontalListSortingStrategy}
        >
          <TabsList className='relative flex min-w-full justify-start gap-x-2'>
            {tabs.map((menu, index) => (
              <TabItemComponent
                key={menu.id}
                {...menu}
                isFirst={index === 0}
                isLast={index === tabs.length - 1}
                hasLeftTabs={index > 0}
                hasRightTabs={index < tabs.length - 1}
                totalTabs={tabs.length}
                onClose={onClose}
                onClick={onClick}
                onRefresh={onRefresh}
                onCloseLeft={onCloseLeft}
                onCloseRight={onCloseRight}
                onCloseOther={onCloseOther}
              />
            ))}
          </TabsList>
        </SortableContext>
      </Tabs>
      <DragOverlay>
        {activeId && activeDragTab ? (
          <DragPreview tab={activeDragTab} isActive={activeTab === activeId} />
        ) : null}
      </DragOverlay>
    </DndContext>
  )
}

function DragPreview({ tab, isActive }: { tab: TabItem; isActive: boolean }) {
  return (
    <div
      className={`flex items-center gap-x-1 rounded-md border px-2 py-1 shadow-lg ${
        isActive ? 'bg-background' : 'bg-muted text-muted-foreground'
      }`}
    >
      {tab.icon}
      {tab.title}
    </div>
  )
}

function TabItemComponent({
  id,
  icon,
  title,
  isFirst,
  hasLeftTabs,
  hasRightTabs,
  totalTabs,
  onClose,
  onClick,
  onRefresh,
  onCloseLeft,
  onCloseRight,
  onCloseOther,
}: TabItem & {
  isFirst?: boolean
  isLast?: boolean
  hasLeftTabs?: boolean
  hasRightTabs?: boolean
  totalTabs?: number
  onClose?: (id: string) => void
  onClick?: (id: string) => void
  onRefresh?: (id: string) => void
  onCloseLeft?: (id: string) => void
  onCloseRight?: (id: string) => void
  onCloseOther?: (id: string) => void
}) {
  const tabRef = React.useRef<HTMLDivElement>(null)
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: id,
    disabled: isFirst, // 第一个tab不能拖拽
    animateLayoutChanges: defaultAnimateLayoutChanges,
  })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.3 : 1,
  }

  const handleClick = (e?: React.MouseEvent) => {
    // 如果点击的是删除按钮，不执行切换和滚动
    if (e && (e.target as HTMLElement).closest('.tab-close-btn')) return
    tabRef.current?.scrollIntoView({
      behavior: 'smooth',
      inline: 'center',
      block: 'nearest',
    })
    onClick?.(id)
  }

  return (
    <ContextMenu>
      <ContextMenuTrigger asChild>
        <div
          ref={(node) => {
            setNodeRef(node)
            tabRef.current = node
          }}
          style={style}
          {...attributes}
          {...(!isFirst ? listeners : {})}
        >
          <TabsTrigger
            key={id}
            value={id}
            width='fit'
            className={`gap-x-1 px-2 py-1 transition-colors ${
              isDragging ? 'cursor-grabbing' : 'cursor-pointer'
            }`}
            onClick={handleClick}
            asChild
          >
            <p>
              {icon}
              {title}
              {!isFirst && (
                <Button
                  variant='ghost'
                  size='icon'
                  className='tab-close-btn size-4 scale-150 rounded-full'
                  onClick={(e) => {
                    e.stopPropagation()
                    e.preventDefault()
                    onClose?.(id)
                  }}
                  onPointerDown={(e) => {
                    e.stopPropagation()
                  }}
                >
                  <IconX className='size-4 scale-75' />
                </Button>
              )}
            </p>
          </TabsTrigger>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent>
        <ContextMenuItem onClick={() => onRefresh?.(id)}>
          <IconRefresh className='mr-2' />
          刷新当前标签页
        </ContextMenuItem>
        <ContextMenuSeparator />
        {hasLeftTabs && (
          <ContextMenuItem onClick={() => onCloseLeft?.(id)}>
            <IconArrowBarLeft className='mr-2' />
            关闭左侧标签页
          </ContextMenuItem>
        )}
        {hasRightTabs && (
          <ContextMenuItem onClick={() => onCloseRight?.(id)}>
            <IconArrowBarRight className='mr-2' />
            关闭右侧标签页
          </ContextMenuItem>
        )}
        <ContextMenuItem onClick={() => onClose?.(id)} disabled={isFirst}>
          <IconX className='mr-2' />
          关闭当前标签页
        </ContextMenuItem>
        {totalTabs && totalTabs > 2 && (
          <ContextMenuItem onClick={() => onCloseOther?.(id)}>
            <IconMinus className='mr-2' />
            关闭其它标签页
          </ContextMenuItem>
        )}
      </ContextMenuContent>
    </ContextMenu>
  )
}
