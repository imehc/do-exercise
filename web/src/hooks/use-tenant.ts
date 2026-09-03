import { useQuery } from '@tanstack/react-query'
import { useAtomValue } from 'jotai'
import { originTokenAtom } from '~/atoms'
import { AuthApi, AvailableTenants } from '~/do-exercise-api'
import { useApi } from './use-api'

/**
 * 匿名引导接口：权限动作词表等静态词表由它一次下发。
 * 这些值由后端常量决定，运行期不会变化，因此整个会话只取一次。
 */
function useAvailableTenants() {
  const authApi = useApi(AuthApi)
  return useQuery<AvailableTenants>({
    queryKey: ['availableTenants'],
    queryFn: () => authApi.availableTenants(),
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
  })
}

/**
 * 菜单按钮的权限动作词表。唯一来源在后端 `global.MenuPermissionActions`，
 * 前端不再各自维护一份：新增动作只要改 Go 常量，界面下拉自动跟随。
 * 词表未到达前返回空数组，由调用方决定兜底（见 action-schema 的
 * fallbackMenuPermissionActions），保证首屏就能渲染出下拉项。
 */
export function useMenuPermissionActions(): string[] {
  const { data } = useAvailableTenants()
  return data?.permissionActions ?? []
}

/**
 * 平台统一维护的资源（菜单 / API）在租户侧只读：非平台超管一律不显示新增、编辑、删除入口。
 * 服务端同样拦截（menuReadonlyForTenant / apiReadonlyForTenant），
 * 这里只是不把注定失败的入口摆在用户面前。
 */
export function usePlatformResourceReadonly(): boolean {
  return !useAtomValue(originTokenAtom).isSuperAdmin
}
