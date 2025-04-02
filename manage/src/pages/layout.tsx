import { Suspense, useEffect, useState } from 'react'
import { Link, Outlet, useLocation } from 'react-router'
import { Menu, ChevronLeft } from 'lucide-react'
import { Loading } from '~/components'
import { useRouterMenus } from '~/provider'
import { useSidebarStore } from '~/store/sidebar'
import { Button } from '~/components/ui/button'
import { cn } from '~/lib/utils'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '~/components/ui/accordion'
import '~/animations.css'
import { MenuItem } from '#/menus'

const LayoutPage: React.FC = () => {
  const menus = useRouterMenus()
  const { collapsed, toggleCollapsed } = useSidebarStore()
  const location = useLocation()
  const [expandedMenu, setExpandedMenu] = useState<string>('')

  // 查找当前路由对应的菜单项及其父级
  const findParent = (menus: MenuItem[], path: string): string[] => {
    for (const menu of menus) {
      if (menu.route === path) {
        return [menu.route]
      }
      if (menu.children) {
        // 递归查找子菜单
        const foundInChildren = findParent(menu.children, path)
        if (foundInChildren.length > 0) {
          return [menu.route, ...foundInChildren]
        }
      }
    }
    return []
  }

  // 将findActiveMenu函数声明提前到使用之前
  const findActiveMenu = (pathname: string) => {
    const parentPath = findParent(menus, pathname)
    return {
      activeParent: parentPath.length > 0 ? parentPath[0] : ''
    }
  }

  const { activeParent } = findActiveMenu(location.pathname)

  useEffect(() => {
    setExpandedMenu(activeParent)
  }, [activeParent])

  return (
    <div className="min-h-screen flex bg-background/80 backdrop-blur-sm">
      {/* 侧边栏 */}
      <aside
        className={cn(
          'h-screen flex flex-col border-r transition-all duration-300 ease-in-out',
          collapsed ? 'w-16' : 'w-64'
        )}
      >
        <div className="h-16 flex items-center justify-between px-4 border-b">
          <span className={cn('font-bold animate-slide-in', collapsed && 'hidden')}>
            菜单菜单菜单
          </span>
          <Button variant="ghost" size="icon" onClick={toggleCollapsed} className="ml-auto">
            <ChevronLeft
              className={cn('transition-transform duration-300', collapsed && 'rotate-180')}
            />
          </Button>
        </div>
        <nav className="flex-1 overflow-y-auto py-4">
          <Accordion
            type="single"
            collapsible
            className="space-y-2 px-2"
            value={expandedMenu}
            onValueChange={setExpandedMenu}
          >
            {menus.map(menu => (
              <AccordionItem key={menu.route} value={menu.route} className="border-none">
                <div className="flex items-center">
                  <Link
                    to={menu.route}
                    className={cn(
                      'flex flex-1 items-center gap-2 px-4 py-2 rounded-md hover:bg-accent/50 transition-colors',
                      !collapsed && 'justify-start',
                      collapsed && 'justify-center',
                      menu.route === location.pathname && 'bg-accent/50'
                    )}
                  >
                    <Menu className="shrink-0" />
                    {!collapsed && <span className="truncate">{menu.name}</span>}
                  </Link>
                  {menu.children && !collapsed && <AccordionTrigger className="px-2" />}
                </div>
                {menu.children && !collapsed && (
                  <AccordionContent>
                    <ul className="pl-6 space-y-1">
                      {menu.children.map(child => (
                        <li key={child.route}>
                          <Link
                            to={child.route}
                            className={cn(
                              'flex items-center gap-2 px-4 py-2 rounded-md hover:bg-accent/50 transition-colors',
                              child.route === location.pathname && 'bg-accent/50'
                            )}
                          >
                            <div className="w-1.5 h-1.5 rounded-full bg-foreground/50" />
                            <span className="truncate">{child.name}</span>
                          </Link>
                          {child.children && (
                            <Accordion type="single" collapsible className="pl-4">
                              <AccordionItem value={child.route} className="border-none">
                                <AccordionContent>
                                  <ul className="pl-6 space-y-1">
                                    {child.children.map(subChild => (
                                      <li key={subChild.route}>
                                        <Link
                                          to={subChild.route}
                                          className={cn(
                                            'flex items-center gap-2 px-4 py-2 rounded-md hover:bg-accent/50 transition-colors',
                                            subChild.route === location.pathname && 'bg-accent/50'
                                          )}
                                        >
                                          <div className="w-1.5 h-1.5 rounded-full bg-foreground/50" />
                                          <span className="truncate">{subChild.name}</span>
                                        </Link>
                                      </li>
                                    ))}
                                  </ul>
                                </AccordionContent>
                              </AccordionItem>
                            </Accordion>
                          )}
                        </li>
                      ))}
                    </ul>
                  </AccordionContent>
                )}
              </AccordionItem>
            ))}
          </Accordion>
        </nav>
      </aside>

      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="h-16 border-b bg-background/50 backdrop-blur-sm flex items-center px-6">
          <h1 className="text-lg font-semibold">Header</h1>
        </header>

        {/* 主内容区域 */}
        <main className="flex-1 overflow-auto p-6">
          <Suspense fallback={<Loading global />}>
            <Outlet />
          </Suspense>
        </main>
      </div>
    </div>
  )
}

export default LayoutPage
