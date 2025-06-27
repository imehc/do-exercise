import React from 'react'
import { IconPencil, IconUser } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { toast } from 'sonner'
import { ensureHttpPrefix } from '~/utils/url'
import { useAvatarUploadSubject } from '~/hooks/use-upload'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { Button } from '~/components/ui/button'

type AvatarUploadProps = {
  value?: string
  onChange?: (value?: string) => void
  disabled?: boolean
}

/**
 * @description https://github.com/shadcn-ui/ui/issues/250#issuecomment-1985951964
 */
export function AvatarUpload({
  value,
  onChange,
  disabled = false,
}: AvatarUploadProps) {
  const { isPending, upload } = useAvatarUploadSubject(onChange, {
    preview: true,
  })

  const inputRef = React.useRef<HTMLInputElement>(null)
  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0]
      if (!file) {
        toast.error(t`未选择文件`)
        return
      }
      upload(file)
      // const base64 = (await toBase64(file)) as string
      // onChange?.(base64)
    }
  }

  return (
    <div className='relative h-40 w-40'>
      <Avatar className='h-full w-full'>
        <AvatarImage src={ensureHttpPrefix(value)} className='object-cover' />
        <AvatarFallback className='bg-secondary'>
          <IconUser className='h-16 w-16' />
        </AvatarFallback>
      </Avatar>
      <Button
        disabled={disabled || isPending}
        variant='ghost'
        size='icon'
        className='bg-secondary-foreground/80 hover:bg-secondary-foreground dark:bg-primary/80 dark:hover:bg-primary absolute right-0 bottom-0 rounded-full p-1'
        onClick={(e) => {
          e.preventDefault()
          inputRef.current?.click()
        }}
      >
        <IconPencil className='h-4 w-4 text-white dark:text-black' />
      </Button>
      <input
        ref={inputRef}
        type='file'
        className='hidden'
        onChange={handleChange}
        accept='image/*'
      />
    </div>
  )
}
