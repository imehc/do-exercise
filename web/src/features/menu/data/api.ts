import { queryOptions } from '@tanstack/react-query'
import { SystemMenuApi } from '~/do-exercise-api'
import { apiInstance } from '~/hooks/use-api'

export const findMenuTree = () => {
  const sysMenuApi = apiInstance(SystemMenuApi)
  return queryOptions({
    queryKey: ['findMenuTree'],
    queryFn: () => sysMenuApi.findMenuTree(),
  })
}
