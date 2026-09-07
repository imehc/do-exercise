import type { MenuPermissionAction } from '../schemas/action-schema'

export function routePermissionKey(route?: string): string {
  return (
    route
      ?.trim()
      .replace(/^\/+|\/+$/g, '')
      .replaceAll('/', '_') ?? ''
  )
}

export function buildMenuPermission(
  route: string | undefined,
  action: MenuPermissionAction
): string {
  const key = routePermissionKey(route)
  return key ? `${key}:${action}` : ''
}
