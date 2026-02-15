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
  <q-dialog v-model="dialogVisible" backdrop-filter="invert(70%)" @hide="handleClose">
    <q-card class="w-[500px]">
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ $t('signature.select.title') }}
      </q-card-section>

      <q-card-section>
        <div class="text-caption text-grey-7 mb-3">
          {{ $t('signature.select.description') }}
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center items-center py-10">
          <q-spinner color="primary" size="2em" />
        </div>

        <!-- 空状态 -->
        <div v-else-if="signatureList.length === 0" class="flex flex-col items-center py-10 text-gray-500">
          <q-icon name="badge" size="40px" class="mb-2 opacity-50" />
          <div class="text-sm mb-2">{{ $t('signature.select.emptyState') }}</div>
          <q-btn flat dense color="primary" @click="goToSignaturePage">
            {{ $t('signature.select.createSignature') }}
          </q-btn>
        </div>

        <!-- 签名列表 -->
        <div v-else class="max-h-96 overflow-y-auto">
          <q-list bordered separator>
            <q-item
              v-for="signature in signatureList"
              :key="signature.id"
              clickable
              v-ripple
              :active="selectedSignatureId === signature.id"
              @click="selectSignature(signature.id)"
            >
              <q-item-section avatar>
                <q-avatar size="40px" rounded>
                  <img
                    v-if="signature.cardImage"
                    :src="getImageUrl(signature.cardImage as unknown as string)"
                    class="object-cover"
                  />
                  <q-icon v-else name="person" size="24px" />
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label>{{ signature.name }}</q-item-label>
                <q-item-label caption lines="1">
                  {{ signature.intro || $t('signature.select.noIntro') }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-icon v-if="selectedSignatureId === signature.id" name="check_circle" color="primary" size="24px" />
              </q-item-section>
            </q-item>
          </q-list>
        </div>

        <!-- 无签名选项 -->
        <div class="mt-3">
          <q-checkbox
            v-model="noSignature"
            :label="$t('signature.select.noSignature')"
            @update:model-value="handleNoSignatureChange"
          />
        </div>
      </q-card-section>

      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="$t('signature.form.cancel')" color="primary" @click="handleClose" />
        <q-btn
          flat
          :label="$t('signature.select.confirm')"
          color="primary"
          @click="handleConfirm"
          :disable="!selectedSignatureId && !noSignature"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import type { Signature } from 'src/types/signature';

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'select', signatureId: string | null): void;
}>();

const q = useQuasar();
const { t: $t } = useI18n();
const router = useRouter();

// 对话框显示状态控制
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

// 数据加载状态 - 绑定加载动画显示
const loading = ref(false);

// 签名列表数据 - 由具体数据流实现填充，绑定到签名列表渲染
const signatureList = ref<Signature[]>([]);

// 当前选中的签名 ID - 绑定到列表项的 active 状态和右侧勾选图标
const selectedSignatureId = ref<string | null>(null);

// 是否选择"无签名"选项 - 绑定到 checkbox，控制 selectedSignatureId 置空
const noSignature = ref(false);

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      loadSignatures();
    }
  }
);

onMounted(() => {
  if (props.modelValue) {
    loadSignatures();
  }
});

/** 加载签名列表数据 */
async function loadSignatures() {
  loading.value = true;
  try {
    // TODO: 实现签名列表加载逻辑
    // 1. 调用 getSignaturesList() 获取加密的签名列表
    // 2. 处理新的 SignatureStorageEntry 结构：{ value: string, sort: { time: number } }
    // 3. 按 sort.time 时间戳排序（升序）
    // 4. 逐个解密并解析签名数据
    // 5. 获取签名图片 URL
    // 6. 填充 signatureList
    // 注意：与 Signature_management_page.vue 中的排序逻辑保持一致
    signatureList.value = [];
  } catch (err) {
    console.error('Failed to load signatures:', err);
    q.notify({
      type: 'negative',
      message: $t('signature.notify.loadFailed'),
      position: 'top',
    });
  } finally {
    loading.value = false;
  }
}

/** 选中一个签名 */
function selectSignature(id: string) {
  selectedSignatureId.value = id;
  noSignature.value = false;
}

/** 处理"无签名"选项变化 */
function handleNoSignatureChange(value: boolean) {
  if (value) {
    selectedSignatureId.value = null;
  }
}

/** 关闭对话框并重置状态 */
function handleClose() {
  selectedSignatureId.value = null;
  noSignature.value = false;
  dialogVisible.value = false;
}

/** 确认选择并发出信号 */
function handleConfirm() {
  emit('select', noSignature.value ? null : selectedSignatureId.value);
  handleClose();
}

/** 获取图片 URL - 由具体业务层实现 */
function getImageUrl(filename: string): string {
  // TODO: 具体 URL 生成逻辑由业务层实现
  return '';
}

/** 导航到签名管理页面 */
function goToSignaturePage() {
  handleClose();
  router.push('/signature-management');
}
</script>

<style scoped>
/* 自定义样式 */
</style>
