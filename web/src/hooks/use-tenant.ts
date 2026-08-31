import { useQuery } from '@tanstack/react-query'
import { useAtomValue } from 'jotai'
import { originTokenAtom } from '~/atoms'
import { AuthApi } from '~/do-exercise-api'
import { useApi } from './use-api'

/**
 * 是否多租户部署。mode 由后端配置决定（single / multi），运行期不会变化，
 * 因此整个会话只取一次。加载完成前返回 false（按单租户处理），避免误藏功能入口。
 */
export function useIsMultiTenant(): boolean {
  const authApi = useApi(AuthApi)
  const { data } = useQuery({
    queryKey: ['availableTenants'],
    queryFn: () => authApi.availableTenants(),
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
  })
  return data?.mode === 'multi'
}

/**
 * 平台统一维护的资源（菜单 / API）在租户侧只读：
 * 多租户模式下的非超管一律不显示新增、编辑、删除入口。
 * 服务端同样拦截（menuReadonlyForTenant / apiReadonlyForTenant），
 * 这里只是不把注定失败的入口摆在用户面前。
 */
export function usePlatformResourceReadonly(): boolean {
  const isMulti = useIsMultiTenant()
  const isSuperAdmin = !!useAtomValue(originTokenAtom).isSuperAdmin
  return isMulti && !isSuperAdmin
}
