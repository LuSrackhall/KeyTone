<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card
      class="signature-picker-dialog"
      style="width: 90%; max-width: 340px; max-height: 85vh; display: flex; flex-direction: column"
    >
      <!-- Header -->
      <q-card-section class="bg-primary text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.pickerDialog.title') }}</div>
      </q-card-section>

      <!-- Content (Scrollable) -->
      <q-card-section class="q-pa-md q-pt-lg col-grow overflow-auto">
        <!-- Description -->
        <div class="text-caption q-mb-md">
          {{ $t('exportFlow.pickerDialog.description') }}
        </div>

        <!-- Search Bar -->
        <div class="q-mb-md">
          <q-input
            v-model="searchQuery"
            filled
            dense
            :placeholder="$t('exportFlow.pickerDialog.search')"
            icon="search"
            clearable
            size="sm"
          />
          <div class="q-mt-xs flex justify-end">
            <q-btn
              size="xs"
              flat
              color="primary"
              icon="add"
              :label="$t('exportFlow.pickerDialog.createSignature')"
              @click="onCreateNew"
            />
          </div>
        </div>

        <!-- Signatures Grid / Empty State -->
        <div class="signatures-container">
          <!-- Empty State -->
          <div v-if="filteredSignatures.length === 0 && !searchQuery" class="text-center q-pa-md">
            <q-icon name="mail" size="32px" color="grey-5" />
            <div class="text-caption text-grey q-mt-sm">
              {{ $t('exportFlow.pickerDialog.emptyState') }}
            </div>
            <q-btn
              flat
              color="primary"
              :label="$t('exportFlow.pickerDialog.createSignature')"
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
              {{ $t('exportFlow.pickerDialog.noResults') }}
            </div>
          </div>

          <!-- Signatures List -->
          <div v-else class="col">
            <!-- Signature Cards -->
            <div v-for="sig in filteredSignatures" :key="sig.id" class="q-mb-sm">
              <q-card
                flat
                bordered
                :class="['signature-card cursor-pointer', selectedId === sig.id ? 'selected' : '']"
                @click="selectSignature(sig.id)"
                style="display: flex; align-items: center; min-height: 70px"
              >
                <!-- Left: Image Area (fixed 60px) -->
                <div
                  class="flex-shrink-0 flex items-center justify-center"
                  style="width: 60px; height: 60px; background-color: #f5f5f5; border-radius: 4px; margin: 0 8px"
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
                <div class="col flex flex-col justify-center" style="padding: 0 8px; min-width: 0">
                  <div class="text-caption text-weight-bold truncate" style="font-size: 0.9rem">
                    {{ sig.name }}
                  </div>
                  <div
                    class="text-caption text-grey"
                    style="
                      font-size: 0.75rem;
                      overflow: hidden;
                      text-overflow: ellipsis;
                      display: -webkit-box;
                      -webkit-line-clamp: 2;
                      line-clamp: 2;
                      -webkit-box-orient: vertical;
                    "
                  >
                    {{ sig.intro || $t('exportFlow.pickerDialog.noIntro') }}
                  </div>
                </div>

                <!-- Right: Selection Indicator -->
                <div v-if="selectedId === sig.id" class="flex-shrink-0" style="margin-left: 8px; margin-right: 8px">
                  <q-icon name="check_circle" size="20px" color="positive" />
                </div>
              </q-card>
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat :label="$t('exportFlow.pickerDialog.cancel')" color="primary" size="sm" @click="onCancel" />
        <q-btn
          unelevated
          :label="$t('exportFlow.pickerDialog.confirm')"
          color="primary"
          size="sm"
          :disable="!selectedId"
          @click="onConfirm"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useSignatureStore } from 'src/stores/signature-store';
import { getSignaturesList, decryptSignatureData, getSignatureImage } from 'boot/query/signature-query';

interface Signature {
  id: string;
  name: string;
  intro?: string;
  image?: string;
}

interface SignaturePickerDialogProps {
  visible: boolean;
  signatures?: Signature[];
}

interface SignaturePickerDialogEmits {
  (e: 'select', signatureId: string): void;
  (e: 'createNew'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<SignaturePickerDialogProps>(), {
  visible: false,
  signatures: () => [],
});

const emit = defineEmits<SignaturePickerDialogEmits>();
useI18n(); // Initialize i18n, $t available in template
const signatureStore = useSignatureStore();

// UI State
const isVisible = ref(false);
const searchQuery = ref('');
const selectedId = ref('');
const loading = ref(false);
const localSignatures = ref<Signature[]>([]);
const imageUrls = ref<Map<string, string>>(new Map()); // cardImage path -> blob URL
const blobUrls = ref<Set<string>>(new Set()); // Track blob URLs for cleanup

// Determine which signatures to use (real data or fallback to props)
const effectiveSignatures = computed(() => {
  if (localSignatures.value.length > 0) {
    return localSignatures.value;
  }
  return props.signatures || [];
});

// Filter signatures based on search query (only by name)
const filteredSignatures = computed(() => {
  if (!searchQuery.value) {
    return effectiveSignatures.value;
  }

  const query = searchQuery.value.toLowerCase();
  return effectiveSignatures.value.filter((sig) => sig.name.toLowerCase().includes(query));
});

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

      // Load real data if not already loaded and no prop signatures provided
      if (localSignatures.value.length === 0 && (!props.signatures || props.signatures.length === 0)) {
        await loadSignaturesRealtime();
      }
    }
  }
);

// Watch isVisible to sync with parent
watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('cancel');
  }
});

// Handlers
const selectSignature = (id: string) => {
  selectedId.value = id;
};

const onCreateNew = () => {
  emit('createNew');
};

const onConfirm = () => {
  if (!selectedId.value) return;
  emit('select', selectedId.value);
  isVisible.value = false;
};

const onCancel = () => {
  isVisible.value = false;
};
</script>

<style scoped lang="scss">
.signature-picker-dialog {
  border-radius: 8px;

  :deep(.q-card__section) {
    padding: 16px;

    &:first-child {
      padding: 12px 16px;
    }
  }

  .signatures-container {
    // Scrollbar styling
    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.2);
      border-radius: 3px;

      &:hover {
        background: rgba(0, 0, 0, 0.3);
      }
    }
  }

  .signature-card {
    position: relative;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      transform: translateY(-2px);
    }

    &.selected {
      border-color: var(--q-primary) !important;
      border-width: 2px;
      box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.1);
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
}
</style>
