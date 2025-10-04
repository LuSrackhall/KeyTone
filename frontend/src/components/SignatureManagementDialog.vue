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
  <q-dialog v-model="showDialog" persistent backdrop-filter="invert(70%)">
    <q-card style="max-width: 360px; width: 100%; max-height: 420px">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.title') }}</div>
      </q-card-section>

      <q-separator />

      <q-card-section style="max-height: 280px; overflow-y: auto">
        <div v-if="signatures.length === 0" class="text-center text-grey-6 q-pa-md">
          <div>{{ $t('signature.emptyState.noSignatures') }}</div>
          <div class="text-caption">{{ $t('signature.emptyState.createFirst') }}</div>
        </div>

        <q-list v-else bordered separator>
          <q-item v-for="sig in signatures" :key="sig.name" clickable @click="editSignature(sig)">
            <q-item-section>
              <q-item-label>{{ sig.name }}</q-item-label>
              <q-item-label caption>{{ sig.intro || $t('signature.signatureIntro') }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row q-gutter-xs">
                <q-btn flat dense icon="file_download" size="sm" @click.stop="exportSignature(sig.name)">
                  <q-tooltip>{{ $t('signature.exportSignature') }}</q-tooltip>
                </q-btn>
                <q-btn flat dense icon="delete" size="sm" @click.stop="deleteSignature(sig.name)">
                  <q-tooltip>{{ $t('signature.deleteSignature') }}</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <q-separator />

      <q-card-actions align="right">
        <q-btn flat :label="$t('signature.createSignature')" color="primary" @click="showCreateDialog" />
        <q-btn flat :label="$t('signature.importSignature')" color="primary" @click="showImportDialog" />
        <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Create Signature Dialog -->
  <q-dialog v-model="createDialog" persistent backdrop-filter="invert(70%)">
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.dialog.createTitle') }}</div>
      </q-card-section>

      <q-card-section>
        <q-input
          v-model="newSignature.name"
          :label="$t('signature.signatureName')"
          outlined
          dense
          :error="!newSignature.name"
          :error-message="$t('signature.notify.nameRequired')"
        />
        <q-input
          v-model="newSignature.intro"
          :label="$t('signature.signatureIntro')"
          outlined
          dense
          type="textarea"
          class="q-mt-sm"
        />
        <q-file
          v-model="newSignature.cardImage"
          :label="$t('signature.cardImage')"
          outlined
          dense
          accept="image/*"
          class="q-mt-sm"
          @update:model-value="handleCardImageSelect"
        >
          <template v-slot:prepend>
            <q-icon name="image" />
          </template>
        </q-file>
        <div v-if="newSignature.cardImagePreview" class="q-mt-sm">
          <div class="text-caption q-mb-xs">{{ $t('signature.cardImagePreview') }}</div>
          <img :src="newSignature.cardImagePreview" class="card-image-preview" @click="showImageDialog = true" />
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" v-close-popup />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="createSignature" />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Import Signature Dialog -->
  <q-dialog v-model="importDialog" persistent backdrop-filter="invert(70%)">
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.dialog.importTitle') }}</div>
      </q-card-section>

      <q-card-section>
        <q-file
          v-model="importFile"
          :label="$t('signature.dialog.selectFile')"
          outlined
          dense
          accept=".ktsign"
          @update:model-value="handleFileSelect"
        >
          <template v-slot:prepend>
            <q-icon name="attach_file" />
          </template>
        </q-file>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" v-close-popup />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="importSignature" :disable="!importFile" />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Edit Signature Dialog -->
  <q-dialog v-model="editDialog" persistent backdrop-filter="invert(70%)">
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.dialog.editTitle') }}</div>
      </q-card-section>

      <q-card-section>
        <q-input
          v-model="editingSignature.intro"
          :label="$t('signature.signatureIntro')"
          outlined
          dense
          type="textarea"
        />
        <q-file
          v-model="editingSignature.cardImage"
          :label="$t('signature.cardImage')"
          outlined
          dense
          accept="image/*"
          class="q-mt-sm"
          @update:model-value="handleEditCardImageSelect"
        >
          <template v-slot:prepend>
            <q-icon name="image" />
          </template>
        </q-file>
        <div v-if="editingSignature.cardImagePreview" class="q-mt-sm">
          <div class="text-caption q-mb-xs">{{ $t('signature.cardImagePreview') }}</div>
          <img :src="editingSignature.cardImagePreview" class="card-image-preview" @click="showEditImageDialog = true" />
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" v-close-popup />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="updateSignature" />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Image Preview Dialog (for create) -->
  <q-dialog v-model="showImageDialog" backdrop-filter="invert(70%)">
    <q-card>
      <q-card-section class="q-pa-none">
        <img :src="newSignature.cardImagePreview" style="max-width: 90vw; max-height: 90vh; width: auto; height: auto;" />
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Image Preview Dialog (for edit) -->
  <q-dialog v-model="showEditImageDialog" backdrop-filter="invert(70%)">
    <q-card>
      <q-card-section class="q-pa-none">
        <img :src="editingSignature.cardImagePreview" style="max-width: 90vw; max-height: 90vh; width: auto; height: auto;" />
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Overwrite Confirmation Dialog -->
  <q-dialog v-model="overwriteDialog" persistent backdrop-filter="invert(70%)">
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.dialog.overwriteTitle') }}</div>
      </q-card-section>

      <q-card-section>
        {{ $t('signature.dialog.overwriteMessage', { name: overwriteSignatureName }) }}
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" @click="cancelOverwrite" />
        <q-btn flat :label="$t('signature.dialog.overwrite')" color="primary" @click="confirmOverwrite" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { api } from 'src/boot/axios';

const $q = useQuasar();
const { t } = useI18n();

const showDialog = ref(false);
const createDialog = ref(false);
const importDialog = ref(false);
const editDialog = ref(false);
const showImageDialog = ref(false);
const showEditImageDialog = ref(false);
const overwriteDialog = ref(false);
const signatures = ref<any[]>([]);
const importFile = ref<File | null>(null);
const overwriteSignatureName = ref('');
let pendingImportData: { fileName: string; blobBase64: string } | null = null;

const newSignature = ref({
  name: '',
  intro: '',
  cardImage: null as File | null,
  cardImagePreview: '',
});

const editingSignature = ref({
  name: '',
  intro: '',
  cardImage: null as File | null,
  cardImagePreview: '',
  originalKey: '',
});

// Open the dialog
const open = () => {
  showDialog.value = true;
  loadSignatures();
};

// Load signatures from backend
const loadSignatures = async () => {
  try {
    const response = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    if (response.data.message === 'ok' && response.data.value) {
      const signatureManager = response.data.value;
      signatures.value = Object.values(signatureManager).filter((v: any) => v && v.name);
    } else {
      signatures.value = [];
    }
  } catch (error) {
    console.error('Failed to load signatures:', error);
    signatures.value = [];
  }
};

// Show create dialog
const showCreateDialog = () => {
  newSignature.value = { name: '', intro: '', cardImage: null, cardImagePreview: '' };
  createDialog.value = true;
};

// Show import dialog
const showImportDialog = () => {
  importFile.value = null;
  importDialog.value = true;
};

// Create signature
const createSignature = async () => {
  if (!newSignature.value.name) {
    $q.notify({
      type: 'negative',
      message: t('signature.notify.nameRequired'),
    });
    return;
  }

  try {
    // Generate a simple protect code (in real implementation, use nanoid or similar)
    const protectCode = 'pc_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    const encryptedKey = 'key_' + protectCode; // Simple encryption placeholder

    const signatureData: any = {
      name: newSignature.value.name,
      intro: newSignature.value.intro,
      createdAt: new Date().toISOString(),
    };

    // Handle card image if present
    if (newSignature.value.cardImage) {
      signatureData.cardImagePath = await saveCardImage(newSignature.value.cardImage, encryptedKey);
    }

    // Get existing signatures
    const getResponse = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    const signatureManager = getResponse.data.value || {};
    signatureManager[encryptedKey] = signatureData;

    // Save back
    await api.post('/store/set', {
      key: 'signature_manager',
      value: signatureManager,
    });

    $q.notify({
      type: 'positive',
      message: t('signature.notify.createSuccess'),
    });

    createDialog.value = false;
    loadSignatures();
  } catch (error) {
    console.error('Failed to create signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.createFailed'),
    });
  }
};

