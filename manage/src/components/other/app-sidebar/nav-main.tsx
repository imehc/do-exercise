import { MenuItem } from '#/menus'
import { ChevronRight } from 'lucide-react'
import { memo, useEffect, useMemo, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router'
import { DynamicIcon } from 'lucide-react/dynamic'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '~/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem
} from '~/components/ui/sidebar'
import { useRouterMenus } from '~/provider'
import { cn } from '~/lib/utils'

interface RecursiveMenuProps {
  /** 当前菜单项的配置信息 */
  menu: MenuItem
  /** 当前展开的菜单路径列表 */
  expandedPaths: string[]
  /** 当前激活的路由路径 */
  currentPath: string
  /** 当前菜单项的层级深度 */
  level: number
  /** 当前菜单项是否展开 */
  isExpanded: boolean
  /** 点击菜单项时的回调函数 */
  onMenuClick: (menuRoute: string) => void
  /** 是否自动跳转到第一个子菜单项，默认为false */
  autoJumpToFirst?: boolean
}

// 查找第一个没有子菜单的节点
const findFirstLeafNode = (menu: MenuItem): string => {
  if (!menu.children || menu.children.length === 0) {
    return menu.route
  }
  return findFirstLeafNode(menu.children[0])
}

const RecursiveMenuComponent = memo(
  ({
    menu,
    expandedPaths,
    currentPath,
    onMenuClick,
    autoJumpToFirst = false
  }: RecursiveMenuProps) => {
    const navigate = useNavigate()
    const hasChildren = menu.children && menu.children.length > 0
    const isActive = currentPath === menu.route

    if (!hasChildren) {
      return (
        <SidebarMenuSubItem key={menu.name}>
          <SidebarMenuSubButton
            asChild
            className={cn('w-full', isActive && 'bg-sidebar-accent text-sidebar-accent-foreground')}
          >
            <Link to={menu.route}>
              <span>{menu.name}</span>
            </Link>
          </SidebarMenuSubButton>
        </SidebarMenuSubItem>
      )
    }

    return (
      <Collapsible
        key={menu.id}
        asChild
        open={expandedPaths.includes(menu.route)}
        className="group/collapsible"
      >
        <SidebarMenuItem>
          <CollapsibleTrigger asChild>
            <SidebarMenuButton
              tooltip={menu.name}
              className={cn(
                expandedPaths.includes(menu.route) &&
                  menu.route !== currentPath &&
                  'bg-sidebar-accent text-sidebar-accent-foreground'
              )}
              onClick={() => {
                onMenuClick(menu.route)
                if (hasChildren && autoJumpToFirst) {
                  const firstLeafRoute = findFirstLeafNode(menu)
                  navigate(firstLeafRoute)
                }
              }}
            >
              <DynamicIcon name="gamepad" />
              <span>{menu.name}</span>
              <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
            </SidebarMenuButton>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <SidebarMenuSub>
              {menu.children?.map(childMenu => (
                <RecursiveMenuComponent
                  key={childMenu.id}
                  menu={childMenu}
                  expandedPaths={expandedPaths}
                  currentPath={currentPath}
                  onMenuClick={onMenuClick}
                  isExpanded={expandedPaths.includes(menu.route)}
                  level={0}
                />
              ))}
            </SidebarMenuSub>
          </CollapsibleContent>
        </SidebarMenuItem>
      </Collapsible>
    )
  }
)

export function NavMain() {
  const menus = useRouterMenus()
  const location = useLocation()
  const [expandedMenus, setExpandedMenus] = useState<Set<string>>(new Set())

  // 更新 findParent 函数以返回完整路径
  const findParent = useMemo(
    () =>
      (menus: MenuItem[], path: string): string[] => {
        for (const menu of menus) {
          if (menu.route === path) {
            return [menu.route] // Return the current route as a single-element array
          }
          if (menu.children) {
            const foundInChildren = findParent(menu.children, path)
            if (foundInChildren.length > 0) {
              return [menu.route, ...foundInChildren] // Prepend the parent route to the child route
            }
          }
        }
        return []
      },
    []
  )

  // 更新 handleMenuClick 函数以确保点击菜单时不影响其它菜单的折叠展开状态
  const handleMenuClick = (menuRoute: string) => {
    setExpandedMenus(prev => {
      const next = new Set(prev)
      if (next.has(menuRoute)) {
        next.delete(menuRoute) // 如果当前菜单已展开，则折叠它
      } else {
        next.add(menuRoute) // 展开当前菜单
      }
      return next
    })
  }

  const expandedPaths = useMemo(
    () => Array.from(expandedMenus), // Use expandedMenus directly for tracking expanded states
    [expandedMenus]
  )
  const currentPath = location.pathname

  useEffect(() => {
    // Auto-expand menus based on the current route on component mount
    const parentPaths = findParent(menus, currentPath)
    if (parentPaths.length > 0) {
      setExpandedMenus(prev => {
        const next = new Set(prev) // Preserve existing expanded state
        parentPaths.forEach(path => next.add(path)) // Add necessary parent paths
        return next
      })
    }
  }, [currentPath, menus, findParent]) // Ensure effect runs when currentPath or menus change

  return (
    <SidebarGroup>
      <SidebarMenu>
        {menus.map(menu => (
          <RecursiveMenuComponent
            key={menu.id}
            menu={menu}
            expandedPaths={expandedPaths}
            currentPath={currentPath}
            level={0}
            isExpanded={expandedMenus.has(menu.route)}
            onMenuClick={handleMenuClick}
          />
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}
