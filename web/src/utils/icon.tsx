import { Suspense, lazy, type ComponentType } from 'react'

export type IconProps = { size?: number; className?: string }

export type IconComponent = ComponentType<IconProps>

const iconPrefix = 'Icon'

/**
 * 按需加载单个图标，避免全量图标库进入主包。
 * 每个图标是独立 chunk，只有用到的才会被下载。
 * import.meta.glob 不解析 alias，需要 node_modules 相对路径（走 bare specifier 解析）。
 */
const iconModules = import.meta.glob<{ default: IconComponent }>(
  '../../node_modules/@tabler/icons-react/dist/esm/icons/Icon*.mjs'
)

/** 图标名(不含 Icon 前缀) -> glob 模块路径 */
const iconPaths = new Map<string, string>(
  Object.keys(iconModules).map((filePath) => {
    const fileName = filePath.slice(
      filePath.lastIndexOf('/') + 1,
      -'.mjs'.length
    )
    return [fileName.slice(iconPrefix.length), filePath]
  })
)

/** 缓存包装后的组件，保证同名图标组件标识稳定，避免重复挂载 */
const componentCache = new Map<string, IconComponent>()

/** 占位，保持加载期间尺寸稳定，避免布局跳动 */
function IconPlaceholder({ size = 24, className }: IconProps) {
  return (
    <span
      className={className}
      style={{ width: size, height: size, display: 'inline-block' }}
      aria-hidden
    />
  )
}

/** 从字符串转化为Ico组件 */
export function toIconComponent(icon?: string): IconComponent | null {
  if (!icon) return null

  const cached = componentCache.get(icon)
  if (cached) return cached

  const filePath = iconPaths.get(icon)
  if (!filePath) return null

  const LazyIcon = lazy(iconModules[filePath])
  const Icon: IconComponent = ({ size, className }) => (
    <Suspense fallback={<IconPlaceholder size={size} className={className} />}>
      <LazyIcon size={size} className={className} />
    </Suspense>
  )
  Icon.displayName = iconPrefix + icon

  componentCache.set(icon, Icon)
  return Icon
}

/** 渲染图标，避免在组件渲染期间创建组件 */
export function renderIcon(icon?: string, className?: string, size?: number) {
  const Icon = toIconComponent(icon)
  if (!Icon) return null
  return <Icon className={className} size={size} />
}

let iconList: { label: string; icon: IconComponent }[] | undefined

/** 获取所有Icon */
export function getIconComponentList() {
  iconList ??= [...iconPaths.keys()].map((label) => ({
    label,
    icon: toIconComponent(label)!,
  }))
  return iconList
}
