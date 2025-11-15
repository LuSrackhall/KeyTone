# 设计：签名应用路由与无签名授权扩展

## Context
- 导出流程存在两个阶段：对话链收集需求 → 真实导出（加密 → 新签名路由 → 导出文件）。
- 目前 `useExportSignatureFlow` 在“无需签名”分支直接结束，无法得知是否需要授权。
- 新增路由只需校验参数并打印日志，因此该版本专注于打通前后端契约，后续可在同一路由内落盘。

## Goals / Non-Goals
- **Goals**
  - 前端在所有分支下都能构造统一的 `ApplySignatureConfigPayload`。
  - 根据 `needSignature` 与 `requireAuthorization` 自动决定是否触发加密与新路由。
  - 后端快速收到请求并记录可观测日志，便于后续实现。
- **Non-Goals**
  - 真正写入专辑配置或签名文件内容。
  - 新的授权 UI（沿用现有对话框）。
  - 针对已有签名专辑的授权门控改造（留待后续）。

## Decisions
1. **路由路径**：沿用 `keytone_pkg` 分组新增 `POST /apply_signature_config`，与加密路由同一命名空间，方便前端调用顺序固定为“加密 → 应用签名 → 导出”。
2. **payload 结构**：

    ```json
    {
      "albumPath": "...",
      "needSignature": true,
      "requireAuthorization": false,
      "signatureId": "enc-id",
      "contactEmail": "user@example.com",
      "contactAdditional": "Discord: foo"
    }
    ```

    - 只要调用该路由就必须携带 `signatureId`，因为“无需签名且无需授权”不会触发调用。
    - `contactEmail` 仅在 `requireAuthorization === true` 时填写，且前端需完成格式校验；`contactAdditional` 为可选描述字段。
3. **加密条件**：`shouldEncrypt = needSignature || requireAuthorization`，确保“无需签名但需要授权”也会生成受保护的配置。
4. **状态机改造**：
   - `handleConfirmSignatureSubmit(false)` 也进入 `auth-requirement`，只有在用户明确选择“无需签名且无需授权”时才 `done`。
   - `handleAuthRequirementSubmit` 根据 `needSignature` 判断是否需要继续到签名选择。


## Risks / Trade-offs

- 目前后端只记录日志，若前端提前依赖真实写入会造成功能缺失 → 需在 UI 提示中标注“测试阶段”。
- `signatureId` 仍为加密 key，后端尚未验证其有效性 → 必须在后续版本补充校验。
- 旧版前端若未更新仍会遗漏授权分支，但新路由为幂等 no-op（只日志），风险可控。