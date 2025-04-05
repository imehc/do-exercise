import { useQuery } from '@tanstack/react-query'
import { createContext, lazy, PropsWithChildren, useContext, useState } from 'react'
import { createBrowserRouter, Navigate, RouteObject, RouterProvider } from 'react-router'
import LayoutPage from '~/pages/layout'
import NotFoundPage from '~/pages/not-found'
import { getAdminMenus, MenuItem } from '#/menus'
import { AuthWrapper, Loading } from '~/components'
import { type IconName } from 'lucide-react/dynamic'

export type RouteHandle = {
  path: string
  name: string
  icon?: IconName
}

const Dashboard = lazy(() => import('~/pages/dashboard/page.tsx'))
const Login = lazy(() => import('~/pages/login/page.tsx'))

const baseRouter: RouteObject[] = [
  {
    id: 'layout',
    path: '/',
    element: (
      <AuthWrapper>
        <LayoutPage />
      </AuthWrapper>
    ),
    children: [
      {
        path: '/',
        element: <Navigate to="/dashboard" />
      },
      {
        id: 'dashboard',
        path: '/dashboard',
        Component: Dashboard,
        handle: {
          path: '/dashboard',
          name: '仪表盘',
          icon: 'layout-dashboard'
        } as RouteHandle
      }
    ]
  },
  {
    id: 'login',
    path: '/login',
    Component: Login
  },
  {
    id: 'not-found',
    path: '*',
    Component: NotFoundPage
  }
]

const modules = import.meta.glob('~/pages/**/page.tsx')
const components = Object.keys(modules).reduce<Record<string, unknown>>((prev, curr) => {
  prev[curr.replace('/src/pages', '')] = modules[curr]
  return prev
}, {})

interface RouterMenuContextProvider {
  menus: MenuItem[]
}

const RouterMenuContext = createContext<RouterMenuContextProvider>({
  menus: []
})

type Component = () => Promise<PropsWithChildren<any>>

/** 路由菜单 */
export default function RouteMenuProvider() {
  const [router, setRouter] = useState<ReturnType<typeof createBrowserRouter>>()
  const { isLoading, data: menus = [] } = useQuery({
    queryKey: ['router-menu'],
    queryFn: async () => {
      return await getAdminMenus()
    },
    onSuccess: res => {
      // 递归处理菜单项，返回扁平化的路由数组
      const flattenMenus = (menus: MenuItem[]): RouteObject[] => {
        return menus.flatMap(menu => {
          const route = {
            id: menu.id,
            path: menu.route,
            Component: menu.filePath ? lazy(components[menu.filePath] as Component) : null,
            handle: {
              path: menu.route,
              name: menu.name,
              icon: menu.icon
            } as RouteHandle
          } as RouteObject
          // 如果有子菜单，则递归处理子菜单，并添加重定向路由
          if (menu.children && menu.children.length > 0) {
            const redirectRoute = {
              id: menu.id,
              path: menu.route,
              element: <Navigate to={menu.children[0].route} />
            }
            return [redirectRoute, ...flattenMenus(menu.children)]
          }
          return [route]
        })
      }

      // 获取菜单后，动态添加路由
      const children = flattenMenus(res)
      const existingIds = new Set((baseRouter[0].children as RouteObject[]).map(r => r.id))
      const uniqueChildren = children.filter(child => !existingIds.has(child.id))
      baseRouter[0].children = [...(baseRouter[0].children as RouteObject[]), ...uniqueChildren]
      setRouter(createBrowserRouter(baseRouter))
    }
  })

  if (isLoading || !router) {
    return <Loading global />
  }

  return (
    <RouterMenuContext.Provider
      value={{
        menus: [
          {
            id: 'dashboard',
            name: '首页',
            route: '/dashboard',
            filePath: '/dashboard/page.tsx',
            children: []
          },
          ...menus
        ]
      }}
    >
      <RouterProvider router={router} />
    </RouterMenuContext.Provider>
  )
}

/** 获取路由菜单 */
export const useRouterMenus = () => {
  const { menus } = useContext(RouterMenuContext)
  return menus
}
