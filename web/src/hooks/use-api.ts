import { toast } from 'sonner'
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
import { useAuthStore } from '~/stores'

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
  const auth = useAuthStore.getState().auth
  const _conf = new Configuration({
    basePath: '/apis',
    // 注意如果是bearer token的话，这里则是accessToken
    apiKey: auth?.accessToken || import.meta.env.VITE_APP_AUTH_KEY,
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
          handleRefreshToken()
        })
      } else if (response.status === 403) {
        const text = await response.text()
        toast.error(text || 'Forbidden')
        // redirectToForbidden()
      } else if (response.status === 429) {
        const text = await response.text()
        toast.error(text || 'Too many requests, please try again later.')
      }
    }

    return response
  },

  onError: async ({ response }: ErrorContext) => {
    return response
  },
} satisfies Middleware

const handleRefreshToken = async () => {
  if (refreshTokenFlag) return
  refreshTokenFlag = true

  const { auth, setAuth } = useAuthStore.getState()
  // 没有 refresh token 时直接跳转登录
  if (!auth?.refreshToken) {
    redirectToLogin()
    return
  }

  try {
    const authApi = new AuthApi(new Configuration({ basePath: '/api' }))
    const token = await authApi.refreshToken({
      refreshToken: auth.refreshToken,
    })

    // 更新 token 到 store
    setAuth(token)

    // 重试之前失败的请求
    for (const queued of requestQueue) {
      const updatedInit = {
        ...queued.init.headers,
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
    // eslint-disable-next-line no-console
    console.error('Token refresh failed:', error)
    redirectToLogin()
  } finally {
    requestQueue = []
    refreshTokenFlag = false
  }
}

// TODO: 使用router跳转，而不是window.location.href
const redirectToLogin = () => {
  requestQueue = [] // 清空队列
  refreshTokenFlag = false
  window.location.href = '/sign-in' // 跳转到登录页
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const redirectToForbidden = () => {
  requestQueue = [] // 清空队列
  refreshTokenFlag = false
  window.location.href = '/403' // 跳转到forbidden页
}
