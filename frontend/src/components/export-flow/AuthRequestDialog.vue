<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card class="auth-request-dialog" style="width: 90%; max-width: 400px">
      <!-- Header -->
      <q-card-section class="bg-deep-purple bg-opacity-90 text-white q-pa-sm sticky top-0 z-10 backdrop-blur-sm">
        <div class="text-subtitle1">{{ t('exportFlow.authRequestDialog.title') }}</div>
      </q-card-section>

      <!-- Content -->
      <q-card-section class="q-pa-md">
        <!-- Step Indicator - 优化版 -->
        <div class="q-mb-md">
          <div class="stepper-container">
            <!-- Step 1 -->
            <div class="step-item" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">
              <div class="step-circle">
                <q-icon v-if="currentStep > 1" name="check" size="14px" />
                <span v-else>1</span>
              </div>
              <div class="step-label">{{ t('exportFlow.authRequestDialog.step1Title') }}</div>
            </div>
            <div class="step-connector" :class="{ active: currentStep > 1 }"></div>
            <!-- Step 2 -->
            <div class="step-item" :class="{ active: currentStep >= 2, completed: currentStep > 2 }">
              <div class="step-circle">
                <q-icon v-if="currentStep > 2" name="check" size="14px" />
                <span v-else>2</span>
              </div>
              <div class="step-label">{{ t('exportFlow.authRequestDialog.step2Title') }}</div>
            </div>
            <div class="step-connector" :class="{ active: currentStep > 2 }"></div>
            <!-- Step 3 -->
            <div class="step-item" :class="{ active: currentStep >= 3 }">
              <div class="step-circle">
                <q-icon v-if="currentStep >= 3" name="check" size="14px" />
                <span v-else>3</span>
              </div>
              <div class="step-label">{{ t('exportFlow.authRequestDialog.step3Title') }}</div>
            </div>
          </div>
        </div>

        <!-- Step 1: Select Unauthorized Signature -->
        <div v-if="currentStep === 1">
          <div class="text-body2 q-mb-sm">
            {{ t('exportFlow.authRequestDialog.step1Description') }}
          </div>

          <!-- No signatures available -->
          <div v-if="unauthorizedSignatures.length === 0" class="text-center q-pa-md">
            <q-icon name="info" size="32px" color="grey-5" />
            <div class="text-caption text-grey q-mt-sm">
              {{
                hasAnySignatures
                  ? t('exportFlow.authRequestDialog.noUnauthorizedSignatures')
                  : t('exportFlow.authRequestDialog.noSignatures')
              }}
            </div>
            <q-btn
              v-if="!hasAnySignatures"
              flat
              color="primary"
              :label="t('exportFlow.authRequestDialog.createSignature')"
              icon="add"
              size="sm"
              class="q-mt-md"
              @click="onCreateSignature"
            />
          </div>

          <!-- Signature list -->
          <div v-else class="signature-list q-mb-sm" style="max-height: 200px; overflow-y: auto">
            <q-card
              v-for="sig in unauthorizedSignatures"
              :key="sig.id"
              flat
              bordered
              :class="['signature-item q-mb-xs cursor-pointer', selectedSignatureId === sig.id ? 'selected' : '']"
              @click="selectSignature(sig.id)"
            >
              <q-card-section class="q-pa-sm flex items-center">
                <div class="flex-shrink-0 q-mr-sm" style="width: 40px; height: 40px">
                  <img
                    v-if="sig.image"
                    :src="sig.image"
                    :alt="sig.name"
                    style="width: 100%; height: 100%; object-fit: cover; border-radius: 4px"
                  />
                  <q-icon v-else name="person" size="40px" color="grey-4" />
                </div>
                <div class="col" style="overflow: hidden; min-width: 0">
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
                    {{ sig.name }}
                  </div>
                  <div
                    class="text-caption text-grey"
                    :class="[
                      /* 对溢出的情况, 采取滚动策略（与签名列表保持一致） */
                      'max-w-full !overflow-x-auto whitespace-nowrap !mt-1.5',
                      // 添加细微滚动条（与签名列表保持一致）
                      'h-4.4 [&::-webkit-scrollbar]:h-0.3 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                    ]"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ sig.intro || t('exportFlow.authRequestDialog.noIntro') }}
                  </div>
                </div>
                <q-icon v-if="selectedSignatureId === sig.id" name="check_circle" color="positive" size="20px" />
              </q-card-section>
            </q-card>
          </div>
        </div>

        <!-- Step 2: Show Contact Info & Export Instructions -->
        <div v-if="currentStep === 2">
          <div class="text-body2 q-mb-sm">
            {{ t('exportFlow.authRequestDialog.step2Description') }}
          </div>

          <!-- Selected Signature Card - 完整展示 -->
          <div class="q-mb-md">
            <div class="text-overline text-grey" style="font-size: 0.7rem">
              {{ t('exportFlow.authRequestDialog.selectedSignature') }}
            </div>
            <q-card flat bordered class="q-mt-xs">
              <q-card-section class="q-pa-sm flex items-center">
                <div class="flex-shrink-0 q-mr-sm" style="width: 48px; height: 48px">
                  <img
                    v-if="selectedSignature?.image"
                    :src="selectedSignature.image"
                    :alt="selectedSignature.name"
                    style="width: 100%; height: 100%; object-fit: cover; border-radius: 6px"
                  />
                  <div
                    v-else
                    class="flex items-center justify-center bg-grey-2"
                    style="width: 100%; height: 100%; border-radius: 6px"
                  >
                    <q-icon name="person" size="24px" color="grey-5" />
                  </div>
                </div>
                <div class="col" style="overflow: hidden; min-width: 0">
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
                    {{ selectedSignatureName }}
                  </div>
                  <div
                    class="text-caption text-grey"
                    :class="[
                      /* 对溢出的情况, 采取滚动策略（与签名列表保持一致） */
                      'max-w-full !overflow-x-auto whitespace-nowrap !mt-1.5',
                      // 添加细微滚动条（与签名列表保持一致）
                      'h-4.4 [&::-webkit-scrollbar]:h-0.3 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
                    ]"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ selectedSignature?.intro || t('exportFlow.authRequestDialog.noIntro') }}
                  </div>
                </div>
              </q-card-section>
              <!-- 资格码指纹展示 -->
              <q-separator v-if="selectedSignature?.qualificationFingerprint" />
              <q-item v-if="selectedSignature?.qualificationFingerprint" dense class="q-pa-xs">
                <q-item-section avatar style="min-width: 24px">
                  <q-icon name="fingerprint" color="grey-6" size="16px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label caption style="font-size: 0.65rem">
                    {{ t('exportFlow.authRequestDialog.qualificationFingerprint') }}
                  </q-item-label>
                  <q-item-label
                    class="text-caption"
                    style="font-family: monospace; word-break: break-all; line-height: 1.3; font-size: 0.65rem"
                  >
                    {{ selectedSignature.qualificationFingerprint }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn
                    flat
                    dense
                    size="xs"
                    icon="content_copy"
                    @click="copyText(selectedSignature.qualificationFingerprint, t('exportFlow.authRequestDialog.qualificationFingerprint'))"
                  />
                </q-item-section>
              </q-item>
            </q-card>
          </div>

          <!-- Author Contact Section - 分开展示邮箱和备用联系方式 -->
          <div class="q-mb-md">
            <div class="text-overline text-grey q-mb-xs" style="font-size: 0.7rem">
              {{ t('exportFlow.authRequestDialog.authorContact') }}
            </div>
            <q-card flat bordered class="bg-grey-1">
              <q-card-section class="q-pa-sm">
                <!-- 邮箱 -->
                <div class="row items-center q-mb-xs">
                  <q-icon name="email" size="16px" color="grey-7" class="q-mr-xs" />
                  <span class="text-caption text-grey-7">{{ t('exportFlow.authRequestDialog.email') }}:</span>
                </div>
                <div class="row items-center justify-between q-mb-sm">
                  <div class="col-grow text-body2" style="word-break: break-all">
                    {{ contactEmail || t('exportFlow.authRequestDialog.noContact') }}
                  </div>
                  <q-btn
                    v-if="contactEmail"
                    flat
                    dense
                    size="xs"
                    icon="content_copy"
                    :title="t('exportFlow.authRequestDialog.copy')"
                    @click="copyText(contactEmail)"
                  />
                </div>
                <!-- 备用联系方式 -->
                <template v-if="contactAdditional">
                  <q-separator class="q-my-xs" />
                  <div class="row items-center q-mb-xs q-mt-sm">
                    <q-icon name="contact_phone" size="16px" color="grey-7" class="q-mr-xs" />
                    <span class="text-caption text-grey-7"
                      >{{ t('exportFlow.authRequestDialog.additionalContact') }}:</span
                    >
                  </div>
                  <div class="row items-center justify-between">
                    <div
                      class="col-grow text-body2"
                      style="word-break: break-all; line-height: 1.4; white-space: pre-line"
                    >
                      {{ contactAdditional }}
                    </div>
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      :title="t('exportFlow.authRequestDialog.copy')"
                      @click="copyText(contactAdditional)"
                    />
                  </div>
                </template>
              </q-card-section>
            </q-card>
          </div>

          <!-- Instructions - 优化说明 -->
          <div class="q-mb-sm">
            <div class="text-overline text-grey q-mb-xs" style="font-size: 0.7rem">
              {{ t('exportFlow.authRequestDialog.instructions') }}
            </div>
            <!-- 重要提示 -->
            <q-banner dense rounded class="bg-amber-1 text-amber-10 q-mb-sm" style="font-size: 0.75rem">
              <template v-slot:avatar>
                <q-icon name="tips_and_updates" color="amber-8" size="18px" />
              </template>
              {{ t('exportFlow.authRequestDialog.instructionTip') }}
            </q-banner>
            <div class="text-caption text-grey-8" style="line-height: 1.5">
              <p class="q-mb-xs">1. {{ t('exportFlow.authRequestDialog.instruction1') }}</p>
              <p class="q-mb-xs">2. {{ t('exportFlow.authRequestDialog.instruction2') }}</p>
              <p class="q-mb-xs">3. {{ t('exportFlow.authRequestDialog.instruction3') }}</p>
            </div>
          </div>
        </div>

        <!-- Step 3: Export Complete -->
        <div v-if="currentStep === 3">
          <div class="text-center q-pa-md">
            <q-icon name="check_circle" size="64px" color="positive" />
            <div class="text-h6 q-mt-md">{{ t('exportFlow.authRequestDialog.exportSuccess') }}</div>
            <div class="text-caption text-grey q-mt-sm">
              {{ t('exportFlow.authRequestDialog.exportSuccessHint') }}
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="t('exportFlow.authRequestDialog.cancel')" color="grey" size="sm" @click="onCancel" />

        <!-- Step 1 actions -->
        <q-btn
          v-if="currentStep === 1"
          unelevated
          :label="t('exportFlow.authRequestDialog.next')"
          color="deep-purple"
          size="sm"
          :disable="!selectedSignatureId"
          @click="goToStep2"
        />

        <!-- Step 2 actions -->
        <q-btn
          v-if="currentStep === 2"
          flat
          :label="t('exportFlow.authRequestDialog.back')"
          color="primary"
          size="sm"
          @click="goToStep1"
        />
        <q-btn
          v-if="currentStep === 2"
          unelevated
          :label="t('exportFlow.authRequestDialog.exportRequest')"
          color="deep-purple"
          size="sm"
          :loading="exporting"
          @click="exportAuthRequest"
        />

        <!-- Step 3 actions -->
        <q-btn
          v-if="currentStep === 3"
          unelevated
          :label="t('exportFlow.authRequestDialog.done')"
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
import { generateAuthRequest } from 'boot/query/signature-query';

interface SignatureItem {
  id: string;
  name: string;
  intro?: string;
  image?: string;
  isAuthorized?: boolean;
  /** 资格码指纹，用于前端展示 */
  qualificationFingerprint?: string;
}

interface AuthRequestDialogProps {
  visible: boolean;
  signatures: SignatureItem[];
  contactEmail: string;
  contactAdditional?: string;
  authorizationUUID: string;
  originalAuthorQualificationCode: string;
}

interface AuthRequestDialogEmits {
  (e: 'cancel'): void;
  (e: 'done'): void;
  (e: 'createSignature'): void;
}

const props = withDefaults(defineProps<AuthRequestDialogProps>(), {
  visible: false,
  signatures: () => [],
  contactEmail: '',
  contactAdditional: '',
  authorizationUUID: '',
  originalAuthorQualificationCode: '',
});

const emit = defineEmits<AuthRequestDialogEmits>();

const { t } = useI18n();
const $q = useQuasar();

// UI State
const isVisible = ref(false);
const currentStep = ref(1);
const selectedSignatureId = ref('');
const exporting = ref(false);

// Computed
const unauthorizedSignatures = computed(() => {
  return props.signatures.filter((sig) => !sig.isAuthorized);
});

const hasAnySignatures = computed(() => {
  return props.signatures.length > 0;
});

const selectedSignature = computed(() => {
  return props.signatures.find((sig) => sig.id === selectedSignatureId.value);
});

const selectedSignatureName = computed(() => {
  return selectedSignature.value?.name || '';
});

// 联系方式计算属性
const contactEmail = computed(() => props.contactEmail);
const contactAdditional = computed(() => props.contactAdditional);

const stepDescription = computed(() => {
  switch (currentStep.value) {
    case 1:
      return t('exportFlow.authRequestDialog.step1Title');
    case 2:
      return t('exportFlow.authRequestDialog.step2Title');
    case 3:
      return t('exportFlow.authRequestDialog.step3Title');
    default:
      return '';
  }
});

// Watch
watch(
  () => props.visible,
  (newVal) => {
    isVisible.value = newVal;
    if (newVal) {
      // Reset state when dialog opens
      currentStep.value = 1;
      selectedSignatureId.value = '';
      exporting.value = false;
    }
  }
);

watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('cancel');
  }
});

