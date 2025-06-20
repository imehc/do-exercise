import { Link } from '@tanstack/react-router'
import { Trans } from '@lingui/react/macro'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card'
import AuthLayout from '../auth-layout'
import { EmailSignInForm } from './components/email-sign-in-form'

export default function EmailSignIn() {
  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader>
          <CardTitle className='text-lg tracking-tight'>
            <Trans>使用邮箱登录</Trans>
          </CardTitle>
          <CardDescription>
            <Trans>
              请输入您的邮箱地址，我们将向您的邮箱发送一个验证码，请使用验证码登录。
            </Trans>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <EmailSignInForm />
        </CardContent>
        <CardFooter className='justify-center'>
          <p className='text-muted-foreground px-8 text-center text-sm'>
            <Trans>更喜欢使用密码？</Trans>
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              <Trans>使用密码登录</Trans>
            </Link>
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
