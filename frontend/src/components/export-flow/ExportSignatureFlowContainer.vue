<template>
  <div class="export-signature-flow-container">
    <!-- Export Signature Policy Dialog -->
    <export-signature-policy-dialog
      :visible="policyDialogVisible"
      :album-has-signature="albumHasSignature"
      @submit="handlePolicySubmit"
      @cancel="handlePolicyCancel"
    />

    <!-- Authorization Gate Dialog -->
    <export-authorization-gate-dialog
      :visible="authGateDialogVisible"
      :author-contact="authorContact"
      @authorized="handleAuthGateAuthorized"
      @cancel="handleAuthGateCancel"
    />

    <!-- Signature Picker Dialog -->
    <signature-picker-dialog
      :visible="pickerDialogVisible"
      :signatures="mockSignatures"
      @select="handlePickerSelect"
      @createNew="handlePickerCreateNew"
      @cancel="handlePickerCancel"
    />

    <!-- Test/Demo Section -->
    <div class="q-pa-md">
      <div class="text-h6 q-mb-md">Export Signature Flow Demo</div>

      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-md-6">
          <q-card flat bordered>
            <q-card-section>
              <div class="text-subtitle2 q-mb-md">Album Info</div>
              <div class="q-mb-md">
                <q-checkbox v-model="albumHasSignature" label="Album has signature(s)" />
              </div>
              <div class="q-mb-md">
                <q-field label="Author Contact">
                  <template #control>
                    <div class="text-body2">{{ authorContact }}</div>
                  </template>
                </q-field>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-6">
          <q-card flat bordered>
            <q-card-section>
              <div class="text-subtitle2 q-mb-md">Flow State</div>
              <div class="q-mb-md">
                <div class="text-caption text-grey">Current Step:</div>
                <div class="text-body2">{{ currentStep }}</div>
              </div>
              <div class="q-mb-md">
                <div class="text-caption text-grey">Result:</div>
                <div class="text-body2">{{ JSON.stringify(flowResult, null, 2) }}</div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <div class="col-12">
          <q-btn color="primary" icon="cloud_upload" label="Start Export Flow" @click="startFlow" />
          <q-btn flat color="primary" label="Reset" @click="resetFlow" class="q-ml-md" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useExportSignatureFlow } from './useExportSignatureFlow';
import ExportSignaturePolicyDialog from '@/components/export-flow/ExportSignaturePolicyDialog.vue';
import ExportAuthorizationGateDialog from '@/components/export-flow/ExportAuthorizationGateDialog.vue';
import SignaturePickerDialog from '@/components/export-flow/SignaturePickerDialog.vue';

interface Signature {
  id: string;
  name: string;
  intro?: string;
  image?: string;
}

// Mock data
const mockSignatures: Signature[] = [
  {
    id: '1',
    name: 'Creator A',
    intro: 'Professional music producer',
    image: undefined,
  },
  {
    id: '2',
    name: 'Artist B',
    intro: 'Sound designer',
    image: undefined,
  },
  {
    id: '3',
    name: 'Musician C',
    intro: 'Composer and keyboardist',
    image: undefined,
  },
];

const { notify } = useQuasar();
const {
  currentStep,
  confirmSignatureDialogVisible: policyDialogVisible,
  authGateDialogVisible,
  pickerDialogVisible,
  start,
  handleConfirmSignatureSubmit: handlePolicySubmit,
  handleConfirmSignatureCancel: handlePolicyCancel,
  handleAuthGateAuthorized,
  handleAuthGateCancel,
  handlePickerSelect,
  handlePickerCreateNew,
  handlePickerCancel,
  getResult,
  reset,
} = useExportSignatureFlow();

const albumHasSignature = ref(false);
const authorContact = ref('contact@example.com');
const flowResult = ref({});

const startFlow = async () => {
  flowResult.value = {};
  await start({ albumHasSignature: albumHasSignature.value });
};

const resetFlow = () => {
  reset();
  flowResult.value = {};
};

const handlePolicySubmitWrapper = (data: any) => {
  handlePolicySubmit(data);

  // If no signature needed, show result
  if (!data.needSignature) {
    flowResult.value = getResult();
    notify({
      type: 'positive',
      message: 'Export flow completed: No signature required',
      position: 'top',
    });
  }
};

const handlePickerSelectWrapper = (signatureId: string) => {
  handlePickerSelect(signatureId);
  flowResult.value = getResult();
  notify({
    type: 'positive',
    message: `Export flow completed: Signature ${signatureId} selected`,
    position: 'top',
  });
};

const handlePickerCreateNewWrapper = () => {
  handlePickerCreateNew();
  notify({
    type: 'info',
    message: 'Create new signature (placeholder)',
    position: 'top',
  });
  // In real app, open signature creation dialog
};
</script>

<style scoped lang="scss">
// Container component styles are inherited from parent layout
</style>
