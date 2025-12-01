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
import { nanoid } from 'nanoid';

export interface ExportSignatureFlowResult {
  needSignature: boolean;
  requireAuthorization?: boolean;
  contactEmail?: string;
  contactAdditional?: string;
  signatureId?: string;
  updateSignatureContent?: boolean;
  /**
   * 授权标识UUID
   *
   * 【生成时机】
   * - 首次导出选择"需要签名"时由前端 nanoid 生成
   * - 无论选择"需要授权"还是"无需授权"都会生成此UUID
   * - 再次导出时传空字符串，SDK会沿用已存储的UUID
   *
   * 【未来用途 - 签名授权导出/导入功能】
   * 授权是"签名+专辑"的特定授权，而非通用签名授权。
   * 用于授权申请文件和授权文件的生成校验。
   */
  authorizationUUID?: string;
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
    | 'auth-gate-from-picker' // 从签名选择页面打开的授权门控（用于导入授权文件）
    | 'picker' // 选择签名
    | 'done';
  flowData?: {
    needSignature?: boolean; // 是否需要签名
    requireAuthorization?: boolean; // 二次创作是否需要作者授权
    contactEmail?: string; // 邮箱
    contactAdditional?: string; // 额外联系方式
    updateSignatureContent?: boolean; // 是否更新签名内容
    /**
     * 授权标识UUID
     * - 首次导出时由nanoid生成（用户选择"需要签名"时）
     * - 再次导出时为空，SDK沿用已存储的UUID
     */
    authorizationUUID?: string;
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
    // 首次导出选择"需要签名"时，生成授权标识UUID
    // 此UUID无论后续选择"需要授权"还是"无需授权"都会被存储
    // 用于未来签名授权导出/导入功能的身份校验
    const authorizationUUID = payload.needSignature ? nanoid() : undefined;

    state.value.flowData = {
      ...(state.value.flowData ?? {}),
      needSignature: payload.needSignature,
      authorizationUUID,
    };

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

    // 优化：无论是否需要授权，都直接进入签名选择
    // 用户只能选择已授权的签名（如果requireAuthorization=true）
    // 如果需要导入授权，可以在签名选择页面点击"导入授权"按钮
    state.value.step = 'picker';
    pickerDialogVisible.value = true;
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

    // After auth, proceed to picker (or return to picker if opened from picker)
    state.value.step = 'picker';
    pickerDialogVisible.value = true;
  };

  /**
   * Handle auth gate dialog cancel.
   * 如果是从签名选择页面打开的授权门控，取消时应该返回签名选择页面
   */
  const handleAuthGateCancel = () => {
    authGateDialogVisible.value = false;
    // 如果当前是从签名选择页面触发的授权门控，返回签名选择页面
    if (state.value.step === 'auth-gate-from-picker') {
      state.value.step = 'picker';
      pickerDialogVisible.value = true;
    } else {
      state.value.step = 'idle';
    }
  };

  /**
   * Open auth gate dialog from signature picker.
   * 从签名选择页面打开授权门控（用于导入授权文件）
   */
  const openAuthGateFromPicker = () => {
    pickerDialogVisible.value = false;
    state.value.step = 'auth-gate-from-picker';
    authGateDialogVisible.value = true;
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
    // 首次导出时 authorizationUUID 由 nanoid 生成（在 handleConfirmSignatureSubmit 中）
    // 再次导出时 authorizationUUID 为 undefined，SDK 会沿用已存储的 UUID
    return {
      needSignature,
      requireAuthorization: state.value.flowData?.requireAuthorization,
      contactEmail: state.value.flowData?.contactEmail,
      contactAdditional: state.value.flowData?.contactAdditional,
      signatureId: state.value.selectedSignatureId,
      updateSignatureContent: state.value.flowData?.updateSignatureContent,
      authorizationUUID: state.value.flowData?.authorizationUUID,
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
    openAuthGateFromPicker,
    handlePickerSelect,
    handlePickerCreateNew,
    handleSignatureCreated,
    handlePickerCancel,
    getResult,
    reset,
  };
}
