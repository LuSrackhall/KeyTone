# 键音专辑组件拆分重构 - 任务清单

> 目标：拆分 `Keytone_album.vue`，不改变任何现有功能/行为。

---

## 当前进度概览

| Phase   | 状态   | 说明                                           |
| ------- | ------ | ---------------------------------------------- |
| Phase 0 | ✅ 完成 | 基线验证通过                                   |
| Phase 1 | ✅ 完成 | 目录结构和类型文件已创建                       |
| Phase 2 | ✅ 完成 | Step1/2/3 已替换为独立组件，Step4 保留在父组件 |
| Phase 3 | ✅ 完成 | 6 个核心 Dialog 组件已创建并集成               |
| Phase 4 | ⏸️ 暂停 | 可选，待后续                                   |
| Phase 5 | ✅ 完成 | 构建验证通过，等待手动回归测试                 |

---

## Phase 0: 准备与基线

- [x] 0.1 记录关键手动回归路径（对应 proposal 的 Validation 最小集）
- [x] 0.2 在拆分前确认 `npm -C frontend run lint` 与 `npm -C frontend run build` 可通过（作为基线）

## Phase 1: 目录结构与兼容入口

- [x] 1.1 新建目录 `frontend/src/components/keytone-album/`（steps/dialogs/composables）
- [x] 1.2 创建类型定义文件 `types.ts`
- [x] 1.3 创建模块入口 `index.ts`
- [x] 1.4 构建验证通过

## Phase 2: Step 组件渐进拆分（模板优先） ✅

### 2.0 父组件 provide 集成 ✅

- [x] 2.0.1 在 `Keytone_album.vue` 中导入 `provide` 和 Context 类型
- [x] 2.0.2 构建 `keytoneAlbumContext` 对象，包含所有状态和方法
- [x] 2.0.3 调用 `provide(KEYTONE_ALBUM_CONTEXT_KEY, keytoneAlbumContext)`
- [x] 2.0.4 构建验证通过

### 2.1 Step1 子组件 ✅

- [x] 创建 `steps/StepLoadAudioFiles.vue` 组件框架
- [x] 定义组件结构，通过 inject 获取 Context
- [x] 添加详细架构注释
- [x] 用子组件替换父组件中的原有 Step1 模板

### 2.2 Step2 子组件 ✅

- [x] 创建 `steps/StepDefineSounds.vue` 组件框架
- [x] 定义组件结构，通过 inject 获取 Context
- [x] 添加详细架构注释
- [x] 用子组件替换父组件中的原有 Step2 模板

### 2.3 Step3 子组件 ✅

- [x] 创建 `steps/StepCraftKeySounds.vue` 组件框架
- [x] 定义组件结构，通过 inject 获取 Context
- [x] 添加详细架构注释
- [x] 用子组件替换父组件中的原有 Step3 模板

### 2.4 Step4 子组件 (保留在父组件)

- [x] 创建 `steps/StepLinkageEffects.vue` 组件框架
- [x] 定义组件结构，通过 inject 获取 Context
- [x] 添加详细架构注释
- [ ] **暂不替换**: Step4 保留在父组件中（复杂度高，包含虚拟键盘+多个嵌套对话框，约1500行）

> **说明**: Step4 包含虚拟键盘、全键声效对话框、单键声效对话框等复杂交互。
> 为降低风险，暂时保留在父组件中，待后续单独迁移。

## Phase 3: Dialog 抽离（提升可复用） ✅

- [x] 3.1 将"新增音频源文件"对话框抽离 → `dialogs/AddAudioFileDialog.vue`
- [x] 3.2 将"管理音频源文件"对话框抽离 → `dialogs/ManageAudioFilesDialog.vue`
- [x] 3.3 将"创建声音"对话框抽离 → `dialogs/CreateSoundDialog.vue`
- [x] 3.4 将"编辑声音"对话框抽离 → `dialogs/EditSoundDialog.vue`
- [x] 3.5 将"创建按键音"对话框抽离 → `dialogs/CreateKeySoundDialog.vue`
- [x] 3.6 将"编辑按键音"对话框抽离 → `dialogs/EditKeySoundDialog.vue`
- [ ] 3.7 将"全键声效"对话框抽离 → `dialogs/EveryKeyEffectDialog.vue` (可选)
- [ ] 3.8 将"单键声效"对话框抽离 → `dialogs/SingleKeyEffectDialog.vue` (可选)

