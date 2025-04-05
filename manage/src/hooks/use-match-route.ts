import React, { useEffect, useState } from 'react'
import { useLocation, useMatches, useOutlet } from 'react-router'
import { type RouteHandle } from '~/provider'

/**
 * 路由匹配信息接口
 * @interface MatchRouteType
 */
interface MatchRouteType {
  /** 菜单名称 */
  label: string
  /** tab对应的url路径 */
  pathname: string
  /** 路由参数 */
  search?: string
  /** 要渲染的组件实例 */
  children: React.ReactNode
  /**
   * 实际路由路径
   * 与pathname的区别是，详情页pathname是 /:id，而routePath是实际的路径，如 /1
   */
  routePath: string
  /** 菜单图标 */
  icon?: string
}

/**
 * 获取当前匹配路由信息的Hook
 * @returns {MatchRouteType | undefined} 返回当前匹配的路由信息，如果未匹配则返回undefined
 */
export function useMatchRoute(): MatchRouteType | undefined {
  // 获取路由组件实例
  const children = useOutlet()
  // 获取所有路由
  const matches = useMatches()
  // 获取当前url
  const { pathname, search } = useLocation()

  const [matchRoute, setMatchRoute] = useState<MatchRouteType | undefined>()

  // 监听pathname变化，重新匹配并返回新路由信息
  useEffect(() => {
    // 获取当前匹配的路由
    const lastRoute = matches.at(-1)

    if (!lastRoute?.handle) return

    setMatchRoute({
      label: (lastRoute?.handle as RouteHandle)?.name,
      pathname,
      search,
      children,
      routePath: lastRoute?.pathname || '',
      icon: (lastRoute?.handle as RouteHandle)?.icon
    })
  }, [pathname])

  return matchRoute
}
