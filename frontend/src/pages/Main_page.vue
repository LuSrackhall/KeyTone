<template>
  <!-- <div>Silent Typist's Friend</div> -->
  <!--
   q-layout 有个 style="min-height: 803.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->
  <q-page style="min-height: 0px">
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
          :icon="setting_store.audioVolumeProcessing.volumeSilent ? 'volume_off' : 'volume_up'"
          @click="isSilent"
        >
        </q-btn>

        <q-slider
          :class="['w-[80%]']"
          v-model="setting_store.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-setting_store.audioVolumeProcessing.volumeAmplify - 5"
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
          v-model="setting_store.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-setting_store.audioVolumeProcessing.volumeAmplify - 5"
          :step="0"
          :markers="markersDebug"
          marker-labels
          label
          :label-value="labelValueDebug"
          color="light-green"
        />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import logoUrl from 'assets/img/KeyTone.png?url';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, watch } from 'vue';

const setting_store = useSettingStore();

const labelValue = computed(() => {
  const percentage = (
    (1 - -setting_store.audioVolumeProcessing.volumeNormal / (setting_store.audioVolumeProcessing.volumeAmplify + 5)) *
    100
  )
    .toFixed(2)
    .split('.');
  return percentage[1] === '00' ? percentage[0] + '%' : percentage[0] + '.' + percentage[1] + '%';
});

watch(
  () => setting_store.audioVolumeProcessing.volumeNormal,
  () => {
    // 当用户拖动音量进度条时, 自动解除静音
    setting_store.audioVolumeProcessing.volumeSilent = false;
  }
);

const labelValueDebug = computed(() => {
  return setting_store.audioVolumeProcessing.volumeNormal.toFixed(2);
});
const markersDebug = computed(() => {
  return (setting_store.audioVolumeProcessing.volumeAmplify + 5) / 1;
});

const isSilent = () => {
  setting_store.audioVolumeProcessing.volumeSilent = !setting_store.audioVolumeProcessing.volumeSilent;
};
</script>

<style lang="scss" scoped></style>
