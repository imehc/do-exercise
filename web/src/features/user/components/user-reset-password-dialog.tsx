import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconArrowBackUp, IconLoader3 } from '@tabler/icons-react'
import { toast } from 'sonner'
import {
  ResetUserPasswordRequest,
  SystemUserApi,
  SysUser,
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
import { PasswordInput } from '~/components/password-input'
import {
  ActionResetPasswordSysUserFormValues,
  resetPasswordSchema,
} from '../schemas/action-schema'

interface Props {
  currentRow: SysUser
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function UserResetPasswordDialog({
  currentRow,
  open,
  onOpenChange,
}: Props) {
  const sysUserApi = useApi(SystemUserApi)

  const form = useForm<ActionResetPasswordSysUserFormValues>({
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
    resolver: zodResolver(resetPasswordSchema),
  })

  const { isPending: isPendingResetPassword, mutate: saveResetPassword } =
    useMutation({
      mutationFn: (values: ResetUserPasswordRequest) =>
        sysUserApi.resetUserPassword(values),
      onSuccess: () => {
        toast.success('重置密码成功')
        form.reset()
        onOpenChange(false, true)
      },
    })

  const onSubmit = ({
    confirmPassword,
    ...values
  }: ActionResetPasswordSysUserFormValues) => {
    saveResetPassword({
      id: currentRow.id,
      resetSysUserPassword: {
        password: values.password,
      },
    })
  }

  const isPending = isPendingResetPassword
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
          <DialogTitle>
            <span>重置用户</span>
            <span className='mx-1 text-gray-500 italic'>
              {currentRow?.username}
            </span>
            <span>密码</span>
          </DialogTitle>
          <DialogDescription>设置新密码完成后点击保存。</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
              <FormLabel className='col-span-2 text-right'>用户名</FormLabel>
              <FormControl>
                <Input
                  placeholder='请输入用户名'
                  className='col-span-8'
                  autoComplete='off'
                  defaultValue={currentRow?.username}
                  readOnly
                  disabled
                />
              </FormControl>
              <FormMessage className='col-span-8 col-start-3' />
            </FormItem>
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-2 text-right'>密码</FormLabel>
                  <FormControl>
                    <PasswordInput
                      placeholder='请输入密码'
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
                    确认密码
                  </FormLabel>
                  <FormControl>
                    <PasswordInput
                      placeholder='请输入确认密码'
                      className='col-span-8'
                      autoComplete='off'
                      {...field}
                      disabled={!isPasswordTouched || isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-8 col-start-3' />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <DialogFooter>
          <Button type='submit' form='menu-form' disabled={isPending}>
            {isPending ? (
              <>
                <IconLoader3 className='animate-spin' />
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

export function ResetPassword({ currentRow }: Pick<Props, 'currentRow'>) {
  const { setOpen, setCurrentRow } = useFormDialog()
  return (
    <Button
      variant='outline'
      size='icon'
      onClick={() => {
        setOpen('reset')
        setCurrentRow(currentRow)
      }}
    >
      <IconArrowBackUp />
    </Button>
  )
}
