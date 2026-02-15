<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 340px">
      <q-card-section class="bg-warning text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.authImpact.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pa-md">
        <div class="text-caption" style="line-height: 1.3">
          {{ $t('exportFlow.authImpact.message') }}
        </div>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-sm q-gutter-xs" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat size="sm" color="primary" :label="$t('exportFlow.common.back')" @click="onBack" />
        <q-btn unelevated size="sm" color="primary" :label="$t('exportFlow.authImpact.confirm')" @click="onConfirm" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ visible: boolean }>();
const emit = defineEmits<{
  (e: 'confirm'): void;
  (e: 'back'): void;
}>();

const dialogVisible = computed({
  get: () => props.visible,
  set: () => {
    /* Dialog is controlled by parent; no local side-effects needed on close */
  },
});
const isVisible = dialogVisible;

function onConfirm() {
  emit('confirm');
}
function onBack() {
  emit('back');
}
</script>
