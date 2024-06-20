<template>
  <!--
   q-layout 有个 style="min-height: 835.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->
  <q-layout view="hHh lpR lFr" style="min-height: 0px" class="sizeChange">
    <q-header
      elevated
      class="bg-primary text-white rounded-t-lg"
      style="width: calc(100% - 10px); right: 5px"
    >
      <q-bar class="q-electron-drag rounded-t-lg">
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <!-- <q-btn flat dense round icon="keyboard_alt"></q-btn> -->
        <div class="flex">
          <div>KeyTone</div>
          <div class="ml-2 text-xs pt-1.5">{{ version }}</div>
        </div>

        <q-space />

        <q-btn dense flat icon="horizontal_rule" @click="minimize" />
        <!-- <q-btn dense flat icon="crop_square" @click="toggleMaximize" /> -->
        <q-btn dense flat icon="close" @click="closeApp" />
      </q-bar>
    </q-header>

    <!-- 实际上的抽屉在q-layout中, 因此将.sizeChange这个css放到那边而不是这里 -->
    <q-drawer
      v-model="leftDrawerOpen"
      side="left"
      overlay
      behavior="desktop"
      elevated
    >
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
const version = process.env.APP_VERSION;

const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

function minimize() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.minimize();
  }
}

function toggleMaximize() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.toggleMaximize();
  }
}

function closeApp() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.close();
  }
}
</script>

<style lang="scss" scoped>
.sizeChange {
  :deep(.q-drawer) {
    height: calc(100% - 42px);
    border-bottom-left-radius: 0.25rem /* 4px */;
  }
}
</style>
