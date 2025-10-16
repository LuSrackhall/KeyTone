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
  <q-dialog v-model="dialogVisible" backdrop-filter="invert(70%)" @hide="handleClose">
    <q-card class="w-96">
      <!-- 标题 -->
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ isEditMode ? $t('signature.form.editTitle') : $t('signature.form.createTitle') }}
      </q-card-section>

      <!-- 编辑提示信息 -->
      <q-card-section v-if="isEditMode" class="q-py-sm bg-blue-1">
        <div class="text-caption text-blue-9">
          {{ $t('signature.form.editWarning') }}
        </div>
      </q-card-section>

      <!-- 可滚动内容区 -->
      <div class="form-content-scroll">
        <!-- 签名名称 -->
        <div class="q-px-md q-pt-md q-pb-sm">
          <div>
            <q-input
              v-model="formData.name"
              :label="$t('signature.form.name')"
              outlined
              dense
              :readonly="isEditMode"
              :disable="isEditMode"
              :rules="[(val) => (val && val.length > 0) || $t('signature.form.nameRequired')]"
              class="name-input-wrapper"
            />
          </div>

          <!-- 个人介绍 -->
          <div class="q-mt-md q-mb-md">
            <q-input
              v-model="formData.intro"
              :label="$t('signature.form.intro')"
              outlined
              dense
              type="textarea"
              autogrow
              class="intro-input-wrapper"
              @input="handleIntroInput"
            />
          </div>

          <!-- 名片图片 -->
          <div class="q-mt-md">
            <div class="text-caption text-grey-7 q-mb-sm">{{ $t('signature.form.cardImage') }}</div>

            <!-- 文件选择器 -->
            <q-file
              v-model="formData.cardImage"
              :label="$t('signature.form.selectImage')"
              outlined
              dense
              accept="image/*"
              @update:model-value="handleImageChange"
            >
              <template v-slot:prepend>
                <q-icon name="image" />
              </template>
            </q-file>
            <div class="text-caption text-grey-6 q-mt-xs">
              {{ $t('signature.form.imageHint') }}
            </div>

            <!-- 图片快速预览（在选择器下方） -->
            <div v-if="imagePreview" class="q-mt-md flex justify-center">
              <div class="relative">
                <img
                  :src="imagePreview"
                  class="w-32 h-32 object-cover rounded-lg border-2 border-grey-4 cursor-pointer hover:opacity-90"
                  @click="showImagePreviewDialog = true"
                  :title="$t('signature.form.clickToPreview')"
                />
                <q-btn
                  round
                  dense
                  flat
                  icon="close"
                  size="sm"
                  class="absolute -top-2 -right-2 bg-red-500 text-white hover:bg-red-600"
                  @click="removeImage"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 固定底部按钮 -->
      <q-card-actions align="right">
        <q-btn flat :label="$t('signature.form.cancel')" color="primary" @click="handleClose" />
        <q-btn
          flat
          :label="isEditMode ? $t('signature.form.update') : $t('signature.form.create')"
          color="primary"
          @click="handleSubmit"
          :loading="loading"
        />
      </q-card-actions>
    </q-card>

    <!-- 图片大图预览对话框 -->
    <q-dialog v-model="showImagePreviewDialog" backdrop-filter="invert(70%)">
      <q-card class="image-preview-card relative" style="background: transparent; max-width: 90vw; max-height: 90vh">
        <q-btn
          icon="close"
          flat
          round
          dense
          color="negative"
          size="md"
          v-close-popup
          class="absolute top-0 right-0 z-10"
          style="background-color: rgba(255, 255, 255, 0.3)"
        />
        <q-card-section class="q-pa-none flex items-center justify-center" style="min-width: 300px; min-height: 300px">
          <img :src="imagePreview" class="max-w-full max-h-full" style="object-fit: contain; display: block" />
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { createSignature, updateSignature, fileToBase64 } from 'boot/query/signature-query';
import type { Signature } from 'src/types/signature';

const props = defineProps<{
  modelValue: boolean;
  signature?: Signature | null;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'success'): void;
}>();

const q = useQuasar();
const { t: $t } = useI18n();

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

const isEditMode = computed(() => !!props.signature);
const loading = ref(false);
const showImagePreviewDialog = ref(false);

const formData = ref<{
  name: string;
  intro: string;
  cardImage: File | null;
}>({
  name: '',
  intro: '',
  cardImage: null,
});

const imagePreview = ref<string>('');

// 监听 signature 属性变化，填充表单
watch(
  () => props.signature,
  (newVal) => {
    if (newVal) {
      formData.value.name = newVal.name;
      formData.value.intro = newVal.intro || '';
      formData.value.cardImage = null;

      // 如果有现有图片，显示预览
      if (newVal.cardImage) {
        const port = (window as any).myWindowAPI?.getBackendPort() || 38888;
        imagePreview.value = `http://127.0.0.1:${port}/signature/image/${newVal.cardImage}`;
      } else {
        imagePreview.value = '';
      }
    } else {
      resetForm();
    }
  },
  { immediate: true }
);

function resetForm() {
  formData.value = {
    name: '',
    intro: '',
    cardImage: null,
  };
  imagePreview.value = '';
}

function handleClose() {
  resetForm();
  dialogVisible.value = false;
}

