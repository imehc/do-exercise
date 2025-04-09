import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import { Card, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import '~/animations.css'
import { useApi } from '~/hooks'
import { AuthApi, type LoginRequest } from '~/do-exercise-api'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { Loading } from '~/components'
import { loginSchema, LoginSchemaType } from './schema.'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '~/components/ui/form'
import ServerErrorPage from '../server-error'
import { DynamicIcon } from 'lucide-react/dynamic'
import JSEncrypt from 'jsencrypt'
import clsx from 'clsx'

export function LoginPage() {
  const form = useForm<LoginSchemaType>({
    resolver: resolver(loginSchema)
  })

  const authApi = useApi(AuthApi)
  const {
    data: publicKeyData,
    isLoading: publicKeyIsLoading,
    isError: publicKeyIsError
  } = useQuery({
    queryFn: async () => {
      const data = await authApi.getPublicKey()
      form.setValue('publicKey', data.publicKey)
      return data
    },
    queryKey: ['getPublicKey'],
    retry: false,
    keepPreviousData: true,
    refetchOnWindowFocus: false
  })
  const {
    data: captchaData,
    isLoading: captchaIsLoading,
    isError: captchaIsError,
    refetch: refreshCaptcha
  } = useQuery({
    queryFn: async () => {
      const data = await authApi.getCaptcha()
      form.setValue('captchaId', data.captchaId)
      return data
    },
    queryKey: ['getCaptcha'],
    retry: false,
    cacheTime: 0,
    refetchOnWindowFocus: false
  })

  const { mutate: signin } = useMutation(
    async (value: LoginRequest) => {
      const res = await authApi.login(value)
      return res
    },
    {
      // TODO: 登录处理
      onSuccess: data => {
        console.log(data)
      },
      onError: error => {
        console.log(error)
      }
    }
  )

  // 表单提交
  const onSubmit = (data: LoginSchemaType) => {
    if (!publicKeyData) return
    const encrypt = new JSEncrypt()
    // 注意：Go 返回的可能是 PEM 格式的 base64，需要先解码
    const publicKey = atob(publicKeyData.publicKey) // 浏览器环境用 atob，Node 用 Buffer
    encrypt.setPublicKey(publicKey)
    const password = encrypt.encrypt(data.password)
    if (!password) {
      return
    }

    signin({
      login: {
        ...data,
        password: password,
        username: data.username
      }
    })
  }

  const isSubmitting = form.formState.isValid && form.formState.isSubmitting
  const isPending = form.formState.isSubmitSuccessful || isSubmitting

  if (publicKeyIsLoading || captchaIsLoading) {
    return <Loading />
  }

  if (publicKeyIsError || captchaIsError) {
    return <ServerErrorPage />
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-amber-50 to-emerald-100 dark:from-amber-900 dark:to-emerald-900 overflow-hidden">
      <div className="absolute inset-0 z-0">
        <div className="absolute inset-0 animate-float">
          {/* 龙猫图标 */}
          <svg
            className="w-32 h-32 absolute top-1/4 left-1/4 text-amber-400/30 dark:text-amber-500/20"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z" />
          </svg>
          {/* 飘浮的树叶 */}
          <svg
            className="w-24 h-24 absolute top-1/3 right-1/3 text-emerald-400/30 dark:text-emerald-500/20 animate-spin-slow"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path d="M17 8C8 10 5.9 16.17 3.82 21.34l1.89.66l.95-2.3c.48.17.98.3 1.5.3C19 20 22 3 22 3c-1 2-8 2.25-13 3.25S2 11.5 2 13.5s1.75 3.75 1.75 3.75C7 8 17 8 17 8z" />
          </svg>
          {/* 额外的装饰元素 */}
          <svg
            className="w-20 h-20 absolute bottom-1/4 right-1/4 text-amber-300/30 dark:text-amber-400/20 animate-bounce-slow"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path d="M12 3c-4.97 0-9 4.03-9 9s4.03 9 9 9s9-4.03 9-9c0-.46-.04-.92-.1-1.36c-.98 1.37-2.58 2.26-4.4 2.26c-3 0-5.44-2.44-5.44-5.44c0-1.81.89-3.41 2.26-4.4C12.92 3.04 12.46 3 12 3z" />
          </svg>
        </div>
      </div>
      <Card className="w-full max-w-md relative z-10 backdrop-blur-md bg-white/60 dark:bg-gray-800/60 shadow-2xl border-opacity-30 border-amber-100 dark:border-amber-700 rounded-3xl transition-all duration-300 hover:shadow-amber-200/50 dark:hover:shadow-amber-700/30">
        <CardHeader className="pb-2">
          <CardTitle className="text-3xl font-bold text-center text-amber-700 dark:text-amber-300 font-ghibli mb-2">
            欢迎回来
          </CardTitle>
          <CardDescription className="text-center text-amber-600/80 dark:text-amber-400/80 text-lg">
            请登录您的账号
          </CardDescription>
        </CardHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="px-6">
            <div className="grid w-full items-center gap-y-2">
              <div className="flex flex-col space-y-1.5">
                <FormField
                  control={form.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>用户名</FormLabel>
                      <FormControl>
                        <Input
                          disabled={isPending}
                          startIcon="user"
                          placeholder="Username"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage namespace="SigninPage" />
                    </FormItem>
                  )}
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>密码</FormLabel>
                      <FormControl>
                        <Input
                          disabled={isPending}
                          startIcon="lock"
                          type="password"
                          placeholder="Password"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage namespace="SigninPage" />
                    </FormItem>
                  )}
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <FormField
                  control={form.control}
                  name="captcha"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>验证码</FormLabel>
                      <FormControl>
                        <div className="w-full gap-x-4 flex justify-between items-center">
                          <Input
                            fullWidth
                            disabled={isPending}
                            startIcon="binary"
                            placeholder="Capthca"
                            {...field}
                          />
                          <div className="h-9 aspect-[3/1] relative rounded-md overflow-hidden border border-[var(--input)] border-solid">
                            {!!captchaData && (
                              <img
                                className={clsx('w-full h-full', [
                                  !isPending
                                    ? 'cursor-pointer pointer-events-auto'
                                    : 'pointer-events-none cursor-none'
                                ])}
                                src={captchaData.picPath}
                                alt="captcha"
                                onClick={() => refreshCaptcha()}
                              />
                            )}
                          </div>
                        </div>
                      </FormControl>
                      <FormMessage namespace="SigninPage" />
                    </FormItem>
                  )}
                />
              </div>
            </div>
            <Button
              size="sm"
              type="submit"
              disabled={isPending}
              className="w-full mt-4 font-medium"
            >
              {isPending && <DynamicIcon name="loader-circle" className="mr-2 animate-spin" />}
              登录
            </Button>
          </form>
        </Form>
      </Card>
    </div>
  )
}

export default LoginPage
