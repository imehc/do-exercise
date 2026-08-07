import { Suspense, lazy, useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { Main } from '~/components/layout/main'

// 两者按语言互斥渲染，懒加载避免日历与农历库进入首屏包
const Calendar = lazy(() =>
  import('~/components/other/calendar').then((m) => ({ default: m.Calendar }))
)
const LunarCalendar = lazy(() =>
  import('~/components/other/lunar').then((m) => ({ default: m.LunarCalendar }))
)

export default function Dashboard() {
  const lang = useAtomValue(languageAtom)
  const newComponent = lang === 'zh-CN' ? 'lunar' : 'calendar'
  const [currentComponent, setCurrentComponent] = useState(newComponent)
  const isLoading = currentComponent !== newComponent

  useEffect(() => {
    if (!isLoading) return
    const timer = setTimeout(() => setCurrentComponent(newComponent), 300)
    return () => clearTimeout(timer)
  }, [isLoading, newComponent])

  return (
    <Main fixed>
      <AnimatePresence mode='wait'>
        {isLoading ? (
          <motion.div
            key='loading'
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className='flex h-full w-full items-center justify-center'
          >
            <div className='border-primary h-8 w-8 animate-spin rounded-full border-b-2'></div>
          </motion.div>
        ) : (
          <motion.div
            key={currentComponent}
            initial={{ opacity: 0, x: 40 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -40 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className='h-full w-full'
          >
            <Suspense
              fallback={
                <div className='flex h-full w-full items-center justify-center'>
                  <div className='border-primary h-8 w-8 animate-spin rounded-full border-b-2'></div>
                </div>
              }
            >
              {currentComponent === 'lunar' ? (
                <LunarCalendar />
              ) : (
                <Calendar events={[]} users={[]} readonly />
              )}
            </Suspense>
          </motion.div>
        )}
      </AnimatePresence>
    </Main>
  )
}
