import { HTMLAttributes, useEffect } from 'react'
import { addSeconds } from 'date-fns'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { keepPreviousData, useMutation, useQuery } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { IconLoader3 } from '@tabler/icons-react'
import { useSetAtom } from 'jotai'
import { JSEncrypt } from 'jsencrypt'
import { originTokenAtom } from '~/atoms'
import { AuthApi, LoginRequest } from '~/do-exercise-api'
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
import { Input } from '~/components/ui/input'
import { PasswordInput } from '~/components/password-input'
import {
  SignInActionFormValues,
  signInActionSchema,
} from '../schemas/action-schema'

type UserAuthFormProps = HTMLAttributes<HTMLFormElement>

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const authApi = useApi(AuthApi)
  const setToken = useSetAtom(originTokenAtom)
  const navigate = useNavigate()

  const form = useForm<SignInActionFormValues>({
    resolver: zodResolver(signInActionSchema),
    defaultValues: {
      username: '',
      password: '',
      captchaId: '',
      captcha: '',
      publicKey: '',
    },
  })

  const {
    data: publicKeyData,
    isLoading: publicKeyDataIsLoading,
    refetch: refetchPublicKey,
  } = useQuery({
    queryKey: ['getPublicKey'],
    queryFn: () => authApi.getPublicKey(),
    retry: false,
    placeholderData: keepPreviousData,
    refetchInterval: 5 * 60 * 1000,
  })

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

  const { mutate: login, isPending: loginIsPending } = useMutation({
    mutationFn: (value: LoginRequest) => authApi.login(value),
    onSuccess: (data) => {
      setToken({
        ...data,
        expireTime: addSeconds(new Date(), data.expireTime).getTime(),
        refreshExpireTime: addSeconds(
          new Date(),
          data.refreshExpireTime
        ).getTime(),
      })
      navigate({
        to: '/',
      })
    },
    onError: () => {
      refetchPublicKey()
      refetchCaptcha()
    },
  })

  function onSubmit(data: SignInActionFormValues) {
    if (!publicKeyData) return
    const encrypt = new JSEncrypt()
    // 注意：Go 返回的可能是 PEM 格式的 base64，需要先解码
    const publicKey = atob(publicKeyData.publicKey) // 浏览器环境用 atob，Node 用 Buffer
    encrypt.setPublicKey(publicKey)
    const password = encrypt.encrypt(data.password)
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
              <FormLabel>用户名</FormLabel>
              <FormControl>
                <Input placeholder='用户名' {...field} />
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
              <FormLabel>密码</FormLabel>
              <FormControl>
                <PasswordInput placeholder='********' {...field} />
              </FormControl>
              <FormMessage />
              <Link
                to='/forgot-password'
                className='text-muted-foreground absolute -top-0.5 right-0 text-sm font-medium hover:opacity-75'
              >
                忘记密码?
              </Link>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='captcha'
          render={({ field }) => (
            <FormItem>
              <FormLabel>验证码</FormLabel>
              <FormControl>
                <div className='flex w-full items-center justify-between gap-x-4'>
                  <Input
                    disabled={isPending}
                    placeholder='Capthca'
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
              <span>登录中...</span>
            </>
          ) : (
            <span>登录</span>
          )}
        </Button>
        {/* 
        <div className='relative my-2'>
          <div className='absolute inset-0 flex items-center'>
            <span className='w-full border-t' />
          </div>
          <div className='relative flex justify-center text-xs uppercase'>
            <span className='bg-background text-muted-foreground px-2'>
              Or continue with
            </span>
          </div>
        </div>

        <div className='grid grid-cols-2 gap-2'>
          <Button variant='outline' type='button' disabled={isLoading}>
            <IconBrandGithub className='h-4 w-4' /> GitHub
          </Button>
          <Button variant='outline' type='button' disabled={isLoading}>
            <IconBrandFacebook className='h-4 w-4' /> Facebook
          </Button>
        </div> */}
      </form>
    </Form>
  )
}
