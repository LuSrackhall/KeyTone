# Design (UI-only): Album Export Signature Flow

## Context

- 目标：在“导出专辑”前，串联三个前端对话框（签名策略→授权门控→签名选择）。
- 范围：仅 UI/UX 与交互状态，不落地真实业务（授权校验、签名写入）。
- 依赖：签名管理模块（读取签名列表与创建对话框）。

## Goals / Non-Goals

- Goals
  - 前端对话框组件化与可复用
  - 明确的状态机与编排，便于后续接入真实逻辑
- Non-Goals
  - 后端 API、授权文件格式与加密
  - 持久化落库与实际导出实现

## Flow Orchestration

```text
Export Click
  ├─ Zero-signature? → Yes → Policy Dialog
  │    ├─ Need signature? → Yes → (If require-authorization & not authorized) → Authorization Gate → Signature Picker
  │    └─ Need signature? → No  → Return to caller (continue export)
  └─ Zero-signature? → No  → (If require-authorization & not authorized) → Authorization Gate → Signature Picker
```

状态（简化）：
- idle → policy → auth-gate → picker → done/cancel

## Components

- ExportSignaturePolicyDialog.vue
  - Props: { visible, albumHasSignature, defaultRequireAuth }
  - Emits: { submit: { needSignature, requireAuthorization, contact }, cancel }
- ExportAuthorizationGateDialog.vue
  - Props: { visible, authorContact }
  - Emits: { authorized, cancel }
- SignaturePickerDialog.vue
  - Props: { visible, signatures }
  - Emits: { select(signatureId), createNew, cancel }

## Error Modes

- 表单校验失败：禁用主按钮 + 内联错误提示
- 授权门控未导入：禁用继续按钮
- 签名列表为空：显示空态与“新建”按钮

## Open Questions

- 首次导出勾选“需授权”但选择“无需签名”时，授权主体如何标识？（当前仅存联系方式，后续与业务澄清）
- 授权文件校验与授权态缓存的策略（由后端/业务侧定义）。
