import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'svg:openlist-logo',
      order: 100,
      title: $t('openlist.title'),
    },
    name: 'Openlist',
    path: '/openlist',
    children: [
      //openlist 服务维护
      {
        name: 'OpenlistService',
        path: 'service',
        component: () => import('#/views/openlist/service/list.vue'),
        meta: {
          title: $t('openlist.server.title'),
          icon: 'lucide:server',
        },
      },
    ],
  },
];

export default routes;
