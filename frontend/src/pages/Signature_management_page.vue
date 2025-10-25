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
      class="overflow-hidden scroll-h-hide"
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
        <TransitionGroup name="list" tag="div" class="space-y-2">
          <q-card
            v-for="signature in signatureList"
            :key="signature.id"
            class="cursor-move hover:shadow-lg transition-all relative draggable-card w-87.2"
            :style="{
              minHeight: '60px',
              zIndex: isDraggingInProgress && draggedSignature?.id === signature.id ? 1000 : 0,
            }"
            draggable="true"
            @dragstart="handleDragStart($event, signature)"
            @dragend="handleDragEnd"
            @dragover.prevent="handleDragOver"
            @drop="handleDrop($event, signature)"
            @dragenter="handleDragEnter"
            @dragleave="handleDragLeave"
          >
            <q-card-section class="q-pa-none" style="display: flex; align-items: center; position: relative">
              <!-- 左侧图片区域 -->
              <div
                class="flex-shrink-0 flex items-center justify-center"
                style="width: 60px; height: 60px; min-width: 60px; background-color: #f5f5f5; border-radius: 4px"
              >
                <!-- 有图片时显示 -->
                <q-img
                  v-if="signature.cardImage"
                  :src="getImageUrl(signature.cardImage as unknown as string)"
                  style="width: 50px; height: 50px; cursor: pointer; border-radius: 4px"
                  class="object-cover"
                  @click.stop="handleImagePreview(signature.cardImage as unknown as string)"
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
                  class="flex items-center justify-center w-full h-full"
                  :title="$t('signature.page.noImageHint')"
                  style="cursor: default"
                >
                  <div
                    class="text-caption text-grey-6 text-center"
                    style="font-size: 0.65rem; line-height: 1.3; padding: 4px"
                  >
                    {{ $t('signature.page.noImageHint') }}
                  </div>
                </div>
              </div>

              <!-- 中间信息区域（点击展开菜单） -->
              <div
                :ref="(el) => { if (el) contextMenuRefs.set(signature.id, el as HTMLElement); }"
                class="flex-1 flex flex-col justify-center cursor-pointer hover:bg-grey-2 rounded transition-colors"
                :style="{ padding: '8px 12px', minWidth: 0 }"
                @click="handleInfoClick(signature, $event)"
                @contextmenu="handleInfoContextMenu(signature, $event)"
              >
                <div
                  class="text-subtitle2 text-weight-bold"
                  :class="[
                    /* 对溢出的情况, 采取滚动策略 */
                    'max-w-full !overflow-x-auto whitespace-nowrap !text-clip',
                    // 添加细微滚动条
                    'h-5.5 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                  ]"
                  style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 0.95rem"
                >
                  {{ signature.name }}
                </div>
                <div
                  v-if="signature.intro"
                  class="text-caption text-grey-7"
                  :class="[
                    /* 对溢出的情况, 采取滚动策略 */
                    'max-w-full !overflow-x-auto whitespace-nowrap !mt-1.5',
                    // 添加细微滚动条
                    'h-4.4 [&::-webkit-scrollbar]:h-0.2 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                  ]"
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
        </TransitionGroup>
      </div>
    </q-scroll-area>

    <!-- 签名表单对话框 -->
    <SignatureFormDialog v-model="showFormDialog" :signature="selectedSignature" @success="handleFormSuccess" />

    <!-- 导入对话框 -->
    <q-dialog v-model="showImportDialogVisible" backdrop-filter="blur(4px)">
      <q-card style="width: 95%; max-width: 480px">
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

    <!-- 导出签名确认对话框 -->
    <q-dialog v-model="showExportConfirmDialog" backdrop-filter="blur(4px)">
      <q-card style="width: 95%; max-width: 360px; max-height: 85vh">
        <!-- 对话框标题 -->
        <q-card-section class="row items-center q-pb-sm q-px-md q-pt-md">
          <div class="text-h6 q-my-none" style="font-size: 1rem">{{ $t('signature.export.confirmTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup size="sm" />
        </q-card-section>

        <!-- 对话框内容（可滚动） -->
        <q-scroll-area style="height: calc(85vh - 140px); min-height: 200px">
          <q-card-section class="q-px-md q-py-sm">
            <!-- 安全警告标题 -->
            <div class="text-subtitle2 text-negative q-mb-md" style="font-size: 0.9rem">
              {{ $t('signature.export.confirmMessage') }}
            </div>

            <!-- 主要警告 -->
            <div
              class="q-pa-sm"
              style="background-color: rgba(244, 67, 54, 0.1); border-radius: 4px; border-left: 3px solid #f44336"
            >
              <div class="text-weight-medium q-mb-sm" style="font-size: 0.85rem">
                {{ $t('signature.export.confirmWarning') }}
              </div>
              <div class="text-grey-7" style="font-size: 0.75rem; line-height: 1.4">
                {{ $t('signature.export.confirmDetail') }}
              </div>
            </div>

            <!-- 风险提示 -->
            <div class="q-mt-md">
              <div class="text-subtitle2 text-weight-medium q-mb-sm" style="font-size: 0.9rem">
                {{ $t('signature.export.confirmRisks') }}
              </div>
              <div class="text-grey-7" style="font-size: 0.75rem; line-height: 1.5">
                <div class="q-mb-xs">{{ $t('signature.export.riskItem1') }}</div>
                <div class="q-mb-xs">{{ $t('signature.export.riskItem2') }}</div>
                <div>{{ $t('signature.export.riskItem3') }}</div>
              </div>
            </div>
          </q-card-section>
        </q-scroll-area>

        <!-- 按钮区域 -->
        <q-card-actions align="right" class="q-px-md q-py-sm" style="gap: 8px">
          <q-btn
            flat
            :label="$t('signature.export.cancelButton')"
            color="primary"
            v-close-popup
            size="sm"
            class="text-caption"
          />
          <q-btn
            unelevated
            :label="$t('signature.export.confirmButton')"
            color="negative"
            @click="
              () => {
                showExportConfirmDialog = false;
                performExport();
              }
            "
            size="sm"
            class="text-caption"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 自定义菜单 -->
    <ContextMenu v-model="contextMenuVisible" :x="menuX" :y="menuY">
      <ContextMenuItem
        icon="edit"
        @click="
          () => {
            contextMenuVisible = false;
            handleEdit();
          }
        "
      >
        {{ $t('signature.page.edit') }}
      </ContextMenuItem>
      <ContextMenuItem
        icon="drive_file_move"
        @click="
          () => {
            contextMenuVisible = false;
            handleExport();
          }
        "
      >
        {{ $t('signature.page.export') }}
      </ContextMenuItem>
      <ContextMenuItem
        icon="delete"
        :is-negative="true"
        :show-divider="false"
        @click="
          () => {
            contextMenuVisible = false;
            handleDelete();
          }
        "
      >
        {{ $t('signature.page.delete') }}
      </ContextMenuItem>
    </ContextMenu>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import SignatureFormDialog from 'src/components/SignatureFormDialog.vue';
import ContextMenu from 'src/components/ContextMenu.vue';
import ContextMenuItem from 'src/components/ContextMenuItem.vue';
import { useSignatureStore } from 'src/stores/signature-store';
import {
  getSignaturesList,
  decryptSignatureData,
  getSignatureImage,
  deleteSignature,
  updateSignatureSort as updateSignatureSortApi,
  exportSignature,
  importSignature,
  confirmImportSignature,
} from 'boot/query/signature-query';
import type { Signature } from 'src/types/signature';

// 声明 File System Access API 相关类型
declare global {
  interface Window {
    showSaveFilePicker: (options?: {
      suggestedName?: string;
      types?: Array<{
        description: string;
        accept: Record<string, string[]>;
      }>;
    }) => Promise<FileSystemFileHandle>;
  }
}

interface FileSystemCreateWritableOptions {
  keepExistingData?: boolean;
}

interface FileSystemFileHandle {
  createWritable: (options?: FileSystemCreateWritableOptions) => Promise<FileSystemWritableFileStream>;
}

interface FileSystemWritableFileStream extends WritableStream {
  write: (data: BufferSource | Blob | string) => Promise<void>;
  close: () => Promise<void>;
}

const q = useQuasar();
const { t: $t } = useI18n();
const signatureStore = useSignatureStore();

// 系统类型判断
const isMacOS = computed(() => {
  if (process.env.MODE === 'electron') {
    return process.platform === 'darwin' || navigator.platform === 'MacIntel';
  }
  return navigator.platform === 'MacIntel';
});

// ========== 列表数据状态 ==========

// 页面数据加载状态 - 绑定加载动画显示
const loading = ref(false);

// 页面数据加载错误状态 - 绑定错误提示显示
const error = ref(false);

// 签名列表数据 - 由具体数据流实现填充，绑定到签名卡片列表渲染
const signatureList = ref<Signature[]>([]);

// 签名排序信息 Map - 存储每个签名 ID 对应的排序时间戳
// 用于排序后持久化到后端
const signatureSortMap = ref<Map<string, number>>(new Map());

// 图片 URL Map - 存储每个图片路径对应的 Blob URL
const imageUrls = ref<Map<string, string>>(new Map());

// 拖动排序中的状态 - 防止重复请求
const isSortingUpdating = ref(false);

// 拖动进行中的状态 - 用于在动画期间保持被拖动元素的高 z-index
const isDraggingInProgress = ref(false);

// ========== 对话框和菜单状态 ==========

// 表单对话框显示状态
const showFormDialog = ref(false);

// 当前编辑/查看的签名对象 - 传递给 SignatureFormDialog，null 表示创建模式
const selectedSignature = ref<Signature | null>(null);

// 导入对话框显示状态
const showImportDialogVisible = ref(false);

// 导入文件 - 绑定文件选择器
const importFile = ref<File | null>(null);

// 导入操作加载状态 - 绑定导入按钮的 loading 属性
const importing = ref(false);

// 图片预览对话框显示状态
const showImagePreview = ref(false);

// 预览图片 URL - 绑定到预览对话框的 img 标签
const previewImageUrl = ref('');

// 导出确认对话框显示状态
const showExportConfirmDialog = ref(false);

// ========== 上下文菜单状态 ==========

// 上下文菜单显示状态
const contextMenuVisible = ref(false);

// 当前上下文菜单指向的签名对象 - 用于菜单操作
const contextMenuSignature = ref<Signature | null>(null);

// 菜单项 DOM 引用映射 - 签名 ID -> DOM 元素
const contextMenuRefs = new Map<string, HTMLElement>();

// 当前打开菜单的签名 ID - 用于防止重复打开
const activeMenuSignatureId = ref<string | null>(null);

// 自定义菜单的位置坐标
const menuX = ref(0);
const menuY = ref(0);

// ========== 生命周期 ==========

onMounted(() => {
  // 加载初始数据
  loadSignatures();

  // 注册 SSE 回调
  signatureStore.registerSseCallback(handleSseUpdate);
});

onUnmounted(() => {
  // 注销 SSE 回调
  signatureStore.unregisterSseCallback();

  // 清理图片 Blob URL
  imageUrls.value.forEach((url) => {
    URL.revokeObjectURL(url);
  });
  imageUrls.value.clear();
});

// ========== 数据加载和同步 ==========

/**
 * 从加密签名数据中解密并构建 Signature 对象集合
 * 这是一个公共方法，被 loadSignatures 和 handleSseUpdate 共享调用
 *
 * @param encryptedSignatures 后端返回的加密签名对象
 * @returns 解密后的签名 Map (id -> Signature)
 */
async function decryptAndBuildSignatureMap(encryptedSignatures: Record<string, any>): Promise<Map<string, Signature>> {
  const signatureMap = new Map<string, Signature>();

  for (const [encryptedId, entry] of Object.entries(encryptedSignatures)) {
    try {
      // 兼容新旧格式
      let encryptedValue: string;

      if (typeof entry === 'string') {
        // 旧格式：直接是字符串
        encryptedValue = entry;
      } else if (typeof entry === 'object' && entry !== null) {
        // 新格式：是 SignatureStorageEntry 对象
        encryptedValue = (entry as any).value || '';
      } else {
        console.warn(`Unrecognized signature entry format for id: ${encryptedId}`);
        continue;
      }

      if (!encryptedValue) {
        console.warn(`Empty encrypted value for signature with id: ${encryptedId}`);
        continue;
      }

      // 解密 value 值（传递encryptedId以使用动态密钥解密）
      const decryptedJson = await decryptSignatureData(encryptedValue, encryptedId);
      if (!decryptedJson) {
        console.warn(`Failed to decrypt signature with id: ${encryptedId}`);
        continue;
      }

      // 解析 JSON
      const signatureData = JSON.parse(decryptedJson);

      // 创建 Signature 对象
      const signature: Signature = {
        id: encryptedId,
        name: signatureData.name,
        intro: signatureData.intro,
        cardImage: signatureData.cardImage || null,
      };

      signatureMap.set(encryptedId, signature);

      // 异步获取图片 URL（不阻塞列表显示）
      if (signatureData.cardImage) {
        loadImageUrl(signatureData.cardImage);
      }
    } catch (err) {
      console.error(`Failed to process signature with id ${encryptedId}:`, err);
    }
  }

  return signatureMap;
}

/**
 * 从签名 Map 中提取排序后的列表
 * @param encryptedSignatures 原始的加密签名对象（用于提取排序信息）
 * @param signatureMap 已解密的签名 Map
 * @returns 排序后的签名数组
 */
function extractSortedSignatures(
  encryptedSignatures: Record<string, any>,
  signatureMap: Map<string, Signature>
): Signature[] {
  // 构建带排序信息的数组
  const signatureWithSort: Array<{ signature: Signature; sortTime: number }> = [];

  for (const [encryptedId, entry] of Object.entries(encryptedSignatures)) {
    const signature = signatureMap.get(encryptedId);
    if (!signature) continue;

    // 提取排序时间戳
    let sortTime = 0;
    if (typeof entry === 'object' && entry !== null) {
      sortTime = (entry as any).sort?.time || 0;
    }

    signatureWithSort.push({ signature, sortTime });
  }

  // 按排序时间戳排序（升序：最早创建的在前面）
  signatureWithSort.sort((a, b) => {
    // 如果排序时间都是 0（旧格式迁移或未初始化），则按 ID 排序保持稳定性
    if (a.sortTime === 0 && b.sortTime === 0) {
      return a.signature.id.localeCompare(b.signature.id);
    }
    return a.sortTime - b.sortTime;
  });

  // 提取排序后的签名列表
  return signatureWithSort.map((item) => item.signature);
}

/**
 * 加载签名列表数据
 * 流程：
 * 1. 从后端获取加密的 key-value 对
 * 2. 逐个解密 value 值
 * 3. 将 JSON 字符串解析为 Signature 对象
 * 4. 获取图片 Blob URL
 */
async function loadSignatures() {
  loading.value = true;
  error.value = false;

  try {
    // 步骤1: 获取加密的签名列表
    const encryptedSignatures = await getSignaturesList();
    if (!encryptedSignatures) {
      error.value = true;
      return;
    }

    // 步骤2: 解密并构建签名 Map
    const signatureMap = await decryptAndBuildSignatureMap(encryptedSignatures);

    // 步骤3: 排序签名列表
    const sortedSignatures = extractSortedSignatures(encryptedSignatures, signatureMap);

    // 步骤4: 赋值给响应式变量
    signatureList.value = sortedSignatures;
  } catch (err) {
    console.error('[loadSignatures] 异常:', err);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

/**
 * SSE 更新回调
 * 当后端配置变化时（通过 SSE），触发此回调
 * 执行增量更新而不是全量重新加载，避免列表闪烁
 *
 * 增量更新的优势：
 * - 避免完整替换数组导致的 DOM 闪烁
 * - 保留现有的排序状态和展开/收起状态
 * - 只更新真正有变化的项目
 * - 支持排序的增量更新
 */
async function handleSseUpdate() {
  // 如果正在进行排序更新，跳过此次 SSE 更新以避免竞态条件导致的数据混乱
  // 排序更新完成后，会触发后端 SSE 推送，此时再执行增量更新
  if (isSortingUpdating.value) {
    console.debug('[SSE] Skipping update during sort, will be handled by sort callback');
    return;
  }

  console.debug('[SSE] Signature list updated, performing incremental update...');
  try {
    // 获取最新的加密签名列表
    const encryptedSignatures = await getSignaturesList();
    if (!encryptedSignatures) {
      console.warn('[SSE] Failed to fetch updated signatures');
      return;
    }

    // 解密并构建新的签名 Map
    const newSignaturesMap = await decryptAndBuildSignatureMap(encryptedSignatures);

    // 执行增量更新（传递加密签名对象以获取排序信息）
    updateSignaturesIncremental(newSignaturesMap, encryptedSignatures);
  } catch (err) {
    console.error('[SSE] Incremental update failed:', err);
    // 降级方案：全量重新加载
    console.debug('[SSE] Falling back to full reload');
    await loadSignatures();
  }
}

/**
 * 从加密签名对象中提取排序信息
 * @param encryptedSignatures 原始的加密签名对象
 * @returns 签名 ID -> 排序时间戳 的 Map
 */
function extractSortTimeMap(encryptedSignatures: Record<string, any>): Map<string, number> {
  const sortTimeMap = new Map<string, number>();

  for (const [encryptedId, entry] of Object.entries(encryptedSignatures)) {
    let sortTime = 0;
    if (typeof entry === 'object' && entry !== null) {
      sortTime = (entry as any).sort?.time || 0;
    }
    sortTimeMap.set(encryptedId, sortTime);
  }

  return sortTimeMap;
}

/**
 * 检查排序是否改变，如果改变则应用新的排序
 * 这个函数不会重新创建数组或元素，只会调整元素的顺序
 *
 * @param newSortTimeMap 新的排序信息 Map (id -> sortTime)
 * @returns 是否发生了排序变化
 */
function checkAndApplySortOrder(newSortTimeMap: Map<string, number>): boolean {
  // 构建当前列表的排序 Map（按当前顺序生成）
  const currentSortTimeMap = new Map<string, number>();
  signatureList.value.forEach((sig, index) => {
    currentSortTimeMap.set(sig.id, index);
  });

  // 构建新的排序顺序（基于新的排序时间戳）
  const itemsWithNewSort: Array<{ id: string; sortTime: number; item: Signature }> = [];

  for (const sig of signatureList.value) {
    const sortTime = newSortTimeMap.get(sig.id) || 0;
    itemsWithNewSort.push({ id: sig.id, sortTime, item: sig });
  }

  // 按新的排序时间戳排序
  itemsWithNewSort.sort((a, b) => {
    // 如果排序时间都是 0，则按当前顺序保持不变
    if (a.sortTime === 0 && b.sortTime === 0) {
      return currentSortTimeMap.get(a.id)! - currentSortTimeMap.get(b.id)!;
    }
    return a.sortTime - b.sortTime;
  });

  // 检查排序是否真的改变了
  let sortOrderChanged = false;
  for (let i = 0; i < signatureList.value.length; i++) {
    if (signatureList.value[i].id !== itemsWithNewSort[i].id) {
      sortOrderChanged = true;
      break;
    }
  }

  // 如果排序改变，更新列表
  if (sortOrderChanged) {
    console.debug('[SSE] Sort order changed, applying new sort order');
    console.debug(
      '[SSE] Old order:',
      signatureList.value.map((s) => s.id)
    );
    console.debug(
      '[SSE] New order:',
      itemsWithNewSort.map((i) => i.id)
    );

    // 重新排列数组元素（保持对象引用，只改变顺序）
    const newList = itemsWithNewSort.map((item) => item.item);
    signatureList.value.splice(0, signatureList.value.length, ...newList);
  }

  return sortOrderChanged;
}

/**
 * 执行签名列表的增量更新
 * 只更新有变化的项，避免整个列表重新渲染
 * 支持增量的排序更新
 *
 * @param newSignaturesMap 新的签名数据 Map
 * @param encryptedSignatures 原始的加密签名对象（用于提取排序信息）
 */
function updateSignaturesIncremental(
  newSignaturesMap: Map<string, Signature>,
  encryptedSignatures?: Record<string, any>
) {
  // 构建当前列表的 ID Set
  const currentIds = new Set(signatureList.value.map((s) => s.id));
  const newIds = new Set(newSignaturesMap.keys());

  // 检测需要删除的签名（存在于当前但不存在于新数据中）
  const toDeleteIds = new Set<string>();
  currentIds.forEach((id) => {
    if (!newIds.has(id)) {
      toDeleteIds.add(id);
    }
  });

  // 检测需要添加的签名（存在于新数据但不存在于当前中）
  const toAddIds = new Set<string>();
  newIds.forEach((id) => {
    if (!currentIds.has(id)) {
      toAddIds.add(id);
    }
  });

  // 检测需要更新的签名（检查数据是否发生变化）
  const toUpdateIds = new Set<string>();
  for (const id of currentIds) {
    if (newIds.has(id) && !toDeleteIds.has(id) && !toAddIds.has(id)) {
      const currentSig = signatureList.value.find((s) => s.id === id);
      const newSig = newSignaturesMap.get(id);
      if (currentSig && newSig) {
        // 比较关键字段，判断是否需要更新
        if (
          currentSig.name !== newSig.name ||
          currentSig.intro !== newSig.intro ||
          currentSig.cardImage !== newSig.cardImage
        ) {
          toUpdateIds.add(id);
        }
      }
    }
  }

  console.debug('[SSE] Incremental update detected:', {
    toAdd: toAddIds.size,
    toDelete: toDeleteIds.size,
    toUpdate: toUpdateIds.size,
  });

  // 如果没有任何变化且无排序信息，直接返回
  if (toAddIds.size === 0 && toDeleteIds.size === 0 && toUpdateIds.size === 0) {
    // 即使没有数据变化，也要检查排序是否改变
    if (encryptedSignatures) {
      const sortTimeMap = extractSortTimeMap(encryptedSignatures);
      const sortChanged = checkAndApplySortOrder(sortTimeMap);
      if (!sortChanged) {
        console.debug('[SSE] No changes detected, skipping update');
        return;
      }
    } else {
      console.debug('[SSE] No changes detected, skipping update');
      return;
    }
  }

  // 执行删除操作
  if (toDeleteIds.size > 0) {
    signatureList.value = signatureList.value.filter((s) => !toDeleteIds.has(s.id));
  }

  // 执行更新操作（保持在原位置）
  if (toUpdateIds.size > 0) {
    signatureList.value.forEach((sig) => {
      if (toUpdateIds.has(sig.id)) {
        const newSig = newSignaturesMap.get(sig.id);
        if (newSig) {
          // 使用 Object.assign 进行原地更新，保持数组引用不变
          Object.assign(sig, newSig);
        }
      }
    });
  }

  // 执行添加操作（添加到列表末尾）
  if (toAddIds.size > 0) {
    const addedSignatures = Array.from(toAddIds)
      .map((id) => newSignaturesMap.get(id))
      .filter((sig): sig is Signature => sig !== undefined);
    signatureList.value.push(...addedSignatures);
  }

  // 检查并应用排序更新
  if (encryptedSignatures) {
    const sortTimeMap = extractSortTimeMap(encryptedSignatures);
    checkAndApplySortOrder(sortTimeMap);
  }

  console.debug('[SSE] Incremental update completed');
}

/**
 * 异步加载图片 URL
 * 获取图片 Blob 并存储为 Blob URL
 */
async function loadImageUrl(imagePath: string) {
  // 如果已加载，跳过
  if (imageUrls.value.has(imagePath)) {
    return;
  }

  try {
    const blob = await getSignatureImage(imagePath);
    if (blob) {
      const url = URL.createObjectURL(blob);
      imageUrls.value.set(imagePath, url);
    }
  } catch (err) {
    console.warn(`Failed to load image from ${imagePath}:`, err);
  }
}

/**
 * 获取图片 URL
 * 从 Map 中查询已加载的 Blob URL
 */
function getImageUrl(filename: string): string {
  return imageUrls.value.get(filename) || '';
}
function handleCreate() {
  selectedSignature.value = null;
  showFormDialog.value = true;
}

/** 处理列表项信息区域点击 - 展开或收起菜单 */
function handleInfoClick(signature: Signature, event: MouseEvent) {
  event.preventDefault();
  event.stopPropagation();

  if (activeMenuSignatureId.value === signature.id && contextMenuVisible.value) {
    contextMenuVisible.value = false;
    activeMenuSignatureId.value = null;
  } else {
    calculateMenuPosition(null, event.clientX, event.clientY);
    contextMenuSignature.value = signature;
    contextMenuVisible.value = true;
    activeMenuSignatureId.value = signature.id;
  }
}

/** 处理列表项信息区域右键 - 展开或收起菜单 */
function handleInfoContextMenu(signature: Signature, event: MouseEvent) {
  event.preventDefault();
  event.stopPropagation();

  if (activeMenuSignatureId.value === signature.id && contextMenuVisible.value) {
    contextMenuVisible.value = false;
    activeMenuSignatureId.value = null;
  } else {
    calculateMenuPosition(null, event.clientX, event.clientY);
    contextMenuSignature.value = signature;
    contextMenuVisible.value = true;
    activeMenuSignatureId.value = signature.id;
  }
}

/** 计算菜单显示位置 */
function calculateMenuPosition(element: HTMLElement | null, clientX: number, clientY: number) {
  // 简单直接地设置菜单位置
  menuX.value = clientX;
  menuY.value = clientY;
}

/** 菜单项点击处理 - 关闭菜单并执行回调 */
function handleMenuItemClick(callback: () => void | Promise<void>) {
  contextMenuVisible.value = false;
  activeMenuSignatureId.value = null;
  contextMenuSignature.value = null;
  // 异步执行回调，确保菜单已关闭
  setTimeout(() => {
    callback();
  }, 0);
}

/** 打开编辑表单 */
function handleEdit() {
  if (contextMenuSignature.value) {
    // 为编辑的列表项创建深拷贝，避免过程中影响到原有列表项的实际内容
    // 深拷贝中的 cardImage 保持为字符串路径，不包含 File 对象
    const clonedSignature: Signature = {
      id: contextMenuSignature.value.id,
      name: contextMenuSignature.value.name,
      intro: contextMenuSignature.value.intro,
      cardImage: contextMenuSignature.value.cardImage,
    };
    selectedSignature.value = clonedSignature;
    showFormDialog.value = true;
  }
}

/** 删除签名 */
async function handleDelete() {
  if (!contextMenuSignature.value) return;

  // 确认删除对话框
  q.dialog({
    title: $t('signature.delete.confirmTitle'),
    message: $t('signature.delete.confirm'),
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    try {
      const success = await deleteSignature(contextMenuSignature.value!.id);
      if (success) {
        q.notify({
          type: 'positive',
          message: $t('signature.notify.deleteSuccess'),
          position: 'top',
        });
        // 直接从列表中移除该项，而不是全量重新加载（避免闪烁）
        // SSE 会推送后端的配置变更来最终确认删除
        const deletedId = contextMenuSignature.value!.id;
        signatureList.value = signatureList.value.filter((s) => s.id !== deletedId);
        console.debug(`[handleDelete] Signature ${deletedId} removed from list, waiting for SSE confirmation...`);
      } else {
        q.notify({
          type: 'negative',
          message: $t('signature.notify.deleteFailed'),
          position: 'top',
        });
      }
    } catch (error) {
      console.error('Failed to delete signature:', error);
      q.notify({
        type: 'negative',
        message: $t('signature.notify.unexpectedError'),
        position: 'top',
      });
    }
  });
}

/** 导出签名 - 降级方案（使用传统的下载方式） */
async function exportSignatureLegacy() {
  if (!contextMenuSignature.value) return;

  const signature = contextMenuSignature.value;

  try {
    // 调用导出 API
    const fileBlob = await exportSignature(signature.id);
    if (!fileBlob) {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.exportFailed'),
        position: 'top',
      });
      return;
    }

    // 使用 Web API 触发文件下载
    const url = URL.createObjectURL(fileBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${signature.name}.ktsign`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    // 清理
    URL.revokeObjectURL(url);

    // 显示成功通知
    q.notify({
      type: 'positive',
      message: $t('signature.notify.exportSuccess'),
      position: 'top',
    });

    console.debug('Signature exported successfully (legacy):', signature.name);
  } catch (error) {
    console.error('Failed to export signature (legacy):', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.unexpectedError'),
      position: 'top',
    });
  }
}

/** 导出签名 - 实际执行导出逻辑 */
async function performExport() {
  // 检查 API 是否可用
  if (typeof window.showSaveFilePicker !== 'function') {
    console.log('Browser does not support File System Access API, falling back to legacy export');
    return exportSignatureLegacy();
  }

  if (!contextMenuSignature.value) return;

  const signature = contextMenuSignature.value;

  try {
    // 调用导出 API
    const fileBlob = await exportSignature(signature.id);
    if (!fileBlob) {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.exportFailed'),
        position: 'top',
      });
      return;
    }

    try {
      // 打开系统的保存文件对话框
      const handle = await window.showSaveFilePicker({
        suggestedName: `${signature.name}.ktsign`,
        types: [
          {
            description: $t('signature.notify') + ' (.ktsign)',
            accept: { 'application/octet-stream': ['.ktsign'] },
          },
        ],
      });

      // 写入文件
      const writable = await handle.createWritable();
      await writable.write(fileBlob);
      await writable.close();

      // 文件成功保存后再通知
      q.notify({
        type: 'positive',
        message: $t('signature.notify.exportSuccess'),
        position: 'top',
      });

      console.debug('Signature exported successfully:', signature.name);
    } catch (err) {
      // 用户取消选择文件时不显示错误
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      throw err;
    }
  } catch (error) {
    console.error('Failed to export signature:', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.unexpectedError'),
      position: 'top',
    });
  }
}

/** 导出签名 - 显示安全确认对话框后执行导出 */
async function handleExport() {
  if (!contextMenuSignature.value) return;

  // 显示安全确认对话框
  showExportConfirmDialog.value = true;
}

/** 打开导入对话框 */
function showImportDialog() {
  importFile.value = null;
  showImportDialogVisible.value = true;
}

/** 导入签名 */
async function handleImport() {
  if (!importFile.value) return;

  importing.value = true;

  try {
    // 1. 读取文件内容
    const fileContent = await importFile.value.text();

    // 2. 调用导入 API
    const result = await importSignature(importFile.value);

    if (result === false) {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.importFailed') || 'Import failed',
        position: 'top',
      });
      return;
    }

    // 3. 处理冲突情况
    if (result.conflict) {
      console.warn('Signature conflict detected:', result.name);

      // 显示冲突对话框，让用户选择是否覆盖
      q.dialog({
        title: $t('signature.import.conflictTitle') || 'Signature Exists',
        message:
          $t('signature.import.conflictMessage', { name: result.name }) ||
          `Signature "${result.name}" already exists. Overwrite?`,
        cancel: true,
        persistent: true,
      }).onOk(async () => {
        // 用户选择覆盖
        const confirmResult = await confirmImportSignature(result.encryptedId, fileContent, true);

        if (confirmResult && confirmResult.success) {
          q.notify({
            type: 'positive',
            message: $t('signature.notify.importSuccess') || 'Signature imported successfully',
            position: 'top',
          });
          showImportDialogVisible.value = false;
          importFile.value = null;
        } else {
          q.notify({
            type: 'negative',
            message: $t('signature.notify.importFailed') || 'Import failed',
            position: 'top',
          });
        }
      });
      return;
    }

    // 4. 导入成功（无冲突）
    if (result.success) {
      q.notify({
        type: 'positive',
        message: $t('signature.notify.importSuccess') || 'Signature imported successfully',
        position: 'top',
      });
      showImportDialogVisible.value = false;
      importFile.value = null;

      console.debug('Signature imported successfully:', result.name);
    } else {
      q.notify({
        type: 'negative',
        message: $t('signature.notify.importFailed') || 'Import failed',
        position: 'top',
      });
    }
  } catch (error) {
    console.error('Failed to import signature:', error);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.unexpectedError') || 'An unexpected error occurred',
      position: 'top',
    });
  } finally {
    importing.value = false;
  }
}

/** 表单成功提交后的回调 - 刷新列表 */
function handleFormSuccess() {
  // 不再直接调用 loadSignatures() 全量重新加载（会导致列表闪烁）
  // 而是依赖 SSE 推送的数据通过 handleSseUpdate() 增量更新
  // SSE 会在后端配置变更后推送新数据，避免前端重复加载
  console.debug('[handleFormSuccess] Waiting for SSE update to sync signature list...');
}

// ========== 拖动排序相关函数 ==========

// 当前被拖动的签名
let draggedSignature: Signature | null = null;
// 拖动时的 DOM 元素
let draggedElement: HTMLElement | null = null;
// 插入位置指示器
let dropIndicator: HTMLElement | null = null;
// 被拖动元素在列表中的原始索引
let draggedElementOriginalIndex = -1;
// 当前显示的边框方向 ('top' | 'bottom' | null)
let currentDropPosition: 'top' | 'bottom' | null = null;
// 上次计算的目标卡片 DOM
let lastTargetCard: HTMLElement | null = null;
// ✨ 关键：记录当前指示的目标卡片及其对应的插入位置
// 这样在 drop 时可以获取到准确的值，不会被其他 dragenter 改变
let currentTargetCardForDrop: HTMLElement | null = null;
let currentInsertPositionForDrop: 'top' | 'bottom' | null = null;

/**
 * 判断鼠标在卡片中的位置 - 返回应该的插入方向
 * 使用简单的中点 50% 分界线，确保实时反馈
 *
 * @param mouseY 鼠标相对于卡片顶部的 Y 坐标
 * @param cardHeight 卡片高度
 * @returns 'top' 表示上方，'bottom' 表示下方
 */
function getDropPosition(mouseY: number, cardHeight: number): 'top' | 'bottom' {
  // 简单的中点判断：越过 50% 就切换
  return mouseY < cardHeight * 0.5 ? 'top' : 'bottom';
}

/** 处理拖动开始 */
function handleDragStart(event: DragEvent, signature: Signature) {
  draggedSignature = signature;
  draggedElement = event.target as HTMLElement;
  draggedElementOriginalIndex = signatureList.value.findIndex((s) => s.id === signature.id);
  isDraggingInProgress.value = true;
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.setData('text/html', (event.target as HTMLElement).innerHTML);
  }
  // 添加拖动中的视觉效果
  const cardElement = (event.target as HTMLElement).closest('.draggable-card') as HTMLElement;
  if (cardElement) {
    cardElement.classList.add('dragging');
  }
}

/** 处理拖动结束 */
function handleDragEnd(event: DragEvent) {
  const target = event.target as HTMLElement;
  const cardElement = target.closest('.draggable-card') as HTMLElement;
  if (cardElement) {
    cardElement.classList.remove('dragging');
  }

  // 移除所有高亮
  document.querySelectorAll('.draggable-card').forEach((card) => {
    card.classList.remove('drag-over-top');
    card.classList.remove('drag-over-bottom');
  });

  // 移除插入指示器
  if (dropIndicator) {
    dropIndicator.remove();
    dropIndicator = null;
  }

  draggedSignature = null;
  draggedElement = null;
  draggedElementOriginalIndex = -1;
  currentDropPosition = null;
  lastTargetCard = null;
  // ✨ 清除 drop 相关的记录
  currentTargetCardForDrop = null;
  currentInsertPositionForDrop = null;

  // 在动画完成后清除拖动状态（300ms 是 list-move 动画时长）
  setTimeout(() => {
    isDraggingInProgress.value = false;
  }, 350);
}

/** 处理拖动悬停 */
function handleDragOver(event: DragEvent) {
  event.preventDefault();
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move';
  }
}

/** 处理拖动进入 */
function handleDragEnter(event: DragEvent) {
  if (!draggedSignature) return;

  const target = event.target as HTMLElement;
  const cardElement = target.closest('.draggable-card');

  if (!cardElement || cardElement === draggedElement) {
    return;
  }

  // 获取鼠标在卡片中的位置
  const rect = cardElement.getBoundingClientRect();
  const mouseY = event.clientY - rect.top;

  // ✨ 简化逻辑：直接判断位置，实时反馈给用户
  const dropPosition = getDropPosition(mouseY, rect.height);

  // 只在位置改变时才更新 DOM（减少重排，提高性能）
  if (currentDropPosition === dropPosition && lastTargetCard === cardElement) {
    return;
  }

  currentDropPosition = dropPosition;
  lastTargetCard = cardElement as HTMLElement;

  // ✨ 记录当前卡片和插入位置，供 drop 时使用
  currentTargetCardForDrop = cardElement as HTMLElement;
  currentInsertPositionForDrop = dropPosition;

  // 移除之前的高亮
  document.querySelectorAll('.draggable-card').forEach((card) => {
    card.classList.remove('drag-over-top');
    card.classList.remove('drag-over-bottom');
  });

  // ✨ 直接显示当前位置对应的边框
  if (dropPosition === 'top') {
    cardElement.classList.add('drag-over-top');
  } else if (dropPosition === 'bottom') {
    cardElement.classList.add('drag-over-bottom');
  }
}

/** 处理拖动离开 */
function handleDragLeave(event: DragEvent) {
  const target = event.target as HTMLElement;
  const cardElement = target.closest('.draggable-card');

  if (cardElement) {
    // 检查是否真的离开了卡片
    if (!cardElement.contains(event.relatedTarget as Node)) {
      cardElement.classList.remove('drag-over-top');
      cardElement.classList.remove('drag-over-bottom');

      // 如果离开的是当前追踪的卡片，清除状态
      if (cardElement === lastTargetCard) {
        currentDropPosition = null;
        lastTargetCard = null;
        // ✨ 如果离开的是 drop 目标卡片，也清除 drop 相关状态
        currentTargetCardForDrop = null;
        currentInsertPositionForDrop = null;
      }
    }
  }
}

/** 处理放下 - 插入排序 */
async function handleDrop(event: DragEvent, targetSignature: Signature) {
  event.preventDefault();
  event.stopPropagation();

  const target = event.target as HTMLElement;
  const cardElement = target.closest('.draggable-card');

  if (cardElement) {
    cardElement.classList.remove('drag-over-top');
    cardElement.classList.remove('drag-over-bottom');
  }

  if (!draggedSignature || draggedSignature.id === targetSignature.id) {
    return;
  }

  // 查找两个签名在列表中的索引
  const draggedIndex = signatureList.value.findIndex((s) => s.id === draggedSignature!.id);
  const targetIndex = signatureList.value.findIndex((s) => s.id === targetSignature.id);

  if (draggedIndex === -1 || targetIndex === -1) {
    return;
  }

  // 获取鼠标在卡片中的位置
  const rect = (event.target as HTMLElement).closest('.draggable-card')!.getBoundingClientRect();
  const mouseY = event.clientY - rect.top;
  // ✨ 关键：使用 dragenter 时记录的准确值
  // 这样可以确保 drop 时使用的是实际显示的边框对应的位置
  let insertAfter: boolean;

  if (currentTargetCardForDrop === cardElement && currentInsertPositionForDrop !== null) {
    // ✨ 使用 dragenter 时记录的插入位置
    insertAfter = currentInsertPositionForDrop === 'bottom';
  } else {
    // 备用方案（如果没有记录上次的值，重新计算）
    const dropPosition = getDropPosition(mouseY, rect.height);
    insertAfter = dropPosition === 'bottom';
  }

  // 计算新位置将会是什么
  let newIndex: number;
  if (draggedIndex < targetIndex) {
    // 从前往后拖动：删除后目标索引会减 1
    newIndex = insertAfter ? targetIndex : targetIndex - 1;
  } else {
    // 从后往前拖动：删除后目标索引不变
    newIndex = insertAfter ? targetIndex + 1 : targetIndex;
  }

  // 检查是否真的改变了位置
  // 如果新位置等于原位置，就不需要更新
  if (newIndex === draggedIndex) {
    return;
  }

  // 从原位置移除
  const [draggedItem] = signatureList.value.splice(draggedIndex, 1);

  // 重新计算新位置（因为移除后索引会改变）
  let finalNewIndex = signatureList.value.findIndex((s) => s.id === targetSignature.id);

  if (insertAfter) {
    finalNewIndex += 1;
  }

  // 插入到新位置
  signatureList.value.splice(finalNewIndex, 0, draggedItem);

  // 只在位置真的改变时才更新排序
  await updateSignatureSort();
}

/** 更新签名排序 - 生成新的时间戳并提交到后端 */
async function updateSignatureSort() {
  if (isSortingUpdating.value) {
    return;
  }

  isSortingUpdating.value = true;

  try {
    // 为排序后的签名列表生成新的排序时间戳
    // 使用递增的时间戳，从当前时间开始，每个签名间隔 1000 毫秒
    const baseTime = Math.floor(Date.now() / 1000);
    const newSortOrder = signatureList.value.map((sig, index) => ({
      id: sig.id,
      sortTime: baseTime + index,
    }));

    // 调用后端 API 更新排序
    const success = await updateSignatureSortApi(newSortOrder);
    if (!success) {
      q.notify({
        type: 'negative',
        message: '排序更新失败，请重试',
        position: 'top',
      });
      // 还原到原始顺序
      await loadSignatures();
    } else {
      q.notify({
        type: 'positive',
        message: '排序已更新',
        position: 'top',
        timeout: 300,
      });
    }
  } catch (error) {
    console.error('Failed to update signature sort:', error);
    q.notify({
      type: 'negative',
      message: '排序更新失败',
      position: 'top',
    });
    // 还原到原始顺序
    await loadSignatures();
  } finally {
    isSortingUpdating.value = false;
  }
}

/** 预览签名图片 */
function handleImagePreview(filename: string) {
  const imageUrl = getImageUrl(filename);
  if (!imageUrl || imageUrl.trim() === '') {
    q.notify({
      type: 'info',
      message: $t('signature.page.noImagePreviewTip'),
      position: 'top',
    });
    return;
  }
  previewImageUrl.value = imageUrl;
  showImagePreview.value = true;
}
</script>

<style lang="scss" scoped>
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

.draggable-card {
  user-select: none;
  transition: opacity 0.2s, transform 0.2s;
  position: relative;
}

.draggable-card.dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

/** 在目标上方插入 - 显示上边框 */
.draggable-card.drag-over-top {
  border-top: 3px solid #1976d2;
  padding-top: 0px;
  background-color: rgba(25, 118, 210, 0.08);
}

/** 在目标下方插入 - 显示下边框 */
.draggable-card.drag-over-bottom {
  border-bottom: 3px solid #1976d2;
  padding-bottom: 0px;
  background-color: rgba(25, 118, 210, 0.08);
}

/** TransitionGroup 动画样式 */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-move {
  transition: transform 0.3s ease;
  opacity: 0.6;
}

.scroll-h-hide {
  :deep(.q-scrollarea__thumb--h) {
    display: none !important;
  }
}
</style>
