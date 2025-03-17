'use client';

import * as React from 'react';
import { SquareTerminal } from 'lucide-react';
import { Icon } from '@iconify/react';

import { NavUser } from '~/components/other/nav-user';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '~/components/ui/sidebar';
import { type MenuCompact } from '~/do-exercise-api';
import Link from 'next/link';

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
  sideMenus: MenuCompact[];
}
export function AppSidebar({ sideMenus, ...props }: Props) {
  return (
    <Sidebar collapsible="icon" variant="floating" {...props}>
      <SidebarHeader>
        <Link href="/dashboard">
          <SidebarMenuButton
            size="lg"
            className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground pointer-events-none"
          >
            <div className="h-8 w-8 rounded-lg flex items-center justify-center shrink-0 border border-solid">
              <Icon
                icon="material-symbols:quiz"
                className="size-5 text-primary"
              />
            </div>
            <div className="grid flex-1 text-left">
              {/* TODO: 多语言 */}
              <div className="font-semibold truncate">题库系统</div>
            </div>
          </SidebarMenuButton>
        </Link>
      </SidebarHeader>
      <SidebarContent>
        {sideMenus.map((menu) => (
          <SidebarGroup key={menu.id}>
            <SidebarGroupLabel>{menu.label}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {menu.children.map((item) => (
                  <SidebarMenuItem key={item.label}>
                    <SidebarMenuButton asChild>
                      <Link href={`/dashboard/${item.route}`}>
                        <Icon
                          icon={item.icon ?? 'material-symbols:square'}
                          className="size-4"
                        />
                        <span>{item.label}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarFooter className="border-t">
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
