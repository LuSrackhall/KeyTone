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
          <q-item v-for="sig in signatures" :key="sig.name" clickable @click="openEditDialog(sig)">
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
  <q-dialog v-model="createDialog" persistent>
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
        
        <!-- Card Image Upload -->
        <div class="q-mt-md">
          <q-file
            v-model="cardImageFile"
            :label="$t('signature.cardImage')"
            outlined
            dense
            accept="image/*"
            @update:model-value="handleCardImageSelect"
          >
            <template v-slot:prepend>
              <q-icon name="image" />
            </template>
            <template v-slot:append v-if="newSignature.cardImageData">
              <q-icon name="close" @click.stop="cardImageFile = null; handleCardImageSelect(null)" class="cursor-pointer" />
            </template>
          </q-file>
          
          <!-- Image Preview -->
          <div v-if="newSignature.cardImageData" class="q-mt-sm">
            <div class="text-caption q-mb-xs">{{ $t('signature.preview') }}</div>
            <div class="relative inline-block">
              <img 
                :src="newSignature.cardImageData" 
                alt="Card preview"
                class="rounded cursor-pointer"
                style="max-width: 200px; max-height: 150px; object-fit: contain;"
                @click="showImagePreview(newSignature.cardImageData)"
              />
              <q-tooltip>{{ $t('signature.clickToZoom') }}</q-tooltip>
            </div>
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" v-close-popup />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="createSignature" />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Import Signature Dialog -->
  <q-dialog v-model="importDialog" persistent>
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

  <!-- Image Preview Dialog -->
  <q-dialog v-model="imagePreviewDialog">
    <q-card style="max-width: 90vw; max-height: 90vh">
      <q-card-section class="q-pa-none">
        <img 
          :src="previewImageUrl" 
          alt="Preview"
          style="max-width: 100%; max-height: 90vh; object-fit: contain;"
        />
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Edit Signature Dialog -->
  <q-dialog v-model="editDialog" persistent>
    <q-card style="max-width: 360px; width: 100%">
      <q-card-section>
        <div class="text-h6">{{ $t('signature.dialog.editTitle') }}</div>
      </q-card-section>

      <q-card-section>
        <q-input
          v-model="editSignature.name"
          :label="$t('signature.signatureName')"
          outlined
          dense
          :error="!editSignature.name"
          :error-message="$t('signature.notify.nameRequired')"
          disable
        />
        <q-input
          v-model="editSignature.intro"
          :label="$t('signature.signatureIntro')"
          outlined
          dense
          type="textarea"
          class="q-mt-sm"
        />
        
        <!-- Card Image Upload -->
        <div class="q-mt-md">
          <q-file
            v-model="editCardImageFile"
            :label="$t('signature.cardImage')"
            outlined
            dense
            accept="image/*"
            @update:model-value="handleEditCardImageSelect"
          >
            <template v-slot:prepend>
              <q-icon name="image" />
            </template>
            <template v-slot:append v-if="editSignature.cardImageData">
              <q-icon name="close" @click.stop="editCardImageFile = null; handleEditCardImageSelect(null)" class="cursor-pointer" />
            </template>
          </q-file>
          
          <!-- Image Preview -->
          <div v-if="editSignature.cardImageData" class="q-mt-sm">
            <div class="text-caption q-mb-xs">{{ $t('signature.preview') }}</div>
            <div class="relative inline-block">
              <img 
                :src="editSignature.cardImageData" 
                alt="Card preview"
                class="rounded cursor-pointer"
                style="max-width: 200px; max-height: 150px; object-fit: contain;"
                @click="showImagePreview(editSignature.cardImageData)"
              />
              <q-tooltip>{{ $t('signature.clickToZoom') }}</q-tooltip>
            </div>
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" v-close-popup />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="updateSignature" />
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
const editDialog = ref(false);
const importDialog = ref(false);
const signatures = ref<any[]>([]);
const importFile = ref<File | null>(null);
const imagePreviewDialog = ref(false);
const previewImageUrl = ref('');
const cardImageFile = ref<File | null>(null);
const editCardImageFile = ref<File | null>(null);

const editSignature = ref({
  originalName: '',
  name: '',
  intro: '',
  cardImagePath: '',
  cardImageData: '',
});

