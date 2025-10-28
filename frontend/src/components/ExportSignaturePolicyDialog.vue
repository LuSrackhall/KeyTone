<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card class="export-policy-dialog" style="width: 90%; max-width: 340px">
      <!-- Header -->
      <q-card-section class="bg-primary text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.policyDialog.title') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-md q-pt-lg">
        <div class="text-caption q-mb-md">
          {{ $t('exportFlow.policyDialog.description') }}
        </div>

        <!-- Signature Requirement Options (only for first export without signatures) -->
        <div v-if="!albumHasSignature" class="q-mb-md">
          <div class="q-mb-sm">
            <q-option-group v-model="formData.needSignature" :options="signatureOptions" color="primary" />
          </div>

          <!-- Signature Option Hints -->
          <div class="q-gutter-sm q-mb-md">
            <div class="text-caption text-grey">📋 {{ signatureOptions[0].description }}</div>
            <div class="text-caption text-grey">✍️ {{ signatureOptions[1].description }}</div>
          </div>
        </div>

        <!-- Authorization Requirement Checkbox -->
        <div class="q-mb-md">
          <q-checkbox
            v-model="formData.requireAuthorization"
            :label="$t('exportFlow.policyDialog.authRequired')"
            color="primary"
            class="text-caption"
            @update:model-value="onAuthorizationToggle"
          />
          <div class="text-caption text-grey q-ml-lg q-mt-xs">
            {{ $t('exportFlow.policyDialog.authRequiredHint') }}
          </div>
        </div>

        <!-- Contact Information Input (visible when authorization is enabled) -->
        <div v-if="formData.requireAuthorization" class="q-mb-md">
          <q-input
            v-model="formData.contact"
            :label="$t('exportFlow.policyDialog.contact')"
            :placeholder="$t('exportFlow.policyDialog.contactPlaceholder')"
            :rules="contactRules"
            filled
            dense
            counter
            maxlength="200"
            type="text"
            @blur="validateContact"
          />
          <div class="text-caption text-grey q-mt-xs">
            {{ $t('exportFlow.policyDialog.contactHint') }}
          </div>
          <div v-if="contactError" class="text-caption text-negative q-mt-xs">
            {{ contactError }}
          </div>
        </div>

        <!-- Tips Section -->
        <q-banner class="bg-info text-white q-mb-md" rounded dense>
          <template #avatar>
            <q-icon name="info" size="sm" />
          </template>
          <div class="text-caption">
            {{ $t('exportFlow.policyDialog.tipsTitle') }}
          </div>
          <div class="text-caption q-mt-xs" style="line-height: 1.3">
            {{ $t('exportFlow.policyDialog.tips') }}
          </div>
        </q-banner>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat :label="$t('exportFlow.policyDialog.cancel')" color="primary" size="sm" @click="onCancel" />
        <q-btn
          unelevated
          :label="$t('exportFlow.policyDialog.continue')"
          color="primary"
          size="sm"
          :disable="!canContinue"
          @click="onSubmit"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

interface PolicyDialogProps {
  visible: boolean;
  albumHasSignature?: boolean;
  defaultRequireAuth?: boolean;
}

interface PolicyDialogEmits {
  (
    e: 'submit',
    data: {
      needSignature: boolean;
      requireAuthorization: boolean;
      contact?: string;
    }
  ): void;
  (e: 'cancel'): void;
}

interface FormData {
  needSignature: boolean;
  requireAuthorization: boolean;
  contact: string;
}

const props = withDefaults(defineProps<PolicyDialogProps>(), {
  visible: false,
  albumHasSignature: false,
  defaultRequireAuth: false,
});

const emit = defineEmits<PolicyDialogEmits>();
const { t } = useI18n();

const isVisible = ref(false);
const formData = ref<FormData>({
  needSignature: true,
  requireAuthorization: false,
  contact: '',
});
const contactError = ref('');

// Signature requirement options (only for no-signature albums)
const signatureOptions = computed(() => [
  {
    label: t('exportFlow.policyDialog.requireSignature'),
    value: true,
    description: t('exportFlow.policyDialog.requireSignatureHint'),
  },
  {
    label: t('exportFlow.policyDialog.noSignature'),
    value: false,
    description: t('exportFlow.policyDialog.noSignatureHint'),
  },
]);

// Contact validation rules
const contactRules = computed(() => [
  (val: string) => {
    if (formData.value.requireAuthorization && !val) {
      return t('exportFlow.policyDialog.contactRequired');
    }
    if (val && val.length > 500) {
      return 'Max 500 characters';
    }
    return true;
  },
]);

// Computed: Can continue button be enabled?
const canContinue = computed(() => {
  if (formData.value.requireAuthorization && !formData.value.contact.trim()) {
    return false;
  }
  return true;
});

// Watch visible prop
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      // Reset form on open
      formData.value = {
        needSignature: true,
        requireAuthorization: props.defaultRequireAuth,
        contact: '',
      };
      contactError.value = '';
    }
  }
);

// Watch isVisible to sync with parent
watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('cancel');
  }
});

// Handlers
const onAuthorizationToggle = () => {
  if (!formData.value.requireAuthorization) {
    formData.value.contact = '';
    contactError.value = '';
  }
};

const validateContact = () => {
  if (formData.value.requireAuthorization && !formData.value.contact.trim()) {
    contactError.value = t('exportFlow.policyDialog.contactRequired');
  } else {
    contactError.value = '';
  }
};

const onSubmit = () => {
  if (!canContinue.value) {
    validateContact();
    return;
  }

  emit('submit', {
    needSignature: formData.value.needSignature,
    requireAuthorization: formData.value.requireAuthorization,
    contact: formData.value.requireAuthorization ? formData.value.contact : undefined,
  });

  isVisible.value = false;
};

const onCancel = () => {
  isVisible.value = false;
};
</script>

<style scoped lang="scss">
.export-policy-dialog {
  border-radius: 8px;

  :deep(.q-card__section) {
    padding: 16px;

    &:first-child {
      padding: 12px 16px;
    }
  }

  .text-caption {
    line-height: 1.5;
  }

  .whitespace-pre-wrap {
    white-space: pre-wrap;
    word-break: break-word;
  }
}
</style>
