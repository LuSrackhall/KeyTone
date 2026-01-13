# 任务清单：提升签名名片图片清理的安全性

## 1. 规格与校验

- [x] 1.1 完成 spec delta：更新 `signature-management` 关于 `CleanupOrphanCardImages` 的安全约束与场景
- [x] 1.2 运行 `openspec validate update-signature-image-cleanup-safety --strict` 并修复全部问题

## 2. SDK 后端实现（Go）

- [x] 2.1 在 `sdk/signature/signature.go` 中实现“保守删除”策略：检测到任意解密/解析失败则跳过删除
- [x] 2.2 增加日志：
  - [x] 记录签名条目总数、成功解密数、失败数
  - [x] 当跳过删除时输出 warn（包含失败原因统计）
- [ ] 2.3（可选）为 `CleanupOrphanCardImages` 添加 Go 单元测试（覆盖：全量失败、部分失败、全量成功）

## 3. 验证

- [ ] 3.1 手工验证：
  - [ ] 构造一个包含签名条目的配置，但使其解密失败（模拟密钥不兼容），确认不会删除 `ConfigPath/signature/` 中任意文件
  - [ ] 构造解密成功且存在孤立文件的场景，确认仅删除孤立文件
- [ ] 3.2 确认启动任务仍可运行且不会影响其他功能（仅日志变化与删除策略变化）
