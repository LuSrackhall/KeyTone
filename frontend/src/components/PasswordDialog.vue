<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card style="min-width: 350px">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">{{ $t('passwordDialog.title') || '输入导出密码' }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section>
        <div class="text-body2 q-mb-md text-grey-7">
          {{ $t('passwordDialog.message') || '该键音包不支持二次导出，请输入密码以获取权限' }}
        </div>
        
        <q-input
          v-model="password"
          :label="$t('passwordDialog.password') || '密码'"
          outlined
          dense
          type="password"
          autofocus
          @keyup.enter="handleConfirm"
          :error="hasError"
          :error-message="errorMessage"
        />
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('common.cancel') || '取消'" @click="closeDialog" />
        <q-btn
          unelevated
          color="primary"
          :label="$t('common.confirm') || '确认'"
          :loading="loading"
          @click="handleConfirm"
          :disable="!password.trim()"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, defineEmits, defineProps } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// Props and Emits
const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  'confirm': [password: string];
}>();

// Reactive data
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const password = ref('');
const loading = ref(false);
const hasError = ref(false);
const errorMessage = ref('');

// Methods
const handleConfirm = () => {
  if (!password.value.trim()) return;
  
  hasError.value = false;
  errorMessage.value = '';
  emit('confirm', password.value);
};

const closeDialog = () => {
  password.value = '';
  hasError.value = false;
  errorMessage.value = '';
  dialogVisible.value = false;
};

const showError = (message: string) => {
  hasError.value = true;
  errorMessage.value = message;
  loading.value = false;
};

// Expose methods for parent component
defineExpose({
  showError,
  setLoading: (value: boolean) => {
    loading.value = value;
  }
});
</script>