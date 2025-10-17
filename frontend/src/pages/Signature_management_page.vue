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
import type { Signature } from 'src/types/signature';

const q = useQuasar();
const { t: $t } = useI18n();

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

onMounted(() => {
  loadSignatures();
  // TODO: 添加 SSE 事件监听器
});

onUnmounted(() => {
  // TODO: 移除 SSE 事件监听器
});

/** 加载签名列表数据 */
async function loadSignatures() {
  loading.value = true;
  error.value = false;

  try {
    // TODO: 具体数据加载逻辑由业务层实现
    signatureList.value = [];
  } catch (err) {
    console.error('[loadSignatures] 异常:', err);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

/** 打开创建签名表单 */
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
  // TODO: 菜单位置计算逻辑
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

  // TODO: 具体删除逻辑由业务层实现
}

/** 导出签名 */
async function handleExport() {
  if (!contextMenuSignature.value) return;

  // TODO: 具体导出逻辑由业务层实现
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

/** 获取签名图片 URL - 由具体业务层实现 */
function getImageUrl(filename: string): string {
  // TODO: 具体 URL 生成逻辑由业务层实现
  return '';
}

/** 预览签名图片 */
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
