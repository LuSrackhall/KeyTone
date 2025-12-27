# 签名授权流程技术设计

## 概述

本文档描述签名授权申请与受理流程的技术实现细节。

## 架构设计

### 模块依赖关系

```text
┌─────────────────────────────────────────────────────────────────┐
│                         前端 (Vue/TypeScript)                     │
├─────────────────────────────────────────────────────────────────┤
│  SignaturePickerDialog  →  AuthRequestDialog                    │
│         ↓ authRequest                  ↓ export                 │
│                                   .ktauthreq 文件               │
├─────────────────────────────────────────────────────────────────┤
│  Signature_management_page  →  AuthGrantDialog                  │
│         ↓ authGrant                    ↓ import/export          │
│                              .ktauthreq → .ktauth               │
├─────────────────────────────────────────────────────────────────┤
│  ExportAuthorizationGateDialog                                  │
│         ↓ import                                                │
│                              .ktauth → authorizedList           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP API                                 │
├─────────────────────────────────────────────────────────────────┤
│  /signature/generate-auth-request                               │
│  /signature/parse-auth-request                                  │
│  /signature/generate-auth-grant                                 │
│  /signature/verify-import-auth-grant                            │
│  /signature/add-to-authorized-list                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         SDK (Go)                                 │
├─────────────────────────────────────────────────────────────────┤
│  signature/authorization.go                                     │
│    - GenerateAuthRequest()                                      │
│    - ParseAuthRequest()                                         │
│    - GenerateAuthGrant()                                        │
│    - VerifyAndImportAuthGrant()                                 │
├─────────────────────────────────────────────────────────────────┤
│  audioPackage/config/signatureConfig.go                         │
│    - AddToAuthorizedList()                                      │
└─────────────────────────────────────────────────────────────────┘
```

## 加密方案

### 密钥定义

| 密钥变量 | 值                              | 用途             |
| -------- | ------------------------------- | ---------------- |
| AuthKeyF | KT_AuthFlowFieldEncryptKey_F01  | 加密请求方签名ID |
| AuthKeyK | KT_AuthFlowFieldEncryptKey_K02  | 预留             |
| AuthKeyN | KT_AuthFlowFieldEncryptKey_N03  | 预留             |
| AuthKeyY | KT_AuthFlowCodePartialEncrypt04 | 加密资格码前15位 |

### 密钥管理与构建注入

- **混淆工具**：`tools/key-obfuscator/main.go`，用于离线生成混淆后的 Hex 值。
- **本地开发/构建**：
  - 配置文件：`sdk/private_keys.env`（从 template 复制，不提交 git）。
  - 自动化脚本：`sdk/setup_build_env.sh`。
    - 功能：读取私钥 -> 调用混淆工具 -> 导出 `EXTRA_LDFLAGS` 环境变量。
    - 用法：`source sdk/setup_build_env.sh`（仅当前终端生效）。
  - 构建支持：
    - `sdk/makefile`：已更新支持读取 `EXTRA_LDFLAGS`。
    - `frontend/src-electron/debug.sh`：已更新支持读取 `EXTRA_LDFLAGS`（用于 `quasar dev` 调试）。
- **GitHub Actions**：
  - 在 Secrets 中配置一个 `EXTRA_LDFLAGS` 变量，包含所有 `-X ...` 注入参数（使用混淆后的 Hex 值）。
  - Workflow 直接使用该 Secret，无需运行 setup 脚本。

## 前端交互细节

- 授权受理对话框（AuthGrantDialog）审核步骤：
  - “申请方”字段文案调整为“申请方名称”，对应请求文件中的 `requesterSignatureName`。
  - “您的签名”区域展示完整签名列表项（图片（若有）+名称+介绍）；若匹配多个则先选择再展示所选项。

### 授权申请文件生成流程

```text
输入：
  - authorizationUUID: 专辑的授权标识UUID
  - originalAuthorQualificationCode: 原始作者的资格码
  - requesterSignatureID: 请求方的签名ID

步骤：
1. authorizationUUIDHash = SHA256(authorizationUUID)[后10位]
2. encryptedOriginalAuthorCode = AES_GCM(originalAuthorQualificationCode[前15位], AuthKeyY)
3. requesterSignatureIDSuffix = AES_GCM(requesterSignatureID, AuthKeyF)

输出：
{
  "authorizationUUIDHash": "...",
  "encryptedOriginalAuthorCode": "...",
  "requesterSignatureIDSuffix": "..."
}
```

