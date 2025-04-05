import { useMemo } from 'react'
import { useNavigate } from 'react-router'

/**
 * 路由导航Hook，提供常用的路由操作方法
 * @returns {Object} 包含路由导航方法的对象
 * - back: 返回上一页
 * - forward: 前进到下一页
 * - reload: 刷新当前页面
 * - push: 导航到指定路径
 * - replace: 替换当前路径
 */
export function useRouter() {
  const navigate = useNavigate()

  const router = useMemo(
    () => ({
      /** 返回上一页 */
      back: () => navigate(-1),
      /** 前进到下一页 */
      forward: () => navigate(1),
      /** 刷新当前页面 */
      reload: () => window.location.reload(),
      /** 导航到指定路径 @param {string} href - 目标路径 */
      push: (href: string) => navigate(href),
      /** 替换当前路径 @param {string} href - 目标路径 */
      replace: (href: string) => navigate(href, { replace: true })
    }),
    [navigate]
  )

  return router
}
