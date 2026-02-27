<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="min-width: 350px">
      <q-card-section class="sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        <div class="text-h6">{{ t('exportFlow.reexportWarning.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        {{ t('exportFlow.reexportWarning.description') }}
      </q-card-section>

      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="t('exportFlow.common.cancel')" color="primary" @click="onCancel" />
        <q-btn flat :label="t('exportFlow.common.confirm')" color="primary" @click="onConfirm" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: 'confirm'): void;
  (e: 'cancel'): void;
}>();

const { t } = useI18n();

const isVisible = computed({
  get: () => props.visible,
  set: (val) => {
    if (!val) emit('cancel');
  },
});

const onConfirm = () => {
  emit('confirm');
};

const onCancel = () => {
  emit('cancel');
};
</script>
