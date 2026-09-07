import { i18n } from '@lingui/core'
import { knownMenuI18nKeys } from './menu-catalog'

/**
 * 菜单行里与展示名相关的最小字段集。
 * 侧边栏、菜单表格、角色授权树用的是三种不同的生成模型（Menu / SysMenu /
 * SysMenuTree），这里按结构取交集，避免为了共用一个函数去挑其中一个当基准。
 */
export interface MenuLabelSource {
  name: string
  i18nKey?: string
}

/**
 * Resolve a menu's stable translation key while keeping old rows readable
 * until their catalog entry is added.
 */
export function getMenuLabel(menu: MenuLabelSource): string {
  const key = menu.i18nKey?.trim()
  if (!key) return menu.name

  // 描述符先落到变量再传入：id 本来就是运行期的动态值，写成字面量参数会让
  // `lingui extract` 反复报「Missing message ID」。真正的 catalog 声明在
  // menu-catalog.ts。
  const descriptor = { id: key, message: menu.name }
  return i18n._(descriptor) || menu.name
}

/**
 * 该翻译键是否已有 catalog 声明（见 menu-catalog.ts）。
 * 自建菜单可以填任意键，界面会回落到菜单名，这里只用于提示运维补齐。
 */
export function hasMenuCatalogEntry(i18nKey?: string): boolean {
  const key = i18nKey?.trim()
  return !!key && knownMenuI18nKeys.has(key)
}
