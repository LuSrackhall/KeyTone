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
  <q-item :class="['h-15 mb-5']">
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid']">
      <div :class="['w-[92%] flex justify-between items-center flex-nowrap  gap-[12px]']">
        <q-input
          dense
          hide-bottom-space
          :class="['w-[66%] h-10.5 ']"
          v-model.number="setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope"
          type="number"
          filled
          :label="$t('setting.mainHome.音量降幅.index')"
          stack-label
          :rules="[(val: number) => { return val >= 5 && val<100000000 || $t('setting.mainHome.音量降幅.rulesErrorInfo'); }]"
        />

        <q-btn
          :class="['min-w-15 min-h-5']"
          color="primary"
          size="10px"
          :label="$t('setting.mainHome.重置')"
          @click="returnToNormalReduceScope()"
        />
      </div>
    </div>
  </q-item>

  <q-item>
    <!-- 左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.mainHome.音量调试滑块.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.mainHome.音量调试滑块.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.mainHome.audioVolumeProcessing.isOpenVolumeDebugSlider" />
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
import { useSettingStore } from 'src/stores/setting-store';

const setting_store = useSettingStore();

const returnToNormalReduceScope = () => {
  setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope = 5.0;
};
</script>

<style lang="scss" scoped>
// 用于修复主页面全局的:global(.q-field__native)中的h-5.8这个样式影响了当前页面中的q-input的问题
:deep(.q-placeholder) {
  // 在这里重置q-input组件的输入样式的高度以修复这个问题
  @apply h-auto;
}

:deep(.q-item__section) {
  @apply text-wrap;
}
</style>
