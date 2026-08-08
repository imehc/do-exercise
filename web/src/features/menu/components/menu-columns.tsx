import { ColumnDef, CellContext } from '@tanstack/react-table'
import { IconChevronDown, IconChevronRight } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { MenuType, SysMenuTree } from '~/do-exercise-api'
import { basicMoreOptions, usePermissions } from '~/provider'
import { toIconComponent } from '~/utils/icon'
import { Button } from '~/components/ui/button'
import { DataTableRowActions } from '~/components/other'
import {
  createColumn,
  createActionColumn,
  createBadgeColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'
import {
  getCallMenuMapping,
  callMenuTypes,
  callVisibleTypes,
} from '../data/data'

const translations = {
  menuName: (): string => t`菜单名称`,
  permission: (): string => t`权限标识`,
  icon: (): string => t`图标`,
  menuType: (): string => t`菜单类型`,
  route: (): string => t`路由`,
  component: (): string => t`组件`,
  sort: (): string => t`排序`,
  visible: (): string => t`是否可见`,
  yes: (): string => t`是`,
  no: (): string => t`否`,
} as const

export const getColumnTitle = (columnId: string): string =>
  translations[columnId as keyof typeof translations]?.() ?? columnId

export const useColumns = () => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasMore = basicMoreOptions.some((p) => permissions.includes(p))

  return [
    createColumn<SysMenuTree>({
      id: 'expander',
      header: () => null,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => {
        return row.getCanExpand() ? (
          <Button
            variant='ghost'
            className='h-6 w-6 p-0 hover:bg-transparent'
            onClick={row.getToggleExpandedHandler()}
          >
            {row.getIsExpanded() ? (
              <IconChevronDown className='h-4 w-4' />
            ) : (
              <IconChevronRight className='h-4 w-4' />
            )}
          </Button>
        ) : null
      },
      options: {
        minSize: 20,
        size: 30,
        maxSize: 40,
      },
    }),
    createColumn<SysMenuTree>({
      key: 'name',
      title: translations.menuName,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => (
        <div
          className='w-fit text-nowrap'
          style={{
            paddingLeft: `${row.depth * 12}px`,
          }}
        >
          {row.original.name}
        </div>
      ),
    }),
    createColumn<SysMenuTree>({
      key: 'permission',
      title: translations.permission,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => (
        <div className='w-fit text-nowrap'>
          {row.original.permission || '-'}
        </div>
      ),
    }),
    createColumn<SysMenuTree>({
      key: 'icon',
      title: translations.icon,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => {
        const SelectedIcon =
          row.original.icon && row.original.type === MenuType.menu
            ? toIconComponent(row.original.icon)
            : null
        if (!SelectedIcon) {
          return <div className='w-fit text-nowrap'>-</div>
        }
        return <SelectedIcon />
      },
    }),
    createBadgeColumn<SysMenuTree>(
      'type',
      translations.menuType,
      (value: unknown) => callMenuTypes.get(value as MenuType),
      (value: unknown) => getCallMenuMapping().get(value as MenuType) ?? '-'
    ),
    createColumn<SysMenuTree>({
      key: 'route',
      title: translations.route,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => (
        <div>{row.original.route || '-'}</div>
      ),
    }),
    createColumn<SysMenuTree>({
      key: 'component',
      title: translations.component,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => (
        <div>{row.original.component || '-'}</div>
      ),
    }),
    createColumn<SysMenuTree>({
      key: 'sort',
      title: translations.sort,
      cell: ({ row }: CellContext<DataTableFeatures, SysMenuTree, unknown>) => (
        <div>{row.original.sort ?? '-'}</div>
      ),
    }),
    createBadgeColumn<SysMenuTree>(
      'visible',
      translations.visible,
      (value: unknown) => callVisibleTypes.get(value as boolean),
      (value: unknown) => (value ? translations.yes() : translations.no())
    ),
    ...[
      hasMore
        ? createActionColumn<SysMenuTree>(({ row }) => (
            <DataTableRowActions
              row={row}
              showEdit={permissions.some((p) => p === 'update')}
              showDelete={permissions.some((p) => p === 'delete')}
              showInfo={permissions.some((p) => p === 'info')}
              showAdd={permissions.some((p) => p === 'create')}
            />
          ))
        : [],
    ].flat(),
  ] satisfies ColumnDef<DataTableFeatures, SysMenuTree>[]
}
