import { HTMLAttributes, useEffect } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import {
  AuthApi,
  ForgetPasswordRequest,
  GetForgetPasswordCodeRequest,
} from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { encryptPassword } from '~/utils/encrypt'
import { useApi } from '~/hooks/use-api'
import { useCountdown } from '~/hooks/use-count-down'
import { usePublicKey } from '~/hooks/use-public-key'
import { Button } from '~/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { PasswordInput } from '~/components/password-input'
import { passwordRule } from '~/features/user/schemas/action-schema'

type ForgotFormProps = HTMLAttributes<HTMLFormElement>

const formSchema = z.object({
  email: z
    .string()
    .min(1, { message: 'Please enter your email' })
    .email({ message: 'Invalid email address' }),
  password: passwordRule,
  confirmPassword: passwordRule,
  code: z
    .string({ required_error: '请输入验证码' })
    .min(1, { message: '请输入验证码' }),
  publicKey: z.string({ required_error: '公钥不能为空' }),
})

export function ForgotPasswordForm({ className, ...props }: ForgotFormProps) {
  const navigate = useNavigate()
  const { isCounting, count, start } = useCountdown({ seconds: 10 })
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { email: '' },
  })

  const {
    data: publicKeyData,
    isLoading: publicKeyDataIsLoading,
    refetch: refetchPublicKey,
  } = usePublicKey()

  useEffect(() => {
    if (!publicKeyData?.publicKey) return
    form.setValue('publicKey', publicKeyData.publicKey)
  }, [form, publicKeyData?.publicKey])

  const authApi = useApi(AuthApi)
  const {
    mutate: getForgetPasswordCode,
    isPending: forgetPasswordCodeIsPending,
  } = useMutation({
    mutationFn: (value: GetForgetPasswordCodeRequest) =>
      authApi.getForgetPasswordCode(value),
    onSuccess: () => {
      toast.success(`发送成功`)
      start()
    },
  })

  const { mutate: forgetPassword, isPending: forgetPasswordIsPending } =
    useMutation({
      mutationFn: (value: ForgetPasswordRequest) =>
        authApi.forgetPassword(value),
      onSuccess: () => {
        toast.success(`重置密码成功`)
        navigate({
          to: '/sign-in',
        })
      },
      onError: () => {
        refetchPublicKey()
      },
    })

  function onSubmit({ confirmPassword, ...data }: z.infer<typeof formSchema>) {
    if (!publicKeyData?.publicKey) return
    const password = encryptPassword(data.password, publicKeyData.publicKey)
    if (!password) return
    forgetPassword({
      forgetPassword: {
        ...data,
        password,
      },
    })
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-2', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='email'
          render={({ field }) => (
            <FormItem>
              <FormLabel> Email</FormLabel>
              <div className='flex justify-between gap-x-2 max-sm:flex-col max-sm:gap-y-2'>
                <Input placeholder='email' {...field} />
                <Button
                  type='button'
                  variant='outline'
                  className='w-1/3 max-sm:w-full'
                  disabled={
                    !form.watch('email') ||
                    isCounting ||
                    forgetPasswordCodeIsPending ||
                    forgetPasswordIsPending ||
                    publicKeyDataIsLoading
                  }
                  onClick={() => {
                    const email = form.watch('email')
                    if (!email) return
                    getForgetPasswordCode({ email })
                  }}
                >
                  {isCounting ? count : '获取验证码'}
                </Button>
              </div>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='code'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Code</FormLabel>
              <FormControl>
                <Input
                  placeholder='code'
                  {...field}
                  disabled={forgetPasswordIsPending}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='password'
                  {...field}
                  disabled={forgetPasswordIsPending}
                />
              </FormControl>
              <FormDescription>This is your new password.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='confirmPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='confirm password'
                  {...field}
                  disabled={forgetPasswordIsPending}
                />
              </FormControl>
              <FormDescription>Enter your new password again.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button className='mt-2' disabled={forgetPasswordIsPending}>
          Continue
        </Button>
      </form>
    </Form>
  )
}
