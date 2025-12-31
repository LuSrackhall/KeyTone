## Context

`Keytone_album.vue` 当前把"步骤 UI + 多层 dialog + SSE 映射 + 配置写回 + 依赖校验提示 + 样式"全部聚合在单文件中。

其核心风险在于：任何微小改动都会触发大范围重编译、难以定位回归，并且跨步骤复用 dialog 的诉求会进一步放大耦合。

## Goals / Non-Goals

### Goals

- **零行为变更**：拆分后用户侧行为保持一致。
- **边界清晰**：按 step 与对话框职责拆分，降低单文件复杂度。
- **可复用**：对话框支持被任意步骤直接拉起（不强绑在某个 step 内部）。
- **渐进迁移**：允许逐步拆分（先拆模板，再拆逻辑），每一步都可编译运行。

### Non-Goals

- 不做功能增强或 UI 重新设计
- 不引入新的状态管理框架
- 不改动后端接口与数据结构

---

## 架构设计详解

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Keytone_album.vue (父组件)                            │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        所有状态 & 副作用                                │  │
│  │  - refs: step, pkgName, soundFileList, soundList, keySoundList...     │  │
│  │  - computed: downSoundList, upSoundList, keyOptions...                │  │
│  │  - watch: SSE 监听、配置自动保存                                        │  │
│  │  - methods: saveSoundConfig, deleteSound, previewSound...             │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                    provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)  ✅ 已实现        │
│                                    │                                         │
│    ┌───────────────────────────────┼───────────────────────────────┐        │
│    │                               │                               │        │
│    ▼                               ▼                               ▼        │
│ ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐       │
│ │ StepLoadAudio    │    │ StepDefineSounds │    │ StepCraftKey     │  ...  │
│ │ Files.vue ✅     │    │ .vue ✅          │    │ Sounds.vue ✅    │       │
│ │                  │    │                  │    │                  │       │
│ │ inject(ctx)      │    │ inject(ctx)      │    │ inject(ctx)      │       │
│ │ 只负责 UI 渲染    │    │ 只负责 UI 渲染    │    │ 只负责 UI 渲染    │       │
│ └────────┬─────────┘    └────────┬─────────┘    └────────┬─────────┘       │
│          │                       │                       │                  │
│          ▼                       ▼                       ▼                  │
│    ┌───────────┐           ┌───────────┐           ┌───────────┐           │
│    │ Dialog ✅ │           │ Dialog ✅ │           │ Dialog ✅ │           │
│    │ 组件      │           │ 组件      │           │ 组件      │           │
│    │ inject()  │           │ inject()  │           │ inject()  │           │
│    └───────────┘           └───────────┘           └───────────┘           │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 为什么使用 provide/inject 而不是 props

| 方案                 | 优点                            | 缺点                                   |
| -------------------- | ------------------------------- | -------------------------------------- |
| **props 传递**       | 显式依赖，类型检查友好          | 参数爆炸（100+ 个状态/方法），维护困难 |
| **Pinia store**      | 跨组件共享自然                  | 状态迁移成本高，容易改变行为           |
| **provide/inject** ✓ | 避免 props 爆炸，保持单一数据源 | 隐式依赖，需要类型定义                 |

### 文件职责说明

