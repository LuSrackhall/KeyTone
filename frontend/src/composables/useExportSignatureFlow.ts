import { ref, computed, Ref } from 'vue';

export interface ExportSignatureFlowResult {
  needSignature: boolean;
  requireAuthorization?: boolean;
  contact?: string;
  signatureId?: string;
}

export interface ExportSignatureFlowOptions {
  albumHasSignature?: boolean;
}

interface State {
  step: 'idle' | 'policy' | 'auth-gate' | 'picker' | 'done';
  policyData?: {
    needSignature: boolean;
    requireAuthorization: boolean;
    contact?: string;
  };
  isAuthorized: boolean;
  selectedSignatureId?: string;
}

/**
 * Composable for orchestrating album export signature flow.
 *
 * Flow:
 * 1. No signatures → show policy dialog
 * 2. Has signatures → skip policy, go to auth check
 * 3. If require auth & not authorized → show auth gate
 * 4. If need signature → show picker
 * 5. Done → return result
 */
export function useExportSignatureFlow() {
  const state: Ref<State> = ref({
    step: 'idle',
    isAuthorized: false,
  });

  // Dialog visibility refs
  const policyDialogVisible = ref(false);
  const authGateDialogVisible = ref(false);
  const pickerDialogVisible = ref(false);

  // Computed properties
  const currentStep = computed(() => state.value.step);

  /**
   * Start the export flow.
   * @param options Configuration for the flow
   */
  const start = async (options: ExportSignatureFlowOptions = {}) => {
    const { albumHasSignature = false } = options;

    state.value.step = 'idle';
    state.value.policyData = undefined;
    state.value.isAuthorized = false;
    state.value.selectedSignatureId = undefined;

    // Step 1: Check if album has signatures
    if (!albumHasSignature) {
      // No signatures → show policy dialog
      state.value.step = 'policy';
      policyDialogVisible.value = true;
    } else {
      // Has signatures → skip policy, go to auth check
      await checkAuthAndProceed();
    }
  };

  /**
   * Handle policy dialog submission.
   */
  const handlePolicySubmit = async (data: {
    needSignature: boolean;
    requireAuthorization: boolean;
    contact?: string;
  }) => {
    state.value.policyData = data;

    if (!data.needSignature) {
      // No signature required → done
      state.value.step = 'done';
      return getResult();
    }

    // Signature required → check auth and proceed to picker
    await checkAuthAndProceed();
  };

  /**
   * Handle policy dialog cancel.
   */
  const handlePolicyCancel = () => {
    policyDialogVisible.value = false;
    state.value.step = 'idle';
  };

  /**
   * Check if authorization is needed, and proceed accordingly.
   */
  const checkAuthAndProceed = async () => {
    const requiresAuth = state.value.policyData?.requireAuthorization || false;

    if (requiresAuth && !state.value.isAuthorized) {
      // Need auth and not authorized → show auth gate
      state.value.step = 'auth-gate';
      authGateDialogVisible.value = true;
    } else {
      // Authorized or no auth required → show picker
      state.value.step = 'picker';
      pickerDialogVisible.value = true;
    }
  };

  /**
   * Handle auth gate dialog authorization.
   */
  const handleAuthGateAuthorized = () => {
    state.value.isAuthorized = true;
    authGateDialogVisible.value = false;

    // After auth, proceed to picker
    state.value.step = 'picker';
    pickerDialogVisible.value = true;
  };

  /**
   * Handle auth gate dialog cancel.
   */
  const handleAuthGateCancel = () => {
    authGateDialogVisible.value = false;
    policyDialogVisible.value = false;
    state.value.step = 'idle';
  };

  /**
   * Handle signature picker dialog selection.
   */
  const handlePickerSelect = (signatureId: string) => {
    state.value.selectedSignatureId = signatureId;
    pickerDialogVisible.value = false;
    state.value.step = 'done';
  };

  /**
   * Handle signature picker dialog create new.
   */
  const handlePickerCreateNew = () => {
    // This will be handled by parent component
    // Parent should open signature creation dialog and call handleSignatureCreated on success
  };

  /**
   * Handle after signature is created in picker.
   */
  const handleSignatureCreated = (newSignatureId: string) => {
    // Reopen picker with updated list, auto-select the new signature
    state.value.selectedSignatureId = newSignatureId;
    pickerDialogVisible.value = true;
  };

  /**
   * Handle signature picker dialog cancel.
   */
  const handlePickerCancel = () => {
    pickerDialogVisible.value = false;
    authGateDialogVisible.value = false;
    policyDialogVisible.value = false;
    state.value.step = 'idle';
  };

  /**
   * Get the final result.
   */
  const getResult = (): ExportSignatureFlowResult => {
    const needSignature = state.value.policyData?.needSignature ?? true;
    return {
      needSignature,
      requireAuthorization: state.value.policyData?.requireAuthorization,
      contact: state.value.policyData?.contact,
      signatureId: state.value.selectedSignatureId,
    };
  };

  /**
   * Reset the flow to idle state.
   */
  const reset = () => {
    state.value = {
      step: 'idle',
      isAuthorized: false,
    };
    policyDialogVisible.value = false;
    authGateDialogVisible.value = false;
    pickerDialogVisible.value = false;
  };

  return {
    // State
    state,
    currentStep,
    policyDialogVisible,
    authGateDialogVisible,
    pickerDialogVisible,

    // Methods
    start,
    handlePolicySubmit,
    handlePolicyCancel,
    handleAuthGateAuthorized,
    handleAuthGateCancel,
    handlePickerSelect,
    handlePickerCreateNew,
    handleSignatureCreated,
    handlePickerCancel,
    getResult,
    reset,
  };
}
