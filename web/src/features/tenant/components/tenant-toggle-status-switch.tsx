import { SwitchThumb } from '@radix-ui/react-switch'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { t } from '@lingui/core/macro'
import { toast } from 'sonner'
import {
  SystemTenantApi,
  Tenant,
  TenantList,
  UpdateTenantRequest,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { Switch } from '~/components/ui/switch'

interface Props {
  tenantId: string
  name: string
  remark?: string
  status: boolean
}

export function TenantToggleStatusSwitch({
  tenantId,
  name,
  remark,
  status,
}: Props) {
  const systemTenantApi = useApi(SystemTenantApi)
  const queryClient = useQueryClient()
  // 乐观更新：先更新所有分页缓存中的该租户状态，失败再回滚
  const { mutate: update, isPending } = useMutation({
    mutationFn: async (values: UpdateTenantRequest) => {
      await systemTenantApi.updateTenant(values)
      return values.updateTenant.status
    },
    onMutate: async (values) => {
      await queryClient.cancelQueries({ queryKey: ['findTenants'] })
      const previousQueries: Array<[readonly unknown[], unknown]> = []
      queryClient
        .getQueryCache()
        .findAll({ queryKey: ['findTenants'] })
        .forEach((query) => {
          previousQueries.push([query.queryKey, query.state.data])
          queryClient.setQueryData(
            query.queryKey,
            (old: TenantList | undefined) => {
              if (!old) return old
              return {
                ...old,
                data: old.data.map((tenant: Tenant) =>
                  tenant.tenantId === values.id
                    ? { ...tenant, status: values.updateTenant.status }
                    : tenant
                ),
              }
            }
          )
        })
      return { previousQueries }
    },
    onError: (_err, _vars, context) => {
      context?.previousQueries.forEach(([key, data]) => {
        queryClient.setQueryData(key, data)
      })
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['findTenants'] })
    },
  })
  return (
    <Switch
      checked={status}
      disabled={isPending}
      onCheckedChange={(value) => {
        update(
          { id: tenantId, updateTenant: { name, remark, status: value } },
          {
            onSuccess: () => {
              toast.success(value ? t`租户已启用` : t`租户已停用`)
            },
          }
        )
      }}
    >
      <SwitchThumb />
    </Switch>
  )
}
