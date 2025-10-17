/**
 * 签名数据结构（前端使用）
 * Signature data structure (frontend use)
 */
export interface Signature {
  /** 唯一标识（nanoid，21字符） / Unique identifier (nanoid, 21 chars) */
  id: string;

  /** 签名名称（必填，1-50字符） / Signature name (required, 1-50 chars) */
  name: string;

  /** 个人介绍（选填，0-500字符） / Personal introduction (optional, 0-500 chars) */
  intro: string;

  /** 名片图片路径（选填） / Card image path (optional) */
  cardImage: string;
}

/**
 * 签名管理器（前端使用）
 * Signature manager (frontend use)
 */
export type SignatureManager = Record<string, Signature>;

/**
 * .ktsign 文件格式
 * .ktsign file format
 */
export interface SignatureFile {
  /** 文件格式版本 / File format version */
  version: string;

  /** 签名数据 / Signature data */
  signature: {
    id: string;
    name: string;
    intro: string;
    /** Base64 编码的图片数据 / Base64 encoded image data */
    cardImage: string;
  };

  /** SHA-256 校验和 / SHA-256 checksum */
  checksum: string;
}

/**
 * 签名错误代码
 * Signature error codes
 */
export enum SignatureErrorCode {
  /** 签名不存在 / Signature not found */
  NOT_FOUND = 'SIGNATURE_NOT_FOUND',

  /** 签名已存在 / Signature already exists */
  ALREADY_EXISTS = 'SIGNATURE_ALREADY_EXISTS',

  /** 无效的签名数据 / Invalid signature data */
  INVALID_DATA = 'INVALID_SIGNATURE_DATA',

  /** 无效的文件格式 / Invalid file format */
  INVALID_FILE = 'INVALID_FILE_FORMAT',

  /** 校验和不匹配 / Checksum mismatch */
  CHECKSUM_MISMATCH = 'CHECKSUM_MISMATCH',

  /** 加密失败 / Encryption failed */
  ENCRYPTION_FAILED = 'ENCRYPTION_FAILED',

  /** 解密失败 / Decryption failed */
  DECRYPTION_FAILED = 'DECRYPTION_FAILED',

  /** 网络错误 / Network error */
  NETWORK_ERROR = 'NETWORK_ERROR',

  /** 未知错误 / Unknown error */
  UNKNOWN_ERROR = 'UNKNOWN_ERROR',
}

/**
 * 签名表单数据（用于创建/编辑）
 * Signature form data (for create/edit)
 */
export interface SignatureFormData {
  /** 签名名称 / Signature name */
  name: string;

  /** 个人介绍 / Personal introduction */
  intro: string;

  /** 名片图片文件对象（非路径字符串） / Card image file object (not path string) */
  cardImage?: File;
}

/**
 * 签名创建请求
 * Signature create request
 */
export interface SignatureCreateRequest {
  /** 签名名称 / Signature name */
  name: string;

  /** 个人介绍 / Personal introduction */
  intro: string;

  /** Base64 编码的图片数据 / Base64 encoded image data */
  cardImage?: string;
}

/**
 * 签名更新请求
 * Signature update request
 */
export interface SignatureUpdateRequest {
  /** 签名ID / Signature ID */
  id: string;

  /** 签名名称 / Signature name */
  name: string;

  /** 个人介绍 / Personal introduction */
  intro: string;

  /** Base64 编码的图片数据（如果更新图片） / Base64 encoded image data (if updating image) */
  cardImage?: string;
}

/**
 * 签名响应
 * Signature response
 */
export interface SignatureResponse {
  /** 是否成功 / Success status */
  success: boolean;

  /** 错误信息 / Error message */
  message?: string;

  /** 签名数据 / Signature data */
  data?: Signature;
}

/**
 * 签名列表响应
 * Signature list response
 */
export interface SignatureListResponse {
  /** 是否成功 / Success status */
  success: boolean;

  /** 错误信息 / Error message */
  message?: string;

  /** 签名数据 / Signature data */
  data?: SignatureManager;
}
