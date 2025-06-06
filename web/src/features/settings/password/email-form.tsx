import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { toast } from 'sonner'
import { ModifyPasswordRequest, UserApi } from '~/do-exercise-api'
import { encryptPassword } from '~/utils/encrypt'
import { useApi } from '~/hooks/use-api'
import { usePublicKey } from '~/hooks/use-public-key'
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
import { PasswordInput } from '~/components/password-input'
import {
  passwordSchema,
  PasswordSchemaFormValues,
} from './schemas/action-schema'

export default function PasswordForm() {
  const form = useForm<PasswordSchemaFormValues>({
    resolver: resolver(passwordSchema),
    mode: 'onChange',
  })

  const userApi = useApi(UserApi)
  const {
    data: publicKeyData,
    isLoading: publicKeyDataIsLoading,
    refetch: refetchPublicKey,
  } = usePublicKey()

  const { mutate: modifyPassword, isPending: modifyPasswordIspending } =
    useMutation({
      mutationFn: (value: ModifyPasswordRequest) =>
        userApi.modifyPassword(value),
      onSuccess: () => {
        toast.success(`更改密码成功`)
        refetchPublicKey()
      },
      onError: () => {
        refetchPublicKey()
      },
    })

  const onSubmit = (data: PasswordSchemaFormValues) => {
    if (!publicKeyData) return
    const oldPassword = encryptPassword(
      data.oldPassword,
      publicKeyData.publicKey
    )
    const password = encryptPassword(data.password, publicKeyData.publicKey)
    if (!oldPassword || !password) {
      return
    }
    modifyPassword({
      modifyPassword: {
        oldPassword: oldPassword,
        password: password,
        publicKey: publicKeyData.publicKey,
      },
    })
  }

  const isPending = modifyPasswordIspending || publicKeyDataIsLoading

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) => onSubmit(data))}
        className='space-y-8'
      >
        <FormField
          control={form.control}
          name='oldPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Old Password</FormLabel>
              <PasswordInput
                placeholder='old password'
                {...field}
                disabled={modifyPasswordIspending}
              />
              <FormDescription>This is your old password.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='password'
                  {...field}
                  disabled={modifyPasswordIspending}
                />
              </FormControl>
              <FormDescription>This is your new password.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='confirmPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='confirm password'
                  {...field}
                  disabled={modifyPasswordIspending}
                />
              </FormControl>
              <FormDescription>Enter your new password again.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {modifyPasswordIspending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>更新中...</span>
            </>
          ) : (
            <span>更新</span>
          )}
        </Button>
      </form>
    </Form>
  )
}
