# 签名管理 SSE 闪烁问题 - 快速总结

## 问题

签名管理页面中，SSE 触发时整个列表闪烁，用户体验差。

## 根本原因

**旧实现**：SSE 触发 → `handleSseUpdate()` → `loadSignatures()` → **整个数组被替换** → 整个列表 DOM 销毁重建 → 闪烁

## 解决方案

**新实现**：SSE 触发 → `handleSseUpdate()` → **增量比对** → **只更新变化的项** → 无闪烁

## 核心改动

### 1. 新增公共函数

```typescript
// 解密并构建签名 Map
async function decryptAndBuildSignatureMap(
  encryptedSignatures: Record<string, any>
): Promise<Map<string, Signature>>

// 从加密数据提取排序后的列表
function extractSortedSignatures(
  encryptedSignatures: Record<string, any>,
  signatureMap: Map<string, Signature>
): Signature[]

// 执行增量更新
function updateSignaturesIncremental(
  newSignaturesMap: Map<string, Signature>
)
```

### 2. 修改的函数

**`loadSignatures()`**
- 改为使用公共函数
- 完整重新加载（页面初次加载）

**`handleSseUpdate()`**
- 改为使用增量更新（避免闪烁）
- 添加异常降级到全量加载

## 增量更新逻辑

```
比对新旧数据
    ↓
检测删除项 → 从列表中过滤删除
检测更新项 → 使用 Object.assign 原地修改
检测添加项 → push 新项到列表末尾
    ↓
保持数组引用和元素引用稳定
    ↓
Vue 仅更新实际变化的项
    ↓
用户看不到闪烁
```

## 关键优势

| 方面     | 旧方案             | 新方案                |
| -------- | ------------------ | --------------------- |
| DOM 操作 | 销毁+重建所有元素  | 只修改变化的元素      |
| 闪烁     | 明显               | 无                    |
| 性能     | 差（频繁重排重绘） | 优（最小化 DOM 操作） |
| 状态保留 | 丢失               | 保留                  |
| 代码复用 | 低                 | 高                    |

## 文件改动

- **修改**：`frontend/src/pages/Signature_management_page.vue`
  - 新增 3 个函数
  - 修改 2 个函数
  - 代码量：约 250 行新增/修改

## 兼容性

- ✅ 保留旧格式兼容
- ✅ 保留新格式支持
- ✅ 异常自动降级
- ✅ 保证最终一致性

## 测试要点

1. ✅ 创建新签名 → 应添加到列表末尾，其他项不闪烁
2. ✅ 编辑现有签名 → 只该项更新，其他项不动
3. ✅ 删除签名 → 平滑移除，其他项保持位置
4. ✅ 多窗口同步 → SSE 更新时其他窗口不闪烁
5. ✅ 大量签名 → 性能良好，响应迅速

## 文档参考

详见：`SOLUTION_FLICKER_ISSUE.md`（完整技术分析）
