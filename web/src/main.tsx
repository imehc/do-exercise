import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { i18n } from '@lingui/core'
import { I18nProvider } from '@lingui/react'
import { Provider } from 'jotai'
import { messages as enMessages } from '~/locales/en-US/messages'
import { messages as zhMessages } from '~/locales/zh-CN/messages'
import { languageAtom, originTokenAtom, store } from './atoms'
import './index.css'
import { SseBridge } from './hooks/use-sse'
import { FontProvider } from './provider/font'
import { ThemeProvider } from './provider/theme'
// Generated Routes
import { routeTree } from './routeTree.gen'

// 仅在开发环境加载，动态导入避免 react-scan 进入生产包
if (import.meta.env.DEV && typeof window !== 'undefined') {
  void import('react-scan').then(({ scan }) => {
    scan({ enabled: true, log: true })
  })
}

i18n.load({
  'en-US': enMessages,
  'zh-CN': zhMessages,
})
i18n.activate(store.get(languageAtom))

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (import.meta.env.DEV) console.log({ failureCount, error })

        if (failureCount >= 0 && import.meta.env.DEV) return false
        if (failureCount > 1 && import.meta.env.PROD) {
          router.navigate({ to: '/sign-in', replace: true })
          return false
        }

        return false
      },
      refetchOnWindowFocus: import.meta.env.PROD,
      staleTime: 10 * 1000, // 10s
    },
    mutations: {
      onError: () => {},
    },
  },
  // 401/403/429 等状态码统一在 ~/hooks/use-api.ts 的 fetch middleware 中处理，
  // 那里能拿到原始请求以便刷新 token 后重试，这里不再重复拦截。
  queryCache: new QueryCache({}),
})

// Create a new router instance
export const router = createRouter({
  routeTree,
  context: { queryClient, token: store.get(originTokenAtom) },
  defaultPreload: 'intent',
  defaultPreloadStaleTime: 0,
})

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

// Render the app
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <I18nProvider i18n={i18n}>
        <Provider store={store}>
          <QueryClientProvider client={queryClient}>
            <ThemeProvider defaultTheme='light' storageKey='vite-ui-theme'>
              <FontProvider>
                <SseBridge />
                <RouterProvider router={router} />
              </FontProvider>
            </ThemeProvider>
          </QueryClientProvider>
        </Provider>
      </I18nProvider>
    </StrictMode>
  )
}
