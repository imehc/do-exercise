import { selectAtom } from 'jotai/utils'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner'
import { originTokenAtom, routeTabCacheAtom, store } from '~/atoms'
import { router } from '~/main'

const SSE_URL = '/apis/sse'
const BASE_RETRY_DELAY = 1_000
const MAX_RETRY_DELAY = 30_000

// 服务端下发的会话吊销/强制下线事件类型，收到即清空本地会话并跳转登录页
const REVOCATION_TYPES = new Set(['session_revoked', 'force_logout', 'token_revoked'])

interface SseMessage {
  type: string
  content?: unknown
}

const emptyToken = {
  accessToken: '',
  expireTime: 0,
  refreshToken: '',
  refreshExpireTime: 0,
  mustChangePassword: false,
}

/** 只投影 accessToken，避免 token 刷新写回整对象时触发无谓的 SSE 重连 */
const accessTokenAtom = selectAtom(originTokenAtom, (token) => token.accessToken)

/** 清空会话并跳转登录页（与退出登录一致，并保留服务端下发的吊销原因） */
function handleRevoked(reason?: unknown) {
  store.set(originTokenAtom, emptyToken)
  store.set(routeTabCacheAtom, {
    activeTab: '',
    tabs: [],
  })
  toast.error(typeof reason === 'string' ? reason : '您的会话已被管理员终止')
  router.navigate({ to: '/sign-in', replace: true })
}

/**
 * SSE 长连接客户端。
 *
 * 用 fetch 而不是 EventSource：后端 token 校验依赖 Authorization 请求头，
 * EventSource 无法自定义请求头。收到吊销类事件后立即断开并重登，
 * 服务端重启/断流时按指数退避自动重连，登录态清除或 401/403 时停止。
 */
export function useSse() {
  const controllerRef = useRef<AbortController | null>(null)
  const retryTimerRef = useRef<number | null>(null)
  const attemptRef = useRef(0)
  // 当前连接所用的 accessToken，用来判断 token 变更是否真需要重连
  const connectedTokenRef = useRef<string | null>(null)

  const stop = () => {
    controllerRef.current?.abort()
    controllerRef.current = null
    connectedTokenRef.current = null
    if (retryTimerRef.current != null) {
      window.clearTimeout(retryTimerRef.current)
      retryTimerRef.current = null
    }
    attemptRef.current = 0
  }

  const connect = () => {
    const accessToken = store.get(accessTokenAtom)
    if (!accessToken) return

    const controller = new AbortController()
    controllerRef.current = controller
    connectedTokenRef.current = accessToken

    fetch(SSE_URL, {
      headers: { Authorization: accessToken },
      signal: controller.signal,
    })
      .then(async (response) => {
        // 会话已失效，停止重连
        if (response.status === 401 || response.status === 403) {
          stop()
          return
        }
        if (!response.ok || !response.body) {
          throw new Error(`SSE connect failed: ${response.status}`)
        }

        attemptRef.current = 0
        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        let buffer = ''
        let eventType = 'message'

        const dispatch = (data: string) => {
          if (eventType === 'ping') return
          try {
            const msg = JSON.parse(data) as SseMessage
            if (REVOCATION_TYPES.has(msg.type)) {
              handleRevoked(msg.content)
              stop()
            }
          } catch {
            // 非 JSON 消息忽略
          }
        }

        const parse = () => {
          let idx: number
          while ((idx = buffer.indexOf('\n\n')) !== -1) {
            const block = buffer.slice(0, idx)
            buffer = buffer.slice(idx + 2)
            let data = ''
            for (const line of block.split('\n')) {
              if (line.startsWith('event:')) {
                eventType = line.slice(6).trim()
              } else if (line.startsWith('data:')) {
                data += `${line.slice(5).trimStart()}\n`
              }
            }
            if (data) {
              dispatch(data.trimEnd())
              eventType = 'message'
            }
          }
        }

        // 读取完成后，若是外部主动断开则不重连
        const pump = async () => {
          for (;;) {
            const { done, value } = await reader.read()
            if (done) break
            buffer += decoder.decode(value, { stream: true })
            parse()
          }
          if (!controller.signal.aborted) scheduleReconnect()
        }
        await pump()
      })
      .catch((err) => {
        if ((err as Error).name === 'AbortError') return
        if (!controller.signal.aborted) scheduleReconnect()
      })
  }

  const scheduleReconnect = () => {
    if (retryTimerRef.current != null) return
    const attempt = attemptRef.current++
    const delay = Math.min(BASE_RETRY_DELAY * 2 ** attempt, MAX_RETRY_DELAY)
    retryTimerRef.current = window.setTimeout(() => {
      retryTimerRef.current = null
      connect()
    }, delay)
  }

  useEffect(() => {
    // 只订阅 accessToken 字段：originTokenAtom 是整对象，刷新 token 会连带
    // expireTime/refreshToken 一起写回，订阅整个 atom 会让每次刷新都触发一次
    // 无谓的断开重连。accessToken 没变就没有重连的必要。
    const unsub = store.sub(accessTokenAtom, () => {
      const accessToken = store.get(accessTokenAtom)
      if (!accessToken) {
        stop()
        return
      }
      // token 变了必须换新连接——旧连接的 Authorization 头仍是过期 token，
      // 服务端下次校验会断流。相同 token 则保持现有连接不动。
      if (accessToken !== connectedTokenRef.current) {
        stop()
        connect()
      }
    })

    if (store.get(accessTokenAtom)) connect()

    return () => {
      unsub()
      stop()
    }
  }, [])

  return null
}

/** 挂在应用顶层，负责建立并维持 SSE 连接 */
export function SseBridge() {
  useSse()
  return null
}
