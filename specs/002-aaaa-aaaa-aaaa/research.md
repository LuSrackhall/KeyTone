# Research: 键音专辑签名系统

## Decisions

- Decision: 统一文件交互通过 Go SDK 实现（除导入/导出对话框）
  - Rationale: 避免前端直接文件系统操作，保持架构边界与跨平台一致性
  - Alternatives: 前端直接 Node/Electron FS → Rejected（违背宪章架构分离）

- Decision: 签名导出/导入使用浏览器文件选择/保存（.ktsign）
  - Rationale: 满足用户交互习惯与安全提示；避免路径硬编码
  - Alternatives: SDK 写死路径 → Rejected（不利于用户控制与权限提示）

- Decision: 保护码由系统自动生成，导入不输入保护码
  - Rationale: 降低用户心智负担；保护码用于校验与去重
  - Alternatives: 用户手输 → Rejected（易出错，体验差）

- Decision: E2E 使用 Playwright
  - Rationale: 跨浏览器、端到端流程验证；仓库已有 MCP 支撑
  - Alternatives: 仅单元测试 → Rejected（不足以覆盖对话框/导入/导出全流程）

## Open Questions (Resolved)

- Clarifications: 用户明确给出 override，允许在无 Clarifications 节下推进 /plan
