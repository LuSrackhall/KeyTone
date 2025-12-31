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
 * 文件说明: keytone-album/types.ts - 键音专辑编辑器类型定义
 * ============================================================================
 *
 * 【文件作用】
 * 本文件定义了"键音专辑编辑器"组件拆分重构所需的所有 TypeScript 类型。
 * 它是整个 keytone-album 模块的"类型契约"，确保父组件与子组件之间的数据传递类型安全。
 *
 * 【核心概念: Context 模式】
 * 本重构采用 Vue 的 provide/inject 机制实现"上下文共享"：
 *
 *   ┌─────────────────────────────────────────────────────────────────┐
 *   │                    Keytone_album.vue (父组件)                    │
 *   │  ┌─────────────────────────────────────────────────────────┐   │
 *   │  │  所有状态 (refs/computed/reactive)                       │   │
 *   │  │  所有方法 (saveSoundConfig, deleteSound, etc.)           │   │
 *   │  │  SSE 监听 & 配置写回逻辑                                  │   │
 *   │  └─────────────────────────────────────────────────────────┘   │
 *   │                          │                                      │
 *   │                    provide(CONTEXT_KEY, ctx)                    │
 *   │                          │                                      │
 *   │    ┌─────────────────────┼─────────────────────┐               │
 *   │    ▼                     ▼                     ▼               │
 *   │ ┌──────────┐      ┌──────────┐         ┌──────────┐           │
 *   │ │ Step1    │      │ Step2    │   ...   │ Dialog   │           │
 *   │ │ inject() │      │ inject() │         │ inject() │           │
 *   │ └──────────┘      └──────────┘         └──────────┘           │
 *   └─────────────────────────────────────────────────────────────────┘
 *
 * 【为什么用 Context 而不是 Props】
 * 1. 原组件有 100+ 个状态/方法，用 props 会导致 "props 爆炸"
 * 2. Dialog 可能被任意 Step 调用，用 Context 更灵活
 * 3. 保持"单一数据源"在父组件，子组件只是 UI 展示层
 *
 * 【关联文件】
 * - index.ts              : 模块入口，导出类型和组件
 * - steps/*.vue           : Step 子组件，通过 inject 获取 Context
 * - dialogs/*.vue         : Dialog 子组件，通过 inject 获取 Context
 * - Keytone_album.vue     : 父组件，通过 provide 提供 Context (待集成)
 *
 * 【类型分类】
 * 1. 基础类型: SoundFileInfo, SoundEntry, KeySoundEntry 等 - 描述数据结构
 * 2. 参数类型: SaveSoundConfigParams 等 - 描述操作函数的参数
 * 3. Context 类型: KeytoneAlbumContext - 父组件向子组件共享的完整上下文
 * 4. 注入 Key: KEYTONE_ALBUM_CONTEXT_KEY - Vue provide/inject 使用的 Symbol
 *
 * ============================================================================
 */

import type { Ref, ComputedRef } from 'vue';
import type { DependencyIssue } from 'src/utils/dependencyValidator';

// ============================================================================
// 基础类型定义 - 描述键音专辑中的核心数据结构
// ============================================================================

/** 音频源文件信息 */
export interface SoundFileInfo {
  sha256: string;
  name_id: string;
  name: string;
  type: string;
}

/** 声音裁剪参数 */
export interface SoundCutParams {
  start_time: number;
  end_time: number;
  volume: number;
}

/** 声音源文件引用 */
export interface SoundSourceFileRef {
  sha256: string;
  name_id: string;
  type: string;
}

/** 声音值结构 */
export interface SoundValue {
  cut: SoundCutParams;
  name: string;
  source_file_for_sound: SoundSourceFileRef;
}

/** 声音条目 */
export interface SoundEntry {
  soundKey: string;
  soundValue: SoundValue;
}

/** 按键音播放模式 */
export type PlayMode = 'single' | 'random' | 'loop';

/** 按键音按压/释放配置 */
export interface KeySoundTriggerConfig {
  mode: PlayMode | { mode: PlayMode };
  value: Array<any>;
}

/** 按键音值结构 */
export interface KeySoundValue {
  name: string;
  down: KeySoundTriggerConfig;
  up: KeySoundTriggerConfig;
}

/** 按键音条目 */
export interface KeySoundEntry {
  keySoundKey: string;
  keySoundValue: KeySoundValue;
}

/** 选项类型标识 */
export type OptionType = 'audio_files' | 'sounds' | 'key_sounds';

