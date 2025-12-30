# 无需授权场景可选联系方式规格（增量）

## Purpose

本增量规格描述在用户选择"需要签名"但"无需授权"时，增加可选联系方式填写步骤的要求，补充 `openspec/specs/export-flow/spec.md` 中的授权策略设置对话链。

---

### Requirement: 可选联系方式对话框

Normative: The `OptionalContactDialog` SHALL 在用户选择"无需授权"后显示，允许用户可选地填写邮箱和附加联系方式；用户可直接跳过或填写后继续进入签名选择。

#### Scenario: 显示可选联系方式对话框

- **GIVEN** 用户在授权策略对话中选择"无需授权"
- **WHEN** 用户点击"继续"
- **THEN** 显示可选联系方式对话框，而非直接进入签名选择

#### Scenario: 填写有效联系方式

- **GIVEN** 可选联系方式对话框可见
- **AND** 用户输入了有效格式的邮箱
- **WHEN** 用户点击"保存并继续"
- **THEN** 联系方式被保存到流程数据，对话框关闭，进入签名选择

#### Scenario: 跳过联系方式填写

- **GIVEN** 可选联系方式对话框可见
- **WHEN** 用户点击"跳过"
- **THEN** 不保存联系方式，对话框关闭，直接进入签名选择

#### Scenario: 仅填写附加联系方式

- **GIVEN** 可选联系方式对话框可见
- **AND** 用户未填写邮箱但填写了附加联系方式
- **WHEN** 用户点击"保存并继续"
- **THEN** 仅附加联系方式被保存，对话框关闭，进入签名选择

#### Scenario: 邮箱格式校验

- **GIVEN** 用户在可选联系方式对话框中输入了邮箱
- **AND** 邮箱格式不正确
- **WHEN** 用户点击"保存并继续"
- **THEN** 显示邮箱格式错误提示，"保存并继续"按钮被禁用

#### Scenario: 取消可选联系方式对话框

- **GIVEN** 可选联系方式对话框可见
- **WHEN** 用户关闭对话框（非通过"跳过"或"保存并继续"）
- **THEN** 对话框关闭，状态机回到 `idle`，导出流程终止

---

### Requirement: 联系方式数据复用

Normative: The optional contact information SHALL 使用与"需要授权"场景相同的数据字段（`contactEmail` 和 `contactAdditional`）存储，确保数据结构一致性。

#### Scenario: 联系方式存储

- **GIVEN** 用户在可选联系方式对话框中填写了联系方式
- **WHEN** 流程完成并导出专辑
- **THEN** 联系方式存储在专辑签名配置的 `authorization.contactEmail` 和 `authorization.contactAdditional` 字段

#### Scenario: 签名信息展示

- **GIVEN** 专辑已导出且包含联系方式（无论是否需要授权）
- **WHEN** 用户查看专辑签名信息
- **THEN** 联系方式在原始作者区块中正确显示