```
frontend/src/components/
├── Keytone_album.vue              # 原始大组件（已改造为父组件）
│                                  # 职责：持有所有状态、副作用、provide Context
│                                  # 状态：✅ 已添加 provide，Step1/2/3 已替换
│
└── keytone-album/                 # 新的模块目录
    ├── index.ts                   # 模块入口 ✅
    │                              # 职责：统一导出，保持向后兼容
    │
    ├── types.ts                   # 类型定义 ✅
    │                              # 职责：定义 Context 接口、数据结构、注入 Key
    │
    ├── steps/                     # Step 子组件目录
    │   ├── StepLoadAudioFiles.vue     # Step 1: 加载音频源文件 ✅ 已集成
    │   ├── StepDefineSounds.vue       # Step 2: 定义声音 ✅ 已集成
    │   ├── StepCraftKeySounds.vue     # Step 3: 制作按键音 ✅ 已集成
    │   └── StepLinkageEffects.vue     # Step 4: 联动声效 ✅ 已集成
    │   # 职责：只负责 UI 渲染，通过 inject 获取状态
    │
    ├── dialogs/                   # Dialog 子组件目录
    │   ├── AddAudioFileDialog.vue     # 添加音频文件对话框 ✅ 已集成
    │   ├── ManageAudioFilesDialog.vue # 管理音频文件对话框 ✅ 已集成
    │   ├── CreateSoundDialog.vue      # 创建声音对话框 ✅ 已集成
    │   ├── EditSoundDialog.vue        # 编辑声音对话框 ✅ 已集成
    │   ├── CreateKeySoundDialog.vue   # 创建按键音对话框 ✅ 已集成
    │   └── EditKeySoundDialog.vue     # 编辑按键音对话框 ✅ 已集成
    │   # 职责：可复用的对话框 UI，可被任意 Step 调用
    │
    └── composables/               # 可复用逻辑 [Phase 4，可选]
        # 职责：抽离纯逻辑（如 SSE 映射、排序、校验）
```

### 数据流图

```
用户操作
    │
    ▼
┌───────────────────┐
│ Step/Dialog 组件   │  ←── inject(ctx) 获取状态和方法
│ (UI 层)           │
└─────────┬─────────┘
          │ 调用 ctx.xxx() 方法
          ▼
┌───────────────────┐
│ 父组件            │  ←── 状态更新
│ (数据层)          │
└─────────┬─────────┘
          │ watch 触发
          ▼
┌───────────────────┐
│ ConfigSet/Delete  │  ←── 写入配置文件
│ (持久化层)        │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ SDK 后端          │  ←── SSE 推送更新
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ 父组件            │  ←── messageAudioPackage 事件监听
│ 更新状态          │
└─────────┬─────────┘
          │ 响应式传播
          ▼
┌───────────────────┐
│ 子组件 UI 自动更新 │
└───────────────────┘
```

---

## Decision: 采用"父组件持有状态 + 子组件承载 UI"的拆分方式

### 选择

- 将现有 `Keytone_album.vue` 的 **状态（refs/computed/watch）与核心副作用（SSE 监听、ConfigSet/ConfigDelete、排序、依赖校验）** 保持在父层。
- 将每个 `q-step` 的模板与局部交互拆到单独的 step 组件。
- 将对话框从 step 内抽离为 `dialogs/` 下的组件，使其可被任意步骤直接拉起。

### 关键原因

- **最小风险**：保留数据流与副作用在一处，减少"拆分导致 reactivity/副作用丢失"的概率。
- **减少 prop 爆炸**：如果完全靠 props 传递大量 refs/方法，会导致 step 组件 API 复杂。

## Decision: 上下文注入（provide/inject）用于共享 ref 与 action ✅ 已实现

为避免 step 组件 props 过长，父组件提供一个 `KeytoneAlbumContext`：

- 状态：`step`、`pkgName`、各列表与选择项 refs、各种 dialog 的 `v-model` 状态等
- 行为：`saveSoundConfig`、`deleteSound`、`previewSound`、`deleteKeySound`、配置写回等

Step 与 Dialog 组件通过 `inject` 获取上下文，达到：

- 子组件更薄，只关心 UI 逻辑
- Dialog 可在任意 step 触发（满足"在某步骤中拉起其他步骤配置对话框"的诉求）

> 这仍然保持"单一数据源"在父组件，不引入新 store。

### 实现细节

**代码位置**: `Keytone_album.vue` script 结尾处

```typescript
// 导入 Context 类型和注入 Key
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from './keytone-album/types';

// 构建 Context 对象（包含所有需要共享的状态和方法）
const keytoneAlbumContext: KeytoneAlbumContext = {
  // Props
  pkgPath: props.pkgPath,
  isCreate: props.isCreate,
  
  // 核心状态
  step, pkgName,
  
  // Step1 相关
  addNewSoundFile, files, editSoundFile, soundFileList, selectedSoundFile,
  
  // Step2 相关
  createNewSound, soundName, sourceFileForSound, soundStartTime, soundEndTime, soundVolume,
  showEditSoundDialog, soundList, selectedSound,
  
  // Step3 相关
  createNewKeySound, keySoundName, /* ... */
  
  // Step4 相关
  isEnableEmbeddedTestSound, showEveryKeyEffectDialog, /* ... */
  
  // 操作函数
  saveSoundConfig, deleteSound, previewSound, saveKeySoundConfig, deleteKeySound,
  saveUnifiedSoundEffectConfig, saveSingleKeySoundEffectConfig,
  
  // 工具函数
  $t, album_options_select_label, naturalSort, /* ... */
};

// 提供 Context 给子组件
provide(KEYTONE_ALBUM_CONTEXT_KEY, keytoneAlbumContext);
```

