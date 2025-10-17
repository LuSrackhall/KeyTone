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

// 表单数据对象 - 绑定到各输入框
// 结构: { name, intro, cardImage }
// - name: 签名名称，编辑模式下禁用
// - intro: 个人介绍，支持多行
// - cardImage: 选择的图片文件
const formData = ref<{
  name: string;
  intro: string;
  cardImage: File | null;
}>({
  name: '',
  intro: '',
  cardImage: null,
});

// 图片预览 URL - 绑定到预览图片和大图预览对话框
// 可以是 Base64 字符串或 HTTP URL
const imagePreview = ref<string>('');

// 监听 props.signature 变化，填充表单数据
watch(
  () => props.signature,
  (newVal) => {
    if (newVal) {
      // TODO: 具体表单填充逻辑由业务层实现
      // 根据 props.signature 填充 formData 和 imagePreview
    } else {
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
  imagePreview.value = '';
}

/** 关闭对话框并重置表单 */
function handleClose() {
  resetForm();
  dialogVisible.value = false;
}

/** 处理图片文件变化 - 生成预览 */
async function handleImageChange(file: File | null) {
  if (file) {
    try {
      // TODO: 具体图片转 Base64 预览生成逻辑由业务层实现
      imagePreview.value = '';
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

/** 移除已选择的图片 */
function removeImage() {
  formData.value.cardImage = null;
  imagePreview.value = '';
}

/** 处理介绍文本输入 */
function handleIntroInput() {
  // 由 autogrow 属性处理自动高度，此处保留扩展空间
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

  loading.value = true;

  try {
    // TODO: 具体创建/更新逻辑由业务层实现
    // 根据 isEditMode 区分创建模式（isEditMode 为 false）和更新模式（isEditMode 为 true）
    // 需要处理图片上传和表单数据提交
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
