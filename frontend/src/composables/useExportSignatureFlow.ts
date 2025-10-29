import { ref, computed, Ref } from 'vue';

export interface ExportSignatureFlowResult {
  needSignature: boolean;
  requireAuthorization?: boolean;
  contact?: string;
  contactEmail?: string;
  contactAdditional?: string;
  signatureId?: string;
}

export interface ExportSignatureFlowOptions {
  albumHasSignature?: boolean;
  /** 临时测试开关：已有签名且原作者要求授权时为 true，用于触发展示授权门控对话框 */
  existingSignatureRequireAuthorization?: boolean;
}

interface State {
  step:
    | 'idle'
    | 'confirm-signature' // 没有任何签名时，先确认是否需要签名
    | 'auth-requirement' // 二次创作是否需要授权（推荐不需要）
    | 'auth-impact-confirm' // 选择需要授权后的二次确认弹窗
    | 'auth-contact' // 需要授权时的联系方式填写
    | 'auth-gate' // 已有签名且原作者要求授权时的授权门控
    | 'picker' // 选择签名
    | 'done';
  flowData?: {
    needSignature?: boolean; // 是否需要签名
    requireAuthorization?: boolean; // 二次创作是否需要作者授权
    contact?: string; // 授权联系方式（格式化字符串）
    contactEmail?: string; // 邮箱
    contactAdditional?: string; // 额外联系方式
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
  const confirmSignatureDialogVisible = ref(false);
  const authRequirementDialogVisible = ref(false);
  const authImpactConfirmDialogVisible = ref(false);
  const authContactDialogVisible = ref(false);
  const authGateDialogVisible = ref(false);
  const pickerDialogVisible = ref(false);

  // Computed properties
  const currentStep = computed(() => state.value.step);

  /**
   * Start the export flow.
   * @param options Configuration for the flow
   */
  const start = async (options: ExportSignatureFlowOptions = {}) => {
    const { albumHasSignature = false, existingSignatureRequireAuthorization = false } = options;

    state.value.step = 'idle';
    state.value.flowData = undefined;
    state.value.isAuthorized = false;
    state.value.selectedSignatureId = undefined;

    // Step 1: Check if album has signatures
    if (!albumHasSignature) {
      // 没有任何签名 → 先询问是否需要签名
      state.value.step = 'confirm-signature';
      confirmSignatureDialogVisible.value = true;
      return;
    }

    // 已有签名：如果要求授权则先显示授权门控；否则直接到签名选择
    if (existingSignatureRequireAuthorization) {
      state.value.step = 'auth-gate';
      authGateDialogVisible.value = true;
      return;
    }

    state.value.step = 'picker';
    pickerDialogVisible.value = true;
  };

  // ========== Step: confirm-signature ==========
  const handleConfirmSignatureSubmit = (payload: { needSignature: boolean }) => {
    state.value.flowData = { ...(state.value.flowData ?? {}), needSignature: payload.needSignature };

    if (!payload.needSignature) {
      // 用户选择无需签名 → 流程直接完成
      confirmSignatureDialogVisible.value = false;
      state.value.step = 'done';
      return;
    }

    // 需要签名 → 进入“是否需要授权”的选择
    confirmSignatureDialogVisible.value = false;
    state.value.step = 'auth-requirement';
    authRequirementDialogVisible.value = true;
  };
  const handleConfirmSignatureCancel = () => {
    confirmSignatureDialogVisible.value = false;
    state.value.step = 'idle';
  };

  // ========== Step: auth-requirement ==========
  const handleAuthRequirementSubmit = (payload: { requireAuthorization: boolean }) => {
    state.value.flowData = { ...(state.value.flowData ?? {}), requireAuthorization: payload.requireAuthorization };
    authRequirementDialogVisible.value = false;

    if (!payload.requireAuthorization) {
      // 不需要授权 → 直接进入签名选择
      state.value.step = 'picker';
      pickerDialogVisible.value = true;
      return;
    }

    // 需要授权 → 二次确认
    state.value.step = 'auth-impact-confirm';
    authImpactConfirmDialogVisible.value = true;
  };
  const handleAuthRequirementCancel = () => {
    authRequirementDialogVisible.value = false;
    state.value.step = 'idle';
  };

  // ========== Step: auth-impact-confirm ==========
  const handleAuthImpactBack = () => {
    authImpactConfirmDialogVisible.value = false;
    state.value.step = 'auth-requirement';
    authRequirementDialogVisible.value = true;
  };
  const handleAuthImpactConfirm = () => {
    // 前往填写联系方式
    authImpactConfirmDialogVisible.value = false;
    state.value.step = 'auth-contact';
    authContactDialogVisible.value = true;
  };

  // ========== Step: auth-contact ==========
  const handleAuthContactSubmit = (payload: { contact: string; email: string; additional?: string }) => {
    state.value.flowData = {
      ...(state.value.flowData ?? {}),
      contact: payload.contact,
      contactEmail: payload.email,
      contactAdditional: payload.additional?.trim() ? payload.additional.trim() : undefined,
    };
    authContactDialogVisible.value = false;
    // 进入签名选择
    state.value.step = 'picker';
    pickerDialogVisible.value = true;
  };
  const handleAuthContactCancel = () => {
    authContactDialogVisible.value = false;
    state.value.step = 'idle';
  };

  /**
   * Check if authorization is needed, and proceed accordingly.
   */
  const checkAuthAndProceed = async () => {
    // 此函数保留，仅供已有签名的专辑分支调用（根据 options 判定）
    if (!state.value.isAuthorized) {
      state.value.step = 'auth-gate';
      authGateDialogVisible.value = true;
      return;
    }
    state.value.step = 'picker';
    pickerDialogVisible.value = true;
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
    confirmSignatureDialogVisible.value = false;
    authRequirementDialogVisible.value = false;
    authImpactConfirmDialogVisible.value = false;
    authContactDialogVisible.value = false;
    state.value.step = 'idle';
  };

  /**
   * Get the final result.
   */
  const getResult = (): ExportSignatureFlowResult => {
    const needSignature = state.value.flowData?.needSignature ?? true;
    return {
      needSignature,
      requireAuthorization: state.value.flowData?.requireAuthorization,
      contact: state.value.flowData?.contact,
      contactEmail: state.value.flowData?.contactEmail,
      contactAdditional: state.value.flowData?.contactAdditional,
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
    confirmSignatureDialogVisible.value = false;
    authRequirementDialogVisible.value = false;
    authImpactConfirmDialogVisible.value = false;
    authContactDialogVisible.value = false;
    authGateDialogVisible.value = false;
    pickerDialogVisible.value = false;
  };

  return {
    // State
    state,
    currentStep,
    confirmSignatureDialogVisible,
    authRequirementDialogVisible,
    authImpactConfirmDialogVisible,
    authContactDialogVisible,
    authGateDialogVisible,
    pickerDialogVisible,

    // Methods
    start,
    handleConfirmSignatureSubmit,
    handleConfirmSignatureCancel,
    handleAuthRequirementSubmit,
    handleAuthRequirementCancel,
    handleAuthImpactBack,
    handleAuthImpactConfirm,
    handleAuthContactSubmit,
    handleAuthContactCancel,
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
