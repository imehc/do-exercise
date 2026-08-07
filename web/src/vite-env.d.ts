/// <reference types="vite/client" />

// @types/react-syntax-highlighter 只覆盖包入口与 cjs 样式路径，
// 这里为按需加载用到的 esm 深层路径补充声明。
declare module 'react-syntax-highlighter/dist/esm/prism-light' {
  import type { ComponentType } from 'react'
  import type { SyntaxHighlighterProps } from 'react-syntax-highlighter'

  const SyntaxHighlighter: ComponentType<SyntaxHighlighterProps> & {
    registerLanguage: (name: string, language: unknown) => void
    alias: (name: string, aliases: string[]) => void
  }
  export default SyntaxHighlighter
}

declare module 'react-syntax-highlighter/dist/esm/languages/prism/json' {
  const language: unknown
  export default language
}

declare module 'react-syntax-highlighter/dist/esm/styles/prism/one-dark' {
  import type { CSSProperties } from 'react'

  const style: { [key: string]: CSSProperties }
  export default style
}

declare module 'react-syntax-highlighter/dist/esm/styles/prism/one-light' {
  import type { CSSProperties } from 'react'

  const style: { [key: string]: CSSProperties }
  export default style
}
