import { useQuery } from '@tanstack/react-query'
import { IconMessage } from '@tabler/icons-react'
import { i18n } from '@lingui/core'
import { Trans } from '@lingui/react/macro'
import { MenuType, SysMenuTree, SystemMenuApi } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
import { toIconComponent } from '~/utils/icon'
import { getMenuLabel } from '~/utils/menu-label'
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
import { StatusRenderer } from '~/components/other'
import { callMethodTypes } from '~/features/api/data/data'
import {
  getCallMenuMapping,
  callMenuTypes,
  callVisibleTypes,
} from '../data/data'
import { getMenuScopeLabel } from '../schemas/action-schema'

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
      setBackgroundColorOnScale
      open={open}
      onOpenChange={(state) => {
        onOpenChange(state)
      }}
      direction='bottom'
    >
      <DrawerContent className='max-h-[85vh]'>
        <DrawerHeader className='text-left'>
          <DrawerTitle className='flex items-center gap-2'>
            <IconMessage /> <Trans>菜单详情</Trans>
          </DrawerTitle>
          <DrawerDescription>
            <Trans>查看菜单详细信息。</Trans>
          </DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoadingMenu} data={data}>
          {(menu) => {
            const SelectedIcon =
              menu.icon && menu.type === MenuType.menu
                ? toIconComponent(menu.icon)
                : null
            return (
              <div className='grid gap-6 overflow-y-auto px-4'>
                <div className='space-y-3'>
                  <h4 className='text-lg font-medium'>
                    <Trans>基本信息</Trans>
                  </h4>
                  <div className='grid gap-3 text-sm'>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>
                        <Trans>菜单名称</Trans>
                      </div>
                      <div className='col-span-2'>{getMenuLabel(menu)}</div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>
                        <Trans>国际化键</Trans>
                      </div>
                      <div className='col-span-2'>{menu.i18nKey || '-'}</div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>
                        <Trans>菜单类型</Trans>
                      </div>
                      <div className='col-span-2'>
                        <Badge
                          variant='outline'
                          className={cn(
                            'capitalize',
                            callMenuTypes.get(menu.type)
                          )}
                        >
                          {getCallMenuMapping().get(menu.type) ?? '-'}
                        </Badge>
                      </div>
                    </div>
                    {menu.type === MenuType.menu && (
                      <>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>
                            <Trans>路由</Trans>
                          </div>
                          <div className='col-span-2'>{menu.route}</div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>
                            <Trans>图标</Trans>
                          </div>
                          <div className='col-span-2'>
                            {SelectedIcon && <SelectedIcon />}
                          </div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>
                            <Trans>是否可见</Trans>
                          </div>
                          <div className='col-span-2'>
                            <Badge
                              variant='outline'
                              className={cn(
                                'capitalize',
                                callVisibleTypes.get(menu.visible ?? false)
                              )}
                            >
                              {menu.visible ? (
                                <Trans>是</Trans>
                              ) : (
                                <Trans>否</Trans>
                              )}
                            </Badge>
                          </div>
                        </div>
                      </>
                    )}
                    {menu.type !== MenuType.button && (
                      <div className='grid grid-cols-3 items-center gap-4'>
                        <div className='font-medium'>
                          <Trans>序号</Trans>
                        </div>
                        <div className='col-span-2'>{menu.sort}</div>
                      </div>
                    )}
                    {menu.type === MenuType.button && (
                      <>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>
                            <Trans>权限标识</Trans>
                          </div>
                          <div className='col-span-2'>{menu.permission}</div>
                        </div>
                        <div className='grid grid-cols-3 items-center gap-4'>
                          <div className='font-medium'>
                            <Trans>关联接口</Trans>
                          </div>
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
                      <div className='font-medium'>
                        <Trans>可见范围</Trans>
                      </div>
                      <div className='col-span-2 flex items-center gap-2'>
                        <Badge variant='outline'>
                          {getMenuScopeLabel(menu.scope ?? 'both')}
                        </Badge>
                        {menu.isSystem && (
                          <Badge variant='outline'>
                            <Trans>系统内置</Trans>
                          </Badge>
                        )}
                      </div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>
                        <Trans>创建时间</Trans>
                      </div>
                      <div className='col-span-2'>
                        {i18n.date(menu.createdAt, {
                          dateStyle: 'short',
                          timeStyle: 'medium',
                        })}
                      </div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>
                        <Trans>更新时间</Trans>
                      </div>
                      <div className='col-span-2'>
                        {i18n.date(menu.updatedAt, {
                          dateStyle: 'short',
                          timeStyle: 'medium',
                        })}
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
            <Button variant='outline'>
              <Trans>关闭</Trans>
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
