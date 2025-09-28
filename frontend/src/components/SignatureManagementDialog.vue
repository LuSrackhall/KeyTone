<!--
/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
-->

<template>
  <q-dialog v-model="showDialog" backdrop-filter="invert(70%)" persistent :style="{ '--i18n_fontSize': i18n_fontSize }">
    <q-card class="signature-management-dialog-card">
      <!-- Header -->
      <q-card-section class="signature-dialog-header">
        <div class="text-h6 text-weight-medium">{{ $t('signatureManagementDialog.title') }}</div>
        <div class="text-caption text-grey-6 q-mt-xs">{{ $t('signatureManagementDialog.subtitle') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-lg">
        <div class="signature-options">
          <!-- Create Signature Option -->
          <q-card 
            flat 
            bordered 
            class="signature-option-card cursor-pointer" 
            @click="handleCreateSignature"
          >
            <q-card-section class="flex items-center q-pa-md">
              <q-icon name="add_circle_outline" size="2.5rem" color="primary" class="q-mr-md" />
              <div class="flex-1">
                <div class="text-subtitle1 text-weight-medium">{{ $t('signatureManagementDialog.createSignature') }}</div>
                <div class="text-caption text-grey-6 q-mt-xs">{{ $t('signatureManagementDialog.createSignatureDesc') }}</div>
              </div>
              <q-icon name="chevron_right" size="1.2rem" color="grey-5" />
            </q-card-section>
          </q-card>

          <!-- Manage Signatures Option -->
          <q-card 
            flat 
            bordered 
            class="signature-option-card cursor-pointer q-mt-md" 
            @click="handleManageSignature"
          >
            <q-card-section class="flex items-center q-pa-md">
              <q-icon name="manage_accounts" size="2.5rem" color="primary" class="q-mr-md" />
              <div class="flex-1">
                <div class="text-subtitle1 text-weight-medium">{{ $t('signatureManagementDialog.manageSignature') }}</div>
                <div class="text-caption text-grey-6 q-mt-xs">{{ $t('signatureManagementDialog.manageSignatureDesc') }}</div>
              </div>
              <q-icon name="chevron_right" size="1.2rem" color="grey-5" />
            </q-card-section>
          </q-card>
        </div>
      </q-card-section>

      <!-- Footer -->
      <q-card-actions align="right" class="q-pa-md">
        <q-btn flat :label="$t('signatureManagementDialog.cancel')" color="grey" @click="cancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

interface Props {
  modelValue: boolean;
  i18nFontSize?: string;
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'createSignature'): void;
  (e: 'manageSignature'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<Props>(), {
  i18nFontSize: '1rem',
});

const emit = defineEmits<Emits>();
const { t } = useI18n();

// Dialog state
const showDialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
});

const i18n_fontSize = computed(() => props.i18nFontSize);

const handleCreateSignature = () => {
  emit('createSignature');
};

const handleManageSignature = () => {
  emit('manageSignature');
};

const cancel = () => {
  emit('cancel');
};
</script>

<style lang="scss" scoped>
.signature-management-dialog-card {
  width: 90vw;
  max-width: 480px;
  min-height: 300px;
}

.signature-dialog-header {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.signature-options {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.signature-option-card {
  transition: all 0.2s ease;
  border-radius: 8px;
  
  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }
  
  &:active {
    transform: translateY(0);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  }
}

/* Ensure inputs work properly on smaller screens */
@media (max-width: 600px) {
  .signature-management-dialog-card {
    width: 95vw;
    max-width: none;
  }
}
</style>