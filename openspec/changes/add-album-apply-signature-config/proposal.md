# 变更：导出签名配置落盘与无签名授权分支

## 为什么

- 导出流程目前只在选择“需要签名”时执行加密与签名配置写入，导致“无需签名但二次创作需要授权”的场景下缺少必要的签名与联系方式记录。
- 需要一个新的后端路由让前端在加密后、真正导出前把所选签名与授权联系方式写入键音专辑配置（当前阶段只需记录日志，后续再落盘）。
- 前端状态机需要覆盖“无需签名”路径的授权判断，以便按需触发加密与新路由，并保证传参结构在所有分支一致。

## 变更内容

- 扩展 `useExportSignatureFlow`：无论用户在“是否需要签名”对话框中选择哪一项，都必须进入“二次创作是否需要授权”的分支；当选择“无需签名 + 需要授权”时，仍需填写联系方式并挑选签名。
- 新增路由 `POST /keytone_pkg/apply_signature_config`，位于加密与实际导出之间；当前实现只打印前端提交的表单 JSON，后续将真正落盘。
- 新的前端 API `applySignatureConfig` 负责向后端发送 `albumPath`、`needSignature`、`requireAuthorization`、`signatureId` 以及（仅在需要授权时）联系人邮箱与可选附加联系方式。
- 导出入口在计算“是否加密/是否调用新路由”时遵循：只要 `needSignature === true` 或 `requireAuthorization === true` 就需要执行加密与新路由；否则直接走原有导出接口。

## 影响范围

- 规格：`openspec/specs/export-flow/spec.md` 及新增变更说明。
- 前端：`frontend/src/components/export-flow/useExportSignatureFlow.ts`、`frontend/src/pages/Keytone_album_page_new.vue`、`frontend/src/boot/query/keytonePkg-query.ts`。
- 后端：`sdk/server/server.go` 以及与日志输出相关的结构。
- 未来任务：将日志输出替换为真正的专辑配置写入逻辑，并在落盘前校验字段。