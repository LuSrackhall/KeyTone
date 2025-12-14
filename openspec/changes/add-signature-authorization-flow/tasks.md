# 签名授权流程任务清单

## 概述

本任务清单用于跟踪签名授权申请与受理流程的实现进度。

## 任务列表

### Phase 1: SDK 核心功能

- [x] 1.1 创建 `sdk/signature/authorization.go` 文件
  - [x] 定义加密密钥常量（AuthKeyF, AuthKeyK, AuthKeyN, AuthKeyY）
  - [x] 实现 GenerateAuthRequest 函数
  - [x] 实现 ParseAuthRequest 函数
  - [x] 实现 GenerateAuthGrant 函数
  - [x] 实现 VerifyAndImportAuthGrant 函数

- [x] 1.2 更新 `sdk/audioPackage/config/signatureConfig.go`
  - [x] 实现 AddToAuthorizedList 函数

- [x] 1.3 验证 SDK 编译通过
  - [x] 运行 `go build ./...` 无错误

- [x] 1.4 修复加密逻辑缺陷
  - [x] 修正 `GenerateAuthGrant` 和 `VerifyAndImportAuthGrant` 中的非确定性加密问题，改为使用固定Nonce的确定性加密

- [x] 1.5 密钥安全管理
  - [x] 将 `sdk/signature/authorization.go` 中的密钥常量改为变量，支持 ldflags 注入
  - [x] 更新 `.gitignore` 忽略私有密钥配置文件
  - [x] 实现密钥混淆工具 `tools/key-obfuscator`
  - [x] 在 `sdk/signature/authorization.go` 中实现运行时 XOR 解混淆逻辑
  - [x] 创建 `sdk/private_keys.template.env` 模板文件
  - [x] 创建 `sdk/setup_build_env.sh` 自动化脚本，用于本地开发/构建时自动混淆并注入环境变量

### Phase 2: HTTP API 端点

- [x] 2.1 更新 `sdk/server/signature_handlers.go`
  - [x] 添加 POST /signature/generate-auth-request 端点
  - [x] 添加 POST /signature/parse-auth-request 端点
  - [x] 添加 POST /signature/generate-auth-grant 端点
  - [x] 添加 POST /signature/verify-import-auth-grant 端点
  - [x] 添加 POST /signature/add-to-authorized-list 端点

- [x] 2.2 验证 API 编译通过
  - [x] 运行 `go build ./...` 无错误

### Phase 3: 前端 API 函数

- [x] 3.1 更新 `frontend/src/boot/query/signature-query.ts`
  - [x] 添加 generateAuthRequest 函数
  - [x] 添加 parseAuthRequest 函数
  - [x] 添加 generateAuthGrant 函数
  - [x] 添加 verifyAndImportAuthGrant 函数
  - [x] 添加 addToAuthorizedList 函数

### Phase 4: 前端组件 - 授权申请

- [x] 4.1 创建 `AuthRequestDialog.vue`
  - [x] 定义 Props 和 Emits 接口
  - [x] 实现步骤 1：选择签名
  - [x] 实现步骤 2：查看联系方式 & 导出申请文件
  - [x] 实现步骤 3：完成提示
  - [x] 添加导出 .ktauthreq 文件功能
  - [x] 优化步骤条 UI（带图标、标题、连接线）
  - [x] 已选签名展示为完整卡片（图片+名称+介绍）
  - [x] 联系方式分开展示（邮箱+备用联系方式）
  - [x] 操作说明添加沟通提示
  - [x] 导出申请使用 File System Access API（showSaveFilePicker），用户确认保存后再提示成功，取消则不提示；不支持时回退下载链接

- [x] 4.2 更新 `SignaturePickerDialog.vue`
  - [x] 添加"授权申请"按钮
  - [x] 添加 @authRequest emit
  - [x] 调整按钮布局：三按钮同行（申请授权左、导入授权中、创建签名右）

- [x] 4.3 更新 `useExportSignatureFlow.ts`
  - [x] 添加 authRequestDialogVisible 状态
  - [x] 添加 openAuthRequestFromPicker 函数
  - [x] 添加 handleAuthRequestDone 函数
  - [x] 添加 handleAuthRequestCancel 函数

