import React, { createContext, useContext } from 'react'
import { useMatch } from '@tanstack/react-router'

type PermissionContextType = {
  permissions: string[]
}

export const PermissionContext = createContext<PermissionContextType>({
  permissions: [],
})

export const PermissionProvider = ({
  permissions,
  children,
}: {
  permissions: string[]
  children: React.ReactNode
}) => (
  <PermissionContext.Provider value={{ permissions }}>
    {children}
  </PermissionContext.Provider>
)

const routeToPermissionKey = (route: string): string => {
  // 去掉首尾斜杠
  const trimmed = route.replace(/^\/|\/$/g, '')
  // 如果为空，返回空字符串
  if (!trimmed) return ''
  // 用下划线拼接
  return trimmed.split('/').join('_')
}

export const usePermissions = () => {
  const { permissions } = useContext(PermissionContext)
  const match = useMatch({ strict: false })
  return permissions
    .filter(
      (item) => item.split(':')[0] === routeToPermissionKey(match.pathname)
    )
    .map((item) => item.split(':')[1]) as PermissionType[]
}

export const useHasPermission = (permission: PermissionType) => {
  const permissions = usePermissions()
  return permissions.includes(permission)
}

export type PermissionType =
  | 'query'
  | 'info'
  | 'create'
  | 'update'
  | 'delete'
  | 'start'
  | 'stop'
  | 'execute'
  | 'reset'

export const basicMoreOptions = [
  'update',
  'delete',
  'info',
] satisfies PermissionType[]

type SignalPermission = {
  permission: PermissionType
  children: React.ReactNode
  fallback?: React.ReactNode
}
type CallbackPermission = {
  children: (permissions: PermissionType[]) => React.ReactNode
}

type PermissionProps = SignalPermission | CallbackPermission

export const WithPermission: React.FC<PermissionProps> = ({
  children,
  ...props
}) => {
  const permissions = usePermissions()
  console.log(permissions)

  if (typeof children === 'function') {
    return children(permissions as PermissionType[])
  }

  if (permissions.includes((props as SignalPermission).permission)) {
    return children
  }

  return <>{(props as SignalPermission).fallback ?? null}</>
}
