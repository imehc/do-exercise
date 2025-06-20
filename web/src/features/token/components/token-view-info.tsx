import { format } from 'date-fns'
import { IconMessage } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
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
            <IconMessage /> <Trans>令牌详情</Trans>
          </DrawerTitle>
          <DrawerDescription>
            <Trans>查看令牌详细信息。</Trans>
          </DrawerDescription>
        </DrawerHeader>

        <StatusRenderer data={currentRow}>
          {(token) => (
            <div className='grid gap-6 overflow-y-auto px-4'>
              <div className='space-y-3'>
                <h4 className='text-lg font-medium'>
                  <Trans>基本信息</Trans>
                </h4>
                <div className='grid gap-3 text-sm'>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>用户ID</Trans>
                    </div>
                    <div className='col-span-2'>{token.userId}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>用户名</Trans>
                    </div>
                    <div className='col-span-2'>{token.username}</div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>令牌</Trans>
                    </div>
                    <div className='col-span-2'>
                      <InlineCopy
                        className='w-fit text-nowrap'
                        text={token.accessToken}
                      />
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>刷新令牌</Trans>
                    </div>
                    <div className='col-span-2'>
                      <InlineCopy
                        className='w-fit text-nowrap'
                        text={token.refreshToken}
                      />
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>创建时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {format(token.createdAt, 'yyyy-MM-dd HH:mm:ss')}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>到期时间</Trans>
                    </div>
                    <div className='col-span-2'>
                      {format(token.expiredAt, 'yyyy-MM-dd HH:mm:ss')}
                    </div>
                  </div>
                  <div className='grid grid-cols-3 items-center gap-4'>
                    <div className='font-medium'>
                      <Trans>禁用状态</Trans>
                    </div>
                    <div className='col-span-2'>
                      <Badge
                        variant='outline'
                        className={cn(
                          'capitalize',
                          callDisabledTypes.get(token.disabled)
                        )}
                      >
                        {!token.disabled ? (
                          <Trans>正常</Trans>
                        ) : (
                          <Trans>已禁用</Trans>
                        )}
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
            <Button variant='outline'>
              <Trans>关闭</Trans>
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