/** 混合选项项 */
export interface MixedOptionItem {
  type: OptionType;
  value: any;
}

// ============================================================================
// 操作函数类型
// ============================================================================

/** 保存声音配置参数 */
export interface SaveSoundConfigParams {
  soundKey?: string;
  source_file_for_sound: SoundSourceFileRef;
  name: string;
  cut: SoundCutParams;
  onSuccess?: () => void;
}

/** 删除声音参数 */
export interface DeleteSoundParams {
  soundKey: string;
  onSuccess?: () => void;
}

/** 预览声音参数 */
export interface PreviewSoundParams {
  source_file_for_sound: SoundSourceFileRef;
  cut: SoundCutParams;
}

/** 保存按键音配置参数 */
export interface SaveKeySoundConfigParams {
  key: string;
  name: string;
  down: { mode: string; value: Array<any> };
  up: { mode: string; value: Array<any> };
}

/** 删除按键音参数 */
export interface DeleteKeySoundParams {
  keySoundKey: string;
  onSuccess?: () => void;
}

/** 保存全局联动声效配置参数 */
export interface SaveUnifiedSoundEffectParams {
  down: any;
  up: any;
}

/** 保存单键联动声效配置参数 */
export interface SaveSingleKeySoundEffectParams {
  singleKeys: Array<number>;
  down: any;
  up: any;
}

// ============================================================================
// Context 类型（provide/inject 共享）
// ============================================================================

/**
 * KeytoneAlbumContext - 父组件向子组件提供的上下文
 *
 * 包含所有需要在 Step 和 Dialog 组件间共享的状态与方法。
 * 子组件通过 inject(KEYTONE_ALBUM_CONTEXT_KEY) 获取。
 */
export interface KeytoneAlbumContext {
  // ============ Props ============
  pkgPath: string;
  isCreate: boolean;

  // ============ 核心状态 ============
  step: Ref<number>;
  pkgName: Ref<string>;

  // ============ Step1: 音频源文件相关 ============
  addNewSoundFile: Ref<boolean>;
  files: Ref<Array<File>>;
  editSoundFile: Ref<boolean>;
  soundFileList: Ref<Array<SoundFileInfo>>;
  selectedSoundFile: Ref<SoundFileInfo>;

  // ============ Step2: 声音定义相关 ============
  createNewSound: Ref<boolean>;
  soundName: Ref<string>;
  sourceFileForSound: Ref<SoundFileInfo>;
  soundStartTime: Ref<number>;
  soundEndTime: Ref<number>;
  soundVolume: Ref<number>;
  showEditSoundDialog: Ref<boolean>;
  soundList: Ref<Array<SoundEntry>>;
  selectedSound: Ref<SoundEntry | undefined>;

  // ============ Step3: 按键音相关 ============
  createNewKeySound: Ref<boolean>;
  keySoundName: Ref<string>;
  configureDownSound: Ref<boolean>;
  configureUpSound: Ref<boolean>;
  selectedSoundsForDown: Ref<Array<any>>;
  playModeForDown: Ref<string>;
  maxSelectionForDown: ComputedRef<number>;
  downTypeGroup: Ref<Array<string>>;
  downSoundList: ComputedRef<Array<MixedOptionItem>>;
  selectedSoundsForUp: Ref<Array<any>>;
  playModeForUp: Ref<string>;
  maxSelectionForUp: ComputedRef<number>;
  upTypeGroup: Ref<Array<string>>;
  upSoundList: ComputedRef<Array<MixedOptionItem>>;
  editExistingKeySound: Ref<boolean>;
  edit_configureDownSound: Ref<boolean>;
  edit_configureUpSound: Ref<boolean>;
  edit_downTypeGroup: Ref<Array<string>>;
  edit_upTypeGroup: Ref<Array<string>>;
  edit_downSoundList: ComputedRef<Array<MixedOptionItem>>;
  edit_upSoundList: ComputedRef<Array<MixedOptionItem>>;
  keySoundList: Ref<Array<KeySoundEntry>>;
  selectedKeySound: Ref<any>;

