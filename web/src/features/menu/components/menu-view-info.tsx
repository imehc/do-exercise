import { format } from 'date-fns'
import { useQuery } from '@tanstack/react-query'
import { Icon, IconMessage } from '@tabler/icons-react'
import * as icons from '@tabler/icons-react'
import { MenuType, SysMenuTree, SystemMenuApi } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
import { useApi } from '~/hooks/use-api'
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
import { iconPrefix, StatusRenderer } from '~/components/other'
import { callMethodTypes } from '~/features/api/data/data'
import { callMenuMapping, callMenuTypes, callVisibleTypes } from '../data/data'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function MenuViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysMenuTree>()
  const sysMenuApi = useApi(SystemMenuApi)
  const { data, isLoading: isLoadingMenu } = useQuery({
    queryKey: ['findMenu', currentRow?.id],
    queryFn: () => sysMenuApi.findMenu({ id: currentRow?.id as number }),
    enabled: !!currentRow?.id,
  })

  return (
    <Drawer
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
      direction='bottom'
    >
      <DrawerContent className='max-h-[85vh]'>
        <DrawerHeader className='text-left'>
          <DrawerTitle className='flex items-center gap-2'>
            <IconMessage /> 菜单详情
          </DrawerTitle>
          <DrawerDescription>
            查看菜单详细信息。这包括 不同菜单类型对应不同的信息。
          </DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoadingMenu} data={data}>
          {(menu) => {
            const SelectedIcon =
              menu.icon && menu.type === MenuType.menu
                ? (icons[
                    (iconPrefix + menu.icon) as keyof typeof icons
                  ] as Icon)
                : null
            return (
              <div className='grid gap-6 overflow-y-auto px-4'>
                <div className='space-y-3'>
                  <h4 className='text-lg font-medium'>基本信息</h4>
                  <div className='grid gap-3 text-sm'>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>菜单名称</div>
                      <div className='col-span-2'>{menu.name}</div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>菜单类型</div>
                      <div className='col-span-2'>
                        <Badge
                          variant='outline'
                          className={cn(
                            'capitalize',
                            callMenuTypes.get(menu.type)
                          )}
                        >
                          {callMenuMapping.get(menu.type) ?? '-'}
                        </Badge>
                      </div>
                    </div>
                    {menu.type === MenuType.menu && (
                      <>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>路由</div>
                          <div className='col-span-2'>{menu.route}</div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>组件</div>
                          <div className='col-span-2'>{menu.component}</div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>图标</div>
                          <div className='col-span-2'>
                            {SelectedIcon && <SelectedIcon />}
                          </div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>是否可见</div>
                          <div className='col-span-2'>
                            <Badge
                              variant='outline'
                              className={cn(
                                'capitalize',
                                callVisibleTypes.get(menu.visible ?? false)
                              )}
                            >
                              {menu.visible ? '是' : '否'}
                            </Badge>
                          </div>
                        </div>
                      </>
                    )}
                    {menu.type !== MenuType.button && (
                      <div className='grid grid-cols-3 items-center gap-4'>
                        <div className='font-medium'>序号</div>
                        <div className='col-span-2'>{menu.sort}</div>
                      </div>
                    )}
                    {menu.type === MenuType.button && (
                      <>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>权限标识</div>
                          <div className='col-span-2'>{menu.permission}</div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>关联API</div>
                          <div className='col-span-2 flex flex-col gap-y-2'>
                            {menu?.apis?.map((api) => (
                              <Badge
                                key={api.id}
                                variant='outline'
                                className={cn(callMethodTypes.get(api.method))}
                              >
                                <span>{api.method}</span>
                                <span className='mx-1'>|</span>
                                <span>{api.path}</span>
                              </Badge>
                            ))}
                          </div>
                        </div>
                      </>
                    )}
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>创建时间</div>
                      <div className='col-span-2'>
                        {format(menu.createdAt, 'yyyy-MM-dd HH:mm:ss')}
                      </div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>更新时间</div>
                      <div className='col-span-2'>
                        {format(menu.updatedAt, 'yyyy-MM-dd HH:mm:ss')}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )
          }}
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
