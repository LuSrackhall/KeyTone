<!--
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
-->

<template>
  <q-page style="min-height: 0px" :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]']">
    <q-scroll-area :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]']">
      <q-list>
        <template v-for="(menuItem, index) in menuList" :key="index">
          <q-expansion-item
            v-model="setting_store.settingItemsOpenedState[index]"
            :icon="menuItem.icon"
            :label="$t(menuItem.label)"
            :caption="$t(menuItem.caption)"
            :header-inset-level="0"
            :content-inset-level="0"
            @click="handleItemClick(menuItem.to)"
          >
            <Language v-if="menuItem.to === '/setting-language'"></Language>
            <MainHome v-if="menuItem.to === '/setting-mainHome'"></MainHome>
            <StartupAndAutoStartup v-if="menuItem.to === '/setting-startupAndAutoStartup'"></StartupAndAutoStartup>
            <VolumeAmplify v-if="menuItem.to === '/setting-volumeAmplify'"></VolumeAmplify>
            <KeytoneAlbumPage v-if="menuItem.to === '/setting-keytoneAlbumPage'"></KeytoneAlbumPage>
          </q-expansion-item>
          <q-separator :key="'sep' + index" v-if="menuItem.separator" />
        </template>
      </q-list>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts">
import Language from 'src/pages/SettingPageChildren/Language_setting.vue';
import MainHome from 'src/pages/SettingPageChildren/MainHome_setting.vue';
import StartupAndAutoStartup from 'src/pages/SettingPageChildren/StartupAndAutoStartup_setting.vue';
import VolumeAmplify from './SettingPageChildren/VolumeAmplify_setting.vue';
import KeytoneAlbumPage from './SettingPageChildren/KeytoneAlbumPage_setting.vue';
import { useSettingStore } from 'src/stores/setting-store';
import { ref, onMounted } from 'vue';

const setting_store = useSettingStore();

const menuList = [
  {
    icon: 'language',
    label: 'setting.language.index',
    caption: 'setting.language.caption',
    separator: true,
    to: '/setting-language',
  },
  {
    icon: 'home',
    label: 'setting.mainHome.mainHome.index',
    caption: 'setting.mainHome.mainHome.caption',
    separator: true,
    to: '/setting-mainHome',
  },
  {
    icon: 'sunny_snowing',
    label: 'setting.启动与自动启动.启动与自动启动.index',
    caption: 'setting.启动与自动启动.启动与自动启动.caption',
    separator: true,
    to: '/setting-startupAndAutoStartup',
  },
  {
    icon: 'volume_up',
    label: 'setting.原始音量增减调节.原始音量增减调节.index',
    caption: 'setting.原始音量增减调节.原始音量增减调节.caption',
    separator: true,
    to: '/setting-volumeAmplify',
  },
  {
    icon: 'library_music',
    label: 'KeyTone.setting.category.keytoneAlbumPage',
    caption: 'KeyTone.setting.category.keytoneAlbumPageCaption',
    separator: true,
    to: '/setting-keytoneAlbumPage',
  },
];

import { useRouter } from 'vue-router';

const router = useRouter();

/**
 * 上一次单击的时间戳（用于双击检测）。
 * 采用时间戳对比而非定时器方案，避免 setTimeout 回调与组件卸载的竞态问题。
 */
let lastClickTime = 0;

/**
 * 在组件挂载时，确保 settingItemsOpenedState 数组长度与菜单项数量一致。
 * 防止因数组长度不足导致 v-model 绑定到 undefined 索引时出现异常。
 */
onMounted(() => {
  while (setting_store.settingItemsOpenedState.length < menuList.length) {
    setting_store.settingItemsOpenedState.push(false);
  }
});

/**
 * 处理设置分组头部的点击/双击事件。
 *
 * - 单击（默认行为）：由 Quasar 的 q-expansion-item 内部处理展开/收起。
 *   我们的 handler 仅记录时间戳，不干预默认行为。
 * - 双击（间隔 < 300ms）：视为"进入该设置组的独立详情页"，触发路由导航。
 *   此时 q-expansion-item 的展开/收起同样会触发两次（开→关 或 关→开），
 *   但由于页面即将跳转，用户无感知。
 *
 * 注意：q-expansion-item 的 @click 是组件事件（仅头部 QItem 触发），
 * 内容区域（如 select、input 等控件）的点击不会传入此 handler。
 */
function handleItemClick(menuItemTo: string) {
  const now = Date.now();
  if (now - lastClickTime < 300) {
    // 检测到双击 → 导航到该设置组的独立页面
    lastClickTime = 0;
    router.push(menuItemTo);
  } else {
    lastClickTime = now;
  }
}

const isMacOS = ref(getMacOSStatus());
function getMacOSStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getMacOSStatus();
  }
  return false;
}
</script>

<style lang="scss" scoped>
// 用于修复主页面全局的:global(.q-field__native)中的h-5.8这个样式影响了当前页面中的q-input的问题
:deep(.q-placeholder) {
  // 在这里重置q-input组件的输入样式的高度以修复这个问题
  @apply h-auto;
}

:deep(.q-item__section) {
  @apply text-wrap;
  @apply overflow-hidden;
}
</style>
