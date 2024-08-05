<template>
  <q-page>
    <q-list>
      <template v-for="(menuItem, index) in menuList" :key="index">
        <q-expansion-item
          default-opened
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
    <!-- WARN:
    // 验证当前页面滚动条溢出界面的bug -->
    <!-- <div class="h-20"></div>
    <div class="h-20"></div>
    <div class="h-20"></div>
    <div class="h-20"></div>
    <div class="h-20"></div>
    <div class="h-20"></div> -->
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
