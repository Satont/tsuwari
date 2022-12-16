import { useMantineTheme } from '@mantine/core';
import type { NextPage } from 'next';
import Head from 'next/head';

const Home: NextPage = () => {
  const theme = useMantineTheme();

  return (
    <div>
      <Head>
        <title>Create Next App</title>
        <meta name="description" content="Generated by create next app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {theme.colorScheme}
    </div>
  );
};

export default Home;
