import { useQuery } from '@tanstack/react-query'
import { IconMessage } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { SystemUserApi, SysUser } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { ensureHttpPrefix } from '~/utils/url'
import { useApi } from '~/hooks/use-api'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from '~/components/ui/drawer'
import { StatusRenderer } from '~/components/other'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function UserViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysUser>()
  const systemUserApi = useApi(SystemUserApi)
  const { data, isLoading: isLoadingUser } = useQuery({
    queryKey: ['findUser', currentRow?.id],
    queryFn: () => systemUserApi.findUser({ id: currentRow?.id as string }),
    enabled: !!currentRow?.id,
  })

  return (
    <Drawer
      shouldScaleBackground
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
      direction='bottom'
    >
      <DrawerContent className='max-h-[85vh]'>
        <DrawerHeader className='text-left'>
          <DrawerTitle className='flex items-center gap-2'>
            <IconMessage /> <Trans>用户详情</Trans>
          </DrawerTitle>
          <DrawerDescription>
            <Trans>查看用户详细信息。</Trans>
          </DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoadingUser} data={data}>
          {(user) => (
            <div className='grid gap-6 overflow-y-auto px-4'>
              <div className='space-y-3'>
                <h4 className='text-lg font-medium'>
                  <Trans>基本信息</Trans>
                </h4>
                <div className='grid gap-3 text-sm'>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>用户名</Trans>
                    </div>
                    <div className='col-span-2'>{user.username}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>头像</Trans>
                    </div>
                    <div className='col-span-2'>
                      {user?.avatar ? (
                        <Avatar>
                          <AvatarImage
                            src={ensureHttpPrefix(user.avatar)}
                            alt={user.username}
                          />
                          <AvatarFallback>
                            {user.username.slice(0, 2)}
                          </AvatarFallback>
                        </Avatar>
                      ) : (
                        '-'
                      )}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>昵称</Trans>
                    </div>
                    <div className='col-span-2'>{user.nickname}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>邮箱</Trans>
                    </div>
                    <div className='col-span-2'>{user.email || '-'}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>关联角色</Trans>
                    </div>
                    <div className='col-span-2'>
                      {user.roles.map((role) => (
                        <Badge key={role.id} className='mr-2 last-of-type:mr-0'>
                          {role.name}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </StatusRenderer>

        <DrawerFooter className='gap-y-2'>
          <DrawerClose asChild>
            <Button variant='outline'>
              <Trans>关闭</Trans>
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
