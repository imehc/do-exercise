import { useCallback, useMemo } from 'react'
import { Subject } from 'rxjs'

export function useSubject<T = unknown>(): [Subject<T>, (val: T) => void] {
  const subject = useMemo(() => new Subject<T>(), [])
  const handler = useCallback(
    (val: T) => {
      subject.next(val)
    },
    [subject]
  )
  return useMemo(() => [subject, handler], [handler, subject])
}

// RxJS Subject 类型区别：
// Subject: 只会推送订阅后 next 的值，订阅前的不会收到。适合事件流。
// BehaviorSubject: 有初始值，订阅时会立即收到当前值（最新的一个）。适合状态管理。
// ReplaySubject: 可缓存指定数量的历史值，新订阅者会收到缓存的所有历史值。适合需要回放历史事件的场景。
