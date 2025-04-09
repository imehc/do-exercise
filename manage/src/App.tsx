import { QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import RouterMenuProvider from './provider/router-menu'
import { ThemeProvider } from './provider'
import './i18n'
import { Suspense } from 'react'
import { Loading } from './components'

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: error => {
      console.error(`Something went wrong: ${error}`)
    }
  })
})

function App() {
  return (
    <Suspense fallback={<Loading global />}>
      <ThemeProvider>
        <QueryClientProvider client={queryClient}>
          <RouterMenuProvider />
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </ThemeProvider>
    </Suspense>
  )
}

export default App