---

## 当前实现状态

### 已完成的文件

| 文件                                               | 说明                            | 状态       |
| -------------------------------------------------- | ------------------------------- | ---------- |
| `Keytone_album.vue`                                | 父组件（已添加 provide）        | ✅ 已完成   |
| `keytone-album/types.ts`                           | 定义 Context 接口和所有数据类型 | ✅ 已完成   |
| `keytone-album/index.ts`                           | 模块入口，统一导出              | ✅ 已完成   |
| `keytone-album/steps/StepLoadAudioFiles.vue`       | Step 1: 加载音频源文件          | ✅ 已集成   |
| `keytone-album/steps/StepDefineSounds.vue`         | Step 2: 定义声音                | ✅ 已集成   |
| `keytone-album/steps/StepCraftKeySounds.vue`       | Step 3: 制作按键音              | ✅ 已集成   |
| `keytone-album/steps/StepLinkageEffects.vue`       | Step 4: 联动声效（框架）        | ✅ 框架完成 |
| `keytone-album/dialogs/AddAudioFileDialog.vue`     | 添加音频文件对话框              | ✅ 已集成   |
| `keytone-album/dialogs/ManageAudioFilesDialog.vue` | 管理音频文件对话框              | ✅ 已集成   |
| `keytone-album/dialogs/CreateSoundDialog.vue`      | 创建声音对话框                  | ✅ 已集成   |
| `keytone-album/dialogs/EditSoundDialog.vue`        | 编辑声音对话框                  | ✅ 已集成   |
| `keytone-album/dialogs/CreateKeySoundDialog.vue`   | 创建按键音对话框                | ✅ 已集成   |
| `keytone-album/dialogs/EditKeySoundDialog.vue`     | 编辑按键音对话框                | ✅ 已集成   |
| `keytone-album/dialogs/EveryKeyEffectDialog.vue`   | 全键声效对话框                  | ✅ 已集成   |
| `keytone-album/dialogs/SingleKeyEffectDialog.vue`  | 单键声效对话框（主对话框）      | ✅ 已集成   |
| `keytone-album/dialogs/AddSingleKeyEffectSubDialog.vue`  | 单键声效：添加子对话框     | ✅ 已集成   |
| `keytone-album/dialogs/EditSingleKeyEffectSubDialog.vue` | 单键声效：编辑子对话框     | ✅ 已集成   |

### 待完成（可选）

- [ ] 将 Step4 替换为子组件（复杂度高，包含虚拟键盘）
- [x] 创建 `EveryKeyEffectDialog.vue` 对话框 ✅
- [x] 创建 `SingleKeyEffectDialog.vue` 单键声效对话框 ✅
- [ ] 抽离 composables（SSE 映射、排序、校验逻辑）

---

## Alternatives considered

1) **全 props 传递**
- 优点：显式依赖
- 缺点：参数爆炸、维护成本高

2) **抽到 Pinia store**
- 优点：跨组件共享更自然
- 缺点：状态迁移更大、更容易引入行为变化；不符合"最小改动"

3) **一次性把所有逻辑拆成 composables**
- 优点：最终结构更干净
- 缺点：迁移风险高；因此采用渐进式（先拆模板/对话框，再逐步抽逻辑）

## Risks / Trade-offs

- **风险：副作用触发时机变化**（watch、debounce、SSE listener）
  - 缓解：副作用与数据写回仍保留在父组件；步骤拆分仅迁移模板与交互入口，避免触发时机漂移。
- **风险：Dialog 的 v-model 绑定与关闭时机**
  - 缓解：每个 dialog 单独抽取时保持 `v-model` 字段与事件处理不变。
