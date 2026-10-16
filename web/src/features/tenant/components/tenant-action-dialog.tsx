import { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  CreateTenant,
  SystemTenantApi,
  Tenant,
  UpdateTenant,
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { Switch } from '~/components/ui/switch'
import { DialogHeaderContent } from '~/components/other/dialog-header-content'
import { PasswordInput } from '~/components/password-input'
import { ActionTenantFormValues, getSchema } from '../schemas/action-schema'

interface Props {
  currentRow?: Tenant
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function TenantActionDialog({ currentRow, open, onOpenChange }: Props) {
  const systemTenantApi = useApi(SystemTenantApi)

  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useChan(
    useForm<ActionTenantFormValues>({
      defaultValues: isEdit
        ? {
            isEdit,
            name: currentRow?.name ?? '',
            code: currentRow?.code ?? '',
            status: currentRow?.status,
            remark: currentRow?.remark ?? '',
          }
        : {
            isEdit,
            name: '',
            code: '',
            adminMode: 'new',
            adminUserId: '',
            adminUsername: '',
            adminPassword: '',
            confirmPassword: '',
            remark: '',
          },
      resolver: zodResolver(getSchema()),
    })
  )

  // 创建租户时可选的现有用户列表（用作租户管理员）
  const { data: assignableAdmins = [] } = useQuery({
    queryKey: ['listAssignableAdmins'],
    queryFn: () => systemTenantApi.listAssignableAdmins(),
    enabled: !isEdit,
  })

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateTenant) =>
      systemTenantApi.createTenant({ createTenant: values }),
    onSuccess: () => {
      toast.success(t`创建成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateTenant) =>
      systemTenantApi.updateTenant({
        id: currentRow?.tenantId as string,
        updateTenant: values,
      }),
    onSuccess: () => {
      toast.success(t`更新成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const onSubmit = ({
    isEdit,
    confirmPassword,
    adminPassword,
    adminUsername,
    adminMode,
    adminUserId,
    code,
    ...values
  }: ActionTenantFormValues) => {
    if (isEdit) {
      saveUpdate({
        name: values.name,
        status: values.status,
        remark: values.remark,
      })
      return
    }
    if (adminMode === 'existing') {
      saveCreate({
        name: values.name,
        code: code as string,
        adminMode: 'existing',
        adminUserId: adminUserId as string,
        remark: values.remark,
      })
      return
    }
    saveCreate({
      name: values.name,
      code: code as string,
      adminMode: 'new',
      adminUsername: adminUsername as string,
      adminPassword: adminPassword as string,
      remark: values.remark,
    })
  }

  const isPending = isPendingCreate || isPendingUpdate
  const isAdminPasswordTouched = !!form.formState.dirtyFields.adminPassword
  const adminMode = form.watch('adminMode')

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
          <DialogHeaderContent isEdit={isEdit} text={<Trans>租户</Trans>} />
        </DialogHeader>
        <Form {...form}>
          <form
            id='tenant-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <FormField
              control={form.control}
              name='name'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    <Trans>租户名称</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入租户名称`}
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
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
                    <Trans>租户编码</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入租户编码`}
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      readOnly={isEdit}
                      disabled={isEdit || isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            {!isEdit && (
              <>
                <FormField
                  control={form.control}
                  name='adminMode'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        <Trans>管理员账号</Trans>
                      </FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value ?? 'new'}
                        disabled={isPending}
                      >
                        <FormControl className='col-span-8 w-full'>
                          <SelectTrigger>
                            <SelectValue placeholder={t`请选择管理员账号`} />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value='new'>
                            <Trans>新建用户</Trans>
                          </SelectItem>
                          <SelectItem value='existing'>
                            <Trans>选择现有用户</Trans>
                          </SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
                {adminMode === 'existing' ? (
                  <FormField
                    control={form.control}
                    name='adminUserId'
                    render={({ field }) => (
                      <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                        <FormLabel className='col-span-2 text-right'>
                          <Trans>选择用户</Trans>
                        </FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value || undefined}
                          disabled={isPending}
                        >
                          <FormControl className='col-span-8 w-full'>
                            <SelectTrigger>
                              <SelectValue placeholder={t`请选择现有用户`} />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {assignableAdmins.map((user) => (
                              <SelectItem key={user.id} value={user.id ?? ''}>
                                {user.username}
                                {user.nickname ? `（${user.nickname}）` : ''}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage className='col-span-8 col-start-3' />
                      </FormItem>
                    )}
                  />
                ) : (
                  <>
                    <FormField
                      control={form.control}
                      name='adminUsername'
                      render={({ field }) => (
                        <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                          <FormLabel className='col-span-2 text-right'>
                            <Trans>管理员用户名</Trans>
                          </FormLabel>
                          <FormControl>
                            <Input
                              placeholder={t`请输入管理员用户名`}
                              className='col-span-8'
                              autoComplete='off'
                              {...field}
                              disabled={isPending}
                            />
                          </FormControl>
                          <FormMessage className='col-span-8 col-start-3' />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name='adminPassword'
                      render={({ field }) => (
                        <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                          <FormLabel className='col-span-2 text-right'>
                            <Trans>管理员初始密码</Trans>
                          </FormLabel>
                          <FormControl>
                            <PasswordInput
                              placeholder={t`请输入管理员初始密码`}
                              className='col-span-8'
                              autoComplete='off'
                              {...field}
                              disabled={isPending}
                            />
                          </FormControl>
                          <FormMessage className='col-span-8 col-start-3' />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name='confirmPassword'
                      render={({ field }) => (
                        <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                          <FormLabel className='col-span-2 text-right'>
                            <Trans>确认密码</Trans>
                          </FormLabel>
                          <FormControl>
                            <PasswordInput
                              placeholder={t`请输入确认密码`}
                              className='col-span-8'
                              autoComplete='off'
                              {...field}
                              disabled={!isAdminPasswordTouched || isPending}
                            />
                          </FormControl>
                          <FormMessage className='col-span-8 col-start-3' />
                        </FormItem>
                      )}
                    />
                  </>
                )}
              </>
            )}
            {isEdit && (
              <FormField
                control={form.control}
                name='status'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>状态</Trans>
                    </FormLabel>
                    <FormControl className='col-span-8'>
                      <div className='flex items-center gap-2'>
                        <Switch
                          checked={!!field.value}
                          onCheckedChange={field.onChange}
                          disabled={isPending}
                        />
                        <span className='text-sm'>
                          {field.value ? <Trans>启用</Trans> : <Trans>停用</Trans>}
                        </span>
                      </div>
                    </FormControl>
                    <FormMessage className='col-span-8 col-start-3' />
                  </FormItem>
                )}
              />
            )}
            <FormField
              control={form.control}
              name='remark'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    <Trans>备注</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入备注`}
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <DialogFooter>
          <Button type='submit' form='tenant-form' disabled={isPending}>
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