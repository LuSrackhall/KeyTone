# Design: 构建注入式对称密钥适配

## Overview

本设计将历史遗留的对称密钥/secret 统一到与授权流相同的构建注入体系：

- 默认值保留在源码中（开源构建无需私钥文件即可工作）
- 私有构建通过 `-ldflags -X` 覆盖变量值
- 覆盖值为 XOR 混淆后的 hex 字符串，运行时自动解混淆为明文

## Injection Format

### Build-Time

- 通过 Go 链接器 `-X 'package.path.VarName=VALUE'` 注入
- `VALUE` 推荐为：`tools/key-obfuscator` 输出的 hex（对任意长度字符串可用；长度非 32 时会提示 warning）

### Runtime

- 若变量值等于默认常量，则直接按默认明文使用
- 否则按以下逻辑尝试解混淆：
  1. `hex.DecodeString(value)`
  2. 对每个字节执行 `b ^ xorMask[i%len(xorMask)]`
  3. 转换为 `string` 或（对 32-byte key 场景）截断/补齐到 32
- 若注入值并非 hex（用户误注入明文），则回退为“直接使用该字符串”

## Key Inventory and Usage

### Signature Keys

- KeyA (`KeyToneSignatureEncryptionKeyA`)
  - 用途：加密签名 ID、派生动态密钥（PBKDF2）
  - 注入点：`KeyTone/signature.KeyToneSignatureEncryptionKeyA`

- KeyB (`KeyToneSignatureEncryptionKeyB`)
  - 用途：`.ktsign` 导入/导出对称加密
  - 注入点：`KeyTone/signature.KeyToneSignatureEncryptionKeyB`

### Album Export Keys (XOR)

- v1/v2：用于 `.ktalbum` 文件体（zip 数据）的 XOR 加/解密
- 文件头 `Version` 指定密钥版本；解密时校验失败会回退尝试 v1
- 注入点：
  - `KeyTone/server.KeytoneEncryptKeyV1`
  - `KeyTone/server.KeytoneEncryptKeyV2`

### Album Config Seed (FixedSecret)

- 用途：派生 AES key：`SHA256(secret + last6(sha1(albumUUID)))`
- 注入点：`KeyTone/audioPackage/enc.FixedSecret`
- 说明：该 secret 非固定 32 字节，因此采用“可变长度解混淆”路径（hex->xor->string）

### Album Signature Field Inner Key

- 用途：专辑配置 `signature` 字段内层 AES-GCM（外层仍由 albumUUID 派生 key 保护）
- 注入点：`KeyTone/signature.KeyToneAlbumSignatureEncryptionKey`

## Build Script Integration

- `sdk/private_keys.template.env` 增加对应 KEY_* 项
- `sdk/setup_build_env.sh` 增加 `KEYS_TO_PROCESS` 映射，使本地私有构建可以自动生成 `EXTRA_LDFLAGS`

## Compatibility Notes

- 未注入：行为与当前开源版本完全一致
- 注入后：
  - 官方/私有构建与社区构建的加密产物可能不互通（符合 BUILD_COMPATIBILITY 设计）
  - 若选择覆盖 v1 key，将导致旧 v1 产物在该私有构建中不可解密（预期行为，需在 proposal 中显式告知）
