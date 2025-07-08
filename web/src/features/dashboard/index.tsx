import { Main } from '~/components/layout/main'
import { LunarCalendar, MainHeader } from '~/components/other'

export default function Dashboard() {
  return (
    <>
      <MainHeader />
      <Main fixed>
        <LunarCalendar />
      </Main>
    </>
  )
}
