<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card
      class="signature-picker-dialog [overflow:overlay] [&::-webkit-scrollbar]:hidden"
      style="width: 90%; max-width: 340px; max-height: 85vh; display: flex; flex-direction: column"
    >
      <!-- Header & Search (Sticky Top) -->
      <div class="sticky-top">
        <!-- Header -->
        <q-card-section class="bg-primary bg-opacity-90 text-white q-pa-sm sticky top-0 z-10 backdrop-blur-sm">
          <div class="text-subtitle1">{{ t('exportFlow.pickerDialog.title') }}</div>
        </q-card-section>

        <!-- Search Bar -->
        <q-card-section class="q-pa-sm">
          <!-- Description -->
          <div class="text-caption q-mb-sm">
            {{ t('exportFlow.pickerDialog.description') }}
          </div>

          <!-- Search Input -->
          <div class="q-mb-sm">
            <q-input
              v-model="searchQuery"
              filled
              dense
              :placeholder="t('exportFlow.pickerDialog.search')"
              icon="search"
              clearable
              size="sm"
            />
            <div class="q-mt-xs flex items-center" :class="requireAuthorization ? 'justify-between' : 'justify-end'">
              <!-- 授权申请按钮：仅在需要授权的再次导出时显示（左侧）-->
              <q-btn
                v-if="requireAuthorization"
                size="xs"
                flat
                color="deep-purple"
                icon="mail_outline"
                :label="t('exportFlow.pickerDialog.authRequest')"
                @click="onAuthRequest"
              />
              <!-- 导入授权按钮：仅在需要授权的再次导出时显示（中间）-->
              <q-btn
                v-if="requireAuthorization"
                size="xs"
                flat
                color="secondary"
                icon="key"
                :label="t('exportFlow.pickerDialog.importAuth')"
                @click="onImportAuth"
              />
              <!-- 创建签名按钮（右侧）-->
              <q-btn
                size="xs"
                flat
                color="primary"
                icon="add"
                :label="t('exportFlow.pickerDialog.createSignature')"
                @click="onCreateNew"
              />
            </div>
          </div>
        </q-card-section>
      </div>

      <q-card-section class="q-pa-md col-grow overflow-auto">
        <!-- Empty State -->
        <div v-if="filteredSignatures.length === 0 && !searchQuery" class="text-center q-pa-md">
          <q-icon name="mail" size="32px" color="grey-5" />
          <div class="text-caption text-grey q-mt-sm">
            {{ t('exportFlow.pickerDialog.emptyState') }}
          </div>
          <q-btn
            flat
            color="primary"
            :label="t('exportFlow.pickerDialog.createSignature')"
            icon="add"
            size="xs"
            class="q-mt-md"
            @click="onCreateNew"
          />
        </div>

        <!-- No Results -->
        <div v-else-if="filteredSignatures.length === 0 && searchQuery" class="text-center q-pa-md">
          <q-icon name="search_off" size="32px" color="grey-5" />
          <div class="text-caption text-grey q-mt-sm">
            {{ t('exportFlow.pickerDialog.noResults') }}
          </div>
        </div>

        <!-- Signatures List -->
        <div v-else class="col">
          <!-- Signature Cards -->
          <div v-for="sig in filteredSignatures" :key="sig.id" class="q-mb-xs">
            <q-card
              flat
              bordered
              :class="[
                'signature-card',
                selectedId === sig.id ? 'selected' : '',
                isSignatureDisabled(sig) ? 'disabled-card' : 'cursor-pointer',
              ]"
              @click="handleSignatureClick(sig.id)"
              style="display: flex; align-items: center; min-height: 70px"
            >
              <!-- Left: Image Area (fixed 60px) -->
              <div
                class="flex-shrink-0 flex items-center justify-center"
                :style="{
                  width: '60px',
                  height: '60px',
                  backgroundColor: '#f5f5f5',
                  borderRadius: '4px',
                  margin: '0 8px',
                  opacity: isSignatureDisabled(sig) ? 0.5 : 1,
                }"
              >
                <img
                  v-if="sig.image"
                  :src="sig.image"
                  :alt="sig.name"
                  style="width: 100%; height: 100%; object-fit: cover; border-radius: 4px"
                />
                <div
                  v-else
                  style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center"
                >
                  <q-icon name="image_not_supported" size="20px" color="grey-5" />
                </div>
              </div>

              <!-- Middle: Info Area (flex-grow) -->
              <div
                class="col flex flex-col justify-center q-py-2xs q-px-sm"
                :style="{ minWidth: 0, opacity: isSignatureDisabled(sig) ? 0.5 : 1 }"
              >
                <!-- Name: single line with horizontal scroll -->
                <div class="name-container text-caption text-weight-bold flex items-center" style="font-size: 0.9rem">
                  <div class="scrollable-x q-mr-xs">{{ sig.name }}</div>
                </div>
                <!-- Intro: max 2 lines with horizontal scroll -->
                <div class="intro-container text-caption text-grey q-mt-2xs" style="font-size: 0.75rem; min-width: 0">
                  <div class="scrollable-x line-clamp-2">
                    {{ sig.intro || t('exportFlow.pickerDialog.noIntro') }}
                  </div>
                </div>
                <div class="name-container text-caption text-weight-bold flex items-center" style="font-size: 0.9rem">
                  <q-badge
                    v-if="sig.isOriginalAuthor"
                    color="purple"
                    text-color="white"
                    :label="t('exportFlow.pickerDialog.originalAuthor')"
                    class="text-xs"
                  />
                  <q-badge
                    v-else-if="sig.isInAlbum"
                    color="teal"
                    text-color="white"
                    :label="t('exportFlow.pickerDialog.contributor')"
                    class="text-xs"
                  />
                  <!-- 未授权标签 -->
                  <q-badge
                    v-if="isSignatureDisabled(sig)"
                    color="grey"
                    text-color="white"
                    :label="t('exportFlow.pickerDialog.unauthorized')"
                    class="text-xs q-ml-xs"
                  />
                </div>
              </div>

              <!-- Right: Selection Indicator with Glow -->
              <div
                v-if="selectedId === sig.id"
                class="flex-shrink-0 selection-indicator-wrapper"
                style="margin-left: 8px; margin-right: 8px"
              >
                <div class="selection-glow"></div>
                <q-icon name="check_circle" size="20px" color="positive" class="selection-icon" />
              </div>
              <!-- Right: Lock icon for disabled signatures -->
              <div
                v-else-if="isSignatureDisabled(sig)"
                class="flex-shrink-0"
                style="margin-left: 8px; margin-right: 8px"
              >
                <q-icon name="lock" size="20px" color="grey" />
              </div>
            </q-card>
          </div>
        </div>
      </q-card-section>

      <!-- Actions (Sticky Bottom) -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs sticky-bottom" :class="['bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="t('exportFlow.pickerDialog.cancel')" color="primary" size="sm" @click="onCancel" />
        <q-btn
          unelevated
          :label="t('exportFlow.pickerDialog.confirm')"
          color="primary"
          size="sm"
          :disable="!selectedId"
          @click="onConfirmClick"
        />
      </q-card-actions>
    </q-card>

    <!-- Update Confirmation Dialog -->
    <q-dialog v-model="updateConfirmDialogVisible" persistent>
      <q-card style="min-width: 300px">
        <q-card-section>
          <div class="text-h6">{{ t('exportFlow.signatureUpdateConfirm.title') }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          {{ t('exportFlow.signatureUpdateConfirm.message') }}
        </q-card-section>

        <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
          <q-btn flat :label="t('exportFlow.common.cancel')" color="primary" v-close-popup />
          <q-btn
            flat
            :label="t('exportFlow.signatureUpdateConfirm.noUpdate')"
            color="primary"
            v-close-popup
            @click="handleUpdateConfirm(false)"
          />
          <q-btn
            flat
            :label="t('exportFlow.signatureUpdateConfirm.update')"
            color="primary"
            v-close-popup
            @click="handleUpdateConfirm(true)"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useSignatureStore } from 'src/stores/signature-store';
import { getSignaturesList, decryptSignatureData, getSignatureImage } from 'boot/query/signature-query';
import { GetAvailableSignaturesForExport, CheckSignatureInAlbum } from 'src/boot/query/keytonePkg-query';
import type { AvailableSignature } from 'src/types/export-flow';

interface Signature {
  id: string;
  name: string;
  intro?: string;
  image?: string;
  isAuthorized?: boolean; // Added for filtering
  isOriginalAuthor?: boolean; // Added for UI tag
  isInAlbum?: boolean; // Added for UI tag
}

interface SignaturePickerDialogProps {
  visible: boolean;
  signatures?: Signature[];
  albumPath?: string; // 专辑路径，用于获取签名状态
  requireAuthorization?: boolean; // 是否仅允许选择已授权的签名
}

interface SignaturePickerDialogEmits {
  (e: 'select', signatureId: string, updateContent: boolean): void;
  (e: 'createNew'): void;
  (e: 'cancel'): void;
  (e: 'importAuth'): void; // 从签名选择页面打开授权门控
  (e: 'authRequest'): void; // 打开授权申请流程
}

const props = withDefaults(defineProps<SignaturePickerDialogProps>(), {
  visible: false,
  signatures: () => [],
  albumPath: '',
  requireAuthorization: false,
});

const emit = defineEmits<SignaturePickerDialogEmits>();
const { t } = useI18n(); // 解构出 t 函数供模板使用
const signatureStore = useSignatureStore();

// UI State
const isVisible = ref(false);
const searchQuery = ref('');
const selectedId = ref('');
const loading = ref(false);
const localSignatures = ref<Signature[]>([]);
const imageUrls = ref<Map<string, string>>(new Map()); // cardImage path -> blob URL
const blobUrls = ref<Set<string>>(new Set()); // Track blob URLs for cleanup
const updateConfirmDialogVisible = ref(false);
const isImportingAuth = ref(false); // 标记是否正在导入授权，避免触发 cancel 事件

// Determine which signatures to use (real data or fallback to props)
const effectiveSignatures = computed(() => {
  if (localSignatures.value.length > 0) {
    return localSignatures.value;
  }
  return props.signatures || [];
});

// Filter signatures based on search query (only by name)
// 不再根据授权状态过滤，而是显示所有签名但禁用未授权的
const filteredSignatures = computed(() => {
  const result = effectiveSignatures.value;

  if (!searchQuery.value) {
    return result;
  }

  const query = searchQuery.value.toLowerCase();
  return result.filter((sig) => sig.name.toLowerCase().includes(query));
});

/**
 * 判断签名是否应该被禁用（需要授权但未授权的签名）
 */
function isSignatureDisabled(sig: Signature): boolean {
  // 如果不需要授权检查，所有签名都可用
  if (!props.requireAuthorization) {
    return false;
  }
  // 需要授权检查时，未授权的签名应该被禁用
  return !sig.isAuthorized;
}

/**
 * Load real signature data from store
 * Decrypts each signature and fetches image URLs
 */
async function loadSignaturesRealtime() {
  loading.value = true;
  try {
    // 1. Get encrypted signature list from backend
    const encryptedSignatures = await getSignaturesList();
    if (!encryptedSignatures) {
      console.warn('[SignaturePicker] Failed to fetch signatures');
      return;
    }

    // 1.5 Get authorization status if albumPath is provided
    const authMap = new Map<string, boolean>();
    const originalAuthorMap = new Map<string, boolean>();
    const inAlbumMap = new Map<string, boolean>();

    if (props.albumPath) {
      try {
        const availableSigs = await GetAvailableSignaturesForExport(props.albumPath);
        availableSigs.forEach((sig: AvailableSignature) => {
          authMap.set(sig.encryptedId, sig.isAuthorized);
          originalAuthorMap.set(sig.encryptedId, sig.isOriginalAuthor);
          inAlbumMap.set(sig.encryptedId, sig.isInAlbum);
        });
      } catch (err) {
        console.error('[SignaturePicker] Failed to fetch available signatures:', err);
      }
    }

    // 2. Decrypt and build signature list
    const tempSignatures: Signature[] = [];
    for (const [encryptedId, entry] of Object.entries(encryptedSignatures)) {
      try {
        // Handle both old (string) and new (object) formats
        let encryptedValue: string;
        if (typeof entry === 'string') {
          encryptedValue = entry;
        } else if (typeof entry === 'object' && entry !== null) {
          encryptedValue = (entry as any).value || '';
        } else {
          continue;
        }

        if (!encryptedValue) continue;

        // Decrypt value
        const decryptedJson = await decryptSignatureData(encryptedValue, encryptedId);
        if (!decryptedJson) {
          console.warn(`[SignaturePicker] Failed to decrypt signature ${encryptedId}`);
          continue;
        }

        // Parse JSON
        const signatureData = JSON.parse(decryptedJson);

        // Create signature object
        const signature: Signature = {
          id: encryptedId,
          name: signatureData.name,
          intro: signatureData.intro,
          image: signatureData.cardImage ? await getImageUrlForSignature(signatureData.cardImage) : undefined,
          isAuthorized: props.albumPath ? authMap.get(encryptedId) ?? false : true,
          isOriginalAuthor: props.albumPath ? originalAuthorMap.get(encryptedId) ?? false : false,
          isInAlbum: props.albumPath ? inAlbumMap.get(encryptedId) ?? false : false,
        };

        tempSignatures.push(signature);

        // Preload image URL if exists
        if (signatureData.cardImage) {
          preloadImageUrl(signatureData.cardImage);
        }
      } catch (err) {
        console.error(`[SignaturePicker] Error processing signature ${encryptedId}:`, err);
      }
    }

    // 3. Sort by sort.time if available
    tempSignatures.sort((a, b) => {
      const timeA = (encryptedSignatures[a.id] as any)?.sort?.time || 0;
      const timeB = (encryptedSignatures[b.id] as any)?.sort?.time || 0;
      if (timeA === 0 && timeB === 0) {
        return a.id.localeCompare(b.id);
      }
      return timeA - timeB;
    });

    localSignatures.value = tempSignatures;
    console.debug('[SignaturePicker] Loaded', tempSignatures.length, 'signatures');
  } catch (err) {
    console.error('[SignaturePicker] Error loading signatures:', err);
  } finally {
    loading.value = false;
  }
}

/**
 * Get or fetch image URL for a card image
 */
async function getImageUrlForSignature(cardImage: string): Promise<string | undefined> {
  if (!cardImage) return undefined;

  // Check cache first
  if (imageUrls.value.has(cardImage)) {
    return imageUrls.value.get(cardImage);
  }

  try {
    const result = await getSignatureImage(cardImage);
    if (result && result instanceof Blob) {
      const url = URL.createObjectURL(result);
      imageUrls.value.set(cardImage, url);
      blobUrls.value.add(url);
      return url;
    }
  } catch (err) {
    console.error('[SignaturePicker] Error fetching image:', err);
  }
  return undefined;
}

/**
 * Preload image URL asynchronously (non-blocking)
 */
function preloadImageUrl(cardImage: string) {
  getImageUrlForSignature(cardImage).catch((err) => {
    console.debug('[SignaturePicker] Preload failed:', err);
  });
}

/**
 * Handle SSE updates from signature store
 */
async function handleSseUpdate() {
  console.debug('[SignaturePicker] SSE update received');
  try {
    await loadSignaturesRealtime();
    // Preserve selected ID if it still exists
    if (selectedId.value && !localSignatures.value.find((s) => s.id === selectedId.value)) {
      selectedId.value = '';
    }
  } catch (err) {
    console.error('[SignaturePicker] Error handling SSE update:', err);
  }
}

// Lifecycle
onMounted(() => {
  // Register SSE callback for updates
  signatureStore.registerSseCallback(handleSseUpdate);
});

onUnmounted(() => {
  // Unregister SSE callback
  signatureStore.unregisterSseCallback();

  // Clean up image Blob URLs
  blobUrls.value.forEach((url) => {
    URL.revokeObjectURL(url);
  });
  blobUrls.value.clear();
  imageUrls.value.clear();
});

// Watch visible prop
watch(
  () => props.visible,
  async (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      searchQuery.value = '';
      selectedId.value = '';

      // Always load real data if no prop signatures provided to ensure status (isOriginalAuthor, isInAlbum) is up to date
      if (!props.signatures || props.signatures.length === 0) {
        await loadSignaturesRealtime();
      }
    }
  }
);

