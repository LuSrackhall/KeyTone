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
  <q-dialog
    v-model="showDialog"
    backdrop-filter="invert(70%)"
    persistent
    :style="{ '--i18n_fontSize': i18n_fontSize }"
  >
    <q-card class="copyright-dialog-card">
      <!-- Fixed Header with Glass Background -->
      <div class="copyright-dialog-header">
        <q-card-section class="row items-center q-pb-none text-h6">
          {{ $t('copyrightDialog.title') }}
        </q-card-section>
        
        <q-card-section class="text-subtitle2 text-grey-7 q-pt-sm q-pb-md">
          {{ $t('copyrightDialog.subtitle') }}
        </q-card-section>
      </div>

      <!-- Scrollable Content -->
      <div class="copyright-dialog-content">
        <q-card-section class="q-py-sm">
          <!-- Author Name (Required) -->
          <q-input
            v-model="authorName"
            :label="$t('copyrightDialog.authorName')"
            :placeholder="$t('copyrightDialog.authorNamePlaceholder')"
            outlined
            stack-label
            dense
            :error="authorNameError"
            :error-message="$t('copyrightDialog.authorNameRequired')"
            @update:model-value="validateForm"
            class="q-mb-sm"
          />

          <!-- Copyright Protection Code (Required) -->
          <q-input
            v-model="protectionCode"
            :label="$t('copyrightDialog.protectionCode')"
            :placeholder="$t('copyrightDialog.protectionCodePlaceholder')"
            outlined
            stack-label
            dense
            type="password"
            :error="protectionCodeError"
            :error-message="$t('copyrightDialog.protectionCodeError')"
            @update:model-value="validateForm"
            class="q-mb-sm"
          />

          <!-- Text Contact Information (Optional) -->
          <q-input
            v-model="textContact"
            :label="$t('copyrightDialog.textContact')"
            :placeholder="$t('copyrightDialog.textContactPlaceholder')"
            outlined
            stack-label
            dense
            type="textarea"
            rows="3"
            class="q-mb-sm"
          />

          <!-- Image Contact Information (Optional) -->
          <div class="q-mb-sm">
            <q-file
              v-model="imageContactFile"
              :label="$t('copyrightDialog.imageContact')"
              outlined
              stack-label
              dense
              accept="image/*"
              :loading="isUploading"
              @update:model-value="handleImageSelect"
            >
              <template v-slot:prepend>
                <q-icon name="attach_file" />
              </template>
            </q-file>
            
            <!-- Image Preview -->
            <div v-if="imageContactPreview" class="q-mt-sm">
              <q-img
                :src="imageContactPreview"
                class="rounded-borders"
                style="max-width: 200px; max-height: 150px"
                fit="contain"
              />
              <div class="flex items-center q-mt-xs">
                <q-btn
                  flat
                  dense
                  icon="close"
                  color="negative"
                  size="sm"
                  @click="removeImage"
                >
                  {{ $t('copyrightDialog.removeImage') }}
                </q-btn>
                <div v-if="isUploading" class="q-ml-sm">
                  <q-spinner color="primary" size="sm" />
                  <span class="q-ml-xs text-caption">{{ $t('copyrightDialog.uploading') }}</span>
                </div>
                <div v-else-if="imageContactPath" class="q-ml-sm text-positive">
                  <q-icon name="check_circle" size="sm" />
                </div>
              </div>
            </div>
          </div>

          <!-- Warning for skip option -->
          <q-banner
            v-if="hasExistingCopyright && showSkipWarning"
            class="bg-orange-1 text-orange-8"
            icon="warning"
          >
            {{ $t('copyrightDialog.skipWarning') }}
          </q-banner>
        </q-card-section>
      </div>

      <!-- Fixed Footer with Glass Background -->
      <div class="copyright-dialog-footer">
        <q-card-actions align="right" class="q-pa-md">
          <q-btn
            flat
            :label="$t('copyrightDialog.cancel')"
            color="grey"
            @click="cancel"
          />
          
          <q-btn
            v-if="!hasExistingCopyright"
            flat
            :label="$t('copyrightDialog.skipAndExport')"
            color="orange"
            @click="skipAndExport"
          />
          
          <q-btn
            :label="$t('copyrightDialog.confirmAndExport')"
            color="primary"
            :disable="!isFormValid"
            @click="confirmAndExport"
          />
        </q-card-actions>
      </div>
    </q-card>
  </q-dialog>

  <!-- Skip Confirmation Dialog -->
  <q-dialog v-model="showSkipConfirmation" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar icon="warning" color="orange" text-color="white" />
        <span class="q-ml-sm">{{ $t('copyrightDialog.skipConfirmTitle') }}</span>
      </q-card-section>

      <q-card-section>
        {{ $t('copyrightDialog.skipConfirmMessage') }}
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('copyrightDialog.cancel')" color="primary" @click="showSkipConfirmation = false" />
        <q-btn :label="$t('copyrightDialog.confirmSkip')" color="orange" @click="confirmSkip" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { UploadCopyrightImage } from 'src/boot/query/keytonePkg-query';

