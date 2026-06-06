import { IconDeviceLaptop, IconMoon, IconSun } from '@tabler/icons-react'
import { msg } from '@lingui/core/macro'

export const themeList = [
  {
    value: 'light',
    icon: <IconSun />,
    label: msg`亮色模式`,
  },
  {
    value: 'dark',
    icon: <IconMoon />,
    label: msg`暗色模式`,
  },
  {
    value: 'system',
    icon: <IconDeviceLaptop />,
    label: msg`跟随系统`,
  },
] as const
