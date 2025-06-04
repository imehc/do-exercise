import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '~/components/ui/sidebar'
import { NavGroup } from '~/components/layout/nav-group'
import {
  TeamSwitcher,
  TeamSwitcherProps,
} from '~/components/layout/team-switcher'
import { NavGroup as NavGroupType } from '~/components/layout/types'
import { NavUser, NavUserProps } from './nav-user'

export function AppSidebar({
  team,
  navGroups = [],
  user,
  ...props
}: React.ComponentProps<typeof Sidebar> & {
  navGroups?: NavGroupType[]
  team?: TeamSwitcherProps
  user?: NavUserProps
}) {
  return (
    <Sidebar collapsible='icon' variant='floating' {...props}>
      <SidebarHeader>
        <TeamSwitcher {...team} />
      </SidebarHeader>
      <SidebarContent>
        {navGroups.map((props) => (
          <NavGroup key={props.title} {...props} />
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser {...user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
