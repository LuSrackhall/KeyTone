# Encrypted Outputs

## Purpose

定义 KeyTone 在“加密产物”上的兼容性边界与构建身份（build identity）的关联。

本能力覆盖：

- 专辑导出文件（`.ktalbum`）加密
- 专辑配置加密（AES-GCM，基于 albumUUID 派生 key）
- 签名导出文件（`.ktsign`）加密
- 构建时密钥注入（不改源码、开源构建默认可用）

## Requirements

### Requirement: 构建时密钥注入

Normative: The system SHALL support overriding selected symmetric keys/secrets at build-time via Go `-ldflags -X` using XOR-obfuscated hex values; when not injected, it MUST fall back to the default hardcoded values so open-source builds remain functional.

#### Scenario: 开源构建未注入

- **GIVEN** 用户从公开源码直接构建且未提供私钥文件/注入参数
- **WHEN** 系统进行涉及加密的功能
- **THEN** 使用源码默认密钥/secret，并保持可用

#### Scenario: 私有构建注入

- **GIVEN** 用户提供私钥（例如通过 `sdk/setup_build_env.sh` 生成 `EXTRA_LDFLAGS`）
- **WHEN** 使用 `-ldflags -X` 注入混淆值构建
- **THEN** 构建产物的加密身份发生变化，加密产物可与未注入构建不兼容（预期）

### Requirement: `.ktalbum` 版本化 XOR 加密

Normative: The system SHALL encrypt the `.ktalbum` zip body using a versioned XOR key; it MUST store the key version in the file header and MUST select the decryption key by header version, with a v1 fallback on checksum mismatch.

Note: A local debug tool (ktalbum-tools) MAY additionally try both the injected key and the default open-source key for the same version to improve inspection compatibility. This does not change the main application's compatibility boundary.

#### Scenario: 导出写入 v2

- **GIVEN** 用户导出专辑
- **WHEN** 系统生成 `.ktalbum`
- **THEN** MUST set `header.Version=2` and encrypt with v2 key

#### Scenario: 导入按版本选择密钥

- **GIVEN** 用户导入 `.ktalbum`
- **WHEN** 系统解析文件头
- **THEN** MUST decrypt using the key corresponding to `header.Version`

### Requirement: 专辑配置 AES-GCM 派生密钥

Normative: The system SHALL derive a 32-byte AES key using `SHA256(secret + last6(sha1(albumUUID)))`, where `secret` is configurable via build-time injection; the encrypted config bytes MUST be stored as `nonce + ciphertext`.

Note: The internal debug utility `sdk/audioPackage/cmd/printconfig` can be used to inspect decrypted config; for private builds it must be run/built with the same `-ldflags -X` injection (e.g. via `EXTRA_LDFLAGS`).

#### Scenario: 使用默认 secret 派生 key

- **GIVEN** 未注入自定义 secret
- **WHEN** 系统派生专辑配置 AES key
- **THEN** 使用默认 secret，派生结果与当前版本一致

#### Scenario: 注入 secret 改变派生 key

- **GIVEN** 构建时注入了自定义 secret
- **WHEN** 系统派生专辑配置 AES key
- **THEN** 派生结果与默认构建不同，导致加密配置不可跨构建身份通用（预期）

### Requirement: 签名导出 KeyB

Normative: The system SHALL encrypt `.ktsign` using KeyB, where KeyB SHALL be build-injectable; when not injected, it MUST use the default hardcoded value.

#### Scenario: KeyB 未注入

- **GIVEN** 未注入 KeyB
- **WHEN** 导出签名
- **THEN** 使用默认 KeyB 加密 `.ktsign`

#### Scenario: KeyB 注入

- **GIVEN** 注入了 KeyB
- **WHEN** 导出签名
- **THEN** `.ktsign` MUST be encrypted with the injected KeyB