- **风险：scoped style / :deep 选择器作用域变化**
  - 缓解：先保留样式在父组件；必要时将通用样式移动到更合适的位置，但需逐步验证。

## 清理状态（与实现同步）

- ✅ 已移除父组件中旧的单键声效对话框实现（原先用于对照的 `v-if="false"` 模板块），当前仅保留 `SingleKeyEffectDialog` 组件入口。

## Phase 4（composables）落地情况

- ✅ SSE 映射逻辑已从父组件抽离至 `keytone-album/composables/useKeytoneAlbumSseSync.ts`，父组件仅负责 attach/detach 与状态承载。
- ✅ 列表映射/自然排序的纯工具已抽离至 `keytone-album/composables/keytoneAlbumMappers.ts`（避免在多个位置重复实现同一映射逻辑）。
- ✅ 依赖校验 computed/watch 已抽离至 `keytone-album/composables/useKeytoneAlbumDependencyIssues.ts`，父组件继续暴露 `dependencyIssues` 给 UI 使用。
- ✅ Phase 4.5：`initData()` 与 `watch(audioFiles)` 中的映射逻辑已改用 `keytoneAlbumMappers` 纯函数，消除重复代码并确保"初始化 / watch / SSE"三条路径的映射行为一致。

---

## 架构导览索引（快速定位入口）

> 本节是一份"速查清单"，帮助你在 review 或排错时快速找到对应文件。

