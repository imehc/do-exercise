import { IconSettings } from '@tabler/icons-react'
import { useTheme } from '~/provider'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { Switch } from '~/components/ui/switch'
import { useCalendar } from '../../contexts/calendar-context'
import { TCalendarView } from '../../types'

export function Settings() {
  const {
    badgeVariant,
    setBadgeVariant,
    use24HourFormat,
    toggleTimeFormat,
    view,
    setView,
    agendaModeGroupBy,
    setAgendaModeGroupBy,
  } = useCalendar()
  const { theme, setTheme } = useTheme()

  const isDarkMode = theme === 'dark'
  const isDotVariant = badgeVariant === 'dot'

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant='outline' size='icon'>
          <IconSettings />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className='w-56'>
        <DropdownMenuLabel>Calendar settings</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem>
            Use dark mode
            <DropdownMenuShortcut>
              <Switch
                checked={isDarkMode}
                onCheckedChange={(checked) =>
                  setTheme(checked ? 'dark' : 'light')
                }
              />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem>
            Use dot badge
            <DropdownMenuShortcut>
              <Switch
                checked={isDotVariant}
                onCheckedChange={(checked) =>
                  setBadgeVariant(checked ? 'dot' : 'colored')
                }
              />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem>
            Use 24 hour format
            <DropdownMenuShortcut>
              <Switch
                checked={use24HourFormat}
                onCheckedChange={toggleTimeFormat}
              />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup className='w-56'>
          <DropdownMenuLabel>Default view</DropdownMenuLabel>
          <DropdownMenuRadioGroup
            value={view}
            onValueChange={(value) => setView(value as TCalendarView)}
          >
            <DropdownMenuRadioItem value='day'>Day</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value='week'>Week</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value='month'>Month</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value='year'>Year</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value='agenda'>Agenda</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuLabel>Agenda view group by</DropdownMenuLabel>
          <DropdownMenuRadioGroup
            value={agendaModeGroupBy}
            onValueChange={(value) =>
              setAgendaModeGroupBy(value as 'date' | 'color')
            }
          >
            <DropdownMenuRadioItem value='date'>Date</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value='color'>Color</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
