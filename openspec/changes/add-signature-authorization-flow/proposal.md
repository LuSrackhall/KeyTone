# 变更：实现签名授权申请与受理流程

## 为什么需要

在"专辑有签名+需要授权"的场景下，第三方用户希望使用自己的签名导出专辑时，需要获得原始作者的授权。本变更实现完整的授权申请和受理流程，使：

1. **请求方**可以选择一个未授权的签名，生成授权申请文件（.ktauthreq）发送给原始作者
2. **原始作者**可以导入授权申请文件，审核后生成授权文件（.ktauth）返回给请求方
3. **请求方**导入授权文件后，该签名被添加到专辑的授权列表中，可正常导出

## 变更内容

### 核心功能

#### SDK 侧（后端）

1. **授权申请文件生成** (`GenerateAuthRequest`)
   - 输入：专辑 authorizationUUID、原始作者资格码、请求方签名 ID
   - 输出：二进制授权申请文件内容
   - 加密字段：
     - `authorizationUUIDHash`: authorizationUUID 的 SHA256 后10位
     - `encryptedOriginalAuthorCode`: 原始作者资格码前15位用密钥y加密
     - `requesterSignatureIDSuffix`: 请求方签名 ID 用密钥f加密

2. **授权申请文件解析** (`ParseAuthRequest`)
   - 输入：授权申请文件内容、本地签名ID列表
   - 输出：解析后的申请信息（专辑哈希、请求方签名、匹配的本地签名）
   - 验证：确认当前用户拥有对应的原始作者签名

3. **授权文件生成** (`GenerateAuthGrant`)
   - 输入：解析后的申请信息
   - 输出：二进制授权文件内容
   - 最终验证码组成：
     1. authorizationUUID 的 SHA256 值的后10位
     2. 请求方签名资格码的前11位
     3. 加密后的资格码（前15位用密钥y进行确定性加密）
     4. 最终组合字符串再进行一次 SHA256 哈希
   - 注意：为了保证验证时的确定性，对原始作者资格码前15位的加密采用固定Nonce的AES-GCM模式。

4. **授权文件验证与导入** (`VerifyAndImportAuthGrant`)
   - 输入：授权文件内容、验证参数（请求方签名ID可选）
   - 输出：验证结果、请求方签名资格码
   - 功能：验证授权文件有效性；若未指定请求方签名ID，自动遍历本地签名寻找匹配项；返回请求方签名的资格码以便添加到授权列表

5. **添加到授权列表** (`AddToAuthorizedList`)
   - 输入：专辑路径、资格码
   - 功能：将资格码添加到专辑配置的 authorizedList

#### HTTP 端点

1. `POST /signature/generate-auth-request` - 生成授权申请文件
2. `POST /signature/parse-auth-request` - 解析授权申请文件
3. `POST /signature/generate-auth-grant` - 生成授权文件
4. `POST /signature/verify-import-auth-grant` - 验证并导入授权文件
5. `POST /signature/add-to-authorized-list` - 添加到授权列表

#### 前端侧

1. **SignaturePickerDialog 增强**
   - 新增"授权申请"按钮（仅在需要授权模式下显示）
   - 点击后进入授权申请流程

2. **AuthRequestDialog 组件**（新增）
   - 3步骤向导：选择签名 → 查看联系方式&导出申请 → 完成
   - 仅显示未授权的签名列表
   - 展示原始作者联系方式（复制功能）
   - 签名列表项名称/介绍支持横向滚动（滚动条样式与签名列表一致），避免对话框溢出
   - **展示选中签名的资格码指纹**（便于授权链路中的相关方检阅）
   - 导出 .ktauthreq 文件

3. **AuthGrantDialog 组件**（新增）
   - 3步骤向导：导入申请文件 → 审核&授权 → 完成
   - 解析并展示申请信息
   - 验证当前用户拥有对应签名
   - 审核步骤将"申请方"标签显示为"申请方名称"
   - **展示申请方的资格码指纹**（从申请文件中解析，便于受理方核实申请方身份）
   - "您的签名"区域展示完整签名列表项（图片（若有）+名称+介绍），而非仅名称
   - **展示匹配签名的资格码指纹**（便于授权链路中的相关方检阅）
   - 长文本（名称/介绍）支持横向滚动（滚动条样式与签名列表一致），避免对话框溢出
   - 导出 .ktauth 文件

4. **签名管理页面增强**
   - 新增"授权受理"按钮
   - 打开 AuthGrantDialog
   - **右键菜单增加"查看资格码指纹"选项**，支持查看任意签名的资格码指纹

