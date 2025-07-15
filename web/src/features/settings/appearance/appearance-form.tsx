import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { IconLoader3 } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans } from '@lingui/react/macro'
import { toast } from 'sonner'
import { fonts } from '~/config/fonts'
import { useFont } from '~/provider/font'
import { useTheme } from '~/provider/theme'
import { useChan } from '~/hooks/use-chan'
import { Button } from '~/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { RadioGroup, RadioGroupItem } from '~/components/ui/radio-group'
import { SelectDropdown } from '~/components/select-dropdown'

const getAppearanceFormSchema = () =>
  z.object({
    theme: z.enum(['light', 'dark'], {
      error: t`请选择一个主题`,
    }),
    font: z.enum(fonts, {
      error: t`选择一种字体`,
    }),
  })

type AppearanceFormValues = z.infer<ReturnType<typeof getAppearanceFormSchema>>

export function AppearanceForm() {
  const { font, setFont } = useFont()
  const { theme, setTheme } = useTheme()
  const [loading, setLoading] = useState(false)

  // This can come from your database or API.
  const defaultValues: Partial<AppearanceFormValues> = {
    theme: theme as 'light' | 'dark',
    font,
  }

  const form = useChan(
    useForm<AppearanceFormValues>({
      resolver: zodResolver(getAppearanceFormSchema()),
      defaultValues,
    })
  )

  function onSubmit(data: AppearanceFormValues) {
    setLoading(true)
    if (data.font != font) setFont(data.font)
    if (data.theme != theme) setTheme(data.theme)
    setLoading(false)
    toast.success(t`已成功更新外观设置`)
    // showSubmittedData(data)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-8'>
        <FormField
          control={form.control}
          name='font'
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                <Trans>字体</Trans>
              </FormLabel>
              <FormControl>
                <SelectDropdown
                  defaultValue={field.value ?? ''}
                  onValueChange={field.onChange}
                  className='w-full'
                  items={fonts
                    .map((item) => ({
                      label: item,
                      value: item,
                    }))
                    .flat()
                    .sort((a, b) => a.label?.localeCompare(b.label))}
                />
              </FormControl>
              <FormDescription className='font-manrope'>
                <Trans>设定您想要在系统中使用的字体。</Trans>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='theme'
          render={({ field }) => (
            <FormItem className='space-y-1'>
              <FormLabel>
                <Trans>主题</Trans>
              </FormLabel>
              <FormDescription>
                <Trans>选择您想要在系统中使用的主题。</Trans>
              </FormDescription>
              <FormMessage />
              <RadioGroup
                onValueChange={field.onChange}
                defaultValue={field.value}
                className='grid max-w-md grid-cols-2 gap-8 pt-2 max-sm:flex max-sm:flex-col max-sm:space-y-1'
              >
                <FormItem>
                  <FormLabel className='[&:has([data-state=checked])>div]:border-primary'>
                    <FormControl>
                      <RadioGroupItem value='light' className='sr-only' />
                    </FormControl>
                    <div className='border-muted hover:border-accent items-center rounded-md border-2 p-1'>
                      <div className='space-y-2 rounded-sm bg-[#ecedef] p-2'>
                        <div className='space-y-2 rounded-md bg-white p-2 shadow-xs'>
                          <div className='h-2 w-[80px] rounded-lg bg-[#ecedef]' />
                          <div className='h-2 w-[100px] rounded-lg bg-[#ecedef]' />
                        </div>
                        <div className='flex items-center space-x-2 rounded-md bg-white p-2 shadow-xs'>
                          <div className='h-4 w-4 rounded-full bg-[#ecedef]' />
                          <div className='h-2 w-[100px] rounded-lg bg-[#ecedef]' />
                        </div>
                        <div className='flex items-center space-x-2 rounded-md bg-white p-2 shadow-xs'>
                          <div className='h-4 w-4 rounded-full bg-[#ecedef]' />
                          <div className='h-2 w-[100px] rounded-lg bg-[#ecedef]' />
                        </div>
                      </div>
                    </div>
                    <span className='block w-full p-2 text-center font-normal'>
                      <Trans>浅色模式</Trans>
                    </span>
                  </FormLabel>
                </FormItem>
                <FormItem>
                  <FormLabel className='[&:has([data-state=checked])>div]:border-primary'>
                    <FormControl>
                      <RadioGroupItem value='dark' className='sr-only' />
                    </FormControl>
                    <div className='border-muted bg-popover hover:bg-accent hover:text-accent-foreground items-center rounded-md border-2 p-1'>
                      <div className='space-y-2 rounded-sm bg-slate-950 p-2'>
                        <div className='space-y-2 rounded-md bg-slate-800 p-2 shadow-xs'>
                          <div className='h-2 w-[80px] rounded-lg bg-slate-400' />
                          <div className='h-2 w-[100px] rounded-lg bg-slate-400' />
                        </div>
                        <div className='flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-xs'>
                          <div className='h-4 w-4 rounded-full bg-slate-400' />
                          <div className='h-2 w-[100px] rounded-lg bg-slate-400' />
                        </div>
                        <div className='flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-xs'>
                          <div className='h-4 w-4 rounded-full bg-slate-400' />
                          <div className='h-2 w-[100px] rounded-lg bg-slate-400' />
                        </div>
                      </div>
                    </div>
                    <span className='block w-full p-2 text-center font-normal'>
                      <Trans>深色模式</Trans>
                    </span>
                  </FormLabel>
                </FormItem>
              </RadioGroup>
            </FormItem>
          )}
        />

        <Button type='submit' className='max-sm:w-full'>
          {loading ? (
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
      </form>
    </Form>
  )
}