### 1. 调用关系一览

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│                              Keytone_album.vue（父组件 / 薄壳）                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────────────┐│
│  │ 持有所有状态 (refs/reactive)                                                                 ││
│  │ provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)                                                      ││
│  │ 生命周期：onBeforeMount → attach SSE；onUnmounted → detach SSE                               ││
│  └─────────────────────────────────────────────────────────────────────────────────────────────┘│
│        │                         │                         │                                    │
│        │ 调用                    │ 调用                    │ 调用                               │
│        ▼                         ▼                         ▼                                    │
│ useKeytoneAlbumSseSync   useKeytoneAlbumDependencyIssues   keytoneAlbumMappers                  │
│ (SSE监听+写入)            (依赖校验computed+watch)          (纯函数映射/排序工具)                  │
│        │                         │                                                              │
│        │ 使用                    │                                                              │
│        └────────────────────────►│                                                              │
│                                  │                                                              │
└─────────────────────────────────────────────────────────────────────────────────────────────────┘
                         ▲                               ▲
                         │ inject(ctx)                   │ inject(ctx)
    ┌────────────────────┴────────────────────┐          │
    │ steps/*.vue                              │          │
    │ - StepLoadAudioFiles.vue (Step 1)        │          │
    │ - StepDefineSounds.vue (Step 2)          │          │
    │ - StepCraftKeySounds.vue (Step 3)        │          │
    │ - StepLinkageEffects.vue (Step 4)        │          │
    └────────────────────┬────────────────────┘          │
                         │ 嵌套/打开                      │
                         ▼                               │
    ┌─────────────────────────────────────────────────────┘
    │ dialogs/*.vue
    │ - AddAudioFileDialog / ManageAudioFilesDialog (Step1)
    │ - CreateSoundDialog / EditSoundDialog (Step2)
    │ - CreateKeySoundDialog / EditKeySoundDialog (Step3)
    │ - EveryKeyEffectDialog / SingleKeyEffectDialog (Step4)
    └──────────────────────────────────────────────────────
```

### 2. 数据流简表

| 事件/触发点                              | 数据流向                                                             | 关键文件                             |
| ---------------------------------------- | -------------------------------------------------------------------- | ------------------------------------ |
| 后端 SSE `messageAudioPackage`           | SDK → eventSource → useKeytoneAlbumSseSync → 父组件 refs → 子组件 UI | `useKeytoneAlbumSseSync.ts`          |
| 用户操作（保存/删除）                    | Dialog/Step → ctx.xxx() → 父组件 watch → ConfigSet/ConfigDelete      | `Keytone_album.vue`                  |
| 列表映射（audio_files→soundFileList 等） | 父组件/SSE 回调 → keytoneAlbumMappers → 父组件 refs                  | `keytoneAlbumMappers.ts`             |
| 依赖校验变化                             | soundFileList/soundList/keySoundList → computed → dependencyIssues   | `useKeytoneAlbumDependencyIssues.ts` |

### 3. 排错入口速查

| 现象                  | 首选排查文件                                     | 次选排查文件                         |
| --------------------- | ------------------------------------------------ | ------------------------------------ |
| SSE 推送后 UI 未刷新  | `useKeytoneAlbumSseSync.ts`                      | `Keytone_album.vue` 的 attach/detach |
| 列表排序不对          | `keytoneAlbumMappers.ts` 的 naturalSort 调用     | 父组件的 `naturalSort` 实现          |
| 全键/单键联动映射不对 | 父组件的 `convertValue` + soundList/keySoundList | `useKeytoneAlbumSseSync.ts`          |
| 依赖警告不显示或误报  | `useKeytoneAlbumDependencyIssues.ts`             | `src/utils/dependencyValidator.ts`   |
| Dialog 打开/关闭异常  | 对应 `dialogs/*.vue` + `StepLinkageEffects.vue`  | 父组件的 v-model 状态（Context 字段） |
| Step 折叠/展开异常    | `Keytone_album.vue` 的 step 状态 + q-stepper     | —                                    |

### 4. 文件职责一句话总结

| 文件                                                           | 一句话职责                                                        |
| -------------------------------------------------------------- | ----------------------------------------------------------------- |
| `Keytone_album.vue`                                            | 父组件：持有所有状态、provide Context、生命周期挂载 SSE           |
| `keytone-album/types.ts`                                       | 类型：定义 Context 接口、数据结构、注入 Key                       |
| `keytone-album/index.ts`                                       | 入口：统一导出，保持向后兼容                                      |
| `keytone-album/steps/*.vue`                                    | Step UI：只负责渲染，通过 inject 获取状态和方法                   |
| `keytone-album/dialogs/*.vue`                                  | Dialog UI：可被任意 Step 拉起的对话框，通过 inject 获取状态和方法 |
| `keytone-album/composables/useKeytoneAlbumSseSync.ts`          | SSE：监听 messageAudioPackage、解析 JSON、写入父组件 refs         |
| `keytone-album/composables/keytoneAlbumMappers.ts`             | 映射：纯函数，Object.entries + 自然排序，供 SSE/初始化复用        |
| `keytone-album/composables/useKeytoneAlbumDependencyIssues.ts` | 校验：computed + watch，输出 dependencyIssues 供 UI 展示          |

---

## 代码注释规范（Code Comment Standards）

> **重要**: 本规范是本次重构的强制要求。详细的注释对于后续代码审查和维护至关重要。

### 注释保留原则

1. **不删除有效注释**：除非注释内容已过时或与当前逻辑不符，否则不得删除原有注释。
2. **迁移时保留**：将代码从父组件迁移到子组件时，必须保留原有的注释说明。
3. **翻译优化**：如果原注释是中文，迁移后保持中文；如果是英文，保持英文。不做无意义的语言转换。

### 必须添加注释的位置

#### 1. 文件头部注释（Vue SFC 文件）
每个新建的 `.vue` 文件必须在 `<template>` 标签之前包含完整的文件头注释：

```vue
<!--
 * This file is part of the KeyTone project.
 * Copyright (C) 2024 LuSrackhall
 * License: GPL-3.0
-->

<!--
============================================================================
文件说明: [目录]/[文件名] - [简短描述]
============================================================================

【文件作用】
[详细描述此文件的职责和功能]

【在整体架构中的位置】
[ASCII 架构图或文字说明]

【数据流】
[说明组件的数据来源和流向]

【关联文件】
- [相关文件1]: [关系说明]
- [相关文件2]: [关系说明]

【当前状态】
[✅ 已集成 / ⚠️ 待集成 / 🔄 开发中]
============================================================================
-->
```

#### 2. Script 区块注释
`<script setup>` 开头必须有组件职责说明：

```typescript
/**
 * [组件名].vue - [简短描述]
 *
 * 【组件职责】
 * - [职责1]
 * - [职责2]
 *
 * 【数据流向】
 * [说明数据如何流入流出此组件]
 *
 * 【注意事项】
 * [如有特殊注意事项，在此说明]
 */
