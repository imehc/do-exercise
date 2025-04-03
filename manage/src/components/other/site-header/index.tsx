import { Separator } from '~/components/ui/separator'
import { SidebarTrigger } from '~/components/ui/sidebar'
import { ModelToggle } from './model-toggle'

export function SiteHeader() {
  return (
    <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mx-2 data-[orientation=vertical]:h-4" />
        <div className="flex-1 flex justify-between items-center">
          <h1 className="text-base font-medium">Documents</h1>
          <div className="flex justify-end items-center gap-x-4">
            <ModelToggle />
          </div>
        </div>
      </div>
    </header>
  )
}
