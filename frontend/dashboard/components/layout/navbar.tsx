import {
  ActionIcon,
  Burger,
  Container,
  createStyles,
  Flex,
  Group,
  Header,
  Loader,
  Menu,
  Text,
  Box, Divider, Button, Badge, Tooltip,
} from '@mantine/core';
import { IconMoonStars, IconSun, IconLanguage, IconNotes } from '@tabler/icons';
import { useStore, useAtom } from 'jotai';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { Dispatch, SetStateAction } from 'react';

import DiscordSvg from '../../public/assets/icons/brands/discord.svg';
import { Profile } from './profile';

import { ChangelogModal } from '@/components/changelog/modal';
import { modalOpenedAtomic } from '@/components/changelog/store';
import { useProfile } from '@/services/api';
import { useLocale, LOCALES } from '@/services/dashboard';
import { useTheme } from '@/services/dashboard';

const useStyles = createStyles((theme) => ({
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    height: '100%',
  },

  hiddenMobile: {
    [theme.fn.smallerThan('sm')]: {
      display: 'none',
    },
  },

  hiddenDesktop: {
    [theme.fn.largerThan('sm')]: {
      display: 'none',
    },
  },
}));

export function NavBar({
  opened,
  setOpened,
}: {
  setOpened: Dispatch<SetStateAction<boolean>>;
  opened: boolean;
}) {
  const store = useStore();
  const router = useRouter();

  const [modalOpened, setModalOpened] = useAtom(modalOpenedAtomic);

  const { classes } = useStyles();
  const { theme, colorScheme, toggleColorScheme } = useTheme();
  const { locale, toggleLocale } = useLocale();
  const { data: userData, isLoading: isLoadingProfile } = useProfile();

  return (
    <Header height={60}>
      <Container maw="unset" className={classes.header}>
        <Flex gap="sm" justify="flex-start" align="center" direction="row">
          <Burger
            className={classes.hiddenDesktop}
            opened={opened}
            onClick={() => setOpened(!opened)}
            size="sm"
            color={theme.colors.gray[6]}
            mr="xl"
          />
          <Box display="flex" className={classes.hiddenMobile}>
            <Image src="/dashboard/TsuwariInCircle.svg" width={30} height={30} alt="Tsuwari Logo" />
            <Text
              component="span"
              ml="sm"
              sx={{
                color: 'white',
                fontFamily: 'Golos Text, sans-serif',
              }}
              fz="xl"
              fw={500}
            >
              Twir
            </Text>
          </Box>
        </Flex>
        <Group position="center">
          <Button className={classes.hiddenMobile} variant={'light'} onClick={() => router.push('/application')}>Application</Button>
          <Tooltip label={'Discord'} withArrow>
            <ActionIcon
              size={'lg'}
              variant={'default'}
              component="a"
              href="https://discord.gg/Q9NBZq3zVV"
              target={'_blank'}
            >
              <DiscordSvg width={20} fill={'#e3e3e4'} />
            </ActionIcon>
          </Tooltip>
          <Tooltip label={'Changelog'} withArrow>
            <ActionIcon
              size={'lg'}
              variant={'default'}
              title={'Changelog'}
              onClick={() => setModalOpened(true)}
            >
              <IconNotes />
            </ActionIcon>
          </Tooltip>

          <Divider orientation="vertical" />

          <ActionIcon
            size="lg"
            variant="default"
            color={colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {colorScheme === 'dark' ? <IconSun size={18} /> : <IconMoonStars size={18} />}
          </ActionIcon>
          <Menu transition="pop" shadow="md" withArrow width={200}>
            <Menu.Target>
              <ActionIcon size="lg" title="Toggle language" variant="default">
                <IconLanguage size={18} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Change language</Menu.Label>
              <Menu.Divider />
              {Array.from(LOCALES.entries()).map(([lang, { icon, name }]) => (
                <Menu.Item
                  style={{ fontWeight: lang === locale ? 'bold' : 'initial' }}
                  icon={icon}
                  key={lang}
                  onClick={() => toggleLocale(lang)}
                >
                  {name}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>
          {isLoadingProfile && <Loader />}
          {!isLoadingProfile && userData && <Profile user={userData} />}
        </Group>
      </Container>

      <ChangelogModal />
    </Header>
  );
}
