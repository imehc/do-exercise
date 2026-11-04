import { useEffect, useMemo, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { SwitchThumb } from '@radix-ui/react-switch'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  CreateMenuRequest,
  CreateSysMenu,
  MenuType,
  SysMenuTree,
  SysMenuWithButton,
  SystemApiApi,
  SystemMenuApi,
  UpdateMenuRequest,
} from '~/do-exercise-api'
import { router } from '~/main'
import { useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
import { getMenuLabel } from '~/utils/menu-label'
import { useApi } from '~/hooks/use-api'
import { useMenuPermissionActions } from '~/hooks/use-tenant'
import { Badge } from '~/components/ui/badge'
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
import { Switch } from '~/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { ConfirmDialog } from '~/components/confirm-dialog'
import {
  IconSelect,
  StatusRenderer,
  TransferList,
  transformData,
  TreeSelect,
} from '~/components/other'
import { DialogHeaderContent } from '~/components/other/dialog-header-content'
import { SelectDropdown } from '~/components/select-dropdown'
import { callMethodTypes } from '~/features/api/data/data'
import { getCallMenuMapping } from '../data/data'
import {
  ActionSysMenuFormValues,
  ActionSysMenuWithDirectoryFormValues,
  ActionSysMenuWithMenuFormValues,
  fallbackMenuPermissionActions,
  getActionSysMenuSchema,
  getMenuPermissionActionLabel,
  getMenuScopeLabel,
  menuScopes,
} from '../schemas/action-schema'
import { buildMenuPermission, routePermissionKey } from '../utils/permission'

function toCamelCase(str?: string) {
  return str
    ?.split('-')
    .map((word, index) => {
      if (index === 0) return word.charAt(0).toUpperCase() + word.slice(1)
      return word.charAt(0).toUpperCase() + word.slice(1)
    })
    .join('')
}

function findMenuRoute(nodes: SysMenuTree[], id: number): string | undefined {
  for (const node of nodes) {
    if (node.id === id) return node.route
    if (node.children?.length) {
      const route = findMenuRoute(node.children, id)
      if (route !== undefined) return route
    }
  }
  return undefined
}

interface Props {
  treeData: SysMenuTree[]
  currentRow?: SysMenuTree
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function MenuActionDialog({
  treeData,
  currentRow,
  open,
  onOpenChange,
}: Props) {
  const { open: openType } = useFormDialog()
  const routes = useMemo<Array<string>>(
    () =>
      Object.values(router.routesByPath)
        .map((item) => item.fullPath.trim())
        .map((item) => (item.length > 1 ? item.replace(/\/+$/, '') : item)),
    []
  )

  const sysMenuApi = useApi(SystemMenuApi)
  const sysApi = useApi(SystemApiApi)

  const backendActions = useMenuPermissionActions()
  // 词表以后端下发为准；未到达时用兜底表，避免首屏下拉为空
  const permissionActions = useMemo(
    () =>
      backendActions.length
        ? backendActions
        : [...fallbackMenuPermissionActions],
    [backendActions]
  )

  const isEdit = useMemo(() => !!currentRow?.id, [currentRow])
  const form = useForm<ActionSysMenuFormValues>({
    resolver: zodResolver(getActionSysMenuSchema(permissionActions)),
  })
  // 平台专属是授权边界变更，提交前二次确认（后端另有审计日志）
  const [pendingPlatformValues, setPendingPlatformValues] =
    useState<ActionSysMenuFormValues | null>(null)

  const menuType = useWatch({ control: form.control, name: 'type' })
  const parentId = useWatch({ control: form.control, name: 'parentId' })
  const permissionAction = useWatch({
    control: form.control,
    name: 'permissionAction',
  })
  const scope = useWatch({ control: form.control, name: 'scope' })
  const permissionPrefix = routePermissionKey(
    typeof parentId === 'number' ? findMenuRoute(treeData, parentId) : undefined
  )
  const generatedPermission =
    permissionPrefix && permissionAction
      ? `${permissionPrefix}:${permissionAction}`
      : ''
  const isButtonType = menuType === MenuType.button
  const isActionOpen =
    openType === 'add' || openType === 'edit' || openType === 'add-child'

  const { data: apis = [], isLoading: isLoadingApis } = useQuery({
    queryKey: ['findAllApis'],
    queryFn: () => sysApi.findAllApis(),
    enabled: isButtonType && isActionOpen,
  })

  const { data: menu, isLoading: isLoadingMenu } = useQuery({
    queryKey: ['findMenu', currentRow?.id],
    queryFn: () => sysMenuApi.findMenu({ id: currentRow?.id as number }),
    enabled: isButtonType && isEdit && isActionOpen,
  })

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateMenuRequest) => sysMenuApi.createMenu(values),
    onSuccess: () => {
      toast.success(t`创建成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateMenuRequest) => sysMenuApi.updateMenu(values),
    onSuccess: () => {
      toast.success(t`更新成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  useEffect(() => {
    if (isLoadingMenu) return
    if (isEdit) {
      switch (currentRow?.type) {
        case MenuType.directory:
          form.reset({
            type: MenuType.directory,
            sort: currentRow.sort,
            parentId: currentRow.parentId,
            name: currentRow.name,
            i18nKey: currentRow.i18nKey,
            scope: currentRow.scope,
            visible: currentRow.visible,
          })
          break
        case MenuType.menu:
          form.reset({
            type: MenuType.menu,
            sort: currentRow.sort,
            parentId: currentRow.parentId,
            name: currentRow.name,
            i18nKey: currentRow.i18nKey,
            scope: currentRow.scope,
            icon: currentRow.icon,
            route: currentRow.route,
            visible: currentRow.visible,
          })
          break
        case MenuType.button:
          form.reset({
            parentId: currentRow.parentId,
            name: currentRow.name,
            i18nKey: currentRow.i18nKey,
            scope: currentRow.scope,
            type: MenuType.button,
            permissionAction: currentRow.permission?.split(':')[1] ?? 'query',
            apiIds:
              (menu as SysMenuWithButton)?.apis?.map((api) => api.id) ?? [],
            visible: currentRow.visible,
          })
      }
      return
    }
    if (openType === 'add-child') {
      form.reset({
        parentId: currentRow?.parentId,
        type: currentRow?.type,
      } as ActionSysMenuFormValues)
      return
    }
    form.reset({
      type: MenuType.directory,
      sort: 0,
    })
  }, [
    currentRow?.icon,
    currentRow?.i18nKey,
    currentRow?.name,
    currentRow?.parentId,
    currentRow?.permission,
    currentRow?.route,
    currentRow?.scope,
    currentRow?.sort,
    currentRow?.type,
    currentRow?.visible,
    form,
    isEdit,
    isLoadingMenu,
    menu,
    openType,
  ])

  const handleFormat = (values: ActionSysMenuFormValues) => {
    switch (values.type) {
      case MenuType.directory:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          i18nKey: values.i18nKey,
          scope: values.scope,
          sort: values.sort,
          visible: Boolean(values.visible),
        } satisfies ActionSysMenuWithDirectoryFormValues
      case MenuType.menu:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          i18nKey: values.i18nKey,
          scope: values.scope,
          icon: values.icon,
          route: values.route,
          sort: values.sort,
          visible: Boolean(values.visible),
        } satisfies ActionSysMenuWithMenuFormValues
      case MenuType.button:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          i18nKey: values.i18nKey,
          scope: values.scope,
          permission: buildMenuPermission(
            findMenuRoute(treeData, values.parentId),
            values.permissionAction
          ),
          apiIds: values.apiIds,
          visible: Boolean(values.visible),
        }
    }
  }

  const submitValues = (values: ActionSysMenuFormValues) => {
    if (isEdit) {
      saveUpdate({
        createSysMenu: handleFormat(values) as unknown as CreateSysMenu,
        id: currentRow?.id as number,
      })
      return
    }
    saveCreate({
      createSysMenu: handleFormat(values) as unknown as CreateSysMenu,
    })
  }

  const onSubmit = (values: ActionSysMenuFormValues) => {
    if (values.type === MenuType.button && !permissionPrefix) {
      toast.error(t`父级菜单必须绑定有效路由`)
      return
    }
    // 只在「新落到 platform」时确认：已经是平台菜单的改名、改排序不该反复打断
    if (values.scope === 'platform' && currentRow?.scope !== 'platform') {
      setPendingPlatformValues(values)
      return
    }
    submitValues(values)
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
          <DialogHeaderContent isEdit={isEdit} text={<Trans>菜单</Trans>} />
          {currentRow?.isSystem && (
            <div className='text-muted-foreground flex items-center gap-2 text-xs'>
              <Badge variant='outline'>
                <Trans>系统内置</Trans>
              </Badge>
              <span>
                <Trans>
                  路由、权限标识、菜单类型属于平台契约，提交后仍保持原值
                </Trans>
              </span>
            </div>
          )}
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <Tabs
              value={menuType?.toString()}
              onValueChange={(value) =>
                form.setValue('type', +value as MenuType)
              }
            >
              <TabsList className='grid w-full grid-cols-3'>
                {Array.from(getCallMenuMapping().entries())
                  .map(([key, value]) => ({
                    key,
                    value,
                  }))
                  .map((item) => (
                    <TabsTrigger key={item.key} value={item.key.toString()}>
                      {item.value}
                    </TabsTrigger>
                  ))}
              </TabsList>
              <FormField
                control={form.control}
                name='parentId'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>父级菜单</Trans>
                    </FormLabel>
                    <FormControl>
                      <TreeSelect
                        className='col-span-8'
                        placeholder={t`请选择父级菜单`}
                        valueMode='parent-only'
                        data={[
                          {
                            name: t`根节点`,
                            value: 0,
                            children: transformData(
                              treeData,
                              (item) => getMenuLabel(item),
                              (item) => item.id
                            ),
                          },
                        ]}
                        value={[field.value].filter(
                          (item) => typeof item === 'number'
                        )}
                        onChange={(value) =>
                          field.onChange(
                            typeof value[0] === 'number' ? value[0] : undefined
                          )
                        }
                      />
                    </FormControl>
                    <FormMessage className='col-span-8 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>菜单名称</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t`请输入菜单名称`}
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
                name='i18nKey'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>国际化键</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t`请输入稳定的翻译键（可选）`}
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
                name='scope'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>可见范围</Trans>
                    </FormLabel>
                    <FormControl>
                      <SelectDropdown
                        defaultValue={field.value ?? 'both'}
                        isControlled
                        onValueChange={field.onChange}
                        placeholder={t`请选择可见范围`}
                        className='col-span-8 w-full'
                        disabled={currentRow?.isSystem}
                        items={menuScopes.map((item) => ({
                          label: getMenuScopeLabel(item),
                          value: item,
                        }))}
                      />
                    </FormControl>
                    <div className='text-muted-foreground col-span-8 col-start-3 text-xs'>
                      {currentRow?.isSystem ? (
                        <Trans>
                          系统内置菜单的可见范围由平台维护，不可在此修改
                        </Trans>
                      ) : scope === 'platform' ? (
                        <Trans>
                          仅平台：该菜单及其权限不会下发给任何业务租户，保存时需要二次确认
                        </Trans>
                      ) : (
                        <Trans>
                          决定该菜单能否下发给业务租户，留空按「平台与租户」处理
                        </Trans>
                      )}
                    </div>
                    <FormMessage className='col-span-8 col-start-3' />
                  </FormItem>
                )}
              />
              <TabsContent
                value={MenuType.menu.toString()}
                className='flex flex-col gap-2'
              >
                <FormField
                  control={form.control}
                  name='route'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>路由</Trans>
                      </FormLabel>
                      <FormControl>
                        <SelectDropdown
                          defaultValue={field.value ?? '/'}
                          onValueChange={field.onChange}
                          placeholder={t`请选择路由`}
                          className='col-span-8 w-full'
                          items={[
                            ...Array.from(new Set(routes), (route) => ({
                              label: route,
                              value: route,
                            })),
                            ...(currentRow?.route &&
                            !routes.includes(currentRow.route)
                              ? [
                                  {
                                    label: currentRow.route,
                                    value: currentRow.route,
                                  },
                                ]
                              : []),
                          ].sort((a, b) => a.label.localeCompare(b.label))}
                        />
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='icon'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>图标</Trans>
                      </FormLabel>
                      <FormControl>
                        <IconSelect
                          className='col-span-8 w-full'
                          placeholder={t`请选择图标`}
                          value={toCamelCase(field.value)}
                          onChange={(value) => field.onChange(value)}
                        />
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
              </TabsContent>
              <TabsContent
                value={MenuType.button.toString()}
                className='flex flex-col gap-2'
              >
                <FormField
                  control={form.control}
                  name='permissionAction'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>权限动作</Trans>
                      </FormLabel>
                      <FormControl>
                        <SelectDropdown
                          defaultValue={field.value}
                          isControlled
                          onValueChange={field.onChange}
                          placeholder={t`请选择权限动作`}
                          className='col-span-8 w-full'
                          items={permissionActions.map((action) => ({
                            label: getMenuPermissionActionLabel(action),
                            value: action,
                          }))}
                        />
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    <Trans>权限标识</Trans>
                  </FormLabel>
                  <Input
                    readOnly
                    disabled
                    value={generatedPermission}
                    placeholder={t`选择父级菜单和权限动作后自动生成`}
                    className='col-span-8'
                  />
                </FormItem>
                <FormField
                  control={form.control}
                  name='apiIds'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>关联接口</Trans>
                      </FormLabel>
                      <FormControl>
                        <StatusRenderer
                          className='col-span-8 h-50'
                          isLoading={isLoadingApis || isLoadingMenu}
                        >
                          <TransferList
                            className='col-span-8'
                            value={field.value || []}
                            onChange={field.onChange}
                            renderLabel={(item) => (
                              <Badge
                                variant='outline'
                                className={cn(callMethodTypes.get(item.method))}
                              >
                                <span>{item.method}</span>
                                <span className='mx-1'>|</span>
                                <span>{item.description}</span>
                              </Badge>
                            )}
                            data={apis
                              .filter((item) => item.group === 'SYSTEM')
                              .map((item) => ({
                                ...item,
                                key: item.id,
                                label: item.description || item.path,
                                selected: field.value?.includes(item.id),
                              }))}
                          />
                        </StatusRenderer>
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
              </TabsContent>
              <FormField
                control={form.control}
                name='visible'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>是否可见</Trans>
                    </FormLabel>
                    <FormControl>
                      <Switch
                        disabled={isPending}
                        checked={field.value || false}
                        onCheckedChange={field.onChange}
                      >
                        <SwitchThumb />
                      </Switch>
                    </FormControl>
                    <FormMessage className='col-span-8 col-start-3' />
                  </FormItem>
                )}
              />
              {form.getValues('type') !== MenuType.button && (
                <FormField
                  control={form.control}
                  name='sort'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>排序</Trans>
                      </FormLabel>
                      <FormControl>
                        <Input
                          disabled={isPending}
                          placeholder={t`请输入排序值`}
                          className='col-span-8'
                          type='number'
                          min={0}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
              )}
            </Tabs>
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
        <ConfirmDialog
          open={!!pendingPlatformValues}
          onOpenChange={(state) => {
            if (!state) setPendingPlatformValues(null)
          }}
          destructive
          isLoading={isPending}
          title={<Trans>确认设为平台专属菜单？</Trans>}
          desc={
            <Trans>
              该菜单及其下的权限将只对平台超级管理员可见，所有业务租户都拿不到，
              已经授予租户角色的相关权限也会失效。此操作会记入平台审计日志。
            </Trans>
          }
          confirmText={<Trans>确认设为平台专属</Trans>}
          handleConfirm={() => {
            const values = pendingPlatformValues
            setPendingPlatformValues(null)
            if (values) submitValues(values)
          }}
        />
      </DialogContent>
    </Dialog>
  )
}
