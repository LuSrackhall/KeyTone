<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="width: 90%; max-width: 400px">
      <!-- Header -->
      <q-card-section class="bg-info text-white q-pa-sm">
        <div class="text-subtitle1">🧪 {{ $t('exportFlow.test.title') }}</div>
      </q-card-section>

      <!-- Description -->
      <q-card-section class="q-pa-md">
        <div class="text-caption q-mb-md" style="line-height: 1.4">
          {{ $t('exportFlow.test.description') }}
        </div>

        <!-- Album State Configuration -->
        <div class="text-overline text-grey q-mb-sm">{{ $t('exportFlow.test.albumState') }}</div>

        <!-- Has Signature Toggle -->
        <q-checkbox
          v-model="testState.albumHasSignature"
          :label="$t('exportFlow.test.hasSignature')"
          color="primary"
          class="text-caption q-mb-sm"
        />
        <div class="text-caption text-grey q-ml-lg q-mb-md">
          {{ $t('exportFlow.test.hasSignatureHint') }}
        </div>

        <!-- Require Authorization Toggle (only when has signature) -->
        <q-checkbox
          v-model="testState.requireAuthorization"
          :disable="!testState.albumHasSignature"
          :label="$t('exportFlow.test.requireAuth')"
          color="primary"
          class="text-caption q-mb-sm"
        />
        <div class="text-caption text-grey q-ml-lg q-mb-md">
          {{ $t('exportFlow.test.requireAuthHint') }}
        </div>

        <!-- Current State Display -->
        <q-separator class="q-my-md" />
        <div class="text-caption text-weight-bold q-mb-xs">{{ $t('exportFlow.test.currentState') }}</div>
        <div class="bg-grey-2 rounded-lg q-pa-sm text-caption" style="line-height: 1.5; font-family: monospace">
          <div>
            Album Has Signature: <strong>{{ testState.albumHasSignature ? '✓ YES' : '✗ NO' }}</strong>
          </div>
          <div>
            Require Authorization: <strong>{{ testState.requireAuthorization ? '✓ YES' : '✗ NO' }}</strong>
          </div>
          <div class="q-mt-xs" style="font-size: 0.8rem; color: #666">
            Expected Flow:
            <div v-if="!testState.albumHasSignature">
              → Confirm Signature Dialog → (Auth Requirement / Auth Contact / Picker)
            </div>
            <div v-else-if="testState.requireAuthorization">→ Auth Gate Dialog → Picker</div>
            <div v-else>→ Picker directly</div>
          </div>
        </div>

        <!-- Flow Explanation -->
        <q-expansion-item
          header-class="text-caption text-weight-bold"
          :label="$t('exportFlow.test.flowExplain')"
          class="q-mt-md bg-blue-1 rounded-lg"
        >
          <div class="text-caption q-pa-md" style="line-height: 1.6">
            <div class="q-mb-sm">
              <strong>No Signature (albumHasSignature = false)</strong>
              <div class="text-grey q-mt-xs">
                1. Show "Confirm Signature" dialog<br />
                &nbsp;&nbsp;- If "No Signature" → Direct export, done<br />
                &nbsp;&nbsp;- If "Require Signature" → Next step<br />
                2. Show "Auth Requirement" dialog<br />
                &nbsp;&nbsp;- If "No Auth" → Go to picker<br />
                &nbsp;&nbsp;- If "Require Auth" → Next step<br />
                3. Show "Auth Impact Confirm" dialog<br />
                4. Show "Auth Contact" dialog (required)<br />
                5. Show "Picker" dialog
              </div>
            </div>
            <div>
              <strong>Has Signature (albumHasSignature = true)</strong>
              <div class="text-grey q-mt-xs">
                1. Skip "Confirm Signature"<br />
                2. If requireAuthorization = true:<br />
                &nbsp;&nbsp;- Show "Auth Gate" dialog → Then picker<br />
                3. If requireAuthorization = false:<br />
                &nbsp;&nbsp;- Show "Picker" directly
              </div>
            </div>
          </div>
        </q-expansion-item>
      </q-card-section>

      <!-- Actions -->
      <q-card-actions align="right" class="q-pa-sm q-gutter-xs">
        <q-btn flat size="sm" color="primary" :label="$t('exportFlow.test.close')" @click="onClose" />
        <q-btn unelevated size="sm" color="info" :label="$t('exportFlow.test.apply')" @click="onApply" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  visible: boolean;
  albumHasSignature: boolean;
  requireAuthorization: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void;
  (e: 'apply', payload: { albumHasSignature: boolean; requireAuthorization: boolean }): void;
}>();

const { t } = useI18n();

const isVisible = computed({
  get: () => props.visible,
  set: (v) => emit('update:visible', v),
});

// 本地测试状态副本，防止用户取消时丢失
const testState = ref({
  albumHasSignature: props.albumHasSignature,
  requireAuthorization: props.requireAuthorization,
});

// 监听 props 变化，同步到本地副本
watch(
  () => ({ hasSignature: props.albumHasSignature, requireAuth: props.requireAuthorization }),
  ({ hasSignature, requireAuth }) => {
    testState.value.albumHasSignature = hasSignature;
    testState.value.requireAuthorization = requireAuth;
  }
);

function onClose() {
  isVisible.value = false;
}

function onApply() {
  // 确保逻辑一致：若无签名，则不能启用"需要授权"
  const finalRequireAuth = testState.value.albumHasSignature ? testState.value.requireAuthorization : false;

  emit('apply', {
    albumHasSignature: testState.value.albumHasSignature,
    requireAuthorization: finalRequireAuth,
  });

  isVisible.value = false;
}
</script>
