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
  cardImage: File;
}
