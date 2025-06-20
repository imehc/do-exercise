import { IconDeviceLaptop, IconMoon, IconSun } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'

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
