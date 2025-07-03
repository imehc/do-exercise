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
}

/**
 * 水印组件（使用TailwindCSS样式，带防删除和防样式篡改机制）
 * @param content 水印内容，默认当前页面URL
 * @param opacity 透明度，默认0.15
 * @param rotate 旋转角度，默认-22
 * @param fontSize 字体大小，默认16
 * @param color 颜色，默认#000
 * @param zIndex 层级，默认9999
 * @param gap 水印间距，默认200
 * @param className 额外的CSS类名
 */
export const Watermark: React.FC<WatermarkProps> = ({
  content = typeof window !== 'undefined' ? window.location.href : '',
  opacity = 0.15,
  rotate = -22,
  fontSize = 16,
  color = '#000',
  zIndex = 9999,
  gap = 200,
  className = '',
}) => {
  const ref = React.useRef<HTMLDivElement>(null)

  // 生成水印canvas base64
  const canvas = React.useMemo(() => {
    const c = document.createElement('canvas')
    c.width = gap
    c.height = gap
    const ctx = c.getContext('2d')!
    ctx.globalAlpha = opacity
    ctx.font = `${fontSize}px sans-serif`
    ctx.fillStyle = color
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.translate(gap / 2, gap / 2)
    ctx.rotate((rotate * Math.PI) / 180)
    ctx.fillText(content, 0, 0)
    return c.toDataURL('image/png')
  }, [content, opacity, rotate, fontSize, color, gap])

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
  React.useEffect(() => {
    const node = ref.current
    if (!node) return
    const expected = {
      display: 'block',
      visibility: 'visible',
      opacity: '1',
      pointerEvents: 'none',
      zIndex: String(zIndex),
      backgroundImage: `url(${canvas})`,
    }
    const checkAndFix = () => {
      let changed = false
      const style = node.style
      if (style.display !== expected.display) {
        style.display = expected.display
        changed = true
      }
      if (style.visibility !== expected.visibility) {
        style.visibility = expected.visibility
        changed = true
      }
      if (style.opacity !== expected.opacity) {
        style.opacity = expected.opacity
        changed = true
      }
      if (style.pointerEvents !== expected.pointerEvents) {
        style.pointerEvents = expected.pointerEvents
        changed = true
      }
      if (style.zIndex !== expected.zIndex) {
        style.zIndex = expected.zIndex
        changed = true
      }
      if (style.backgroundImage !== expected.backgroundImage) {
        style.backgroundImage = expected.backgroundImage
        changed = true
      }
      // 其他关键样式可按需补充
    }
    const timer = setInterval(checkAndFix, 1000)
    return () => clearInterval(timer)
  }, [canvas, zIndex])

  return (
    <div
      ref={ref}
      className={cn(`pointer-events-none fixed inset-0 w-screen h-screen select-none ${className}`)}
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