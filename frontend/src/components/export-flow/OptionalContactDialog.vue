<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 360px">
      <!-- Header -->
      <q-card-section class="bg-teal text-white q-pa-sm">
        <div class="text-subtitle1">{{ t('exportFlow.optionalContact.title') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-md">
        <div class="text-caption text-grey-7 q-mb-md">
          {{ t('exportFlow.optionalContact.description') }}
        </div>

        <!-- Email Input -->
        <q-input
          v-model="formData.email"
          :label="t('exportFlow.contact.emailLabel')"
          :placeholder="t('exportFlow.contact.emailPlaceholder')"
          filled
          dense
          type="email"
          counter
          maxlength="200"
          @blur="validateEmail"
        />
        <div v-if="emailError" class="text-caption text-negative q-mt-xs">
          {{ emailError }}
        </div>

        <!-- Additional Contact Input -->
        <q-input
          v-model="formData.additional"
          :label="t('exportFlow.contact.additionalLabel')"
          :placeholder="t('exportFlow.contact.additionalPlaceholder')"
          filled
          dense
          type="textarea"
          autogrow
          counter
          maxlength="500"
          class="q-mt-sm"
        />

        <div class="text-caption text-grey q-mt-sm">
          {{ t('exportFlow.optionalContact.hint') }}
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat :label="t('exportFlow.optionalContact.skip')" color="grey" size="sm" @click="onSkip" />
        <q-btn
          unelevated
          :label="t('exportFlow.optionalContact.continue')"
          color="teal"
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

interface OptionalContactDialogProps {
  visible: boolean;
}

interface OptionalContactDialogEmits {
  (e: 'submit', data: { email?: string; additional?: string }): void;
  (e: 'skip'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<OptionalContactDialogProps>(), {
  visible: false,
});

const emit = defineEmits<OptionalContactDialogEmits>();
const { t } = useI18n();

const isVisible = ref(false);
const formData = ref({
  email: '',
  additional: '',
});
const emailError = ref('');
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Validate email format (only if email is not empty)
const isValidEmail = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return true; // Empty is valid (optional)
  return emailPattern.test(trimmed);
};

// Computed: Can continue button be enabled?
const canContinue = computed(() => {
  // If email is entered, it must be valid format
  return isValidEmail(formData.value.email);
});

// Watch visible prop
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      // Reset form on open
      formData.value = {
        email: '',
        additional: '',
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

// Validate email on blur
const validateEmail = () => {
  const trimmed = formData.value.email.trim();
  if (!trimmed) {
    emailError.value = '';
    return;
  }

  if (!isValidEmail(trimmed)) {
    emailError.value = t('exportFlow.contact.emailInvalid');
    return;
  }

  emailError.value = '';
};

// Skip - proceed without contact info
const onSkip = () => {
  emit('skip');
  isVisible.value = false;
};

// Submit - proceed with contact info
const onSubmit = () => {
  if (!canContinue.value) {
    validateEmail();
    return;
  }

  const email = formData.value.email.trim();
  const additional = formData.value.additional.trim();

  emit('submit', {
    email: email || undefined,
    additional: additional || undefined,
  });

  isVisible.value = false;
};
</script>
