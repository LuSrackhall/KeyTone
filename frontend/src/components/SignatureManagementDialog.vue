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
          <q-item v-for="sig in signatures" :key="sig.name" clickable @click="showEditDialog(sig)">
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
        <q-btn flat :label="$t('signature.createSignature')" color="primary" @click="showCreateDialogHandler" />
        <q-btn flat :label="$t('signature.importSignature')" color="primary" @click="showImportDialogHandler" />
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
          v-model="newSignature.cardImageFile"
          :label="$t('signature.cardImage')"
          outlined
          dense
          accept="image/*"
          class="q-mt-sm"
          clearable
          @update:model-value="handleCardImageSelect"
        >
          <template v-slot:prepend>
            <q-icon name="image" />
          </template>
          <template v-slot:hint v-if="newSignature.cardImagePreview">
            <div class="q-mt-xs">
              <img :src="newSignature.cardImagePreview" style="max-width: 100%; max-height: 100px" />
            </div>
          </template>
        </q-file>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" @click="cancelCreateDialog" />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="createSignature" />
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
          v-model="editSignature.intro"
          :label="$t('signature.signatureIntro')"
          outlined
          dense
          type="textarea"
        />
        <q-file
          v-model="editSignature.cardImageFile"
          :label="$t('signature.cardImage')"
          outlined
          dense
          accept="image/*"
          class="q-mt-sm"
          clearable
          @update:model-value="handleEditCardImageSelect"
        >
          <template v-slot:prepend>
            <q-icon name="image" />
          </template>
          <template v-slot:hint v-if="editSignature.cardImagePreview">
            <div class="q-mt-xs">
              <img :src="editSignature.cardImagePreview" style="max-width: 100%; max-height: 100px" />
            </div>
          </template>
        </q-file>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('KeyToneAlbum.cancel')" color="primary" @click="cancelEditDialog" />
        <q-btn flat :label="$t('KeyToneAlbum.confirm')" color="primary" @click="updateSignature" />
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

const newSignature = ref({
  name: '',
  intro: '',
  cardImageFile: null as File | null,
  cardImagePreview: null as string | null,
});

const editSignature = ref({
  name: '',
  intro: '',
  cardImageFile: null as File | null,
  cardImagePreview: null as string | null,
  originalKey: '',
  originalData: null as any,
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
const showCreateDialogHandler = () => {
  newSignature.value = { 
    name: '', 
    intro: '',
    cardImageFile: null,
    cardImagePreview: null,
  };
  createDialog.value = true;
};

// Show edit dialog
const showEditDialog = async (sig: any) => {
  // Find the signature key in signature_manager
  try {
    const response = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    if (response.data.message === 'ok' && response.data.value) {
      const signatureManager = response.data.value;
      
      // Find the key for this signature
      for (const [key, value] of Object.entries(signatureManager)) {
        if ((value as any)?.name === sig.name) {
          editSignature.value = {
            name: sig.name,
            intro: sig.intro || '',
            cardImageFile: null,
            cardImagePreview: sig.cardImagePath || null,
            originalKey: key,
            originalData: value,
          };
          editDialog.value = true;
          break;
        }
      }
    }
  } catch (error) {
    console.error('Failed to load signature for editing:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.loadFailed'),
    });
  }
};

// Show import dialog
const showImportDialogHandler = () => {
  importFile.value = null;
  importDialog.value = true;
};

// Cancel create dialog
const cancelCreateDialog = () => {
  if (newSignature.value.cardImagePreview) {
    URL.revokeObjectURL(newSignature.value.cardImagePreview);
  }
  newSignature.value = {
    name: '',
    intro: '',
    cardImageFile: null,
    cardImagePreview: null,
  };
  createDialog.value = false;
};

// Cancel edit dialog
const cancelEditDialog = () => {
  if (editSignature.value.cardImagePreview && editSignature.value.cardImageFile) {
    URL.revokeObjectURL(editSignature.value.cardImagePreview);
  }
  editSignature.value = {
    name: '',
    intro: '',
    cardImageFile: null,
    cardImagePreview: null,
    originalKey: '',
    originalData: null,
  };
  editDialog.value = false;
};

// Handle card image selection for create
const handleCardImageSelect = (file: File | null) => {
  if (newSignature.value.cardImagePreview) {
    URL.revokeObjectURL(newSignature.value.cardImagePreview);
  }
  
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      newSignature.value.cardImagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  } else {
    newSignature.value.cardImagePreview = null;
  }
};

// Handle card image selection for edit
const handleEditCardImageSelect = (file: File | null) => {
  if (editSignature.value.cardImagePreview && editSignature.value.cardImageFile) {
    URL.revokeObjectURL(editSignature.value.cardImagePreview);
  }
  
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      editSignature.value.cardImagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  } else if (!file && editSignature.value.originalData?.cardImagePath) {
    // If cleared and original had image, restore original preview
    editSignature.value.cardImagePreview = editSignature.value.originalData.cardImagePath;
  } else {
    editSignature.value.cardImagePreview = null;
  }
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

    // Add card image if provided
    if (newSignature.value.cardImagePreview) {
      signatureData.cardImagePath = newSignature.value.cardImagePreview;
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

    cancelCreateDialog();
    loadSignatures();
  } catch (error) {
    console.error('Failed to create signature:', error);
    $q.notify({
      type: 'negative',
      message: t('signature.notify.createFailed'),
    });
  }
};

// Update signature
const updateSignature = async () => {
  try {
    // Get existing signatures
    const getResponse = await api.get('/store/get', {
      params: { key: 'signature_manager' }
    });

    const signatureManager = getResponse.data.value || {};
    
    // Update the signature data
    if (signatureManager[editSignature.value.originalKey]) {
      signatureManager[editSignature.value.originalKey] = {
        ...signatureManager[editSignature.value.originalKey],
        intro: editSignature.value.intro,
      };

      // Update card image if changed
      if (editSignature.value.cardImageFile) {
        signatureManager[editSignature.value.originalKey].cardImagePath = editSignature.value.cardImagePreview;
      }
    }

    // Save back
    await api.post('/store/set', {
      key: 'signature_manager',
      value: signatureManager,
    });

    $q.notify({
      type: 'positive',
      message: t('signature.notify.updateSuccess'),
    });

    cancelEditDialog();
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
        if (error.response?.data?.error === 'exists_without_overwrite') {
          // Show overwrite confirmation
          $q.dialog({
            title: t('signature.dialog.overwriteTitle'),
            message: t('signature.dialog.overwriteMessage', { name: error.response.data.name || '' }),
            cancel: true,
            persistent: true,
            ok: {
              label: t('signature.dialog.overwrite'),
              color: 'primary',
            },
            cancel: {
              label: t('KeyToneAlbum.cancel'),
              color: 'primary',
              flat: true,
            },
          }).onOk(async () => {
            try {
              // Retry with overwrite
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
    ok: {
      label: t('KeyToneAlbum.confirm'),
      color: 'negative',
    },
    cancel: {
      label: t('KeyToneAlbum.cancel'),
      color: 'primary',
      flat: true,
    },
  }).onOk(async () => {
    try {
      // Get existing signatures
      const getResponse = await api.get('/store/get', {
        params: { key: 'signature_manager' }
      });

      if (getResponse.data.message !== 'ok') {
        throw new Error('Failed to get signature manager');
      }

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
      const setResponse = await api.post('/store/set', {
        key: 'signature_manager',
        value: signatureManager,
      });

      if (setResponse.data.message !== 'ok') {
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
</style>