## Phase 4: 逻辑域 composables（可选、低风险逐步）

- [ ] 4.1 抽离 SSE 映射逻辑到 `composables/useKeytoneAlbumSseSync.ts`
- [ ] 4.2 抽离列表映射/自然排序逻辑
- [ ] 4.3 抽离依赖校验逻辑

## Phase 5: 清理与验证 ✅

- [x] 5.1 确保无重复事件监听、watch 行为未改变
- [x] 5.2 确保样式作用域不回归
- [x] 5.3 运行 `npm -C frontend run lint` (warnings only, no errors)
- [x] 5.4 运行 `npm -C frontend run build` ✅ 通过
- [ ] 5.5 完成最小手动回归路径 (待用户测试)

---

## 已创建/修改的文件清单

| 文件路径                                           | 说明                           | 状态 |
| -------------------------------------------------- | ------------------------------ | ---- |
| `Keytone_album.vue`                                | 父组件（已添加 provide）       | ✅    |
| `keytone-album/types.ts`                           | Context 接口和所有数据类型定义 | ✅    |
| `keytone-album/index.ts`                           | 模块入口，统一导出             | ✅    |
| `keytone-album/steps/StepLoadAudioFiles.vue`       | Step 1: 加载音频源文件         | ✅    |
| `keytone-album/steps/StepDefineSounds.vue`         | Step 2: 定义声音               | ✅    |
| `keytone-album/steps/StepCraftKeySounds.vue`       | Step 3: 制作按键音             | ✅    |
| `keytone-album/steps/StepLinkageEffects.vue`       | Step 4: 联动声效               | ✅    |
| `keytone-album/dialogs/AddAudioFileDialog.vue`     | 添加音频文件对话框             | ✅    |
| `keytone-album/dialogs/ManageAudioFilesDialog.vue` | 管理音频文件对话框             | ✅    |
| `keytone-album/dialogs/CreateSoundDialog.vue`      | 创建声音对话框                 | ✅    |
| `keytone-album/dialogs/EditSoundDialog.vue`        | 编辑声音对话框                 | ✅    |
| `keytone-album/dialogs/CreateKeySoundDialog.vue`   | 创建按键音对话框               | ✅    |
| `keytone-album/dialogs/EditKeySoundDialog.vue`     | 编辑按键音对话框               | ✅    |

---

## 下一步行动

**当前状态**: 
- ✅ Step1/2/3 已替换为独立组件（父组件减少约 1750 行）
- ✅ Step4 保留原始实现（复杂度高，包含虚拟键盘，约 1500 行）
- ✅ 6 个核心 Dialog 组件已创建并集成
- ✅ 构建验证通过（无错误，仅有 warnings）
- ✅ 代码注释规范已添加到 design.md

**推荐下一步**: 手动功能回归测试

按照 Validation 最小集进行测试：
1. Step1：上传/管理音频源文件
2. Step2：创建/编辑/删除/预览声音
3. Step3：创建/编辑/删除键音；依赖警告展示
4. Step4：切换 embedded test sound；打开全局设置对话框并选择声效

**可选后续工作**:
- 将 Step4 的复杂内容（虚拟键盘、EveryKeyEffect 对话框、SingleKeyEffect 对话框）拆分为独立组件
- 抽离 composables（SSE 映射、排序、校验逻辑）

---

## 变更历史

| 日期       | 变更内容                                       |
| ---------- | ---------------------------------------------- |
| 2024-12-31 | 添加代码注释规范到 design.md                   |
| 2024-12-31 | 修复 EditKeySoundDialog 子对话框缺失问题       |
| 2024-12-31 | 补充完整注释到所有 Dialog 组件                 |
| 2024-12-30 | **模板替换完成**: Step1/2/3 替换为独立组件     |
| 2024-12-30 | 创建 4 个 Dialog 组件（声音/按键音的创建编辑） |
| 2024-12-30 | 创建 Step2/3/4 组件框架，构建验证通过          |
| 2024-12-30 | 添加 provide 到父组件 `Keytone_album.vue`      |
| 2024-12-30 | 创建 Step1 组件和两个 Dialog 组件              |
| 2024-12-30 | 创建目录结构和基础文件                         |
| 2024-12-30 | 初始化任务清单                                 |
