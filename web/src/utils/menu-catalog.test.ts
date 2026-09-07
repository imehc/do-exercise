import { describe, expect, it } from 'vitest'
import { knownMenuI18nKeys, menuMessageDescriptors } from './menu-catalog'

/**
 * 这份键表必须与 server/migration/menu_seed.go 的 menuI18nKeys 完全一致：
 * 后端回填的键在前端没有声明，界面就会静默回落到数据库里的中文名，
 * 而这种回落在中文环境下看不出问题，只有切到 en-US 才会暴露。
 */
const goSeedKeys = [
  'menu.system',
  'menu.api',
  'menu.menu',
  'menu.role',
  'menu.user',
  'menu.operation-log',
  'menu.token',
  'menu.system-info',
  'menu.task',
  'menu.tenant',
  'menu.action.query',
  'menu.action.info',
  'menu.action.create',
  'menu.action.update',
  'menu.action.delete',
  'menu.action.start',
  'menu.action.stop',
  'menu.action.execute',
  'menu.action.reset',
]

describe('menu-catalog', () => {
  it('与后端种子的翻译键一一对应', () => {
    expect([...knownMenuI18nKeys].sort()).toEqual([...goSeedKeys].sort())
  })

  it('每条声明都带 id 与源文案', () => {
    for (const descriptor of menuMessageDescriptors) {
      expect(descriptor.id).toBeTruthy()
      expect(descriptor.message).toBeTruthy()
    }
  })

  it('没有重复键', () => {
    expect(knownMenuI18nKeys.size).toBe(menuMessageDescriptors.length)
  })
})
