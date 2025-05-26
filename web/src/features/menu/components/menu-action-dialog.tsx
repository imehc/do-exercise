import { useEffect, useMemo } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { SwitchThumb } from '@radix-ui/react-switch'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Loader2 } from 'lucide-react'
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
import { SelectDropdown } from '~/components/select-dropdown'
import { callMenuMapping } from '../data/data'
import {
  ActionSysMenuFormValues,
  actionSysMenuSchema,
  ActionSysMenuWithButton,
  ActionSysMenuWithDirectory,
  ActionSysMenuWithMenu,
} from '../schemas/action-schema'

function toCamelCase(str: string) {
  return str
    .split('-')
    .map((word, index) => {
      if (index === 0) return word.charAt(0).toUpperCase() + word.slice(1)
      return word.charAt(0).toUpperCase() + word.slice(1)
    })
    .join('')
}

interface Props {
  currentRow?: SysMenuTree
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

const modules = Object.keys(import.meta.glob('~/features/**/index.tsx')).map(
  (item) => item.replace('/src/features', '')
)

export function MenuActionDialog({ currentRow, open, onOpenChange }: Props) {
  const sysMenuApi = useApi(SystemMenuApi)
  const sysApi = useApi(SystemApiApi)

  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useForm<ActionSysMenuFormValues>({
    resolver: zodResolver(actionSysMenuSchema),
  })

  const { data = [], isLoading } = useQuery({
    queryKey: ['findMenuTree'],
    queryFn: () => sysMenuApi.findMenuTree(),
  })
  const { data: apis = [], isLoading: isLoadingApis } = useQuery({
    queryKey: ['findAllApis'],
    queryFn: () => sysApi.findAllApis(),
    enabled: form.watch('type') === MenuType.button,
  })
  const { data: menu, isLoading: isLoadingMenu } = useQuery({
    queryKey: ['findMenu', currentRow?.id],
    queryFn: () => sysMenuApi.findMenu({ id: currentRow?.id as number }),
    enabled: form.watch('type') === MenuType.button && isEdit,
  })

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateMenuRequest) => sysMenuApi.createMenu(values),
    onSuccess: () => {
      toast.success('创建成功')
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateMenuRequest) => sysMenuApi.updateMenu(values),
    onSuccess: () => {
      toast.success('更新成功')
      form.reset()
      onOpenChange(false, true)
    },
  })

  useEffect(() => {
    if (isEdit) {
      switch (currentRow?.type) {
        case MenuType.directory:
          form.reset({
            type: MenuType.directory,
            sort: currentRow.sort,
            parentId: currentRow.parentId,
            name: currentRow.name,
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
              (menu as SysMenuWithButton).apis?.map((api) => api.id) ?? [],
          })
      }
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
    menu,
  ])

  const handleFormat = (values: ActionSysMenuFormValues) => {
    switch (values.type) {
      case MenuType.directory:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          sort: values.sort,
        } satisfies z.infer<typeof ActionSysMenuWithDirectory>
      case MenuType.menu:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          icon: values.icon,
          route: values.route,
          component: values.component,
          visible: Boolean(values.visible),
          sort: values.sort,
        } satisfies z.infer<typeof ActionSysMenuWithMenu>
      case MenuType.button:
        return {
          type: values.type,
          parentId: values.parentId,
          name: values.name,
          permission: values.permission,
          apiIds: values.apiIds,
        } satisfies z.infer<typeof ActionSysMenuWithButton>
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
          <DialogTitle>{isEdit ? '修改菜单' : '创建菜单'}</DialogTitle>
          <DialogDescription>
            {isEdit ? '更新菜单相关信息。' : '创建菜单相关信息。'}
            完成后点击保存。
          </DialogDescription>
        </DialogHeader>
        <StatusRenderer isLoading={isLoading}>
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
                  {Array.from(callMenuMapping.entries())
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
                        父级菜单
                      </FormLabel>
                      <FormControl>
                        <TreeSelect
                          className='col-span-8'
                          placeholder='请选择父级菜单'
                          data={[
                            {
                              name: '根节点',
                              value: '0',
                              children: transformData(
                                data,
                                (item) => item.name,
                                (item) => item.id.toString()
                              ),
                            },
                          ]}
                          value={
                            typeof field.value === 'number'
                              ? [field.value.toString()]
                              : []
                          }
                          onChange={(value) =>
                            field.onChange(value[0] ? Number(value[0]) : '')
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
                        菜单名称
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder='请输入菜单名称'
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
                          路由
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder='请输入路由地址'
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
                    name='component'
                    render={({ field }) => (
                      <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                        <FormLabel className='col-span-2 text-right'>
                          组件
                        </FormLabel>
                        <FormControl>
                          <SelectDropdown
                            defaultValue={field.value ?? ''}
                            onValueChange={field.onChange}
                            placeholder='请选择组件路径'
                            className='col-span-8 w-full'
                            items={[
                              ...modules.map((item) => ({
                                label: item,
                                value: item,
                              })),
                              modules.some(
                                (item) => item === currentRow?.component
                              )
                                ? []
                                : {
                                    label: currentRow?.component as string,
                                    value: currentRow?.component as string,
                                  },
                            ]
                              .flat()
                              .sort((a, b) => a.label.localeCompare(b.label))}
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
                          图标
                        </FormLabel>
                        <FormControl>
                          <IconSelect
                            className='col-span-8 w-full'
                            placeholder='请选择图标'
                            value={toCamelCase(field.value)}
                            onChange={(value) => field.onChange(value)}
                          />
                        </FormControl>
                        <FormMessage className='col-span-8 col-start-3' />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name='visible'
                    render={({ field }) => (
                      <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                        <FormLabel className='col-span-2 text-right'>
                          是否隐藏
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
                          权限标识
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder='请输入权限标识'
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
                          关联API
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
                              data={apis.map((item) => ({
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
                {form.getValues('type') !== MenuType.button && (
                  <FormField
                    control={form.control}
                    name='sort'
                    render={({ field }) => (
                      <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                        <FormLabel className='col-span-2 text-right'>
                          排序
                        </FormLabel>
                        <FormControl>
                          <Input
                            disabled={isPending}
                            placeholder='请输入排序值'
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
        </StatusRenderer>
        <DialogFooter>
          <Button type='submit' form='menu-form' disabled={isPending}>
            {isPending ? (
              <>
                <Loader2 className='animate-spin' />
                <span>保存中...</span>
              </>
            ) : (
              <span>保存</span>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
