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
  <!-- <q-item>
    <div :class="['ml-0 rounded-full  mr-6 h-6 self-center']"></div>
    <div :class="['w-full flex justify-between items-center']">
      <q-badge :class="['bg-cyan-700 h-5']"> Volume: {{ setting_store.audioVolumeProcessing.volumeAmplify }} </q-badge>

      <q-btn :class="['w-15 h-5 mr-5']" color="primary" size="10px" label="重置" @click="returnToNormal()" />
    </div>
  </q-item> -->

  <q-item :class="['pt-8 h-50']">
    <!-- select左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid justify-items-center']">
      <div :class="['min-w-[90%] flex justify-between items-center -ml-6 mb-10']">
        <q-badge :class="['bg-cyan-700 h-5 ml-1.5']">
          {{ $t('setting.原始音量增减调节.音量') }}:{{ setting_store.audioVolumeProcessing.volumeAmplify }}
        </q-badge>
        <q-btn
          :class="['min-w-15 min-h-5 ml-2']"
          color="primary"
          size="10px"
          :label="$t('setting.原始音量增减调节.重置')"
          @click="returnToNormal()"
        />
      </div>

      <div :class="['w-[100%] grid justify-items-center']">
        <q-slider
          :class="['w-[88%] -ml-5']"
          v-model="setting_store.audioVolumeProcessing.volumeAmplify"
          :max="max"
          :min="min"
          :step="0"
          :markers="markers"
          marker-labels
          label
          label-always
          :label-value="labelValue"
          color="light-green"
        />
      </div>
    </div>
  </q-item>
  <q-item :class="['h-15 mb-12']">
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid justify-items-center']">
      <div :class="['min-w-[90%] flex justify-between items-center -ml-1.5']">
        <q-input
          dense
          hide-bottom-space
          :class="['min-w-[55%] h-10.5 mr-1']"
          v-model.number="setting_store.audioVolumeProcessing.volumeAmplifyLimit"
          type="number"
          filled
          :label="$t('setting.原始音量增减调节.AmplifyLimit.index')"
          stack-label
          :rules="[(val: number) => { return val > 0 && val<100000000 || $t('setting.原始音量增减调节.AmplifyLimit.rulesErrorInfo'); }]"
        />

        <q-btn
          :class="['min-w-15 min-h-6.5 mr-2.5']"
          color="primary"
          size="10px"
          :label="$t('setting.原始音量增减调节.重置')"
          @click="returnToDefaultLimit()"
        />
      </div>
    </div>
  </q-item>
</template>

<script setup lang="ts">
import { useSettingStore } from 'src/stores/setting-store';
import { ref, watch } from 'vue';
import { debounce } from 'lodash';

const setting_store = useSettingStore();

const max = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit);
const min = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1);
const markers = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit / 1);
const labelValue = ref(setting_store.audioVolumeProcessing.volumeAmplify.toFixed(2));

const debounced = debounce(
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit <
      Math.abs(setting_store.audioVolumeProcessing.volumeAmplify)
    ) {
      if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
        setting_store.audioVolumeProcessing.volumeAmplify = setting_store.audioVolumeProcessing.volumeAmplifyLimit;
      } else {
        setting_store.audioVolumeProcessing.volumeAmplify = setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1;
      }
    }
  },
  800,
  { trailing: true }
);

watch(
  () => setting_store.audioVolumeProcessing.volumeAmplifyLimit,
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit > 0 &&
      setting_store.audioVolumeProcessing.volumeAmplifyLimit < 100000000
    ) {
      markers.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit / 1;
      max.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit;
      min.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1;

      debounced.cancel;
      debounced();
    }
  }
);

// TIPS: 这里的watch处理, 是解决控件的一个bug。(涉及到多个控件的结合使用, 错误原因不明, 不过这样解藕后, 便可解决)
//       > 输入框控件中持续执行`backspace按键`的删除操作时,会引发的`q-slider 控件`的 `:label-value="labelValue"`报错。
watch(
  () => setting_store.audioVolumeProcessing.volumeAmplify,
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit > 0 &&
      setting_store.audioVolumeProcessing.volumeAmplifyLimit < 100000000
    ) {
      labelValue.value = setting_store.audioVolumeProcessing.volumeAmplify.toFixed(2);
    } else {
      // 如果setting_store.audioVolumeProcessing.volumeAmplifyLimit不符合规范, 则在用户操作音量条时将最后一次符合规范的值, 重新给到它
      setting_store.audioVolumeProcessing.volumeAmplifyLimit = max.value;
    }
  }
);

function returnToNormal() {
  setting_store.audioVolumeProcessing.volumeAmplify = 0.0;
}

function returnToDefaultLimit() {
  setting_store.audioVolumeProcessing.volumeAmplifyLimit = 10.0;
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
}
</style>
