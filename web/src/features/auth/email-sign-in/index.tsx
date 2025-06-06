import { Link } from '@tanstack/react-router'
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
            Sign in with Email
          </CardTitle>
          <CardDescription>
            Enter your email address and we’ll send you a login link to access
            your account.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <EmailSignInForm />
        </CardContent>
        <CardFooter>
          <p className='text-muted-foreground px-8 text-center text-sm'>
            Prefer using a password?{' '}
            <Link
              to='/sign-in'
              className='hover:text-primary underline underline-offset-4'
            >
              Sign in with password
            </Link>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
