import React from 'react'
import json from 'react-syntax-highlighter/dist/esm/languages/prism/json'
import SyntaxHighlighter from 'react-syntax-highlighter/dist/esm/prism-light'
import oneDark from 'react-syntax-highlighter/dist/esm/styles/prism/one-dark'
import oneLight from 'react-syntax-highlighter/dist/esm/styles/prism/one-light'

// 只注册用到的语言，避免打包 refractor 全部语言定义
SyntaxHighlighter.registerLanguage('json', json)

interface CodeBlockHighlighterProps {
  formatted: string
  language: 'json' | 'text'
  dark: boolean
}

const CodeBlockHighlighter: React.FC<CodeBlockHighlighterProps> = ({
  formatted,
  language,
  dark,
}) => (
  <SyntaxHighlighter
    language={language}
    style={dark ? oneDark : oneLight}
    customStyle={{
      background: 'transparent',
      padding: '1rem',
      margin: 0,
      borderRadius: '0.5rem',
    }}
  >
    {formatted}
  </SyntaxHighlighter>
)

export default CodeBlockHighlighter
