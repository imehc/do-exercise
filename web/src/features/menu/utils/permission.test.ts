import { describe, expect, it } from 'vitest'
import { buildMenuPermission, routePermissionKey } from './permission'

describe('menu permission generation', () => {
  it('normalizes nested routes to the frontend permission key', () => {
    expect(routePermissionKey('/operation/log/')).toBe('operation_log')
    expect(routePermissionKey('/')).toBe('')
  })

  it('generates an action permission from the parent route', () => {
    expect(buildMenuPermission('/role', 'update')).toBe('role:update')
    expect(buildMenuPermission(undefined, 'query')).toBe('')
  })
})
