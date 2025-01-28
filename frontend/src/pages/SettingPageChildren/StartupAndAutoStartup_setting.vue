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
  <!-- 启动时是否隐藏窗口 -->
  <q-item>
    <!-- 左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.启动与自动启动.启动时隐藏窗口.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.启动与自动启动.启动时隐藏窗口.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.startup.isHideWindows" />
    </q-item-section>
  </q-item>

  <!-- 是否开机自动启动功能(appx版) -->
  <q-item v-if="isWindowsStore">
    <!-- 左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-12 self-center']"></div>
    <q-item-section>
      <q-item-label class="flex flew-row items-center bg-opacity-80">
        {{ $t('setting.启动与自动启动.自动启动.index') }}
        <q-badge class="ml-2 bg-opacity-70 bg-purple-500 text-[0.66rem]" text-color="white" rounded>
          {{ $t('setting.启动与自动启动.自动启动.appxIndex') }}
        </q-badge>
      </q-item-label>
      <q-item-label caption class="flex flex-col bg-opacity-80">
        {{ $t('setting.启动与自动启动.自动启动.caption') }}
        <q-badge class="self-start mt-1 bg-opacity-70 bg-purple-500 text-[0.66rem]" text-color="white" rounded>
          {{ $t('setting.启动与自动启动.自动启动.appxCaption') }}
        </q-badge>
      </q-item-label>
    </q-item-section>
  </q-item>

  <!-- 是否开启开机自动启动 -->
  <q-item v-if="!isWindowsStore">
    <!-- 左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.启动与自动启动.自动启动.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.启动与自动启动.自动启动.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.autoStartup.isAutoRun" />
    </q-item-section>
  </q-item>

  <!-- 开机自动启动时, 是否隐藏窗口 -->
  <q-item :disable="!setting_store.autoStartup.isAutoRun" v-if="!isWindowsStore">
    <!-- 左边的竖线 -->
    <div :class="['ml-13 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.启动与自动启动.自动启动时隐藏窗口.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.启动与自动启动.自动启动时隐藏窗口.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.autoStartup.isHideWindows" :disable="!setting_store.autoStartup.isAutoRun" />
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
import { useSettingStore } from 'src/stores/setting-store';
import { ref } from 'vue';

const setting_store = useSettingStore();

const isWindowsStore = ref(getWindowsStoreStatus());

function getWindowsStoreStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getWindowsStoreStatus();
  }
  return false;
}
</script>

<style lang="scss" scoped></style>
