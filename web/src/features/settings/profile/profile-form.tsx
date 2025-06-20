import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { UpdateUserProfileRequest, UserApi } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
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
import { AvatarUpload } from '~/components/other'
import {
  getProfileSchema,
  ProfileSchemaFormValues,
} from './schemas/action-schema'

export default function ProfileForm() {
  const {
    data: userProfile,
    isLoading: userProfileIsLoading,
    refetch,
  } = useUserProfile()
  const form = useChan(
    useForm<ProfileSchemaFormValues>({
      resolver: resolver(getProfileSchema()),
      mode: 'onChange',
    })
  )

  useEffect(() => {
    form.reset({
      nickname: userProfile?.nickname,
      avatar: userProfile?.avatar,
    })
  }, [form, userProfile?.avatar, userProfile?.nickname])

  const userApi = useApi(UserApi)
  const { mutate: updateUserProfile, isPending: updateUserProfileIspending } =
    useMutation({
      mutationFn: (value: UpdateUserProfileRequest) =>
        userApi.updateUserProfile(value),
      onSuccess: () => {
        toast.success(t`更新成功`)
        refetch()
      },
    })

  const isPending = userProfileIsLoading || updateUserProfileIspending

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) =>
          updateUserProfile({ updateUserProfile: data })
        )}
        className='space-y-8'
      >
        <FormItem>
          <FormLabel>
            <Trans>用户名</Trans>
          </FormLabel>
          <FormControl>
            <Input
              placeholder={t`请输入用户名`}
              defaultValue={userProfile?.username}
              readOnly
              disabled
            />
          </FormControl>
          <FormDescription>
            <Trans>这是您的公开显示名称。</Trans>
          </FormDescription>
          <FormMessage />
        </FormItem>

        <FormField
          control={form.control}
          name='nickname'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>昵称</Trans>
              </FormLabel>
              <Input
                placeholder={t`请输入昵称`}
                {...field}
                disabled={updateUserProfileIspending}
              />
              <FormDescription>
                <Trans>这是您的昵称。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='avatar'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>头像</Trans>
              </FormLabel>
              <FormControl>
                <AvatarUpload
                  // TODO: 根据前缀判断是否需要拼接完整的图片地址
                  disabled={updateUserProfileIspending}
                  value={field.value}
                  onChange={field.onChange}
                />
              </FormControl>
              <FormDescription>
                <Trans>这是您的头像。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {updateUserProfileIspending ? (
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
