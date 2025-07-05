import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { TokenInfo } from '~/do-exercise-api'
import { PermissionType, usePermissions, WithPermission } from '~/provider'
import { DataTableRowActions, InlineCopy } from '~/components/other'
import {
  createColumn,
  createDateColumn,
  createActionColumn,
} from '~/components/other/data-table/column-utils'
import { ToggleDisabledSwitch } from './token-toggle-disabled-switch'

const columnTitleMap = {
  id: (): string => t`序号`,
  userId: (): string => t`用户ID`,
  username: (): string => t`用户名`,
  accessToken: (): string => t`令牌`,
  disabled: (): string => t`禁用状态`,
  accessTokenCreated: (): string => t`创建时间`,
  accessTokenExpired: (): string => t`到期时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<TokenInfo>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasUpdate = permissions.some((p) => p === 'update')
  const hasMore = (['update', 'delete'] satisfies PermissionType[]).some((p) =>
    permissions.includes(p)
  )
  return [
    // 选择框列
    // {
    //   id: 'select',
    //   header: ({ table }) => {
    //     const isAllPageRowsSelected = table.getIsAllPageRowsSelected()
    //     const isSomePageRowsSelected = table.getIsSomePageRowsSelected()
    //     return (
    //       <Checkbox
    //         checked={
    //           isSomePageRowsSelected && !isAllPageRowsSelected
    //             ? 'indeterminate'
    //             : isAllPageRowsSelected
    //         }
    //         onCheckedChange={(value) =>
    //           table.toggleAllPageRowsSelected(!!value)
    //         }
    //         aria-label='Select all'
    //         className='translate-y-[2px]'
    //       />
    //     )
    //   },
    //   cell: ({ row }) => (
    //     <Checkbox
    //       checked={row.getIsSelected()}
    //       onCheckedChange={(value) => row.toggleSelected(!!value)}
    //       aria-label='Select row'
    //       className='translate-y-[2px]'
    //     />
    //   ),
    //   enableSorting: false,
    //   enableHiding: false,
    //   meta: {
    //     className: 'w-10 text-center',
    //   },
    // },
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
    createColumn<TokenInfo>({
      key: 'userId',
      title: columnTitleMap.userId,
      cell: ({ row }: CellContext<TokenInfo, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.userId}</div>
      ),
    }),
    createColumn<TokenInfo>({
      key: 'username',
      title: columnTitleMap.username,
      cell: ({ row }: CellContext<TokenInfo, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.username}</div>
      ),
    }),
    createColumn<TokenInfo>({
      key: 'accessToken',
      title: columnTitleMap.accessToken,
      cell: ({ row }: CellContext<TokenInfo, unknown>) => (
        <InlineCopy
          className='w-fit text-nowrap'
          text={row.original.accessToken}
        />
      ),
    }),
    ...[
      hasUpdate
        ? createColumn<TokenInfo>({
            key: 'disabled',
            title: columnTitleMap.disabled,
            cell: ({ row }: CellContext<TokenInfo, unknown>) => {
              const token = row.original
              return (
                <ToggleDisabledSwitch
                  accessToken={token.accessToken}
                  disabled={token.disabled}
                />
              )
            },
          })
        : [],
    ].flat(),
    createDateColumn<TokenInfo>(
      'accessTokenCreated',
      columnTitleMap.accessTokenCreated
    ),
    createDateColumn<TokenInfo>(
      'accessTokenExpired',
      columnTitleMap.accessTokenExpired
    ),
    ...[
      hasMore
        ? createActionColumn<TokenInfo>(({ row }) => (
            <WithPermission>
              {(permissions) => (
                <DataTableRowActions
                  row={row}
                  showDelete={permissions.some((p) => p === 'delete')}
                  showInfo={permissions.some((p) => p === 'info')}
                />
              )}
            </WithPermission>
          ))
        : [],
    ].flat(),
  ]
}
