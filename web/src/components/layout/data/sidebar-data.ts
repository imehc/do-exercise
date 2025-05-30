import {
  IconCommand,
  IconSpacingVertical,
  IconWaveSine,
} from '@tabler/icons-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'satnaing',
    email: 'satnaingdev@gmail.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Shadcn Admin',
      logo: IconCommand,
      plan: 'Vite + ShadcnUI',
    },
    {
      name: 'Acme Inc',
      logo: IconSpacingVertical,
      plan: 'Enterprise',
    },
    {
      name: 'Acme Corp.',
      logo: IconWaveSine,
      plan: 'Startup',
    },
  ],
  navGroups: [],
}
