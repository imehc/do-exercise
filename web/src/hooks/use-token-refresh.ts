import { useEffect, useRef } from 'react'
import { originTokenAtom, store } from '~/atoms'
import { refreshAccessToken } from './use-api'

// 在 access token 过期前主动刷新，避免断网/空闲时被动等 401。
// 与 use-api.ts 的 401 拦截器共用同一个单飞刷新（refreshAccessToken），
// 不会发生并发刷新撞上 refresh token 轮转的问题。
const REFRESH_LEAD_MS = 60_000 // 提前 60s 刷新
const MIN_DELAY_MS = 5_000 // 最小延迟，避免对近过期 token 打高频请求

/**
 * 本地过期定时器（SSE 专项前端清单第 4 项）：
 * 登录拿到 expireTime 后本地排程，到期前主动 refresh_token。
 * token 每次变化（刷新/登出）都会重新排程。
 */
export function useTokenRefresh() {
  const timerRef = useRef<number | null>(null)

  useEffect(() => {
    const schedule = () => {
      if (timerRef.current != null) {
        window.clearTimeout(timerRef.current)
        timerRef.current = null
      }

      const token = store.get(originTokenAtom)
      if (!token.accessToken || token.expireTime <= 0) return

      const remaining = token.expireTime - Date.now()
      // access token 已过期时交给 401 拦截器主导（refresh 仍可用 refresh token 换发）
      if (remaining <= 0) return

      const delay = Math.max(remaining - REFRESH_LEAD_MS, MIN_DELAY_MS)
      timerRef.current = window.setTimeout(() => {
        timerRef.current = null
        void refreshAccessToken()
      }, delay)
    }

    schedule()
    const unsub = store.sub(originTokenAtom, schedule)

    return () => {
      unsub()
      if (timerRef.current != null) {
        window.clearTimeout(timerRef.current)
      }
    }
  }, [])

  return null
}

/** 挂在应用顶层，负责排程本地 token 过期刷新 */
export function TokenRefreshBridge() {
  useTokenRefresh()
  return null
}