  // ============ Step4: 联动声效相关 ============
  isEnableEmbeddedTestSound: { down: boolean; up: boolean };
  showEveryKeyEffectDialog: Ref<boolean>;
  keyDownUnifiedSoundEffectSelect: Ref<any>;
  keyUpUnifiedSoundEffectSelect: Ref<any>;
  unifiedTypeGroup: Ref<Array<string>>;
  keyUnifiedSoundEffectOptions: ComputedRef<Array<MixedOptionItem>>;
  isShowUltimatePerfectionKeySoundAnchoring: ComputedRef<boolean>;
  isAnchoringUltimatePerfectionKeySound: Ref<boolean>;
  showSingleKeyEffectDialog: Ref<boolean>;
  isShowAddOrSettingSingleKeyEffectDialog: Ref<boolean>;
  selectedSingleKeys: Ref<Array<number>>;
  isRecordingSingleKeys: Ref<boolean>;
  keyOptions: ComputedRef<Array<number>>;
  filterOptions: Ref<Array<number>>;
  isGetsFocused: Ref<boolean>;
  isDownSoundEffectSelectEnabled: Ref<boolean>;
  isUpSoundEffectSelectEnabled: Ref<boolean>;
  keyDownSingleKeySoundEffectSelect: Ref<any>;
  keyUpSingleKeySoundEffectSelect: Ref<any>;
  singleKeyTypeGroup: Ref<Array<string>>;
  keySingleKeySoundEffectOptions: ComputedRef<Array<MixedOptionItem>>;
  isShowUltimatePerfectionKeySoundAnchoring_singleKey: ComputedRef<boolean>;
  isAnchoringUltimatePerfectionKeySound_singleKey: Ref<boolean>;
  keysWithSoundEffect: Ref<Map<string, any>>;
  isShowSingleKeySoundEffectEditDialog: Ref<boolean>;
  currentEditingKey: Ref<number | null>;
  currentEditingKey_old: number | null; // 非 ref，用于记录旧值
  currentEditingKeyOfName: ComputedRef<string>;
  keyDownSingleKeySoundEffectSelect_edit: Ref<any>;
  keyUpSingleKeySoundEffectSelect_edit: Ref<any>;
  keyDownSingleKeySoundEffectSelect_edit_old: any; // 非 ref，用于记录旧值
  keyUpSingleKeySoundEffectSelect_edit_old: any; // 非 ref，用于记录旧值
  singleKeyTypeGroup_edit: Ref<Array<string>>;
  keySingleKeySoundEffectOptions_edit: ComputedRef<Array<MixedOptionItem>>;
  isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit: ComputedRef<boolean>;
  isAnchoringUltimatePerfectionKeySound_singleKey_edit: Ref<boolean>;

  // ============ 依赖校验 ============
  dependencyIssues: Ref<DependencyIssue[]>;

  // ============ 工具函数 ============
  /** 选项标签显示函数 */
  album_options_select_label: (item: any) => any;
  /** 自然排序函数 */
  naturalSort: (a: string, b: string) => number;
  /** 防止录制时默认键盘行为 */
  preventDefaultKeyBehaviorWhenRecording: (event: KeyboardEvent) => void;
  /** 防止录制时默认鼠标行为 */
  preventDefaultMouseWhenRecording: (event: MouseEvent) => void;
  /** 单键录制：设置 clear_flag（避免录制瞬间误记录鼠标行为） */
  setSingleKeyRecordingClearFlag: () => void;
  /** 值转换函数（配置文件格式 -> UI格式） */
  convertValue: (item: any) => any;

  // ============ 操作函数 ============
  saveSoundConfig: (params: SaveSoundConfigParams) => void;
  deleteSound: (params: DeleteSoundParams) => void;
  previewSound: (params: PreviewSoundParams) => void;
  saveKeySoundConfig: (params: SaveKeySoundConfigParams, onSuccess?: () => void) => void;
  deleteKeySound: (params: DeleteKeySoundParams) => void;
  saveUnifiedSoundEffectConfig: (params: SaveUnifiedSoundEffectParams, onSuccess?: () => void) => void;
  saveSingleKeySoundEffectConfig: (params: SaveSingleKeySoundEffectParams, onSuccess?: () => void) => void;

  // ============ i18n ============
  $t: (key: string, params?: any) => string;

  // ============ 样式相关 ============
  i18n_fontSize: ComputedRef<string>;
  step_introduce_fontSize: ComputedRef<string>;
  isMacOS: Ref<boolean>;

  // ============ 选项常量 ============
  options: Array<{ label: string; value: string; label_0: string }>;
  playModeOptions: string[];
  playModeLabels: Map<string, string>;
}

/** Context 注入 key */
export const KEYTONE_ALBUM_CONTEXT_KEY = Symbol('KeytoneAlbumContext');