async function handleImageChange(file: File | null) {
  if (file) {
    try {
      const base64 = await fileToBase64(file);
      imagePreview.value = base64;
    } catch (error) {
      console.error('Failed to read image file:', error);
      q.notify({
        type: 'negative',
        message: $t('signature.notify.imageReadFailed'),
        position: 'top',
      });
    }
  } else {
    imagePreview.value = '';
  }
}

function removeImage() {
  formData.value.cardImage = null;
  imagePreview.value = '';
}

function handleIntroInput() {
  // 使用 autogrow 模式，无需手动调整行数
  // 此函数保留为兼容性，但功能由 autogrow 属性处理
}

async function handleSubmit() {
  // 验证表单
  if (!formData.value.name || formData.value.name.length === 0) {
    q.notify({
      type: 'warning',
      message: $t('signature.form.nameRequired'),
      position: 'top',
    });
    return;
  }

  loading.value = true;

  try {
    let cardImageBase64 = '';

    // 如果有新图片，转换为 Base64
    if (formData.value.cardImage) {
      cardImageBase64 = await fileToBase64(formData.value.cardImage);
    }

    if (isEditMode.value && props.signature) {
      // 更新模式
      const result = await updateSignature({
        id: props.signature.id,
        name: formData.value.name,
        intro: formData.value.intro,
        cardImage: cardImageBase64 || undefined,
      });

      if (result) {
        q.notify({
          type: 'positive',
          message: $t('signature.notify.updateSuccess'),
          position: 'top',
        });
        emit('success');
        handleClose();
      } else {
        q.notify({
          type: 'negative',
          message: $t('signature.notify.updateFailed'),
          position: 'top',
        });
      }
    } else {
      // 创建模式
      const result = await createSignature({
        name: formData.value.name,
        intro: formData.value.intro,
        cardImage: cardImageBase64 || undefined,
      });

      if (result) {
        q.notify({
          type: 'positive',
          message: $t('signature.notify.createSuccess'),
          position: 'top',
        });
        emit('success');
        handleClose();
      } else {
        q.notify({
          type: 'negative',
          message: $t('signature.notify.createFailed'),
          position: 'top',
        });
      }
    }
  } catch (error) {
    console.error('Failed to submit signature:', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.unexpectedError'),
      position: 'top',
    });
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* 可滚动的表单内容区 */
.form-content-scroll {
  max-height: calc(90vh - 180px); /* 给标题、提示、按钮留出空间 */
  overflow-y: auto;
  overflow-x: hidden;
}

/* 滚动条样式 */
.form-content-scroll::-webkit-scrollbar {
  width: 4px;
}

.form-content-scroll::-webkit-scrollbar-track {
  background: rgba(228, 228, 228, 0.2);
  border-radius: 2px;
}

.form-content-scroll::-webkit-scrollbar-thumb {
  background: rgba(51, 65, 85, 0.3);
  border-radius: 2px;
}

.form-content-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(51, 65, 85, 0.5);
}

.name-input-wrapper :deep(.q-field__control) {
  /* 确保输入框容器能够溢出 */
  overflow: hidden;
}

.name-input-wrapper :deep(input) {
  /* 单行模式 */
  white-space: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  /* 参考主页选择框的滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: rgba(203, 213, 225, 0.4) transparent;
}

/* 自定义滚动条样式（webkit浏览器） - 参考主页风格 */
.name-input-wrapper :deep(input::-webkit-scrollbar) {
  height: 4px;
}

.name-input-wrapper :deep(input::-webkit-scrollbar-track) {
  background: rgba(228, 228, 228, 0.2);
  border-radius: 2px;
}

.name-input-wrapper :deep(input::-webkit-scrollbar-thumb) {
  background: rgba(51, 65, 85, 0.3);
  border-radius: 2px;
}

.name-input-wrapper :deep(input::-webkit-scrollbar-thumb:hover) {
  background: rgba(51, 65, 85, 0.5);
}

/* 介绍文本框自动高度 */
.intro-input-wrapper :deep(.q-field__control) {
  display: flex;
  align-items: flex-start;
}

.intro-input-wrapper :deep(textarea) {
  /* autogrow 会自动调整高度，我们设置最大高度限制 */
  max-height: calc(1.5em * 3 + 8px); /* 最多3行的高度 */
  resize: none;
  overflow-y: auto;
  /* 自定义滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: rgba(203, 213, 225, 0.4) transparent;
  /* 确保文本可换行显示 */
  word-wrap: break-word;
  white-space: pre-wrap;
}

/* 介绍文本框的滚动条样式 */
.intro-input-wrapper :deep(textarea::-webkit-scrollbar) {
  width: 4px;
}

.intro-input-wrapper :deep(textarea::-webkit-scrollbar-track) {
  background: rgba(228, 228, 228, 0.2);
  border-radius: 2px;
}

.intro-input-wrapper :deep(textarea::-webkit-scrollbar-thumb) {
  background: rgba(51, 65, 85, 0.3);
  border-radius: 2px;
}

.intro-input-wrapper :deep(textarea::-webkit-scrollbar-thumb:hover) {
  background: rgba(51, 65, 85, 0.5);
}

/* 大图预览卡片样式 */
.image-preview-card {
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 90vw;
  max-height: 90vh;
}

.image-preview-card img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  display: block;
  border-radius: 8px;
}
</style>

```
