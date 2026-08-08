import React, { Suspense, lazy, useState } from 'react'
import { IconCheck, IconCopy, IconSelector } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useTheme } from 'next-themes'
import { cn } from '~/lib/utils'
import { useCopyToClipboard } from '~/hooks/use-copy-to-clipboard'
import { Button } from '~/components/ui/button'

const CodeBlockHighlighter = lazy(() => import('./code-block-highlighter'))

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
  const [expanded, setExpanded] = useState(false)
  const { copied, copy } = useCopyToClipboard()

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
            onClick={() => copy(formatted)}
            className='bg-muted/50 hover:bg-muted h-auto px-2 py-1 backdrop-blur'
          >
            {copied ? (
              <IconCheck className='h-4 w-4 text-green-500' />
            ) : (
              <IconCopy className='h-4 w-4' />
            )}
          </Button>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => setExpanded(!expanded)}
            className='bg-muted/50 hover:bg-muted h-auto gap-1 px-2 py-1 text-xs backdrop-blur'
          >
            {expanded ? <Trans>折叠</Trans> : <Trans>展开</Trans>}
            <IconSelector className='h-4 w-4' />
          </Button>
        </div>
      </div>

      <div
        className={cn(
          'overflow-x-auto rounded-md transition-all duration-300 ease-in-out',
          expanded
            ? 'max-h-[60vh] translate-y-0 opacity-100 shadow-lg'
            : `${collapsedHeight} translate-y-2 opacity-95 shadow-md`
        )}
      >
        <Suspense
          fallback={
            <pre className='m-0 overflow-x-auto rounded-md p-4 font-mono text-sm'>
              {formatted}
            </pre>
          }
        >
          <CodeBlockHighlighter
            formatted={formatted}
            language={language}
            dark={theme === 'dark'}
          />
        </Suspense>
      </div>
    </div>
  )
}
