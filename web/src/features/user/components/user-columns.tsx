import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom, originTokenAtom } from '~/atoms'
import { SysUser } from '~/do-exercise-api'
import {
  basicMoreOptions,
  PermissionType,
  usePermissions,
  WithPermission,
} from '~/provider'
import { formatTenant } from '~/utils/tenant'
import { ensureHttpPrefix } from '~/utils/url'
import { useUserProfile } from '~/hooks/use-user'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { DataTableRowActions } from '~/components/other'
import {
  createColumn,
  createDateColumn,
  createActionColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'
import { AssignTenant } from './user-assign-tenant-dialog'
import { ResetPassword } from './user-reset-password-dialog'

const columnTitleMap = {
  id: (): string => t`序号`,
  username: (): string => t`用户名`,
  nickname: (): string => t`昵称`,
  email: (): string => t`邮箱`,
  avatar: (): string => t`头像`,
  tenantName: (): string => t`所属租户`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<DataTableFeatures, SysUser>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  // 当前操作者是否为平台超级管理员（超管才能在用户列表执行「分配租户」）
  const isSuperAdmin = !!useAtomValue(originTokenAtom).isSuperAdmin
  // 当前登录账号，用于禁用「删除自己」——服务端同样拦截（cannotDeleteSelf），
  // 这里只是不把注定失败的入口摆在用户面前。
  const { data: profile } = useUserProfile()
  const hasMore = (
    [...basicMoreOptions, 'reset'] satisfies PermissionType[]
  ).some((p) => permissions.includes(p))

  return [
    {
      accessorKey: 'id',
      header: () => columnTitleMap.id(),
      cell: ({ row, table }) => {
        const pagination = table.atoms.pagination.get()
        return (
          (pagination.pageIndex ?? 0) * (pagination.pageSize ?? 0) +
          row.index +
          1
        )
      },
    },
    createColumn<SysUser>({
      key: 'username',
      title: columnTitleMap.username,
      cell: ({ row }: CellContext<DataTableFeatures, SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.username}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'nickname',
      title: columnTitleMap.nickname,
      cell: ({ row }: CellContext<DataTableFeatures, SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.nickname || '-'}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'email',
      title: columnTitleMap.email,
      cell: ({ row }: CellContext<DataTableFeatures, SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.email || '-'}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'avatar',
      title: columnTitleMap.avatar,
      cell: ({ row }: CellContext<DataTableFeatures, SysUser, unknown>) => {
        const avatar = row.original.avatar
        if (!avatar) {
          return '-'
        }
        return (
          <div className='w-fit text-nowrap'>
            <Avatar>
              <AvatarImage
                src={ensureHttpPrefix(avatar)}
                alt={row.original.username}
              />
              <AvatarFallback>
                {row.original.username.slice(0, 2)}
              </AvatarFallback>
            </Avatar>
          </div>
        )
      },
    }),
    // 所属租户仅对平台超级管理员有意义：受限管理员的列表已被行级隔离锁在本租户内，
    // 每行都是同一个租户，展示出来只是噪音。
    ...(isSuperAdmin
      ? [
          createColumn<SysUser>({
            key: 'tenantName',
            title: columnTitleMap.tenantName,
            cell: ({
              row,
            }: CellContext<DataTableFeatures, SysUser, unknown>) => (
              <div className='w-fit text-nowrap'>
                {formatTenant(row.original.tenantId, row.original.tenantName)}
              </div>
            ),
          }),
        ]
      : []),
    createDateColumn<SysUser>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysUser>('updatedAt', columnTitleMap.updatedAt),
    ...[
      hasMore
        ? createActionColumn<SysUser>(({ row }) => (
            <DataTableRowActions
              row={row}
              showEdit={permissions.some((p) => p === 'update')}
              showDelete={
                permissions.some((p) => p === 'delete') &&
                row.original.id !== profile?.id
              }
              showInfo={permissions.some((p) => p === 'info')}
            >
              {(user) => (
                <>
                  <WithPermission permission='reset'>
                    <ResetPassword currentRow={user} />
                  </WithPermission>
                  {isSuperAdmin && <AssignTenant currentRow={user} />}
                </>
              )}
            </DataTableRowActions>
          ))
        : [],
    ].flat(),
  ]
}
