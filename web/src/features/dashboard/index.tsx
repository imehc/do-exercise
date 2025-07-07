import { Suspense } from 'react'
import {
  Calendar,
  CalendarSkeleton,
  MainBody,
  MainHeader,
} from '~/components/other'
import { getEvents, getUsers } from '~/components/other/calendar/requests'

export default function Dashboard() {
  return (
    <>
      <MainHeader />
      <MainBody>
        <Suspense fallback={<CalendarSkeleton />}>
          <Calendar events={getEvents()} users={getUsers()} readonly />
        </Suspense>
      </MainBody>
    </>
  )
}
