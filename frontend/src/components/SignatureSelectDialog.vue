<!--
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
-->

<template>
  <q-dialog v-model="showDialog" persistent>
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.exportFlow.selectSignature') }}</div>
      </q-card-section>

      <q-separator />

      <q-card-section>
        <div class="text-body2 q-mb-md">{{ $t('signature.exportFlow.selectPrompt') }}</div>
        
        <q-select
          v-model="selectedSignature"
          :options="signatureOptions"
          :label="$t('signature.signatureName')"
          outlined
          dense
          option-label="name"
          option-value="name"
          :error="required && !selectedSignature"
          :error-message="$t('signature.exportFlow.signatureRequired')"
        >
          <template v-slot:no-option>
            <q-item>
              <q-item-section class="text-grey">
                {{ $t('signature.emptyState.noSignatures') }}
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </q-card-section>

      <q-separator />

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" @click="handleCancel" />
        <q-btn 
          flat 
          :label="$t('KeyToneAlbum.confirm')" 
          color="primary" 
          @click="handleConfirm" 
          :disable="required && !selectedSignature"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { api } from 'src/boot/axios';

const $q = useQuasar();
const { t } = useI18n();

const showDialog = ref(false);
const selectedSignature = ref<any>(null);
const signatureOptions = ref<any[]>([]);
const required = ref(false);

let resolveCallback: ((signature: any) => void) | null = null;
let rejectCallback: (() => void) | null = null;

// Load signatures from backend
const loadSignatures = async () => {
  try {
    const response = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    if (response.data.message === 'ok' && response.data.value) {
      const signatureManager = response.data.value;
      signatureOptions.value = Object.values(signatureManager).filter((v: any) => v && v.name);
    } else {
      signatureOptions.value = [];
    }
  } catch (error) {
    console.error('Failed to load signatures:', error);
    signatureOptions.value = [];
  }
};

// Open the dialog
const open = (isRequired: boolean = false): Promise<any> => {
  return new Promise((resolve, reject) => {
    required.value = isRequired;
    selectedSignature.value = null;
    resolveCallback = resolve;
    rejectCallback = reject;
    showDialog.value = true;
    loadSignatures();
  });
};

// Handle confirm
const handleConfirm = () => {
  if (required.value && !selectedSignature.value) {
    $q.notify({
      type: 'warning',
      message: t('signature.exportFlow.signatureRequired'),
    });
    return;
  }

  showDialog.value = false;
  if (resolveCallback) {
    resolveCallback(selectedSignature.value);
  }
};

// Handle cancel
const handleCancel = () => {
  showDialog.value = false;
  if (rejectCallback) {
    rejectCallback();
  }
};

// Expose open method
defineExpose({ open });
</script>

<style scoped>
.q-card {
  border-radius: 8px;
}
</style>
