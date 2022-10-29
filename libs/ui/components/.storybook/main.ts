import type { StorybookViteConfig } from '@storybook/builder-vite';
import { resolve } from 'node:path';
import { mergeConfig, UserConfig } from 'vite';

const config: StorybookViteConfig = {
  stories: [
    '../src/components/**/*.stories.mdx',
    '../src/components/**/*.stories.@(js|jsx|ts|tsx)',
  ],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-actions',
    '@storybook/addon-interactions',
  ],
  framework: '@storybook/vue3',
  core: {
    builder: '@storybook/builder-vite',
  },
  features: {
    storyStoreV7: true,
  },
  typescript: {
    check: false,
  },
  async viteFinal(config) {
    return mergeConfig(config, {
      resolve: {
        alias: [{ find: '@', replacement: resolve(__dirname, '..', 'src') }],
      },
    } as UserConfig);
  },
};

module.exports = config;
