import { AppSidebar, SiteHeader, SiteTabs } from '~/components'
import '~/animations.css'
import { SidebarInset, SidebarProvider } from '~/components/ui/sidebar'

const LayoutPage: React.FC = () => {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <SiteHeader />
        <SiteTabs />
      </SidebarInset>
    </SidebarProvider>
  )
}

export default LayoutPage
