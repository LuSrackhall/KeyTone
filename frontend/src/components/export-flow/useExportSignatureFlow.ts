/**
 * Centralized state machine for the album export signature experience.
 *
 * This composable intentionally lives alongside the export-flow UI components and
 * encodes the production-ready sequence of dialogs (policy → authorization → picker).
 * Integrated with real backend APIs for signature management.
 */
import { ref, computed, Ref } from 'vue';
import { GetAlbumSignatureInfo } from 'src/boot/query/keytonePkg-query';
import type { AlbumSignatureInfo } from 'src/types/export-flow';

export interface ExportSignatureFlowResult {
  needSignature: boolean;
  requireAuthorization?: boolean;
  contactEmail?: string;
  contactAdditional?: string;
  signatureId?: string;
  updateSignatureContent?: boolean;
}

export interface ExportSignatureFlowOptions {
  /** 专辑路径，用于获取签名信息 */
  albumPath: string;
}

interface State {
  step:
    | 'idle'
    | 'confirm-signature' // 没有任何签名时，先确认是否需要签名
    | 're-export-warning' // 再次导出时的警告提示
    | 'auth-requirement' // 二次创作是否需要授权（推荐不需要）
    | 'auth-impact-confirm' // 选择需要授权后的二次确认弹窗
    | 'auth-contact' // 需要授权时的联系方式填写
    | 'auth-gate' // 已有签名且原作者要求授权时的授权门控
    | 'picker' // 选择签名
    | 'done';
  flowData?: {
    needSignature?: boolean; // 是否需要签名
    requireAuthorization?: boolean; // 二次创作是否需要作者授权
    contactEmail?: string; // 邮箱
    contactAdditional?: string; // 额外联系方式
    updateSignatureContent?: boolean; // 是否更新签名内容
  };
  isAuthorized: boolean;
  selectedSignatureId?: string;
  signatureInfo?: AlbumSignatureInfo;
  albumPath?: string;
}
/**
 * Composable for orchestrating album export signature flow.
 *
 * Flow:
 * 1. No signatures → show policy dialog, then always evaluate authorization requirements
 * 2. Has signatures → show re-export warning
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
  const reExportWarningDialogVisible = ref(false);
  const authRequirementDialogVisible = ref(false);
  const authImpactConfirmDialogVisible = ref(false);
  const authContactDialogVisible = ref(false);
  const authGateDialogVisible = ref(false);
  const pickerDialogVisible = ref(false);

  // Computed properties
  const currentStep = computed(() => state.value.step);
  const requireAuthorizationForPicker = computed(() => {
    return state.value.signatureInfo?.originalAuthor?.requireAuthorization ?? false;
  });
  const currentAlbumPath = computed(() => state.value.albumPath);

  /**
   * Start the export flow.
   * @param options Configuration for the flow
   */
  const start = async (options: ExportSignatureFlowOptions) => {
    const { albumPath } = options;

    state.value.step = 'idle';
    state.value.flowData = undefined;
    state.value.isAuthorized = false;
    state.value.selectedSignatureId = undefined;
    state.value.signatureInfo = undefined;
    state.value.albumPath = albumPath;

    // 新逻辑：使用真实API获取专辑签名信息
    try {
      const signatureInfo = await GetAlbumSignatureInfo(albumPath);
      state.value.signatureInfo = signatureInfo;

      // 情况1：专辑无签名 → 首次导出流程
      if (!signatureInfo.hasSignature) {
        state.value.step = 'confirm-signature';
        confirmSignatureDialogVisible.value = true;
        return;
      }

      // 情况2 & 3：专辑有签名 → 再次导出流程
      // 显示再次导出警告对话框
      state.value.step = 're-export-warning';
      reExportWarningDialogVisible.value = true;
    } catch (error) {
      console.error('获取专辑签名信息失败:', error);
      // 出错时默认按首次导出处理
      state.value.step = 'confirm-signature';
      confirmSignatureDialogVisible.value = true;
    }
  };

  // ========== Step: confirm-signature ==========
  const handleConfirmSignatureSubmit = (payload: { needSignature: boolean }) => {
    state.value.flowData = { ...(state.value.flowData ?? {}), needSignature: payload.needSignature };

    confirmSignatureDialogVisible.value = false;

    // 如果选择"无需签名"，直接完成，不进入授权流程
    if (!payload.needSignature) {
      state.value.step = 'done';
      return;
    }

    // 选择"需要签名"，进入授权判断
    state.value.step = 'auth-requirement';
    authRequirementDialogVisible.value = true;
  };
  const handleConfirmSignatureCancel = () => {
    confirmSignatureDialogVisible.value = false;
    state.value.step = 'idle';
  };

  // ========== Step: re-export-warning ==========
  const handleReExportConfirm = () => {
    reExportWarningDialogVisible.value = false;

    // 检查是否需要授权
    const requireAuthorization = state.value.signatureInfo?.originalAuthor?.requireAuthorization;

    if (requireAuthorization) {
      // 需要授权，进入授权门控
      state.value.step = 'auth-gate';
      authGateDialogVisible.value = true;
    } else {
      // 无需授权，直接进入签名选择
      state.value.step = 'picker';
      pickerDialogVisible.value = true;
    }
  };

  const handleReExportCancel = () => {
    reExportWarningDialogVisible.value = false;
    state.value.step = 'idle';
  };

  // ========== Step: auth-requirement ==========
  const handleAuthRequirementSubmit = (payload: { requireAuthorization: boolean }) => {
    state.value.flowData = { ...(state.value.flowData ?? {}), requireAuthorization: payload.requireAuthorization };
    authRequirementDialogVisible.value = false;

    if (!payload.requireAuthorization) {
      const needSignature = state.value.flowData?.needSignature ?? true;
      if (needSignature) {
        // 仍需落地签名，继续选择
        state.value.step = 'picker';
        pickerDialogVisible.value = true;
      } else {
        // 无需签名且无需授权 → 就地完成
        state.value.step = 'done';
      }
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
  const handleAuthContactSubmit = (payload: { email: string; additional?: string }) => {
    state.value.flowData = {
      ...(state.value.flowData ?? {}),
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
  const handlePickerSelect = (signatureId: string, updateContent = true) => {
    state.value.selectedSignatureId = signatureId;
    state.value.flowData = {
      ...(state.value.flowData ?? {}),
      updateSignatureContent: updateContent,
    };
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
      contactEmail: state.value.flowData?.contactEmail,
      contactAdditional: state.value.flowData?.contactAdditional,
      signatureId: state.value.selectedSignatureId,
      updateSignatureContent: state.value.flowData?.updateSignatureContent,
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
    requireAuthorizationForPicker,
    currentAlbumPath,
    confirmSignatureDialogVisible,
    reExportWarningDialogVisible,
    authRequirementDialogVisible,
    authImpactConfirmDialogVisible,
    authContactDialogVisible,
    authGateDialogVisible,
    pickerDialogVisible,

    // Methods
    start,
    handleConfirmSignatureSubmit,
    handleConfirmSignatureCancel,
    handleReExportConfirm,
    handleReExportCancel,
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
