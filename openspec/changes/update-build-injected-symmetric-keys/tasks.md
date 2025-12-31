# 旧有对称加密密钥接入构建注入体系 - 任务清单

## 状态

**当前状态**: ✅ 已实现（代码与文档已同步）

## 1. 发现与清点（留痕）

- [x] 在 `sdk/**/*.go` 中定位所有包含硬编码对称密钥/secret 的位置
- [x] 在 `tools/**/*.go` 中定位所有包含硬编码对称密钥/secret 的位置
- [x] 逐一标注用途、调用链、数据格式（32-byte key / 可变长度 secret / XOR key）

## 2. SDK：密钥注入适配

- [x] `sdk/signature/encryption.go`：KeyA/KeyB 改为可注入 `var`，默认保持原值
- [x] `sdk/signature/album.go`：专辑 signature 字段密钥改为可注入 `var`
- [x] `sdk/audioPackage/enc/enc.go`：FixedSecret 改为可注入 `var` 并支持可变长度解混淆
- [x] `sdk/server/server.go`：专辑导出 XOR v1/v2 key 改为可注入 `var`，使用时解混淆
- [x] `sdk/server/server.go`：移除/替换不必要的硬编码签名清理密钥，改用 `signature.GetKeyA()`

## 3. 工具链：ktalbum-tools 适配

- [x] `tools/ktalbum-tools/utils/header.go`：增加 v1/v2 key 与注入点、提供 `GetEncryptKeyByVersion`
- [x] `tools/ktalbum-tools/commands/extract.go`：按版本解密 + 校验失败回退 v1
- [x] `tools/ktalbum-tools/commands/info.go`：按版本解密 + 校验失败回退 v1

## 4. 构建入口与模板

- [x] `sdk/setup_build_env.sh`：追加 KEY_A/KEY_B/KEY_ALBUM_* 映射
- [x] `sdk/private_keys.template.env`：追加新增 KEY_* 项

## 5. 文档与规格同步

- [x] `BUILD_COMPATIBILITY.md`：补充 Build-Time Injected Keys
- [x] 新增能力规格：`openspec/specs/encrypted-outputs/spec.md`
- [x] 变更增量规格：`openspec/changes/update-build-injected-symmetric-keys/specs/encrypted-outputs/spec.md`

## 6. 验证

- [x] `go build ./...`（sdk）
- [ ] 可选：在本地提供 `sdk/private_keys.env` 并运行 `source sdk/setup_build_env.sh` 后构建，验证注入 keys 生效
- [ ] 可选：为 `tools/ktalbum-tools` 构建时传入 `-ldflags -X`，验证能解密对应构建身份产物
