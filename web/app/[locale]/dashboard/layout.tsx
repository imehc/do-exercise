import { SidebarInset, SidebarProvider } from '~/components/ui/sidebar';
import { AppSidebar } from '~/components/other/app-sidebar';
import { cookies } from 'next/headers';
import { apiInstance } from '~/helper/api';
import { MenuApi } from '~/do-exercise-api';
import { Header } from '~/components/other/header';

type LayoutProps = {
  children: React.ReactNode;
};

export default async function DashboardLayout({ children }: LayoutProps) {
  const cookieStore = await cookies();
  const defaultOpen = cookieStore.get('sidebar_state')?.value === 'true';
  const menuApi = await apiInstance(MenuApi);
  const sideMenus = await menuApi.getMenuTreeCompact();

  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <AppSidebar sideMenus={sideMenus} />
      <SidebarInset>
        <main>
          <Header />
          {children}
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
