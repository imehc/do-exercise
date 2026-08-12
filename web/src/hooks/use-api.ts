import { addSeconds } from 'date-fns'
import { toast } from 'sonner'
import { originTokenAtom, store } from '~/atoms'
import {
  AuthApi,
  Configuration,
  ConfigurationParameters,
  ErrorContext,
  Middleware,
  RequestContext,
  ResponseContext,
  ValidationError,
} from '~/do-exercise-api'
import { router } from '~/main'

let refreshTokenFlag = false
type RequestQueue = {
  init: RequestInit
  url: string
  resolve: (res: Response) => void
  reject: (err: unknown) => void
}
let requestQueue: RequestQueue[] = []

export function apiInstance<
  T extends new (conf?: Configuration) => InstanceType<T>,
>(Api: T, conf?: ConfigurationParameters) {
  const accessToken = store.get(originTokenAtom).accessToken
  const _conf = new Configuration({
    basePath: '/apis',
    // 注意如果是bearer token的话，这里则是accessToken
    apiKey: accessToken, //|| import.meta.env.VITE_APP_AUTH_KEY,
    headers: conf?.headers,
    middleware: [middleware],
    ...conf,
  })

  const instance: InstanceType<T> = new Api(_conf)

  return instance
}

export function useApi<T extends new (conf?: Configuration) => InstanceType<T>>(
  Api: T,
  conf?: ConfigurationParameters
) {
  return apiInstance(Api, conf)
}

const middleware = {
  pre: async ({ url, init }: RequestContext) => {
    return { url, init }
  },

  post: async (ctx: ResponseContext) => {
    const { response, init, url } = ctx
    if (response.status >= 400 && response.status < 500) {
      if (response.status === 400) {
        const data = (await response.json()) as ValidationError
        toast.error(data.message)
      } else if (response.status === 401) {
        // 暂停请求，处理 token 刷新
        return await new Promise<Response>((resolve, reject) => {
          requestQueue.push({ url, init, resolve, reject })
          refreshAccessToken()
        })
      } else if (response.status === 403) {
        // 默认口令强制轮换——未改密时非白名单接口返回 403，强制回到改密页
        if (store.get(originTokenAtom).mustChangePassword) {
          redirectToChangePassword()
        } else {
          redirectToForbidden()
        }
      } else if (response.status === 404) {
        // 目标记录不存在（如 token 指向的账号已被删除）。不触发 token 刷新——
        // 刷新能成功但重试仍会 404，只会来回循环。
        const data = (await response.json()) as ValidationError
        toast.error(data.message)
      } else if (response.status === 429) {
        const text = await response.text()
        toast.error(text || 'Too many requests, please try again later.')
      }
    } else if (response.status >= 500) {
      // 5xx（含网关 502/503/504）：后端或代理故障，必须显式提示而不是静默失败，
      // 否则用户会以为操作已生效。body 可能是 nginx 的 HTML（非 JSON），不能假设可解析。
      const text = await response.text().catch(() => '')
      const message =
        /"message"\s*:/.test(text)
          ? (text.match(/"message"\s*:\s*"([^"]*)"/)?.[1] ?? '')
          : ''
      toast.error(
        message || `服务暂时不可用（HTTP ${response.status}），请稍后重试`
      )
    }

    return response
  },

  onError: async ({ response }: ErrorContext) => {
    return response
  },
} satisfies Middleware

// refreshAccessToken 静默刷新 access token，供 401 拦截器与本地过期定时器共用。
// 单飞标志（refreshTokenFlag）保证两者并发触发时只发一次刷新请求——
// refresh token 是轮转的，并发刷新会让后到者拿到已消费的旧 refresh token。
export const refreshAccessToken = async () => {
  if (refreshTokenFlag) return
  refreshTokenFlag = true

  const refreshExpireTime = store.get(originTokenAtom).refreshExpireTime
  const refreshToken = store.get(originTokenAtom).refreshToken
  // 没有 refresh token 时直接跳转登录
  if (!refreshToken || Date.now() > refreshExpireTime) {
    redirectToLogin()
    return
  }

  try {
    const authApi = new AuthApi(new Configuration({ basePath: '/apis' }))
    // 注意入参要包一层 refreshTokenRequest：openapi.yaml 里该 requestBody 是内联匿名
    // schema，生成器会额外包装成 RefreshTokenOperationRequest。直接传 { refreshToken }
    // 会在客户端同步抛 RequiredError，刷新请求发不出去，401 直接退化成跳登录页。
    const token = await authApi.refreshToken({
      refreshTokenRequest: { refreshToken },
    })

    // 更新 token 到 store
    store.set(originTokenAtom, {
      ...token,
      expireTime: addSeconds(new Date(), token.expireTime).getTime(),
      refreshExpireTime: addSeconds(
        new Date(),
        token.refreshExpireTime
      ).getTime(),
    })

    // 重试之前失败的请求
    // 注意：必须展开 queued.init（完整 RequestInit），而不是 queued.init.headers。
    // 只展开 headers 会丢掉 method 与 body，导致刷新 token 后重试的
    // POST/PUT/DELETE 全部退化成无 body 的 GET，请求静默失效。
    const pending = requestQueue
    requestQueue = []
    for (const queued of pending) {
      const updatedInit: RequestInit = {
        ...queued.init,
        headers: {
          ...queued.init.headers,
          Authorization: token.accessToken,
        },
      }

      try {
        const retryResponse = await fetch(queued.url, updatedInit)
        queued.resolve(retryResponse)
      } catch (err) {
        queued.reject(err)
      }
    }
  } catch (error) {
    console.error('Token refresh failed:', error)
    redirectToLogin()
  } finally {
    // 兜底：刷新期间新入队的请求必须被拒绝，否则其 Promise 永远不会 settle，
    // 调用方的 await 会一直挂起。
    rejectQueued('Token refresh finished before this request was retried')
    refreshTokenFlag = false
  }
}

// rejectQueued 结算所有排队中的请求，避免 Promise 永久挂起
const rejectQueued = (reason: string) => {
  const pending = requestQueue
  requestQueue = []
  for (const queued of pending) {
    queued.reject(new Error(reason))
  }
}

const redirectToLogin = () => {
  rejectQueued('Session expired')
  refreshTokenFlag = false
  store.set(originTokenAtom, {
    accessToken: '',
    expireTime: 0,
    refreshToken: '',
    refreshExpireTime: 0,
    mustChangePassword: false,
  })
  router.navigate({ to: '/sign-in', replace: true })
}

const redirectToForbidden = () => {
  rejectQueued('Forbidden')
  refreshTokenFlag = false
  window.location.href = '/403' // 跳转到forbidden页
}

// redirectToChangePassword 未完成强制改密时，把用户送回改密页
const redirectToChangePassword = () => {
  rejectQueued('Password change required')
  refreshTokenFlag = false
  window.location.href = '/settings/password'
}
