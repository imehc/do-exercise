import { SwitchThumb } from '@radix-ui/react-switch'
import { useQueryClient, useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  ModityTokenStatus,
  SystemTokenApi,
  UpdateTokenStatusRequest,
  TokenInfo,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Switch } from '~/components/ui/switch'

export function ToggleDisabledSwitch({
  accessToken,
  disabled,
}: ModityTokenStatus) {
  const sysTokenApi = useApi(SystemTokenApi)
  const queryClient = useQueryClient()
  // 乐观更新
  const { mutate: update, isPending } = useMutation({
    mutationFn: async (values: UpdateTokenStatusRequest) => {
      await sysTokenApi.updateTokenStatus(values)
      return values.modityTokenStatus.disabled
    },
    onMutate: async ({ modityTokenStatus: { disabled, accessToken } }) => {
      await queryClient.cancelQueries({ queryKey: ['findAllToken'] })
      const previousUsers =
        queryClient.getQueryData<TokenInfo[]>(['findAllToken']) || []
      queryClient.setQueryData<TokenInfo[]>(['findAllToken'], (old = []) =>
        old.map((token) =>
          token.accessToken === accessToken ? { ...token, disabled } : token
        )
      )
      return { previousUsers }
    },
    onError: (_err, _vars, context) => {
      if (context?.previousUsers) {
        queryClient.setQueryData(['findAllToken'], context.previousUsers)
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['findAllToken'] })
    },
  })
  return (
    <Switch
      checked={disabled}
      disabled={isPending}
      onCheckedChange={(value) => {
        update(
          { modityTokenStatus: { accessToken, disabled: value } },
          {
            onSuccess: () => {
              toast.success(`令牌已${value ? '禁用' : '启用'}`)
            },
            onError: () => {
              // toast.error('切换失败，已恢复原状态')
            },
          }
        )
      }}
    >
      <SwitchThumb />
    </Switch>
  )
}
