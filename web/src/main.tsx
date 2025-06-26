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
import { FontProvider } from './provider/font'
import { ThemeProvider } from './provider/theme'
// Generated Routes
import { routeTree } from './routeTree.gen'

i18n.load({
  'en-US': enMessages,
  'zh-CN': zhMessages,
})
i18n.activate(store.get(languageAtom))

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        // eslint-disable-next-line no-console
        if (import.meta.env.DEV) console.log({ failureCount, error })

        if (failureCount >= 0 && import.meta.env.DEV) return false
        if (failureCount > 1 && import.meta.env.PROD) {
          router.navigate({ to: '/sign-in', replace: true })
          return false
        }

        // return !(
        //   error instanceof ResponseError &&
        //   [401, 403].includes(error.response?.status ?? 0)
        // )
        return true
      },
      refetchOnWindowFocus: import.meta.env.PROD,
      staleTime: 10 * 1000, // 10s
    },
    mutations: {
      onError: () => { },
    },
  },
  queryCache: new QueryCache({
    // onError: async (error) => {
    // if (error instanceof ResponseError) {
    // if (error.response?.status === 401) {
    //   toast.error('Session expired!')
    //   useAuthStore.getState().clearAuth()
    //   const redirect = `${router.history.location.href}`
    //   router.navigate({ to: '/sign-in', search: { redirect } })
    // }
    // if (error.response?.status === 500) {
    //   toast.error('Internal Server Error!')
    //   router.navigate({ to: '/500' })
    // }
    // if (error.response?.status === 403) {
    //   // router.navigate("/forbidden", { replace: true });
    // }
    // }
    // },
  }),
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
                <RouterProvider router={router} />
              </FontProvider>
            </ThemeProvider>
          </QueryClientProvider>
        </Provider>
      </I18nProvider>
    </StrictMode>
  )
}
