import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { SysUser } from '~/do-exercise-api'
import {
  basicMoreOptions,
  PermissionType,
  usePermissions,
  WithPermission,
} from '~/provider'
import { ensureHttpPrefix } from '~/utils/url'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { DataTableRowActions } from '~/components/other'
import {
  createColumn,
  createDateColumn,
  createActionColumn,
} from '~/components/other/data-table/column-utils'
import { ResetPassword } from './user-reset-password-dialog'

const columnTitleMap = {
  id: (): string => t`序号`,
  username: (): string => t`用户名`,
  nickname: (): string => t`昵称`,
  email: (): string => t`邮箱`,
  avatar: (): string => t`头像`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<SysUser>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasMore = (
    [...basicMoreOptions, 'reset'] satisfies PermissionType[]
  ).some((p) => permissions.includes(p))

  return [
    {
      accessorKey: 'id',
      header: () => columnTitleMap.id(),
      cell: ({ row, table }) => {
        const pagination = table.getState().pagination
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
      cell: ({ row }: CellContext<SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.username}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'nickname',
      title: columnTitleMap.nickname,
      cell: ({ row }: CellContext<SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.nickname || '-'}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'email',
      title: columnTitleMap.email,
      cell: ({ row }: CellContext<SysUser, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.email || '-'}</div>
      ),
    }),
    createColumn<SysUser>({
      key: 'avatar',
      title: columnTitleMap.avatar,
      cell: ({ row }: CellContext<SysUser, unknown>) => {
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
    createDateColumn<SysUser>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysUser>('updatedAt', columnTitleMap.updatedAt),
    ...[
      hasMore
        ? createActionColumn<SysUser>(({ row }) => (
            <DataTableRowActions
              row={row}
              showEdit={permissions.some((p) => p === 'update')}
              showDelete={permissions.some((p) => p === 'delete')}
              showInfo={permissions.some((p) => p === 'info')}
            >
              {(user) => (
                <WithPermission permission='reset'>
                  <ResetPassword currentRow={user} />
                </WithPermission>
              )}
            </DataTableRowActions>
          ))
        : [],
    ].flat(),
  ]
}
