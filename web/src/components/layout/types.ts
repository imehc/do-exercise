import { LinkProps } from '@tanstack/react-router'

interface BaseNavItem {
  title: string
  badge?: string
  icon?: string
}

type NavLink = BaseNavItem & {
  url: LinkProps['to']
  items?: never
}

type NavCollapsible = BaseNavItem & {
  items: (BaseNavItem & { url: LinkProps['to'] })[]
  url?: never
}

type NavItem = NavCollapsible | NavLink

interface NavGroup {
  title: string
  items: NavItem[]
}

export type { NavGroup, NavItem, NavCollapsible, NavLink }
