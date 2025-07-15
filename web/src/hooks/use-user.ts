import { useMutation, useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { useSetAtom } from 'jotai'
import { originTokenAtom, routeTabCacheAtom } from '~/atoms'
import { AuthApi, ResponseError, UserApi } from '~/do-exercise-api'
import { useApi } from './use-api'

/** 退出登录 */
export const useLogout = () => {
  const setOriginToken = useSetAtom(originTokenAtom)
  const setRouteTabCache = useSetAtom(routeTabCacheAtom)
  const navigate = useNavigate()
  const authApi = useApi(AuthApi)
  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      setRouteTabCache({
        activeTab: '',
        tabs: [],
      })
      setOriginToken({
        accessToken: '',
        expireTime: 0,
        refreshToken: '',
        refreshExpireTime: 0,
      })
      navigate({ to: '/sign-in', replace: true })
    },
  })
}

/** 获取用户信息 */
export const useUserProfile = () => {
  const userApi = useApi(UserApi)
  return useQuery({
    queryKey: ['getUserProfile'],
    queryFn: () => userApi.getUserProfile(),
    retry: (failureCount, error) => {
      if (error instanceof ResponseError) {
        if (error.response.status == 401 || error.response?.status === 403) {
          return false
        }
      }
      // 其他错误，最多重试 3 次
      return failureCount < 3
    },
  })
}
