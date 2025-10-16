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
  <q-page
    style="min-height: 0px"
    :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]', 'overflow-hidden']"
  >
    <!-- 页面顶部标题栏 -->
    <div class="q-pa-md">
      <div class="flex justify-between items-center">
        <h1 class="text-h5 q-my-none">{{ $t('signature.page.title') }}</h1>
        <div class="q-gutter-md">
          <q-btn
            flat
            dense
            round
            icon="file_upload"
            @click="showImportDialog"
            size="md"
            :title="$t('signature.page.import')"
          >
            <q-tooltip>{{ $t('signature.page.import') }}</q-tooltip>
          </q-btn>
          <q-btn
            flat
            dense
            round
            icon="add"
            color="primary"
            @click="handleCreate"
            size="md"
            :title="$t('signature.page.create')"
          >
            <q-tooltip>{{ $t('signature.page.create') }}</q-tooltip>
          </q-btn>
        </div>
      </div>
    </div>

    <!-- 滚动容器 -->
    <q-scroll-area
      class="overflow-hidden"
      :style="{ height: isMacOS ? 'calc(458.5px - 80px)' : 'calc(458.5px - 80px)' }"
    >
      <!-- 加载状态 -->
      <div v-if="loading" class="q-pa-lg flex flex-col items-center justify-center" style="min-height: 300px">
        <q-spinner color="primary" size="50px" />
        <p class="q-mt-md">{{ $t('signature.page.title') }}...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="q-pa-lg flex flex-col items-center justify-center" style="min-height: 300px">
        <q-icon name="error_outline" size="64px" color="negative" />
        <p class="q-mt-md">{{ $t('signature.page.loadError') }}</p>
        <q-btn :label="$t('signature.page.retry')" color="primary" @click="loadSignatures" />
      </div>

      <!-- 空状态 -->
      <div
        v-else-if="signatureList.length === 0"
        class="q-pa-lg flex flex-col items-center justify-center"
        style="min-height: 300px"
      >
        <q-icon name="badge" size="64px" color="grey-5" />
        <p class="q-mt-md">{{ $t('signature.page.emptyState') }}</p>
        <q-btn :label="$t('signature.page.createFirst')" color="primary" icon="add" @click="handleCreate" />
      </div>

      <!-- 签名列表网格 -->
      <div v-else class="q-pa-md">
        <div class="row q-col-gutter-md">
          <div v-for="signature in signatureList" :key="signature.id" class="col-xs-12 col-sm-6 col-md-4 col-lg-3">
            <q-card clickable @click="handleEdit(signature)" class="cursor-pointer hover:shadow-lg transition-all">
              <!-- 图片 -->
              <q-img
                v-if="signature.cardImage"
                :src="getImageUrl(signature.cardImage)"
                ratio="1"
                class="cursor-pointer"
                @click.stop="handleImagePreview(signature.cardImage)"
              >
                <div class="absolute-top-left q-pa-xs">
                  <q-btn round flat dense icon="close" color="red" size="sm" @click.stop="handleDelete(signature)" />
                </div>
                <div class="absolute-top-right q-pa-xs">
                  <q-btn
                    round
                    flat
                    dense
                    icon="file_download"
                    color="blue"
                    size="sm"
                    @click.stop="handleExport(signature)"
                  />
                </div>
              </q-img>
              <div v-else class="flex items-center justify-center bg-grey-2" style="aspect-ratio: 1">
                <q-icon name="person" size="80px" color="grey-5" />
                <div class="absolute-top-left q-pa-xs">
                  <q-btn round flat dense icon="close" color="red" size="sm" @click.stop="handleDelete(signature)" />
                </div>
                <div class="absolute-top-right q-pa-xs">
                  <q-btn
                    round
                    flat
                    dense
                    icon="file_download"
                    color="blue"
                    size="sm"
                    @click.stop="handleExport(signature)"
                  />
                </div>
              </div>

              <!-- 签名信息 -->
              <q-card-section class="q-pa-md">
                <div class="text-subtitle2 text-weight-bold ellipsis">{{ signature.name }}</div>
                <div v-if="signature.intro" class="text-caption text-grey-7 ellipsis-2-lines q-mt-xs">
                  {{ signature.intro }}
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </div>
    </q-scroll-area>

    <!-- 签名表单对话框 -->
    <SignatureFormDialog v-model="showFormDialog" :signature="selectedSignature" @success="handleFormSuccess" />

    <!-- 导入对话框 -->
    <q-dialog v-model="showImportDialogVisible" backdrop-filter="blur(4px)">
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ $t('signature.import.title') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-file
            v-model="importFile"
            :label="$t('signature.import.selectFile')"
            outlined
            accept=".ktsign"
            dense
            @update:model-value="() => {}"
          />
          <div class="text-caption text-grey-7 q-mt-md">{{ $t('signature.import.fileHint') }}</div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="$t('signature.form.cancel')" color="primary" v-close-popup />
          <q-btn
            flat
            :label="$t('signature.import.import')"
            color="primary"
            @click="handleImport"
            :loading="importing"
            :disable="!importFile"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 图片预览对话框 -->
    <q-dialog v-model="showImagePreview" backdrop-filter="blur(4px)">
      <q-card>
        <q-card-section class="q-pa-none">
          <img :src="previewImageUrl" class="max-w-full" style="max-height: 80vh" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('signature.form.close')" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { useSignatureStore } from 'src/stores/signature-store';
