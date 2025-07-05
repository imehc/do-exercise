import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'
import { JobStatus, SysJob } from '~/do-exercise-api'
import { basicMoreOptions, usePermissions, WithPermission } from '~/provider'
import { DataTableRowActions, InlineCopy } from '~/components/other'
import {
  createColumn,
  createDateColumn,
  createActionColumn,
  createBadgeColumn,
} from '~/components/other/data-table/column-utils'
import { Row } from '~/features/api/components/api-columns'
import { scheduleStatusTypes } from '../data/data'
import { TaskActionButton } from './task-action-button'

const columnTitleMap = {
  id: (): string => t`序号`,
  name: (): string => t`任务名称`,
  jobGroup: (): string => t`任务分组`,
  cronExpression: (): string => t`cron表达式`,
  command: (): string => t`执行命令`,
  status: (): string => t`状态`,
  description: (): string => t`任务描述`,
  createdAt: (): string => t`创建时间`,
  updatedAt: (): string => t`更新时间`,
}

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (refresh?: () => void): ColumnDef<SysJob>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  const hasMore = basicMoreOptions.some((p) => permissions.includes(p))

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
    createColumn<SysJob>({
      key: 'name',
      title: columnTitleMap.name,
      cell: ({ row }: CellContext<SysJob, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.name}</div>
      ),
    }),
    createColumn<SysJob>({
      key: 'jobGroup',
      title: columnTitleMap.jobGroup,
      cell: ({ row }: CellContext<SysJob, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.jobGroup || '-'}</div>
      ),
    }),
    createColumn<SysJob>({
      key: 'cronExpression',
      title: columnTitleMap.cronExpression,
      cell: ({ row }: CellContext<SysJob, unknown>) => (
        <InlineCopy
          className='w-fit text-nowrap'
          text={row.original.cronExpression}
        />
      ),
    }),
    createColumn<SysJob>({
      key: 'command',
      title: columnTitleMap.command,
      cell: ({ row }: CellContext<SysJob, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.command || '-'}</div>
      ),
    }),
    createBadgeColumn<SysJob>(
      'status',
      columnTitleMap.status,
      (value: unknown) => scheduleStatusTypes.get(value as JobStatus),
      (value: unknown) =>
        value === JobStatus.normal ? <Trans>正常</Trans> : <Trans>暂停</Trans>,
      {
        filterFn: (row: Row, id: string, value: Array<string>) =>
          value.includes(row.getValue(id)),
        enableHiding: false,
        enableSorting: false,
      }
    ),
    createDateColumn<SysJob>('createdAt', columnTitleMap.createdAt),
    createDateColumn<SysJob>('updatedAt', columnTitleMap.updatedAt),
    {
      accessorKey: 'start',
      header: '',
      cell: ({ row }) => (
        <WithPermission permission='start'>
          <TaskActionButton row={row.original} type='start' refresh={refresh} />
        </WithPermission>
      ),
    },
    {
      accessorKey: 'stop',
      header: '',
      cell: ({ row }) => (
        <WithPermission permission='stop'>
          <TaskActionButton row={row.original} type='stop' refresh={refresh} />
        </WithPermission>
      ),
    },
    {
      accessorKey: 'execute',
      header: '',
      cell: ({ row }) => (
        <WithPermission permission='execute'>
          <TaskActionButton
            row={row.original}
            type='execute'
            refresh={refresh}
          />
        </WithPermission>
      ),
    },
    ...[
      hasMore
        ? createActionColumn<SysJob>(({ row }) => (
          <WithPermission>
            {(permissions) => (
              <DataTableRowActions
                row={row}
                showEdit={permissions.some((p) => p === 'update')}
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
