# 键音专辑组件拆分重构提案（不改变功能）

## 概述

当前前端组件 `frontend/src/components/Keytone_album.vue` 体积过大（约 5k+ 行），导致：

- 新功能迭代与修复困难（理解成本高、容易引入回归）
- 模板、状态、业务逻辑、对话框互相耦合，难以复用
- 目录结构不清晰，边界不明确

本变更提案目标是在 **完全不影响现有功能与行为** 的前提下，将该组件按“步骤条”及对话框/逻辑域进行拆分，使目录结构清晰、组件可复用、后续迭代成本更低。

## 为什么要做

- `Keytone_album.vue` 同时承担：步骤 UI、对话框、数据映射（SSE → UI state）、写回配置（ConfigSet/ConfigDelete）、依赖校验提示（DependencyWarning）等职责。
- 当前步骤条包含至少 4 个主要步骤（加载音频源文件、定义声音、制作按键音、联动声效设置）。这些步骤的 UI 结构相对独立，适合拆分。

## 变更范围

### 做什么（What Changes）

- 将 `Keytone_album.vue` 以最小风险方式拆分为：
  - 4 个 step 子组件（与现有步骤对应），承载各步骤模板与局部交互
  - 若干可复用 dialog 子组件（例如“新增音频源文件”“编辑声音”“新增键音”“全局联动设置”等）
  - 以“领域划分”的 composables（可选、循序渐进）承载：SSE 映射、写回配置、排序、依赖检查等逻辑
- 建立新的目录结构以承载上述拆分产物，提升可维护性与边界清晰度。

### 不做什么（Non-goals）

- 不新增 UI/交互能力（不新增按钮、页面、动画、额外配置项）
- 不更改现有业务行为（SSE 监听、ConfigSet 写回、step=99 折叠策略、校验/提示文案等保持不变）
- 不调整 i18n key、不改动后端 SDK API
- 不做设计重构（仅做“结构重构”）

## 代码可读性约束（Review-friendly）

- 本次重构新增的 step/dialog/composable/mapper 文件 SHALL 在文件头部提供清晰的说明注释：
  - 文件作用与边界（做什么/不做什么）
  - 与哪些文件/状态/事件配合使用
  - 行为不变的关键约束（避免 review 时误以为可随意“优化时序”）
  - Debug 定位入口（出问题优先看哪里）

## 现状要点（用于约束“不能变”）

以下要点在拆分后必须保持一致：

- 使用 `q-stepper` + `step=99` 的“折叠/展开”策略；点击 header 可在对应 step 与 99 间切换。
- Step1：音频文件上传/管理相关对话框与列表行为。
- Step2：声音裁剪定义、保存/删除/预览等行为。
- Step3：按键音（KeySound）制作与编辑，对依赖问题展示 `DependencyWarning`。
- Step4：联动声效（embedded test sound、全局/单键设置）及其对话框可直接拉起。
- `messageAudioPackage` SSE 事件驱动的配置映射流程、以及若干 watch → ConfigSet 的写回闭环。

## 建议的目标目录结构（提案）

> 目标是更清晰地管理“键音专辑编辑器”的功能边界；具体文件拆分以实现阶段为准。

