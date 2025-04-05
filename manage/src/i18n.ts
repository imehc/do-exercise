import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import languageDetector from 'i18next-browser-languagedetector'

i18n
  .use(languageDetector) // 检测用户语言
  .use(initReactI18next)
  .init({
    debug: import.meta.env.MODE === 'development',
    fallbackLng: 'en', // 当切换语言时，未定义的key则使用该语言的key
    interpolation: {
      escapeValue: false
    },
    resources: {
      en: {
        system: {
          loading: 'Loading...',
          refresh: 'Refresh',
          close: 'Close',
          closeOther: 'Close Other'
        },
        theme: {
          light: 'Light',
          dark: 'Dark',
          system: 'System'
        }
      },
      zh: {
        system: {
          loading: '加载中...',
          refresh: '刷新',
          close: '关闭',
          closeOther: '关闭其他'
        },
        theme: {
          light: '浅色',
          dark: '深色',
          system: '系统'
        }
      }
    }
  })

export default i18n
