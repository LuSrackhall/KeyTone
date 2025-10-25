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
import { nanoid } from 'nanoid';
import { createSignature, updateSignature, decryptSignatureData, getSignatureImage } from 'boot/query/signature-query';
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

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

// 是否为编辑模式（由 props.signature 是否存在决定）
const isEditMode = computed(() => !!props.signature);

// 表单提交加载状态 - 绑定提交按钮的 loading 属性
const loading = ref(false);

// 图片预览大图对话框显示状态
const showImagePreviewDialog = ref(false);

// 原始签名的深拷贝 - 用于编辑模式，确保编辑过程不影响原列表项
const editingSignatureClone = ref<Signature | null>(null);

// 表单数据对象 - 绑定到各输入框
// 结构: { name, intro, cardImage }
// - name: 签名名称，编辑模式下禁用
// - intro: 个人介绍，支持多行，允许清空
// - cardImage: 选择的图片文件（File | null）
const formData = ref<{
  name: string;
  intro: string;
  cardImage: File | null;
}>({
  name: '',
  intro: '',
  cardImage: null,
});

// 原始图片预览 URL（编辑模式下，保存从后端加载的原图片预览）
const originalImageUrl = ref<string>('');

// 当前图片预览 URL - 绑定到预览图片和大图预览对话框
// 可以是 Base64 字符串（新选择的图片）或 HTTP URL（原图片）
const imagePreview = ref<string>('');

// 标记图片是否有变化（用于判断是否需要上传新图片）
const imageChanged = ref(false);

// 编辑模式下，保存原始的表单数据（用于检测是否有变更）
const originalFormData = ref<{
  name: string;
  intro: string;
}>({
  name: '',
  intro: '',
});

/**
 * 检测编辑模式下是否有变更
 * @returns true 表示有变更，false 表示无变更
 */
function hasChanges(): boolean {
  if (!isEditMode.value || !editingSignatureClone.value) {
    return true; // 创建模式总是返回 true
  }

  // 检查名称是否改变（编辑模式下名称禁用，所以名称不应改变）
  if (formData.value.name !== originalFormData.value.name) {
    return true;
  }

  // 检查介绍是否改变
  if (formData.value.intro !== originalFormData.value.intro) {
    return true;
  }

  // 检查图片是否改变
  if (imageChanged.value) {
    return true;
  }

  return false;
}

/**
 * 创建签名对象的深拷贝
 * @param signature 原始签名对象
 * @returns 深拷贝后的签名对象
 */
function deepCloneSignature(signature: Signature): Signature {
  return {
    id: signature.id,
    name: signature.name,
    intro: signature.intro,
    cardImage: signature.cardImage, // cardImage 是文件路径字符串或 File 对象，无需深拷贝
  };
}

/**
 * 将 Blob 转换为 File 对象，使用国际化的固定文件名（无后缀）
 * @param blob Blob 对象
 * @param mimeType MIME 类型，用于后端识别文件类型
 * @returns File 对象
 */
function convertBlobToFile(blob: Blob, mimeType: string): File {
  // 获取国际化的文件名（不包含后缀）
  const fileName = $t('signature.form.cardImageFileName');

  // 将 Blob 转换为 File 对象
  // 注：文件名不包含后缀，后端通过 MIME 类型识别文件格式
  return new File([blob], fileName, { type: mimeType });
}

/**
 * 加载原始签名的图片预览（编辑模式）
 * @param imagePath 图片文件路径
 */
async function loadOriginalImagePreview(imagePath: string) {
  if (!imagePath) {
    originalImageUrl.value = '';
    return;
  }

  try {
    const blob = await getSignatureImage(imagePath);
    if (blob) {
      originalImageUrl.value = URL.createObjectURL(blob);
      imagePreview.value = originalImageUrl.value;

      // 将原始图片 Blob 转换为 File 对象并赋值给表单
      // 这样在提交时，即使用户没有修改图片，也能正确上传
      const mimeType = blob.type || 'image/jpeg'; // 默认使用 image/jpeg
      formData.value.cardImage = convertBlobToFile(blob, mimeType);
    }
  } catch (error) {
    console.warn('Failed to load original image preview:', error);
  }
}

// 监听 props.signature 变化，填充表单数据
watch(
  () => props.signature,
  async (newVal) => {
    if (newVal) {
      // 编辑模式：创建深拷贝，防止编辑过程影响原列表项
      editingSignatureClone.value = deepCloneSignature(newVal);

      // 填充表单数据
      formData.value = {
        name: newVal.name,
        intro: newVal.intro,
        cardImage: null, // 编辑模式下，cardImage 初始为 null（表示不上传新图片）
      };

      // 保存原始表单数据（用于变更检测）
      originalFormData.value = {
        name: newVal.name,
        intro: newVal.intro,
      };

      // 重置图片变化标志
      imageChanged.value = false;

      // 如果有原始图片，加载预览
      if (newVal.cardImage && typeof newVal.cardImage === 'string') {
        await loadOriginalImagePreview(newVal.cardImage);
      } else {
        originalImageUrl.value = '';
        imagePreview.value = '';
      }
    } else {
      // 创建模式：重置所有状态
      editingSignatureClone.value = null;
      resetForm();
    }
  },
  { immediate: true }
);

