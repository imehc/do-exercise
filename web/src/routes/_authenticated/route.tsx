import { useMemo } from 'react'
import { QueryClient, queryOptions, useQuery } from '@tanstack/react-query'
import {
  createFileRoute,
  LinkProps,
  Outlet,
  redirect,
} from '@tanstack/react-router'
import { useAtomValue } from 'jotai'
import { languageAtom, originTokenAtom, store } from '~/atoms'
import { MenuType, UserApi } from '~/do-exercise-api'
import { PermissionProvider } from '~/provider'
import { SearchProvider } from '~/provider/search'
import { cn } from '~/lib/utils'
import { handleToMenuTree } from '~/utils/handle-menu-tree'
import { getMenuLabel } from '~/utils/menu-label'
import { apiInstance } from '~/hooks/use-api'
import { useUserProfile } from '~/hooks/use-user'
import { SidebarProvider } from '~/components/ui/sidebar'
import { AppSidebar } from '~/components/layout/app-sidebar'
import { NavGroup } from '~/components/layout/types'
import { LoadingSpinner, MainHeader, Watermark } from '~/components/other'

const CHANGE_PASSWORD_PATH = '/settings/password' as const

const getUserMenu = () => {
  const userApi = apiInstance(UserApi)
  return queryOptions({
    queryKey: ['getUserMenu'],
    queryFn: () => userApi.getUserMenu(),
    staleTime: 60 * 1000,
  })
}

const getUserMenuData = (queryClient: QueryClient) =>
  queryClient.ensureQueryData(getUserMenu())

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context: { queryClient }, location }) => {
    // 强制改密优先于菜单鉴权：未改密时后端只放行白名单接口，
    // 此时任何业务路由都拿不到数据，必须先回改密页而不是判成「无权限」跳 403。
    // 注意读 store 而不是 router context——context.token 是建 router 时的快照，不会更新。
    if (store.get(originTokenAtom).mustChangePassword) {
      if (location.pathname === CHANGE_PASSWORD_PATH) return
      throw redirect({ to: CHANGE_PASSWORD_PATH })
    }

    if (
      location.pathname === '/' ||
      location.pathname.startsWith('/settings')
    ) {
      return
    }

    let menus: Awaited<ReturnType<typeof getUserMenuData>>
    try {
      menus = await getUserMenuData(queryClient)
    } catch (error) {
      // 菜单拿不到时不能静默放行——否则空菜单会被后续判成「无权限」跳 403，
      // 掩盖真正的失败原因（401 已由 use-api 中间件处理，这里只兜剩下的）。
      console.error('获取用户菜单失败，无法完成路由鉴权', error)
      throw error
    }

    const hasPermission = menus
      .filter((item) => item.type === MenuType.menu && item.visible)
      .some((item) => item.route === location.pathname)
    if (!hasPermission) {
      throw redirect({ to: '/403' })
    }
  },
  component: RouteComponent,
  pendingComponent: () => <LoadingSpinner isScreen />,
})

function RouteComponent() {
  // 语言切换不会改变菜单 API 数据，因此显式订阅语言状态以重新计算标签。
  const language = useAtomValue(languageAtom)
  const { data: userProfile, isLoading: userProfileIsLoading } =
    useUserProfile()

  const { data = [], isLoading: userMenuIsLoading } = useQuery(getUserMenu())
  const navGroups = useMemo(() => {
    const { directory = [], menu = [] } = handleToMenuTree(data)
    if (directory.length > 0) {
      return directory
        .filter((item) => item.visible)
        .map(
          (item) =>
            ({
              title: getMenuLabel(item),
              items:
                item.children
                  ?.filter((child) => child.visible)
                  .map((child) => ({
                    title: getMenuLabel(child),
                    url: child.route as LinkProps['to'],
                    icon: child.icon,
                  })) ?? [],
            }) as NavGroup
        )
    }
    return menu
      .filter((item) => item.visible)
      .map(
        (item) =>
          ({
            title: getMenuLabel(item),
            items: [],
          }) as NavGroup
      )
  }, [data, language])

  const permissions = useMemo(
    () =>
      data
        .filter((item) => item.type === MenuType.button)
        .map((item) => item.permission!),
    [data]
  )
  const menus = useMemo(
    () => data.filter((item) => item.type === MenuType.menu),
    [data]
  )

  const isLoading = userProfileIsLoading || userMenuIsLoading

  if (isLoading) {
    return <LoadingSpinner isScreen />
  }

  return (
    <PermissionProvider permissions={permissions} menus={menus}>
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
