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
// 专辑选择器相关组件导出入口
//
// 目录结构：
//   album-selector/
//   ├── index.ts                       <-- 当前文件：模块入口
//   ├── AlbumSignatureBadge.vue        <-- 签名徽章组件
//   └── AlbumSignatureHoverCard.vue    <-- 悬停详情卡片组件
//
// 使用示例：
//   import { AlbumSignatureBadge, AlbumSignatureHoverCard } from 'src/components/album-selector';
// ============================================================================

export { default as AlbumSignatureBadge } from './AlbumSignatureBadge.vue';
export { default as AlbumSignatureHoverCard } from './AlbumSignatureHoverCard.vue';
