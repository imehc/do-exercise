import { HTMLAttributes, useEffect, useState } from 'react'
import { z } from 'zod'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useSetAtom } from 'jotai'
import { toast } from 'sonner'
import { originTokenAtom } from '~/atoms'
import {
  AuthApi,
  LoginResult,
  LoginWithEmailRequest,
  TenantOption,
} from '~/do-exercise-api'
import { Route } from '~/routes/(auth)/otp'
import { applyToken } from '~/lib/token'
import { cn } from '~/lib/utils'
import { useApi } from '~/hooks/use-api'
import { Button } from '~/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from '~/components/ui/input-otp'
import { TenantSelectDialog } from '~/features/auth/components/tenant-select-dialog'

type OtpFormProps = HTMLAttributes<HTMLFormElement>

const getFormSchema = () =>
  z.object({
    code: z.string().min(1, { error: t`请输入邮箱验证码` }),
    email: z.email({
      error: (issue) =>
        issue.input === undefined ? t`请输入您的邮箱` : t`邮箱无效`,
    }),
  })

type FormSchemaValues = z.infer<ReturnType<typeof getFormSchema>>

export function OtpForm({ className, ...props }: OtpFormProps) {
  const setToken = useSetAtom(originTokenAtom)
  const { email } = Route.useSearch()

  const form = useForm<FormSchemaValues>({
    resolver: zodResolver(getFormSchema()),
    defaultValues: { code: '', email: '' },
  })

  useEffect(() => {
    if (!email?.trim()) {
      toast.error(t`邮箱无效`)
      return
    }
    form.setValue('email', email)
  })

  const code = useWatch({ control: form.control, name: 'code' })

  const authApi = useApi(AuthApi)

  // 同一邮箱在多个租户下各有一个账号时，验证码只能证明「邮箱是你的」，
  // 证明不了要进哪个租户。后端返回一次性登录会话与候选租户，选定后再发正式 token，
  // 与口令登录走的是同一套 select_tenant 机制。
  const [tenantSelection, setTenantSelection] = useState<{
    loginSessionId: string
    tenants: TenantOption[]
  } | null>(null)

  const applyLoginResult = (data: LoginResult) => {
    setToken(applyToken(data))
    // 默认口令强制轮换——首次登录成功后强制跳转到改密页
    window.location.href = data.mustChangePassword ? '/settings/password' : '/'
  }

  const { mutate: loginWithEmail, isPending } = useMutation({
    mutationFn: (value: LoginWithEmailRequest) => authApi.loginWithEmail(value),
    onSuccess: (data) => {
      if (data.requiresTenantSelection) {
        setTenantSelection({
          loginSessionId: data.loginSessionId ?? '',
          tenants: data.availableTenants,
        })
        return
      }
      applyLoginResult(data)
    },
  })

  const { mutate: selectTenant, isPending: selectTenantIsPending } =
    useMutation({
      mutationFn: (tenant: TenantOption) =>
        authApi.selectTenant({
          selectTenantRequest: {
            loginSessionId: tenantSelection?.loginSessionId ?? '',
            tenantId: tenant.tenantId ?? '',
          },
        }),
      onSuccess: (data) => {
        setTenantSelection(null)
        applyLoginResult(data)
      },
      onError: () => setTenantSelection(null),
    })

  function onSubmit(data: FormSchemaValues) {
    loginWithEmail({ loginWithEmail: data })
  }

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className={cn('grid gap-2', className)}
          {...props}
        >
          <FormField
            control={form.control}
            name='code'
            render={({ field }) => (
              <FormItem>
                <FormLabel className='sr-only'>
                  <Trans>邮箱验证码</Trans>
                </FormLabel>
                <FormControl>
                  <InputOTP
                    disabled={isPending}
                    maxLength={6}
                    {...field}
                    containerClassName='justify-between sm:[&>[data-slot="input-otp-group"]>div]:w-12'
                  >
                    <InputOTPGroup>
                      <InputOTPSlot index={0} />
                      <InputOTPSlot index={1} />
                    </InputOTPGroup>
                    <InputOTPSeparator />
                    <InputOTPGroup>
                      <InputOTPSlot index={2} />
                      <InputOTPSlot index={3} />
                    </InputOTPGroup>
                    <InputOTPSeparator />
                    <InputOTPGroup>
                      <InputOTPSlot index={4} />
                      <InputOTPSlot index={5} />
                    </InputOTPGroup>
                  </InputOTP>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button className='mt-2' disabled={code.length < 6 || isPending}>
            <Trans>验证</Trans>
          </Button>
        </form>
      </Form>
      <TenantSelectDialog
        open={!!tenantSelection}
        tenants={tenantSelection?.tenants ?? []}
        isPending={selectTenantIsPending}
        onSelect={selectTenant}
      />
    </>
  )
}
