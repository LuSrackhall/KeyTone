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
 * 这是一个轻量级的数据结构，避免加载完整的签名数据。
 *
 * @property hasSignature - 是否有签名
 * @property directExportAuthorName - 直接导出作者的名称
 * @property directExportAuthorImage - 直接导出作者的图片路径（相对于专辑目录）
 */
export interface AlbumSignatureSummary {
  /** 是否有签名 */
  hasSignature: boolean;
  /** 直接导出作者名称 */
  directExportAuthorName: string;
  /** 直接导出作者图片路径（相对于专辑目录） */
  directExportAuthorImage: string;
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
