import { HTMLAttributes } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { i18n } from '@lingui/core'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { AuthApi, GetForgetPasswordCodeRequest } from '~/do-exercise-api'
import { Route } from '~/routes/(auth)/email-sign-in'
import { cn } from '~/lib/utils'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
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

type ForgotFormProps = HTMLAttributes<HTMLFormElement>

const getFormSchema = () =>
  z.object({
    email: z
      .string()
      .min(1, { message: t`请输入您的邮箱` })
      .email({ message: t`邮箱无效` }),
  })

type FormSchemaValues = z.infer<ReturnType<typeof getFormSchema>>

export function EmailSignInForm({ className, ...props }: ForgotFormProps) {
  const navigate = Route.useNavigate()
  const form = useChan(
    useForm<FormSchemaValues>({
      resolver: zodResolver(getFormSchema()),
      defaultValues: { email: '' },
    })
  )

  const authApi = useApi(AuthApi)
  const { mutate: getForgetPasswordCode, isPending } = useMutation({
    mutationFn: async (value: GetForgetPasswordCodeRequest) => {
      await authApi.getCodeWithEmail(value)
      return value
    },
    onSuccess: ({ email }) => {
      navigate({
        to: '/otp',
        search: { email },
      })
    },
  })

  function onSubmit(data: FormSchemaValues) {
    getForgetPasswordCode(data)
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
          name='email'
          render={({ field }) => (
            <FormItem className='space-y-1'>
              <FormLabel>
                <Trans>邮箱</Trans>
              </FormLabel>
              <FormControl>
                <Input placeholder={i18n._(t`请输入邮箱`)} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isPending}>
          <Trans>继续</Trans>
        </Button>
      </form>
    </Form>
  )
}
