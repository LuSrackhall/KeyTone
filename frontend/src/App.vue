<template>
  <router-view />
</template>

<script setup lang="ts">
import { onBeforeMount } from 'vue';
import { useSettingStore } from 'src/stores/setting-store';

const setting_store = useSettingStore();

onBeforeMount(async () => {
  setting_store.settingInitAndRealTimeStorage();
});
</script>
