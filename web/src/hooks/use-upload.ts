import { useEffect, useRef, useState } from 'react'
import { t } from '@lingui/core/macro'
import { EMPTY, catchError, concatMap, defer, finalize } from 'rxjs'
import { toast } from 'sonner'
import { OssApi } from '~/do-exercise-api'
import { toBase64 } from '~/lib/utils'
import { useApi } from '~/hooks/use-api'
import { useSubject } from './use-rx'

export function useAvatarUploadSubject(
  onChange?: (url?: string) => void,
  options?: { preview?: boolean }
) {
  const [isPending, setIsPending] = useState(false)
  const [upload$, handleUpload] = useSubject<File>()
  const ossApi = useApi(OssApi)
  const stableRef = useRef({ onChange, options, ossApi })
  useEffect(() => {
    stableRef.current = { onChange, options, ossApi }
  }, [onChange, options, ossApi])

  useEffect(() => {
    const task = upload$
      .pipe(
        concatMap((file) => {
          const controller = new AbortController()
          return defer(async () => {
            setIsPending(true)
            if (stableRef.current.options?.preview) {
              const preview = await toBase64(file)
              if (!controller.signal.aborted)
                stableRef.current.onChange?.(preview as string)
            }
            const signed = await stableRef.current.ossApi.getPresignedUrl(
              { fileName: file.name, fileSize: file.size, purpose: 'avatar' },
              { signal: controller.signal }
            )
            const body = new FormData()
            Object.entries(signed.formFields).forEach(([key, value]) =>
              body.append(key, value)
            )
            body.append('file', file)
            const response = await fetch(signed.uploadUrl, {
              method: 'POST',
              body,
              signal: controller.signal,
            })
            if (!response.ok || !signed.getObjectUrl)
              throw new Error('upload failed')
            const completed = await stableRef.current.ossApi.completeUpload(
              { ossObject: { objectKey: signed.objectKey } },
              { signal: controller.signal }
            )
            if (!completed.getObjectUrl) throw new Error('upload incomplete')
            return completed.getObjectUrl
          }).pipe(
            catchError(() => {
              if (!controller.signal.aborted) {
                toast.error(t`上传失败`)
                if (stableRef.current.options?.preview)
                  stableRef.current.onChange?.(undefined)
              }
              return EMPTY
            }),
            finalize(() => {
              controller.abort()
              setIsPending(false)
            })
          )
        })
      )
      .subscribe((url) => stableRef.current.onChange?.(url))
    return () => task.unsubscribe()
  }, [upload$])
  return { isPending, upload: handleUpload }
}