import { useAppStore } from 'src/stores/app-store';
import SignatureFormDialog from 'src/components/SignatureFormDialog.vue';
import {
  getAllSignatures,
  deleteSignature,
  exportSignature,
  importSignature,
  getSignatureImageUrl,
} from 'boot/query/signature-query';
import type { Signature, SignatureFile } from 'src/types/signature';

const q = useQuasar();
const { t: $t } = useI18n();
const signature_store = useSignatureStore();
const app_store = useAppStore();

// SSE 消息监听函数
const sseMessageListener = (e: Event) => {
  const messageEvent = e as MessageEvent;
  try {
    const data = JSON.parse(messageEvent.data);
    console.log('[Signature_management_page] SSE 消息接收:', data);

    // 处理全量配置数据推送
    if (data.key === 'get_all_value' && data.value) {
      const value = data.value;
      console.log('[Signature_management_page] 检查 signature_manager 字段');

      // 从全量配置中提取签名数据
      if (value.signature_manager) {
        console.log('[Signature_management_page] 检测到 signature_manager 字段，触发重新加载签名数据');
        // SSE 中的 signature_manager 是加密后的键值对（未解密），前端需要调用后端 API 获取解密后的完整数据
        // 这是正常的工作流程：SSE 通知"配置变化了"，前端再通过 API 拉取解密后的签名列表
        loadSignatures();
      } else {
        console.log('[Signature_management_page] 配置中未找到 signature_manager');
      }
    }
  } catch (error) {
    console.error('[Signature_management_page] SSE 消息处理异常:', error);
  }
};

// 获取系统类型
const isMacOS = computed(() => {
  if (process.env.MODE === 'electron') {
    return process.platform === 'darwin' || navigator.platform === 'MacIntel';
  }
  return navigator.platform === 'MacIntel';
});

const loading = ref(false);
const error = ref(false);
const showFormDialog = ref(false);
const selectedSignature = ref<Signature | null>(null);
const showImportDialogVisible = ref(false);
const importFile = ref<File | null>(null);
const importing = ref(false);
const showImagePreview = ref(false);
const previewImageUrl = ref('');

