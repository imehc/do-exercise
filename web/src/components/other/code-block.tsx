import React, { useState } from 'react'
import { CopyIcon, CheckIcon, ChevronsUpDownIcon } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import {
  oneDark,
  oneLight,
} from 'react-syntax-highlighter/dist/cjs/styles/prism'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'

interface CodeBlockProps {
  content: string | object
  language?: 'json' | 'text'
  className?: string
  collapsedHeight?: string
}

export const CodeBlock: React.FC<CodeBlockProps> = ({
  content,
  language = 'text',
  className,
  collapsedHeight = 'max-h-64',
}) => {
  const { theme } = useTheme()
  const [copied, setCopied] = useState(false)
  const [expanded, setExpanded] = useState(false)

  const formatted =
    language === 'json'
      ? (() => {
          try {
            const parsed =
              typeof content === 'string' ? JSON.parse(content) : content
            return JSON.stringify(parsed, null, 2)
          } catch {
            return content + ''
          }
        })()
      : String(content)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(formatted)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div
      className={cn(
        'bg-muted/30 relative overflow-x-auto rounded-xl border p-4 font-mono text-sm shadow-sm',
        'min-w-0 backdrop-blur-sm',
        className
      )}
    >
      <div className='mb-2 flex items-center justify-end'>
        <div className='flex gap-2'>
          <Button
            variant='ghost'
            size='sm'
            onClick={handleCopy}
            className='bg-muted/50 hover:bg-muted h-auto px-2 py-1 backdrop-blur'
          >
            {copied ? (
              <CheckIcon className='h-4 w-4' />
            ) : (
              <CopyIcon className='h-4 w-4' />
            )}
          </Button>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => setExpanded(!expanded)}
            className='bg-muted/50 hover:bg-muted h-auto gap-1 px-2 py-1 text-xs backdrop-blur'
          >
            {expanded ? '折叠' : '展开'}
            <ChevronsUpDownIcon className='h-4 w-4' />
          </Button>
        </div>
      </div>

      <div
        className={cn(
          'overflow-x-auto rounded-md transition-all duration-300 ease-in-out',
          expanded ? 'max-h-[60vh] translate-y-0 opacity-100 shadow-lg' : `${collapsedHeight} translate-y-2 opacity-95 shadow-md`
        )}
      >
        <SyntaxHighlighter
          language={language}
          style={theme === 'dark' ? oneDark : oneLight}
          customStyle={{
            background: 'transparent',
            padding: '1rem',
            margin: 0,
            borderRadius: '0.5rem',
          }}
        >
          {formatted}
        </SyntaxHighlighter>
      </div>
    </div>
  )
}
