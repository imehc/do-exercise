import { HTMLAttributes, useEffect } from 'react'
import { z } from 'zod'
import { addSeconds } from 'date-fns'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { useSetAtom } from 'jotai'
import { toast } from 'sonner'
import { originTokenAtom } from '~/atoms'
import { AuthApi, LoginWithEmailRequest } from '~/do-exercise-api'
import { Route } from '~/routes/(auth)/otp'
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

type OtpFormProps = HTMLAttributes<HTMLFormElement>

const formSchema = z.object({
  code: z.string().min(1, { message: 'Please enter your email code.' }),
  email: z
    .string()
    .min(1, { message: 'Please enter your email' })
    .email({ message: 'Invalid email address' }),
})

export function OtpForm({ className, ...props }: OtpFormProps) {
  const setToken = useSetAtom(originTokenAtom)
  const navigate = useNavigate()
  const { email } = Route.useSearch()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { code: '', email: '' },
  })

  useEffect(() => {
    if (!email?.trim()) {
      toast.error('email invalid')
      return
    }
    form.setValue('email', email)
  })

  const authApi = useApi(AuthApi)
  const { mutate: loginWithEmail, isPending } = useMutation({
    mutationFn: (value: LoginWithEmailRequest) => authApi.loginWithEmail(value),
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
  })

  function onSubmit(data: z.infer<typeof formSchema>) {
    loginWithEmail({ loginWithEmail: data })
  }

  return (
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
              <FormLabel className='sr-only'>One-Time Password</FormLabel>
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
        <Button
          className='mt-2'
          disabled={form.watch('code').length < 6 || isPending}
        >
          Verify
        </Button>
      </form>
    </Form>
  )
}
