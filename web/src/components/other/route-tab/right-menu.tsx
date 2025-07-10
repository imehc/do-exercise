import { useState } from 'react'
import {
  IconMinimize,
  IconMaximize,
  IconRefresh,
  IconDots,
  IconX,
  IconMinus,
} from '@tabler/icons-react'
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

export function RightMenu({ activeTab, totalTabs, onRefresh, onClose, onCloseAll }: RightMenuProps) {
  const [fullscreen, setFullscreen] = useState(false)

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
            onClick={() => setFullscreen((f) => !f)}
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

function MoreMenu({ activeTab, totalTabs, onClose, onCloseAll }: {
  activeTab?: string
  totalTabs?: number
  onClose?: (id: string) => void
  onCloseAll?: () => void
}) {
  // 判断是否是第一个tab（不能关闭）
  const isFirstTab = activeTab === '1' // 假设第一个tab的id是'1'
  
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
          <IconX className='mr-2' />
          关闭当前标签页
        </DropdownMenuItem>
        {totalTabs && totalTabs > 1 && (
          <DropdownMenuItem onClick={() => onCloseAll?.()}>
            <IconMinus className='mr-2' />
            关闭全部标签页
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
