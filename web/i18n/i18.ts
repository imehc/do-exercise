export const languages = [
  {
    code: 'en-US',
    lang: 'en',
    label: 'English',
  },
  {
    code: 'zh-CN',
    lang: 'zh',
    label: '简体中文',
  },
] as const

export type LocaleType = ((typeof languages)[number])['lang']

export const locales = languages.map((lang) => lang.lang);