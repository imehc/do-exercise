import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import { Card, CardHeader, CardTitle, CardDescription } from '~/components/ui/card'
import '~/animations.css'
import { useApi } from '~/hooks'
import { AuthApi } from '~/do-exercise-api'
import { useQuery } from '@tanstack/react-query'
import { useForm, Controller } from "react-hook-form";
import { zodResolver as resolver } from "@hookform/resolvers/zod";
import { Loading } from '~/components'
import { loginSchema, LoginSchemaType } from './schema.'

export function LoginPage() {
  const authApi = useApi(AuthApi)
  const { data, isLoading } = useQuery({
    queryFn: async () => await authApi.getPublicKey(),
    queryKey: ['getPublicKey'],
    retry: false,
  })

  const { formState: { errors }, handleSubmit, control } = useForm<LoginSchemaType>({
    resolver: resolver(loginSchema)
  })

  const onSubmit = (data: LoginSchemaType) => {
    console.log(data)
  }

  // TOOD: 处理没有数据和过期以及报错情况
  if (isLoading || !data) {
    return <Loading />
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
        <form className="space-y-6 p-6 pt-2" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <div className="relative group">
              <Controller
                control={control}
                name="username"
                render={({ field }) => (
                  <div className="relative h-[4.5rem]">
                    <Input
                      {...field}
                      type="text"
                      placeholder="用户名"
                      className={`w-full h-11 transition-all duration-200 bg-white/50 dark:bg-gray-800/50 border-amber-200 dark:border-amber-700/50 hover:border-amber-400 focus:border-amber-400 focus:ring-amber-400/30 text-lg pl-10 ${errors.username ? 'border-red-500 focus:border-red-500 focus:ring-red-500/30' : ''}`}
                    />
                    <div className="absolute left-3 top-[1.375rem] -translate-y-1/2 pointer-events-none">
                      <svg
                        className="w-5 h-5 text-amber-400/70 dark:text-amber-500/70"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                        />
                      </svg>
                    </div>
                    {errors.username && (
                      <p className="text-sm text-red-500 mt-1">{errors.username.message}</p>
                    )}
                  </div>
                )}
              />
            </div>
            <div className="relative group">
              <Controller
                control={control}
                name="password"
                render={({ field }) => (
                  <div className="relative h-[4.5rem]">
                    <Input
                      {...field}
                      type="password"
                      placeholder="密码"
                      className={`w-full h-11 transition-all duration-200 bg-white/50 dark:bg-gray-800/50 border-amber-200 dark:border-amber-700/50 hover:border-amber-400 focus:border-amber-400 focus:ring-amber-400/30 text-lg pl-10 ${errors.password ? 'border-red-500 focus:border-red-500 focus:ring-red-500/30' : ''}`}
                    />
                    <div className="absolute left-3 top-[1.375rem] -translate-y-1/2 pointer-events-none">
                      <svg
                        className="w-5 h-5 text-amber-400/70 dark:text-amber-500/70"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                        />
                      </svg>
                    </div>
                    {errors.password && (
                      <p className="text-sm text-red-500 mt-1">{errors.password.message}</p>
                    )}
                  </div>
                )}
              />
            </div>
            <div className="relative group flex gap-4">
              <div className="relative flex-1">
                <Input
                  type="text"
                  placeholder="验证码"
                  className="w-full h-11 transition-all duration-200 bg-white/50 dark:bg-gray-800/50 border-amber-200 dark:border-amber-700/50 hover:border-amber-400 focus:border-amber-400 focus:ring-amber-400/30 text-lg pl-10"
                />
                <div className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none">
                  <svg
                    className="w-5 h-5 text-amber-400/70 dark:text-amber-500/70"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
              </div>
              <div className="w-32 h-11 bg-gradient-to-r from-amber-400/20 to-emerald-400/20 rounded-lg flex items-center justify-center text-lg font-mono tracking-wider text-amber-700 dark:text-amber-300 select-none cursor-pointer hover:from-amber-400/30 hover:to-emerald-400/30 transition-all duration-300">
                AB12CD
              </div>
            </div>
          </div>
          <div className="flex items-center justify-between text-sm">
            <label className="flex items-center space-x-2 cursor-pointer group">
              <input
                type="checkbox"
                className="rounded border-amber-300 text-amber-500 focus:ring-amber-400/30 transition-colors"
              />
              <span className="text-gray-600 dark:text-gray-400 group-hover:text-amber-500 transition-colors">
                记住我
              </span>
            </label>
            <a
              href="#"
              className="text-amber-500 hover:text-amber-400 transition-colors hover:underline"
            >
              忘记密码？
            </a>
          </div>
          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-amber-500 to-amber-400 text-white font-medium transition-all hover:from-amber-400 hover:to-amber-300 active:scale-[0.98] transform duration-200 shadow-lg hover:shadow-amber-300/50 dark:hover:shadow-amber-700/30 rounded-xl py-4 px-6"
          >
            登录
          </Button>
        </form>
      </Card>
    </div>
  )
}

export default LoginPage
