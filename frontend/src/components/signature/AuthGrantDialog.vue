<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card class="auth-grant-dialog" style="width: 90%; max-width: 400px">
      <!-- Header -->
      <q-card-section class="bg-teal bg-opacity-90 text-white q-pa-sm sticky top-0 z-10 backdrop-blur-sm">
        <div class="text-subtitle1">{{ t('signature.authGrant.title') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-md">
        <!-- Step 1: Import auth request file -->
        <div v-if="currentStep === 1">
          <div class="text-body2 q-mb-md">
            {{ t('signature.authGrant.step1Description') }}
          </div>

          <!-- File Import -->
          <q-card flat bordered :class="fileImported ? 'bg-positive-1' : 'bg-grey-1'">
            <q-card-section class="q-pa-sm">
              <div class="row items-center justify-between q-gutter-xs">
                <div class="col-grow">
                  <div class="text-caption">
                    {{ fileImported ? t('signature.authGrant.fileImported') : t('signature.authGrant.selectFile') }}
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
                  accept=".ktauthreq"
                  class="hidden-file-input"
                  @change="onFileInputChange"
                  style="display: none"
                />
                <q-btn flat icon="attach_file" color="primary" size="xs" @click="triggerFileInput" :loading="parsing">
                  {{ t('signature.authGrant.importButton') }}
                </q-btn>
              </div>
            </q-card-section>
          </q-card>

          <!-- Error message -->
          <div v-if="errorMessage" class="text-caption text-negative q-mt-sm">
            {{ errorMessage }}
          </div>
        </div>

        <!-- Step 2: Review and approve -->
        <div v-if="currentStep === 2 && parsedRequest">
          <div class="text-body2 q-mb-md">
            {{ t('signature.authGrant.step2Description') }}
          </div>

          <!-- Album UUID Hash -->
          <div class="q-mb-md">
            <div class="text-overline text-grey" style="font-size: 0.7rem">
              {{ t('signature.authGrant.albumId') }}
            </div>
            <q-card flat bordered class="bg-grey-1">
              <q-card-section class="q-pa-sm">
                <div class="text-caption break-words" style="word-break: break-all; font-family: monospace">
                  {{ parsedRequest.authorizationUUIDHash }}
                </div>
              </q-card-section>
            </q-card>
          </div>

          <!-- Requester Info -->
          <div class="q-mb-md">
            <div class="text-overline text-grey" style="font-size: 0.7rem">
              {{ t('signature.authGrant.requester') }}
            </div>
            <q-card flat bordered class="bg-grey-1">
              <q-card-section class="q-pa-sm">
                <div
                  class="text-body2 text-weight-medium"
                  :class="[
                    /* 对溢出的情况, 采取滚动策略（与签名列表保持一致） */
                    'max-w-full !overflow-x-auto whitespace-nowrap !text-clip',
                    // 添加细微滚动条（与签名列表保持一致）
                    'h-5.5 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                  ]"
                  style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                >
                  {{ parsedRequest.requesterSignatureName }}
                </div>
              </q-card-section>
              <!-- 申请方资格码指纹展示 -->
              <q-separator v-if="parsedRequest.requesterQualificationFingerprint" />
              <q-item v-if="parsedRequest.requesterQualificationFingerprint" dense class="q-pa-xs">
                <q-item-section avatar style="min-width: 24px">
                  <q-icon name="fingerprint" color="grey-6" size="16px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label caption style="font-size: 0.65rem">
                    {{ t('signature.authGrant.requesterFingerprint') }}
                  </q-item-label>
                  <q-item-label
                    class="text-caption"
                    style="font-family: monospace; word-break: break-all; line-height: 1.3; font-size: 0.65rem"
                  >
                    {{ parsedRequest.requesterQualificationFingerprint }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn
                    flat
                    dense
                    size="xs"
                    icon="content_copy"
                    @click="copyToClipboard(parsedRequest.requesterQualificationFingerprint, t('signature.authGrant.requesterFingerprint'))"
                  />
                </q-item-section>
              </q-item>
            </q-card>
          </div>

          <!-- Matched Signatures (original author's signatures that match) -->
          <div class="q-mb-md">
            <div class="text-overline text-grey q-mb-xs" style="font-size: 0.7rem">
              {{ t('signature.authGrant.yourSignature') }}
            </div>
            <div v-if="matchedSignatures.length === 0" class="text-caption text-negative">
              {{ t('signature.authGrant.noMatchingSignature') }}
            </div>
            <template v-else>
              <q-select
                v-if="matchedSignatures.length > 1"
                v-model="selectedMatchedSignatureId"
                :options="matchedSignatureOptions"
                emit-value
                map-options
                dense
                outlined
                :label="t('signature.authGrant.selectSignature')"
              />

              <!-- 选中的签名详情：展示完整签名列表项（图片+名称+介绍），若本地没有则回退名称 -->
              <q-card
                v-if="matchedSignatures.length === 1 || selectedMatchedSignatureId"
                flat
                bordered
                class="bg-grey-1 q-mt-sm"
                style="overflow: hidden"
              >
                <q-item dense style="overflow: hidden">
                  <q-item-section avatar v-if="selectedLocalSignature?.cardImage">
                    <q-avatar size="36px">
                      <q-img :src="props.getImageUrl(String(selectedLocalSignature.cardImage))" />
                    </q-avatar>
                  </q-item-section>
                  <q-item-section style="overflow: hidden; min-width: 0">
                    <q-item-label
                      :class="[
                        /* 对溢出的情况, 采取滚动策略（与签名列表保持一致） */
                        'max-w-full !overflow-x-auto whitespace-nowrap !text-clip',
                        // 添加细微滚动条（与签名列表保持一致）
                        'h-5.5 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                      ]"
                      style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                    >
                      {{ selectedLocalSignature?.name || selectedMatchedSignatureFallback?.name || '' }}
                    </q-item-label>
                    <q-item-label
                      caption
                      v-if="selectedLocalSignature?.intro"
                      :class="[
                        /* 对溢出的情况, 采取滚动策略（与签名列表保持一致） */
                        'max-w-full !overflow-x-auto whitespace-nowrap !mt-1.5',
                        // 添加细微滚动条（与签名列表保持一致）
                        'h-4.4 [&::-webkit-scrollbar]:h-0.3 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                      ]"
                      style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                    >
                      {{ selectedLocalSignature.intro }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
                <!-- 资格码指纹展示 -->
                <q-separator />
                <q-item dense class="q-pa-xs">
                  <q-item-section avatar style="min-width: 24px">
                    <q-icon name="fingerprint" color="grey-6" size="16px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption style="font-size: 0.65rem">
                      {{ t('signature.authGrant.qualificationFingerprint') }}
                    </q-item-label>
                    <q-item-label
                      class="text-caption"
                      style="font-family: monospace; word-break: break-all; line-height: 1.3; font-size: 0.65rem"
                    >
                      {{ selectedMatchedSignatureFingerprint }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="copyToClipboard(selectedMatchedSignatureFingerprint, t('signature.authGrant.qualificationFingerprint'))"
                    />
                  </q-item-section>
                </q-item>
              </q-card>
            </template>
          </div>

          <!-- Warning -->
          <q-banner class="bg-warning-1 q-mb-sm" rounded dense>
            <template #avatar>
              <q-icon name="warning" color="warning" />
            </template>
            <div class="text-caption">{{ t('signature.authGrant.warning') }}</div>
          </q-banner>
        </div>

        <!-- Step 3: Success -->
        <div v-if="currentStep === 3">
          <div class="text-center q-pa-md">
            <q-icon name="check_circle" size="64px" color="positive" />
            <div class="text-h6 q-mt-md">{{ t('signature.authGrant.exportSuccess') }}</div>
            <div class="text-caption text-grey q-mt-sm">
              {{ t('signature.authGrant.exportSuccessHint') }}
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="t('signature.authGrant.cancel')" color="grey" size="sm" @click="onCancel" />

        <!-- Step 1 actions -->
        <q-btn
          v-if="currentStep === 1"
          unelevated
          :label="t('signature.authGrant.next')"
          color="teal"
          size="sm"
          :disable="!fileImported"
          @click="goToStep2"
        />

        <!-- Step 2 actions -->
        <q-btn
          v-if="currentStep === 2"
          flat
          :label="t('signature.authGrant.back')"
          color="primary"
          size="sm"
          @click="goToStep1"
        />
        <q-btn
          v-if="currentStep === 2"
          unelevated
          :label="t('signature.authGrant.approve')"
          color="teal"
          size="sm"
          :disable="matchedSignatures.length === 0 || (matchedSignatures.length > 1 && !selectedMatchedSignatureId)"
          :loading="exporting"
          @click="approveAndExport"
        />

        <!-- Step 3 actions -->
        <q-btn
          v-if="currentStep === 3"
          unelevated
          :label="t('signature.authGrant.done')"
          color="positive"
          size="sm"
          @click="onDone"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import type { Signature } from 'src/types/signature';
import {
  parseAuthRequest,
  verifyAuthRequestOwner,
  generateAuthGrant,
  type ParsedAuthRequest,
  type MatchedSignature,
} from 'boot/query/signature-query';

interface AuthGrantDialogProps {
  visible: boolean;

  // 签名管理页本地已解密的签名列表（用于展示完整签名卡片）
  localSignatures?: Signature[];

  // 解析签名图片路径 -> 预览 URL（由父组件复用既有缓存逻辑提供）
  getImageUrl?: (imagePath: string) => string;
}

interface AuthGrantDialogEmits {
  (e: 'cancel'): void;
  (e: 'done'): void;
}

const props = withDefaults(defineProps<AuthGrantDialogProps>(), {
  visible: false,
  localSignatures: () => [],
  getImageUrl: () => '',
});

const emit = defineEmits<AuthGrantDialogEmits>();

const { t } = useI18n();
const $q = useQuasar();

// Refs
const hiddenFileInput = ref<HTMLInputElement | null>(null);

// UI State
const isVisible = ref(false);
const currentStep = ref(1);
const fileImported = ref(false);
const selectedFileName = ref('');
const errorMessage = ref('');
const parsing = ref(false);
const exporting = ref(false);

// Data State
const parsedRequest = ref<ParsedAuthRequest | null>(null);
const matchedSignatures = ref<MatchedSignature[]>([]);
const selectedMatchedSignatureId = ref('');
const fileContent = ref<Uint8Array | null>(null);

// Computed
const matchedSignatureOptions = computed(() => {
  return matchedSignatures.value.map((sig) => ({
    label: sig.name,
    value: sig.encryptedId,
  }));
});

const selectedDisplayEncryptedId = computed(() => {
  if (matchedSignatures.value.length === 1) {
    return matchedSignatures.value[0].encryptedId;
  }
  return selectedMatchedSignatureId.value;
});

const selectedLocalSignature = computed(() => {
  const encryptedId = selectedDisplayEncryptedId.value;
  if (!encryptedId) return null;
  return props.localSignatures.find((s) => s.id === encryptedId) || null;
});

const selectedMatchedSignatureFallback = computed(() => {
  const encryptedId = selectedDisplayEncryptedId.value;
  if (!encryptedId) return null;
  return matchedSignatures.value.find((s) => s.encryptedId === encryptedId) || null;
});

/**
 * 选中签名的资格码指纹
 */
const selectedMatchedSignatureFingerprint = computed(() => {
  const matched = selectedMatchedSignatureFallback.value;
  return matched?.qualificationFingerprint || '';
});

/**
 * 复制到剪贴板
 */
function copyToClipboard(text: string, label: string) {
  navigator.clipboard.writeText(text).then(
    () => {
      $q.notify({
        type: 'positive',
        message: t('signature.authGrant.copySuccess', { label }),
        position: 'top',
        timeout: 1500,
      });
    },
    () => {
      $q.notify({
        type: 'negative',
        message: t('signature.authGrant.copyFailed'),
        position: 'top',
        timeout: 1500,
      });
    }
  );
}

// Watch
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      // Reset state when dialog opens
      resetState();
    }
  }
);

watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('cancel');
  }
});