### Phase 5: 前端组件 - 授权受理

- [x] 5.1 创建 `AuthGrantDialog.vue`
  - [x] 定义 Props 和 Emits 接口
  - [x] 实现步骤 1：导入申请文件
  - [x] 实现步骤 2：审核 & 授权
  - [x] 实现步骤 3：完成提示
  - [x] 添加导出 .ktauth 文件功能
  - [x] 导出授权使用 File System Access API（showSaveFilePicker），用户确认保存后再提示成功，取消则不提示；不支持时回退下载链接

- [x] 5.2 更新 `Signature_management_page.vue`
  - [x] 添加"授权受理"按钮
  - [x] 集成 AuthGrantDialog 组件

### Phase 6: 前端组件 - 授权导入

- [x] 6.1 更新 `ExportAuthorizationGateDialog.vue`
  - [x] 完善授权文件导入逻辑
  - [x] 调用 verifyAndImportAuthGrant API
  - [x] 获取请求方资格码并调用 addToAuthorizedList API 更新授权列表
  - [x] 添加成功/失败提示
  - [x] 联系方式分开展示（邮箱+备用联系方式）

### Phase 7: 国际化

- [x] 7.1 更新 `en-US/index.json`
  - [x] 添加授权申请相关翻译键
  - [x] 添加授权受理相关翻译键
  - [x] 添加联系方式分块展示翻译键

- [x] 7.2 更新 `zh-CN/index.json`
  - [x] 添加授权申请相关翻译键
  - [x] 添加授权受理相关翻译键
  - [x] 添加联系方式分块展示翻译键

### Phase 8: 页面集成

- [x] 8.1 更新 `Keytone_album_page_new.vue`
  - [x] 导入 AuthRequestDialog 组件
  - [x] 添加 @authRequest 事件处理
  - [x] 传递必要的 props
  - [x] 修复 authRequestSignatures 计算属性，实现真实签名数据加载
  - [x] 添加 watch 监听器，在对话框打开时自动加载签名及授权状态
  - [x] 提取签名加载逻辑为独立函数 `loadAuthRequestSignatures()`
  - [x] 签名创建成功后自动刷新授权申请对话框的签名列表
  - [x] 移除重复的签名创建成功提示（SignatureFormDialog 内部已显示）

### Phase 9: 规格文档

- [x] 9.1 创建变更文档
  - [x] proposal.md - 变更提案
  - [x] design.md - 技术设计
  - [x] tasks.md - 任务清单

- [x] 9.2 创建规格增量文档
  - [x] specs/signature-management/spec.md - 签名管理增量
  - [x] specs/export-flow/spec.md - 导出流程增量

### Phase 10: 测试验证

- [x] 10.1 编译验证
  - [x] SDK 编译通过
  - [x] 前端 TypeScript 无类型错误

- [ ] 10.2 端到端测试（待人工测试）
  - [ ] 授权申请流程测试
  - [ ] 授权受理流程测试
  - [ ] 授权导入流程测试
  - [ ] 完整流程端到端测试

## 回滚方案

如需回滚，可按以下步骤操作：

1. 删除新增文件：
   - `sdk/signature/authorization.go`
   - `frontend/src/components/export-flow/AuthRequestDialog.vue`
   - `frontend/src/components/signature/AuthGrantDialog.vue`

2. 恢复修改的文件：
   - `git checkout -- sdk/server/signature_handlers.go`
   - `git checkout -- sdk/audioPackage/config/signatureConfig.go`
   - `git checkout -- frontend/src/components/export-flow/SignaturePickerDialog.vue`
   - `git checkout -- frontend/src/components/export-flow/useExportSignatureFlow.ts`
   - `git checkout -- frontend/src/pages/Signature_management_page.vue`
   - `git checkout -- frontend/src/boot/query/signature-query.ts`
   - `git checkout -- frontend/src/i18n/en-US/index.json`
   - `git checkout -- frontend/src/i18n/zh-CN/index.json`

## 注意事项

1. 加密密钥必须保持固定，更改会导致现有授权文件失效
2. 调试日志包含敏感信息，生产环境需要控制日志级别
3. 文件扩展名 `.ktauthreq` 和 `.ktauth` 是自定义格式，需要前端正确处理