// Watch isVisible to sync with parent
watch(isVisible, (newVal) => {
  if (!newVal) {
    // 如果是导入授权操作，不触发 cancel 事件
    if (isImportingAuth.value) {
      isImportingAuth.value = false;
      return;
    }
    emit('cancel');
  }
});

// Handlers
const handleSignatureClick = (id: string) => {
  // 查找签名对象
  const sig = filteredSignatures.value.find((s) => s.id === id);
  if (!sig) return;

  // 如果签名被禁用（需要授权但未授权），不允许选择
  if (isSignatureDisabled(sig)) {
    return;
  }

  // Toggle: 若点击已选项则取消选择，否则选中
  if (selectedId.value === id) {
    selectedId.value = '';
  } else {
    selectedId.value = id;
  }
};

const onCreateNew = () => {
  emit('createNew');
};

const onImportAuth = () => {
  isImportingAuth.value = true; // 标记正在导入授权，避免触发 cancel
  emit('importAuth');
};

const onAuthRequest = () => {
  isImportingAuth.value = true; // 复用标记，避免触发 cancel
  emit('authRequest');
};

const onConfirmClick = async () => {
  if (!selectedId.value) return;

  // Check if signature is in album
  if (props.albumPath) {
    try {
      const result = await CheckSignatureInAlbum(props.albumPath, selectedId.value);
      if (result.isInAlbum) {
        // 智能检测：只有当内容有变更时才弹出确认框
        if (result.hasChanges) {
          updateConfirmDialogVisible.value = true;
          return;
        } else {
          // 内容无变更，直接选择"不更新"模式
          console.log('Signature content unchanged, skipping update confirmation');
          emit('select', selectedId.value, false);
          isVisible.value = false;
          return;
        }
      }
    } catch (err) {
      console.error('CheckSignatureInAlbum failed:', err);
    }
  }

  // Not in album or check failed, proceed with default (update=true)
  emit('select', selectedId.value, true);
  isVisible.value = false;
};

