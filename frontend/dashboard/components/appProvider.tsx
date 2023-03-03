import { AppShell, ColorScheme, useMantineTheme } from '@mantine/core';
import { useContext, useEffect, useState } from 'react';

import { NavBar } from '@/components/layout/navbar';
import { SideBar } from '@/components/layout/sidebar';
import { FetcherError, useDashboards, useProfile } from '@/services/api';
import { useObsModule } from '@/services/api/modules';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type Props = React.PropsWithChildren<{
  colorScheme: ColorScheme;
}>;

export const AppProvider: React.FC<Props> = (props) => {
  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const dashboardContext = useContext(SelectedDashboardContext);
  const { error: profileError, data: profileData } = useProfile();
  const { data: dashboards } = useDashboards();

  useEffect(() => {
    if (!profileData || !dashboards) return;
    if (!dashboardContext.id) {
      dashboardContext.setId(profileData.id);
    } else {
      const selectedDashboard = dashboards.find((d) => d.id === dashboardContext.id);
      if (!selectedDashboard) {
        dashboardContext.setId(profileData.id);
      }
    }
  }, [profileData, dashboards]);

  useEffect(() => {
    if (profileError) {
      if (profileError instanceof FetcherError && profileError.status === 403) {
        window.location.replace(`/api/auth?state=${window.btoa(window.location.origin)}`);
      } else {
        window.location.replace(`${window.location.origin}`);
      }
    }
  }, [profileError]);

  const theme = useMantineTheme();
  const [sidebarOpened, setSidebarOpened] = useState(false);

  return (
    <AppShell
      styles={{
        main: {
          background: props.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
          width: '100%',
        },
      }}
      navbar={<SideBar opened={sidebarOpened} setOpened={setSidebarOpened} />}
      header={<NavBar setOpened={setSidebarOpened} opened={sidebarOpened} />}
    >
        {props.children}
    </AppShell>
  );
};
