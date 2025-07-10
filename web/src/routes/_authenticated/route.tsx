import { useMemo } from 'react'
import { queryOptions, useQuery } from '@tanstack/react-query'
import {
  createFileRoute,
  LinkProps,
  Outlet,
  redirect,
} from '@tanstack/react-router'
import { MenuType, UserApi } from '~/do-exercise-api'
import { PermissionProvider } from '~/provider'
import { SearchProvider } from '~/provider/search'
import { cn } from '~/lib/utils'
import { handleToMenuTree } from '~/utils/handle-menu-tree'
import { apiInstance } from '~/hooks/use-api'
import { useUserProfile } from '~/hooks/use-user'
import { SidebarProvider } from '~/components/ui/sidebar'
import { AppSidebar } from '~/components/layout/app-sidebar'
import { NavGroup } from '~/components/layout/types'
import { LoadingSpinner, MainHeader, Watermark } from '~/components/other'

const getUserMenu = () => {
  const userApi = apiInstance(UserApi)
  return queryOptions({
    queryKey: ['getUserMenu'],
    queryFn: () => userApi.getUserMenu(),
    staleTime: 60 * 1000,
  })
}

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context: { queryClient }, location }) => {
    try {
      const menus = await queryClient.ensureQueryData(getUserMenu())
      if (
        location.pathname === '/' ||
        location.pathname.startsWith('/settings')
      ) {
        return
      }
      const hasPermission = menus
        .filter((item) => item.type === MenuType.menu)
        .some((item) => item.route === location.pathname)
      if (!hasPermission) {
        return redirect({ to: '/403' })
      }
    } catch (error) {
      console.error('error', error)
    }
  },
  component: RouteComponent,
  pendingComponent: () => <LoadingSpinner isScreen />,
})

function RouteComponent() {
  const { data: userProfile, isLoading: userProfileIsLoading } =
    useUserProfile()

  const { data: menus = [], isLoading: userMenuIsLoading } =
    useQuery(getUserMenu())
  const navGroups = useMemo(() => {
    const { directory = [], menu = [] } = handleToMenuTree(menus)
    if (directory.length > 0) {
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
    }
    return menu.map(
      (item) =>
        ({
          title: item.name,
          items: [],
        }) as NavGroup
    )
  }, [menus])

  const permissions = useMemo(
    () =>
      menus
        .filter((item) => item.type === MenuType.button)
        .map((item) => item.permission!),
    [menus]
  )

  const isLoading = userProfileIsLoading || userMenuIsLoading

  if (isLoading) {
    return <LoadingSpinner isScreen />
  }

  return (
    <PermissionProvider permissions={permissions}>
      <SearchProvider navGroups={navGroups}>
        <SidebarProvider>
          {/* <SkipToMain /> */}
          <AppSidebar navGroups={navGroups} user={userProfile} />
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
            <MainHeader />
            <Outlet />
          </div>
          <Watermark content={userProfile?.nickname || userProfile?.username} />
        </SidebarProvider>
      </SearchProvider>
    </PermissionProvider>
  )
}
