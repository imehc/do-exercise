import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  BindEmailRequest,
  GetBindEmailCodeRequest,
  GetRebindEmailCodeRequest,
  RebindEmailRequest,
  UserApi,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
import { useCountdown } from '~/hooks/use-count-down'
import { useUserProfile } from '~/hooks/use-user'
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
import { getEmailSchema, EmailSchemaFormValues } from './schemas/action-schema'

export default function EmailForm() {
  const { isCounting, count, start } = useCountdown()

  const {
    data: userProfile,
    isLoading: userProfileIsLoading,
    refetch,
  } = useUserProfile()
  const form = useChan(
    useForm<EmailSchemaFormValues>({
      resolver: resolver(
        getEmailSchema().check(({ issues, value }) => {
          if (userProfile?.email === value.email.trim()) {
            issues.push({
              code: 'custom',
              message: t`邮箱不能与当前邮箱相同`,
              path: ['email'],
              input: value,
            })
          }
        })
      ),
      mode: 'onChange',
    })
  )

  const userApi = useApi(UserApi)
  const { mutate: bindEmail, isPending: bindEmailIspending } = useMutation({
    mutationFn: (value: BindEmailRequest) => userApi.bindEmail(value),
    onSuccess: () => {
      toast.success(t`绑定成功`)
      form.reset({
        email: '',
        code: '',
      })
      refetch()
    },
  })

  const { mutate: rebindEmail, isPending: rebindEmailIspending } = useMutation({
    mutationFn: (value: RebindEmailRequest) => userApi.rebindEmail(value),
    onSuccess: () => {
      toast.success(t`更换成功`)
      form.reset({
        email: '',
        code: '',
      })
      refetch()
    },
  })

  const { mutate: getBindEmailCode, isPending: bindEmailCodeIspending } =
    useMutation({
      mutationFn: (value: GetBindEmailCodeRequest) =>
        userApi.getBindEmailCode(value),
      onSuccess: () => {
        toast.success(t`发送成功`)
        start()
      },
    })

  const { mutate: getRebindEmailCode, isPending: rebindEmailCodeIspending } =
    useMutation({
      mutationFn: (value: GetRebindEmailCodeRequest) =>
        userApi.getRebindEmailCode(value),
      onSuccess: () => {
        toast.success(t`发送成功`)
        start()
      },
    })

  const hasBind = !userProfile?.email

  const isPending =
    userProfileIsLoading || bindEmailIspending || rebindEmailIspending

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) =>
          hasBind
            ? bindEmail({ bindEmail: data })
            : rebindEmail({ bindEmail: data })
        )}
        className='space-y-8'
      >
        <FormItem>
          <FormLabel>
            <Trans>用户名</Trans>
          </FormLabel>
          <FormControl>
            <Input
              placeholder='username'
              defaultValue={userProfile?.username}
              readOnly
              disabled
            />
          </FormControl>
          <FormDescription>
            <Trans>这是您的公开显示名称。</Trans>
          </FormDescription>
          <FormMessage />
        </FormItem>

        {userProfile?.email && (
          <FormItem>
            <FormLabel>
              <Trans>邮箱</Trans>
            </FormLabel>
            <FormControl>
              <Input
                placeholder={t`请输入您的邮箱`}
                defaultValue={userProfile.email}
                readOnly
                disabled
              />
            </FormControl>
            <FormDescription>
              <Trans>这是您的邮箱。</Trans>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}

        <FormField
          control={form.control}
          name='email'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                {hasBind ? <Trans>绑定邮箱</Trans> : <Trans>换绑邮箱</Trans>}
              </FormLabel>
              <div className='flex justify-between gap-x-2'>
                <Input placeholder={t`请输入您的邮箱`} {...field} />
                <Button
                  type='button'
                  className='w-1/3'
                  variant='outline'
                  disabled={
                    !form.watch('email') ||
                    form.watch('email')?.trim() === userProfile?.email ||
                    isCounting ||
                    bindEmailCodeIspending ||
                    rebindEmailCodeIspending
                  }
                  onClick={() => {
                    const email = form.watch('email')
                    if (!email) return
                    if (hasBind) {
                      getBindEmailCode({ email })
                      return
                    }
                    getRebindEmailCode({ email })
                  }}
                >
                  {isCounting ? count : t`获取验证码`}
                </Button>
              </div>
              <FormDescription>
                <Trans>这是您绑定的新电子邮件地址。</Trans>
              </FormDescription>
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
                  placeholder={t`请输入邮箱验证码`}
                  {...field}
                  disabled={bindEmailIspending || rebindEmailIspending}
                />
              </FormControl>
              <FormDescription>
                <Trans>请输入您在电子邮件中收到的验证码。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {bindEmailIspending || rebindEmailIspending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>
                {hasBind ? (
                  <>
                    <Trans>绑定中</Trans>...
                  </>
                ) : (
                  <>
                    <Trans>换绑中</Trans>...
                  </>
                )}
              </span>
            </>
          ) : (
            <span>{hasBind ? <Trans>绑定</Trans> : <Trans>换绑</Trans>}</span>
          )}
        </Button>
      </form>
    </Form>
  )
}
