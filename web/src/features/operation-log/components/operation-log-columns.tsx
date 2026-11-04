import { ColumnDef, CellContext } from '@tanstack/react-table'
import { t } from '@lingui/core/macro'
import { useAtomValue } from 'jotai'
import { languageAtom, originTokenAtom } from '~/atoms'
import { SysOperationLog } from '~/do-exercise-api'
import { usePermissions, WithPermission } from '~/provider'
import { formatTenant } from '~/utils/tenant'
import { DataTableRowActions, Status } from '~/components/other'
import {
  createColumn,
  createActionColumn,
  createBadgeColumn,
  createDateColumn,
} from '~/components/other/data-table/column-utils'
import { DataTableFeatures } from '~/components/other/data-table/features'
import type { Row } from '~/features/api/components/api-columns'
import { callMethodTypes } from '~/features/api/data/data'
import { callCodeTypes } from '../data/data'
import { OperationLogViewDialog } from './operation-log-view-dialog'

const columnTitleMap = {
  id: (): string => t`序号`,
  username: (): string => t`请求用户`,
  tenantId: (): string => t`所属租户`,
  ip: (): string => t`IP地址`,
  address: (): string => t`IP属地`,
  os: (): string => t`系统`,
  browser: (): string => t`浏览器`,
  path: (): string => t`请求路径`,
  method: (): string => t`请求方法`,
  success: (): string => t`是否成功`,
  code: (): string => t`状态码`,
  message: (): string => t`错误信息`,
  params: (): string => t`请求参数`,
  body: (): string => t`请求体`,
  result: (): string => t`响应结果`,
  startTime: (): string => t`请求开始时间`,
  endTime: (): string => t`请求结束时间`,
  latency: (): string => t`耗时`,
  successText: (): string => t`成功`,
  failureText: (): string => t`失败`,
} as const

export const getColumnTitle = (columnId: string): string =>
  columnTitleMap[columnId as keyof typeof columnTitleMap]?.() ?? columnId

export const useColumns = (): ColumnDef<
  DataTableFeatures,
  SysOperationLog
>[] => {
  useAtomValue(languageAtom)
  const permissions = usePermissions()
  // 超管的日志列表跨租户，需要一列标明操作来源；受限管理员被行级隔离锁在本租户，
  // 该列恒为同一个值。
  const isSuperAdmin = !!useAtomValue(originTokenAtom).isSuperAdmin
  const hasMore = permissions.some((p) => p === 'info')
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
    createColumn<SysOperationLog>({
      key: 'username',
      title: columnTitleMap.username,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.username || '-'}</div>
      ),
    }),
    ...(isSuperAdmin
      ? [
          createColumn<SysOperationLog>({
            key: 'tenantId',
            title: columnTitleMap.tenantId,
            cell: ({
              row,
            }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
              <div className='w-fit text-nowrap'>
                {formatTenant(row.original.tenantId)}
              </div>
            ),
          }),
        ]
      : []),
    createColumn<SysOperationLog>({
      key: 'ip',
      title: columnTitleMap.ip,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.ip}</div>
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'address',
      title: columnTitleMap.address,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>
          {row.original.address
            ?.split('|')
            .filter((part) => part && !/^\d+$/.test(part))
            .slice(0, -1)
            .join(' ') || '-'}
        </div>
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'os',
      title: columnTitleMap.os,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.os}</div>
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'browser',
      title: columnTitleMap.browser,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.browser}</div>
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'path',
      title: columnTitleMap.path,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div className='w-fit text-nowrap'>{row.original.path}</div>
      ),
    }),
    createBadgeColumn<SysOperationLog>(
      'method',
      columnTitleMap.method,
      (value: unknown) => callMethodTypes.get(value as string),
      (value: unknown) => value as string,
      {
        enableSorting: true,
      }
    ),
    createColumn<SysOperationLog>({
      key: 'success',
      title: columnTitleMap.success,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => {
        const { success } = row.original
        return (
          <div className='flex space-x-2'>
            <Status
              color={success ? 'success' : 'error'}
              label={
                success
                  ? columnTitleMap.successText()
                  : columnTitleMap.failureText()
              }
            />
          </div>
        )
      },
      options: {
        filterFn: (row: Row, id: string, value: Array<string>) =>
          value.includes(row.getValue(id)),
        enableHiding: false,
        enableSorting: false,
      },
    }),
    createBadgeColumn<SysOperationLog>(
      'code',
      columnTitleMap.code,
      (value: unknown) => callCodeTypes.get(value as number),
      (value: unknown) => value as string
    ),
    createColumn<SysOperationLog>({
      key: 'message',
      title: columnTitleMap.message,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <OperationLogViewDialog data={row.original} type='view-msg' />
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'params',
      title: columnTitleMap.params,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <OperationLogViewDialog data={row.original} type='view-params' />
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'body',
      title: columnTitleMap.body,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <OperationLogViewDialog data={row.original} type='view-body' />
      ),
    }),
    createColumn<SysOperationLog>({
      key: 'result',
      title: columnTitleMap.result,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <OperationLogViewDialog data={row.original} type='view-result' />
      ),
    }),
    createDateColumn<SysOperationLog>('startTime', columnTitleMap.startTime),
    createDateColumn<SysOperationLog>('endTime', columnTitleMap.endTime),
    createColumn<SysOperationLog>({
      key: 'latency',
      title: columnTitleMap.latency,
      cell: ({
        row,
      }: CellContext<DataTableFeatures, SysOperationLog, unknown>) => (
        <div>{row.original.latency}ms</div>
      ),
    }),
    ...[
      hasMore
        ? createActionColumn<SysOperationLog>(({ row }) => (
            <WithPermission>
              {(permissions) => (
                <DataTableRowActions
                  row={row}
                  showInfo={permissions.some((p) => p === 'info')}
                />
              )}
            </WithPermission>
          ))
        : [],
    ].flat(),
  ]
}
