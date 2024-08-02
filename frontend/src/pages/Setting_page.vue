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
const menuList = [
  {
    icon: 'language',
    label: 'language.index',
    caption: 'language.caption',
    separator: true,
    to: '/setting-language',
  },
];
</script>

<style lang="scss" scoped></style>
