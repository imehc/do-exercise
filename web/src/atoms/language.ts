import { atomWithStorage } from 'jotai/utils'
import { getLang } from '~/utils/lang'

// import localforage from 'localforage'

export const languages = [
  {
    label: 'English',
    value: 'en-US',
  },
  {
    label: '中文',
    value: 'zh-CN',
  },
] as const

type Language = (typeof languages)[number]['value']

export const languageAtom = atomWithStorage<Language>(
  'lang',
  getLang() as Language
)
