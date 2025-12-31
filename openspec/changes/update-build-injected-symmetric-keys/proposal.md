# Change: 旧有对称加密密钥接入构建注入体系

## Why

项目已在“授权流密钥”中采用构建时注入（`-ldflags -X` + XOR+hex 混淆）将私钥注入到构建结果中，并保持开源构建可用。

但项目中仍存在更早实现的对称加密能力，其密钥/secret 直接硬编码在源码中，导致无法与“构建身份”机制统一，且私有构建难以做到与开源构建的加密产物隔离。

## What Changes

- 将以下对称密钥/secret 从 `const` 改为可注入 `var`，默认值保持原硬编码字符串不变：
  - 签名管理 KeyA / KeyB
  - 专辑配置 `signature` 字段内层加密密钥
  - 专辑导出文件 XOR 密钥（v1/v2）
  - 专辑配置加密派生 secret（FixedSecret）
- 统一注入方式：注入值为 XOR 混淆后的 hex 字符串；运行时自动解混淆。
- 更新构建脚本与模板：`sdk/setup_build_env.sh` 与 `sdk/private_keys.template.env` 增加新 key 项。
- 为本地调试工具 ktalbum-tools 补齐注入脚本：`tools/ktalbum-tools/setup_build_env.sh`（默认复用 `sdk/private_keys.env`）。
- ktalbum-tools 构建脚本会自动应用 `EXTRA_LDFLAGS`：`tools/ktalbum-tools/build.sh`。
- 更新文档：`BUILD_COMPATIBILITY.md` 补充“Build-Time Injected Keys”列表。

## Non-Goals

- 不改变默认开源构建行为：未提供私钥注入时，仍使用源码默认值。
- 不更换加密算法（仅调整密钥来源/注入方式）。

## Impact

- Affected code:
  - `sdk/signature/encryption.go`
  - `sdk/signature/album.go`
  - `sdk/audioPackage/enc/enc.go`
  - `sdk/server/server.go`
  - `sdk/setup_build_env.sh`
  - `sdk/private_keys.template.env`
  - `tools/ktalbum-tools/utils/header.go`
  - `tools/ktalbum-tools/commands/*.go`
  - `tools/ktalbum-tools/setup_build_env.sh`
  - `tools/ktalbum-tools/build.sh`
  - `tools/ktalbum-tools/private_keys.template.env`
  - `BUILD_COMPATIBILITY.md`
  - `BUILD_COMPATIBILITY.zh-CN.md`

- Affected specs:
  - New capability: `openspec/specs/encrypted-outputs/spec.md`

## Review Notes / Audit Trail

| Key/Secret              | Default（源码）                                                              | 注入变量（Go -ldflags -X）                             | 用途摘要                                                |
| ----------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------- |
| 签名 KeyA               | `KeyTone2024Signature_KeyA_SecureEncryptionKeyForIDEncryption`               | `KeyTone/signature.KeyToneSignatureEncryptionKeyA`     | 加密签名ID、派生动态密钥（PBKDF2）                      |
| 签名 KeyB               | `KeyTone2024Signature_KeyB_SuperSecureEncryptionKeyForExportImportOperation` | `KeyTone/signature.KeyToneSignatureEncryptionKeyB`     | `.ktsign` 导入/导出加密                                 |
| 专辑 signature 字段密钥 | `KeyTone2024Album_Signature_Field_EncryptionKey_32Bytes`                     | `KeyTone/signature.KeyToneAlbumSignatureEncryptionKey` | 专辑配置中 `signature` 字段内层 AES-GCM                 |
| 专辑导出 XOR v1         | `KeyTone2024SecretKey`                                                       | `KeyTone/server.KeytoneEncryptKeyV1`                   | `.ktalbum` v1 加/解密（兼容）                           |
| 专辑导出 XOR v2         | `KeyTone2025AlbumSecureEncryptionKeyV2`                                      | `KeyTone/server.KeytoneEncryptKeyV2`                   | `.ktalbum` v2 加/解密（当前）                           |
| 专辑配置派生 secret     | `LuSrackhall_KeyTone_2024_Signature_66688868686688`                          | `KeyTone/audioPackage/enc.FixedSecret`                 | 派生 AES key：`SHA256(secret + last6(sha1(albumUUID)))` |

注：工具链 `tools/ktalbum-tools` 也提供同等注入点（模块路径不同）。

另外：ktalbum-tools 为本地查看/调试用途，解密 `.ktalbum` 时会按顺序尝试“注入密钥 → 开源默认密钥”，并在校验失败时回退尝试 v1，以便同一构建可同时兼容开源与私有两类产物。
