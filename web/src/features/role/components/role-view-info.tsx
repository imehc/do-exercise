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
import { StatusRenderer } from '~/components/other'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function RoleViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysRole>()

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
            <IconMessage /> 角色详情
          </DrawerTitle>
          <DrawerDescription>查看角色详细信息。</DrawerDescription>
        </DrawerHeader>

        <StatusRenderer isLoading={false} data={currentRow}>
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
