/**
 * 签名排序元数据（仅用于本地排序，不参与导出）
 * Signature sort metadata (for local sorting only, not for export)
 */
export interface SignatureSortMetadata {
  /** Unix 时间戳，仅在创建或导入时生成；更新签名不会改变此值；但可通过拖动排序操作更改 */
  time: number;
}

/**
 * 签名数据结构（前端使用）
 * Signature data structure (frontend use)
 */
export interface Signature {
  /** 唯一标识（nanoid，21字符） / Unique identifier (nanoid, 21 chars) */
  id: string;

  /** 签名名称(必填) */
  name: string;

  /** 个人介绍(选填) */
  intro: string;

  /** 名片图片路径（选填） / Card image path (optional) */
  cardImage: File | null;
}

/**
 * 签名在配置文件中的存储结构
 * Signature storage entry structure
 */
export interface SignatureStorageEntry {
  value: string; // 加密的签名数据
  sort: SignatureSortMetadata; // 排序元数据
}
