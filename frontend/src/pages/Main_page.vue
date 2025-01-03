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
  <!-- <div>Silent Typist's Friend</div> -->
  <!--
   q-layout 有个 style="min-height: 803.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->
  <q-page style="min-height: 0px" class="w-[379px] h-[458.5px]">
    <div
      :class="[
        '',
        'w-auto h-50 flex flex-col items-center',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <q-avatar
        :class="[
          '',
          // 设置整体头像的尺寸大小。 至于内部相关文字的尺寸大小(如果有的话  可通过设置text-[8rem] leading-[10rem]实现)
          'size-56',
          '',
        ]"
      >
        <!-- draggable="false"使得图片无法拖动, 免得影响界面体验 -->
        <img :src="logoUrl" draggable="false" />
      </q-avatar>
    </div>
    <!--
      * TIPS: 对于主页的音量调整, 理念是最终原始音频, 即最大音量等于原始音频的正常音量(或是最大音量), 也就是说保持为0。
              > 在此基础上, 其最小值始终保持, 可调整至正常音量的对数 g <目前默认g为5, 写死, 后续如果有必要, 会在设置界面对其新增相关的设置项>。
              > > 当我们通过设置界面的音量强来增大原始音频音量时, 其音量的原始正常值将不再为0, 而是变为对音量增强的指数取相同基数后的对数值。
              > > * 即实际上音频的原始音量对应的数值是 `-setting_store.audioVolumeProcessing.volumeAmplify`
              > > 因此, 此时我们的最小值, 要想保持正常音量的对数 g , 需要 使用`-setting_store.audioVolumeProcessing.volumeAmplify-5`这个数值
              > >   - completed(已完成)   FIXME:
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify > 0 时, 当然没问题不用过分讨论。
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify < 0 且  > -5 时, 使用这个算法就发生了问题(或者说永远保持正常音量的对数g的方案发生了问题)
              > >   * 此时的缩小是仍旧在以缩小的基础上进行的(即此时的0), 而此时仍保证真实正常音量的对数g的话, 是不现实的, 因为主页的音量调整是一个只能缩小音量的滑块, 这要保持只能缩短滑块范围。
              > >   * 拿-2来举例子, 貌似得出的最小值-3是更靠近正常音量了, 但实际上得出的这个最小值在处理时, 并不是一个更靠近正常音量的值(因为 -2 这个值, 在Amplify时已经被实际处理过了, 而这里在添一个-3, 仍旧是-5)
              > >     * 或者说, 此时的0对应的已经是缩小后的音量了(而不是像正数时放大的音量), 我们无法通过只能缩减音量的滑块获得正常的音量, 此时要想进一步保持对数5, 就只能缩短滑块进度条的实际范围了。
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify < -5时, 问题就更大了。
              > >   * 由于要保持真实正常音量的对数g, 我们不得不缩短进度条, 但当这个 setting_store.audioVolumeProcessing.volumeAmplify = -5 时, 这个进度条实际已经不存在了, 此时只能是正常音量的对数g。
              > >   * 而当其进一步缩小的<-5时, 不但最小音量无法保持正常音量的对数g, 甚至滑块无法表示目前真实缩小的正常音量的对数值。(算成百分比时, 甚至出现了大于100%的离谱情况, 实际上这些只是更小的负的百分数)
              > >  FIXME的完成小记: 放弃最小音量值, 始终要想保持正常音量的对数 g的方案。因为当 volumeAmplify 小于0 后, 这个方案是无意义的
              > >                  * 此时, 我们最多只能保证其是 当前已被缩小音量<即此时的0对应的音量>的对数g。或者说此时可缩小的范围固定为g<因为我们的最大值固定为0这个数>。
              > >                    * 即 当 volumeAmplify 小于0 后 的最新算法是 -5, 或者说-g
     -->
    <div
      :class="[
        '',
        'w-full flex flex-col items-center pr-5 mt-10',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <div :class="['w-56 flex justify-between items-center']">
        <q-btn
          dense
          round
          flat
          :icon="setting_store.mainHome.audioVolumeProcessing.volumeSilent ? 'volume_off' : 'volume_up'"
          @click="isSilent"
        >
        </q-btn>

        <q-slider
          :class="['w-[80%]']"
          v-model="setting_store.mainHome.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-min"
          :step="0"
          label
          :label-value="labelValue"
          color="light-green"
        />
      </div>
      <div :class="['w-56  flex justify-end items-center mt-5']">
        <q-slider
          :class="['w-[80%]']"
          v-if="setting_store.mainHome.audioVolumeProcessing.isOpenVolumeDebugSlider"
          v-model="setting_store.mainHome.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-min"
          :step="0"
          :markers="markersDebug"
          marker-labels
          label
          label-always
          :label-value="labelValueDebug"
          color="light-green"
        />
      </div>
    </div>
    <div
      :class="[
        'text-center text-xs text-gray-500 absolute bottom-8.5 w-full',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <div class="flex justify-center items-center">
        Copyright © 2024&nbsp;
        <div class="cursor-pointer hover:bg-gray-100" @click="openExternal('https://github.com/LuSrackhall')">
          LuSrackhall
        </div>
      </div>
      <div class="flex justify-center items-center">
        KeyTone is licensed under&nbsp;
        <div
          class="cursor-pointer hover:bg-gray-100"
          @click="openExternal('https://choosealicense.com/licenses/gpl-3.0/')"
        >
          GNU GPLv3
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import logoUrl from 'assets/img/KeyTone.png?url';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, watch } from 'vue';

const setting_store = useSettingStore();

const min = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      setting_store.audioVolumeProcessing.volumeAmplify +
      setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope
    );
  } else {
    //  FIXME的完成小记: 对应步骤小记, 当setting_store.audioVolumeProcessing.volumeAmplify<0时, 我们使用固定的g值作为缩小幅度(或者说已被缩小的值<即0>为参考的对数g, 因为我们的最大值固定为0, 因此也是固定的g值)
    //                   * g值 = setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope。
    return setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;
  }
});

const labelValue = computed(() => {
  const percentage = ((1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value) * 100)
    .toFixed(2)
    .split('.');
  return percentage[1] === '00' ? percentage[0] + '%' : percentage[0] + '.' + percentage[1] + '%';
});

watch(
  () => setting_store.mainHome.audioVolumeProcessing.volumeNormal,
  () => {
    // 当用户拖动音量进度条时, 自动解除静音
    setting_store.mainHome.audioVolumeProcessing.volumeSilent = false;
  }
);

const labelValueDebug = computed(() => {
  return setting_store.mainHome.audioVolumeProcessing.volumeNormal.toFixed(2);
});
const markersDebug = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      (setting_store.audioVolumeProcessing.volumeAmplify +
        setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope) /
      1
    );
  } else {
    return setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope / 1;
  }
});

const isSilent = (event: any) => {
  if (event.detail === 0) {
    // 由键盘触发，不处理
    return;
  }
  setting_store.mainHome.audioVolumeProcessing.volumeSilent =
    !setting_store.mainHome.audioVolumeProcessing.volumeSilent;
};

function openExternal(url: string) {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.openExternal(url);
  }
}
</script>

<style lang="scss" scoped></style>
