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