5. **ExportAuthorizationGateDialog 增强**
   - 完善授权文件导入逻辑
   - 调用验证 API 并更新授权列表

### 资格码指纹展示

在授权链路中，资格码指纹用于相关方核实签名身份：

- **授权申请方**（AuthRequestDialog）：展示选中签名的资格码指纹
- **授权受理方**（AuthGrantDialog）：
  - 展示申请方的资格码指纹（从申请文件解析）
  - 展示匹配到的本地签名的资格码指纹
- **签名管理页面**：右键菜单支持查看任意签名的资格码指纹
- **指纹来源**：SDK 端动态计算，前端直接使用返回的 `qualificationFingerprint` 字段

### 加密密钥设计与混淆注入

- **代码实现**：密钥在代码中以 `var` 形式存放，运行时通过 XOR 解混淆；默认值为开源占位符。
- **构建注入**：通过 `-ldflags -X ...` 注入混淆后的 Hex 值。
- **本地工作流**：
  - 使用 `sdk/setup_build_env.sh` 脚本自动处理 `private_keys.env` 中的明文密钥，并在当前终端导出 `EXTRA_LDFLAGS`。
  - 支持 `make` 构建和 `quasar dev` 调试（通过 `debug.sh` 透传）。
- **CI 工作流**：
  - GitHub Actions 中直接配置 `EXTRA_LDFLAGS` Secret（预先混淆好的值），无需运行脚本。

### 数据格式

**授权申请文件 (.ktauthreq)**：

```json
{
  "authorizationUUIDHash": "<SHA256后10位>",
  "encryptedOriginalAuthorCode": "<密钥y加密的资格码前15位>",
  "requesterSignatureIDSuffix": "<密钥f加密的请求方签名ID>",
  "requesterSignatureName": "<请求方签名名称（可选，便于识别）>"
}
```

**授权文件 (.ktauth)**：

```json
{
  "authorizationUUIDHash": "<SHA256后10位>",
  "verificationCode": "<最终验证码的SHA256>",
  "grantedAt": "<授权时间戳>"
}
```

### 调试日志

所有加密操作都输出详细的调试日志：
- 原始值
- 变换方式（SHA256、对称加密等）
- 使用的密钥及其来源变量名
- 变换的最终结果

## 影响范围

### 新增能力

- 签名授权申请文件生成
- 签名授权文件生成
- 签名授权文件验证与导入

### 受影响代码

**新增**：
- `sdk/signature/authorization.go` - 授权流程核心逻辑
- `frontend/src/components/export-flow/AuthRequestDialog.vue` - 授权申请对话框
- `frontend/src/components/signature/AuthGrantDialog.vue` - 授权受理对话框

**修改**：
- `sdk/server/signature_handlers.go` - 新增 5 个 HTTP 端点
- `sdk/audioPackage/config/signatureConfig.go` - 新增 AddToAuthorizedList 函数
- `frontend/src/components/export-flow/SignaturePickerDialog.vue` - 新增授权申请按钮
- `frontend/src/components/export-flow/ExportAuthorizationGateDialog.vue` - 完善导入逻辑
- `frontend/src/components/export-flow/useExportSignatureFlow.ts` - 新增授权申请状态
- `frontend/src/pages/Signature_management_page.vue` - 新增授权受理按钮
- `frontend/src/boot/query/signature-query.ts` - 新增 5 个 API 函数
- `frontend/src/i18n/en-US/index.json` - 新增翻译键
- `frontend/src/i18n/zh-CN/index.json` - 新增翻译键

## 安全措施

1. **加密保护**：所有敏感字段使用独立的加密密钥加密
2. **双向验证**：授权文件生成和验证使用相同的算法，确保一致性
3. **哈希保护**：使用 SHA256 保护 authorizationUUID 和资格码
4. **部分暴露**：仅暴露哈希值的部分位数，降低碰撞攻击风险

## 不在范围内

- 授权文件的在线传输机制
- 授权撤销功能
- 批量授权功能

## 测试要点

1. **授权申请流程**
   - 选择未授权签名
   - 查看原始作者联系方式
   - 导出 .ktauthreq 文件
   - 保存对话框需选择路径后才提示成功（支持 File System Access API，失败/取消不提示成功）

2. **授权受理流程**
   - 导入 .ktauthreq 文件
   - 验证签名匹配
   - 导出 .ktauth 文件
   - 保存对话框需选择路径后才提示成功（支持 File System Access API，失败/取消不提示成功）

3. **授权导入流程**
   - 导入 .ktauth 文件
   - 验证通过后更新授权列表
   - 签名状态更新为已授权

4. **端到端测试**
   - 完整走完申请→受理→导入→导出流程
