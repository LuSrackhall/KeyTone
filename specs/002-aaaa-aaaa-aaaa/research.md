# Research: 键音专辑签名系统（Phase 0）

## Decisions

- Decision: 采用混合测试策略（契约/核心逻辑测试优先 + UI 少量烟雾与组件单测）
  - Rationale: Electron + 本地后端导致全链路 TDD 成本高；先保证契约稳定和关键逻辑正确，UI 以手动验证为主、少量自动化兜底。
  - Alternatives: 严格端到端 TDD（高成本、脆弱）；仅手动测试（缺少回归保护）。

- Decision: REST 契约为前后端边界（gin 提供 /sdk/*），为后续 Wails 版本迁移保持协议稳定
  - Rationale: 与现有项目一致，易于编写契约测试与前端 axios 封装。
  - Alternatives: IPC/bridge 专用协议；JSON-RPC（后续可评估）。

- Decision: 签名文件（.ktsign）内容第一阶段最小化
  - Rationale: 优先实现导入/导出流转；校验逻辑逐步增强。
  - Alternatives: 一次性定义完整加密/授权体系（超范围）。

## Unknowns Resolved

- 专辑内签名结构（最小）：签名者 uuid、名称、创建时间、保护码哈希、导出时间戳数组、允许二次导出（bool）。
- .ktsign 文件：JSON 打包 + 名片图片引用（相对路径或内嵌 base64，优先相对路径）。第一阶段可不做强校验，仅保留哈希与时间戳。
- 错误码规范：HTTP 状态 + JSON { code, message }，重复/冲突时返回 409；参数校验 400。
- Electron 集成：前端仅经 REST 调用；文件保存/选择使用浏览器端接口或由后端返回目标路径后存储。
- 动态端口/运行时发现：Electron 主进程拉起 SDK 并输出 KEYTONE_PORT，渲染进程通过 `window.myWindowAPI.getBackendPort()` 获取端口并调用 `UpdateApi(port)`；必要时刷新以重建 SSE 链接。

## Best Practices

- gin：使用 binding 校验与统一错误响应中间件；路由前缀 /sdk/ 收敛。
- axios：统一封装 baseURL（避免到处硬编码），请求/响应拦截器统一错误提示；超时与重试策略最小化，避免 UI 卡顿。
- Playwright（Electron 开发流）：仅保留 1-2 条烟测（加载页面、签名管理入口存在、后端探针 200）。失败时提供 HTML report 与截图/视频。

## Notes

- 后续增强：签名校验与授权策略、跨平台路径与文件对话框封装、性能基线与回归。
