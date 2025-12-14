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
            v-model="formData.contactEmail"
            :label="$t('exportFlow.contact.emailLabel')"
            :placeholder="$t('exportFlow.contact.emailPlaceholder')"
            filled
            dense
            type="email"
            counter
            maxlength="200"
            @blur="validateEmail"
          />
          <q-input
            v-model="formData.contactAdditional"
            :label="$t('exportFlow.contact.additionalLabel')"
            :placeholder="$t('exportFlow.contact.additionalPlaceholder')"
            filled
            dense
            type="textarea"
            autogrow
            counter
            maxlength="500"
            class="q-mt-sm"
          />
          <div class="text-caption text-grey q-mt-xs">
            {{ $t('exportFlow.contact.hint') }}
          </div>
          <div v-if="emailError" class="text-caption text-negative q-mt-xs">
            {{ emailError }}
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
      contactEmail?: string;
      contactAdditional?: string;
    }
  ): void;
  (e: 'cancel'): void;
}

interface FormData {
  needSignature: boolean;
  requireAuthorization: boolean;
  contactEmail: string;
  contactAdditional: string;
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
  contactEmail: '',
  contactAdditional: '',
});
const emailError = ref('');
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

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

const isValidEmail = (value: string) => emailPattern.test(value.trim());

// Computed: Can continue button be enabled?
const canContinue = computed(() => {
  if (formData.value.requireAuthorization && !isValidEmail(formData.value.contactEmail)) {
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
        contactEmail: '',
        contactAdditional: '',
      };
      emailError.value = '';
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
    formData.value.contactEmail = '';
    formData.value.contactAdditional = '';
    emailError.value = '';
  }
};

const validateEmail = () => {
  if (!formData.value.requireAuthorization) {
    emailError.value = '';
    return;
  }

  const trimmed = formData.value.contactEmail.trim();
  if (!trimmed) {
    emailError.value = t('exportFlow.contact.emailRequired');
    return;
  }

  if (!isValidEmail(trimmed)) {
    emailError.value = t('exportFlow.contact.emailInvalid');
    return;
  }

  emailError.value = '';
};

const onSubmit = () => {
  if (!canContinue.value) {
    validateEmail();
    return;
  }

  const email = formData.value.contactEmail.trim();
  const additional = formData.value.contactAdditional.trim();
  emit('submit', {
    needSignature: formData.value.needSignature,
    requireAuthorization: formData.value.requireAuthorization,
    contactEmail: formData.value.requireAuthorization ? email : undefined,
    contactAdditional: formData.value.requireAuthorization && additional ? additional : undefined,
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
