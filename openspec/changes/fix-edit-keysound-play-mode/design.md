# Design: 修复编辑播放模式不生效

## Context

编辑已有高级键音的对话框读取 `keySoundValue.down/up.mode` 时，UI 组件预期为字符串，但保存逻辑与校验逻辑误以为该字段为 `{ mode: string }` 结构，导致保存时写回 `undefined`。

## Goals / Non-Goals

- Goals:
  - 保证编辑对话框保存播放模式时写回字符串模式值
  - 兼容可能存在的旧结构（对象形式）避免断言失败
  - 在代码中明确数据结构与转换原因
- Non-Goals:
  - 不改变播放模式枚举范围（仍为 single/random/loop）
  - 不调整其它与按键音无关的 UI 或数据结构

## Decisions

- 决定在编辑对话框中使用“模式解析函数”以兼容字符串/对象两种结构。
- 决定在父组件中对编辑态数据进行规范化（将 mode 统一为字符串），减少 UI 层误判。

## Risks / Trade-offs

- 需要同时更新模板校验与保存逻辑，避免出现“UI 显示正确但保存错误”的二次回归。

## Migration Plan

- 无需迁移配置文件；仅在前端读取/编辑路径上做兼容处理。

## Open Questions

- 无。当前问题明确，修复后可回归验证。
