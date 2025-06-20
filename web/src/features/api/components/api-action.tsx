'use client'

import { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { SwitchThumb } from '@radix-ui/react-switch'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { SysApi, SystemApiApi, UpdateApiRequest } from '~/do-exercise-api'
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
import { Switch } from '~/components/ui/switch'
import { DialogHeaderContent } from '~/components/other/dialog-header-content'
import {
  ApiActionFormValues,
  getApiActionSchema,
} from '../schemas/action-schema'

interface Props {
  currentRow?: SysApi
  open: boolean
  onOpenChange: (open: boolean, hasRefresh: boolean) => void
}

export function ApiActionDialog({ currentRow, open, onOpenChange }: Props) {
  const isEdit = useMemo(() => !!currentRow, [currentRow])
  const form = useChan(
    useForm<ApiActionFormValues>({
      resolver: zodResolver(getApiActionSchema()),
      defaultValues: isEdit
        ? currentRow
        : {
            description: '',
            group: '',
            disabled: false,
            sort: 0,
          },
    })
  )

  const sysApi = useApi(SystemApiApi)
  const { isPending, mutate: saveChange } = useMutation({
    mutationFn: (values: UpdateApiRequest) => sysApi.updateApi(values),
    onSuccess: () => {
      toast.success(t`更新成功`)
      form.reset()
      onOpenChange(false, true)
    },
  })

  const onSubmit = (values: ApiActionFormValues) => {
    if (isEdit) {
      saveChange({ id: currentRow!.id, updateSysApi: values })
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        form.reset()
        onOpenChange(state, false)
      }}
    >
      <DialogContent className='sm:max-w-lg'>
        <DialogHeader className='text-left'>
          <DialogHeaderContent isEdit={isEdit} text={<Trans>接口</Trans>} />
        </DialogHeader>
        <div className='-mr-4 h-[26.25rem] w-full overflow-y-auto py-1 pr-4'>
          <Form {...form}>
            <form
              id='api-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-4 p-0.5'
            >
              <FormField
                control={form.control}
                name='path'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>请求路径</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        disabled
                        readOnly
                        placeholder={t`请输入请求路径`}
                        className='col-span-4'
                        autoComplete='off'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='method'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>请求方法</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        disabled
                        readOnly
                        placeholder={t`请输入请求方法`}
                        className='col-span-4'
                        autoComplete='off'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>描述</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        disabled={isPending}
                        placeholder={t`请输入描述`}
                        className='col-span-4'
                        autoComplete='off'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='group'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>分组</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        disabled={isPending}
                        placeholder={t`请填写分组名称`}
                        className='col-span-4'
                        autoComplete='off'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='disabled'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>禁用状态</Trans>
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
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='sort'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-right'>
                      <Trans>排序</Trans>
                    </FormLabel>
                    <FormControl>
                      <Input
                        disabled={isPending}
                        placeholder={t`请输入排序值`}
                        className='col-span-4'
                        type='number'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </div>
        <DialogFooter>
          <Button type='submit' form='api-form' disabled={isPending}>
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
