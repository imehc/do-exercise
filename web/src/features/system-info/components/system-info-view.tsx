import React from 'react'
import { Trans } from '@lingui/react/macro'
import { SysInfo } from '~/do-exercise-api'
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card'
import { Progress } from '~/components/ui/progress'

interface SystemDashboardProps {
  data: SysInfo
}

export const SysInfoView: React.FC<SystemDashboardProps> = ({
  data: sysData,
}) => {
  const memoryUsagePercent = (sysData.ram.used / sysData.ram.total) * 100

  // CPU数据转换
  const cpuData = sysData.cpu.cpus.map((usage, index) => ({
    index: index + 1,
    usageRate: usage,
  }))

  return (
    <div className='grid grid-cols-1 gap-4 p-4 md:grid-cols-2 lg:grid-cols-3'>
      {/* 操作系统信息卡片 */}
      <Card className='bg-background col-span-full rounded-lg shadow-sm'>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle className='text-lg font-medium'>
            <Trans>系统概览</Trans>
          </CardTitle>
        </CardHeader>
        <CardContent className='grid grid-cols-2 gap-4 md:grid-cols-5'>
          <div className='space-y-1'>
            <p className='text-muted-foreground text-sm'>
              <Trans>操作系统</Trans>
            </p>
            <p className='font-medium'>{sysData.os.goos}</p>
          </div>
          <div className='space-y-1'>
            <p className='text-muted-foreground text-sm'>
              <Trans>Go版本</Trans>
            </p>
            <p className='font-medium'>{sysData.os.goVersion}</p>
          </div>
          <div className='space-y-1'>
            <p className='text-muted-foreground text-sm'>
              <Trans>协程数</Trans>
            </p>
            <p className='font-medium'>{sysData.os.numGoroutine}</p>
          </div>
          <div className='space-y-1'>
            <p className='text-muted-foreground text-sm'>
              <Trans>编译器</Trans>
            </p>
            <p className='font-medium'>{sysData.os.compiler}</p>
          </div>
          <div className='space-y-1'>
            <p className='text-muted-foreground text-sm'>
              <Trans>CPU核心</Trans>
            </p>
            <p className='font-medium'>{sysData.cpu.cores}</p>
          </div>
        </CardContent>
      </Card>

      {/* CPU使用率卡片 */}
      <Card className='bg-background col-span-full rounded-lg shadow-sm md:col-span-1'>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle className='text-lg font-medium'>
            <Trans>CPU使用率</Trans>
          </CardTitle>
        </CardHeader>
        <CardContent className='space-y-4'>
          {cpuData.map((core) => (
            <div key={core.index} className='space-y-2'>
              <div className='flex justify-between text-sm'>
                <span className='text-muted-foreground'>
                  <Trans>核心 {core.index}</Trans>
                </span>
                <span className='font-medium'>
                  {core.usageRate.toFixed(1)}%
                </span>
              </div>
              <Progress value={core.usageRate} className='h-2' />
            </div>
          ))}
        </CardContent>
      </Card>

      {/* 内存使用情况卡片 */}
      <Card className='bg-background col-span-full rounded-lg shadow-sm md:col-span-1'>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle className='text-lg font-medium'>
            <Trans>内存使用情况</Trans>
          </CardTitle>
        </CardHeader>
        <CardContent className='space-y-4'>
          <div className='space-y-2'>
            <div className='flex justify-between text-sm'>
              <span className='text-muted-foreground'>
                <Trans>总内存</Trans>
              </span>
              <span className='font-medium'>{sysData.ram.total} MB</span>
            </div>
            <div className='flex justify-between text-sm'>
              <span className='text-muted-foreground'>
                <Trans>已使用</Trans>
              </span>
              <span className='font-medium'>{sysData.ram.used} MB</span>
            </div>
            <div className='flex justify-between text-sm'>
              <span className='text-muted-foreground'>
                <Trans>使用率</Trans>
              </span>
              <span className='font-medium'>
                {memoryUsagePercent.toFixed(1)}%
              </span>
            </div>
            <Progress value={memoryUsagePercent} className='h-2' />
          </div>
        </CardContent>
      </Card>

      {/* 磁盘使用情况卡片 */}
      <Card className='bg-background col-span-full rounded-lg shadow-sm md:col-span-1'>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle className='text-lg font-medium'>
            <Trans>磁盘使用情况</Trans>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex flex-col gap-y-3'>
            {sysData.disks.map((disk) => {
              const usagePercent = (disk.used / disk.total) * 100
              return (
                <div key={disk.mountPoint} className='space-y-2'>
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>
                      <Trans>挂载点</Trans>
                    </span>
                    <span className='font-medium'>{disk.mountPoint}</span>
                  </div>
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>
                      <Trans>总容量</Trans>
                    </span>
                    <span className='font-medium'>{disk.total} MB</span>
                  </div>
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>
                      <Trans>已使用</Trans>
                    </span>
                    <span className='font-medium'>{disk.used} MB</span>
                  </div>
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>
                      <Trans>使用率</Trans>
                    </span>
                    <span className='font-medium'>
                      {usagePercent.toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={usagePercent} className='h-2' />
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
