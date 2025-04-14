import { Link, useLocation } from 'react-router'
import { useMemo } from 'react'
import { useRouterMenus } from '~/provider'
import type { MenuItem } from '#/menus'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
} from '~/components/ui/breadcrumb'

// TODO: 面包屑
export function HeaderBreadcrumb() {
  const location = useLocation()
  // const router = useRouter()
  const { menus } = useRouterMenus()

  const findMenuByPath = (menus: MenuItem[], path: string): MenuItem | undefined => {
    for (const menu of menus) {
      if (menu.path === path) {
        return menu
      }
      if (menu.children.length > 0) {
        const found = findMenuByPath(menu.children, path)
        if (found) {
          return found
        }
      }
    }
    return undefined
  }

 /**  const paths = */ useMemo(() => {
    const pathSegments = location.pathname.split('/').filter(Boolean)
    return pathSegments.map((segment, index) => {
      const path = '/' + pathSegments.slice(0, index + 1).join('/')
      const menu = findMenuByPath(menus, path)
      return {
        name: index === pathSegments.length - 1 ? menu?.name || segment : (menu?.name ?? segment),
        path
      }
    })
  }, [location.pathname, menus])

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink asChild>
            <Link to="dashboard">首页</Link>
          </BreadcrumbLink>
        </BreadcrumbItem>
        {/* {paths.map((item, index) => (
          <BreadcrumbItem key={item.path}>
            <BreadcrumbLink asChild>
              {item.path === location.pathname ? (
                <span className="opacity-70 pointer-events-none">{item.name}</span>
              ) : (
                <Link to={item.path}>{item.name}</Link>
              )}
            </BreadcrumbLink>
            {index < paths.length - 1 && paths.length > 1 && (
              <BreadcrumbSeparator>/</BreadcrumbSeparator>
            )}
          </BreadcrumbItem>
        ))} */}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
