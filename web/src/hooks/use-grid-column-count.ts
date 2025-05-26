import { useEffect, useState, useCallback } from 'react'

export function useGridColumnCount(
  ref: React.RefObject<HTMLElement | null>,
  itemWidthIncludingGap: number = 48 // 每个 item 宽度 + gap
) {
  const [columnCount, setColumnCount] = useState<number>(0)

  const updateColumnCount = useCallback(() => {
    if (ref.current) {
      const width = ref.current.clientWidth
      const count = Math.floor(width / itemWidthIncludingGap)
      setColumnCount(Math.max(count, 1)) // 至少显示 1 列
    }
  }, [ref, itemWidthIncludingGap])

  useEffect(() => {
    if (!ref.current) return

    const observer = new ResizeObserver(updateColumnCount)
    observer.observe(ref.current)

    return () => {
      observer.disconnect()
    }
  }, [ref, updateColumnCount])

  return columnCount
}
