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
  <q-page>
    <q-scroll-area class="w-[379px] h-[458.5px]">
      <q-list>
        <template v-for="(menuItem, index) in menuList" :key="index">
          <q-expansion-item
            v-model="setting_store.settingItemsOpenedState[index]"
            :icon="menuItem.icon"
            :label="$t(menuItem.label)"
            :caption="$t(menuItem.caption)"
            :to="menuItem.to"
            :header-inset-level="0"
            :content-inset-level="0"
          >
            <Language v-if="menuItem.to === '/setting-language'"></Language>
            <MainHome v-if="menuItem.to === '/setting-mainHome'"></MainHome>
            <StartupAndAutoStartup v-if="menuItem.to === '/setting-startupAndAutoStartup'"></StartupAndAutoStartup>
            <VolumeAmplify v-if="menuItem.to === '/setting-volumeAmplify'"></VolumeAmplify>
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
import { useSettingStore } from 'src/stores/setting-store';

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
    icon: 'volume_up',
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
];
</script>

<style lang="scss" scoped></style>