/** 重置表单到初始状态 */
function resetForm() {
  formData.value = {
    name: '',
    intro: '',
    cardImage: null,
  };
  originalFormData.value = {
    name: '',
    intro: '',
  };
  imagePreview.value = '';
  originalImageUrl.value = '';
  imageChanged.value = false;
  editingSignatureClone.value = null;
}

/** 关闭对话框并重置表单 */
function handleClose() {
  // 清理 Blob URL（如果是加载的原始图片）
  if (originalImageUrl.value) {
    URL.revokeObjectURL(originalImageUrl.value);
  }
  resetForm();
  dialogVisible.value = false;
}

/** 处理图片文件变化 - 生成预览 */
async function handleImageChange(file: File | null) {
  if (file) {
    try {
      // 图片转 Base64 预览生成
      const reader = new FileReader();
      reader.onload = (e) => {
        imagePreview.value = e.target?.result as string;
        imageChanged.value = true; // 标记图片有变化
      };
      reader.readAsDataURL(file);
    } catch (error) {
      console.error('Failed to read image file:', error);
      q.notify({
        type: 'negative',
        message: $t('signature.notify.imageReadFailed'),
        position: 'top',
      });
    }
  } else {
    // 文件被清空
    imagePreview.value = '';
    imageChanged.value = true; // 标记图片有变化（可能是删除操作）
  }
}

/** 移除已选择的图片 */
function removeImage() {
  formData.value.cardImage = null;
  imagePreview.value = '';
  imageChanged.value = true; // 标记图片有变化（用户主动删除）
}

/** 处理介绍文本输入 */
function handleIntroInput() {
  // 由 autogrow 属性处理自动高度，此处保留扩展空间
  // 允许用户清空介绍文本
}

/** 提交表单 - 创建或更新签名 */
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

  // 编辑模式下，检测是否有变更
  if (isEditMode.value && !hasChanges()) {
    q.notify({
      type: 'info',
      message: $t('signature.notify.noChanges') || 'No changes detected',
      position: 'top',
    });
    return;
  }

  loading.value = true;

  try {
    if (isEditMode.value) {
      // 编辑模式：更新签名
      if (!editingSignatureClone.value) {
        throw new Error('Editing signature clone not found');
      }

      // 构建更新数据
      // 在编辑模式下，formData.value.cardImage 已经被设置为：
      // 1. 如果用户选择了新图片，则为新图片的 File 对象
      // 2. 如果没有选择新图片，则为原始图片转换后的 File 对象（已在 loadOriginalImagePreview 中设置）
      // 3. 如果用户主动删除图片，则为 null 且 imageChanged 为 true
      let cardImage: File | null = formData.value.cardImage;

      // 如果没有图片且用户主动删除，则设置为空 File
      if (!cardImage && imageChanged.value) {
        cardImage = new File([], '');
      }

      const updateData: Signature = {
        id: editingSignatureClone.value.id, // 使用加密的 ID
        name: formData.value.name,
        intro: formData.value.intro,
        cardImage: cardImage || new File([], ''),
        imageChanged: imageChanged.value, // 传递图片是否发生变更的标记
      };

      const success = await updateSignature(updateData);
      if (success) {
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
      // 创建模式：新建签名
      const signatureData: Signature = {
        id: nanoid(21),
        name: formData.value.name,
        intro: formData.value.intro,
        cardImage: formData.value.cardImage || new File([], ''),
      };

      const success = await createSignature(signatureData);
      if (success) {
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

<style lang="scss" scoped>
/* 可滚动的表单内容区 */
.form-content-scroll {
  max-height: calc(90vh - 180px);
  overflow-y: auto;
  overflow-x: hidden;
}

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

.name-input-wrapper {
  :deep(.q-field__native) {
    @apply h-auto;
  }
}

.name-input-wrapper :deep(.q-field__control) {
  overflow: hidden;
}

.name-input-wrapper :deep(input) {
  white-space: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: thin;
  scrollbar-color: rgba(203, 213, 225, 0.4) transparent;
}

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

.intro-input-wrapper :deep(.q-field__control) {
  display: flex;
  align-items: flex-start;
}

.intro-input-wrapper :deep(textarea) {
  max-height: calc(1.5em * 3 + 8px);
  resize: none;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(203, 213, 225, 0.4) transparent;
  word-wrap: break-word;
  white-space: pre-wrap;
}

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