```

#### 3. 代码分区注释
使用分隔线注释对代码进行逻辑分区：

```typescript
// ============================================================================
// [区块名称，如：注入父组件提供的上下文]
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// DOM 引用
// ============================================================================
const downSoundSelectDom = ref<any>(null);

// ============================================================================
// 计算属性
// ============================================================================

// ============================================================================
// 事件处理函数
// ============================================================================
```

#### 4. 函数注释
每个函数（包括事件处理函数）必须有 JSDoc 风格注释：

```typescript
/**
 * 保存按键音配置
 *
 * 【重要】
 * selectedKeySound.keySoundValue.down/up 的结构已在 watch 中被转换，
 * 其 mode 变成了 { mode: string } 对象形式。
 *
 * @param onSuccess - 可选的成功回调
 */
function handleSave(onSuccess?: () => void) {
  // 实现...
}
```

#### 5. 模板内注释
对于复杂的模板结构，使用 HTML 注释说明：

```vue
<template>
  <!-- 主对话框：创建新按键音 -->
  <q-dialog ...>
    <!-- 对话框标题 -->
    <q-card-section>...</q-card-section>

    <!-- 按键音名称输入 -->
    <q-card-section>...</q-card-section>

    <!-- 配置按钮组 -->
    <div>
      <!-- 配置按下声音按钮 -->
      <q-btn ... />

      <!-- 配置按下声音子对话框 -->
      <q-dialog ...>...</q-dialog>
    </div>
  </q-dialog>
</template>
```

#### 6. 样式注释
`<style>` 区块中的样式规则需要注释说明用途：

```scss
<style lang="scss" scoped>
/**
 * [组件名] 组件样式
 *
 * 【样式说明】
 * [说明样式的整体策略]
 */

// 按钮样式 - 统一按钮外观
.q-btn {
  @apply text-xs;
}

// 选择器样式 - 处理溢出
:deep(.q-field__native) {
  @apply max-w-full overflow-auto;
}

// 按键音选择器专用样式 - 用于多选芯片选择框
.zl-ll {
  // ...
}
</style>
```

### 注释检查清单

在提交代码前，检查以下项目：

- [ ] 文件头部注释完整（GPL 声明 + 文件说明）
- [ ] Script 区块有组件职责说明
- [ ] 代码有逻辑分区注释
- [ ] 每个函数有 JSDoc 注释
- [ ] 复杂模板结构有 HTML 注释
- [ ] 样式规则有用途说明
- [ ] 原有的有效注释未被删除
- [ ] TODO 注释标明了后续计划

> 补充：Phase 4 新增的 composables/mappers 文件也必须包含“文件头部说明注释”（用途/边界/关联文件/调试入口），
> 以保证后续 review 与定位问题成本可控。

### 注释示例参考

可参考以下已完成的文件作为注释规范的示例：
- `keytone-album/steps/StepLoadAudioFiles.vue`
- `keytone-album/steps/StepDefineSounds.vue`
- `keytone-album/dialogs/AddAudioFileDialog.vue`
- `keytone-album/types.ts`

---

## Migration Plan（渐进拆分）

1. ✅ 创建新目录结构与上下文类型（不改变功能）。
2. ✅ 在父组件中添加 provide，导出 Context。
3. 🔄 逐个把 Step1~4 的模板迁移到 step 子组件：
   - 每迁移一个 step：编译/运行，手动点测关键路径。
4. 将可复用对话框从 step 中抽离到 dialogs（保持 v-model 与事件函数一致）。
5. 可选：把"纯函数/映射/校验"逐步迁移到 composables，父层只组装。

## Validation

- 静态检查：`npm -C frontend run lint`
- 构建验证：`npm -C frontend run build` ✅ 已通过
- 手动回归（最小集）：
  - Step1：上传/管理音频源文件
  - Step2：创建/编辑/删除/预览声音
  - Step3：创建/编辑/删除键音；依赖警告展示
  - Step4：切换 embedded test sound；打开全局设置对话框并选择声效
  - Step header 折叠/展开（step=99）与 Continue/Back 导航
