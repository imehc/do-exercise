import { useMemo } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  CreateRoleRequest,
  MenuType,
  SysMenuTree,
  SysRole,
  SystemRoleApi,
  UpdateRoleRequest,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
} from '~/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { StatusRenderer, transformData, TreeSelect } from '~/components/other'
import { DialogHeaderContent } from '~/components/other/dialog-header-content'
import { findMenuTree } from '~/features/menu/data/api'
import { getMenuLabel } from '~/utils/menu-label'
import { ActionSysRoleFormValues, getSchema } from '../schemas/action-schema'
import {
  collectSelectedApis,
  collectSelectedPermissions,
  diffMenuSelection,
  findMissingAncestors,
  findPagesWithoutActions,
  indexMenuTree,
} from '../utils/permission-diff'
import { RolePermissionSummary } from './role-permission-summary'

/** 按钮节点在树里直接带上权限标识，授权人不必再点进菜单详情去核对 */
function menuTreeLabel(node: SysMenuTree): string {
  const label = getMenuLabel(node)
  if (node.type === MenuType.button && node.permission) {
    return `${label} · ${node.permission}`
  }
  return label
}

interface Props {
  currentRow?: SysRole
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function RoleActionDialog({ currentRow, open, onOpenChange }: Props) {
  const sysRoleApi = useApi(SystemRoleApi)

  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useChan(
    useForm<ActionSysRoleFormValues>({
      defaultValues: isEdit
        ? {
            ...currentRow,
            menuIds: currentRow?.menus?.map((item) => item.id) ?? [],
          }
        : { name: '', code: '' },
      resolver: zodResolver(getSchema()),
    })
  )

  const { data = [], isLoading } = useQuery(findMenuTree())

  // 保存前的授权基线：编辑时才有，用来算「本次变更」
  const originalMenuIds = useMemo(
    () => currentRow?.menus?.map((item) => item.id) ?? [],
    [currentRow?.menus]
  )
  const selectedMenuIds = useWatch({ control: form.control, name: 'menuIds' })
  const menuIndex = useMemo(() => indexMenuTree(data, menuTreeLabel), [data])
  const selected = useMemo(() => selectedMenuIds ?? [], [selectedMenuIds])
  const missingAncestors = useMemo(
    () => findMissingAncestors(selected, menuIndex),
    [selected, menuIndex]
  )
  const pagesWithoutActions = useMemo(
    () => findPagesWithoutActions(selected, menuIndex),
    [selected, menuIndex]
  )
  const selectedPermissions = useMemo(
    () => collectSelectedPermissions(selected, menuIndex),
    [selected, menuIndex]
  )
  const selectedApis = useMemo(
    () => collectSelectedApis(selected, menuIndex),
    [selected, menuIndex]
  )
  const diff = useMemo(
    () =>
      isEdit
        ? diffMenuSelection(
            originalMenuIds,
            selected,
            menuIndex,
            (id) => t`已失效菜单 #${id}`
          )
        : null,
    [isEdit, originalMenuIds, selected, menuIndex]
  )

  const handleFixAncestors = () => {
    form.setValue(
      'menuIds',
      [...new Set([...selected, ...missingAncestors.map((item) => item.id)])],
      { shouldDirty: true, shouldValidate: true }
    )
  }

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateRoleRequest) => sysRoleApi.createRole(values),
    onSuccess: () => {
      toast.success(t`创建成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateRoleRequest) => sysRoleApi.updateRole(values),
    onSuccess: () => {
      toast.success(t`更新成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const onSubmit = (values: ActionSysRoleFormValues) => {
    if (isEdit) {
      saveUpdate({
        updateSysRole: { name: values.name, menuIds: values.menuIds },
        id: currentRow?.id as number,
      })
      return
    }
    saveCreate({
      createSysRole: values,
    })
  }

  const isPending = isPendingCreate || isPendingUpdate

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        form.reset()
        onOpenChange(state, false)
      }}
    >
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader className='text-left'>
          <DialogHeaderContent isEdit={isEdit} text={<Trans>角色</Trans>} />
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='max-h-[65vh] space-y-4 overflow-y-auto p-0.5'
          >
            <FormField
              control={form.control}
              name='name'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    <Trans>角色名称</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入角色名称`}
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      value={field.value ?? ''}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='code'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    <Trans>角色编码</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      readOnly={isEdit}
                      disabled={isEdit}
                      placeholder={t`请输入角色编码`}
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      value={field.value ?? ''}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            <StatusRenderer isLoading={isLoading} data={data}>
              {(treeData) => (
                <>
                  <FormField
                    control={form.control}
                    name='menuIds'
                    render={({ field }) => (
                      <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                        <FormLabel className='col-span-2 text-right'>
                          <Trans>关联菜单</Trans>
                        </FormLabel>
                        <FormControl>
                          <TreeSelect
                            className='col-span-8'
                            mode='inline'
                            multiple
                            data={transformData(
                              treeData,
                              menuTreeLabel,
                              (item) => item.id
                            )}
                            value={field.value ?? []}
                            onChange={field.onChange}
                          />
                        </FormControl>
                        <FormMessage className='col-span-8 col-start-3' />
                      </FormItem>
                    )}
                  />
                  <div className='grid grid-cols-10 gap-x-4'>
                    <div className='col-span-2 pt-3 text-right text-sm'>
                      <Trans>影响面</Trans>
                    </div>
                    <RolePermissionSummary
                      permissions={selectedPermissions}
                      apis={selectedApis}
                      missingAncestors={missingAncestors}
                      pagesWithoutActions={pagesWithoutActions}
                      diff={diff}
                      userCount={isEdit ? (currentRow?.userCount ?? 0) : null}
                      onFixAncestors={handleFixAncestors}
                    />
                  </div>
                </>
              )}
            </StatusRenderer>
          </form>
        </Form>
        <DialogFooter>
          <Button type='submit' form='menu-form' disabled={isPending}>
            {isPending ? (
              <>
                <IconLoader3 className='animate-spin' />
                <span>
                  <Trans>保存中</Trans>...
                </span>
              </>
            ) : (
              <span>
                <Trans>保存</Trans>
              </span>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
