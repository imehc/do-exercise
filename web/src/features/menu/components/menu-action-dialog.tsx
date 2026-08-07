import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
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
import { useApi } from '~/hooks/use-api'
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
  ActionSysMenuWithButtonFormValues,
  ActionSysMenuWithDirectoryFormValues,
  ActionSysMenuWithMenuFormValues,
  getActionSysMenuSchema,
} from '../schemas/action-schema'

function toCamelCase(str?: string) {
  return str
    ?.split('-')
    .map((word, index) => {
      if (index === 0) return word.charAt(0).toUpperCase() + word.slice(1)
      return word.charAt(0).toUpperCase() + word.slice(1)
    })
    .join('')
}

interface Props {
  treeData: SysMenuTree[]
  currentRow?: SysMenuTree
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

const modules = Object.keys(import.meta.glob('~/features/**/index.tsx')).map(
  (item) => item.replace('/src/features', '')
)

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

  const isEdit = useMemo(() => !!currentRow?.id, [currentRow])
  const form = useForm<ActionSysMenuFormValues>({
    resolver: zodResolver(getActionSysMenuSchema()),
  })

  const { data: apis = [], isLoading: isLoadingApis } = useQuery({
    queryKey: ['findAllApis'],
    queryFn: () => sysApi.findAllApis(),
    enabled:
      form.watch('type') === MenuType.button &&
      (openType === 'add' || openType === 'edit' || openType === 'add-child'),
  })

  const { data: menu, isLoading: isLoadingMenu } = useQuery({
    queryKey: ['findMenu', currentRow?.id],
    queryFn: () => sysMenuApi.findMenu({ id: currentRow?.id as number }),
    enabled:
      form.watch('type') === MenuType.button &&
      isEdit &&
      (openType === 'add' || openType === 'edit' || openType === 'add-child'),
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
            visible: currentRow.visible,
          })
          break
        case MenuType.menu:
          form.reset({
            type: MenuType.menu,
            sort: currentRow.sort,
            parentId: currentRow.parentId,
            name: currentRow.name,
            icon: currentRow.icon,
            route: currentRow.route,
            component: currentRow.component,
            visible: currentRow.visible,
          })
          break
        case MenuType.button:
          form.reset({
            parentId: currentRow.parentId,
            name: currentRow.name,
            type: MenuType.button,
            permission: currentRow.permission,
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
    currentRow?.component,
    currentRow?.icon,
    currentRow?.name,
    currentRow?.parentId,
    currentRow?.permission,
    currentRow?.route,
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
          sort: values.sort,
          visible: Boolean(values.visible),
        } satisfies ActionSysMenuWithDirectoryFormValues
      case MenuType.menu:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          icon: values.icon,
          route: values.route,
          component: values.component,
          sort: values.sort,
          visible: Boolean(values.visible),
        } satisfies ActionSysMenuWithMenuFormValues
      case MenuType.button:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          permission: values.permission,
          apiIds: values.apiIds,
          visible: Boolean(values.visible),
        } satisfies ActionSysMenuWithButtonFormValues
    }
  }

  const onSubmit = (values: ActionSysMenuFormValues) => {
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
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <Tabs
              value={form.watch('type')?.toString()}
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
                              (item) => item.name,
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
                  name='component'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>组件</Trans>
                      </FormLabel>
                      <FormControl>
                        <SelectDropdown
                          defaultValue={field.value ?? ''}
                          onValueChange={field.onChange}
                          placeholder={t`请选择组件路径`}
                          className='col-span-8 w-full'
                          items={[
                            ...modules.map((item) => ({
                              label: item,
                              value: item,
                            })),
                            ...(currentRow?.component &&
                            !modules.includes(currentRow.component)
                              ? [
                                  {
                                    label: currentRow.component,
                                    value: currentRow.component,
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
                  name='permission'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>权限标识</Trans>
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t`请输入权限标识`}
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
                            data={apis.map((item) => ({
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
      </DialogContent>
    </Dialog>
  )
}
