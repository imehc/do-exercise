import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  createFileRoute,
  LinkProps,
  Outlet,
  redirect,
} from '@tanstack/react-router'
import { SearchProvider } from '~/provider/search'
import { cn } from '~/lib/utils'
import { handleMenuTree } from '~/utils/handle-menu-tree'
import { SidebarProvider } from '~/components/ui/sidebar'
import { AppSidebar } from '~/components/layout/app-sidebar'
import { NavGroup } from '~/components/layout/types'
import { LoadingSpinner } from '~/components/other'
import SkipToMain from '~/components/skip-to-main'
import { findMenuTree } from '~/features/menu/data/api'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context: { queryClient }, location }) => {
    try {
      const { menu } = handleMenuTree(
        await queryClient.ensureQueryData(findMenuTree())
      )
      if (location.pathname === '/') {
        return
      }
      const hasPermission = menu.some(
        (item) => item.route === location.pathname
      )
      if (!hasPermission) {
        return redirect({ to: '/403' })
      }
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('error', error)
    }
  },
  component: RouteComponent,
  pendingComponent: () => <LoadingSpinner isScreen />,
})

function RouteComponent() {
  const { data: menusData = [] } = useQuery(findMenuTree())
  const navGroups = useMemo(() => {
    const { directory = [] } = handleMenuTree(menusData)
    return directory.map(
      (item) =>
        ({
          title: item.name,
          items: item.children.map((child) => ({
            title: child.name,
            url: child.route as LinkProps['to'],
            icon: child.icon,
          })),
        }) as NavGroup
    )
  }, [menusData])

  return (
    <SearchProvider navGroups={navGroups}>
      <SidebarProvider>
        <SkipToMain />
        <AppSidebar navGroups={navGroups} />
        <div
          id='content'
          className={cn(
            'ml-auto w-full max-w-full',
            'peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]',
            'peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]',
            'sm:transition-[width] sm:duration-200 sm:ease-linear',
            'flex h-svh flex-col',
            'group-data-[scroll-locked=1]/body:h-full',
            'has-[main.fixed-main]:group-data-[scroll-locked=1]/body:h-svh'
          )}
        >
          <Outlet />
        </div>
      </SidebarProvider>
    </SearchProvider>
  )
}
