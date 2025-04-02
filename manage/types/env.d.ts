interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_ENV: 'development' | 'production'
  readonly VITE_APP_HOMEPAGE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