interface CopyrightData {
  authorName: string;
  textContact: string;
  imageContactPath?: string;
  protectionCode: string;
}

interface Props {
  modelValue: boolean;
  hasExistingCopyright?: boolean;
  i18nFontSize?: string;
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'confirm', data: CopyrightData): void;
  (e: 'skip'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<Props>(), {
  hasExistingCopyright: false,
  i18nFontSize: '1rem',
});

const emit = defineEmits<Emits>();
const { t } = useI18n();

// Dialog state
const showDialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
});

const showSkipConfirmation = ref(false);
const showSkipWarning = ref(false);

// Form data
const authorName = ref('');
const protectionCode = ref('');
const textContact = ref('');
const imageContactFile = ref<File | null>(null);
const imageContactPreview = ref<string | null>(null);
const imageContactPath = ref<string>(''); // Store the uploaded image path
const isUploading = ref(false);

// Form validation
const authorNameError = ref(false);
const protectionCodeError = ref(false);

const i18n_fontSize = computed(() => props.i18nFontSize);

const isFormValid = computed(() => {
  return authorName.value.trim().length > 0 && 
         protectionCode.value.length >= 6 && 
         !authorNameError.value && 
         !protectionCodeError.value;
});

const validateForm = () => {
  authorNameError.value = authorName.value.trim().length === 0;
  protectionCodeError.value = protectionCode.value.length < 6;
};

const handleImageSelect = async (file: File | null) => {
  if (file) {
    // Set preview immediately
    const url = URL.createObjectURL(file);
    imageContactPreview.value = url;
    
    // Upload the file
    isUploading.value = true;
    try {
      const result = await UploadCopyrightImage(file);
      if (result.success) {
        imageContactPath.value = result.path;
        console.log('Image uploaded successfully:', result.path);
      } else {
        console.error('Image upload failed:', result.error);
        // Still keep the preview but clear the path
        imageContactPath.value = '';
      }
    } catch (error) {
      console.error('Error uploading image:', error);
      imageContactPath.value = '';
    } finally {
      isUploading.value = false;
    }
  } else {
    imageContactPreview.value = null;
    imageContactPath.value = '';
  }
};

const removeImage = () => {
  if (imageContactPreview.value) {
    URL.revokeObjectURL(imageContactPreview.value);
  }
  imageContactFile.value = null;
  imageContactPreview.value = null;
  imageContactPath.value = '';
};

const resetForm = () => {
  authorName.value = '';
  protectionCode.value = '';
  textContact.value = '';
  removeImage();
  authorNameError.value = false;
  protectionCodeError.value = false;
  showSkipWarning.value = false;
};

const cancel = () => {
  resetForm();
  emit('cancel');
};

const skipAndExport = () => {
  if (props.hasExistingCopyright) {
    showSkipWarning.value = true;
    return;
  }
  showSkipConfirmation.value = true;
};

const confirmSkip = () => {
  showSkipConfirmation.value = false;
  resetForm();
  emit('skip');
};

const confirmAndExport = () => {
  if (!isFormValid.value) {
    validateForm();
    return;
  }

  const data: CopyrightData = {
    authorName: authorName.value.trim(),
    textContact: textContact.value.trim(),
    protectionCode: protectionCode.value,
  };

  if (imageContactPath.value) {
    data.imageContactPath = imageContactPath.value;
  }

  resetForm();
  emit('confirm', data);
};

// Watch for dialog close to cleanup
watch(showDialog, (newValue) => {
  if (!newValue) {
    // Clean up any blob URLs when dialog is closed
    if (imageContactPreview.value) {
      URL.revokeObjectURL(imageContactPreview.value);
      imageContactPreview.value = null;
    }
  }
});
</script>

<style scoped>
.copyright-dialog-card {
  width: 90vw;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Fixed Header with Glass Background */
.copyright-dialog-header {
  position: sticky;
  top: 0;
  z-index: 10;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

/* Scrollable Content Area */
.copyright-dialog-content {
  flex: 1;
  overflow-y: auto;
  /* Custom scrollbar styling */
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
}

.copyright-dialog-content::-webkit-scrollbar {
  width: 6px;
}

.copyright-dialog-content::-webkit-scrollbar-track {
  background: transparent;
}

.copyright-dialog-content::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.copyright-dialog-content::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

/* Fixed Footer with Glass Background */
.copyright-dialog-footer {
  position: sticky;
  bottom: 0;
  z-index: 10;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

/* Compact spacing for inputs */
.copyright-dialog-content .q-field {
  margin-bottom: 8px;
}

/* Ensure inputs work properly on smaller screens */
@media (max-width: 600px) {
  .copyright-dialog-card {
    width: 95vw;
    max-width: none;
    max-height: 85vh;
  }
}
</style>