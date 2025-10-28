<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 340px">
      <q-card-section class="bg-primary text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.contact.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pa-md">
        <div class="text-caption q-mb-xs">{{ $t('exportFlow.contact.description') }}</div>
        <q-input v-model="contactLocal" :label="$t('exportFlow.contact.label')" dense outlined :rules="rules" />
        <div class="text-caption text-grey q-mt-xs">{{ $t('exportFlow.contact.hint') }}</div>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat size="sm" color="primary" :label="$t('exportFlow.common.cancel')" @click="onCancel" />
        <q-btn
          unelevated
          size="sm"
          color="primary"
          :disable="!canContinue"
          :label="$t('exportFlow.common.continue')"
          @click="onSubmit"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{ visible: boolean }>();
const emit = defineEmits<{
  (e: 'submit', payload: { contact: string }): void;
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

const contactLocal = ref('');
const rules = [(val: string) => !!val || t('exportFlow.contact.required')];
const canContinue = computed(() => !!contactLocal.value);

function onCancel() {
  emit('cancel');
}
function onSubmit() {
  emit('submit', { contact: contactLocal.value });
}
</script>
