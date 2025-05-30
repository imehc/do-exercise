import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '~/components/ui/sidebar'
import { NavGroup } from '~/components/layout/nav-group'
import { TeamSwitcher } from '~/components/layout/team-switcher'
import { NavGroup as NavGroupType } from '~/components/layout/types'
import { sidebarData } from './data/sidebar-data'
import { NavUser } from './nav-user'

export function AppSidebar({
  navGroups,
  ...props
}: React.ComponentProps<typeof Sidebar> & { navGroups: NavGroupType[] }) {
  return (
    <Sidebar collapsible='icon' variant='floating' {...props}>
      <SidebarHeader>
        <TeamSwitcher teams={sidebarData.teams} />
      </SidebarHeader>
      <SidebarContent>
        {navGroups.map((props) => (
          <NavGroup key={props.title} {...props} />
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={sidebarData.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