// Methods
function resetState() {
  currentStep.value = 1;
  fileImported.value = false;
  selectedFileName.value = '';
  errorMessage.value = '';
  parsing.value = false;
  exporting.value = false;
  parsedRequest.value = null;
  matchedSignatures.value = [];
  selectedMatchedSignatureId.value = '';
  fileContent.value = null;
}

function triggerFileInput() {
  hiddenFileInput.value?.click();
}

async function onFileInputChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];

  if (!file) return;

  selectedFileName.value = file.name;
  errorMessage.value = '';
  parsing.value = true;

  try {
    // Read file content
    const arrayBuffer = await file.arrayBuffer();
    fileContent.value = new Uint8Array(arrayBuffer);

    // Parse the auth request file
    const parsed = await parseAuthRequest(fileContent.value);
    if (!parsed) {
      throw new Error('Failed to parse auth request file');
    }

    parsedRequest.value = parsed;

    // Verify ownership
    const verifyResult = await verifyAuthRequestOwner(parsed.originalAuthorQualCodeHash);
    if (!verifyResult) {
      throw new Error('Failed to verify auth request');
    }

    if (!verifyResult.hasPermission) {
      errorMessage.value = t('signature.authGrant.noPermission');
      fileImported.value = false;
      return;
    }

    matchedSignatures.value = verifyResult.matchedSignatures;
    if (matchedSignatures.value.length === 1) {
      selectedMatchedSignatureId.value = matchedSignatures.value[0].encryptedId;
    }

    fileImported.value = true;
  } catch (err) {
    console.error('Parse auth request failed:', err);
    errorMessage.value = t('signature.authGrant.parseError');
    fileImported.value = false;
  } finally {
    parsing.value = false;
    // Reset input
    target.value = '';
  }
}

