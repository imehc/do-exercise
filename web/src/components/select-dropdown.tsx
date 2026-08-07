import { useMemo } from 'react'
import { IconLoader } from '@tabler/icons-react'
import { cn } from '~/lib/utils'
import { FormControl } from '~/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'

interface SelectDropdownProps {
  onValueChange?: (value: string) => void
  defaultValue: string | undefined
  placeholder?: string
  isPending?: boolean
  items: { label: string; value: string }[] | undefined
  disabled?: boolean
  className?: string
  isControlled?: boolean
}

export function SelectDropdown({
  defaultValue,
  onValueChange,
  isPending,
  items,
  placeholder,
  disabled,
  className = '',
  isControlled = false,
}: SelectDropdownProps) {
  const defaultState = isControlled
    ? { value: defaultValue, onValueChange }
    : { defaultValue, onValueChange }
  // Radix 用 value 作为内部原生 option 的 key，且不接受空串，
  // 因此这里剔除空值并按 value 去重，避免上游脏数据触发 key 警告
  const validItems = useMemo(() => {
    const seen = new Set<string>()
    return (items ?? []).filter(({ value }) => {
      if (!value || seen.has(value)) return false
      seen.add(value)
      return true
    })
  }, [items])
  return (
    <Select {...defaultState}>
      <FormControl>
        <SelectTrigger
          disabled={disabled}
          className={cn(className, 'text-base md:text-sm')}
        >
          <SelectValue placeholder={placeholder ?? 'Select'} />
        </SelectTrigger>
      </FormControl>
      <SelectContent>
        {isPending ? (
          <SelectItem disabled value='loading' className='h-14'>
            <div className='flex items-center justify-center gap-2'>
              <IconLoader className='h-5 w-5 animate-spin' />
              Loading...
            </div>
          </SelectItem>
        ) : validItems.length ? (
          validItems.map(({ label, value }) => (
            <SelectItem key={value} value={value}>
              {label}
            </SelectItem>
          ))
        ) : (
          <SelectItem value='none' disabled>
            <div className='flex items-center justify-center gap-2'>
              No options
            </div>
          </SelectItem>
        )}
      </SelectContent>
    </Select>
  )
}
