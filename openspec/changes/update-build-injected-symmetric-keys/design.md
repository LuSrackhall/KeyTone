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

### Compatibility Behavior (重要)

为保证历史 `private_keys.env`（未包含新增 KEY_*）以及“仅开源构建”场景可正常运行：

- `sdk/setup_build_env.sh` 采用 best-effort 策略：
  - 若未找到 `private_keys.env`，脚本不会报错退出，而是设置 `EXTRA_LDFLAGS=""` 并继续（等价于不注入）。
  - 若某个 `KEY_*` 缺失，或值仍为模板占位符（如 `PLACEHOLDER_*` / `REPLACE_ME`），脚本会跳过该 key 的注入，让运行时回退到源码默认值。

该行为的目标是：在不要求用户立刻补齐新增环境变量的前提下，依旧保持与适配前版本的兼容性。

### Shell Robustness Notes

为了保证在 `dev.sh` 的严格模式（`set -euo pipefail`）下也稳定：

- 当脚本需要在“缺失私钥文件”场景下提前结束时，必须在脚本顶层直接使用 `return ... || exit ...`，而不是在函数内封装 `return`（函数内 `return` 只会返回函数，无法中止被 `source` 的脚本继续执行）。
- 在包含中文标点的输出字符串中，变量展开建议统一使用 `${VAR}` 形式，避免某些 locale/编码情况下 `set -u` 将紧邻的非 ASCII 字节误判为变量名的一部分。

### ktalbum-tools（本地调试工具）

ktalbum-tools 的目标是便于本地查看/调试 `.ktalbum` 文件，因此它需要在“一个二进制”中尽量兼容两类产物：

- 开源默认密钥产物
- 私有构建注入密钥产物

实现方式：

- 构建注入脚本：`tools/ktalbum-tools/setup_build_env.sh`（默认复用 `sdk/private_keys.env`，也可使用本地 `tools/ktalbum-tools/private_keys.env`）
- 一键构建：`tools/ktalbum-tools/build.sh` 会自动合并并应用 `EXTRA_LDFLAGS`
- 运行时解密：按顺序尝试“注入密钥 → 默认密钥”，并在校验失败时按 SDK 策略回退尝试 v1

### printconfig（SDK 内部调试工具）

`sdk/audioPackage/cmd/printconfig` 是用于解密查看键音专辑配置的内部工具。

- 位于 SDK 模块内，构建时自动继承 `EXTRA_LDFLAGS` 注入的 `FixedSecret`
- 无需单独配置注入脚本
- 私有构建后可解密私有产物的加密配置

使用方式补充：

- 可直接 `go run ./audioPackage/cmd/printconfig --path ...`（开源默认密钥）
- 若要解密私有构建产物，需要 `go run -ldflags "$EXTRA_LDFLAGS" ./audioPackage/cmd/printconfig --path ...`
- 注入发生在编译/链接阶段，因此同一次运行无法同时兼容两套密钥（需分别运行）

## Compatibility Notes

- 未注入：行为与当前开源版本完全一致
- 注入后：
  - 官方/私有构建与社区构建的加密产物可能不互通（符合 BUILD_COMPATIBILITY 设计）
  - 若选择覆盖 v1 key，将导致旧 v1 产物在该私有构建中不可解密（预期行为，需在 proposal 中显式告知）
