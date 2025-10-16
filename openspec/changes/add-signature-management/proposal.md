# 提案：添加签名管理功能

## Why

KeyTone 用户在 itch.io 社区分享键音专辑时，缺少对作品原创性和创作者身份的标识机制，导致版权难以追溯、创作者贡献得不到认可。需要提供数字签名功能，让创作者在导出专辑时附加个人签名信息，实现作品溯源、版权保护和社区信用建设。

## What Changes

- **新增**：签名管理页面和对话框组件，支持独立页面和对话框两种使用模式
- **新增**：签名 CRUD 功能（创建、编辑、删除、查看列表）
- **新增**：签名数据结构（名称、个人介绍、名片图片、保护码）
- **新增**：签名文件导入/导出功能（.ktsign 文件格式）
- **新增**：AES-256 对称加密存储机制
- **新增**：专辑导出流程中的签名选择步骤
- **新增**：后端图片访问 API（`GET /signature/image/:filename`）
- **新增**：图片预览功能（快速预览 + 大图预览）
- **修改**：侧边栏导航（新增"签名管理"入口）
- **修改**：专辑文件格式（扩展 signatures 字段）
- **修改**：配置文件结构（新增 signature_manager 字段）
- **复用**：现有 SSE 全量配置数据推送机制（无需新增代码）

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
  - 配置：
    - `KeyToneSetting.json`（运行时配置文件结构）

- **Dependencies**:
  - 前端：无新增依赖（移除 crypto-js 和 nanoid，加密逻辑在后端）
  - 后端新增：`github.com/jaevor/go-nanoid`（用于生成保护码）
  - 后端使用标准库：`crypto/aes`, `crypto/cipher`（用于加密）

---

## 详细设计说明

### 数据结构

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
  "signature_manager": {
    "encrypted_protect_code_1": "encrypted_signature_data_1",
    "encrypted_protect_code_2": "encrypted_signature_data_2",
    ...
  }
}
```

### 关键技术决策

1. **加密方案**：使用对称加密（AES-256），保护码经加密后作为 key
2. **存储方式**：复用现有 `/store/get` 和 `/store/set` API
3. **数据同步**：通过现有 SSE 全量配置数据推送机制实时同步（无需额外开发）
4. **图片存储**：
   - 前端通过选择器获取 File 对象，表单中保存 File 对象（非路径字符串）
   - 提交时转为 Base64 通过 HTTP 传输到后端
   - 后端解码 Base64，计算 SHA-256 哈希作为文件名保存到配置目录
   - 配置文件中仅存储路径字符串
   - 列表渲染时，前端将路径字符串转为 HTTP URL 访问图片资源
   - 导出/导入时，图片在 Base64 和文件系统之间互转
5. **文件格式**：`.ktsign` 文件为 JSON 格式，包含 Base64 编码的图片数据

### 用户交互流程

1. **创建签名**：用户填写表单（图片为 File 对象）→ 后端生成保护码 → 图片转 Base64 传输 → 后端保存并加密 → 现有 SSE 自动推送全量配置
2. **编辑签名**：获取签名数据 → 修改（新图片为 File 对象）→ 图片转 Base64 传输 → 保存 → 现有 SSE 自动推送全量配置
3. **导出签名**：选择签名 → 后端将图片转 Base64 → 生成 `.ktsign` 文件 → 下载
4. **导入签名**：选择文件 → 解析验证 → Base64 图片解码保存 → 现有 SSE 自动推送全量配置
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

