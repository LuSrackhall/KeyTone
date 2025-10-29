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
              >
                <!-- Image -->
                <div class="signature-image" style="height: 80px; overflow: hidden">
                  <img
                    v-if="sig.image"
                    :src="sig.image"
                    :alt="sig.name"
                    class="full-width full-height"
                    style="object-fit: cover"
                  />
                  <div v-else class="full-width full-height bg-grey-2 flex flex-center">
                    <q-icon name="image_not_supported" size="24px" color="grey-5" />
                  </div>
                </div>

                <!-- Info -->
                <q-card-section class="q-pa-xs">
                  <div class="text-caption text-weight-bold truncate">
                    {{ sig.name }}
                  </div>
                  <div class="text-caption text-grey truncate-2 q-mt-xs" style="font-size: 0.75rem">
                    {{ sig.intro || $t('exportFlow.pickerDialog.noIntro') }}
                  </div>
                </q-card-section>

                <!-- Selection Indicator -->
                <div v-if="selectedId === sig.id" class="selection-indicator">
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
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';

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
const { t } = useI18n();

const isVisible = ref(false);
const searchQuery = ref('');
const selectedId = ref('');

// Filter signatures based on search query (only by name)
const filteredSignatures = computed(() => {
  if (!searchQuery.value) {
    return props.signatures;
  }

  const query = searchQuery.value.toLowerCase();
  return props.signatures.filter((sig) => sig.name.toLowerCase().includes(query));
});

// Watch visible prop
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      searchQuery.value = '';
      selectedId.value = '';
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

    .signature-image {
      width: 100%;
      height: 120px;
      overflow: hidden;
      background-color: #f5f5f5;
    }

    .selection-indicator {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: rgba(255, 255, 255, 0.95);
      border-radius: 50%;
      width: 64px;
      height: 64px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
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
