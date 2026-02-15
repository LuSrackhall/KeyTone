<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 340px">
      <q-card-section class="bg-primary bg-opacity-90 text-white q-pa-sm sticky top-0 z-10 backdrop-blur-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.confirmSignature.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pa-md">
        <div class="text-caption q-mb-sm">
          {{ $t('exportFlow.confirmSignature.description') }}
        </div>
        <q-option-group v-model="needSignatureLocal" :options="options" color="primary" />
      </q-card-section>

      <q-card-actions align="right" class="q-pa-sm q-gutter-xs" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat color="primary" size="sm" :label="$t('exportFlow.common.cancel')" @click="onCancel" />
        <q-btn unelevated color="primary" size="sm" :label="$t('exportFlow.common.continue')" @click="onSubmit" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

interface Props {
  visible: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'submit', payload: { needSignature: boolean }): void;
  (e: 'cancel'): void;
}>();

const { t } = useI18n();

const dialogVisible = computed({
  get: () => props.visible,
  set: (v: boolean) => {
    if (!v) emit('cancel');
  },
});

const needSignatureLocal = ref<boolean>(true);

const options = computed(() => [
  { label: t('exportFlow.confirmSignature.needSignature'), value: true },
  { label: t('exportFlow.confirmSignature.noSignature'), value: false },
]);

function onCancel() {
  emit('cancel');
}

function onSubmit() {
  emit('submit', { needSignature: needSignatureLocal.value });
}

const isVisible = dialogVisible;
</script>
