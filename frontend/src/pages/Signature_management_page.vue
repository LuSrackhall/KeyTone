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
          <q-btn flat dense round icon="input" @click="showImportDialog" size="md" :title="$t('signature.page.import')">
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

      <!-- 签名列表 -->
      <div v-else class="q-pa-md">
        <div class="space-y-2">
          <q-card
            v-for="signature in signatureList"
            :key="signature.id"
            class="cursor-pointer hover:shadow-lg transition-all relative"
            :style="{ minHeight: signature.cardImage ? '60px' : '50px' }"
          >
            <q-card-section class="q-pa-none" style="display: flex; align-items: center; position: relative">
              <!-- 左侧图片区域 -->
              <div
                class="flex-shrink-0 flex items-center justify-center"
                style="width: 60px; height: 60px; min-width: 60px; background-color: #f5f5f5; border-radius: 4px"
              >
                <q-img
                  v-if="signature.cardImage"
                  :src="getImageUrl(signature.cardImage)"
                  style="width: 50px; height: 50px; cursor: pointer; border-radius: 4px"
                  class="object-cover"
                  @click.stop="handleImagePreview(signature.cardImage)"
                >
                  <template v-slot:loading>
                    <q-spinner color="primary" size="24px" />
                  </template>
                  <template v-slot:error>
                    <q-icon name="image_not_supported" size="24px" color="grey-5" />
                  </template>
                </q-img>
                <!-- 无图片占位符 -->
                <div
                  v-else
                  class="flex flex-col items-center justify-center"
                  :title="$t('signature.page.noImage')"
                  style="width: 50px; height: 50px; cursor: default"
                >
                  <q-icon name="image_not_supported" size="28px" color="grey-4" />
                  <div class="text-caption text-grey-4" style="font-size: 0.65rem; margin-top: 2px">
                    {{ $t('signature.page.noImageHint') }}
                  </div>
                </div>
              </div>

              <!-- 中间信息区域（点击展开菜单） -->
              <div
                :ref="(el) => { if (el) contextMenuRefs.set(signature.id, el as HTMLElement); }"
                class="flex-1 flex flex-col justify-center cursor-pointer hover:bg-grey-2 rounded transition-colors"
                :style="{ padding: signature.cardImage ? '8px 12px' : '8px 12px 8px 0', minWidth: 0 }"
                @click="handleInfoClick(signature, $event)"
                @contextmenu="handleInfoContextMenu(signature, $event)"
              >
                <div
                  class="text-subtitle2 text-weight-bold"
                  style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 0.95rem"
                >
                  {{ signature.name }}
                </div>
                <div
                  v-if="signature.intro"
                  class="text-caption text-grey-7"
                  :style="{
                    marginTop: '2px',
                    display: '-webkit-box',
                    WebkitBoxOrient: 'vertical',
                    WebkitLineClamp: '1',
                    lineClamp: '1',
                    overflow: 'hidden',
                    lineHeight: '1.3',
                    fontSize: '0.8rem',
                  }"
                >
                  {{ signature.intro }}
                </div>
              </div>
            </q-card-section>
          </q-card>
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
      <q-card class="relative" style="background: transparent">
        <q-btn
          icon="close"
          flat
          round
          dense
          color="negative"
          size="md"
          v-close-popup
          class="absolute top-0 right-0 z-10"
          style="background-color: rgba(255, 255, 255, 0.5)"
        />
        <q-card-section class="q-pa-md">
          <img :src="previewImageUrl" class="max-w-full" style="max-height: 80vh; border-radius: 4px" />
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- 虚拟菜单参考元素（用于精确定位菜单到点击位置） -->
    <div ref="virtualMenuRef" style="position: fixed; width: 1px; height: 1px; pointer-events: none; z-index: -1" />

    <!-- 上下文菜单 -->
    <q-menu
      v-model="contextMenuVisible"
      :target="(virtualMenuRef as any)"
      :anchor="menuAnchor"
      :self="menuSelf"
      no-parent-event
      @hide="
        () => {
          activeMenuSignatureId = null;
        }
      "
    >
      <q-list dense style="min-width: 120px">
        <q-item clickable v-close-popup @click="handleEdit(contextMenuSignature!)">
          <q-item-section avatar>
            <q-icon name="edit" />
          </q-item-section>
          <q-item-section>{{ $t('signature.page.edit') }}</q-item-section>
        </q-item>
        <q-item clickable v-close-popup @click="handleExport(contextMenuSignature!)">
          <q-item-section avatar>
            <q-icon name="drive_file_move" />
          </q-item-section>
          <q-item-section>{{ $t('signature.page.export') }}</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable v-close-popup @click="handleDelete(contextMenuSignature!)" class="text-negative">
          <q-item-section avatar>
            <q-icon name="delete" />
          </q-item-section>
          <q-item-section>{{ $t('signature.page.delete') }}</q-item-section>
        </q-item>
      </q-list>
    </q-menu>
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
const contextMenuVisible = ref(false);
const contextMenuTarget = ref<Element | undefined>();
const contextMenuSignature = ref<Signature | null>(null);
const contextMenuRefs = new Map<string, HTMLElement>();
const activeMenuSignatureId = ref<string | null>(null);
const virtualMenuRef = ref<HTMLElement | null>(null);
const menuAnchor = ref<
  | 'bottom left'
  | 'bottom right'
  | 'top left'
  | 'top right'
  | 'center left'
  | 'center right'
  | 'bottom middle'
  | 'top middle'
  | 'center middle'
