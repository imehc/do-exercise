import { useEffect, useState, useRef, useCallback } from 'react'

type UseCountdownOptions = {
  seconds?: number // 倒计时时长（秒），默认 60 秒
  onFinish?: () => void // 倒计时结束后的回调
}

export function useCountdown({
  seconds = 60,
  onFinish,
}: UseCountdownOptions = {}) {
  const [count, setCount] = useState(0)
  const timerRef = useRef<NodeJS.Timeout | null>(null)

  const start = useCallback(() => {
    if (count > 0) return // 正在倒计时中
    setCount(seconds)
  }, [count, seconds])

  const reset = useCallback(() => {
    setCount(0)
    if (timerRef.current) {
      clearInterval(timerRef.current)
      timerRef.current = null
    }
  }, [])

  useEffect(() => {
    if (count === 0) {
      if (timerRef.current) {
        clearInterval(timerRef.current)
        timerRef.current = null
      }
      onFinish?.()
      return
    }

    timerRef.current = setInterval(() => {
      setCount((prev) => {
        if (prev <= 1) {
          clearInterval(timerRef.current!)
          timerRef.current = null
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current)
        timerRef.current = null
      }
    }
  }, [count, onFinish])

  return {
    count, // 剩余秒数，0 表示未在倒计时
    isCounting: count > 0,
    start,
    reset,
  }
}