// Save card image to storage
const saveCardImage = async (file: File, key: string): Promise<string> => {
  // Convert file to base64 and save it
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = async () => {
      try {
        const base64 = reader.result as string;
        // Store in a separate key for images
        await api.post('/store/set', {
          key: `signature_card_image_${key}`,
          value: base64,
        });
        resolve(`signature_card_image_${key}`);
      } catch (error) {
        reject(error);
      }
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};

// Handle card image selection for create
const handleCardImageSelect = (file: File | null) => {
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      newSignature.value.cardImagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  } else {
    newSignature.value.cardImagePreview = '';
  }
};

// Handle card image selection for edit
const handleEditCardImageSelect = (file: File | null) => {
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      editingSignature.value.cardImagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  } else if (!editingSignature.value.cardImagePreview) {
    // Only clear if there's no existing preview
    editingSignature.value.cardImagePreview = '';
  }
};

// Edit signature
const editSignature = async (sig: any) => {
  // Find the key for this signature
  const getResponse = await api.get('/store/get', {
    params: { key: 'signature_manager' }
  });
  
  const signatureManager = getResponse.data.value || {};
  let foundKey = '';
  
  for (const [key, value] of Object.entries(signatureManager)) {
    if ((value as any)?.name === sig.name) {
      foundKey = key;
      break;
    }
  }
  
  if (!foundKey) return;
  
  editingSignature.value = {
    name: sig.name,
    intro: sig.intro || '',
    cardImage: null,
    cardImagePreview: '',
    originalKey: foundKey,
  };
  
  // Load existing card image if present
  if (sig.cardImagePath) {
    try {
      const imageResponse = await api.get('/store/get', {
        params: { key: sig.cardImagePath }
      });
      if (imageResponse.data.value) {
        editingSignature.value.cardImagePreview = imageResponse.data.value;
      }
    } catch (error) {
      console.error('Failed to load card image:', error);
    }
  }
  
  editDialog.value = true;
};

