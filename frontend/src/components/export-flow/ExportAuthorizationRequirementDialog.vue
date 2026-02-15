<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 340px">
      <q-card-section class="bg-primary bg-opacity-90 text-white q-pa-sm sticky top-0 z-10 backdrop-blur-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.authRequire.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pa-md q-pt-md">
        <div class="text-caption q-mb-sm">{{ $t('exportFlow.authRequire.description') }}</div>
        <q-option-group v-model="requireAuthLocal" :options="options" color="primary" />
        <div class="text-caption text-grey q-mt-sm">
          {{ $t('exportFlow.authRequire.recommend') }}
        </div>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-sm q-gutter-xs" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat size="sm" color="primary" :label="$t('exportFlow.common.cancel')" @click="onCancel" />
        <q-btn unelevated size="sm" color="primary" :label="$t('exportFlow.common.continue')" @click="onSubmit" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{ visible: boolean }>();
const emit = defineEmits<{
  (e: 'submit', payload: { requireAuthorization: boolean }): void;
  (e: 'cancel'): void;
}>();

const { t } = useI18n();
const dialogVisible = computed({
  get: () => props.visible,
  set: (v) => {
    if (!v) emit('cancel');
  },
});
const isVisible = dialogVisible;

const requireAuthLocal = ref<boolean>(false);
const options = computed(() => [
  { label: t('exportFlow.authRequire.no'), value: false },
  { label: t('exportFlow.authRequire.yes'), value: true },
]);

function onCancel() {
  emit('cancel');
}
function onSubmit() {
  emit('submit', { requireAuthorization: requireAuthLocal.value });
}
</script>
