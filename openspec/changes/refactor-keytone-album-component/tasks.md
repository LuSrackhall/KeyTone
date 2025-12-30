# 键音专辑组件拆分重构 - 任务清单

> 目标：拆分 `Keytone_album.vue`，不改变任何现有功能/行为。

## Phase 0: 准备与基线

- [ ] 0.1 记录关键手动回归路径（对应 proposal 的 Validation 最小集）
- [ ] 0.2 在拆分前确认 `npm -C frontend run lint` 与 `npm -C frontend run build` 可通过（作为基线）

## Phase 1: 目录结构与兼容入口

- [ ] 1.1 新建目录 `frontend/src/components/keytone-album/`（steps/dialogs/composables/types）
- [ ] 1.2 保持现有引用路径兼容：
  - [ ] 方案A：保留 `frontend/src/components/Keytone_album.vue` 作为薄壳，内部渲染新组件
  - [ ] 或方案B：重导出新组件并确保所有引用点已替换（风险更高，默认不选）

## Phase 2: Step 组件渐进拆分（模板优先）

- [ ] 2.1 提取 Step1（加载音频源文件）到 `steps/StepLoadAudioFiles.vue`
  - [ ] 保持上传/管理逻辑、对话框开关字段、按钮点击行为一致
- [ ] 2.2 提取 Step2（定义声音）到 `steps/StepDefineSounds.vue`
  - [ ] 保持裁剪字段、校验提示、保存/删除/预览一致
- [ ] 2.3 提取 Step3（制作按键音）到 `steps/StepCraftKeySounds.vue`
  - [ ] 保持配置 down/up、依赖警告 `DependencyWarning`、保存/删除一致
- [ ] 2.4 提取 Step4（联动声效）到 `steps/StepLinkageEffects.vue`
  - [ ] 保持 embedded test sound、全局设置对话框、禁用/完成条件一致

## Phase 3: Dialog 抽离（提升可复用）

- [ ] 3.1 将“新增/管理音频源文件”相关对话框抽离到 `dialogs/`
- [ ] 3.2 将“创建/编辑声音”相关对话框抽离到 `dialogs/`
- [ ] 3.3 将“创建/编辑键音”相关对话框抽离到 `dialogs/`
- [ ] 3.4 将“全局联动设置/单键设置”等对话框抽离到 `dialogs/`

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
