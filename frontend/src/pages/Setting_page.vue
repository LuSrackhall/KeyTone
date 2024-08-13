<template>
  <q-page>
    <q-scroll-area class="w-[379px] h-[458.5px]">
      <q-list>
        <template v-for="(menuItem, index) in menuList" :key="index">
          <q-expansion-item
            :icon="menuItem.icon"
            :label="$t('setting.' + menuItem.label)"
            :caption="$t('setting.' + menuItem.caption)"
            :to="menuItem.to"
            :header-inset-level="0"
            :content-inset-level="0"
          >
            <language v-if="menuItem.to === '/setting-language'"></language>
            <StartupAndAutoStartup v-if="menuItem.to === '/setting-startupAndAutoStartup'"></StartupAndAutoStartup>
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

const menuList = [
  {
    icon: 'language',
    label: 'language.index',
    caption: 'language.caption',
    separator: true,
    to: '/setting-language',
  },
  {
    icon: 'sunny_snowing',
    label: '启动与自动启动.启动与自动启动.index',
    caption: '启动与自动启动.启动与自动启动.caption',
    separator: true,
    to: '/setting-startupAndAutoStartup',
  },
];
</script>

<style lang="scss" scoped></style>
