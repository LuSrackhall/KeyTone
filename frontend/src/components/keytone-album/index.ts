/*
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

/**
 * ============================================================================
 * 文件说明: keytone-album/index.ts - 键音专辑编辑器模块入口
 * ============================================================================
 *
 * 【文件作用】
 * 本文件是 keytone-album 模块的统一出口，负责：
 * 1. 导出所有公共类型定义（来自 types.ts）
 * 2. 导出主组件（当前阶段仍指向旧的 Keytone_album.vue）
 * 3. 未来：导出拆分后的子组件供外部按需使用
 *
 * 【向后兼容策略】
 * 重构采用"渐进迁移"策略，确保现有代码不受影响：
 *
 *   迁移前：import KeytoneAlbum from 'src/components/Keytone_album.vue'  ✓ 继续可用
 *   迁移后：import { KeytoneAlbum } from 'src/components/keytone-album'  ✓ 推荐方式
 *
 * 【目录结构】
 *   keytone-album/
 *   ├── index.ts          <-- 当前文件：模块入口
 *   ├── types.ts          <-- 类型定义
 *   ├── steps/            <-- Step 子组件目录
 *   │   ├── StepLoadAudioFiles.vue      (Step 1: 加载音频源文件)
 *   │   ├── StepDefineSounds.vue        (Step 2: 定义声音) [待创建]
 *   │   ├── StepCraftKeySounds.vue      (Step 3: 制作按键音) [待创建]
 *   │   └── StepLinkageEffects.vue      (Step 4: 联动声效) [待创建]
 *   ├── dialogs/          <-- Dialog 子组件目录
 *   │   ├── AddAudioFileDialog.vue      (添加音频文件对话框)
 *   │   ├── ManageAudioFilesDialog.vue  (管理音频文件对话框)
 *   │   └── ... (更多对话框待创建)
 *   └── composables/      <-- 可复用逻辑 [待创建]
 *
 * 【关联文件】
 * - types.ts                  : 被本文件导出的类型定义
 * - ../Keytone_album.vue      : 当前导出的主组件（原始大组件）
 *
 * ============================================================================
 */

// 导出所有类型定义，供外部使用
export * from './types';

// ============================================================================
// 主组件导出
// ============================================================================
//
// 【当前状态】阶段 1 - 保持兼容
// 导出原始的 Keytone_album.vue，确保现有引用不受影响。
//
// 【未来状态】阶段 2+ - 完成迁移后
// 切换为导出重构后的新组件 KeytoneAlbumEditor.vue。
//
export { default as KeytoneAlbum } from '../Keytone_album.vue';

// ============================================================================
// 子组件导出（迁移完成后启用）
// ============================================================================
//
// 以下导出在完成集成后取消注释，供需要单独使用子组件的场景：
//
// export { default as StepLoadAudioFiles } from './steps/StepLoadAudioFiles.vue';
// export { default as StepDefineSounds } from './steps/StepDefineSounds.vue';
// export { default as StepCraftKeySounds } from './steps/StepCraftKeySounds.vue';
// export { default as StepLinkageEffects } from './steps/StepLinkageEffects.vue';
