import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { TabsTrigger } from '~/components/ui/tabs'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger
} from '~/components/ui/context-menu'
import { useTranslation } from 'react-i18next'
import { type Tab } from '~/store/tabs'

interface SortableTabProps {
  tab: Tab
  isActive: boolean
  onClick: () => void
  onContextMenu: (e: React.MouseEvent) => void
  onClose: () => void
  onRefresh: () => void
  onCloseOther: () => void
  showCloseOptions: boolean
}

export function SortableTab({
  tab,
  isActive,
  onClick,
  onContextMenu,
  onClose,
  onRefresh,
  onCloseOther,
  showCloseOptions
}: SortableTabProps) {
  const { t } = useTranslation('system')
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: tab.id
  })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    cursor: 'grab'
  }

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners} className="relative group">
      <TabsTrigger
        value={tab.routePath}
        className="flex items-center justify-center px-2 truncate text-ellipsis overflow-hidden"
        onClick={onClick}
        onContextMenu={onContextMenu}
      >
        {isActive ? (
          <ContextMenu>
            <ContextMenuTrigger>
              <span className="line-clamp-1">{tab.label}</span>
            </ContextMenuTrigger>
            <ContextMenuContent>
              <ContextMenuItem onClick={onRefresh}>{t('refresh')}</ContextMenuItem>
              {showCloseOptions && (
                <>
                  <ContextMenuItem onClick={onClose}>{t('close')}</ContextMenuItem>
                  <ContextMenuItem onClick={onCloseOther}>{t('closeOther')}</ContextMenuItem>
                </>
              )}
            </ContextMenuContent>
          </ContextMenu>
        ) : (
          <span className="line-clamp-1">{tab.label}</span>
        )}
      </TabsTrigger>
    </div>
  )
}
