/**
 * 在导出 KeyTone 专辑之前，在前端和 SDK 之间共享的有效负载，用于应用签名决策。
 */
export interface ApplySignatureConfigPayload {
  albumPath: string; // 签名专辑的路径
  needSignature: boolean; // 是否需要签名
  requireAuthorization: boolean; // 是否需要二次创作授权
  signatureId: string; // 选定的签名 ID（调用此路由时必填）
  contactEmail?: string; // 联系人邮箱（需要授权时必填）
  contactAdditional?: string; // 联系人补充信息（可选）
}
