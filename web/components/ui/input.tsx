import { Icon } from "@iconify/react"
import clsx from "clsx"
import * as React from "react"

import { cn } from "~/lib/utils"

const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<"input"> & { startIcon?: string }>(
  ({ className, type, startIcon, ...props }, ref) => {
    return (
      <div className="relative flex items-center max-w-2xl">
        {startIcon && (
          <Icon
            className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 transform"
            icon={startIcon}
          />
        )}
        <input
          type={type}
          className={cn(
            clsx("flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",{ "pl-8": startIcon }),
            className
          )}
          ref={ref}
          {...props}
        />
      </div>
    )
  }
)
Input.displayName = "Input"

export { Input }
