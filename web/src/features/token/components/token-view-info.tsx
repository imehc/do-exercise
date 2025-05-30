import { format } from 'date-fns'
import { IconMessage } from '@tabler/icons-react'
import { TokenInfo } from '~/do-exercise-api'
import { useFormDialog } from '~/provider'
import { cn } from '~/lib/utils'
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
import { InlineCopy, StatusRenderer } from '~/components/other'
import { callDisabledTypes } from '~/features/api/data/data'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function TokenViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<TokenInfo>()

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
            <IconMessage /> 令牌详情
          </DrawerTitle>
          <DrawerDescription>查看令牌详细信息。</DrawerDescription>
        </DrawerHeader>

        <StatusRenderer data={currentRow}>
          {(token) => (
            <div className='grid gap-6 overflow-y-auto px-4'>
              <div className='space-y-3'>
                <h4 className='text-lg font-medium'>基本信息</h4>
                <div className='grid gap-3 text-sm'>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>用户ID</div>
                    <div className='col-span-2'>{token.userId}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>用户名</div>
                    <div className='col-span-2'>{token.username}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>令牌</div>
                    <div className='col-span-2'>
                      <InlineCopy
                        className='w-fit text-nowrap'
                        text={token.accessToken}
                      />
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>刷新令牌</div>
                    <div className='col-span-2'>
                      <InlineCopy
                        className='w-fit text-nowrap'
                        text={token.refreshToken}
                      />
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>创建时间</div>
                    <div className='col-span-2'>
                      {format(token.createdAt, 'yyyy-MM-dd HH:mm:ss')}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>到期时间</div>
                    <div className='col-span-2'>
                      {format(token.expiredAt, 'yyyy-MM-dd HH:mm:ss')}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>禁用状态</div>
                    <div className='col-span-2'>
                      <Badge
                        variant='outline'
                        className={cn(
                          'capitalize',
                          callDisabledTypes.get(token.disabled)
                        )}
                      >
                        {!token.disabled ? '正常' : '已禁用'}
                      </Badge>
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
