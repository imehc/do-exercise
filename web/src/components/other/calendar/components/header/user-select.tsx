import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { Avatar, AvatarImage, AvatarFallback } from '~/components/ui/avatar'
import { AvatarGroup } from '~/components/ui/avatar-group'
import {
  Select,
  SelectItem,
  SelectContent,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { useCalendar } from '../../contexts/calendar-context'

export function UserSelect() {
  const { users, selectedUserId, filterEventsBySelectedUser } = useCalendar()

  return (
    <Select value={selectedUserId!} onValueChange={filterEventsBySelectedUser}>
      <SelectTrigger className='w-full'>
        <SelectValue placeholder={t`全部`} />
      </SelectTrigger>
      <SelectContent align='end'>
        <SelectItem value='all'>
          <AvatarGroup className='mx-2 flex items-center' max={3}>
            {users.map((user) => (
              <Avatar key={user.id} className='text-xxs size-6'>
                <AvatarImage
                  src={user.picturePath ?? undefined}
                  alt={user.name}
                />
                <AvatarFallback className='text-xxs'>
                  {user.name[0]}
                </AvatarFallback>
              </Avatar>
            ))}
          </AvatarGroup>
          <Trans>全部</Trans>
        </SelectItem>

        {users.map((user) => (
          <SelectItem
            key={user.id}
            value={user.id}
            className='flex-1 cursor-pointer'
          >
            <div className='flex items-center gap-2'>
              <Avatar key={user.id} className='size-6'>
                <AvatarImage
                  src={user.picturePath ?? undefined}
                  alt={user.name}
                />
                <AvatarFallback className='text-xxs'>
                  {user.name[0]}
                </AvatarFallback>
              </Avatar>

              <p className='truncate'>{user.name}</p>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