const newSignature = ref({
  name: '',
  intro: '',
  cardImagePath: '',
  cardImageData: '',
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
  newSignature.value = { name: '', intro: '', cardImagePath: '', cardImageData: '' };
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

    const signatureData = {
      name: newSignature.value.name,
      intro: newSignature.value.intro,
      cardImagePath: newSignature.value.cardImagePath,
      cardImageData: newSignature.value.cardImageData,
      createdAt: new Date().toISOString(),
    };

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
    cardImageFile.value = null;
    loadSignatures();
  } catch (error) {
    console.error('Failed to create signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.createFailed'),
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

      // Trigger download
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = response.data.fileNameSuggested || `${name}.ktsign`;
      link.click();
      window.URL.revokeObjectURL(url);

      $q.notify({
        type: 'positive',
        message: t('signature.notify.exportSuccess'),
      });
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
        // First try without overwrite
        const response = await api.post('/sdk/signatures/import', {
          fileName: importFile.value!.name,
          blobBase64: base64,
          overwrite: false,
        });

        $q.notify({
          type: 'positive',
          message: t('signature.notify.importSuccess'),
        });

        importDialog.value = false;
        loadSignatures();
      } catch (error: any) {
        if (error.response?.status === 409 || error.response?.data?.error === 'exists_without_overwrite') {
          // Signature exists, show overwrite confirmation
          const signatureName = error.response.data.name || importFile.value!.name;
          
          $q.dialog({
            title: t('signature.dialog.overwriteTitle'),
            message: t('signature.dialog.overwriteMessage', { name: signatureName }),
            cancel: {
              label: t('signature.dialog.cancel'),
              flat: true,
              color: 'primary',
            },
            ok: {
              label: t('signature.dialog.overwrite'),
              flat: true,
              color: 'negative',
            },
            persistent: true,
          }).onOk(async () => {
            // Retry with overwrite
            try {
              await api.post('/sdk/signatures/import', {
                fileName: importFile.value!.name,
                blobBase64: base64,
                overwrite: true,
              });
              
              $q.notify({
                type: 'positive',
                message: t('signature.notify.importSuccess'),
              });

              importDialog.value = false;
              loadSignatures();
            } catch (retryError) {
              console.error('Failed to import signature with overwrite:', retryError);
              $q.notify({
                type: 'negative',
                message: t('signature.notify.importFailed'),
              });
            }
          });
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
        $q.notify({
          type: 'negative',
          message: t('signature.notify.deleteFailed'),
        });
        return;
      }

      // Save back
      const setResponse = await api.post('/store/set', {
        key: 'signature_manager',
        value: signatureManager,
      });

      // Check if save was successful
      if (setResponse.data.message !== 'ok') {
        throw new Error('Save failed: ' + setResponse.data.message);
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

// Handle card image selection
const handleCardImageSelect = (file: File | null) => {
  if (!file) {
    newSignature.value.cardImagePath = '';
    newSignature.value.cardImageData = '';
    return;
  }

  // Check if it's an image
  if (!file.type.startsWith('image/')) {
    $q.notify({
      type: 'warning',
      message: t('signature.notify.invalidImageFormat'),
    });
    return;
  }

  // Read file as base64
  const reader = new FileReader();
  reader.onload = (e) => {
    const result = e.target?.result as string;
    newSignature.value.cardImageData = result;
    newSignature.value.cardImagePath = file.name;
  };
  reader.readAsDataURL(file);
};

// Show image preview in zoom dialog
const showImagePreview = (imageData: string) => {
  previewImageUrl.value = imageData;
  imagePreviewDialog.value = true;
};

// Open edit dialog
const openEditDialog = (sig: any) => {
  editSignature.value = {
    originalName: sig.name,
    name: sig.name,
    intro: sig.intro || '',
    cardImagePath: sig.cardImagePath || '',
    cardImageData: sig.cardImageData || '',
  };
  editCardImageFile.value = null;
  editDialog.value = true;
};

// Handle edit card image selection
const handleEditCardImageSelect = (file: File | null) => {
  if (!file) {
    editSignature.value.cardImagePath = '';
    editSignature.value.cardImageData = '';
    return;
  }

  // Check if it's an image
  if (!file.type.startsWith('image/')) {
    $q.notify({
      type: 'warning',
      message: t('signature.notify.invalidImageFormat'),
    });
    return;
  }

  // Read file as base64
  const reader = new FileReader();
  reader.onload = (e) => {
    const result = e.target?.result as string;
    editSignature.value.cardImageData = result;
    editSignature.value.cardImagePath = file.name;
  };
  reader.readAsDataURL(file);
};

// Update signature
const updateSignature = async () => {
  if (!editSignature.value.name) {
    $q.notify({
      type: 'negative',
      message: t('signature.notify.nameRequired'),
    });
    return;
  }

  try {
    // Get existing signatures
    const getResponse = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    const signatureManager = getResponse.data.value || {};
    
    // Find and update the signature
    let found = false;
    for (const [key, value] of Object.entries(signatureManager)) {
      if ((value as any)?.name === editSignature.value.originalName) {
        signatureManager[key] = {
          name: editSignature.value.name,
          intro: editSignature.value.intro,
          cardImagePath: editSignature.value.cardImagePath,
          cardImageData: editSignature.value.cardImageData,
          createdAt: (value as any).createdAt || new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        found = true;
        break;
      }
    }

    if (!found) {
      $q.notify({
        type: 'negative',
        message: t('signature.notify.updateFailed'),
      });
      return;
    }

    // Save back
    const setResponse = await api.post('/store/set', {
      key: 'signature_manager',
      value: signatureManager,
    });

    // Check if save was successful
    if (setResponse.data.message !== 'ok') {
      throw new Error('Save failed: ' + setResponse.data.message);
    }

    $q.notify({
      type: 'positive',
      message: t('signature.notify.updateSuccess'),
    });

    editDialog.value = false;
    editCardImageFile.value = null;
    loadSignatures();
  } catch (error) {
    console.error('Failed to update signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.updateFailed'),
    });
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
