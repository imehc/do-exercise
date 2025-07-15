import { useState, useRef } from 'react'
import { motion, AnimatePresence, Variants } from 'framer-motion'
import { SolarDay } from 'tyme4ts'
import { Badge } from '../../ui/badge'
import { Button } from '../../ui/button'
import { LunarDayDetailDrawer } from './lunar-day-detail-drawer'
import type { LunarCell } from './lunar-types'

// 动画配置，与日历组件保持一致
const transition = {
  type: 'spring' as const,
  stiffness: 200,
  damping: 20,
}

const staggerContainer: Variants = {
  animate: {
    transition: {
      staggerChildren: 0.05,
    },
  },
}

const fadeIn: Variants = {
  initial: { opacity: 0, y: 10 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -10 },
}

const badgeAnimation: Variants = {
  initial: { scale: 0, opacity: 0 },
  animate: { scale: 1, opacity: 1 },
  whileHover: { scale: 1.1 },
}

const weekDays = ['一', '二', '三', '四', '五', '六', '日']

// 中文数字转换函数
function toChinese(num: number): string {
  const chinese = [
    '',
    '一',
    '二',
    '三',
    '四',
    '五',
    '六',
    '七',
    '八',
    '九',
    '十',
  ]
  if (num <= 10) return chinese[num]
  if (num < 20) return '十' + chinese[num - 10]
  if (num === 20) return '二十'
  if (num < 30) return '二十' + chinese[num - 20]
  return '三十'
}

// 农历日期显示函数
function getLunarDayText(day: number): string {
  if (day <= 10) return '初' + toChinese(day)
  if (day < 20) return '十' + toChinese(day - 10)
  if (day === 20) return '二十'
  if (day < 30) return '廿' + toChinese(day - 20)
  return '三十'
}

function getMonthData(year: number, month: number): LunarCell[] {
  const daysInMonth = new Date(year, month, 0).getDate()
  const result: LunarCell[] = []
  for (let d = 1; d <= daysInMonth; d++) {
    const solar = SolarDay.fromYmd(year, month, d)
    const lunar = solar.getLunarDay()
    const ganzhi = lunar.getSixtyCycle?.()?.getName?.() || ''
    const holidayFestival = solar.getLegalHoliday()
    const solarFestival = solar.getFestival?.()?.getName?.() || ''
    const lunarFestival = lunar.getFestival?.()?.getName?.() || ''
    const term = solar.getTerm?.()
    const jieQi = term?.getName?.() || ''
    const isJieQi = solar.getTermDay().getDayIndex() === 0 // 节气的第0天
    const now = new Date()
    const isToday =
      year === now.getFullYear() &&
      month === now.getMonth() + 1 &&
      d === now.getDate()
    const week = new Date(year, month - 1, d).getDay()
    result.push({
      solar: { year, month, day: d, week },
      lunar: {
        year: lunar.getYear(),
        month: lunar.getMonth(),
        day: lunar.getDay(),
        ganzhi,
        solarFestival,
        lunarFestival,
        holidayFestival,
        jieQi,
        isJieQi,
      },
      isToday,
    })
  }
  return result
}

