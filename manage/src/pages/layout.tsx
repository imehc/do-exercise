import { Outlet } from 'react-router'
import { AppSidebar, SiteHeader } from '~/components'
import '~/animations.css'
import { SidebarInset, SidebarProvider } from '~/components/ui/sidebar'

const LayoutPage: React.FC = () => {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <SiteHeader />
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  )
}

export default LayoutPage