```
frontend/src/components/
├── Keytone_album.vue              # 原始大组件 → 改造为父组件（薄壳）
│                                  # 职责：持有所有状态、副作用、provide Context
│
└── keytone-album/                 # 新的模块目录
    │
    ├── index.ts                   # 模块入口
    │                              # - 导出类型定义
    │                              # - 导出主组件（当前指向旧组件，迁移后切换）
    │                              # - 关联文件: types.ts, ../Keytone_album.vue
    │
    ├── types.ts                   # 类型定义文件
    │                              # - 基础数据类型 (SoundFileInfo, SoundEntry...)
    │                              # - 操作参数类型 (SaveSoundConfigParams...)
    │                              # - Context 接口 (KeytoneAlbumContext)
    │                              # - 注入 Key (KEYTONE_ALBUM_CONTEXT_KEY)
    │                              # - 关联文件: 被所有 step/dialog 组件引用
    │
    ├── steps/                     # Step 子组件目录
    │   │                          # 职责：只负责 UI 渲染，通过 inject 获取状态
    │   │
    │   ├── StepLoadAudioFiles.vue     # Step 1: 加载音频源文件
    │   │                              # - 显示添加/管理音频文件按钮
    │   │                              # - 内嵌 AddAudioFileDialog, ManageAudioFilesDialog
    │   │                              # - 关联文件: types.ts, ../dialogs/*
    │   │
    │   ├── StepDefineSounds.vue       # Step 2: 定义声音
    │   │                              # - 显示创建/编辑声音 UI
    │   │                              # - 调用 ctx.saveSoundConfig, ctx.previewSound
    │   │
    │   ├── StepCraftKeySounds.vue     # Step 3: 制作按键音
    │   │                              # - 显示创建/编辑按键音 UI
    │   │                              # - 显示依赖警告 (DependencyWarning)
    │   │
    │   └── StepLinkageEffects.vue     # Step 4: 联动声效
    │                                  # - 全局联动设置
    │                                  # - 单键联动设置
    │
    ├── dialogs/                   # Dialog 子组件目录
    │   │                          # 职责：可复用对话框 UI，可被任意 Step 调用
    │   │
    │   ├── AddAudioFileDialog.vue     # 添加音频文件对话框
    │   │                              # - 文件选择器 (支持拖拽)
    │   │                              # - 调用 SendFileToServer API
    │   │                              # - 关联文件: types.ts, keytonePkg-query.ts
    │   │
    │   ├── ManageAudioFilesDialog.vue # 管理音频文件对话框
    │   │                              # - 音频文件列表选择
    │   │                              # - 重命名/删除功能
    │   │                              # - 关联文件: types.ts, keytonePkg-query.ts
    │   │
    │   └── ... (其余对话框见 dialogs/ 目录)
    │
    └── composables/               # 可复用逻辑 [Phase 4，可选]
                                   # 职责：抽离纯逻辑（如 SSE 映射、排序、校验）
```

### 文件关系图

```
                    ┌─────────────────────┐
                    │  Keytone_album.vue  │  (父组件)
                    │  - 持有所有状态      │
                    │  - provide(ctx)     │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐
  │ types.ts    │◄─────│ index.ts    │      │ steps/*.vue │
  │ (类型定义)   │      │ (模块入口)   │      │ (UI组件)    │
  └──────┬──────┘      └─────────────┘      └──────┬──────┘
         │                                         │
         │                                         ▼
         │                                  ┌─────────────┐
         └─────────────────────────────────►│dialogs/*.vue│
                                            │ (对话框组件) │
                                            └─────────────┘
```

> 兼容性策略：对外保持 `frontend/src/components/Keytone_album.vue` 的现有路径可继续工作（可通过“包装导出”或保留文件并内部引用新组件实现）。

## 影响面

- 影响的代码主要在前端：组件拆分、文件移动/新增。
- 不涉及后端 SDK 行为变更。

## 验收标准

- 视觉与交互一致：各步骤 UI、按钮、对话框、提示、禁用条件与原行为一致。
- 数据一致：SSE 更新、列表排序、选择项值与之前一致；ConfigSet/ConfigDelete 写回路径与触发时机不变。
- 构建一致：`npm -C frontend run lint` 与 `npm -C frontend run build` 通过。

---

## 已落地实现（截至 2025-12-31）

- Step1/2/3/4 已替换为独立组件；父组件作为薄壳持有状态并 provide Context。
- Dialog 抽离：
  - `EveryKeyEffectDialog`（全键声效）已抽离并集成。
  - `SingleKeyEffectDialog`（单键声效）已抽离并集成（内部拆为“添加/编辑”两个子对话框组件）。
- 清理：已删除父组件中旧的单键声效对话框实现（原先用于对照的 `v-if="false"` 块），避免冗余与模板误改风险。
- Phase 4（composables）：
  - SSE 映射监听与数据写入已抽离为 `useKeytoneAlbumSseSync`。
  - 列表映射/自然排序纯工具已抽离为 `keytoneAlbumMappers`。
  - 依赖校验逻辑已抽离为 `useKeytoneAlbumDependencyIssues`。
  - Phase 4.5：`initData()` 与 `watch(audioFiles)` 中的映射逻辑已改用 `keytoneAlbumMappers` 纯函数，消除重复代码，确保"初始化 / watch / SSE"三条路径的映射行为一致。
- 构建验证：`npm -C frontend run build` 通过。
