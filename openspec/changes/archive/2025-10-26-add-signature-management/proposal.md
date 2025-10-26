# 提案：添加签名管理功能

## Why

KeyTone 用户在 itch.io 社区分享键音专辑时，缺少对作品原创性和创作者身份的标识机制，导致版权难以追溯、创作者贡献得不到认可。需要提供数字签名功能，让创作者在导出专辑时附加个人签名信息，实现作品溯源、版权保护和社区信用建设。

## What Changes

- **新增**：签名管理页面和对话框组件，支持独立页面和对话框两种使用模式
- **新增**：签名 CRUD 功能（创建、编辑、删除、查看列表）
- **新增**：签名数据结构（名称、个人介绍、名片图片、保护码）
- **新增**：签名文件导入/导出功能（.ktsign 文件格式）
- **新增**：AES-256-GCM 对称加密存储机制（KeyA/KeyB + 动态密钥派生）
- **新增**：专辑导出流程中的签名选择步骤
- **新增**：后端图片读取 API（`POST /signature/get-image`，按路径读取并返回二进制流）
- **新增**：图片预览功能（快速预览 + 大图预览）
- **修改**：侧边栏导航（新增"签名管理"入口）
- **修改**：专辑文件格式（扩展 signatures 字段）
- **修改**：配置文件结构：新增 `signature` 字段，键为加密ID（KeyA），值为 `{ value: <用动态密钥加密>, sort: { time } }`
- **复用**：现有 SSE 全量配置数据推送机制（无需新增代码）

### 新增目标：缺陷修复与稳定性

为解决近期反馈的关键问题，本变更追加以下非破坏性修复目标（兼容现有接口与数据结构）：

1. 创建成功后签名列表被意外清空
- 修复创建签名后，列表状态被错误重置为“空”的问题；确保创建成功后，列表包含原有项与新建项，不需重启即可可见。
- 加固 SSE 全量配置更新与本地列表状态的同步与合并策略，避免清空或闪烁。

1. 签名图片目录存储到当前执行目录而非 ConfigPath
- 在应用启动（main.go）阶段初始化签名模块的存储根路径为 ConfigPath（如：ConfigPath/signatures/card_images）。
- 所有签名图片的写入/读取均以该初始化路径为根，避免在当前工作目录新建 signatures/。

1. 删除签名总是失败的通知问题
- 校正前端删除调用参数、方法与后端路由（DELETE /signature/delete/:id）的一致性；
- 后端返回码与错误消息规范化；前端根据响应正确提示成功/失败。

1. 更新签名总是失败的通知问题
- 校正前端 PUT /signature/update 的负载与字段（必须包含 id，图片可选 Base64）；
- 后端对无变更/部分字段更新正确处理，返回成功码；
- 前端根据结果与 SSE 同步刷新列表。

1. 编辑对话框二次打开数据为空
- 修复编辑对话框在关闭后再次打开时的表单复位/初始化问题；
- 确保以签名 id 为 key 的受控表单在每次打开时重新填充数据，或在对话框挂载时拉取当前项数据；
- 二次打开应显示正确的签名名称、介绍与图片预览。

## 关键 UI/UX 变更

- **签名列表项布局**：
  - 名片图片仅在左侧展示小图像（约 80x80px），不是主要布局
  - 注意力布局改为签名名称和介绍（右侧为主要信息区域）
  - 点击小图像打开预览，查看完整大图
  - 点击签名名称/介绍区域进入编辑模式

- **编辑/创建对话框中的图片管理**：
  - 图片选择器在表单字段下方
  - 图片快速预览显示在选择器下方（而不是上方）
  - 快速预览可点击，打开大图预览对话框
  - 图片格式不受严格限制，仅支持 webview 支持的格式均可

- **图片格式灵活性**：
  - 移除严格的图片格式限制
  - 允许用户上传 webview 支持的任何图片格式（PNG、JPG、WebP、GIF 等）

## Impact

- **Affected specs**:
  - `signature-management`（新增 capability）
  - `album-export`（修改现有 capability）
  - `config-storage`（扩展现有 capability）
  
- **Affected code**:
  - 前端：
    - `frontend/src/layouts/Main_layout.vue`（侧边栏导航）
    - `frontend/src/pages/Signature_management_page.vue`（新增）
    - `frontend/src/components/SignatureManagementDialog.vue`（新增）
    - `frontend/src/components/SignatureFormDialog.vue`（新增）
    - `frontend/src/components/SignatureSelectDialog.vue`（新增）
    - `frontend/src/services/signature-service.ts`（新增）
    - `frontend/src/utils/encryption.ts`（新增）
    - `frontend/src/types/signature.ts`（新增）
    - 专辑导出相关组件（待确认具体路径）
  - 后端：
    - `sdk/server/server.go`（新增 API 端点）
    - `sdk/config/config.go`（扩展配置字段支持）
    - `sdk/signature/encryption.go`（新增加密/解密逻辑）
    - `tools/ktalbum-tools/utils/album.go`（专辑格式扩展）
    - `sdk/main.go`（初始化签名存储路径，传入 ConfigPath 到签名模块）
    - `sdk/server/signature_handlers.go`（删除与更新端点返回码/错误处理一致性修正）
  - 配置：
    - `KeyToneSetting.json`（运行时配置文件结构）