export function LunarCalendar() {
  const now = new Date()
  const [cur, setCur] = useState({
    year: now.getFullYear(),
    month: now.getMonth() + 1,
  })
  // 记录切换方向（-1: 上个月，1: 下个月，0: 首次）
  const direction = useRef(0)

  // 新增：详情抽屉状态
  const [selectedCell, setSelectedCell] = useState<LunarCell | null>(null)
  const [showDrawer, setShowDrawer] = useState(false)

  const monthData = getMonthData(cur.year, cur.month)

  // 计算本月第一天是周几（0=周日，1=周一...）
  const firstWeekDay = new Date(cur.year, cur.month - 1, 1).getDay() // 0=周日
  // 以周一为首，调整
  const firstDayIdx = (firstWeekDay + 6) % 7 // 0=周一

  // 获取上个月的数据
  const prevMonth = cur.month === 1 ? 12 : cur.month - 1
  const prevYear = cur.month === 1 ? cur.year - 1 : cur.year
  const prevMonthData = getMonthData(prevYear, prevMonth)
  const prevMonthLastDays =
    firstDayIdx > 0
      ? prevMonthData
          .slice(-firstDayIdx)
          .map((day) => ({ ...day, isOtherMonth: true }))
      : []

  // 获取下个月的数据
  const nextMonth = cur.month === 12 ? 1 : cur.month + 1
  const nextYear = cur.month === 12 ? cur.year + 1 : cur.year
  const nextMonthData = getMonthData(nextYear, nextMonth)

  // 填充前导为上月，后导为下月
  const cells = prevMonthLastDays.concat(
    monthData.map((day) => ({ ...day, isOtherMonth: false }))
  )
  while (cells.length % 7 !== 0) {
    const idx = cells.length - prevMonthLastDays.length - monthData.length
    const nextDay = nextMonthData[idx] || nextMonthData[0]
    cells.push({ ...nextDay, isOtherMonth: true })
  }
  const rows = []
  for (let i = 0; i < cells.length; i += 7) rows.push(cells.slice(i, i + 7))

  // 切换月份
  const changeMonth = (offset: number) => {
    let m = cur.month + offset
    let y = cur.year
    if (m < 1) {
      m = 12
      y--
    }
    if (m > 12) {
      m = 1
      y++
    }
    direction.current = offset
    setCur({ year: y, month: m })
  }

  function getCellBgClass(cell: (typeof cells)[number]) {
    if (
      cell.lunar.holidayFestival &&
      typeof cell.lunar.holidayFestival === 'object' &&
      typeof cell.lunar.holidayFestival.isWork === 'function'
    ) {
      if (cell.lunar.holidayFestival.isWork() === false)
        return 'bg-destructive/10'
      if (cell.lunar.holidayFestival.isWork() === true) return 'bg-primary/10'
    }
    return cell.isOtherMonth
      ? 'bg-muted/50 hover:bg-muted/80'
      : 'bg-background hover:bg-accent'
  }

  return (
    <div className='bg-background grid h-full w-full grid-rows-[auto_1fr] rounded-xl font-sans shadow-sm select-none'>
      {/* 星期栏 */}
      <motion.div
        className='bg-muted grid h-16 grid-cols-7'
        variants={staggerContainer}
        initial='initial'
        animate='animate'
      >
        {weekDays.map((w, i) => (
          <motion.div
            key={w}
            variants={fadeIn}
            transition={{ ...transition, delay: i * 0.05 }}
            className={`flex h-full items-center justify-center text-base font-bold ${i === 5 || i === 6 ? 'text-destructive' : 'text-muted-foreground'}`}
          >
            {w}
          </motion.div>
        ))}
      </motion.div>
      {/* 日历格子区域（带动画和按钮） */}
      <div className='relative h-full w-full'>
        <AnimatePresence mode='wait' initial={false}>
          <motion.div
            key={cur.year + '-' + cur.month}
            className='pointer-events-none absolute top-1/2 left-1/2 z-1 -translate-x-1/2 -translate-y-1/2 text-9xl font-bold whitespace-nowrap select-none'
            initial={{
              opacity: 0,
              x: direction.current === 0 ? 0 : direction.current > 0 ? 40 : -40,
            }}
            animate={{ opacity: 0.05, x: 0 }}
            exit={{
              opacity: 0,
              x: direction.current === 0 ? 0 : direction.current > 0 ? -40 : 40,
            }}
            transition={{ ...transition, duration: 0.3 }}
          >
            {cur.year}年{cur.month}月
          </motion.div>
        </AnimatePresence>
        {/* 左侧切换按钮 */}
        <Button
          variant='ghost'
          size='icon'
          onClick={() => changeMonth(-1)}
          className='absolute top-1/2 left-0 z-10 h-10 w-10 -translate-y-1/2'
          aria-label='上个月'
        >
          <span className='text-2xl'>&#x2039;</span>
        </Button>
        {/* 右侧切换按钮 */}
        <Button
          variant='ghost'
          size='icon'
          onClick={() => changeMonth(1)}
          className='absolute top-1/2 right-0 z-10 h-10 w-10 -translate-y-1/2'
          aria-label='下个月'
        >
          <span className='text-2xl'>&#x203A;</span>
        </Button>
        {/* 日历格子动画区域 */}
        <AnimatePresence mode='wait' initial={false}>
          <motion.div
            key={cur.year + '-' + cur.month}
            className='grid h-full min-h-0 grid-cols-7'
            variants={staggerContainer}
            initial={{
              opacity: 0,
              x: direction.current === 0 ? 0 : direction.current > 0 ? 40 : -40,
            }}
            animate={{ opacity: 1, x: 0 }}
            exit={{
              opacity: 0,
              x: direction.current === 0 ? 0 : direction.current > 0 ? -40 : 40,
            }}
            transition={{ ...transition, duration: 0.3 }}
          >
            {rows.flat().map((cell, idx) =>
              cell ? (
                cell.isToday ? (
                  <motion.div
                    key={idx}
                    variants={fadeIn}
                    transition={{ ...transition, delay: idx * 0.01 }}
                    className='bg-primary relative flex h-full min-h-[64px] w-full cursor-pointer flex-col items-center justify-center px-1 py-2 shadow-lg'
                    onClick={() => {
                      setSelectedCell(cell)
                      setShowDrawer(true)
                    }}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    {/* 右上角“今”徽标 */}
                    <motion.span
                      variants={badgeAnimation}
                      transition={{ ...transition, delay: idx * 0.01 + 0.2 }}
                      className='bg-primary-foreground text-primary absolute top-1 right-1 rounded-full px-1 text-xs leading-tight font-bold shadow md:top-2 md:right-2 md:px-2 md:text-sm'
                    >
                      今
                    </motion.span>
                    <motion.div
                      className='text-primary-foreground mb-1 text-2xl leading-none font-extrabold md:text-3xl'
                      initial={{ scale: 0.8 }}
                      animate={{ scale: 1 }}
                      transition={{ ...transition, delay: idx * 0.01 + 0.1 }}
                    >
                      {cell.solar.day}
                    </motion.div>
                    <motion.div
                      className='text-primary-foreground flex min-h-4 w-full items-center justify-center gap-x-2 text-xs'
                      initial={{ opacity: 0, y: 5 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ ...transition, delay: idx * 0.01 + 0.15 }}
                    >
                      {cell.isOtherMonth ? (
                        cell.lunar.day === 1 ? (
                          `${cell.lunar.month < 0 ? '闰' : ''}${toChinese(Math.abs(cell.lunar.month))}月`
                        ) : (
                          getLunarDayText(cell.lunar.day)
                        )
                      ) : cell.lunar.solarFestival ? (
                        cell.lunar.solarFestival
                      ) : cell.lunar.lunarFestival ? (
                        cell.lunar.lunarFestival
                      ) : (
                        <>
                          <span>
                            {cell.lunar.isJieQi
                              ? cell.lunar.jieQi
                              : cell.lunar.day === 1
                                ? `${cell.lunar.month < 0 ? '闰' : ''}${toChinese(Math.abs(cell.lunar.month))}月`
                                : getLunarDayText(cell.lunar.day)}
                          </span>
                          <span>
                            {cell.lunar.ganzhi
                              ? cell.lunar.ganzhi.slice(0, 2)
                              : ''}
                          </span>
                        </>
                      )}
                    </motion.div>
                  </motion.div>
                ) : (
                  <motion.div
                    key={idx}
                    variants={fadeIn}
                    transition={{ ...transition, delay: idx * 0.01 }}
                    className={`group relative flex h-full min-h-[64px] w-full cursor-pointer flex-col items-center justify-center px-1 py-2 transition-all ${getCellBgClass(cell)} ${
                      cell.solar.week === 6 || cell.solar.week === 0
                        ? cell.isOtherMonth
                          ? 'text-muted-foreground/60'
                          : 'text-destructive'
                        : cell.isOtherMonth
                          ? 'text-muted-foreground/60'
                          : 'text-foreground'
                    }`}
                    onClick={() => {
                      setSelectedCell(cell)
                      setShowDrawer(true)
                    }}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    {/* 班/休 badge */}
                    {cell.lunar.holidayFestival &&
                      typeof cell.lunar.holidayFestival === 'object' &&
                      typeof cell.lunar.holidayFestival.isWork === 'function' &&
                      !cell.isOtherMonth && (
                        <motion.div
                          variants={badgeAnimation}
                          transition={{
                            ...transition,
                            delay: idx * 0.01 + 0.2,
                          }}
                        >
                          <Badge
                            variant={
                              cell.lunar.holidayFestival.isWork()
                                ? 'secondary'
                                : 'destructive'
                            }
                            className={`absolute top-1 right-1 rounded-full px-1.5 py-0.5 text-xs font-bold shadow md:top-2 md:right-2 md:px-2 md:text-sm ${cell.lunar.holidayFestival.isWork() ? 'bg-primary-foreground text-black dark:text-white' : 'bg-destructive text-white'}`}
                          >
                            {cell.lunar.holidayFestival.isWork() ? '班' : '休'}
                          </Badge>
                        </motion.div>
                      )}
                    <motion.div
                      className={`mb-1 text-2xl leading-none font-extrabold md:text-3xl ${cell.isOtherMonth ? 'text-muted-foreground/60' : ''}`}
                      initial={{ scale: 0.8 }}
                      animate={{ scale: 1 }}
                      transition={{ ...transition, delay: idx * 0.01 + 0.1 }}
                    >
                      {cell.solar.day}
                    </motion.div>
                    <motion.div
                      className={`flex min-h-4 w-full items-center justify-center gap-x-2 text-xs ${
                        cell.isOtherMonth
                          ? 'text-muted-foreground/40'
                          : 'text-gray-500'
                      }`}
                      initial={{ opacity: 0, y: 5 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ ...transition, delay: idx * 0.01 + 0.15 }}
                    >
                      {cell.isOtherMonth ? (
                        <>
                          <span>
                            {cell.lunar.day === 1
                              ? `${cell.lunar.month < 0 ? '闰' : ''}${toChinese(Math.abs(cell.lunar.month))}月`
                              : getLunarDayText(cell.lunar.day)}
                          </span>
                          <span>
                            {cell.lunar.ganzhi
                              ? cell.lunar.ganzhi.slice(0, 2)
                              : ''}
                          </span>
                        </>
                      ) : cell.lunar.solarFestival ? (
                        cell.lunar.solarFestival
                      ) : cell.lunar.lunarFestival ? (
                        cell.lunar.lunarFestival
                      ) : (
                        <>
                          {cell.lunar.isJieQi
                            ? cell.lunar.jieQi
                            : cell.lunar.day === 1
                              ? `${cell.lunar.month < 0 ? '闰' : ''}${toChinese(Math.abs(cell.lunar.month))}月`
                              : getLunarDayText(cell.lunar.day)}
                          <span>
                            {cell.lunar.ganzhi
                              ? cell.lunar.ganzhi.slice(0, 2)
                              : ''}
                          </span>
                        </>
                      )}
                    </motion.div>
                  </motion.div>
                )
              ) : null
            )}
          </motion.div>
        </AnimatePresence>
      </div>
      {/* 集成 LunarDayDetailDrawer */}
      {selectedCell && (
        <LunarDayDetailDrawer
          open={showDrawer}
          onOpenChange={setShowDrawer}
          cell={selectedCell}
        />
      )}
    </div>
  )
}
