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
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ isEditMode ? $t('signature.form.editTitle') : $t('signature.form.createTitle') }}
      </q-card-section>

      <q-card-section>
        <!-- 签名名称 -->
        <q-input
          v-model="formData.name"
          :label="$t('signature.form.name')"
          :hint="$t('signature.form.nameHint')"
          outlined
          dense
          :rules="[
            (val) => (val && val.length > 0) || $t('signature.form.nameRequired'),
            (val) => (val && val.length <= 50) || $t('signature.form.nameTooLong'),
          ]"
          class="mb-4"
        />

        <!-- 个人介绍 -->
        <q-input
          v-model="formData.intro"
          :label="$t('signature.form.intro')"
          :hint="$t('signature.form.introHint')"
          outlined
          dense
          type="textarea"
          rows="3"
          :rules="[(val) => !val || val.length <= 500 || $t('signature.form.introTooLong')]"
          class="mb-4"
        />

        <!-- 名片图片 -->
        <div class="mb-4">
          <div class="text-caption text-grey-7 mb-2">{{ $t('signature.form.cardImage') }}</div>

          <!-- 图片预览 -->
          <div v-if="imagePreview" class="mb-2 flex justify-center">
            <div class="relative">
              <img :src="imagePreview" class="w-32 h-32 object-cover rounded-lg border-2 border-grey-4" />
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

          <!-- 文件选择器 -->
          <q-file
            v-model="formData.cardImage"
            :label="$t('signature.form.selectImage')"
            outlined
            dense
            accept="image/png,image/jpeg,image/jpg,image/webp"
            @update:model-value="handleImageChange"
          >
            <template v-slot:prepend>
              <q-icon name="image" />
            </template>
          </q-file>
          <div class="text-caption text-grey-6 mt-1">
            {{ $t('signature.form.imageHint') }}
          </div>
        </div>
      </q-card-section>

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

  if (formData.value.name.length > 50) {
    q.notify({
      type: 'warning',
      message: $t('signature.form.nameTooLong'),
      position: 'top',
    });
    return;
  }

  if (formData.value.intro && formData.value.intro.length > 500) {
    q.notify({
      type: 'warning',
      message: $t('signature.form.introTooLong'),
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
/* 自定义样式 */
</style>
