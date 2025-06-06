import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { toast } from 'sonner'
import {
  BindEmailRequest,
  GetBindEmailCodeRequest,
  GetRebindEmailCodeRequest,
  RebindEmailRequest,
  UserApi,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { useCountdown } from '~/hooks/use-count-down'
import { useUserProfile } from '~/hooks/use-user'
import { Button } from '~/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { emailSchema, EmailSchemaFormValues } from './schemas/action-schema'

export default function EmailForm() {
  const { isCounting, count, start } = useCountdown()

  const {
    data: userProfile,
    isLoading: userProfileIsLoading,
    refetch,
  } = useUserProfile()
  const form = useForm<EmailSchemaFormValues>({
    resolver: resolver(
      emailSchema.superRefine(({ email }, ctx) => {
        if (userProfile?.email === email.trim()) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ['email'],
            message: 'The mailbox cannot be the same as the current mailbox.',
          })
        }
      })
    ),
    mode: 'onChange',
  })

  const userApi = useApi(UserApi)
  const { mutate: bindEmail, isPending: bindEmailIspending } = useMutation({
    mutationFn: (value: BindEmailRequest) => userApi.bindEmail(value),
    onSuccess: () => {
      toast.success(`绑定成功`)
      form.reset({
        email: '',
        code: '',
      })
      refetch()
    },
  })

  const { mutate: rebindEmail, isPending: rebindEmailIspending } = useMutation({
    mutationFn: (value: RebindEmailRequest) => userApi.rebindEmail(value),
    onSuccess: () => {
      toast.success(`更换成功`)
      form.reset({
        email: '',
        code: '',
      })
      refetch()
    },
  })

  const { mutate: getBindEmailCode, isPending: bindEmailCodeIspending } =
    useMutation({
      mutationFn: (value: GetBindEmailCodeRequest) =>
        userApi.getBindEmailCode(value),
      onSuccess: () => {
        toast.success(`发送成功`)
        start()
      },
    })

  const { mutate: getRebindEmailCode, isPending: rebindEmailCodeIspending } =
    useMutation({
      mutationFn: (value: GetRebindEmailCodeRequest) =>
        userApi.getRebindEmailCode(value),
      onSuccess: () => {
        toast.success(`发送成功`)
        start()
      },
    })

  const hasBind = !userProfile?.email

  const isPending =
    userProfileIsLoading || bindEmailIspending || rebindEmailIspending

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) =>
          hasBind
            ? bindEmail({ bindEmail: data })
            : rebindEmail({ bindEmail: data })
        )}
        className='space-y-8'
      >
        <FormItem>
          <FormLabel>Username</FormLabel>
          <FormControl>
            <Input
              placeholder='username'
              defaultValue={userProfile?.username}
              readOnly
              disabled
            />
          </FormControl>
          <FormDescription>This is your public display name.</FormDescription>
          <FormMessage />
        </FormItem>

        {userProfile?.email && (
          <FormItem>
            <FormLabel>Email</FormLabel>
            <FormControl>
              <Input
                placeholder='email'
                defaultValue={userProfile.email}
                readOnly
                disabled
              />
            </FormControl>
            <FormDescription>This is your current email.</FormDescription>
            <FormMessage />
          </FormItem>
        )}

        <FormField
          control={form.control}
          name='email'
          render={({ field }) => (
            <FormItem>
              <FormLabel>{hasBind ? 'Bind' : 'Rebind'} Email</FormLabel>
              <div className='flex justify-between gap-x-2'>
                <Input placeholder='email' {...field} />
                <Button
                  type='button'
                  className='w-1/3'
                  variant='outline'
                  disabled={
                    !form.watch('email') ||
                    form.watch('email')?.trim() === userProfile?.email ||
                    isCounting ||
                    bindEmailCodeIspending ||
                    rebindEmailCodeIspending
                  }
                  onClick={() => {
                    const email = form.watch('email')
                    if (!email) return
                    if (hasBind) {
                      getBindEmailCode({ email })
                      return
                    }
                    getRebindEmailCode({ email })
                  }}
                >
                  {isCounting ? count : '获取验证码'}
                </Button>
              </div>
              <FormDescription>
                This is the new email address you have bound.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='code'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Code</FormLabel>
              <FormControl>
                <Input
                  placeholder='code'
                  {...field}
                  disabled={bindEmailIspending || rebindEmailIspending}
                />
              </FormControl>
              <FormDescription>
                Please enter the code you received in your email.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {bindEmailIspending || rebindEmailIspending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>{hasBind ? '绑定中...' : '换绑中...'}</span>
            </>
          ) : (
            <span>{hasBind ? '绑定' : '换绑'}</span>
          )}
        </Button>
      </form>
    </Form>
  )
}
