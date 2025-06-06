import { format } from 'date-fns'
import { IconMessage } from '@tabler/icons-react'
import { SysOperationLog } from '~/do-exercise-api'
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
import { CodeBlock, Status } from '~/components/other'
import { callMethodTypes } from '~/features/api/data/data'
import {
  callCodeTypes,
  callIsBotTypes,
  callIsInternalIpTypes,
  callIsMobileTypes,
} from '../data/data'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function OperationLogViewInfoDialog({ open, onOpenChange }: Props) {
  const { currentRow } = useFormDialog<SysOperationLog>()

  if (!currentRow) return null

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
            <IconMessage /> 操作日志详情
          </DrawerTitle>
          <DrawerDescription>
            查看操作详细信息和相关信息。这包括 请求参数、时间戳和执行状态。
          </DrawerDescription>
        </DrawerHeader>

        <div className='grid gap-6 overflow-y-auto px-4'>
          <div className='space-y-3'>
            <h4 className='text-lg font-medium'>基本信息</h4>
            <div className='grid gap-3 text-sm'>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>请求用户</div>
                <div className='col-span-2'>{currentRow.username}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>IP地址</div>
                <div className='col-span-2'>{currentRow.username}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>是否内网IP</div>
                <div className='col-span-2'>
                  <Badge
                    variant='outline'
                    className={cn(
                      'capitalize',
                      callIsInternalIpTypes.get(currentRow.isInternalIp)
                    )}
                  >
                    {currentRow.isInternalIp ? '是' : '否'}
                  </Badge>
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>IP属地</div>
                <div className='col-span-2'>{currentRow.address || '-'}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>系统</div>
                <div className='col-span-2'>{currentRow.os}</div>
              </div>
              <div className='w-full space-y-3 overflow-hidden'>
                <h4 className='font-medium'>代理</h4>
                <div className='bg-muted w-full rounded-lg p-3'>
                  <CodeBlock content={currentRow?.userAgent} />
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>浏览器</div>
                <div className='col-span-2'>{currentRow.browser}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>浏览器版本</div>
                <div className='col-span-2'>{currentRow.browserVersion}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>是否是移动端</div>
                <div className='col-span-2'>
                  <Badge
                    variant='outline'
                    className={cn(
                      'capitalize',
                      callIsMobileTypes.get(currentRow.isMobile)
                    )}
                  >
                    {currentRow.isMobile ? '是' : '否'}
                  </Badge>
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>是否是机器人</div>
                <div className='col-span-2'>
                  <Badge
                    variant='outline'
                    className={cn(
                      'capitalize',
                      callIsBotTypes.get(currentRow.isBot)
                    )}
                  >
                    {currentRow.isBot ? '是' : '否'}
                  </Badge>
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>请求路径</div>
                <div className='col-span-2'>{currentRow.path}</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>请求方法</div>
                <div className='col-span-2'>
                  <Badge
                    variant='outline'
                    className={cn(
                      'capitalize',
                      callMethodTypes.get(currentRow.method)
                    )}
                  >
                    {currentRow.method}
                  </Badge>
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>是否成功</div>
                <div className='col-span-2'>
                  <Status
                    color={currentRow.success ? 'success' : 'error'}
                    label={currentRow.success ? '成功' : '失败'}
                  />
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>状态码</div>
                <div className='col-span-2'>
                  <Badge
                    variant='outline'
                    className={cn(
                      'capitalize',
                      callCodeTypes.get(currentRow.code)
                    )}
                  >
                    {currentRow.code}
                  </Badge>
                </div>
              </div>
            </div>
          </div>

          <div className='space-y-3'>
            <h4 className='text-lg font-medium'>错误信息</h4>
            {currentRow?.message ? (
              <div className='bg-muted rounded-lg p-3'>
                <CodeBlock content={currentRow?.message ?? ''} />
              </div>
            ) : (
              '-'
            )}
          </div>

          <div className='space-y-3'>
            <h4 className='text-lg font-medium'>请求体</h4>
            {currentRow?.body ? (
              <div className='bg-muted rounded-lg p-3'>
                <CodeBlock language='json' content={currentRow?.body ?? ''} />
              </div>
            ) : (
              '-'
            )}
          </div>

          <div className='space-y-3'>
            <h4 className='text-lg font-medium'>请求参数</h4>
            {currentRow?.params ? (
              <div className='bg-muted rounded-lg p-3'>
                <CodeBlock content={currentRow?.params ?? ''} />
              </div>
            ) : (
              '-'
            )}
          </div>

          <div className='space-y-3'>
            <h4 className='text-lg font-medium'>响应结果</h4>
            {currentRow?.result ? (
              <div className='bg-muted rounded-lg p-3'>
                <CodeBlock language='json' content={currentRow?.result ?? ''} />
              </div>
            ) : (
              '-'
            )}
          </div>

          <div className='space-y-3'>
            <div className='grid gap-3 text-sm'>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>请求开始时间</div>
                <div className='col-span-2'>
                  {format(currentRow.startTime, 'yyyy-MM-dd HH:mm:ss')}
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>请求结束时间</div>
                <div className='col-span-2'>
                  {format(currentRow.endTime, 'yyyy-MM-dd HH:mm:ss')}
                </div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>耗时</div>
                <div className='col-span-2'>{currentRow.latency}ms</div>
              </div>
              <div className='grid grid-cols-3 items-center gap-4'>
                <div className='font-medium'>创建时间</div>
                <div className='col-span-2'>
                  {format(currentRow.endTime, 'yyyy-MM-dd HH:mm:ss')}
                </div>
              </div>
            </div>
          </div>
        </div>

        <DrawerFooter className='gap-y-2'>
          <DrawerClose asChild>
            <Button variant='outline'>关闭</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