### 授权文件生成流程

```text
输入：
  - authorizationUUIDHash: 从申请文件获取
  - requesterSignatureIDSuffix: 从申请文件获取
  - originalAuthorSignatureID: 原始作者选择的签名ID

步骤：
1. requesterQualificationCode = SHA256(Decrypt(requesterSignatureIDSuffix, AuthKeyF))
2. originalAuthorQualificationCode = SHA256(originalAuthorSignatureID)
3. encryptedOriginalCode = AES_GCM_Deterministic(originalAuthorQualificationCode[前15位], AuthKeyY)
4. combinedString = authorizationUUIDHash + requesterQualificationCode[前11位] + encryptedOriginalCode
5. verificationCode = SHA256(combinedString)
6. encryptedAuthToken = AES_GCM(verificationCode, AuthKeyN)

输出：
{
  "encryptedAuthToken": "...",
  "version": "1.0"
}
```

### 授权验证流程

```text
输入：
  - authFile: 授权文件内容
  - authorizationUUID: 专辑的授权标识UUID
  - requesterEncryptedSignatureID: 请求方加密的签名ID
  - originalAuthorQualificationCode: 原始作者的资格码

步骤：
1. 验证 authorizationUUIDHash 与 SHA256(authorizationUUID)[后10位] 匹配
2. 重新计算 verificationCode 并比对
3. 返回验证结果
```

## 前端组件设计

### AuthRequestDialog

**Props**:
- `visible: boolean` - 对话框可见性
- `signatures: Signature[]` - 签名列表（含授权状态）
- `contactEmail: string` - 原始作者邮箱
- `contactAdditional?: string` - 原始作者备用联系方式
- `authorizationUUID: string` - 专辑授权UUID
- `originalAuthorQualificationCode: string` - 原始作者资格码

**数据加载**:
- 签名数据加载逻辑提取为独立函数 `loadAuthRequestSignatures()`
- 当对话框打开时（visible 变为 true），父组件通过 watch 监听自动调用此函数
- 函数内部调用：
  - `GetAvailableSignaturesForExport(albumPath)` 获取授权状态
  - `getSignaturesList()` 获取加密签名列表
  - `decryptSignatureData()` 解密签名数据
  - `getSignatureImage()` 获取签名图片
- 筛选逻辑：仅显示 `isAuthorized=false` 的签名
- 联系方式从专辑签名配置中获取 `authorization.contactEmail` 和 `authorization.contactAdditional`

**签名创建后刷新**:
- 当用户在授权申请对话框中点击"创建签名"并成功创建后
- `onSignatureFormSuccess` 回调检测 `authRequestDialogVisible` 状态
- 如果对话框打开，调用 `loadAuthRequestSignatures()` 刷新列表
- 注意：不再显示重复的成功提示（`SignatureFormDialog` 内部已显示）

**Emits**:
- `done` - 流程完成
- `cancel` - 用户取消
- `createSignature` - 创建新签名

**步骤**:
1. 选择一个未授权的签名
2. 查看原始作者联系方式（邮箱+备用，各带复制按钮），查看已选签名卡片，阅读操作说明（含沟通提示），导出授权申请文件
3. 完成提示

**导出申请文件保存流程**:
- 优先使用 `window.showSaveFilePicker` 弹出保存对话框，`suggestedName: auth-request-<timestamp>.ktauthreq`
- 写入完成后再进入成功步骤并提示成功；用户取消保存则不提示成功
- 兼容回退：若 API 不可用，则使用浏览器下载链接方案触发下载

**步骤条 UI 设计**:
- 使用自定义步骤条组件，每个步骤包含：
  - 圆形指示器（28x28px）：显示步骤编号或勾选图标
  - 步骤标题：显示在圆形下方
  - 连接线：连接相邻步骤
- 状态样式：
  - 未激活：灰色背景
  - 当前步骤：紫色背景 + 阴影
  - 已完成：绿色背景 + 勾选图标

**已选签名展示**:
- 使用卡片布局，包含：
  - 48x48px 图片（圆角6px）
  - 签名名称（加粗）
  - 签名介绍（最多2行，CSS 省略）

**联系方式展示**:
- 分两个区块：
  - 邮箱区块：邮件图标 + 标签 + 邮箱内容 + 复制按钮
  - 备用联系方式区块（可选）：联系人图标 + 标签 + 内容 + 复制按钮
  - 两区块之间使用分隔线

