# 任务：应用签名配置

## 1. 前端交互

- [x] 1.1 扩展 `useExportSignatureFlow`，在“无需签名”路径也进入授权策略分支，并在 `requireAuthorization=false`+`needSignature=false` 时直接完成。
- [x] 1.2 更新授权结果与签名选择逻辑，确保 `requireAuthorization=true` 时总是收集签名 ID + 联系邮箱（可附带额外方式）。
- [x] 1.3 在 `Keytone_album_page_new.vue` 中根据 `needSignature`/`requireAuthorization` 计算加密与新路由调用条件，并串联 API。
- [x] 1.4 新增 `applySignatureConfig` API 包装器，统一表单字段与错误提示。

## 2. 后端 API

- [x] 2.1 在 `sdk/server/server.go` 中注册 `POST /keytone_pkg/apply_signature_config`，验证 `albumPath`/签名 ID/授权字段并记录日志。
- [x] 2.2 暴露结构体占位符，为后续真正写入配置留接口。

## 3. 规格与校验

- [x] 3.1 在 `openspec/changes/add-album-apply-signature-config/specs/export-flow/spec.md` 记录状态机与路由新增要求。
- [x] 3.2 更新 `proposal.md`、`design.md` 与本文件，保持与实现同步。
- [x] 3.3 手动验证四种分支：
  - 无需签名 + 无需授权 → 直接导出。
  - 无需签名 + 需要授权 → 触发加密 + 新路由。
  - 需要签名 + 无需授权 → 加密 + 新路由，仅带签名 ID。
  - 需要签名 + 需要授权 → 完整表单链路。