import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { AuthApi } from '~/do-exercise-api'
import { useApi } from './use-api'

export function usePublicKey() {
  const authApi = useApi(AuthApi)

  return useQuery({
    queryKey: ['getPublicKey'],
    queryFn: () => authApi.getPublicKey(),
    retry: false,
    placeholderData: keepPreviousData,
    refetchInterval: 5 * 60 * 1000,
  })
}