// Update signature
const updateSignature = async () => {
  try {
    // Get existing signatures
    const getResponse = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    const signatureManager = getResponse.data.value || {};
    const existingData = signatureManager[editingSignature.value.originalKey] || {};
    
    const signatureData: any = {
      ...existingData,
      intro: editingSignature.value.intro,
    };

    // Handle card image if changed
    if (editingSignature.value.cardImage) {
      signatureData.cardImagePath = await saveCardImage(editingSignature.value.cardImage, editingSignature.value.originalKey);
    }

    signatureManager[editingSignature.value.originalKey] = signatureData;

    // Save back
    await api.post('/store/set', {
      key: 'signature_manager',
      value: signatureManager,
    });

    $q.notify({
      type: 'positive',
      message: t('signature.notify.updateSuccess'),
    });

    editDialog.value = false;
    loadSignatures();
  } catch (error) {
    console.error('Failed to update signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.updateFailed'),
    });
  }
};

// Export signature
const exportSignature = async (name: string) => {
  try {
    const response = await api.post(`/sdk/signatures/${encodeURIComponent(name)}/export`, {});

    if (response.data.blobBase64) {
      // Convert base64 to blob and download
      const byteCharacters = atob(response.data.blobBase64);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);
      const blob = new Blob([byteArray], { type: 'application/octet-stream' });

      // Check if File System Access API is available
      if (typeof window.showSaveFilePicker === 'function') {
        try {
          const handle = await (window as any).showSaveFilePicker({
            suggestedName: response.data.fileNameSuggested || `${name}.ktsign`,
            types: [
              {
                description: 'KeyTone Signature File',
                accept: { 'application/octet-stream': ['.ktsign'] },
              },
            ],
          });

          const writable = await handle.createWritable();
          await writable.write(blob);
          await writable.close();

          // Only show success after file is actually written
          $q.notify({
            type: 'positive',
            message: t('signature.notify.exportSuccess'),
          });
        } catch (err: any) {
          // User cancelled - don't show error
          if (err.name === 'AbortError') {
            return;
          }
          throw err;
        }
      } else {
        // Fallback to traditional download
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = response.data.fileNameSuggested || `${name}.ktsign`;
        link.click();
        window.URL.revokeObjectURL(url);

        // For legacy download, show success immediately after triggering
        $q.notify({
          type: 'positive',
          message: t('signature.notify.exportSuccess'),
        });
      }
    }
  } catch (error) {
    console.error('Failed to export signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.exportFailed'),
    });
  }
};

