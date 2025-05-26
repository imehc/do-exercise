import { MenuType } from '~/do-exercise-api'

export const callMenuTypes = new Map<MenuType, string>([
  [
    MenuType.directory,
    'bg-yellow-100/30 text-yellow-900 dark:text-yellow-200 border-yellow-200',
  ],
  [
    MenuType.menu,
    'bg-purple-100/30 text-purple-900 dark:text-purple-200 border-purple-200',
  ],
  [
    MenuType.button,
    'bg-cyan-100/30 text-cyan-900 dark:text-cyan-200 border-cyan-200',
  ],
])

export const callMenuMapping = new Map<MenuType, string>([
  [MenuType.directory, '目录'],
  [MenuType.menu, '菜单'],
  [MenuType.button, '页面元素'],
])

export const callVisibleTypes = new Map<boolean, string>([
  [true, 'bg-green-100/30 text-green-900 dark:text-green-200 border-green-200'],
  [false, 'bg-gray-100/30 text-gray-900 dark:text-gray-200 border-gray-200'],
])
