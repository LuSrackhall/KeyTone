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
        <TransitionGroup name="list" tag="div" class="space-y-2">
          <q-card
            v-for="signature in signatureList"
            :key="signature.id"
            class="cursor-move hover:shadow-lg transition-all relative draggable-card"
            :style="{ minHeight: '60px' }"
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
        </TransitionGroup>
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
        <q-item clickable v-close-popup @click="handleEdit">
          <q-item-section avatar>
            <q-icon name="edit" />
          </q-item-section>
          <q-item-section>{{ $t('signature.page.edit') }}</q-item-section>
        </q-item>
        <q-item clickable v-close-popup @click="handleExport">
          <q-item-section avatar>
            <q-icon name="drive_file_move" />
          </q-item-section>
          <q-item-section>{{ $t('signature.page.export') }}</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable v-close-popup @click="handleDelete" class="text-negative">
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
import SignatureFormDialog from 'src/components/SignatureFormDialog.vue';
import { useSignatureStore } from 'src/stores/signature-store';
import {
  getSignaturesList,
  decryptSignatureData,
  getSignatureImage,
  deleteSignature,
  updateSignatureSort as updateSignatureSortApi,
} from 'boot/query/signature-query';
import type { Signature } from 'src/types/signature';

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

// ========== 上下文菜单状态 ==========

// 上下文菜单显示状态
const contextMenuVisible = ref(false);

// 当前上下文菜单指向的签名对象 - 用于菜单操作
const contextMenuSignature = ref<Signature | null>(null);

// 菜单项 DOM 引用映射 - 签名 ID -> DOM 元素
const contextMenuRefs = new Map<string, HTMLElement>();

// 当前打开菜单的签名 ID - 用于防止重复打开
const activeMenuSignatureId = ref<string | null>(null);

// 虚拟菜单参考元素引用 - 用于精确定位菜单
const virtualMenuRef = ref<HTMLElement | null>(null);

// 菜单锚点位置 - 决定菜单相对于虚拟点的附着位置
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

// 菜单自身位置 - 决定菜单的哪一边与锚点对齐
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

      // 解密 value 值
      const decryptedJson = await decryptSignatureData(encryptedValue);
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

/** 计算菜单最佳显示位置 */
function calculateMenuPosition(element: HTMLElement | null, clientX: number, clientY: number) {
  if (!virtualMenuRef.value) return;

  // 设置虚拟菜单位置到点击的确切坐标
  virtualMenuRef.value.style.left = clientX + 'px';
  virtualMenuRef.value.style.top = clientY + 'px';

  // 判断菜单是否可能超出视口边界，动态调整锚点和自身位置
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const menuWidth = 150; // 菜单估计宽度
  const menuHeight = 120; // 菜单估计高度

  // 判断是否需要向左展开
  if (clientX + menuWidth > viewportWidth) {
    menuAnchor.value = 'bottom right';
    menuSelf.value = 'top right';
  } else {
    menuAnchor.value = 'bottom left';
    menuSelf.value = 'top left';
  }

  // 判断是否需要向上展开
  if (clientY + menuHeight > viewportHeight) {
    // 如果已经是向右展开，改为右上
    if (menuAnchor.value === 'bottom right') {
      menuAnchor.value = 'top right';
      menuSelf.value = 'bottom right';
    } else {
      menuAnchor.value = 'top left';
      menuSelf.value = 'bottom left';
    }
  }
}

/** 打开编辑表单 */
function handleEdit() {
  if (contextMenuSignature.value) {
    selectedSignature.value = contextMenuSignature.value;
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
        // 重新加载列表
        await loadSignatures();
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

/** 导出签名 */
async function handleExport() {
  if (!contextMenuSignature.value) return;

  // 暂不实现具体逻辑，仅展示UI交互
  q.notify({
    type: 'info',
    message: $t('signature.export.comingSoon') || 'Export feature coming soon',
    position: 'top',
  });
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
    // TODO: 具体导入逻辑由业务层实现
  } catch (error) {
    console.error('Failed to import signature:', error);
  } finally {
    importing.value = false;
  }
}

/** 表单成功提交后的回调 - 刷新列表 */
function handleFormSuccess() {
  loadSignatures();
}

// ========== 拖动排序相关函数 ==========

// 当前被拖动的签名
let draggedSignature: Signature | null = null;
// 拖动时的 DOM 元素
let draggedElement: HTMLElement | null = null;
// 插入位置指示器
let dropIndicator: HTMLElement | null = null;

/** 处理拖动开始 */
function handleDragStart(event: DragEvent, signature: Signature) {
  draggedSignature = signature;
  draggedElement = event.target as HTMLElement;
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.setData('text/html', (event.target as HTMLElement).innerHTML);
  }
  // 添加拖动中的视觉效果
  const cardElement = (event.target as HTMLElement).closest('.draggable-card');
  if (cardElement) {
    cardElement.classList.add('dragging');
  }
}

/** 处理拖动结束 */
function handleDragEnd(event: DragEvent) {
  const target = event.target as HTMLElement;
  target.closest('.draggable-card')?.classList.remove('dragging');

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
  const midpoint = rect.height / 2;
  const mouseY = event.clientY - rect.top;

  // 移除之前的高亮
  document.querySelectorAll('.draggable-card').forEach((card) => {
    card.classList.remove('drag-over-top');
    card.classList.remove('drag-over-bottom');
  });

  // 根据鼠标位置添加高亮
  if (mouseY < midpoint) {
    // 上方插入
    cardElement.classList.add('drag-over-top');
  } else {
    // 下方插入
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
  const midpoint = rect.height / 2;
  const mouseY = event.clientY - rect.top;

  // 判断是否在下方插入
  const insertAfter = mouseY >= midpoint;

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

.draggable-card {
  user-select: none;
  transition: opacity 0.2s, transform 0.2s;
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
}
</style>
