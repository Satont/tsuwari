import { Group, Avatar, Text, Button, Flex, Alert } from '@mantine/core';
import { IconLogout, IconLogin } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { IntegrationCard } from './card';


import { useStreamlabs } from '@/services/api/integrations';

export const StreamlabsIntegration: React.FC = () => {
  const manager = useStreamlabs();
  const { data } = manager.getIntegration;
  const { t } = useTranslation('integrations');

  async function login() {
    const { data } = manager.getAuthLink;
    if (data) {
      window.location.replace(data);
    }
  }

  return (
    <IntegrationCard
      title="Streamlabs"
      header={
        <Flex direction="row" gap="sm">
          {data && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={() => manager.logout.mutate()}
            >
              {t('logout')}
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            {t('login')}
          </Button>
        </Flex>
      }
    >
      {!data && <Alert>{t('notLoggedIn')}</Alert>}
      {data && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {data.name}
          </Text>
          <Avatar src={data.avatar} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}
    </IntegrationCard>
  );
};
