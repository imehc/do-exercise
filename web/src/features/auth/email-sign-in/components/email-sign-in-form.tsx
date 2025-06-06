import { HTMLAttributes } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { AuthApi, GetForgetPasswordCodeRequest } from '~/do-exercise-api'
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

type ForgotFormProps = HTMLAttributes<HTMLFormElement>

const formSchema = z.object({
  email: z
    .string()
    .min(1, { message: 'Please enter your email' })
    .email({ message: 'Invalid email address' }),
})

export function EmailSignInForm({ className, ...props }: ForgotFormProps) {
  const navigate = useNavigate()
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { email: '' },
  })

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

  function onSubmit(data: z.infer<typeof formSchema>) {
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
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder='email' {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isPending}>
          Continue
        </Button>
      </form>
    </Form>
  )
}
