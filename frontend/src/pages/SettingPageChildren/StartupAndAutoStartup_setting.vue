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

  <!-- 是否开启开机自动启动 -->
  <q-item>
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
  <q-item :disable="!setting_store.autoStartup.isAutoRun">
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

const setting_store = useSettingStore();
</script>

<style lang="scss" scoped></style>
