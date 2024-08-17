<template>
  <q-item :class="['h-15 mb-5']">
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid']">
      <div :class="['w-[88%] flex justify-between items-center']">
        <q-input
          dense
          hide-bottom-space
          :class="['w-[72%] h-10.5 ']"
          v-model.number="setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope"
          type="number"
          filled
          :label="$t('setting.mainHome.音量降幅.index')"
          stack-label
          :rules="[(val: number) => { return val >= 5 && val<100000000 || $t('setting.mainHome.音量降幅.rulesErrorInfo'); }]"
        />

        <q-btn
          :class="['w-15 h-5']"
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

<style lang="scss" scoped></style>
