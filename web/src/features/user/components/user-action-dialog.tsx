import { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import {
  CreateSysUser,
  CreateUserRequest,
  SystemRoleApi,
  SystemUserApi,
  SysUser,
  UpdateUserRequest,
} from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
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
import { AvatarUpload, MultiSelect, StatusRenderer } from '~/components/other'
import { PasswordInput } from '~/components/password-input'
import { ActionSysUserFormValues, schema } from '../schemas/action-schema'

interface Props {
  currentRow?: SysUser
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function UserActionDialog({ currentRow, open, onOpenChange }: Props) {
  const { open: openType } = useFormDialog()
  const sysUserApi = useApi(SystemUserApi)
  const sysRoleApi = useApi(SystemRoleApi)

  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useForm<ActionSysUserFormValues>({
    defaultValues: isEdit
      ? {
          ...currentRow,
          isEdit,
          roleIds: currentRow?.roles?.map((item) => item.id) ?? [],
        }
      : {
          isEdit,
          username: '',
          nickname: '',
          email: '',
          avatar: '',
          roleIds: [],
        },
    resolver: zodResolver(schema),
  })

  const { data = [], isLoading } = useQuery({
    queryKey: ['findAllRoles'],
    queryFn: () => sysRoleApi.findAllRoles(),
    enabled: openType === 'add' || openType === 'edit',
  })

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateUserRequest) => sysUserApi.createUser(values),
    onSuccess: () => {
      toast.success('创建成功')
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateUserRequest) => sysUserApi.updateUser(values),
    onSuccess: () => {
      toast.success('更新成功')
      form.reset()
      onOpenChange(false, true)
    },
  })

  const onSubmit = ({
    isEdit,
    confirmPassword,
    ...values
  }: ActionSysUserFormValues) => {
    if (isEdit) {
      saveUpdate({
        updateSysUser: {
          nickname: values.nickname,
          avatar: values.avatar,
          roleIds: values.roleIds,
        },
        id: currentRow?.id as number,
      })
      return
    }
    saveCreate({
      createSysUser: values as CreateSysUser,
    })
  }

  const isPending = isPendingCreate || isPendingUpdate
  const isPasswordTouched = !!form.formState.dirtyFields.password

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
          <DialogTitle>{isEdit ? '修改用户' : '创建用户'}</DialogTitle>
          <DialogDescription>
            {isEdit ? '更新用户相关信息。' : '创建用户相关信息。'}
            完成后点击保存。
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <FormField
              control={form.control}
              name='username'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    用户名
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder='请输入用户名'
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      readOnly={isEdit || isPending}
                      disabled={isEdit || isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='nickname'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>昵称</FormLabel>
                  <FormControl>
                    <Input
                      placeholder='请输入昵称'
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
              name='email'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>邮箱</FormLabel>
                  <FormControl>
                    <Input
                      placeholder='请输入邮箱'
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
              name='avatar'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>头像</FormLabel>
                  <FormControl>
                    <AvatarUpload
                      // TODO: 根据前缀判断是否需要拼接完整的图片地址
                      value={field.value}
                      onChange={field.onChange}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='roleIds'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>
                    关联角色
                  </FormLabel>
                  <FormControl>
                    <StatusRenderer
                      isLoading={isLoading}
                      className='col-span-8 h-10'
                    >
                      <MultiSelect
                        className='col-span-8'
                        modalPopover
                        options={data.map((item) => ({
                          label: item.name,
                          value: item.id,
                        }))}
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                        placeholder='请选择关联角色'
                        variant='inverted'
                      />
                    </StatusRenderer>
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
            {!isEdit && (
              <>
                <FormField
                  control={form.control}
                  name='password'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-right'>
                        密码
                      </FormLabel>
                      <FormControl>
                        <PasswordInput
                          placeholder='请输入密码'
                          className='col-span-8'
                          autoComplete='off'
                          {...field}
                          disabled={isEdit || isPending}
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
                        确认密码
                      </FormLabel>
                      <FormControl>
                        <PasswordInput
                          placeholder='请输入确认密码'
                          className='col-span-8'
                          autoComplete='off'
                          {...field}
                          disabled={!isPasswordTouched || isEdit || isPending}
                        />
                      </FormControl>
                      <FormMessage className='col-span-8 col-start-3' />
                    </FormItem>
                  )}
                />
              </>
            )}
          </form>
        </Form>
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
