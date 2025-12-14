# Proposal: Add Album Export Signature Flow

## Why

为“导出键音专辑”流程新增签名相关的前端 UI 与交互步骤：
- 在导出前提供“是否需要签名”的选择入口（仅首次导出可选跳过签名）。
- 首次导出时可设置“二次创作需授权”的策略，并要求填写联系方式用于后续授权提示。
- 在需要签名的场景里，提供“选择签名”的对话框，支持从签名管理中新建签名的入口。
- 当策略要求授权且当前未获授权时，在签名选择前显示“授权门控”对话框，引导导入授权文件并展示原作者联系方式。

以上仅提供 UI 与前端交互编排，不落地真实业务逻辑（鉴权、文件校验、落库等后续实现）。

## What Changes

- 新增导出前“签名策略与要求”对话框：
  - 首次导出（无签名痕迹）时可选择“需要签名/不需要签名”。
  - 可勾选“二次创作必须通过当前签名作者授权”，勾选后强制填写联系方式（打包入配置）。
  - 非首次导出（已有签名）不允许选择“不需要签名”。
- 新增“授权门控”对话框：当专辑要求授权且当前未授权时，先展示该对话框（显示原作者联系方式、导入授权文件按钮）再进入签名选择。
- 新增“签名选择”对话框：
  - 拉取签名管理页面已有数据进行选择。
  - 内置“新建签名”入口，弹出签名创建对话框（沿用签名管理 UI）。
- 交互顺序编排：导出按钮 →（条件）签名策略 →（条件）授权门控 → 签名选择 → 继续导出。

## Impact

- Affected specs: `specs/album-export/spec.md`（新增能力：导出前签名策略、授权门控、签名选择 UI 与编排）。
- Related changes: 复用 `changes/add-signature-management` 定义的签名概念与 UI 约束。
- Affected code (front-end only, planned):
  - Vue/Quasar 组件：`ExportSignaturePolicyDialog.vue`、`ExportAuthorizationGateDialog.vue`、`SignaturePickerDialog.vue`
  - 交互编排：`useExportSignatureFlow.ts`
  - i18n 与样式：`i18n` 词条、`uno/tailwind` 样式类

## Non-Goals

- 不实现后端接口、不做真实授权校验与签名写入；仅提供可验收的 UI 与交互流程（可用 Mock）。

