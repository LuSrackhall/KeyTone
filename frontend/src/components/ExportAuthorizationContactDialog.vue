<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 340px">
      <q-card-section class="bg-primary text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.contact.title') }}</div>
      </q-card-section>

      <q-card-section class="q-pa-md">
        <div class="text-caption q-mb-xs">{{ $t('exportFlow.contact.description') }}</div>
        <q-input
          v-model="emailLocal"
          type="email"
          dense
          outlined
          :label="$t('exportFlow.contact.emailLabel')"
          :placeholder="$t('exportFlow.contact.emailPlaceholder')"
          :rules="emailRules"
        />
        <q-input
          v-model="additionalLocal"
          type="textarea"
          autogrow
          dense
          outlined
          class="q-mt-sm"
          :label="$t('exportFlow.contact.additionalLabel')"
          :placeholder="$t('exportFlow.contact.additionalPlaceholder')"
        />
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
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{ visible: boolean }>();
const emit = defineEmits<{
  (e: 'submit', payload: { contact: string; email: string; additional?: string }): void;
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

const emailLocal = ref('');
const additionalLocal = ref('');

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const emailRules = [
  (val: string) => (!!val && !!val.trim()) || t('exportFlow.contact.emailRequired'),
  (val: string) => !val || emailPattern.test(val.trim()) || t('exportFlow.contact.emailInvalid'),
];

const emailValid = computed(() => emailPattern.test(emailLocal.value.trim()));
const canContinue = computed(() => emailValid.value);

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      emailLocal.value = '';
      additionalLocal.value = '';
    }
  }
);

function onCancel() {
  emit('cancel');
}
function onSubmit() {
  const email = emailLocal.value.trim();
  if (!emailValid.value) {
    return;
  }
  const additional = additionalLocal.value.trim();
  const contact = additional ? `${email}\n${additional}` : email;
  emit('submit', { contact, email, additional: additional || undefined });
}
</script>
