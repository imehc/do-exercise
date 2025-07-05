import React from 'react'
import { cn } from '~/lib/utils'

interface WatermarkProps {
  content?: string
  opacity?: number
  rotate?: number
  fontSize?: number
  color?: string
  zIndex?: number
  gap?: number
  className?: string
  lightColor?: string // 新增
  darkColor?: string // 新增
}

/**
 * 水印组件（使用TailwindCSS样式，带防删除和防样式篡改机制，字体跟随shadui全局设置）
 */
export const Watermark: React.FC<WatermarkProps> = ({
  content = typeof window !== 'undefined' ? window.location.href : '',
  opacity = 0.15,
  rotate = -22,
  fontSize = 16,
  color, // 优先级最低
  zIndex = 9999,
  gap = 200,
  className = '',
  lightColor = '#000', // 默认白天黑色
  darkColor = '#fff', // 默认夜间白色
}) => {
  const ref = React.useRef<HTMLDivElement>(null)

  // 检测主题
  const getTheme = React.useCallback(() => {
    if (typeof window === 'undefined') return 'light'
    return document.documentElement.classList.contains('dark')
      ? 'dark'
      : 'light'
  }, [])

  // 获取shadui全局字体
  const getFontFamily = React.useCallback(() => {
    if (typeof window === 'undefined') return 'Inter, sans-serif'
    // 优先取html/body的font-family
    const htmlFont = window.getComputedStyle(
      document.documentElement
    ).fontFamily
    const bodyFont = window.getComputedStyle(document.body).fontFamily
    return bodyFont || htmlFont || 'Inter, sans-serif'
  }, [])

  // 监听主题切换，强制刷新 canvas
  const [theme, setTheme] = React.useState(getTheme())
  React.useEffect(() => {
    if (typeof window === 'undefined') return
    const observer = new MutationObserver(() => {
      setTheme(getTheme())
    })
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    })
    return () => observer.disconnect()
  }, [getTheme])

  // 生成水印canvas base64
  const canvas = React.useMemo(() => {
    const c = document.createElement('canvas')
    c.width = gap
    c.height = gap
    const ctx = c.getContext('2d')!
    ctx.globalAlpha = opacity
    ctx.font = `${fontSize}px ${getFontFamily()}`
    // 主题优先级：props.color > darkColor/lightColor
    let themeColor = color
    if (!themeColor) {
      themeColor = theme === 'dark' ? darkColor : lightColor
    }
    ctx.fillStyle = themeColor
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.translate(gap / 2, gap / 2)
    ctx.rotate((rotate * Math.PI) / 180)
    ctx.fillText(content, 0, 0)
    return c.toDataURL('image/png')
  }, [
    content,
    opacity,
    rotate,
    fontSize,
    color,
    gap,
    getFontFamily,
    theme,
    lightColor,
    darkColor,
  ])

  // 防删除机制：使用MutationObserver
  React.useEffect(() => {
    const node = ref.current
    if (!node) return
    const observer = new MutationObserver(() => {
      if (!document.body.contains(node)) {
        document.body.appendChild(node)
      }
    })
    observer.observe(document.body, { childList: true })
    return () => observer.disconnect()
  }, [])

  // 防样式篡改机制：定时检测并修复关键样式
  // React.useEffect(() => {
  //   const node = ref.current
  //   if (!node) return
  //   const expected = {
  //     display: 'block',
  //     visibility: 'visible',
  //     opacity: '1',
  //     pointerEvents: 'none',
  //     zIndex: String(zIndex),
  //     backgroundImage: `url(${canvas})`,
  //   }
  //   const checkAndFix = () => {
  //     let changed = false
  //     const style = node.style
  //     if (style.display !== expected.display) {
  //       style.display = expected.display
  //       changed = true
  //     }
  //     if (style.visibility !== expected.visibility) {
  //       style.visibility = expected.visibility
  //       changed = true
  //     }
  //     if (style.opacity !== expected.opacity) {
  //       style.opacity = expected.opacity
  //       changed = true
  //     }
  //     if (style.pointerEvents !== expected.pointerEvents) {
  //       style.pointerEvents = expected.pointerEvents
  //       changed = true
  //     }
  //     if (style.zIndex !== expected.zIndex) {
  //       style.zIndex = expected.zIndex
  //       changed = true
  //     }
  //     if (style.backgroundImage !== expected.backgroundImage) {
  //       style.backgroundImage = expected.backgroundImage
  //       changed = true
  //     }
  //     // 其他关键样式可按需补充
  //   }
  //   const timer = setInterval(checkAndFix, 1000)
  //   return () => clearInterval(timer)
  // }, [canvas, zIndex])

  return (
    <div
      ref={ref}
      className={cn(
        `pointer-events-none fixed inset-0 h-screen w-screen select-none ${className}`
      )}
      style={{
        zIndex,
        backgroundImage: `url(${canvas})`,
        backgroundRepeat: 'repeat',
        backgroundPosition: '0 0',
        display: 'block',
        visibility: 'visible',
        opacity: 1,
        pointerEvents: 'none',
      }}
    />
  )
}
