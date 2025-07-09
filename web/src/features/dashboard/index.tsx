import { useAtomValue } from 'jotai'
import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { languageAtom } from '~/atoms'
import { Main } from '~/components/layout/main'
import { Calendar, LunarCalendar, MainHeader } from '~/components/other'

export default function Dashboard() {
  const lang = useAtomValue(languageAtom)
  const [isLoading, setIsLoading] = useState(false)
  const [currentComponent, setCurrentComponent] = useState<'lunar' | 'calendar'>(
    lang === 'zh-CN' ? 'lunar' : 'calendar'
  )

  useEffect(() => {
    const newComponent = lang === 'zh-CN' ? 'lunar' : 'calendar'
    if (newComponent !== currentComponent) {
      setIsLoading(true)
      const timer = setTimeout(() => {
        setCurrentComponent(newComponent)
        setIsLoading(false)
      }, 300)
      return () => clearTimeout(timer)
    }
  }, [lang, currentComponent])

  return (
    <>
      <MainHeader />
      <Main fixed>
        <AnimatePresence mode="wait">
          {isLoading ? (
            <motion.div
              key="loading"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.3 }}
              className="flex h-full w-full items-center justify-center"
            >
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </motion.div>
          ) : (
            <motion.div
              key={currentComponent}
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.3, ease: 'easeInOut' }}
              className="h-full w-full"
            >
              {currentComponent === 'lunar' ? (
                <LunarCalendar />
              ) : (
                <Calendar events={[]} users={[]} readonly />
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </Main>
    </>
  )
}
