import { useQuery } from '@tanstack/react-query'
import { IconMessage } from '@tabler/icons-react'
import { SystemUserApi, SysUser } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
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
            <IconMessage /> 用户详情
          </DrawerTitle>
          <DrawerDescription>查看用户详细信息。</DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoadingUser} data={data}>
          {(user) => (
            <div className='grid gap-6 overflow-y-auto px-4'>
              <div className='space-y-3'>
                <h4 className='text-lg font-medium'>基本信息</h4>
                <div className='grid gap-3 text-sm'>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>用户名</div>
                    <div className='col-span-2'>{user.username}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>头像</div>
                    <div className='col-span-2'>
                      {
                        // TODO: 根据前缀判断是否需要拼接完整的图片地址
                        user?.avatar ? (
                          <Avatar>
                            <AvatarImage
                              src={user.avatar}
                              alt={user.username}
                            />
                            <AvatarFallback>
                              {user.username.slice(0, 2)}
                            </AvatarFallback>
                          </Avatar>
                        ) : (
                          '-'
                        )
                      }
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>昵称</div>
                    <div className='col-span-2'>{user.nickname}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>邮箱</div>
                    <div className='col-span-2'>{user.email || '-'}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>关联角色</div>
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
            <Button variant='outline'>关闭</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
