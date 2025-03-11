'use client';

import * as React from 'react';
import {
  LucideIcon,
  SquareTerminal,
} from 'lucide-react';

import { NavMain } from '~/components/other/nav-main';
import { NavUser } from '~/components/other/nav-user';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '~/components/ui/sidebar';
import { type MenuCompact } from '~/do-exercise-api';

// This is sample data.
const data = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Acme Inc',
      logo: 'mdi:company',
      plan: 'Enterprise',
    },
  ],
  navMain: [
    {
      title: 'Playground',
      url: '#',
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: 'History',
          url: '#',
        },
        {
          title: 'Starred',
          url: '#',
        },
        {
          title: 'Settings',
          url: '#',
        },
      ],
    },
  ],
  // projects: [
  //   {
  //     name: "Design Engineering",
  //     url: "#",
  //     icon: Frame,
  //   },
  // ],
};

interface Props extends React.ComponentProps<typeof Sidebar> {
  sideMenus: MenuCompact[]
}
type NavMainItem = Parameters<typeof NavMain>[number]['items'][number]

export function AppSidebar({ sideMenus, ...props }: Props) {
  const transformMenu = (menu: MenuCompact): NavMainItem => ({
    title: menu.label,
    url: menu.route,
    icon: menu.icon as unknown as LucideIcon,
    items: menu.children?.map(transformMenu)
  })
  const items = sideMenus.map(transformMenu)

  return (
    <Sidebar collapsible="icon" variant="inset" {...props}>
      <SidebarHeader>
        <div className="flex justify-center items-center">dashboard</div>
        {/* <TeamSwitcher teams={data.teams} /> */}
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={items} />
        {/* <NavProjects projects={data.projects} /> */}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
