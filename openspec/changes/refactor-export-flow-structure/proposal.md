# Proposal: refactor-export-flow-structure

## Why

- 导出签名流程的 UI 组件分散在 `components/` 根目录中，可读性差。
- `useExportSignatureFlow` 的职责不够直观，新同事难以分辨其是否为临时脚本。

## What Changes

- 将导出签名流程相关组件集中到 `components/export-flow/` 目录。
- 为 `useExportSignatureFlow` 增加文件级注释，明确其作为生产级状态机的角色。
- 更新规范，记录组件归档和文档化要求。

## Impact

- Affected specs: album-export
- Affected code: frontend/src/components/export-flow/*, frontend/src/composables/useExportSignatureFlow.ts
