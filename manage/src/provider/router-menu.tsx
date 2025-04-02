import { useQuery } from '@tanstack/react-query'
import { createContext, lazy, PropsWithChildren, useContext, useState } from 'react'
import { createBrowserRouter, Navigate, RouteObject, RouterProvider } from 'react-router'
import LayoutPage from '~/pages/layout'
import NotFoundPage from '~/pages/not-found'
import { getUserMenus, MenuItem } from '#/menus'
import { AuthWrapper } from '~/components'

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
        Component: Dashboard
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

/** 路由菜单 */
export default function RouteMenuProvider() {
  const [router, setRouter] = useState<ReturnType<typeof createBrowserRouter>>()
  const { isLoading, data: menus = [] } = useQuery({
    queryKey: ['router-menu'],
    queryFn: async () => {
      return await getUserMenus()
    },
    onSuccess: res => {
      // 获取菜单后，动态添加路由
      const children = res.map(menu => ({
        id: menu.id + '',
        path: menu.route,
        Component: lazy(components[menu.filePath] as () => Promise<PropsWithChildren<any>>)
      }))
      baseRouter[0].children = [...(baseRouter[0].children as RouteObject[]), ...children]
      setRouter(createBrowserRouter(baseRouter))
    }
  })

  if (isLoading || !router) {
    return <div>loading...</div>
  }

  return (
    <RouterMenuContext.Provider value={{ menus: menus }}>
      <RouterProvider router={router} />
    </RouterMenuContext.Provider>
  )
}

/** 获取路由菜单 */
export const useRouterMenus = () => {
  const { menus } = useContext(RouterMenuContext)
  return menus
}