function goToStep1() {
  currentStep.value = 1;
}

function goToStep2() {
  if (fileImported.value && parsedRequest.value) {
    currentStep.value = 2;
  }
}

async function approveAndExport() {
  if (!parsedRequest.value || !selectedMatchedSignatureId.value) {
    return;
  }

  exporting.value = true;

  try {
    // Generate auth grant file
    const grantFileContent = await generateAuthGrant(
      parsedRequest.value.authorizationUUIDHash,
      parsedRequest.value.requesterSignatureIDSuffix,
      selectedMatchedSignatureId.value
    );

    if (!grantFileContent) {
      throw new Error('Failed to generate auth grant file');
    }

    const blob = new Blob([new Uint8Array(grantFileContent)], { type: 'application/octet-stream' });

    // 使用 File System Access API 提供保存对话框，确保用户选择路径后再提示成功
    if (typeof window.showSaveFilePicker === 'function') {
      try {
        const handle = await window.showSaveFilePicker({
          suggestedName: `auth-grant-${Date.now()}.ktauth`,
          types: [
            {
              description: 'Authorization Grant (.ktauth)',
              accept: { 'application/octet-stream': ['.ktauth'] },
            },
          ],
        });

        const writable = await handle.createWritable();
        await writable.write(blob);
        await writable.close();
      } catch (err) {
        if (err instanceof Error && err.name === 'AbortError') {
          // 用户取消保存，不提示成功
          return;
        }
        throw err;
      }
    } else {
      // 兼容性回退方案
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `auth-grant-${Date.now()}.ktauth`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    }

    // Go to success step
    currentStep.value = 3;

    $q.notify({
      type: 'positive',
      message: t('signature.authGrant.exportSuccess'),
      position: 'top',
    });
  } catch (err) {
    console.error('Export auth grant failed:', err);
    $q.notify({
      type: 'negative',
      message: t('signature.authGrant.exportFailed'),
      position: 'top',
    });
  } finally {
    exporting.value = false;
  }
}

function onCancel() {
  isVisible.value = false;
}

function onDone() {
  isVisible.value = false;
  emit('done');
}
</script>

<style scoped lang="scss">
.auth-grant-dialog {
  border-radius: 8px;

  .truncate {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .break-words {
    word-break: break-all;
  }
}
</style>
