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
      <!-- Single scrollable content area with sticky elements -->
      <q-card-section class="copyright-dialog-content">
        <!-- Sticky Header with Glass Background -->
        <div class="copyright-dialog-header">
          <div class="text-h6 text-weight-medium q-pa-sm">{{ $t('copyrightDialog.title') }}</div>
          <div class="text-caption text-grey-6 q-px-sm q-pb-sm">{{ $t('copyrightDialog.subtitle') }}</div>
        </div>
        
        <!-- Content -->
        <div class="content-wrapper">
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
            
            <!-- Enhanced Image Preview -->
            <div v-if="imageContactPreview" class="image-preview-container q-mt-sm">
              <q-card class="image-preview-card">
                <q-card-section class="q-pa-sm">
                  <div class="image-preview-wrapper">
                    <q-img
                      :src="imageContactPreview"
                      class="rounded-borders image-preview"
                      fit="contain"
                      loading="lazy"
                    >
                      <template v-slot:loading>
                        <q-inner-loading showing>
                          <q-spinner-gears size="50px" color="primary" />
                        </q-inner-loading>
                      </template>
                      <template v-slot:error>
                        <div class="absolute-full flex flex-center bg-negative text-white">
                          <q-icon name="broken_image" size="24px" />
                        </div>
                      </template>
                    </q-img>
                    
                    <!-- Status Overlay -->
                    <div class="image-status-overlay">
                      <div v-if="isUploading" class="status-indicator uploading">
                        <q-spinner color="white" size="16px" />
                        <span class="q-ml-xs text-white text-caption">{{ $t('copyrightDialog.uploading') }}</span>
                      </div>
                      <div v-else-if="imageContactPath" class="status-indicator success">
                        <q-icon name="check_circle" color="white" size="16px" />
                        <span class="q-ml-xs text-white text-caption">{{ $t('copyrightDialog.loaded') }}</span>
                      </div>
                      <div v-else class="status-indicator error">
                        <q-icon name="error" color="white" size="16px" />
                        <span class="q-ml-xs text-white text-caption">{{ $t('copyrightDialog.uploadFailed') }}</span>
                      </div>
                    </div>
                  </div>
                </q-card-section>
                
                <!-- Image Actions -->
                <q-card-actions align="between" class="q-pa-sm">
                  <div class="text-caption text-grey-6">
                    {{ imageContactFile?.name }}
                  </div>
                  <div>
                    <q-btn
                      flat
                      dense
                      icon="visibility"
                      color="blue"
                      size="sm"
                      @click="previewImage"
                    >
                      <q-tooltip>{{ $t('copyrightDialog.previewImage') }}</q-tooltip>
                    </q-btn>
                    <q-btn
                      flat
                      dense
                      icon="close"
                      color="negative"
                      size="sm"
                      @click="removeImage"
                    >
                      <q-tooltip>{{ $t('copyrightDialog.removeImage') }}</q-tooltip>
                    </q-btn>
                  </div>
                </q-card-actions>
              </q-card>
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
        </div>

        <!-- Sticky Footer with Glass Background -->
        <div class="copyright-dialog-footer">
          <q-card-actions align="right" class="q-pa-sm">
            <q-btn
              flat
              :label="$t('copyrightDialog.cancel')"
              color="grey"
              size="sm"
              @click="cancel"
            />
            
            <q-btn
              v-if="!hasExistingCopyright"
              flat
              :label="$t('copyrightDialog.skipAndExport')"
              color="orange"
              size="sm"
              @click="skipAndExport"
            />
            
            <q-btn
              :label="$t('copyrightDialog.confirmAndExport')"
              color="primary"
              size="sm"
              :disable="!isFormValid"
              @click="confirmAndExport"
            />
          </q-card-actions>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>

  <!-- Full-size Image Preview Dialog -->
  <q-dialog v-model="showImagePreview" @hide="showImagePreview = false">
    <q-card class="image-preview-fullscreen">
      <q-card-section class="q-pa-none">
        <q-img
          :src="imageContactPreview"
          fit="contain"
          class="fullscreen-image"
        >
          <template v-slot:loading>
            <q-inner-loading showing>
              <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
          </template>
        </q-img>
      </q-card-section>
      <q-card-actions align="right" class="bg-black/50 text-white q-pa-sm">
        <q-btn flat icon="close" color="white" @click="showImagePreview = false" />
      </q-card-actions>
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
const showImagePreview = ref(false);

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

const previewImage = () => {
  showImagePreview.value = true;
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
  position: relative;
  overflow: hidden;
}

/* Content area - scrolls normally with sticky elements */
.copyright-dialog-content {
  max-height: 80vh;
  overflow-y: auto;
  padding: 0; /* No padding, let content flow naturally */
  display: flex;
  flex-direction: column;
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

/* Sticky Header with Glass Background - 70% transparency */
.copyright-dialog-header {
  position: sticky;
  top: 0;
  z-index: 10;
  background: rgba(255, 255, 255, 0.3); /* 70% transparency */
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* Content wrapper for proper spacing */
.content-wrapper {
  padding: 16px;
  flex: 1; /* Allow content to expand and push footer to bottom */
}

/* Sticky Footer with Glass Background - 70% transparency */
.copyright-dialog-footer {
  position: sticky;
  bottom: 0;
  z-index: 10;
  background: rgba(255, 255, 255, 0.3); /* 70% transparency */
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  margin-top: auto; /* Push footer to bottom */
}

/* Enhanced Image Preview Styles */
.image-preview-container {
  width: 100%;
}

.image-preview-card {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  overflow: hidden;
}

.image-preview-wrapper {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
}

.image-preview {
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.image-status-overlay {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
}

.status-indicator {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
}

.status-indicator.uploading {
  background: rgba(0, 0, 0, 0.7);
}

.status-indicator.success {
  background: rgba(76, 175, 80, 0.9);
}

.status-indicator.error {
  background: rgba(244, 67, 54, 0.9);
}

/* Full-screen Image Preview */
.image-preview-fullscreen {
  width: 90vw;
  height: 90vh;
  max-width: none;
  max-height: none;
}

.fullscreen-image {
  width: 100%;
  height: calc(90vh - 60px); /* Account for action bar */
}

/* Compact spacing for inputs */
.content-wrapper .q-field {
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