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

- `frontend/src/components/keytone-album/`
  - `KeytoneAlbum.vue`（薄壳：持有状态/组合式逻辑/提供上下文；保持外部引用稳定）
  - `steps/`
    - `StepLoadAudioFiles.vue`
    - `StepDefineSounds.vue`
    - `StepCraftKeySounds.vue`
    - `StepLinkageEffects.vue`
  - `dialogs/`（从 step 中抽离出来，便于跨步骤复用/直接拉起）
  - `composables/`（可选，逐步抽离）
  - `types.ts`

> 兼容性策略：对外保持 `frontend/src/components/Keytone_album.vue` 的现有路径可继续工作（可通过“包装导出”或保留文件并内部引用新组件实现）。

## 影响面

- 影响的代码主要在前端：组件拆分、文件移动/新增。
- 不涉及后端 SDK 行为变更。

## 验收标准

- 视觉与交互一致：各步骤 UI、按钮、对话框、提示、禁用条件与原行为一致。
- 数据一致：SSE 更新、列表排序、选择项值与之前一致；ConfigSet/ConfigDelete 写回路径与触发时机不变。
- 构建一致：`npm -C frontend run lint` 与 `npm -C frontend run build` 通过。
