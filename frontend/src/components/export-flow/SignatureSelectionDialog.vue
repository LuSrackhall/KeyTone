<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card style="min-width: 700px; max-width: 900px">
      <q-card-section class="row items-center q-pb-none sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        <div class="text-h6">选择签名</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="cancel" />
      </q-card-section>

      <q-card-section v-if="loading" class="flex flex-center" style="min-height: 300px">
        <q-spinner color="primary" size="3em" />
        <div class="text-body2 text-grey-7 q-mt-md">正在加载可用签名...</div>
      </q-card-section>

      <q-card-section v-else-if="signatures.length === 0" class="text-center q-pa-lg">
        <q-icon name="inbox" size="64px" color="grey-5" />
        <div class="text-h6 text-grey-7 q-mt-md">暂无可用签名</div>
        <div class="text-body2 text-grey-6 q-mt-sm">请先创建签名后再导出专辑</div>
      </q-card-section>

      <q-card-section v-else class="q-pt-none">
        <!-- 筛选选项 -->
        <div class="row q-mb-md q-gutter-sm">
          <q-toggle v-model="showOnlyAuthorized" label="仅显示已授权" color="primary" />
          <q-toggle v-model="showOnlyInAlbum" label="仅显示已在专辑中" color="blue" />
        </div>

        <q-separator class="q-mb-md" />

        <!-- 签名列表 -->
        <div class="signatures-grid">
          <q-card
            v-for="sig in filteredSignatures"
            :key="sig.qualificationCode"
            flat
            bordered
            :class="{
              'signature-card': true,
              disabled: !sig.isAuthorized,
              selected: selectedSignature?.qualificationCode === sig.qualificationCode,
              'in-album': sig.isInAlbum,
            }"
            @click="selectSignature(sig)"
          >
            <q-card-section horizontal>
              <q-card-section class="col">
                <div class="row items-center q-mb-sm">
                  <div class="text-h6">{{ sig.name }}</div>
                </div>

                <div class="text-body2 text-grey-7 q-mb-md">
                  {{ sig.intro || '暂无介绍' }}
                </div>

                <!-- 徽章区域 -->
                <div class="badges-container">
                  <q-badge v-if="sig.isOriginalAuthor" color="amber-7" class="q-mr-xs">
                    <q-icon name="star" size="14px" class="q-mr-xs" />
                    原始作者
                  </q-badge>

                  <q-badge v-if="sig.isInAlbum" color="blue-7" class="q-mr-xs">
                    <q-icon name="check_circle" size="14px" class="q-mr-xs" />
                    已在专辑中
                  </q-badge>

                  <q-badge v-if="!sig.isAuthorized" color="negative" class="q-mr-xs" :title="'需要原始作者授权'">
                    <q-icon name="lock" size="14px" class="q-mr-xs" />
                    需要授权
                  </q-badge>

                  <q-badge v-else-if="sig.isAuthorized && !sig.isOriginalAuthor" color="positive">
                    <q-icon name="check" size="14px" class="q-mr-xs" />
                    已授权
                  </q-badge>
                </div>

                <!-- 资格码（小字显示） -->
                <div class="text-caption text-grey-5 q-mt-sm" style="word-break: break-all">
                  资格码: {{ sig.qualificationCode.substring(0, 16) }}...
                </div>
              </q-card-section>

              <!-- 选中标记 -->
              <q-card-section
                v-if="selectedSignature?.qualificationCode === sig.qualificationCode"
                class="col-auto flex flex-center"
              >
                <q-icon name="check_circle" color="primary" size="32px" />
              </q-card-section>
            </q-card-section>

            <!-- 未授权遮罩 -->
            <div v-if="!sig.isAuthorized" class="unauthorized-overlay">
              <q-icon name="lock" size="48px" color="white" />
              <div class="text-white text-body2 q-mt-sm">需要原始作者授权</div>
            </div>
          </q-card>
        </div>

        <!-- 未授权提示 -->
        <q-banner
          v-if="hasUnauthorizedSignatures && !showOnlyAuthorized"
          class="bg-orange-1 text-orange-9 q-mt-md"
          rounded
        >
          <template v-slot:avatar>
            <q-icon name="info" color="orange" />
          </template>
          部分签名需要原始作者授权后才能使用。请联系原始作者获取授权文件。
        </q-banner>
      </q-card-section>

      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat label="取消" color="grey-7" @click="cancel" />
        <q-btn unelevated label="确认选择" color="primary" :disable="!selectedSignature" @click="confirm" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { GetAvailableSignatures } from 'src/boot/query/keytonePkg-query';
import type { AvailableSignature } from 'src/types/export-flow';
import { Notify } from 'quasar';

interface Props {
  albumPath: string;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'confirm', signatureId: string): void;
  (e: 'cancel'): void;
}>();

const dialogVisible = ref(false);
const loading = ref(false);
const signatures = ref<AvailableSignature[]>([]);
const selectedSignature = ref<AvailableSignature | null>(null);
const showOnlyAuthorized = ref(false);
const showOnlyInAlbum = ref(false);

/**
 * 筛选后的签名列表
 */
const filteredSignatures = computed(() => {
  let result = signatures.value;

  if (showOnlyAuthorized.value) {
    result = result.filter((sig) => sig.isAuthorized);
  }

  if (showOnlyInAlbum.value) {
    result = result.filter((sig) => sig.isInAlbum);
  }

  return result;
});

/**
 * 是否有未授权的签名
 */
const hasUnauthorizedSignatures = computed(() => {
  return signatures.value.some((sig) => !sig.isAuthorized);
});

/**
 * 打开对话框
 */
async function open() {
  dialogVisible.value = true;
  loading.value = true;
  selectedSignature.value = null;

  try {
    signatures.value = await GetAvailableSignatures(props.albumPath);
    console.log('加载的签名列表:', signatures.value);
  } catch (err: any) {
    console.error('获取可用签名失败:', err);
    Notify.create({
      type: 'negative',
      message: '加载签名列表失败: ' + (err.message || '未知错误'),
    });
    dialogVisible.value = false;
  } finally {
    loading.value = false;
  }
}

/**
 * 选择签名
 */
function selectSignature(sig: AvailableSignature) {
  if (!sig.isAuthorized) {
    Notify.create({
      type: 'warning',
      message: '此签名需要原始作者授权后才能使用',
      caption: '请联系原始作者获取授权文件',
    });
    return;
  }

  selectedSignature.value = sig;
}

/**
 * 确认选择
 */
function confirm() {
  if (!selectedSignature.value) {
    Notify.create({
      type: 'warning',
      message: '请先选择一个签名',
    });
    return;
  }

  emit('confirm', selectedSignature.value.encryptedId);
  dialogVisible.value = false;
}

/**
 * 取消
 */
function cancel() {
  emit('cancel');
  dialogVisible.value = false;
}

// 暴露方法
defineExpose({
  open,
});
</script>

<style scoped lang="scss">
.signatures-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  max-height: 500px;
  overflow-y: auto;
  padding: 4px;
}

.signature-card {
  cursor: pointer;
  transition: all 0.3s;
  position: relative;

  &:hover:not(.disabled) {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }

  &.selected {
    border: 2px solid $primary;
    box-shadow: 0 4px 12px rgba($primary, 0.3);
  }

  &.in-album {
    border-left: 4px solid $blue-7;
  }

  &.disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.unauthorized-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  opacity: 0;
  transition: opacity 0.3s;

  .signature-card.disabled:hover & {
    opacity: 1;
  }
}

.badges-container {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
</style>
