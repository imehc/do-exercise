import { useEffect, useMemo, useRef } from 'react'
import { FieldValues, UseFormReturn } from 'react-hook-form'
import { useAtomValue } from 'jotai'
import { languageAtom } from '~/atoms'

/**
 * 用于语言切换后同步语言变化
 * @param form
 * @returns
 */
export function useChan<T extends FieldValues = FieldValues>(
  form: UseFormReturn<T>,
  noVerification = false
): UseFormReturn<T> {
  const locale = useAtomValue(languageAtom)

  const hasError = useMemo(() => {
    return Object.keys(form.formState.errors).length > 0
  }, [form.formState.errors])

  const isFirst = useRef(true)

  useEffect(() => {
    if (!noVerification) return
    if (isFirst.current) {
      isFirst.current = false
      return
    }

    if (hasError) {
      form.clearErrors()
      form.trigger()
    }
  }, [form, hasError, locale, noVerification])

  return form
}
