import { HTMLAttributes, useEffect } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  AuthApi,
  ForgetPasswordRequest,
  GetForgetPasswordCodeRequest,
} from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { encryptPassword } from '~/utils/encrypt'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
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
import { getPasswordRule } from '~/features/user/schemas/action-schema'

type ForgotFormProps = HTMLAttributes<HTMLFormElement>

const getFormSchema = () =>
  z.object({
    email: z
      .string()
      .min(1, { message: t`请输入邮箱` })
      .email({ message: t`邮箱无效` }),
    password: getPasswordRule(),
    confirmPassword: getPasswordRule(),
    code: z
      .string({ required_error: t`请输入验证码` })
      .min(1, { message: t`请输入验证码` }),
    publicKey: z.string({ required_error: t`公钥不能为空` }),
  })

type FormSchemaValues = z.infer<ReturnType<typeof getFormSchema>>

export function ForgotPasswordForm({ className, ...props }: ForgotFormProps) {
  const navigate = useNavigate()
  const { isCounting, count, start } = useCountdown({ seconds: 10 })
  const form = useChan(
    useForm<FormSchemaValues>({
      resolver: zodResolver(getFormSchema()),
      defaultValues: { email: '' },
    })
  )

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
      toast.success(t`发送成功`)
      start()
    },
  })

  const { mutate: forgetPassword, isPending: forgetPasswordIsPending } =
    useMutation({
      mutationFn: (value: ForgetPasswordRequest) =>
        authApi.forgetPassword(value),
      onSuccess: () => {
        toast.success(t`重置密码成功`)
        navigate({
          to: '/sign-in',
        })
      },
      onError: () => {
        refetchPublicKey()
      },
    })

  function onSubmit({ confirmPassword, ...data }: FormSchemaValues) {
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
              <FormLabel>
                <Trans>邮箱</Trans>
              </FormLabel>
              <div className='flex justify-between gap-x-2 max-sm:flex-col max-sm:gap-y-2'>
                <Input placeholder={t`请输入邮箱`} {...field} />
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
                  {isCounting ? count : <Trans>发送验证码</Trans>}
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
              <FormLabel>
                <Trans>验证码</Trans>
              </FormLabel>
              <FormControl>
                <Input
                  placeholder={t`请输入验证码`}
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
              <FormLabel>
                <Trans>密码</Trans>
              </FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder={t`请输入密码`}
                  {...field}
                  disabled={forgetPasswordIsPending}
                />
              </FormControl>
              <FormDescription>
                <Trans>输入您的新密码。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='confirmPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>确认密码</Trans>
              </FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder={t`请再次输入密码`}
                  {...field}
                  disabled={forgetPasswordIsPending}
                />
              </FormControl>
              <FormDescription>
                <Trans>再次输入您的新密码。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button className='mt-2' disabled={forgetPasswordIsPending}>
          <Trans>继续</Trans>
        </Button>
      </form>
    </Form>
  )
}
