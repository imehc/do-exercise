import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver as resolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { toast } from 'sonner'
import { UpdateUserProfileRequest, UserApi } from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
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
import { profileSchema, ProfileSchemaFormValues } from './schemas/action-schema'

export default function ProfileForm() {
  const {
    data: userProfile,
    isLoading: userProfileIsLoading,
    refetch,
  } = useUserProfile()
  const form = useForm<ProfileSchemaFormValues>({
    resolver: resolver(profileSchema),
    mode: 'onChange',
  })

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
        toast.success(`更新成功`)
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

        <FormField
          control={form.control}
          name='nickname'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Nickname</FormLabel>
              <Input
                placeholder='nickname'
                {...field}
                disabled={updateUserProfileIspending}
              />
              <FormDescription>This is your nickname.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='avatar'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Avatar</FormLabel>
              <FormControl>
                <AvatarUpload
                  // TODO: 根据前缀判断是否需要拼接完整的图片地址
                  disabled={updateUserProfileIspending}
                  value={field.value}
                  onChange={field.onChange}
                />
              </FormControl>
              <FormDescription>This is your avatar.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit' disabled={isPending} className='max-sm:w-full'>
          {updateUserProfileIspending ? (
            <>
              <IconLoader3 className='animate-spin' />
              <span>保存中...</span>
            </>
          ) : (
            <span>保存</span>
          )}
        </Button>
      </form>
    </Form>
  )
}
