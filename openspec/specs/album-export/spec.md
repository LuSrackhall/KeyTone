# 专辑导出签名流程规格说明

## Purpose

描述 `useExportSignatureFlow` 及配套对话框在专辑导出前指导用户确认签名、授权和选择签名的交互要求，确保导出流程一致且可扩展。

## Requirements

### Requirement: 导出状态机编排

Normative: The composable `useExportSignatureFlow` SHALL 根据专辑状态按顺序驱动签名确认、授权策略、授权门控及签名选择；取消任何步骤 MUST 终止本次导出并重置状态。

#### Scenario: 无签名且无需签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“无需签名”并继续
- **THEN** 状态机立即进入 `done`，调用方可直接触发真实导出

#### Scenario: 无签名且需要签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“需要签名”
- **THEN** 状态机依次展示“授权策略选择 → 风险确认 → 联系方式填写”，完成后进入签名选择

#### Scenario: 已有签名记录

- **GIVEN** 专辑已存在签名
- **WHEN** 用户点击导出
- **THEN** 状态机跳过“签名确认”，根据策略直接进入授权门控（如需）或签名选择

---

### Requirement: 签名确认对话框

Normative: The confirmation dialog SHALL 仅在专辑无签名时显示，默认选中“需要签名”，并允许用户选择“无需签名”；关闭对话框 MUST 恢复空闲状态。

#### Scenario: 选择需要签名

- **GIVEN** 对话框可见且默认选中“需要签名”
- **WHEN** 用户点击“继续”
- **THEN** 对话框关闭并通知状态机进入授权策略步骤

#### Scenario: 选择无需签名

- **GIVEN** 对话框可见
- **WHEN** 用户切换到“无需签名”并点击“继续”
- **THEN** 对话框关闭并直接完成导出前置流程

---

### Requirement: 授权策略设置对话链

Normative: The system SHALL 提供授权策略对话（默认推荐无需授权）、风险提示对话及联系方式收集对话；当用户选择需要授权时，邮箱输入 MUST 校验格式且为空时禁用继续。

#### Scenario: 保持无需授权

- **GIVEN** 授权策略对话默认选择“无需授权”
- **WHEN** 用户直接点击“继续”
- **THEN** 状态机跳过风险提示与联系方式步骤，进入签名选择

#### Scenario: 需要授权并填写联系方式

- **GIVEN** 用户切换到“需要授权”并确认风险提示
- **WHEN** 用户在联系方式对话中输入有效邮箱（可选附加信息）并点击继续
- **THEN** 状态机记录联系方式并进入签名选择

---

### Requirement: 授权门控对话框

Normative: The authorization gate dialog SHALL 在策略要求授权且会话尚未授权时显示作者联系方式、支持复制操作，并在导入授权文件前禁用“继续”按钮。

#### Scenario: 导入授权后继续

- **GIVEN** 授权门控对话框可见且尚未导入文件
- **WHEN** 用户选择授权文件并确认
- **THEN** “继续”按钮启用，点击后关闭对话框并进入签名选择

#### Scenario: 取消授权门控

- **GIVEN** 授权门控对话框可见
- **WHEN** 用户点击“取消”
- **THEN** 对话框关闭，状态机回到 `idle`，导出流程终止

---

### Requirement: 签名选择对话框

Normative: The picker dialog SHALL 展示由父组件提供的签名列表（名称、简介、缩略图），支持名称搜索、空态提示，并提供“新建签名”入口调用签名创建对话框。

#### Scenario: 从列表中选择

- **GIVEN** 传入的签名数组非空
- **WHEN** 用户点击某个签名卡片并确认
- **THEN** 对话框关闭并将选中签名 ID 返回状态机

#### Scenario: 列表为空

- **GIVEN** 传入签名数组为空
- **WHEN** 对话框显示空态提示
- **THEN** 用户可点击“新建签名”触发父级打开签名创建对话框

---

### Requirement: 导出组件目录约束

Normative: All export-flow UI components SHALL 存放于 `frontend/src/components/export-flow/`，状态机逻辑 MUST 位于 `useExportSignatureFlow.ts` 并保持文件头注释说明用途；后续新增步骤组件 SHALL 复用该目录与状态机对接。

#### Scenario: 新增导出组件

- **GIVEN** 开发者实现新的导出步骤对话框
- **WHEN** 组件被提交
- **THEN** 文件路径位于 `components/export-flow/` 并通过 `useExportSignatureFlow` 集成

#### Scenario: 维护状态机

- **GIVEN** 开发者阅读 `useExportSignatureFlow.ts`
- **WHEN** 打开文件
- **THEN** 能看到文件头注释说明状态机用途，从而理解其为生产级编排逻辑

---

### Requirement: 可访问性与反馈

Normative: All export flow dialogs SHALL 支持键盘导航（Tab/Shift+Tab、Enter、Esc），关键操作 MUST 提供禁用或加载态，输入校验错误 MUST 就地提示。

#### Scenario: 键盘操作

- **GIVEN** 任意导出流程对话框处于焦点状态
- **WHEN** 用户使用键盘在控件之间切换并按 Enter/Esc
- **THEN** 主按钮与取消行为与点击操作一致

#### Scenario: 校验提示

- **GIVEN** 用户在联系方式对话中勾选需要授权但邮箱为空
- **WHEN** 系统执行校验
- **THEN** “继续”按钮保持禁用并显示邮箱必填提示

