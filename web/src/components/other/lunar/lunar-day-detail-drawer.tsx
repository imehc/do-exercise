import { useMemo } from 'react'
import { IconCalendar } from '@tabler/icons-react'
import { SixtyCycle, SolarDay, SolarTime, SolarYear } from 'tyme4ts'
import { Badge } from '../../ui/badge'
import { Button } from '../../ui/button'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
  DrawerFooter,
  DrawerClose,
} from '../../ui/drawer'
import type { LunarCell, God, Taboo } from './lunar-types'

// 获取日干
const hourNames = [
  '子',
  '丑',
  '寅',
  '卯',
  '辰',
  '巳',
  '午',
  '未',
  '申',
  '酉',
  '戌',
  '亥',
]
const hourStart = [0, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21]

interface LunarDayDetailDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  cell: LunarCell // 传入日历格子的 cell 数据
}

export function LunarDayDetailDrawer({
  open,
  onOpenChange,
  cell,
}: LunarDayDetailDrawerProps) {
  // 解析农历与公历信息
  const { solar, lunar } = cell || {}
  const solarObj = useMemo(
    () => SolarDay.fromYmd(solar.year, solar.month, solar.day),
    [solar]
  )
  const lunarObj = useMemo(() => solarObj?.getLunarDay?.(), [solarObj])
  // 干支
  const sixty = solarObj.getSixtyCycleDay()
  // 干支年
  const y = sixty.getYear().getName()
  // 干支月
  const m = sixty.getMonth().getName()
  // 干支日
  const d = sixty.getSixtyCycle().getName()
  const sy = SolarYear.fromYear(solar.year)
  // 生肖
  const rabByungYear = sy.getRabByungYear().getZodiac().getName()
  // 逐日胎神
  const fetus = lunarObj.getFetusDay().getName()
  // 建除十二神
  const duty = lunarObj.getDuty().getName()
  // 二十八星宿
  const star = lunarObj.getTwentyEightStar()
  const name = star.getName()
  const sevenStar = star.getSevenStar().getName()
  const animal = star.getAnimal().getName()
  const luck = star.getLuck().getName()
  const starInfo = `${name}${sevenStar}${animal}${luck}`
  // 彭祖百忌
  const sixtyCycle = SixtyCycle.fromIndex(lunar.day)
  const pengZu = sixtyCycle.getPengZu().getName()
  // 纳音
  const sound = sixtyCycle.getSound().getName()
  // 冲煞
  const dayBranch = sixty.getSixtyCycle().getEarthBranch()
  const clashBranch = dayBranch.getOpposite()
  const clashZodiac = clashBranch.getZodiac().getName()
  const clashDirection = clashBranch.getDirection().getName()
  const clashInfo = `冲${clashZodiac} 煞${clashDirection}`
  // 值神
  const twelveStar = lunarObj.getTwelveStar()
  const ecliptic = twelveStar.getEcliptic()
  const dutyGodLuck = ecliptic.getLuck()
  const dutyGodInfo = `${twelveStar.getName()}(${dutyGodLuck.getName()})`

  // 宜、忌
  const recommends: Taboo[] = lunarObj?.getRecommends?.() || []
  const avoids: Taboo[] = lunarObj?.getAvoids?.() || []
  // 吉神宜趋/凶神宜忌
  const gods: God[] = lunarObj?.getGods?.() || []
  const goodGods: God[] = gods.filter((g) => g.getLuck?.().getName?.() === '吉')
  const badGods: God[] = gods.filter((g) => g.getLuck?.().getName?.() !== '吉')

  const hoursInfo = useMemo(() => {
    if (!solar) return []
    return hourStart.map((h, idx) => {
      try {
        const hourObj = SolarTime?.fromYmdHms?.(
          solar.year,
          solar.month,
          solar.day,
          h,
          0,
          0
        )?.getLunarHour?.()

        if (!hourObj) {
          return {
            label: '',
            branch: hourNames[idx],
            luck: '吉',
            star: '',
            idx,
          }
        }

        // 天干地支
        const ganzhi = hourObj.getSixtyCycle?.().getName?.() || ''

        // 使用tyme4ts的时辰吉凶计算
        let luck = '吉'
        let star = ''
        try {
          const twelveStar = hourObj.getTwelveStar()
          star = twelveStar.getName()
          const ecliptic = twelveStar.getEcliptic()
          luck = ecliptic.getLuck().getName()
        } catch (e) {
          // 如果获取失败，使用默认值
          luck = '吉'
          console.error(e)
        }

        return {
          label: ganzhi,
          branch: hourNames[idx],
          luck,
          star,
          idx,
        }
      } catch (e) {
        console.error(e)
        return { label: '', branch: hourNames[idx], luck: '吉', star: '', idx }
      }
    })
  }, [solar])

  if (!cell) return null

  return (
    <Drawer
      shouldScaleBackground
      open={open}
      onOpenChange={onOpenChange}
      direction='bottom'
    >
      <DrawerContent className='max-h-[85vh]'>
        <DrawerHeader className='text-left'>
          <DrawerTitle className='flex items-center gap-2'>
            <IconCalendar />{' '}
            {`${solar?.year}年${solar?.month}月${solar?.day}日 星期${['日', '一', '二', '三', '四', '五', '六'][solar?.week || 0]}`}
          </DrawerTitle>
          <DrawerDescription>
            <span className='text-primary text-base font-semibold'>{`农历${lunar?.month}月${lunar?.day}日`}</span>
            <span className='text-muted-foreground ml-2'>{`${y}(${rabByungYear})年 ${m}月 ${d}日`}</span>
          </DrawerDescription>
        </DrawerHeader>
        <div className='grid gap-2 overflow-y-auto px-4'>
          {/* 宜忌区块 */}
          <div>
            {!!recommends.length && (
              <div className='mb-2 flex items-center gap-4'>
                <span className='text-base font-medium'>宜</span>
                <div className='flex flex-wrap gap-1'>
                  {recommends.length ? (
                    recommends.map((r) => (
                      <Badge key={r.getName()}>{r.getName()}</Badge>
                    ))
                  ) : (
                    <span className='text-muted-foreground'>—</span>
                  )}
                </div>
              </div>
            )}
            {!!avoids.length && (
              <div className='flex items-center gap-4'>
                <span className='text-base font-medium'>忌</span>
                <div className='flex flex-wrap gap-1'>
                  {avoids.length ? (
                    avoids.map((a) => (
                      <Badge key={a.getName()} variant='destructive'>
                        {a.getName()}
                      </Badge>
                    ))
                  ) : (
                    <span className='text-muted-foreground'>—</span>
                  )}
                </div>
              </div>
            )}
          </div>
          <div className='grid min-h-120 grid-cols-8 grid-rows-13 gap-0 overflow-hidden rounded-lg border border-[#e5e5e5] bg-white'>
            {/* 纳音 */}
            <div className='col-span-3 row-span-2 border-r border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-lg font-bold text-[#b98c4b]'>纳音</div>
                <div className='mt-1 text-base font-semibold text-gray-700'>
                  {sound}
                </div>
              </div>
            </div>
            {/* 冲煞 */}
            <div className='col-span-2 row-span-2 border-r border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-lg font-bold text-[#b98c4b]'>冲煞</div>
                <div className='mt-1 text-base font-semibold text-gray-700'>
                  {clashInfo}
                </div>
              </div>
            </div>
            {/* 值神（最右，无右边框） */}
            <div className='col-span-3 row-span-2 border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-lg font-bold text-[#b98c4b]'>值神</div>
                <div className='mt-1 text-base font-semibold text-gray-700'>
                  {dutyGodInfo}
                </div>
              </div>
            </div>
            {/* 时辰吉凶 */}
            <div className='col-span-1 row-span-3 border-r border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-center text-base leading-tight font-bold text-[#b98c4b]'>
                  <span>时辰</span>
                  <br />
                  <span> 吉凶</span>
                </div>
              </div>
            </div>
            {/* 时辰吉凶内容（最右，无右边框） */}
            <div className='col-span-7 row-span-3 border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full w-full items-center justify-around'>
                {hoursInfo.map((h) => (
                  <div
                    key={h.idx}
                    className='flex h-full w-full flex-col items-center justify-center'
                  >
                    <div className='vertical-rl text-sm tracking-[0.5em] text-gray-500'>
                      {h.label}
                      {h.luck}
                    </div>
                    {/* <div className="text-gray-500 text-lg">{h.branch}</div> */}
                  </div>
                ))}
              </div>
            </div>
            {/* 建除十二神（竖排，最左，无左边框） */}
            <div className='col-span-1 row-span-8 border-r border-[#e5e5e5] bg-white'>
              <div className='flex h-full items-center justify-center py-2'>
                <div className='vertical-rl text-base font-bold tracking-[0.5em] text-[#b98c4b]'>
                  建除十二神
                </div>
                <div className='vertical-rl text-base font-semibold tracking-[0.5em] text-gray-700'>
                  {duty}
                </div>
              </div>
            </div>
            {/* 吉神宜趋 */}
            <div className='col-span-2 row-span-5 border-r border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-base font-bold text-[#b98c4b]'>
                  吉神宜趋
                </div>
                <div className='mt-1'>
                  {' '}
                  {goodGods.length ? (
                    goodGods.map((g) => (
                      <Badge
                        key={g.getName?.()}
                        variant='secondary'
                        className='mx-0.5'
                      >
                        {g.getName?.()}
                      </Badge>
                    ))
                  ) : (
                    <span className='text-gray-400'>—</span>
                  )}
                </div>
              </div>
            </div>
            {/* 今日胎神 */}
            <div className='col-span-2 row-span-5 border-r border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-base font-bold text-[#b98c4b]'>
                  今日胎神
                </div>
                <div className='mt-1 text-base font-semibold text-gray-700'>
                  {fetus}
                </div>
              </div>
            </div>
            {/* 凶神宜忌（最右，无右边框） */}
            <div className='col-span-2 row-span-5 border-b border-[#e5e5e5] bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-base font-bold text-[#b98c4b]'>
                  凶神宜忌
                </div>
                <div className='mt-1'>
                  {badGods.length ? (
                    badGods.map((g) => (
                      <Badge
                        key={g.getName?.()}
                        variant='destructive'
                        className='mx-0.5'
                      >
                        {g.getName?.()}
                      </Badge>
                    ))
                  ) : (
                    <span className='text-gray-400'>—</span>
                  )}
                </div>
              </div>
            </div>
            {/* 二十八星宿（竖排，最右，无右边框） */}
            <div className='col-span-1 row-span-8 border-l border-[#e5e5e5] bg-white'>
              <div className='flex h-full items-center justify-center py-2'>
                <div className='vertical-rl text-base font-bold tracking-[0.5em] text-[#b98c4b]'>
                  二十八星宿
                </div>
                <div className='vertical-rl text-base font-semibold tracking-[0.5em] text-gray-700'>
                  {starInfo}
                </div>
              </div>
            </div>
            {/* 彭祖百忌（最下，无下边框） */}
            <div className='col-span-6 row-span-3 bg-white'>
              <div className='flex h-full flex-col items-center justify-center'>
                <div className='text-base font-bold text-[#b98c4b]'>
                  彭祖百忌
                </div>
                <div className='mt-1 text-base font-semibold text-gray-700'>
                  {pengZu.split(' ').map((item) => (
                    <div key={item}>{item}</div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
        <DrawerFooter className='gap-y-2'>
          <DrawerClose asChild>
            <Button variant='outline'>关闭</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
