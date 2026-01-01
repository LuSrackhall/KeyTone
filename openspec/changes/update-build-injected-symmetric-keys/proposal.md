# Proposal: 构建注入式对称密钥适配

## Why

项目中存在多处对称加密密钥/secret 以 `const` 或硬编码字符串形式存在于源码中。
这会导致：

- 官方/私有构建与社区构建可能共享同一套密钥，无法形成“构建身份差异”
- 无法复用现有授权流（Key F/K/Y/N）已经采用的“构建时注入 + 运行时解混淆”机制

本变更的目标是：在**不破坏开源默认行为**的前提下，让这些历史遗留对称密钥也支持本地私有构建注入。

## What Changes

- 将目标密钥从 `const` 改为 `var`，保留原始默认值为 `Default*` 常量
- 构建时通过 Go `-ldflags -X` 注入**XOR 混淆后的 hex**（由 `tools/key-obfuscator` 生成）
- 运行时按“默认值直用 / 注入值解混淆”的规则取值
- ktalbum-tools 增强：按文件头版本选择 key，校验失败时回退尝试 v1；并支持“注入 key → 默认 key”候选顺序
- 文档同步：补充 BUILD_COMPATIBILITY 对“可选对称密钥注入”的说明，并提供 OpenSpec spec delta

## Scope

### 适配对象（密钥清单）

- 签名密钥：KeyA / KeyB
- `.ktalbum` XOR 密钥：v1 / v2
- 专辑配置加密种子：`FixedSecret`
- 专辑配置 `signature` 字段内层 key：`KeyToneAlbumSignatureEncryptionKey`
- 调试工具（ktalbum-tools）对应的 v1/v2 XOR key

### 非目标

- 不改变任何默认密钥的明文内容（开源构建保持原行为）
- 不引入远程 KMS / 在线密钥管理
- 不修改加密算法本身（仅变更密钥来源）

## Impact

- 受影响的能力规格：
  - `signature-management`
  - `export-flow`

- 受影响的代码范围（关键路径）：
  - `sdk/signature/*`（KeyA/KeyB、album signature key）
  - `sdk/audioPackage/enc`（FixedSecret）
  - `sdk/server`（ktalbum v1/v2 keys）
  - `tools/ktalbum-tools`（调试解密回退）

## Compatibility

- 未注入：与当前开源版本行为一致
- 注入后：不同构建身份之间的加密产物可能不互通（符合 BUILD_COMPATIBILITY 设计）
