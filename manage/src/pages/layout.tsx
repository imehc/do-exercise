import { Suspense } from 'react'
import { Link, Outlet } from 'react-router'
import { Menu, ChevronLeft } from 'lucide-react'
import { Loading } from '~/components'
import { useRouterMenus } from '~/provider'
import { useSidebarStore } from '~/store/sidebar'
import { Button } from '~/components/ui/button'
import { cn } from '~/lib/utils'
import '~/animations.css'

const LayoutPage: React.FC = () => {
  const menus = useRouterMenus()
  const { collapsed, toggleCollapsed } = useSidebarStore()

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
          <span className={cn('font-bold animate-slide-in', collapsed && 'hidden')}>菜单菜单菜单</span>
          <Button variant="ghost" size="icon" onClick={toggleCollapsed} className="ml-auto">
            <ChevronLeft
              className={cn('transition-transform duration-300', collapsed && 'rotate-180')}
            />
          </Button>
        </div>
        <nav className="flex-1 overflow-y-auto py-4">
          <ul className="space-y-2 px-2">
            {menus.map(menu => (
              <li key={menu.route} className="space-y-1">
                <Link
                  to={menu.route}
                  className={cn(
                    'flex items-center gap-2 px-4 py-2 rounded-md hover:bg-accent/50 transition-colors',
                    !collapsed && 'justify-start',
                    collapsed && 'justify-center'
                  )}
                >
                  <Menu className="shrink-0" />
                  {!collapsed && <span className="truncate">{menu.name}</span>}
                </Link>
                {menu.children && !collapsed && (
                  <ul className="pl-6 space-y-1">
                    {menu.children.map(child => (
                      <li key={child.route}>
                        <Link
                          to={child.route}
                          className="flex items-center gap-2 px-4 py-2 rounded-md hover:bg-accent/50 transition-colors"
                        >
                          <div className="w-1.5 h-1.5 rounded-full bg-foreground/50" />
                          <span className="truncate">{child.name}</span>
                        </Link>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
          </ul>
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
