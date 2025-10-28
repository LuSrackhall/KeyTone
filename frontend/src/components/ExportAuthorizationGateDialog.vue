<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card class="auth-gate-dialog" style="width: 90%; max-width: 340px">
      <!-- Header -->
      <q-card-section class="bg-warning text-white q-pa-sm">
        <div class="text-subtitle1">{{ $t('exportFlow.authGateDialog.title') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-md q-pt-lg">
        <!-- Description -->
        <div class="text-caption q-mb-md text-center">
          {{ $t('exportFlow.authGateDialog.description') }}
        </div>

        <!-- Author Contact Section -->
        <div class="q-mb-md">
          <div class="text-overline text-grey q-mb-xs" style="font-size: 0.7rem">
            {{ $t('exportFlow.authGateDialog.authorContact') }}
          </div>
          <q-card flat bordered class="bg-grey-1">
            <q-card-section class="q-pa-sm">
              <div class="row items-center justify-between q-gutter-xs">
                <div class="col-grow text-caption break-words" style="word-break: break-all; line-height: 1.3">
                  {{ authorContact }}
                </div>
                <q-btn
                  flat
                  dense
                  size="xs"
                  icon="content_copy"
                  :title="$t('exportFlow.authGateDialog.copy')"
                  @click="copyContact"
                />
              </div>
            </q-card-section>
          </q-card>
          <div class="text-caption text-grey q-mt-xs">
            {{ $t('exportFlow.authGateDialog.contactHint') }}
          </div>
        </div>

        <!-- Authorization File Import Section -->
        <div class="q-mb-md">
          <div class="text-overline text-grey q-mb-xs" style="font-size: 0.7rem">
            {{ $t('exportFlow.authGateDialog.importAuthFile') }}
          </div>
          <div class="text-caption text-grey q-mb-sm">
            {{ $t('exportFlow.authGateDialog.importHint') }}
          </div>

          <!-- File Selection Card -->
          <q-card flat bordered :class="fileSelected ? 'bg-positive-1' : 'bg-grey-1'">
            <q-card-section class="q-pa-sm">
              <div class="row items-center justify-between q-gutter-xs">
                <div class="col-grow">
                  <div class="text-caption">
                    {{
                      fileSelected
                        ? $t('exportFlow.authGateDialog.fileSelected')
                        : $t('exportFlow.authGateDialog.selectFile')
                    }}
                  </div>
                  <div
                    v-if="selectedFileName"
                    class="text-caption text-grey q-mt-xs truncate"
                    style="font-size: 0.8rem"
                  >
                    {{ selectedFileName }}
                  </div>
                </div>
                <input
                  ref="hiddenFileInput"
                  type="file"
                  accept=".json,.ktauth"
                  class="hidden-file-input"
                  @change="onFileInputChange"
                  style="display: none"
                />
                <q-btn flat icon="attach_file" color="primary" size="xs" @click="triggerFileInput">
                  {{ $t('exportFlow.authGateDialog.importButton') }}
                </q-btn>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Help Section -->
        <div class="text-center q-mt-sm">
          <q-btn flat size="xs" color="primary" :label="$t('exportFlow.authGateDialog.help')" icon="help_outline" />
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat :label="$t('exportFlow.authGateDialog.cancel')" color="primary" size="sm" @click="onCancel" />
        <q-btn
          unelevated
          :label="$t('exportFlow.authGateDialog.continue')"
          color="primary"
          size="sm"
          :disable="!fileSelected"
          @click="onAuthorized"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';

interface AuthGateDialogProps {
  visible: boolean;
  authorContact?: string;
}

interface AuthGateDialogEmits {
  (e: 'authorized'): void;
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<AuthGateDialogProps>(), {
  visible: false,
  authorContact: '',
});

const emit = defineEmits<AuthGateDialogEmits>();
const { t } = useI18n();
const { notify } = useQuasar();

const isVisible = ref(false);
const authFile = ref<File | null>(null);
const selectedFileName = ref('');
const fileInput = ref();

const authorContact = computed(() => props.authorContact || t('exportFlow.authGateDialog.authorContact'));
const fileSelected = computed(() => !!authFile.value);

// Watch visible prop
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      authFile.value = null;
      selectedFileName.value = '';
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
const hiddenFileInput = ref<HTMLInputElement>();

const triggerFileInput = () => {
  hiddenFileInput.value?.click();
};

const onFileInputChange = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    authFile.value = file;
    selectedFileName.value = file.name;
  }
};

const onFileSelected = (file: File | null) => {
  if (file) {
    authFile.value = file;
    selectedFileName.value = file.name;
  } else {
    authFile.value = null;
    selectedFileName.value = '';
  }
};

const copyContact = () => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(authorContact.value).then(() => {
      notify({
        type: 'positive',
        message: t('exportFlow.notify.contactCopied'),
        position: 'top',
      });
    });
  }
};

const onAuthorized = () => {
  if (!fileSelected.value) return;

  // Simulate file import success (placeholder for real logic)
  notify({
    type: 'positive',
    message: t('exportFlow.notify.authFileImported'),
    position: 'top',
  });

  emit('authorized');
  isVisible.value = false;
};

const onCancel = () => {
  isVisible.value = false;
};
</script>

<style scoped lang="scss">
.auth-gate-dialog {
  border-radius: 8px;

  :deep(.q-card__section) {
    padding: 16px;

    &:first-child {
      padding: 12px 16px;
    }
  }

  .hidden-file-input {
    display: none;
  }

  .truncate {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .break-words {
    word-break: break-word;
    overflow-wrap: break-word;
  }
}
</style>