// 从 store 获取签名列表
const signatureList = computed(() => {
  console.log('[signatureList] computed 重新计算，store.signatureManager:', signature_store.signatureManager);
  const list = Object.values(signature_store.signatureManager as any) as Signature[];
  // 过滤掉任何不完整或加密未解密的项（即 name 为 undefined 的项）
  const validList = list.filter((item) => {
    if (!item || typeof item !== 'object') {
      console.warn('[signatureList] 跳过非对象项:', item);
      return false;
    }
    if (!item.name) {
      console.warn('[signatureList] 跳过 name 为空的项:', item);
      return false;
    }
    return true;
  });
  return validList.sort((a, b) => a.name.localeCompare(b.name));
});

onMounted(() => {
  console.log('[Signature_management_page] onMounted 被调用');
  loadSignatures();

  // 添加 SSE 消息监听
  console.log('[Signature_management_page] 添加 SSE 消息监听器');
  app_store.eventSource.addEventListener('message', sseMessageListener);
});

onUnmounted(() => {
  console.log('[Signature_management_page] onUnmounted 被调用，移除 SSE 监听器');
  if (sseMessageListener) {
    app_store.eventSource.removeEventListener('message', sseMessageListener);
  }
});

async function loadSignatures() {
  loading.value = true;
  error.value = false;

  try {
    console.log('[loadSignatures] 开始加载');
    const result = await getAllSignatures();
    console.log('[loadSignatures] API 返回:', result);

    if (result !== false && result) {
      console.log('[loadSignatures] 使用 updateFromSSE 更新 store');
      signature_store.updateFromSSE(result);
      console.log('[loadSignatures] 更新完成，当前列表长度:', signatureList.value.length);
    } else {
      console.error('[loadSignatures] API 返回 false');
      error.value = true;
    }
  } catch (err) {
    console.error('[loadSignatures] 异常:', err);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

function handleCreate() {
  selectedSignature.value = null;
  showFormDialog.value = true;
}

function handleEdit(signature: Signature) {
  selectedSignature.value = signature;
  showFormDialog.value = true;
}

async function handleDelete(signature: Signature) {
  q.dialog({
    title: $t('signature.delete.confirmTitle'),
    message: $t('signature.delete.confirmMessage', { name: signature.name }),
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    const result = await deleteSignature(signature.id);
    if (result) {
      q.notify({
        type: 'positive',
        message: $t('signature.notify.deleteSuccess'),
        position: 'top',
      });
      await loadSignatures();
    } else {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.deleteFailed'),
        position: 'top',
      });
    }
  });
}

async function handleExport(signature: Signature) {
  try {
    const fileData = await exportSignature(signature.id);
    if (!fileData) {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.exportFailed'),
        position: 'top',
      });
      return;
    }

    const blob = new Blob([JSON.stringify(fileData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${signature.name}.ktsign`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    q.notify({
      type: 'positive',
      message: $t('signature.notify.exportSuccess'),
      position: 'top',
    });
  } catch (error) {
    console.error('Failed to export signature:', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.exportFailed'),
      position: 'top',
    });
  }
}

function showImportDialog() {
  importFile.value = null;
  showImportDialogVisible.value = true;
}

async function handleImport() {
  if (!importFile.value) return;

  importing.value = true;

  try {
    const text = await importFile.value.text();
    const fileData: SignatureFile = JSON.parse(text);

    if (!fileData.version || !fileData.signature || !fileData.checksum) {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.invalidFileFormat'),
        position: 'top',
      });
      return;
    }

    const result = await importSignature(fileData);
    if (result) {
      q.notify({
        type: 'positive',
        message: $t('signature.notify.importSuccess'),
        position: 'top',
      });
      showImportDialogVisible.value = false;
      await loadSignatures();
    } else {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.importFailed'),
        position: 'top',
      });
    }
  } catch (error) {
    console.error('Failed to import signature:', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.importFailed'),
      position: 'top',
    });
  } finally {
    importing.value = false;
  }
}

function handleFormSuccess() {
  loadSignatures();
}

function getImageUrl(filename: string): string {
  return getSignatureImageUrl(filename);
}

function handleImagePreview(filename: string) {
  previewImageUrl.value = getImageUrl(filename);
  showImagePreview.value = true;
}
</script>

<style scoped>
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
