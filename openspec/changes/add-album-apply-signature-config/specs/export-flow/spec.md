## MODIFIED Requirements

### Requirement: 导出状态机编排

Normative: The composable `useExportSignatureFlow` SHALL 根据专辑状态按顺序驱动签名确认、授权策略、授权门控及签名选择；无论用户是否选择“需要签名”，流程都 MUST 在完成授权策略判断后才进入 `done`，取消任何步骤 MUST 终止本次导出并重置状态。

#### Scenario: 无签名且无需签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“无需签名”，并在后续授权策略对话中保持“无需授权”
- **THEN** 状态机跳过风险提示与联系方式步骤，直接进入 `done`，调用方可立即触发导出

#### Scenario: 无签名且需要签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“需要签名”
- **THEN** 状态机依次展示“授权策略选择 → 风险确认 → 联系方式填写”，完成后进入签名选择

#### Scenario: 无签名但需授权

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“无需签名”，但在后续授权策略对话中勾选“需要授权”
- **THEN** 状态机按照需要授权的路径继续（风险确认 → 联系方式 → 签名选择），并在签名选择完成后进入 `done`

#### Scenario: 已有签名记录

- **GIVEN** 专辑已存在签名
- **WHEN** 用户点击导出
- **THEN** 状态机跳过“签名确认”，根据策略直接进入授权门控（如需）或签名选择

---

### Requirement: 签名确认对话框

Normative: The confirmation dialog SHALL 仅在专辑无签名时显示，默认选中“需要签名”，并允许用户选择“无需签名”；关闭对话框 MUST 恢复空闲状态或进入授权策略对话。

#### Scenario: 选择需要签名

- **GIVEN** 对话框可见且默认选中“需要签名”
- **WHEN** 用户点击“继续”
- **THEN** 对话框关闭并通知状态机进入授权策略步骤

#### Scenario: 选择无需签名

- **GIVEN** 对话框可见
- **WHEN** 用户切换到“无需签名”并点击“继续”
- **THEN** 对话框关闭且状态机进入“二次创作是否需要授权”对话，而非直接完成

---

### Requirement: 授权策略设置对话链

Normative: The system SHALL 提供授权策略对话（默认推荐无需授权）、风险提示对话及联系方式收集对话；当用户选择需要授权时，邮箱输入 MUST 校验格式且为空时禁用继续；当 `needSignature=false` 且用户最终选择“无需授权”时，状态机 MUST 直接完成而不进入签名选择。

#### Scenario: 保持无需授权

- **GIVEN** 授权策略对话默认选择“无需授权”
- **WHEN** 用户直接点击“继续”
- **THEN**
  - 若当前流程 `needSignature=true`，状态机跳过风险提示与联系方式步骤，进入签名选择；
  - 若 `needSignature=false`，状态机直接进入 `done`，允许导出

#### Scenario: 需要授权并填写联系方式

- **GIVEN** 用户切换到“需要授权”并确认风险提示
- **WHEN** 用户在联系方式对话中输入有效邮箱（可选附加信息）并点击继续
- **THEN** 状态机记录联系方式并进入签名选择

#### Scenario: 无签名但需要授权

- **GIVEN** 用户在签名确认对话中选择“无需签名”
- **WHEN** 用户随后在授权策略对话中选择“需要授权”
- **THEN** 系统仍需展示风险提示、联系方式和签名选择；最终结果携带 `needSignature=false` 但 `requireAuthorization=true`

---

## ADDED Requirements

### Requirement: 签名配置应用路由

Normative: The exporter SHALL 调用新的 `POST /keytone_pkg/apply_signature_config` 路由，并在加密完成与文件导出之间提交如下字段：
- `albumPath`：目标专辑目录；
- `needSignature`：布尔值，用于区分真实签名写入或授权用途；
- `requireAuthorization`：布尔值，指示是否收集联系方式并强制加密；
- `signatureId`：调用该路由时 MUST 提供（因为仅在 `needSignature` 或 `requireAuthorization` 为真时才会触发调用）；
- `contactEmail`、`contactAdditional`：仅在 `requireAuthorization=true` 时提供，其中 `contactEmail` MUST 通过前端校验且后端也需校验非空，`contactAdditional` 为可选补充字段。

#### Scenario: 需要签名

- **GIVEN** 用户在导出流程中选择“需要签名”
- **WHEN** 加密完成且导出尚未执行
- **THEN** 前端调用 `apply_signature_config`，payload 包含签名 ID 与可选授权联系方式；后端记录日志并返回 200 以允许继续导出

#### Scenario: 无需签名但需要授权

- **GIVEN** 用户选择“无需签名”但在授权策略中启用“需要授权”
- **WHEN** 导出流程进入加密阶段
- **THEN** 前端仍需调用 `apply_signature_config` 并携带签名 ID 与授权联系方式；后端根据 `needSignature=false` 记录该场景，随后允许导出
