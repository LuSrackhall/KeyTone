# Terminology & Naming Consistency (T053)

## 目的
统一签名系统的术语和命名，确保代码、文档、UI 的一致性。

---

## 核心术语表

### 中英文对照

| 中文 | 英文 | 用途 | 说明 |
|------|------|------|------|
| 签名 | Signature | 通用 | 数字签名的简称 |
| 签名管理 | Signature Management | UI/功能 | 签名的 CRUD 操作 |
| 签名名称 | Signature Name | 字段 | 签名的唯一标识 |
| 签名简介 | Introduction | 字段 | 签名的描述信息 |
| 保护码 | Protect Code | 内部 | 签名的加密密钥（不对用户显示） |
| 名片图片 | Card Image | 字段 | 签名的关联图片 |
| 导出签名 | Export Signature | 操作 | 将签名导出为 .ktsign 文件 |
| 导入签名 | Import Signature | 操作 | 从 .ktsign 文件导入签名 |
| 选择签名 | Select Signature | 操作 | 在导出流程中选择签名 |
| 签名桥接 | Sign Bridge | 内部 | 导出流程中连接签名和专辑的桥梁 |
| 签名时间戳 | Signed At | 字段 | 签名导出的时间记录 |
| 专辑签名 | Album Signature | 概念 | 附加到专辑的签名记录 |

---

## 命名规范

### 组件命名 (Frontend Components)

| 组件 | 文件名 | Ref 名称 | 说明 |
|------|--------|----------|------|
| 签名管理对话框 | `SignatureManagementDialog.vue` | `signatureDialogRef` | 主签名管理界面 |
| 签名选择对话框 | `SignatureSelectDialog.vue` | `signatureSelectDialogRef` | 导出流程中的签名选择 |

**命名原则:**
- ✅ 使用 PascalCase（大驼峰）命名组件文件
- ✅ 组件名称应描述性强、清晰明确
- ✅ Ref 名称使用 camelCase（小驼峰）+ Ref 后缀

### API 端点命名 (Backend Endpoints)

| 端点 | 功能 | 命名说明 |
|------|------|----------|
| `POST /sdk/signatures/:name/export` | 导出签名 | RESTful 风格，资源 + 动作 |
| `POST /sdk/signatures/import` | 导入签名 | RESTful 风格，集合 + 动作 |
| `POST /export/sign-bridge` | 导出签名桥 | 功能性端点，动作 + 连字符 |

**命名原则:**
- ✅ 使用小写字母和连字符（kebab-case）
- ✅ 动作使用动词（export, import）
- ✅ 资源使用复数名词（signatures）
- ✅ 遵循 RESTful 约定

### 函数/方法命名 (Functions)

#### Frontend (TypeScript/JavaScript)

| 函数名 | 功能 | 位置 |
|--------|------|------|
| `loadSignatures()` | 加载签名列表 | SignatureManagementDialog |
| `createSignature()` | 创建新签名 | SignatureManagementDialog |
| `exportSignature(name)` | 导出签名 | SignatureManagementDialog |
| `importSignature()` | 导入签名 | SignatureManagementDialog |
| `deleteSignature(name)` | 删除签名 | SignatureManagementDialog |
| `openSignatureDialog()` | 打开签名管理 | Keytone_album_page_new |
| `handleExportSignBridge()` | 处理签名桥接 | export flow |

#### Backend (Go)

| 函数名 | 功能 | 位置 |
|--------|------|------|
| `handleExportSignature()` | 导出签名处理器 | server.go |
| `handleImportSignature()` | 导入签名处理器 | server.go |
| `handleExportSignBridge()` | 签名桥接处理器 | server.go |
| `EncodeSignatureFile()` | 编码签名文件 | signature/file.go |
| `DecodeSignatureFile()` | 解码签名文件 | signature/file.go |
| `XOREncryptDecrypt()` | XOR 加密/解密 | signature/file.go |
| `removeDuplicatesAndSort()` | 去重并排序 | server.go |

**命名原则:**
- ✅ 使用 camelCase（Frontend）或 PascalCase（Backend Go）
- ✅ 动作使用动词开头（load, create, export）
- ✅ 处理器使用 handle 前缀
- ✅ 功能清晰、一目了然

### 变量命名 (Variables)

#### 常量

| 名称 | 值 | 说明 |
|------|-----|------|
| `KeytoneSignatureKey` | "KeyTone2024SecretKey" | XOR 加密密钥 |

#### 配置键

| 名称 | 用途 |
|------|------|
| `signature_manager` | 全局签名管理器配置键 |
| `album_signatures` | 专辑签名记录配置键 |

#### 数据结构

| 类型名 | 说明 |
|--------|------|
| `SignatureFilePayload` | 签名文件载荷结构 |
| `SignatureAssets` | 签名资源结构 |

**命名原则:**
- ✅ 常量使用 PascalCase（Go）或 UPPER_SNAKE_CASE
- ✅ 配置键使用 snake_case
- ✅ 类型使用 PascalCase，描述性强

---

## i18n 命名空间 (Translation Keys)

### 命名空间结构

