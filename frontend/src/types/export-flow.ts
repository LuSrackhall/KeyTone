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

/**
 * 签名作者信息
 */
export interface SignatureAuthorInfo {
  qualificationCode: string; // 资格码（签名ID的SHA256哈希）
  name: string; // 签名名称
  intro: string; // 个人介绍
  cardImagePath: string; // 名片图片路径
  isOriginalAuthor: boolean; // 是否为原始作者
  requireAuthorization?: boolean; // 是否需要授权（仅原始作者有效）
  authorizedList?: string[]; // 已授权的资格码列表（仅原始作者有效）
}

/**
 * 专辑签名配置中的单个签名条目
 */
export interface AlbumSignatureEntry {
  name: string;
  intro: string;
  cardImagePath: string;
  authorization?: {
    requireAuthorization: boolean;
    contactEmail: string;
    contactAdditional?: string;
    authorizedList: string[];
    directExportAuthor: string;
  };
}

/**
 * 专辑签名信息（完整）
 */
export interface AlbumSignatureInfo {
  hasSignature: boolean; // 专辑是否包含签名
  originalAuthor?: SignatureAuthorInfo; // 原始作者签名信息
  contributorAuthors: SignatureAuthorInfo[]; // 历史贡献作者列表
  directExportAuthor?: SignatureAuthorInfo; // 直接导出作者信息
  allSignatures: Record<string, AlbumSignatureEntry>; // 所有签名条目
}

/**
 * 可用于导出的签名信息
 */
export interface AvailableSignature {
  encryptedId: string; // 加密的签名ID
  qualificationCode: string; // 资格码
  name: string; // 签名名称
  intro: string; // 个人介绍
  isInAlbum: boolean; // 是否已在专辑中
  isAuthorized: boolean; // 是否有导出授权
  isOriginalAuthor: boolean; // 是否为原始作者
}
