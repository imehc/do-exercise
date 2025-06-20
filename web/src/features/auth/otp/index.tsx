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
import { OtpForm } from './components/otp-form'

export default function Otp() {
  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader>
          <CardTitle className='text-base tracking-tight'>
            <Trans>双因素认证</Trans>
          </CardTitle>
          <CardDescription>
            <Trans>
              请输入身份验证代码。
              <br />
              我们已经发送了 验证码到您的电子邮件。
            </Trans>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <OtpForm />
        </CardContent>
        <CardFooter>
          {/* <p className='text-muted-foreground px-8 text-center text-sm'>
            Haven't received it?{' '}
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              Resend a new code.
            </Link>
            .
          </p> */}
          <p className='text-muted-foreground px-8 text-center text-sm'>
            <Trans>更喜欢使用密码？</Trans>{' '}
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              <Trans>密码登录</Trans>
            </Link>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
