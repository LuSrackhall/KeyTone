/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// ============================================================================
// 专辑选择器相关类型定义
// 用于在专辑选择器中展示签名作者信息
// ============================================================================

/**
 * AlbumSignatureSummary - 专辑签名摘要信息
 *
 * 用于在专辑选择器列表中展示签名作者的基本信息，
 * 包含原始作者和直接导出作者的完整信息，供悬停卡片使用。
 *
 * 设计说明：
 * - 当 isSameAuthor 为 true 时，表示原始作者与直接导出作者是同一人
 * - 此时前端应只展示一个作者区块，避免信息重复
 *
 * @property hasSignature - 是否有签名
 * @property originalAuthorName - 原始作者的名称
 * @property originalAuthorImage - 原始作者的图片路径（相对于专辑目录）
 * @property originalAuthorIntro - 原始作者的介绍
 * @property directExportAuthorName - 直接导出作者的名称
 * @property directExportAuthorImage - 直接导出作者的图片路径（相对于专辑目录）
 * @property directExportAuthorIntro - 直接导出作者的介绍
 * @property isSameAuthor - 原始作者与直接导出作者是否为同一人
 */
export interface AlbumSignatureSummary {
  /** 是否有签名 */
  hasSignature: boolean;

  // ============================================================================
  // 原始作者信息
  // ============================================================================
  /** 原始作者名称 */
  originalAuthorName: string;
  /** 原始作者图片路径（相对于专辑目录） */
  originalAuthorImage: string;
  /** 原始作者介绍 */
  originalAuthorIntro: string;

  // ============================================================================
  // 直接导出作者信息
  // ============================================================================
  /** 直接导出作者名称 */
  directExportAuthorName: string;
  /** 直接导出作者图片路径（相对于专辑目录） */
  directExportAuthorImage: string;
  /** 直接导出作者介绍 */
  directExportAuthorIntro: string;

  // ============================================================================
  // 是否为同一作者标记
  // ============================================================================
  /** 原始作者与直接导出作者是否为同一人 */
  isSameAuthor: boolean;
}

/**
 * GetAudioPackageListResponse - 获取专辑列表 API 的响应类型
 *
 * 扩展了原有的列表响应，增加了签名摘要信息，
 * 用于在专辑列表展示时同时显示签名作者信息。
 */
export interface GetAudioPackageListResponse {
  /** 响应消息 */
  message: string;
  /** 专辑路径列表 */
  list: string[];
  /** 签名摘要信息，key 为专辑路径 */
  signatureInfo: Record<string, AlbumSignatureSummary>;
}
