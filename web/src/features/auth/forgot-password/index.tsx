import { Link } from '@tanstack/react-router'
import { Trans } from '@lingui/react/macro'
import { ForgotPasswordForm } from './components/forgot-password-form'

export default function ForgotPassword() {
  return (
    <div className='relative container grid h-svh flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0'>
      <div className='bg-muted relative hidden h-full flex-col p-10 text-white lg:flex dark:border-r'>
        <div className='absolute inset-0 bg-zinc-900' />
        <div className='relative z-20 flex items-center text-lg font-medium'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeWidth='2'
            strokeLinecap='round'
            strokeLinejoin='round'
            className='mr-2 h-6 w-6'
          >
            <path d='M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3' />
          </svg>
          <Trans>敷了管理系统</Trans>
        </div>

        <img
          src='https://img.daisyui.com/images/daisyui/mark-rotating.svg'
          className='relative m-auto'
          width={480}
          height={480}
          alt='daisyui'
        />

        <div className='relative z-20 mt-auto'>
          <blockquote className='space-y-2'>
            <p className='text-lg'>
              &ldquo;<Trans>用于学习基本的权限管理系统及流程</Trans>&rdquo;
            </p>
            <footer className='text-sm'>
              <Trans>Imehc</Trans>
            </footer>
          </blockquote>
        </div>
      </div>
      <div className='lg:p-8'>
        <div className='mx-auto flex w-full flex-col justify-center space-y-2 sm:w-[350px]'>
          <div className='flex flex-col space-y-2 text-left'>
            <h1 className='text-2xl font-semibold tracking-tight'>
              <Trans>忘记密码</Trans>
            </h1>
            <p className='text-muted-foreground text-sm'>
              输入您的注册电子邮件地址，我们将向您的电子邮件发送重置验证码。
            </p>
          </div>
          <ForgotPasswordForm />
          <p className='text-muted-foreground px-8 text-center text-sm'>
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              <Trans>登录</Trans>
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