**操作说明**:
- 醒目提示条（amber 背景）：建议先与原始作者沟通
- 3个步骤说明，包含文件扩展名提示

### ExportAuthorizationGateDialog

**Props**:
- `visible: boolean` - 对话框可见性
- `contactEmail?: string` - 原始作者邮箱
- `contactAdditional?: string` - 原始作者备用联系方式
- `albumPath?: string` - 专辑路径
- `authorizationUUID?: string` - 授权UUID
- `requesterEncryptedSignatureID?: string` - 请求方签名ID
- `originalAuthorQualificationCode?: string` - 原始作者资格码

**联系方式展示**:
- 与 AuthRequestDialog 相同的分块展示方式

**授权文件导出保存流程**:
- 使用 `window.showSaveFilePicker` 弹出保存对话框，`suggestedName: auth-grant-<timestamp>.ktauth`
- 写入完成后再提示成功；若用户取消保存则不提示成功
- 兼容回退：若 API 不可用，回退为浏览器下载链接方案

### SignaturePickerDialog 按钮布局

**布局规则**:
- 当 `requireAuthorization=true` 时：
  - 使用 `justify-between` 布局
  - 左侧：申请授权按钮
  - 右侧容器：导入授权按钮 + 创建签名按钮
- 当 `requireAuthorization=false` 时：
  - 使用 `justify-end` 布局
  - 仅显示创建签名按钮

### AuthGrantDialog

**Props**:
- `visible: boolean` - 对话框可见性
- `localSignatures: Signature[]` - 本地已解密签名列表（用于展示完整签名列表项）
- `getImageUrl: (imagePath: string) => string` - 图片路径转预览 URL（复用签名管理页的缓存逻辑）

**Emits**:
- `done` - 流程完成
- `cancel` - 用户取消

**步骤**:
1. 导入 .ktauthreq 授权申请文件
2. 审核申请信息，选择签名授权
3. 导出 .ktauth 授权文件

## API 设计

### POST /signature/generate-auth-request

**请求**:

```json
{
  "authorizationUUID": "string",
  "originalAuthorQualificationCode": "string",
  "requesterSignatureID": "string",
  "requesterSignatureName": "string"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "fileContent": "base64编码的二进制内容"
  }
}
```

### POST /signature/parse-auth-request

**请求**:

```json
{
  "fileContent": [二进制数组],
  "localSignatureIDs": ["id1", "id2", ...]
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "authorizationUUIDHash": "string",
    "requesterSignatureIDSuffix": "string",
    "matchedLocalSignatureID": "string"
  }
}
```

### POST /signature/generate-auth-grant

**请求**:

```json
{
  "authorizationUUIDHash": "string",
  "requesterSignatureIDSuffix": "string",
  "originalAuthorSignatureID": "string"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "fileContent": "base64编码的二进制内容"
  }
}
```

### POST /signature/verify-import-auth-grant

**请求**:

```json
{
  "fileContent": [二进制数组],
  "authorizationUUID": "string",
  "requesterEncryptedSignatureID": "string", // 可选，若为空则后端自动遍历本地签名匹配
  "originalAuthorQualificationCode": "string"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "valid": true,
    "requesterQualificationCode": "string"
  }
}
```

### POST /signature/add-to-authorized-list

**请求**:

```json
{
  "albumPath": "string",
  "qualificationCode": "string"
}
```

**响应**:

```json
{
  "success": true,
  "message": "授权添加成功"
}
```

## 调试日志格式

```text
[授权申请生成] 开始
  - authorizationUUID: xxx
  - 变换: SHA256 → 取后10位
  - 结果: authorizationUUIDHash = xxx

  - originalAuthorQualificationCode: xxx
  - 取前15位: xxx
  - 变换: AES-GCM(密钥=AuthKeyY)
  - 结果: encryptedOriginalAuthorCode = xxx

  - requesterSignatureID: xxx
  - 变换: AES-GCM(密钥=AuthKeyF)
  - 结果: requesterSignatureIDSuffix = xxx

[授权申请生成] 完成
```

## 错误处理

| 错误场景         | HTTP状态码 | 错误信息                  |
| ---------------- | ---------- | ------------------------- |
| 缺少必填参数     | 400        | 缺少必填参数: xxx         |
| 解析申请文件失败 | 400        | 解析授权申请文件失败: xxx |
| 无匹配签名       | 400        | 未找到匹配的本地签名      |
| 验证失败         | 400        | 授权验证失败              |
| 内部错误         | 500        | 内部错误: xxx             |
