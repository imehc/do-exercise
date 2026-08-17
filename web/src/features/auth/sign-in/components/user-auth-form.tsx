import { HTMLAttributes, useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { IconLoader3, IconMail } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { useSetAtom } from 'jotai'
import { originTokenAtom } from '~/atoms'
import {
  AuthApi,
  LoginRequest,
  LoginResult,
  TenantOption,
} from '~/do-exercise-api'
import { cn } from '~/lib/utils'
import { encryptPassword } from '~/utils/encrypt'
import { applyToken } from '~/lib/token'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
import { usePublicKey } from '~/hooks/use-public-key'
import { Button } from '~/components/ui/button'
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
import { TenantSelectDialog } from './tenant-select-dialog'
import {
  getSignInActionSchema,
  SignInActionFormValues,
} from '../schemas/action-schema'

type UserAuthFormProps = HTMLAttributes<HTMLFormElement>

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const authApi = useApi(AuthApi)
  const setToken = useSetAtom(originTokenAtom)

  const form = useChan(
    useForm<SignInActionFormValues>({
      resolver: zodResolver(getSignInActionSchema()),
      defaultValues: {
        username: '',
        password: '',
        captchaId: '',
        captcha: '',
        publicKey: '',
        tenantId: '',
      },
    })
  )

  const {
    data: publicKeyData,
    isLoading: publicKeyDataIsLoading,
    refetch: refetchPublicKey,
  } = usePublicKey()

  const {
    data: captchaData,
    isLoading: captchaDataIsLoading,
    refetch: refetchCaptcha,
  } = useQuery({
    queryKey: ['getCaptcha'],
    queryFn: () => authApi.getCaptcha(),
    retry: false,
    refetchInterval: 60 * 1000,
  })

  useEffect(() => {
    if (publicKeyData) {
      form.setValue('publicKey', publicKeyData.publicKey)
    }
    if (captchaData) {
      form.setValue('captchaId', captchaData.captchaId)
    }
  }, [captchaData, form, publicKeyData])

  // 多租户登录：账号归属多个启用租户时，登录接口返回 requires_tenant_selection，
  // 暂存会话信息并弹出租户选择框，选定后经 select_tenant 完成登录。
  const [tenantSelection, setTenantSelection] = useState<{
    loginSessionId: string
    tenants: TenantOption[]
  } | null>(null)

  const applyLoginResult = (data: LoginResult) => {
    setToken(applyToken(data))
    // 默认口令强制轮换——首次登录成功后强制跳转到改密页
    window.location.href = data.mustChangePassword
      ? '/settings/password'
      : '/'
  }

  const { mutate: login, isPending: loginIsPending } = useMutation({
    mutationFn: (value: LoginRequest) => authApi.login(value),
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
    onError: () => {
      refetchPublicKey()
      refetchCaptcha()
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
      onError: () => {
        setTenantSelection(null)
        refetchPublicKey()
        refetchCaptcha()
      },
    })

  function onSubmit(data: SignInActionFormValues) {
    if (!publicKeyData) return
    const password = encryptPassword(data.password, publicKeyData.publicKey)
    if (!password) {
      return
    }

    login({
      login: {
        ...data,
        password: password,
        username: data.username,
      },
    })
  }

  const isPending =
    publicKeyDataIsLoading || captchaDataIsLoading || loginIsPending

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className={cn('grid gap-3', className)}
          {...props}
        >
          <FormField
            control={form.control}
            name='username'
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Trans>用户名</Trans>
                </FormLabel>
                <FormControl>
                  <Input placeholder={t`请输入用户名`} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem className='relative'>
              <FormLabel>
                <Trans>密码</Trans>
              </FormLabel>
              <FormControl>
                <PasswordInput placeholder={t`请输入密码`} {...field} />
              </FormControl>
              <FormMessage />
              <Link
                to='/forgot-password'
                className='text-muted-foreground absolute -top-0.5 right-0 text-sm font-medium hover:opacity-75'
              >
                <Trans>忘记密码?</Trans>
              </Link>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='captcha'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>验证码</Trans>
              </FormLabel>
              <FormControl>
                <div className='flex w-full items-center justify-between gap-x-4'>
                  <Input
                    disabled={isPending}
                    placeholder={t`验证码`}
                    {...field}
                  />
                  <div className='relative aspect-[3/1] h-9 overflow-hidden rounded-md border border-solid border-[var(--input)]'>
                    {!!captchaData && (
                      <img
                        className={cn(
                          'box-border h-full w-full px-1',
                          !isPending
                            ? 'pointer-events-auto cursor-pointer'
                            : 'pointer-events-none cursor-none'
                        )}
                        src={captchaData.picPath}
                        alt='captcha'
                        onClick={() => refetchCaptcha()}
                      />
                    )}
                  </div>
                </div>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isPending}>
          {loginIsPending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>
                <Trans>登录</Trans>...
              </span>
            </>
          ) : (
            <span>
              <Trans>登录</Trans>
            </span>
          )}
        </Button>

        <div className='relative my-2'>
          <div className='absolute inset-0 flex items-center'>
            <span className='w-full border-t' />
          </div>
          <div className='relative flex justify-center text-xs uppercase'>
            <span className='bg-background text-muted-foreground px-2'>
              <Trans>或</Trans>
            </span>
          </div>
        </div>

        <div className='grid grid-cols-1 gap-2'>
          <Link to='/email-sign-in' disabled={isPending}>
            <Button
              variant='outline'
              type='button'
              disabled={isPending}
              className='w-full'
            >
              <IconMail className='h-4 w-4' />
              <Trans>邮箱</Trans>
              <Trans>登录</Trans>
            </Button>
          </Link>
          {/* <Button variant='outline' type='button' disabled={isLoading}>
            <IconBrandFacebook className='h-4 w-4' /> Facebook
          </Button> */}
        </div>
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
