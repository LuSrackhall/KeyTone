# Album Export (Maintainability)

## ADDED Requirements

### Requirement: 导出签名流程组件分层与文档化

系统在维护“导出键音专辑”签名流程时 SHALL 将所有对话框/控件组件集中在 `components/export-flow/` 命名空间内，便于识别与复用；负责编排流程的 `components/export-flow/useExportSignatureFlow.ts` MUST 保持文件级注释，说明其为生产级状态机而非临时脚本。任何新增的导出步骤组件 MUST 继续存放在同一目录并复用该状态机。

#### Scenario: 新增导出流程组件

- **WHEN** 新的导出签名对话框或辅助控件被实现
- **THEN** 其文件位于 `components/export-flow/`，并通过统一导出的状态机或现有组件完成集成

#### Scenario: 开发者阅读状态机

- **WHEN** 开发者打开 `useExportSignatureFlow.ts`
- **THEN** 能在文件头部看到关于状态机用途的注释，清晰知晓其为导出流程的正式编排逻辑