// Methods
function selectSignature(id: string) {
  selectedSignatureId.value = id;
}

function goToStep1() {
  currentStep.value = 1;
}

function goToStep2() {
  if (selectedSignatureId.value) {
    currentStep.value = 2;
  }
}

async function copyText(text: string, label?: string) {
  try {
    await navigator.clipboard.writeText(text);
    $q.notify({
      type: 'positive',
      message: label
        ? t('exportFlow.authRequestDialog.copySuccess', { label })
        : t('exportFlow.notify.contactCopied'),
      position: 'top',
      timeout: 2000,
    });
  } catch (err) {
    console.error('Failed to copy:', err);
    $q.notify({
      type: 'negative',
      message: t('exportFlow.authRequestDialog.copyFailed'),
      position: 'top',
      timeout: 2000,
    });
  }
}

async function exportAuthRequest() {
  if (!selectedSignatureId.value || !props.authorizationUUID || !props.originalAuthorQualificationCode) {
    $q.notify({
      type: 'negative',
      message: t('exportFlow.authRequestDialog.missingParams'),
      position: 'top',
    });
    return;
  }

  exporting.value = true;

  try {
    const fileContent = await generateAuthRequest(
      props.authorizationUUID,
      selectedSignatureId.value,
      props.originalAuthorQualificationCode,
      selectedSignatureName.value
    );

    if (!fileContent) {
      throw new Error('Failed to generate auth request file');
    }

    const blob = new Blob([new Uint8Array(fileContent)], { type: 'application/octet-stream' });

    // 优先使用 File System Access API，确保用户选择保存路径后再提示成功
    if (typeof window.showSaveFilePicker === 'function') {
      try {
        const handle = await window.showSaveFilePicker({
          suggestedName: `auth-request-${Date.now()}.ktauthreq`,
          types: [
            {
              description: 'Authorization Request (.ktauthreq)',
              accept: { 'application/octet-stream': ['.ktauthreq'] },
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
      // 兼容不支持 File System Access API 的环境，使用下载链接回退方案
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `auth-request-${Date.now()}.ktauthreq`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    }

    // Go to success step
    currentStep.value = 3;

    $q.notify({
      type: 'positive',
      message: t('exportFlow.authRequestDialog.exportSuccess'),
      position: 'top',
    });
  } catch (err) {
    console.error('Export auth request failed:', err);
    $q.notify({
      type: 'negative',
      message: t('exportFlow.authRequestDialog.exportFailed'),
      position: 'top',
    });
  } finally {
    exporting.value = false;
  }
}

function onCreateSignature() {
  emit('createSignature');
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
.auth-request-dialog {
  border-radius: 8px;

  // 新的步骤条样式
  .stepper-container {
    display: flex;
    align-items: flex-start;
    justify-content: center;
  }

  .step-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 60px;

    .step-circle {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      background: #e8e8e8;
      color: #9e9e9e;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 600;
      transition: all 0.3s ease;
      border: 2px solid transparent;
    }

    .step-label {
      margin-top: 6px;
      font-size: 0.7rem;
      color: #9e9e9e;
      text-align: center;
      transition: color 0.3s ease;
      white-space: nowrap;
    }

    &.active {
      .step-circle {
        background: var(--q-deep-purple);
        color: white;
        border-color: var(--q-deep-purple);
        box-shadow: 0 2px 8px rgba(103, 58, 183, 0.3);
      }

      .step-label {
        color: var(--q-deep-purple);
        font-weight: 500;
      }
    }

    &.completed {
      .step-circle {
        background: var(--q-positive);
        color: white;
        border-color: var(--q-positive);
      }

      .step-label {
        color: var(--q-positive);
      }
    }
  }

  .step-connector {
    flex: 1;
    height: 2px;
    background: #e8e8e8;
    margin: 14px 8px 0;
    max-width: 40px;
    transition: background 0.3s ease;

    &.active {
      background: var(--q-positive);
    }
  }

  .signature-item {
    transition: all 0.2s ease;

    &:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    &.selected {
      border-color: var(--q-positive) !important;
      background: rgba(76, 175, 80, 0.05);
    }
  }

  .truncate {
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
    text-overflow: ellipsis;
  }
}
</style>
