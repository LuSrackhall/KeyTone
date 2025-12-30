# 无需授权场景可选联系方式 - 任务清单

## 状态

**当前状态**: ✅ 已完成

## 任务列表

### Phase 1: 组件实现 ✅

- [x] 1.1 创建 `OptionalContactDialog.vue` 组件
  - [x] Props: `visible`
  - [x] Emits: `submit`, `skip`, `cancel`
  - [x] 邮箱输入（可选，如填写则校验格式）
  - [x] 附加联系方式输入（可选）
  - [x] 跳过按钮
  - [x] 保存并继续按钮

### Phase 2: 状态机扩展 ✅

- [x] 2.1 新增状态 `optional-contact`
- [x] 2.2 新增 `optionalContactDialogVisible` ref
- [x] 2.3 新增处理函数
  - [x] `handleOptionalContactSubmit`
  - [x] `handleOptionalContactSkip`
  - [x] `handleOptionalContactCancel`
- [x] 2.4 修改 `handleAuthRequirementSubmit` 流程
- [x] 2.5 更新 `handlePickerCancel` 重置对话框
- [x] 2.6 更新 `reset` 函数
- [x] 2.7 导出新增变量和函数

### Phase 3: 页面集成 ✅

- [x] 3.1 在 `Keytone_album_page_new.vue` 导入组件
- [x] 3.2 在模板中添加对话框使用
- [x] 3.3 绑定事件处理函数

### Phase 4: 国际化 ✅

- [x] 4.1 中文翻译
  - [x] `optionalContact.title`
  - [x] `optionalContact.description`
  - [x] `optionalContact.hint`
  - [x] `optionalContact.skip`
  - [x] `optionalContact.continue`
- [x] 4.2 英文翻译

### Phase 5: 规格文档 ✅

- [x] 5.1 更新 `openspec/specs/export-flow/spec.md`
- [x] 5.2 创建变更文档
  - [x] `proposal.md`
  - [x] `design.md`
  - [x] `tasks.md`

## 完成标准

1. ✅ 选择"无需授权"后显示可选联系方式对话框
2. ✅ 点击"跳过"直接进入签名选择
3. ✅ 填写联系方式后点击"保存并继续"，联系方式被保存
4. ✅ 邮箱格式校验（仅当填写时）
5. ✅ 代码无 TypeScript 错误
6. ✅ 规格文档与代码同步

## 测试要点

- [ ] 选择"无需授权"后显示可选联系方式对话框
- [ ] 点击"跳过"直接进入签名选择，联系方式不被保存
- [ ] 输入有效邮箱后点击"保存并继续"，联系方式被保存
- [ ] 输入无效邮箱格式时显示错误提示，禁用继续
- [ ] 仅输入附加联系方式（不输入邮箱）也可保存
- [ ] 联系方式最终随签名配置保存到专辑
- [ ] 导出的专辑签名信息中可查看联系方式
