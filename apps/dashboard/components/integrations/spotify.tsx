import { Group, Avatar, Text, Button, Flex } from '@mantine/core';
import { IconBrandSpotify, IconLogin, IconLogout } from '@tabler/icons';

import { IntegrationCard } from './card';

export const SpotifyIntegration: React.FC = () => {
  return (
    <IntegrationCard
      title="Spotify"
      icon={IconBrandSpotify}
      iconColor="green"
      header={
        <Flex direction="row" gap="sm">
          <Button compact leftIcon={<IconLogout />} variant="outline" color="red">
            Logout
          </Button>
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green">
            Login
          </Button>
        </Flex>
      }
    >
      <Group position="apart" mt={10}>
        <Text weight={500} size={30}>
          Satont WorldWide
        </Text>
        <Avatar
          src={
            'https://images.unsplash.com/photo-1527004013197-933c4bb611b3?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=720&q=80'
          }
          h={150}
          w={150}
          style={{ borderRadius: 900 }}
        />
      </Group>
    </IntegrationCard>
  );
};
