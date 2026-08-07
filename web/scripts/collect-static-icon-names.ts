import fs from 'fs'
import path from 'path'

/** 匹配 `IconCheck` 这类图标标识符 */
const ICON_IDENTIFIER = /\bIcon[A-Z][A-Za-z0-9]*/g

const SOURCE_FILE = /\.(tsx?|jsx?)$/

/**
 * 收集业务代码里按名字静态引入的图标（如 `import { IconCheck } from '@tabler/icons-react'`）。
 *
 * 用于 vite.config.ts 的图标分包：图标选择器会通过 icon.tsx 的 glob 动态加载 6000+ 图标，
 * 逐个成包会产生大量碎片请求，因此按首字母合并成 26 个 chunk。但业务代码具名引入的图标
 * 必须排除在外——它们散落在十几个首字母上，一旦并进字母包，首屏就要为几十个图标下载整包
 * （实测首屏图标 gzip 从 9KB 涨到 342KB）。
 *
 * 名单从源码扫描得出，新增图标引用无需手工维护。
 */
export function collectStaticIconNames(dir: string) {
  const names = new Set<string>()

  const walk = (current: string) => {
    for (const entry of fs.readdirSync(current, { withFileTypes: true })) {
      const target = path.join(current, entry.name)
      if (entry.isDirectory()) {
        walk(target)
      } else if (SOURCE_FILE.test(entry.name)) {
        const source = fs.readFileSync(target, 'utf8')
        for (const matched of source.matchAll(ICON_IDENTIFIER)) {
          names.add(matched[0])
        }
      }
    }
  }
  walk(dir)

  return names
}