>('bottom left');
const menuSelf = ref<
  | 'bottom left'
  | 'bottom right'
  | 'top left'
  | 'top right'
  | 'center left'
  | 'center right'
  | 'bottom middle'
  | 'top middle'
  | 'center middle'
>('top left');

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

function showContextMenu(signature: Signature, event: MouseEvent) {
  contextMenuSignature.value = signature;
  contextMenuTarget.value = event.currentTarget as Element;
  contextMenuVisible.value = true;
}

function handleInfoClick(signature: Signature, event: MouseEvent) {
  // 左键点击展开菜单或收起菜单
  event.preventDefault();
  event.stopPropagation();

  if (activeMenuSignatureId.value === signature.id && contextMenuVisible.value) {
    // 重复点击同一项，收起菜单
    contextMenuVisible.value = false;
    activeMenuSignatureId.value = null;
  } else {
    // 展开菜单 - 使用点击位置计算最佳菜单位置
    calculateMenuPosition(null, event.clientX, event.clientY);
    contextMenuSignature.value = signature;
    contextMenuVisible.value = true;
    activeMenuSignatureId.value = signature.id;
  }
}

function handleInfoContextMenu(signature: Signature, event: MouseEvent) {
  // 右键点击展开菜单或收起菜单
  event.preventDefault();
  event.stopPropagation();

  // 检查是否已经打开了相同签名的菜单
  if (activeMenuSignatureId.value === signature.id && contextMenuVisible.value) {
    // 重复右键点击同一项，收起菜单
    contextMenuVisible.value = false;
    activeMenuSignatureId.value = null;
  } else {
    // 展开菜单 - 使用点击位置计算最佳菜单位置
    calculateMenuPosition(null, event.clientX, event.clientY);
    contextMenuSignature.value = signature;
    contextMenuVisible.value = true;
    activeMenuSignatureId.value = signature.id;
  }
}

function calculateMenuPosition(element: HTMLElement | null, clientX: number, clientY: number) {
  // 将虚拟参考元素定位到点击位置，使菜单精确显示在点击点附近
  if (virtualMenuRef.value) {
    virtualMenuRef.value.style.left = `${clientX}px`;
    virtualMenuRef.value.style.top = `${clientY}px`;
  }

  // 菜单的预计大小（根据项数）
  const menuHeight = 130; // 约3项的高度
  const menuWidth = 120;

  // 计算从点击位置到视口边界的距离
  const spaceBelow = window.innerHeight - clientY;
  const spaceAbove = clientY;
  const spaceRight = window.innerWidth - clientX;
  const spaceLeft = clientX;

  // 决定垂直位置 - 优先使用下方，如果空间不足则使用上方
  let verticalAnchor: 'top' | 'bottom' = 'bottom';
  let verticalSelf: 'top' | 'bottom' = 'top';

  if (spaceBelow >= menuHeight) {
    // 下方有足够空间，菜单在点击点下方展开
    verticalAnchor = 'bottom';
    verticalSelf = 'top';
  } else if (spaceAbove >= menuHeight) {
    // 下方空间不足但上方充足，菜单在点击点上方展开
    verticalAnchor = 'top';
    verticalSelf = 'bottom';
  } else {
    // 两方都不足（很少见），默认向下，可能被视口裁剪
    verticalAnchor = 'bottom';
    verticalSelf = 'top';
  }

  // 决定水平位置 - 优先使用右方，如果空间不足则使用左方
  let horizontalAnchor: 'left' | 'right' = 'left';
  let horizontalSelf: 'left' | 'right' = 'left';

  if (spaceRight >= menuWidth) {
    // 右方有足够空间，菜单在点击点右方展开
    horizontalAnchor = 'left';
    horizontalSelf = 'left';
  } else if (spaceLeft >= menuWidth) {
    // 右方空间不足但左方充足，菜单在点击点左方展开
    horizontalAnchor = 'right';
    horizontalSelf = 'right';
  } else {
    // 两方都不足（空间特别拥挤），默认向右
    horizontalAnchor = 'left';
    horizontalSelf = 'left';
  }

  // 组合垂直和水平位置
  menuAnchor.value = `${verticalAnchor} ${horizontalAnchor}` as any;
  menuSelf.value = `${verticalSelf} ${horizontalSelf}` as any;
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