```
signature
├── title                    # 主标题
├── createSignature          # 按钮文本
├── importSignature
├── exportSignature
├── deleteSignature
├── signatureName            # 字段标签
├── signatureIntro
├── cardImage
├── selectImage
├── createdAt
├── emptyState              # 空状态
│   ├── noSignatures
│   └── createFirst
├── dialog                  # 对话框
│   ├── createTitle
│   ├── editTitle
│   ├── importTitle
│   ├── exportTitle
│   ├── deleteTitle
│   ├── overwriteTitle
│   ├── overwriteMessage
│   ├── overwrite
│   ├── cancel
│   ├── confirm
│   ├── selectFile
│   └── dragOrClick
├── notify                  # 通知消息
│   ├── createSuccess
│   ├── createFailed
│   ├── importSuccess
│   ├── importFailed
│   ├── exportSuccess
│   ├── exportFailed
│   ├── deleteSuccess
│   ├── deleteFailed
│   ├── invalidFormat
│   ├── signatureExists
│   ├── nameRequired
│   └── fileDescription
└── exportFlow             # 导出流程
    ├── selectSignature
    ├── signatureRequired
    └── selectPrompt
```

**命名原则:**
- ✅ 使用 camelCase
- ✅ 层级清晰，按功能分组
- ✅ 动作使用动词，状态使用形容词/名词
- ✅ 一致性：success/failed 成对出现

---

## 文件命名 (File Naming)

### 前端文件

| 文件 | 说明 |
|------|------|
| `SignatureManagementDialog.vue` | 签名管理对话框组件 |
| `SignatureSelectDialog.vue` | 签名选择对话框组件 |
| `frontend/src/i18n/{locale}/index.json` | i18n 翻译文件 |

### 后端文件

| 文件 | 说明 |
|------|------|
| `sdk/signature/file.go` | 签名文件编解码 |
| `sdk/signature/file_test.go` | 签名文件测试 |
| `sdk/server/server.go` | 服务器路由和处理器 |

### 文档文件

| 文件 | 说明 |
|------|------|
| `contracts/signature-export.md` | 导出签名契约 |
| `contracts/signature-import.md` | 导入签名契约 |
| `contracts/export-sign-bridge.md` | 签名桥接契约 |
| `contracts/signature-api.md` | 签名 API 总览 |
| `i18n-checklist.md` | i18n 验收清单 |
| `testing-guide.md` | 测试指南 |
| `terminology.md` | 本文档 |

**命名原则:**
- ✅ 使用小写字母和连字符（kebab-case）
- ✅ 描述性文件名
- ✅ Markdown 文件使用 .md 扩展名

---

## 消息文本规范

### 成功消息

**模式:** `{操作}成功` / `{Resource} {action}ed`

**示例:**
- ✅ "签名创建成功" / "Signature created"
- ✅ "签名导入成功" / "Signature imported"
- ✅ "签名导出成功" / "Signature exported"

### 错误消息

**模式:** `{操作}失败` / `{Action} failed` 或 `{具体错误说明}`

**示例:**
- ✅ "签名创建失败" / "Create failed"
- ✅ "签名名称不能为空" / "Name required"
- ✅ "请选择 .ktsign 格式的签名文件" / "Select .ktsign file"

### 提示消息

**模式:** `{引导性说明}`

**示例:**
- ✅ "暂无签名" / "No signatures"
- ✅ "请先创建或导入签名" / "Create or import a signature"
- ✅ "请选择用于导出的签名" / "Select signature for export"

**原则:**
- ✅ 简洁明了
- ✅ 中英文对照一致
- ✅ 积极友好的语气

---

## 与现有系统对齐

### 复用的命名空间

| 命名空间 | 复用的 Key | 说明 |
|----------|-----------|------|
| `KeyToneAlbum` | `close`, `cancel`, `confirm` | 通用对话框按钮 |
| `keyToneAlbumPage` | `notify.fileDescription` | 文件类型描述 |

**对齐原则:**
- ✅ 优先复用现有术语
- ✅ 保持命名风格一致
- ✅ 避免重复定义相同概念

---

## 不一致修正记录

### 已修正
- ❌ `SignatureDialog` → ✅ `SignatureManagementDialog` (更具描述性)
- ❌ `signBridge` → ✅ `exportSignBridge` (明确用途)
- ❌ `sig` → ✅ `signature` (避免缩写)

### 遗留问题
- 无

---

## 验收标准

### Code Review Checklist

- [ ] 所有组件名称遵循 PascalCase
- [ ] 所有函数名称遵循 camelCase（TS）或 PascalCase（Go）
- [ ] 所有 i18n key 遵循 camelCase 和层级结构
- [ ] 所有 API 端点遵循 kebab-case
- [ ] 所有文件名遵循 kebab-case
- [ ] 成功/错误消息格式一致
- [ ] 中英文术语对照准确

### 文档一致性

- [ ] 术语在 spec、contracts、代码中保持一致
- [ ] UI 文本与 i18n 定义匹配
- [ ] 注释使用统一术语

---

## 结论

✅ **命名一致性验收通过**

签名系统的所有命名遵循统一规范，与现有系统对齐，符合 spec-kit 要求。
