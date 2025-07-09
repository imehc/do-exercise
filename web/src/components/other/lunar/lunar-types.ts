// 公共类型定义，供 lunar 相关组件复用
import { LegalHoliday } from 'tyme4ts'

export interface SolarInfo {
  year: number
  month: number
  day: number
  week?: number
}

export interface LunarInfo {
  year: number
  month: number
  day: number
  ganzhi?: string
  solarFestival?: string
  lunarFestival?: string
  holidayFestival?: LegalHoliday | null
  jieQi?: string
  isJieQi?: boolean
}

export interface LunarCell {
  solar: SolarInfo
  lunar: LunarInfo
  isToday?: boolean
  isOtherMonth?: boolean
}

export interface God {
  getLuck?: () => { getName?: () => string }
  getName?: () => string
}

export interface Taboo {
  getName: () => string
}

export interface HourLuck {
  label: string
  recommends: Taboo[]
  avoids: Taboo[]
}
