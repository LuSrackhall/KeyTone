# Encrypted Outputs（增量）

## MODIFIED Requirements

### Requirement: 构建身份与对称密钥来源

Normative: The system SHALL keep default symmetric keys/secrets hardcoded in source for open-source builds, and SHALL allow overriding them at build-time via Go `-ldflags -X` using XOR-obfuscated hex values; at runtime it MUST deobfuscate injected values before use.

#### Scenario: 未注入私钥时保持原行为

- **GIVEN** 构建过程中未提供任何 `-ldflags -X` 覆盖值
- **WHEN** 系统进行签名管理加解密、专辑导出/导入、专辑配置加解密
- **THEN** 系统使用源码默认密钥/secret，行为与当前开源版本一致

#### Scenario: 注入密钥后自动解混淆

- **GIVEN** 构建过程中通过 `-ldflags -X` 注入了 XOR 混淆后的 hex
- **WHEN** 系统在运行时读取该变量并用于加解密
- **THEN** 系统先执行 `hex -> xorMask -> plaintext` 解混淆，再进行加解密

### Requirement: 专辑导出版本化密钥选择

Normative: The system SHALL encrypt `.ktalbum` file body using the current version key (v2), store `header.Version=2`, and SHALL select decryption key by `header.Version` with a v1 fallback on checksum mismatch.

#### Scenario: 导出使用 v2

- **GIVEN** 用户导出专辑
- **WHEN** 系统生成 `.ktalbum`
- **THEN** `header.Version` MUST be `2` and the file body MUST be encrypted with the v2 key

#### Scenario: 导入按版本解密并回退

- **GIVEN** 用户导入 `.ktalbum`
- **WHEN** 系统解密 zip body
- **THEN** 系统按 `header.Version` 选择密钥；若校验失败且版本不是 v1，则 MUST 尝试 v1 回退
