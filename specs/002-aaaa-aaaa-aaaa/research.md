# Research: 键音专辑签名系统（Phase 0）

## Decisions

- Decision: 采用混合测试策略（契约/核心逻辑测试优先 + UI 少量烟雾与组件单测）
  - Rationale: Electron + 本地后端导致全链路 TDD 成本高；先保证契约稳定和关键逻辑正确，UI 以手动验证为主、少量自动化兜底。
  - Alternatives: 严格端到端 TDD（高成本、脆弱）；仅手动测试（缺少回归保护）。

- Decision: REST 契约为前后端边界（gin 提供 /sdk/*），为后续 Wails 版本迁移保持协议稳定
  - Rationale: 与现有项目一致，易于编写契约测试与前端 axios 封装。
  - Alternatives: IPC/bridge 专用协议；JSON-RPC（后续可评估）。

- Decision: 签名文件（.ktsign）内容第一阶段最小化
  - Rationale: 优先实现导入/导出流转；校验逻辑逐步增强。
  - Alternatives: 一次性定义完整加密/授权体系（超范围）。

- Decision: 签名存储对称加密（AES-256-GCM + Base64 编码）
  - Rationale: AES-GCM 提供认证加密（AEAD），保证完整性与保密性；Base64 编码保证 JSON 序列化兼容；256 位密钥强度足够。
  - Alternatives: AES-CBC（需独立 HMAC）；ChaCha20-Poly1305（标准库支持待确认）；非对称加密（第二阶段授权场景考虑）。

- Decision: 全局配置签名以加密 protectCode 为 key、明文 JSON 为 value；专辑签名以加密 sha256(protectCode+name) 为 key、加密 JSON 为 value
  - Rationale: 全局配置允许客户端快速索引（明文 value 便于前端渲染）；专辑配置双重加密防止泄露签名逻辑与内容。
  - Alternatives: 全局配置也加密 value（增加前端解密复杂度）；专辑配置不加密 value（安全性不足）。

- Decision: 密钥管理第一阶段使用应用级固定密钥（32 字节随机值，硬编码于二进制或环境变量）
  - Rationale: 实现简单，满足第一阶段基本保护需求；后续可扩展为用户密码派生（PBKDF2/Argon2）或密钥链存储。
  - Alternatives: 用户密码派生（需增加密钥恢复与重置机制）；操作系统密钥链（需跨平台封装）；无加密（不符合安全要求）。

- Decision: nanoid 用于生成 UUID（protectCode、cardImage 文件名）
  - Rationale: 字母数字混合，可读性优于 UUID v4；碰撞概率满足需求；无需独立 id 字段，name 作为签名唯一标识。
  - Alternatives: UUID v4（更标准，但更长）；自增 ID（不适合分布式场景）。

- Decision: 签名无独立 id 字段，name 作为唯一标识
  - Rationale: 简化数据模型，name 本身已唯一且不可变；避免 id 与 name 双重维护；protectCode 仅用于加密/哈希，不作为业务标识。
  - Alternatives: 添加独立 id（增加复杂度，无实际收益）；使用 protectCode 作为标识（不符合业务语义）。

- Decision: 专辑配置键派生采用 encrypt(sha256(decrypt(protectCode) + name))
  - Rationale: 解密 protectCode 后与 name 拼接计算哈希，确保同一签名在不同专辑中键值一致；SHA-256 提供单向性与碰撞抗性；最终加密保证配置文件层面不泄露哈希规则。
  - Alternatives: 直接使用 encrypt(name)（无法防止专辑间签名关联分析）；encrypt(sha256(protectCode + name)) 不解密（逻辑错误，protectCode 本身已加密存储）。

## Unknowns Resolved

- 签名结构（最小）：name（唯一标识）、intro、cardImagePath、protectCode（nanoid 生成）、createdAt、signedAt 数组（专辑签名时）、authorization（仅原始作者）。
- .ktsign 文件：JSON 打包 + 名片图片引用（相对路径或内嵌 base64，优先相对路径）。第一阶段可不做强校验，仅保留哈希与时间戳。
- 错误码规范：HTTP 状态 + JSON { code, message }，重复/冲突时返回 409；参数校验 400。
- Electron 集成：前端仅经 REST 调用；文件保存/选择使用浏览器端接口或由后端返回目标路径后存储。
- 动态端口/运行时发现：Electron 主进程拉起 SDK 并输出 KEYTONE_PORT，渲染进程通过 `window.myWindowAPI.getBackendPort()` 获取端口并调用 `UpdateApi(port)`；必要时刷新以重建 SSE 链接。

## Best Practices

- gin：使用 binding 校验与统一错误响应中间件；路由前缀 /sdk/ 收敛。
- axios：统一封装 baseURL（避免到处硬编码），请求/响应拦截器统一错误提示；超时与重试策略最小化，避免 UI 卡顿。
- Playwright（Electron 开发流）：仅保留 1-2 条烟测（加载页面、签名管理入口存在、后端探针 200）。失败时提供 HTML report 与截图/视频。

## Notes

- 后续增强：签名校验与授权策略、跨平台路径与文件对话框封装、性能基线与回归。