const handleUpdateConfirm = (update: boolean) => {
  if (!selectedId.value) return;
  emit('select', selectedId.value, update);
  isVisible.value = false;
  updateConfirmDialogVisible.value = false;
};

const onCancel = () => {
  isVisible.value = false;
};
</script>

<style scoped lang="scss">
.signature-picker-dialog {
  border-radius: 8px;

  // Sticky top header and search area
  .sticky-top {
    position: sticky;
    top: 0;
    left: 0;
    right: 0;
    z-index: 5;
    background: white;
  }

  :deep(.q-card__section) {
    padding: 16px;

    &:first-child {
      padding: 12px 16px;
    }
  }

  .signature-card {
    position: relative;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.2s ease;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);

    &:hover:not(.disabled-card) {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-1px);
    }

    &.selected {
      border: 2px solid var(--q-primary) !important;
      box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.1), 0 4px 16px rgba(33, 150, 243, 0.15);
    }

    // 禁用状态的签名卡片样式
    &.disabled-card {
      cursor: not-allowed;
      background-color: rgba(0, 0, 0, 0.02);
      border-color: rgba(0, 0, 0, 0.08);

      &:hover {
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
        transform: none;
      }
    }

    // Content area optimization
    .name-container,
    .intro-container {
      width: 100%;
      overflow: hidden;

      .scrollable-x {
        overflow-x: auto;
        overflow-y: hidden;
        white-space: nowrap;
        padding-right: 4px;

        // Custom scrollbar for horizontal scroll areas
        &::-webkit-scrollbar {
          height: 3px;
        }

        &::-webkit-scrollbar-track {
          background: transparent;
        }

        &::-webkit-scrollbar-thumb {
          background: rgba(0, 0, 0, 0.12);
          border-radius: 2px;

          &:hover {
            background: rgba(0, 0, 0, 0.2);
          }
        }
      }
    }

    .line-clamp-2 {
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    // Selection indicator with glow effect
    .selection-indicator-wrapper {
      position: relative;
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;

      .selection-glow {
        position: absolute;
        width: 100%;
        height: 100%;
        border-radius: 50%;
        background: radial-gradient(circle, rgba(33, 150, 243, 0.3), transparent);
        animation: glow-pulse 1.5s ease-in-out infinite;
      }

      .selection-icon {
        position: relative;
        z-index: 1;
      }
    }

    :deep(.q-card__section) {
      padding: 12px;

      .truncate {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .truncate-2 {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }

    &.create-new-card {
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 200px;

      &:hover {
        background-color: rgba(33, 150, 243, 0.02);
      }

      :deep(.q-card__section) {
        width: 100%;
      }
    }
  }

  // Sticky bottom action bar with frosted glass effect (统一为 4px 模糊，半透明 ~0.3)
  .sticky-bottom {
    position: sticky;
    bottom: 0;
    left: 0;
    right: 0;

    // Frosted glass effect: blur background + semi-transparent overlay
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
    background: rgba(255, 255, 255, 0.3);

    // Visual separation from content
    border-top: 1px solid rgba(0, 0, 0, 0.08);

    z-index: 10;
    transition: all 0.2s ease;
  }
}

// Animation for selection glow effect
@keyframes glow-pulse {
  0% {
    opacity: 0.8;
    transform: scale(1);
  }

  50% {
    opacity: 0.4;
  }

  100% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}
</style>
