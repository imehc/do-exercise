import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { ModifyPasswordRequest, UserApi } from '~/do-exercise-api'
import { originTokenAtom, store } from '~/atoms'
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
  getPasswordSchema,
  PasswordSchemaFormValues,
} from './schemas/action-schema'

export default function PasswordForm() {
  const form = useForm<PasswordSchemaFormValues>({
    resolver: resolver(getPasswordSchema()),
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
        toast.success(t`更改密码成功`)
        // 改密成功后后端保留当前会话并清除了强制改密标记，前端同步清除本地标记，
        // 否则 use-api 的 403 处理仍会把用户拽回改密页。
        const token = store.get(originTokenAtom)
        store.set(originTokenAtom, { ...token, mustChangePassword: false })
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
              <FormLabel>
                <Trans>原密码</Trans>
              </FormLabel>
              <PasswordInput
                placeholder={t`请输入原密码`}
                {...field}
                disabled={modifyPasswordIspending}
              />
              <FormDescription>
                <Trans>这是您之前的设置的密码</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>新密码</Trans>
              </FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder={t`请输入新密码`}
                  {...field}
                  disabled={modifyPasswordIspending}
                />
              </FormControl>
              <FormDescription>
                <Trans>这是您设置的新密码</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='confirmPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel>确认密码</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder={t`请再次输入您的新密码`}
                  {...field}
                  disabled={modifyPasswordIspending}
                />
              </FormControl>
              <FormDescription>
                <Trans>这是您设置的新密码</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {modifyPasswordIspending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>
                <Trans>保存中</Trans>...
              </span>
            </>
          ) : (
            <span>
              <Trans>保存</Trans>
            </span>
          )}
        </Button>
      </form>
    </Form>
  )
}
