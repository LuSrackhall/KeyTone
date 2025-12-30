# 键音专辑组件拆分重构 - 任务清单

> 目标：拆分 `Keytone_album.vue`，不改变任何现有功能/行为。

---

## 当前进度概览

| Phase   | 状态     | 说明                     |
| ------- | -------- | ------------------------ |
| Phase 0 | ✅ 完成   | 基线验证通过             |
| Phase 1 | ✅ 完成   | 目录结构和类型文件已创建 |
| Phase 2 | 🔄 进行中 | Step1 框架完成，待集成   |
| Phase 3 | 🔄 进行中 | 部分 Dialog 已创建       |
| Phase 4 | ⏸️ 暂停   | 可选，待后续             |
| Phase 5 | ⏳ 待开始 | 验证                     |

---

## Phase 0: 准备与基线

- [x] 0.1 记录关键手动回归路径（对应 proposal 的 Validation 最小集）
- [x] 0.2 在拆分前确认 `npm -C frontend run lint` 与 `npm -C frontend run build` 可通过（作为基线）

## Phase 1: 目录结构与兼容入口

- [x] 1.1 新建目录 `frontend/src/components/keytone-album/`（steps/dialogs/composables）
- [x] 1.2 创建类型定义文件 `types.ts`
  - [x] 定义 `KeytoneAlbumContext` 接口
  - [x] 定义 `KEYTONE_ALBUM_CONTEXT_KEY` 注入 key
  - [x] 定义基础数据类型 (SoundFileInfo, SoundEntry, KeySoundEntry 等)
  - [x] 定义操作参数类型 (SaveSoundConfigParams 等)
- [x] 1.3 创建模块入口 `index.ts`
  - [x] 导出类型定义
  - [x] 导出主组件（当前指向旧组件）
- [x] 1.4 构建验证通过

### 已创建的文件清单

| 文件路径                                           | 说明                           | 状态 |
| -------------------------------------------------- | ------------------------------ | ---- |
| `keytone-album/types.ts`                           | Context 接口和所有数据类型定义 | ✅    |
| `keytone-album/index.ts`                           | 模块入口，统一导出             | ✅    |
| `keytone-album/steps/StepLoadAudioFiles.vue`       | Step 1 UI 组件框架             | ✅    |
| `keytone-album/dialogs/AddAudioFileDialog.vue`     | 添加音频文件对话框             | ✅    |
| `keytone-album/dialogs/ManageAudioFilesDialog.vue` | 管理音频文件对话框             | ✅    |

## Phase 2: Step 组件渐进拆分（模板优先）

### 总体策略

```
1. 父组件 (Keytone_album.vue) 添加 provide() 提供 Context
2. 逐个创建 Step 组件，通过 inject() 获取 Context
3. 用新 Step 组件替换原有模板
4. 每完成一个 Step 就验证构建和功能
```

### 任务列表

- [x] 2.1 创建 Step1（加载音频源文件）框架 `steps/StepLoadAudioFiles.vue`
  - [x] 定义组件结构
  - [x] 通过 inject 获取 Context
  - [x] 添加详细架构注释
  - [ ] **待完成**: 在父组件中集成

- [ ] 2.2 创建 Step2（定义声音）框架 `steps/StepDefineSounds.vue`
  - [ ] 保持裁剪字段、校验提示、保存/删除/预览一致

- [ ] 2.3 创建 Step3（制作按键音）框架 `steps/StepCraftKeySounds.vue`
  - [ ] 保持配置 down/up、依赖警告 `DependencyWarning`、保存/删除一致

- [ ] 2.4 创建 Step4（联动声效）框架 `steps/StepLinkageEffects.vue`
  - [ ] 保持 embedded test sound、全局设置对话框、禁用/完成条件一致

## Phase 3: Dialog 抽离（提升可复用）

- [x] 3.1 将"新增音频源文件"对话框抽离 → `dialogs/AddAudioFileDialog.vue`
- [x] 3.2 将"管理音频源文件"对话框抽离 → `dialogs/ManageAudioFilesDialog.vue`
- [ ] 3.3 将"创建/编辑声音"相关对话框抽离到 `dialogs/`
- [ ] 3.4 将"创建/编辑键音"相关对话框抽离到 `dialogs/`
- [ ] 3.5 将"全局联动设置/单键设置"等对话框抽离到 `dialogs/`

> 要求：dialog 组件 `v-model`、按钮行为、关闭行为、校验与通知保持一致。

## Phase 4: 逻辑域 composables（可选、低风险逐步）

- [ ] 4.1 抽离 SSE 映射逻辑到 `composables/useKeytoneAlbumSseSync.ts`
- [ ] 4.2 抽离列表映射/自然排序逻辑（audio_files/sounds/key_sounds）
- [ ] 4.3 抽离依赖校验逻辑（dependencyValidator）

## Phase 5: 清理与验证

- [ ] 5.1 确保无重复事件监听、watch 行为未改变
- [ ] 5.2 确保样式作用域不回归（尤其是 `:deep` 与滚动条样式）
- [ ] 5.3 运行 `npm -C frontend run lint`
- [ ] 5.4 运行 `npm -C frontend run build`
- [ ] 5.5 完成最小手动回归路径（Phase 0.1）

---

## 下一步行动

**当前阻塞点**: 在父组件 `Keytone_album.vue` 中集成新的子组件

**需要执行的操作**:
1. 在 `Keytone_album.vue` 中添加 `provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)`
2. 导入并使用 `StepLoadAudioFiles` 组件替换原有 Step1 模板
3. 验证功能正常后，继续创建 Step2/3/4

**风险提示**: 
- 修改 5k+ 行的主组件有一定风险
- 建议逐步替换，每替换一个 Step 就验证一次
