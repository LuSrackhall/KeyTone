/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/Main_layout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/Main_page.vue'),
      },
      // TIPS: (1). 注意, 这里我们的Setting_page.vue想要共用Main_layout布局, 但是写在这个位置(即'/'路径的children)是不对的
      //            * 因为这样无法使进一步嵌套的children共用layout组件。 或者说, 它在继承共用layout组件的同时继承了setting_page.vue。 因此结果不符合预期。
      //            * 而我们希望的是, 仅继承layout页面。 想要如此, 我们需要从全局开始重新写setting路由, 而不是从'/'路由路径的children写。
      // {
      //   path: '/setting',
      //   component: () => import('pages/Setting_page.vue'),
      //   children: [
      //     // TIPS: vue-router中, children中的路由路径, 是独立的路径, 即:
      //     //       /setting-language              // yes
      //     //       /setting/setting-language      // no
      //     // TIPS: 注意, 我的children使用的都是Page组件
      //     { path: '/setting-language', component: () => import('pages/SettingPageChildren/Language_setting.vue') },
      //   ],
      // },
    ],
  },
  // 注意, 这里我们想要共用Main_layout布局, 就要这样写。 (TIPS: (2). 不是从'/'路径的children写, 而是想现在这样从全局开写)
  {
    path: '/setting',
    component: () => import('layouts/Main_layout.vue'),
    children: [
      { path: '', component: () => import('pages/Setting_page.vue') },
      // TIPS: vue-router中, children中的路由路径, 是独立的路径, 即:
      //       /setting-language              // yes
      //       /setting/setting-language      // no
      // TIPS: 注意, 我的children使用的都是Page组件
      { path: '/setting-language', component: () => import('pages/SettingPageChildren/Language_setting.vue') },
      { path: '/setting-mainHome', component: () => import('pages/SettingPageChildren/MainHome_setting.vue') },
      {
        path: '/setting-startupAndAutoStartup',
        component: () => import('pages/SettingPageChildren/StartupAndAutoStartup_setting.vue'),
      },
      {
        path: '/setting-volumeAmplify',
        component: () => import('pages/SettingPageChildren/VolumeAmplify_setting.vue'),
      },
    ],
  },
  {
    path: '/keytone_pkg',
    component: () => import('layouts/Main_layout.vue'),
    children: [
      { path: '', component: () => import('pages/Keytone_pkg_page.vue') },
      { path: '/keytone_pkg-create', component: () => import('pages/KeytonePkgPageChildren/Create_keytone_pkg.vue') },
    ],
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
