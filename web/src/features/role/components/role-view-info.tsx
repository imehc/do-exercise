import { useQuery } from '@tanstack/react-query'
import { IconMessage } from '@tabler/icons-react'
import { SysRole } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
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
import { StatusRenderer, transformData, TreeSelect } from '~/components/other'
import { findMenuTree } from '~/features/menu/data/api'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function RoleViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysRole>()
  const { data = [], isLoading } = useQuery(findMenuTree())

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
            <IconMessage /> 角色详情
          </DrawerTitle>
          <DrawerDescription>查看角色详细信息。</DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={isLoading} data={currentRow}>
          {(role) => {
            return (
              <div className='grid gap-6 overflow-y-auto px-4'>
                <div className='space-y-3'>
                  <h4 className='text-lg font-medium'>基本信息</h4>
                  <div className='grid gap-3 text-sm'>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>角色名称</div>
                      <div className='col-span-2'>{role.name}</div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>角色编码</div>
                      <div className='col-span-2'>{role.code}</div>
                    </div>
                    <div className='grid grid-cols-3 items-center gap-4'>
                      <div className='font-medium'>关联菜单</div>
                      <div className='col-span-2'>
                        <TreeSelect
                          className='col-span-8'
                          mode='view'
                          multiple
                          readonly
                          data={transformData(
                            data,
                            (item) => item.name,
                            (item) => item.id
                          )}
                          value={role.menus.map((item) => item.id) ?? []}
                        />
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