// Import signature
const importSignature = async () => {
  if (!importFile.value) return;

  try {
    // Read file as base64
    const reader = new FileReader();
    reader.onload = async (e) => {
      const arrayBuffer = e.target?.result as ArrayBuffer;
      const bytes = new Uint8Array(arrayBuffer);
      const base64 = btoa(String.fromCharCode(...bytes));

      try {
        const response = await api.post('/sdk/signatures/import', {
          fileName: importFile.value!.name,
          blobBase64: base64,
          overwrite: false, // First try without overwrite
        });

        $q.notify({
          type: 'positive',
          message: t('signature.notify.importSuccess'),
        });

        importDialog.value = false;
        loadSignatures();
      } catch (error: any) {
        if (error.response?.data?.error === 'exists_without_overwrite') {
          // Show overwrite confirmation dialog
          overwriteSignatureName.value = error.response.data.name || '';
          pendingImportData = {
            fileName: importFile.value!.name,
            blobBase64: base64,
          };
          overwriteDialog.value = true;
        } else {
          throw error;
        }
      }
    };
    reader.readAsArrayBuffer(importFile.value);
  } catch (error) {
    console.error('Failed to import signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.importFailed'),
    });
  }
};

// Confirm overwrite
const confirmOverwrite = async () => {
  if (!pendingImportData) return;
  
  try {
    await api.post('/sdk/signatures/import', {
      fileName: pendingImportData.fileName,
      blobBase64: pendingImportData.blobBase64,
      overwrite: true,
    });
    
    $q.notify({
      type: 'positive',
      message: t('signature.notify.importSuccess'),
    });
    
    overwriteDialog.value = false;
    importDialog.value = false;
    pendingImportData = null;
    loadSignatures();
  } catch (error) {
    console.error('Failed to import signature with overwrite:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.importFailed'),
    });
  }
};

// Cancel overwrite
const cancelOverwrite = () => {
  overwriteDialog.value = false;
  pendingImportData = null;
};

// Delete signature
const deleteSignature = async (name: string) => {
  $q.dialog({
    title: t('signature.deleteSignature'),
    message: t('KeyToneAlbum.notify.confirmDelete', { name }),
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    try {
      // Get existing signatures
      const getResponse = await api.get('/store/get', {
        params: { key: 'signature_manager' }
      });

      const signatureManager = getResponse.data.value || {};
      
      // Find and remove the signature
      let found = false;
      for (const [key, value] of Object.entries(signatureManager)) {
        if ((value as any)?.name === name) {
          delete signatureManager[key];
          found = true;
          break;
        }
      }

      if (!found) {
        throw new Error('Signature not found');
      }

      // Save back
      const saveResponse = await api.post('/store/set', {
        key: 'signature_manager',
        value: signatureManager,
      });

      // Check if save was successful
      if (saveResponse.data.message !== 'ok') {
        throw new Error('Failed to save signature manager');
      }

      $q.notify({
        type: 'positive',
        message: t('signature.notify.deleteSuccess'),
      });

      loadSignatures();
    } catch (error) {
      console.error('Failed to delete signature:', error);
      $q.notify({
        type: 'negative',
        message: t('signature.notify.deleteFailed'),
      });
    }
  });
};

const handleFileSelect = (file: File | null) => {
  if (file && !file.name.endsWith('.ktsign')) {
    $q.notify({
      type: 'warning',
      message: t('signature.notify.invalidFormat'),
    });
    importFile.value = null;
  }
};

// Expose open method
defineExpose({ open });
</script>

<style scoped>
.q-card {
  border-radius: 8px;
}

.card-image-preview {
  max-width: 100%;
  max-height: 200px;
  object-fit: contain;
  cursor: pointer;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
}
</style>
