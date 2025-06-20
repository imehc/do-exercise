import { useEffect } from 'react'
import {
  IconCheck,
  IconDeviceLaptop,
  IconMoon,
  IconSun,
} from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useTheme } from '~/provider/theme'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

export const themeList = [
  {
    value: 'light',
    icon: <IconSun />,
    label: <Trans>亮色模式</Trans>,
  },
  {
    value: 'dark',
    icon: <IconMoon />,
    label: <Trans>暗色模式</Trans>,
  },
  {
    value: 'system',
    icon: <IconDeviceLaptop />,
    label: <Trans>跟随系统</Trans>,
  },
] as const

export function ThemeSwitch() {
  const { theme, setTheme } = useTheme()

  /* Update theme-color meta tag
   * when theme is updated */
  useEffect(() => {
    const themeColor = theme === 'dark' ? '#020817' : '#fff'
    const metaThemeColor = document.querySelector("meta[name='theme-color']")
    if (metaThemeColor) metaThemeColor.setAttribute('content', themeColor)
  }, [theme])

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild>
        <Button variant='ghost' size='icon' className='scale-95 rounded-full'>
          <IconSun className='size-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90' />
          <IconMoon className='absolute size-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0' />
          <span className='sr-only'>
            <Trans>切换主题</Trans>
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end'>
        {themeList.map((item) => (
          <DropdownMenuItem
            key={item.value}
            onClick={() => setTheme(item.value)}
          >
            {item.icon}
            {item.label}
            <IconCheck
              size={14}
              className={cn('ml-auto', theme !== item.value && 'hidden')}
            />
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
