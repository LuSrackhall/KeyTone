# Tasks: 构建注入式对称密钥适配

## 1. 签名加密密钥适配 (sdk/signature/encryption.go)

- [x] 1.1 将 KeyA/KeyB 从 `const` 改为 `var`（保留 Default 常量）
- [x] 1.2 GetKeyA/GetKeyB：未注入用默认；注入后走 XOR 解混淆

## 2. 专辑导出密钥适配 (sdk/server/server.go)

- [x] 2.1 将 `KeytoneEncryptKeyV1/V2` 从 `const` 改为 `var`（保留 Default 常量）
- [x] 2.2 添加 `GetEncryptKeyV1/V2/Current()`：未注入用默认；注入后走 XOR 解混淆
- [x] 2.3 更新版本选择逻辑与回退：checksum 失败时回退尝试 v1
- [x] 2.4 移除无意义的硬编码清理 key：`CleanupOrphanCardImages(nil)`

## 3. 专辑配置加密密钥适配

- [x] 3.1 适配 `FixedSecret` (sdk/audioPackage/enc/enc.go)
  - [x] const -> var + Default
  - [x] `GetFixedSecret()` 支持可变长度解混淆
  - [x] `DeriveKey()` 使用 `GetFixedSecret()`
- [x] 3.2 适配 `KeyToneAlbumSignatureEncryptionKey` (sdk/signature/album.go)
  - [x] const -> var + Default
  - [x] `GetAlbumSignatureKey()` 支持解混淆逻辑

## 4. 构建脚本更新（SDK）

- [x] 4.1 更新 `sdk/private_keys.template.env` 增加 KEY_A/KEY_B/KEY_V1/KEY_V2/FIXED_SECRET/KEY_ALBUM_SIG
- [x] 4.2 更新 `sdk/setup_build_env.sh` 扩展 `KEYS_TO_PROCESS` 映射

## 5. 工具适配（tools/ktalbum-tools）

- [x] 5.1 `utils/header.go`：加入 v1/v2 key + 注入支持 + 候选 key 列表
- [x] 5.2 `commands/info.go` / `commands/extract.go`：按候选 key 解密并用 checksum 验证
- [x] 5.3 `setup_build_env.sh`：为 ktalbum-tools 生成并导出 `EXTRA_LDFLAGS`
- [x] 5.4 `build.sh`：应用 `EXTRA_LDFLAGS`

## 6. 规格/文档同步

- [x] 6.1 创建 OpenSpec spec delta：
  - [x] `specs/signature-management/spec.md`
  - [x] `specs/export-flow/spec.md`
- [x] 6.2 更新 BUILD_COMPATIBILITY 文档，说明对称密钥可注入

## 7. 验证清单

- [x] 7.1 `go test ./signature ./server ./audioPackage/enc`（SDK）
- [x] 7.2 `go test ./...`（ktalbum-tools）
- [ ] 7.3 `go test ./...`（SDK 全量）
  - 说明：当前全量测试失败包含与本变更无关的编译/测试问题（见终端输出：`keySound` 的 fmt.Sprintf 类型不匹配、`audioPackage/config` 测试用例文件路径问题）。