- **Dependencies**:
  - 前端：无新增依赖（移除 crypto-js 和 nanoid，加密逻辑在后端）
  - 后端无外部依赖变更（使用标准库 crypto/aes、cipher、rand、encoding/hex 等；动态密钥使用 x/crypto/pbkdf2）
  - 后端使用标准库：`crypto/aes`, `crypto/cipher`（用于加密）

---

## 详细设计说明

### 数据结构（更新）

```typescript
// 前端签名数据结构
interface Signature {
  name: string;           // 签名名称（必填）
  intro?: string;         // 个人介绍（选填）
  cardImage?: string;     // 名片图片路径（选填）
  protectCode: string;    // 保护码（前端自动生成，nanoid）
  createdAt: string;      // 创建时间
  id: string;             // 唯一标识（前端生成，nanoid）
}

// 后端存储结构（配置文件中）
{
  "signature": {
    "<encryptedId>": {
      "value": "<encryptedSignatureJSON>",
      "sort": { "time": 1730000000 }
    }
  }
}
```

### 关键技术决策

1. **加密方案**：AES-256-GCM；ID 用 KeyA 加密作为键；Value 使用基于 ID 后7位与 KeyA 的 PBKDF2 动态密钥加密；导出/导入整体使用 KeyB 加解密
2. **存储方式**：通过 config.SetValue/GetValue 维护 `signature` 字段
3. **数据同步**：通过现有 SSE 全量配置数据推送机制实时同步（无需额外开发）
4. **图片存储**：
- 上传采用 multipart/form-data（字段 cardImage）；后端写入 ConfigPath/signature，文件名为 SHA-1(id|name|originalName|timestamp)
- 配置中存储的是绝对路径；前端通过后端 `POST /signature/get-image` 或读取为 Blob URL 展示
- 导出/导入时图片在二进制与十六进制字符串之间转换（内部 JSON 用十六进制表示图片数据），外层整体再用 KeyB 加密
1. **文件格式**：`.ktsign` 文件为 JSON 格式，包含 Base64 编码的图片数据

### 用户交互流程

1. **创建签名**：用户填写表单（multipart）→ 后端保存图片并加密存储 → SSE 自动推送全量配置
2. **编辑签名**：表单（multipart）支持 removeImage/imageChanged → 后端条件更新并加密存储 → SSE 自动推送全量配置
3. **导出签名**：选择签名 → 后端生成内部 JSON（图片十六进制）→ 用 KeyB 加密 → 以二进制流下载 `.ktsign`
4. **导入签名**：上传 `.ktsign` → 后端用 KeyB 解密并解析 → 冲突返回 409，确认后走 `/signature/import-confirm` → SSE 自动推送全量配置
5. **专辑签名**：导出专辑时 → 选择签名 → 嵌入签名信息到专辑文件

## 风险与缓解

### 技术风险

- **风险**：加密密钥泄露可能导致签名数据被破解
  - **缓解**：使用行业标准加密算法，密钥存储在代码中（开源项目的常见做法）

- **风险**：签名数据与专辑数据不一致
  - **缓解**：在导出时进行完整性校验，确保签名数据有效

### 兼容性风险

- **风险**：旧版本无法识别带签名的专辑
  - **缓解**：签名数据作为可选字段，不影响核心功能；提供版本检测

### 用户体验风险

- **风险**：签名管理流程过于复杂
  - **缓解**：采用简洁的表单设计，提供明确的操作指引

## 验收标准

### 功能完整性

- [ ] 用户可以创建、编辑、删除签名
- [ ] 用户可以导入和导出 `.ktsign` 文件
- [ ] 签名列表可以正确显示和刷新
- [ ] 专辑导出时可以选择签名
- [ ] 签名数据安全存储和加密

### 性能要求

- [ ] 签名列表加载时间 < 500ms
- [ ] 签名操作响应时间 < 1s
- [ ] 图片上传和处理 < 2s

### 用户体验

- [ ] 界面符合 Quasar 设计规范
- [ ] 支持中英文国际化
- [ ] 错误提示清晰友好
- [ ] 操作流程符合用户直觉

### 安全性

- [ ] 签名数据加密存储
- [ ] 文件导入时进行格式验证
- [ ] 防止 XSS 和注入攻击

## 实施计划

### 阶段 1：基础架构（3-5 天）

- 定义数据模型和 API 接口
- 实现加密/解密工具函数
- 搭建签名管理页面框架

### 阶段 2：核心功能（5-7 天）

- 实现签名 CRUD 操作
- 实现导入/导出功能
- 集成 SSE 数据同步

### 阶段 3：专辑集成（3-5 天)

- 修改专辑导出流程
- 实现签名选择和嵌入
- 添加签名验证逻辑

### 阶段 4：测试与优化（2-3 天）

- 单元测试和集成测试
- UI/UX 优化
- 性能调优

## 相关文档

- 详细设计：`design.md`
- 实施任务：`tasks.md`
- 功能规范：`specs/signature-management/spec.md`

## 变更历史

- 2025-10-15: 初始提案创建

