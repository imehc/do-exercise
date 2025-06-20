import * as React from 'react'
import {
  IconCheck,
  IconChevronDown,
  IconCircleX,
  IconSearch,
  IconX,
} from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '~/components/ui/popover'
import { Separator } from '~/components/ui/separator'

/**
 * Variants for the multi-select component to handle different styles.
 * Uses class-variance-authority (cva) to define different styles based on "variant" prop.
 */
const multiSelectVariants = cva(
  'm-1 transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-110 duration-300',
  {
    variants: {
      variant: {
        default:
          'border-foreground/10 text-foreground bg-card hover:bg-card/80',
        secondary:
          'border-foreground/10 bg-secondary text-secondary-foreground hover:bg-secondary/80',
        destructive:
          'border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80',
        inverted: 'inverted',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
)

type Value = string | number

/**
 * Props for MultiSelect component
 * @description https://github.com/sersavan/shadcn-multi-select-component
 */
interface MultiSelectProps
  extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, 'defaultValue'>,
    VariantProps<typeof multiSelectVariants> {
  /**
   * An array of option objects to be displayed in the multi-select component.
   * Each option object has a label, value, and an optional icon.
   */
  options: {
    /** The text to display for the option. */
    label: string
    /** The unique value associated with the option. */
    value: Value
    /** Optional icon component to display alongside the option. */
    icon?: React.ComponentType<{ className?: string }>
  }[]

  /**
   * Callback function triggered when the selected values change.
   * Receives an array of the new selected values.
   */
  onValueChange: (value: Value[]) => void

  /** The default selected values when the component mounts. */
  defaultValue?: Value[]

  /**
   * Placeholder text to be displayed when no values are selected.
   * Optional, defaults to "Select options".
   */
  placeholder?: string

  /**
   * Maximum number of items to display. Extra selected items will be summarized.
   * Optional, defaults to 3.
   */
  maxCount?: number

  /**
   * The modality of the popover. When set to true, interaction with outside elements
   * will be disabled and only popover content will be visible to screen readers.
   * Optional, defaults to false.
   */
  modalPopover?: boolean

  /**
   * If true, renders the multi-select component as a child of another component.
   * Optional, defaults to false.
   */
  asChild?: boolean

  /**
   * Additional class names to apply custom styles to the multi-select component.
   * Optional, can be used to add custom styles.
   */
  className?: string
}

export const MultiSelect = React.forwardRef<
  HTMLButtonElement,
  MultiSelectProps
>(
  (
    {
      options,
      onValueChange,
      variant,
      defaultValue = [],
      placeholder = 'Select options',
      maxCount = 3,
      modalPopover = false,
      // asChild = false,
      className,
      ...props
    },
    ref
  ) => {
    const [selectedValues, setSelectedValues] =
      React.useState<Value[]>(defaultValue)
    const [isPopoverOpen, setIsPopoverOpen] = React.useState(false)
    const [searchQuery, setSearchQuery] = React.useState('')

    const handleInputKeyDown = (
      event: React.KeyboardEvent<HTMLInputElement>
    ) => {
      if (event.key === 'Enter') {
        setIsPopoverOpen(true)
      } else if (event.key === 'Backspace' && !event.currentTarget.value) {
        const newSelectedValues = [...selectedValues]
        newSelectedValues.pop()
        setSelectedValues(newSelectedValues)
        onValueChange(newSelectedValues)
      }
    }

    const toggleOption = (option: Value) => {
      const newSelectedValues = selectedValues.includes(option)
        ? selectedValues.filter((value) => value !== option)
        : [...selectedValues, option]
      setSelectedValues(newSelectedValues)
      onValueChange(newSelectedValues)
    }

    const handleClear = () => {
      setSelectedValues([])
      onValueChange([])
    }

    const handleTogglePopover = () => {
      setIsPopoverOpen((prev) => !prev)
    }

    const clearExtraOptions = () => {
      const newSelectedValues = selectedValues.slice(0, maxCount)
      setSelectedValues(newSelectedValues)
      onValueChange(newSelectedValues)
    }

    const toggleAll = () => {
      if (selectedValues.length === options.length) {
        handleClear()
      } else {
        const allValues = options.map((option) => option.value)
        setSelectedValues(allValues)
        onValueChange(allValues)
      }
    }

    return (
      <Popover
        open={isPopoverOpen}
        onOpenChange={setIsPopoverOpen}
        modal={modalPopover}
      >
        <PopoverTrigger asChild>
          <Button
            ref={ref}
            {...props}
            onClick={handleTogglePopover}
            className={cn(
              'flex h-auto min-h-10 w-full items-center justify-between rounded-md border bg-inherit p-1 hover:bg-inherit [&_svg]:pointer-events-auto',
              className
            )}
          >
            {selectedValues.length > 0 ? (
              <div className='flex w-full items-center justify-between'>
                <div className='flex flex-wrap items-center'>
                  {selectedValues.slice(0, maxCount).map((value) => {
                    const option = options.find((o) => o.value === value)
                    const IconComponent = option?.icon
                    return (
                      <Badge
                        key={value}
                        className={cn(multiSelectVariants({ variant }))}
                      >
                        {IconComponent && (
                          <IconComponent className='mr-2 h-4 w-4' />
                        )}
                        {option?.label}
                        <span
                          onClick={(event) => {
                            event.stopPropagation()
                            toggleOption(value)
                          }}
                        >
                          <IconCircleX className='ml-2 h-4 w-4 cursor-pointer' />
                        </span>
                      </Badge>
                    )
                  })}
                  {selectedValues.length > maxCount && (
                    <Badge
                      className={cn(
                        'text-foreground border-foreground/1 bg-transparent hover:bg-transparent',
                        multiSelectVariants({ variant })
                      )}
                    >
                      {`+ ${selectedValues.length - maxCount} more`}
                      <IconCircleX
                        className='ml-2 h-4 w-4 cursor-pointer'
                        onClick={(event) => {
                          event.stopPropagation()
                          clearExtraOptions()
                        }}
                      />
                    </Badge>
                  )}
                </div>
                <div className='flex items-center justify-between'>
                  <IconX
                    className='text-muted-foreground mx-2 h-4 cursor-pointer'
                    onClick={(event) => {
                      event.stopPropagation()
                      handleClear()
                    }}
                  />
                  <Separator
                    orientation='vertical'
                    className='flex h-full min-h-6'
                  />
                  <IconChevronDown className='text-muted-foreground mx-2 h-4 cursor-pointer' />
                </div>
              </div>
            ) : (
              <div className='mx-auto flex w-full items-center justify-between'>
                <span className='text-muted-foreground mx-2 text-base font-normal md:text-sm'>
                  {placeholder}
                </span>
                <IconChevronDown className='text-muted-foreground mx-2 h-4 cursor-pointer' />
              </div>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent
          className='w-[var(--radix-popover-trigger-width)] p-0'
          align='start'
          onEscapeKeyDown={() => setIsPopoverOpen(false)}
        >
          <div className='flex items-center border-b px-3'>
            <IconSearch className='mr-2 size-4 shrink-0 opacity-50' />
            <input
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={handleInputKeyDown}
              className='placeholder:text-muted-foreground flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50'
              placeholder='Search icons...'
            />
          </div>
          <div className='bg-popover text-popover-foreground pointer-events-auto h-56 !overflow-y-auto rounded-md p-1 shadow-md'>
            <div
              className='hover:bg-accent flex cursor-pointer items-center rounded-sm px-2 py-1.5 text-sm outline-none'
              onClick={toggleAll}
            >
              <div
                className={cn(
                  'border-primary mr-2 flex h-4 w-4 items-center justify-center rounded-sm border',
                  selectedValues.length === options.length
                    ? 'bg-primary text-primary-foreground'
                    : 'opacity-50 [&_svg]:invisible'
                )}
              >
                <IconCheck className='h-4 w-4' />
              </div>
              <Trans>(选择全部)</Trans>
            </div>
            {options
              .filter((option) =>
                option.label.toLowerCase().includes(searchQuery.toLowerCase())
              )
              .map((option) => {
                const isSelected = selectedValues.includes(option.value)
                const IconComponent = option.icon
                return (
                  <div
                    key={option.value}
                    className='hover:bg-accent flex cursor-pointer items-center rounded-sm px-2 py-1.5 text-sm outline-none'
                    onClick={() => toggleOption(option.value)}
                  >
                    <div
                      className={cn(
                        'border-primary mr-2 flex h-4 w-4 items-center justify-center rounded-sm border',
                        isSelected
                          ? 'bg-primary text-primary-foreground'
                          : 'opacity-50 [&_svg]:invisible'
                      )}
                    >
                      <IconCheck className='h-4 w-4' />
                    </div>
                    {IconComponent && (
                      <IconComponent className='text-muted-foreground mr-2 h-4 w-4' />
                    )}
                    <span>{option.label}</span>
                  </div>
                )
              })}
          </div>
        </PopoverContent>
      </Popover>
    )
  }
)

MultiSelect.displayName = 'MultiSelect'
