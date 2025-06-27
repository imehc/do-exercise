import { useEffect, useRef, useState } from 'react'
import { t } from '@lingui/core/macro'
import { from, mergeMap, map, iif, of, tap, finalize } from 'rxjs'
import { fromFetch } from 'rxjs/fetch'
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
        mergeMap((file) =>
          iif(
            () =>
              !!(
                stableRef.current.options?.preview && stableRef.current.onChange
              ),
            from(toBase64(file)).pipe(
              tap((base64) =>
                of(stableRef.current.onChange?.(base64 as string))
              ),
              map(() => file)
            ),
            of(file)
          ).pipe(tap(() => setIsPending(true)))
        ),
        mergeMap((file) =>
          from(
            stableRef.current.ossApi.getPresignedUrl({ fileName: file.name })
          ).pipe(map((oss) => ({ file, oss })))
        ),
        mergeMap(({ file, oss }) =>
          fromFetch(oss.putObjectUrl, {
            method: 'PUT',
            body: file,
          }).pipe(
            map((response) => {
              if (!(response.ok && response.status === 200)) {
                throw new Error('upload failed')
              }
              return oss.getObjectUrl
            })
          )
        ),
        finalize(() => {
          setIsPending(false)
          if (stableRef.current.options?.preview) {
            stableRef.current.onChange?.(undefined)
          }
        })
      )
      .subscribe({
        next: (url) => {
          setIsPending(false)
          stableRef.current.onChange?.(url)
        },
        error: () => {
          toast.error(t`上传失败`)
        },
      })

    return () => task.unsubscribe()
  }, [upload$])

  return { isPending, upload: handleUpload }
}
