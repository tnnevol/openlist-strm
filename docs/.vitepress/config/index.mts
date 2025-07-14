import { withPwa } from '@vite-pwa/vitepress';
import { defineConfigWithTheme } from 'vitepress';

import { en } from './en.mts';
import { shared } from './shared.mts';
import { zh } from './zh.mts';

export default withPwa(
  defineConfigWithTheme({
    ...shared,
    // 移除多语言配置，只保留中文
    ...zh,
  }),
);
