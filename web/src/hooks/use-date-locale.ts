import { i18n } from '@lingui/core'
import { enUS, zhCN } from 'date-fns/locale'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms/language'

export function useDateLocale() {
  const language = useAtomValue(languageAtom)

  return language === 'zh-CN' ? zhCN : enUS
}

export function useDateFormat() {
  const language = useAtomValue(languageAtom)
  const locale = useDateLocale()

  return {
    locale,
    language,
    // 使用 Lingui 的 i18n.date() 来格式化日期
    formatDate: (date: Date | number, options?: Intl.DateTimeFormatOptions) => {
      return i18n.date(new Date(date), options)
    },
    // 格式化星期几的简写
    formatWeekday: (date: Date | number) => {
      return i18n.date(new Date(date), { weekday: 'short' })
    },
    // 格式化日期数字
    formatDay: (date: Date | number) => {
      return i18n.date(new Date(date), { day: 'numeric' })
    },
    // 格式化月份和日期
    formatMonthDay: (date: Date | number) => {
      return i18n.date(new Date(date), {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
    },
    // 格式化日期时间
    formatDateTime: (date: Date | number, use24Hour: boolean = true) => {
      return i18n.date(new Date(date), {
        dateStyle: 'medium',
        timeStyle: 'short',
        hour12: !use24Hour,
      })
    },
    // 格式化时间
    formatTime: (date: Date | number, use24Hour: boolean = true) => {
      return i18n.date(new Date(date), {
        hour: 'numeric',
        minute: '2-digit',
        hour12: !use24Hour,
      })
    },
    // 格式化日期加星期几
    formatDateWithWeek: (date: Date | number) => {
      return i18n.date(new Date(date), {
        weekday: 'short',
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
    },
    // 格式化几月
    formatMonth: (date: Date | number) => {
      return i18n.date(new Date(date), {
        month: 'short',
      })
    },
    // 格式化年份和月份
    formatYearAndMonth: (date: Date | number) => {
      return i18n.date(new Date(date), {
        month: 'short',
        year: 'numeric',
      })
    },
  }
}
