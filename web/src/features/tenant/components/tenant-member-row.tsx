import { useMemo, useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { IconDeviceFloppy, IconLoader3, IconTrash } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { SystemRoleBrief, SystemUserApi, SysUser } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Button } from '~/components/ui/button'
import { TableCell, TableRow } from '~/components/ui/table'
import { ConfirmDialog } from '~/components/confirm-dialog'
import { MultiSelect } from '~/components/other'

interface Props {
  member: SysUser
  /** 目标租户下可分配的角色（由对话框统一查询后传入，避免每行重复请求） */
  roles: SystemRoleBrief[]
  /** 角色调整或删除账号成功后刷新成员列表 */
  onChanged(): void
}

export function TenantMemberRow({ member, roles, onChanged }: Props) {
  const sysUserApi = useApi(SystemUserApi)
  const initialRoleIds = useMemo(
    () => (member.roles ?? []).map((role) => role.id),
    [member.roles]
  )
  const [roleIds, setRoleIds] = useState<number[]>(initialRoleIds)
  const [confirmRemove, setConfirmRemove] = useState(false)

  // 只有角色集合真的变了才允许提交，避免误触发全量覆盖
  const dirty =
    roleIds.length !== initialRoleIds.length ||
    roleIds.some((id) => !initialRoleIds.includes(id))

  const { isPending: isSaving, mutate: saveRoles } = useMutation({
    mutationFn: () =>
      sysUserApi.updateUser({
        id: member.id,
        // 更新接口整体覆盖，昵称与头像必须原样回传，否则会被清空
        updateSysUser: {
          nickname: member.nickname,
          avatar: member.avatar,
          roleIds,
        },
      }),
    onSuccess: () => {
      toast.success(t`角色已更新`)
      onChanged()
    },
    onError: () => {
      toast.error(t`角色更新失败`)
    },
  })

  // 「删除」在这里是字面意思：删掉该租户下的账号行（软删），不存在「解除成员关系」
  // 这回事——当前模型里一个人在每个租户下各有一个独立账号（P2-7）。
  const { isPending: isRemoving, mutate: removeMember } = useMutation({
    mutationFn: () => sysUserApi.deleteUser({ id: member.id }),
    onSuccess: () => {
      toast.success(t`已删除该租户下的账号`)
      setConfirmRemove(false)
      onChanged()
    },
    onError: () => {
      toast.error(t`删除账号失败`)
    },
  })

  const isPending = isSaving || isRemoving
  // 提到消息外面：lingui 的 msgid 用内联表达式会退化成位置占位符 {0}，
  // 命名变量才能让译者从 msgid 本身看出这里填的是什么
  const username = member.username

  return (
    <TableRow>
      <TableCell className='text-nowrap'>{member.username}</TableCell>
      <TableCell className='text-nowrap'>{member.nickname || '-'}</TableCell>
      <TableCell className='min-w-64'>
        <MultiSelect
          modalPopover
          options={roles.map((role) => ({
            label: role.name,
            value: role.id,
          }))}
          onValueChange={(values) =>
            setRoleIds(
              values.filter(
                (value): value is number => typeof value === 'number'
              )
            )
          }
          defaultValue={initialRoleIds}
          placeholder={t`请选择关联角色`}
          variant='inverted'
        />
      </TableCell>
      <TableCell>
        <div className='flex items-center gap-2'>
          <Button
            variant='outline'
            size='icon'
            title={t`保存角色`}
            disabled={!dirty || isPending}
            onClick={() => saveRoles()}
          >
            {isSaving ? (
              <IconLoader3 size={16} className='animate-spin' />
            ) : (
              <IconDeviceFloppy size={16} />
            )}
          </Button>
          <Button
            variant='outline'
            size='icon'
            title={t`删除该租户下的账号`}
            disabled={isPending}
            onClick={() => setConfirmRemove(true)}
          >
            <IconTrash size={16} className='text-red-500!' />
          </Button>
        </div>
      </TableCell>
      <ConfirmDialog
        open={confirmRemove}
        onOpenChange={setConfirmRemove}
        isLoading={isRemoving}
        handleConfirm={() => removeMember()}
        destructive
        title={<Trans>删除该租户下的账号</Trans>}
        desc={t`该操作删除的是账号本身，不是解除一段成员关系：账号「${username}」在本租户下的记录会被删除，其会话立即失效。同名账号在其他租户下是独立的账号，不受影响。`}
        confirmText={<Trans>删除账号</Trans>}
      />
    </TableRow>
  )
}
