# 变更：专辑配置加密与透明降级加载

## 为什么需要

当专辑需要签名导出时，明文 `package.json` 会暴露与签名相关的敏感元数据。我们需要一个最小化、向后兼容的机制，仅在导出时选择"需要签名"的专辑加密配置，同时允许未签名专辑继续使用现有的明文 JSON 工作流。加载必须优雅降级：先尝试常规 JSON 读取；若失败则尝试解密并挂载临时解密副本供 Viper 使用。

## 变更内容

- 引入基于固定密钥 + 专辑 UUID 派生后缀的确定性对称加密。
- 仅当导出流程确认"需要签名"时，将明文配置转存为同目录下的二进制文件 `core`（无后缀），`package.json` 则写入最小化的指示 JSON（仅包含 `_keytone_core` 等元数据）。
- 加载时：读取 `package.json`，若检测到指示 JSON 或无法解析明文，则读取 `core`、解密、写入临时目录、通过 Viper 加载解密后的 JSON 并设置监听器以在写回时重新加密写回 `core`。
- 配置变更时（Viper 写入 / 文件监听）：重新加密更新后的 JSON 并原子覆盖源文件。
- 导入流程：检测加密的专辑配置（JSON 解析失败）并透明处理解密。
- 提供独立调试 CLI（`audioPackage/cmd/printconfig`）用于打印解密后的配置，可识别新的 `core` 文件存储格式。
- 提供单元测试覆盖密钥派生、加解密往返、降级加载和写回重新加密。
- 不破坏现有未签名专辑（它们保持明文）。
- 新增独立 API 路由 `POST /encrypt_album_config` 接受 `albumPath` 参数，前端在选择"需要签名"时调用以加密对应专辑；支持旧版十六进制密文自动迁移至 `core` 存储。
- 导入（保存为新专辑）时，若源配置已加密（含 `_keytone_core` 指示 JSON），系统自动使用原 UUID 解密，更新 `audio_pkg_uuid` 字段，使用新 UUID 重新加密并写回 `core`，避免 AES-GCM 认证失败。

## 影响范围

- 新增能力：`album-config-encryption`。
- 受影响代码：`sdk/audioPackage/config/audioPackageConfig.go`、专辑导出逻辑新增独立加密路由于 `sdk/server/server.go`、潜在的导入路径、新增 CLI 于 `sdk/audioPackage/cmd/printconfig`。
- 安全性：防止随手窥视；不主张强安全性（固定密钥嵌入代码）。需文档说明局限性。

## 不在范围内

- 密钥轮换。
- 字段级加密或部分混淆。
- 远程密钥管理 / KMS 集成。

## 待确认问题

- **已确认**：加密配置分离为指示 JSON `package.json` + 二进制 `core` 文件；临时文件使用 OS 临时目录自动清理。
- **已确认**：导入时自动识别加密指示 JSON，使用原 UUID 解密，按新 ID 重新加密。
