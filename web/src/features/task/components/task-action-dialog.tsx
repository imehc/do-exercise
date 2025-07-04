import { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { SwitchThumb } from '@radix-ui/react-switch'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import {
  CreateSysJobRequest,
  SysJob,
  SystemJobApi,
  UpdateSysJobRequest,
  JobStatus,
} from '~/do-exercise-api'
import { useApi } from '~/hooks/use-api'
import { useChan } from '~/hooks/use-chan'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
} from '~/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { Switch } from '~/components/ui/switch'
import { Textarea } from '~/components/ui/textarea'
import { DialogHeaderContent } from '~/components/other/dialog-header-content'
import { ActionSysJobFormValues, getSchema } from '../schemas/action-schema'

interface Props {
  currentRow?: SysJob
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function TaskActionDialog({ currentRow, open, onOpenChange }: Props) {
  const systemJobApi = useApi(SystemJobApi)

  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useChan(
    useForm<ActionSysJobFormValues>({
      defaultValues: isEdit
        ? {
          ...currentRow,
          isEdit,
        }
        : {
          isEdit,
          name: '',
          jobGroup: '',
          cronExpression: '',
          command: '',
          status: JobStatus.normal,
          description: '',
          concurrent: false,
          retryTimes: undefined,
          retryInterval: undefined,
          timeout: undefined,
        },
      resolver: zodResolver(getSchema()),
    })
  )

  const { isPending: isPendingCreate, mutate: saveCreate } = useMutation({
    mutationFn: (values: CreateSysJobRequest) =>
      systemJobApi.createSysJob(values),
    onSuccess: () => {
      toast.success(t`创建成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const { isPending: isPendingUpdate, mutate: saveUpdate } = useMutation({
    mutationFn: (values: UpdateSysJobRequest) =>
      systemJobApi.updateSysJob(values),
    onSuccess: () => {
      toast.success(t`更新成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const onSubmit = (values: ActionSysJobFormValues) => {
    if (form.getValues('isEdit')) {
      saveUpdate({
        updateSysJob: {
          name: values.name,
          jobGroup: values.jobGroup,
          cronExpression: values.cronExpression,
          command: values.command,
          status: values.status as JobStatus,
          description: values.description,
          concurrent: values.concurrent,
          retryTimes: values.retryTimes,
          retryInterval: values.retryInterval,
          timeout: values.timeout,
          id: currentRow?.id as number,
        },
        id: currentRow?.id as number,
      })
      return
    }
    saveCreate({
      createSysJob: {
        name: values.name,
        jobGroup: values.jobGroup,
        cronExpression: values.cronExpression,
        command: values.command,
        status: values.status as JobStatus,
        description: values.description,
        concurrent: values.concurrent,
        retryTimes: values.retryTimes,
        retryInterval: values.retryInterval,
        timeout: values.timeout,
      },
    })
  }

  const isPending = isPendingCreate || isPendingUpdate

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        form.reset()
        onOpenChange(state, false)
      }}
    >
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader className='text-left'>
          <DialogHeaderContent isEdit={isEdit} text={<Trans>定时任务</Trans>} />
        </DialogHeader>
        <Form {...form}>
          <form
            id='menu-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 p-0.5'
          >
            <FormField
              control={form.control}
              name='name'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>任务名称</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入任务名称`}
                      className='col-span-7'
                      autoComplete='off'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='jobGroup'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>任务分组</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入任务分组`}
                      className='col-span-7'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='cronExpression'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>cron表达式</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t`请输入cron表达式`}
                      className='col-span-7'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='command'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>执行命令</Trans>
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t`请输入执行命令`}
                      className='col-span-7'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='status'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>状态</Trans>
                  </FormLabel>
                  <Select
                    onValueChange={(v) => field.onChange(+v as JobStatus)}
                    defaultValue={field.value?.toString()}
                  >
                    <FormControl className='col-span-7 w-full'>
                      <SelectTrigger>
                        <SelectValue placeholder='请选择状态' />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.values(JobStatus).map((status) => (
                        <SelectItem key={status} value={status.toString()}>
                          {status === JobStatus.normal ? (
                            <Trans>正常</Trans>
                          ) : (
                            <Trans>暂停</Trans>
                          )}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='description'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>描述</Trans>
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t`请输入描述`}
                      className='col-span-7'
                      autoComplete='off'
                      {...field}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='concurrent'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>并发执行</Trans>
                  </FormLabel>
                  <FormControl>
                    <Switch
                      disabled={isPending}
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    >
                      <SwitchThumb />
                    </Switch>
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='retryTimes'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>重试次数</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      placeholder={t`请输入重试次数`}
                      className='col-span-7'
                      {...field}
                      onChange={(v) => field.onChange(+v.target.value)}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='retryInterval'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>重试间隔(秒)</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      placeholder={t`请输入重试间隔(秒)`}
                      className='col-span-7'
                      {...field}
                      onChange={(v) => field.onChange(+v.target.value)}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='timeout'
              render={({ field }) => (
                <FormItem className='grid grid-cols-10 items-center space-y-0 gap-x-4 gap-y-1'>
                  <FormLabel className='col-span-3 text-right'>
                    <Trans>超时时间(秒)</Trans>
                  </FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      placeholder={t`请输入超时时间(秒)`}
                      className='col-span-7'
                      {...field}
                      onChange={(v) => field.onChange(+v.target.value)}
                      disabled={isPending}
                    />
                  </FormControl>
                  <FormMessage className='col-span-7 col-start-4' />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <DialogFooter>
          <Button type='submit' form='menu-form' disabled={isPending}>
            {isPending ? (
              <>
                <IconLoader3 className='animate-spin' />
                <span>
                  <Trans>保存中</Trans>...
                </span>
              </>
            ) : (
              <span>
                <Trans>保存</Trans>
              </span>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
