import * as icons from '@tabler/icons-react'

const iconPrefix = 'Icon'

/** 从字符串转化为Ico组件 */
export function toIconComponent(icon?: string) {
  if (!icon) return null
  try {
    return icons[
      (iconPrefix + icon) as keyof typeof icons
    ] as React.ComponentType<{ size?: number; className?: string }>
  } catch (error) {
    console.error(error)
    return null
  }
}

/** 获取所有Icon */
export function getIconComponentList() {
  return Object.entries(icons)
    .filter(([name]) => name.startsWith('Icon'))
    .map(([name]) => ({
      label: name.replace(iconPrefix, ''),
      icon: icons[name as keyof typeof icons] as icons.Icon,
    }))
}
