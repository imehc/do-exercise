import { useState, useEffect } from 'react'
import {
  IconMinimize,
  IconMaximize,
  IconRefresh,
  IconDots,
  IconX,
  IconMinus,
} from '@tabler/icons-react'
import screenfull from 'screenfull'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { Separator } from '~/components/ui/separator'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '~/components/ui/tooltip'

interface RightMenuProps {
  activeTab?: string
  totalTabs?: number
  onRefresh?: (id: string) => void
  onClose?: (id: string) => void
  onCloseAll?: () => void
}

export function RightMenu({
  activeTab,
  totalTabs,
  onRefresh,
  onClose,
  onCloseAll,
}: RightMenuProps) {
  const [fullscreen, setFullscreen] = useState(false)

  // 监听全屏状态变化，确保 UI 状态同步
  useEffect(() => {
    if (!screenfull.isEnabled) return
    const handler = () => setFullscreen(screenfull.isFullscreen)
    screenfull.on('change', handler)
    return () => {
      screenfull.off('change', handler)
    }
  }, [])

  // 进入/退出全屏
  const toggleFullscreen = () => {
    if (!screenfull.isEnabled) return
    const root = document.getElementById('root') || document.documentElement
    if (!screenfull.isFullscreen) {
      screenfull.request(root)
    } else {
      screenfull.exit()
    }
  }

  return (
    <div className='bg-background/80 z-10 flex items-center gap-1 pl-2'>
      <Separator orientation='vertical' className='!h-6' />
      {/* 更多 */}
      <Tooltip>
        <TooltipTrigger asChild>
          <MoreMenu
            activeTab={activeTab}
            totalTabs={totalTabs}
            onClose={onClose}
            onCloseAll={onCloseAll}
          />
        </TooltipTrigger>
        <TooltipContent>更多操作</TooltipContent>
      </Tooltip>
      {/* 刷新 */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant='ghost'
            size='icon'
            className='size-9 scale-95'
            onClick={() => activeTab && onRefresh?.(activeTab)}
          >
            <IconRefresh className='size-[1.2rem]' />
          </Button>
        </TooltipTrigger>
        <TooltipContent>刷新当前标签页</TooltipContent>
      </Tooltip>
      {/* 全屏 */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant='ghost'
            size='icon'
            onClick={toggleFullscreen}
            className='size-9 scale-95'
          >
            {fullscreen ? (
              <IconMinimize className='size-[1.2rem]' />
            ) : (
              <IconMaximize className='size-[1.2rem]' />
            )}
          </Button>
        </TooltipTrigger>
        <TooltipContent>{fullscreen ? '退出全屏' : '全屏'}</TooltipContent>
      </Tooltip>
    </div>
  )
}

function MoreMenu({
  activeTab,
  totalTabs,
  onClose,
  onCloseAll,
}: {
  activeTab?: string
  totalTabs?: number
  onClose?: (id: string) => void
  onCloseAll?: () => void
}) {
  const isFirstTab = activeTab === 'dashboard'

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant='ghost' size='icon' className='ml-1 size-9 scale-95'>
          <IconDots className='size-[1.2rem]' />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end'>
        <DropdownMenuItem
          onClick={() => activeTab && onClose?.(activeTab)}
          disabled={isFirstTab}
        >
          <IconX className='mr-1' />
          关闭当前标签页
        </DropdownMenuItem>
        {totalTabs && totalTabs > 1 && (
          <DropdownMenuItem onClick={() => onCloseAll?.()}>
            <IconMinus className='mr-1' />
            关闭全部标签页
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
