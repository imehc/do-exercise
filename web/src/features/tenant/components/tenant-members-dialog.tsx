import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { IconUsersGroup } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { SystemRoleApi, SystemUserApi, Tenant } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from '~/components/ui/table'
import { LoadingSpinner } from '~/components/other'
import { TenantMemberRow } from './tenant-member-row'

interface Props {
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
  currentRow: Tenant
}

const PAGE_SIZE = 10

/**
 * 租户成员管理：列出目标租户下的全部成员并就地调整角色 / 删除账号。
 * 复用用户接口的 tenant_id 筛选参数（仅平台超级管理员可用），
 * 角色候选同样按目标租户查询，避免把别的租户的角色分配过去。
 */
export function TenantMembersDialog({ open, onOpenChange, currentRow }: Props) {
  const sysUserApi = useApi(SystemUserApi)
  const sysRoleApi = useApi(SystemRoleApi)
  const tenantId = currentRow.tenantId ?? ''
  const [page, setPage] = useState(1)

  // 换租户或重新打开时回到第一页，避免沿用上一个租户的页码。
  // 用渲染期比对而不是 useEffect：effect 里 setState 会先用旧页码发一次请求再重渲染。
  const [pageScope, setPageScope] = useState({ open, tenantId })
  if (pageScope.open !== open || pageScope.tenantId !== tenantId) {
    setPageScope({ open, tenantId })
    setPage(1)
  }

  const {
    data: members,
    isLoading: membersIsLoading,
    refetch,
  } = useQuery({
    queryKey: ['findUsers', tenantId, page],
    queryFn: () =>
      sysUserApi.findUsers({ page, pageSize: PAGE_SIZE, tenantId }),
    enabled: open && !!tenantId,
  })

  const { data: roles = [], isLoading: rolesIsLoading } = useQuery({
    queryKey: ['findAllRoles', tenantId],
    queryFn: () => sysRoleApi.findAllRoles({ tenantId }),
    enabled: open && !!tenantId,
  })

  const total = members?.meta.total ?? 0
  const rows = members?.data ?? []
  const maxPage = Math.max(1, Math.ceil(total / PAGE_SIZE))
  const isLoading = membersIsLoading || rolesIsLoading
  // 提到消息外面：内联表达式会让 msgid 退化成位置占位符 {0}
  const tenantName = currentRow.name ?? ''

  return (
    <Dialog open={open} onOpenChange={(state) => onOpenChange(state, false)}>
      <DialogContent className='sm:max-w-4xl'>
        <DialogHeader className='text-left'>
          <DialogTitle className='flex items-center gap-2'>
            <IconUsersGroup size={18} />
            <Trans>成员管理</Trans>
          </DialogTitle>
          <DialogDescription>
            <Trans>
              查看租户「{tenantName}
              」下的账号，并调整其角色或删除该租户下的账号。
            </Trans>
          </DialogDescription>
        </DialogHeader>
        <div className='-mr-4 max-h-[26.25rem] w-full overflow-y-auto py-1 pr-4'>
          {isLoading ? (
            <div className='flex h-40 items-center justify-center'>
              <LoadingSpinner />
            </div>
          ) : rows.length === 0 ? (
            <div className='text-muted-foreground flex h-40 items-center justify-center text-sm'>
              <Trans>该租户下暂无成员</Trans>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>
                    <Trans>用户名</Trans>
                  </TableHead>
                  <TableHead>
                    <Trans>昵称</Trans>
                  </TableHead>
                  <TableHead>
                    <Trans>关联角色</Trans>
                  </TableHead>
                  <TableHead>
                    <Trans>操作</Trans>
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((member) => (
                  <TenantMemberRow
                    // 角色调整后重新拉取，key 带上角色集合以重置行内选择状态
                    key={`${member.id}-${(member.roles ?? [])
                      .map((role) => role.id)
                      .join('_')}`}
                    member={member}
                    roles={roles}
                    // 租户列表本身不展示成员信息，成员变动只需刷新对话框内的列表
                    onChanged={() => refetch()}
                  />
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        <DialogFooter className='sm:justify-between'>
          <div className='text-muted-foreground flex items-center gap-2 text-sm'>
            <Trans>共 {total} 人</Trans>
            <Button
              variant='outline'
              size='sm'
              disabled={page <= 1 || isLoading}
              onClick={() => setPage((prev) => Math.max(1, prev - 1))}
            >
              <Trans>上一页</Trans>
            </Button>
            <span>
              {page} / {maxPage}
            </span>
            <Button
              variant='outline'
              size='sm'
              disabled={page >= maxPage || isLoading}
              onClick={() => setPage((prev) => Math.min(maxPage, prev + 1))}
            >
              <Trans>下一页</Trans>
            </Button>
          </div>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false, false)}
          >
            <Trans>关闭</Trans>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
