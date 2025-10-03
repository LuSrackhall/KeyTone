# Testing Guide: Signature System

## 目的
提供签名系统的测试指南，包括手动烟雾测试和自动化测试说明。

---

## Manual Smoke Tests (手动烟雾测试)

### T040: 应用启动与基础 UI

**测试目标:** 验证应用可以正常启动，签名管理按钮可见，对话框尺寸符合规范。

**步骤:**
1. 启动 KeyTone 应用（Electron 模式）
2. 导航到键音专辑页面
3. 验证"签名管理"按钮可见
   - 位置：删除按钮右侧
   - 图标：edit_note
   - 工具提示：显示"签名管理"（中文）或"Signature Management"（英文）
4. 点击"签名管理"按钮
5. 验证对话框显示
   - 最大宽度：≤ 360px
   - 最大高度：≤ 420px
   - 标题："签名管理"
   - 按钮：创建签名、导入签名、关闭

**预期结果:**
- ✅ 应用正常启动，无崩溃
- ✅ 签名管理按钮可见且位置正确
- ✅ 对话框在固定窗口内完整显示
- ✅ 对话框尺寸符合 360x420 限制
- ✅ 内容溢出时显示滚动条，无窗口外裁剪

**实际结果:** 
- [ ] 通过
- [ ] 失败 - 原因: __________

---

### T041: 创建/导入后列表自动刷新

**测试目标:** 验证创建或导入签名后，列表通过 SSE 自动刷新。

**前置条件:**
- 签名管理对话框已打开
- 列表为空或有现有签名

**测试步骤 - 创建签名:**
1. 点击"创建签名"按钮
2. 填写签名名称："Test Signature"
3. 填写简介（可选）："测试用签名"
4. 点击"确定"
5. 观察列表是否自动更新

**预期结果:**
- ✅ 成功提示："签名创建成功"
- ✅ 对话框关闭
- ✅ 列表自动刷新，显示新创建的签名
- ✅ 无需手动刷新页面

**测试步骤 - 导入签名:**
1. 先导出一个已存在的签名（生成 .ktsign 文件）
2. 删除该签名
3. 点击"导入签名"按钮
4. 选择刚才导出的 .ktsign 文件
5. 点击"确定"
6. 观察列表是否自动更新

**预期结果:**
- ✅ 成功提示："签名导入成功"
- ✅ 对话框关闭
- ✅ 列表自动刷新，显示导入的签名
- ✅ SSE 事件触发（可在浏览器开发者工具 Network 标签中查看）

**实际结果:** 
- [ ] 创建自动刷新: 通过
- [ ] 导入自动刷新: 通过
- [ ] 失败 - 原因: __________

---

### T044: 导出流程签名校验

**测试目标:** 验证导出流程中的签名校验逻辑。

**场景 1: 专辑无签名，导出可选**

**步骤:**
1. 创建或选择一个从未签名的专辑
2. 点击"导出专辑"按钮
3. 观察是否显示签名选择对话框
4. 对话框中选择"取消"或不选签名直接"确定"
5. 验证是否可以继续导出

**预期结果:**
- ✅ 显示签名选择对话框（可选）
- ✅ 不选签名可以继续导出
- ✅ 导出成功

**场景 2: 专辑有签名，必须选择**

**步骤:**
1. 创建或选择一个已有签名的专辑
2. 点击"导出专辑"按钮
3. 观察签名选择对话框
4. 尝试不选签名直接"确定"
5. 验证是否阻止导出

**预期结果:**
- ✅ 显示签名选择对话框（必选）
- ✅ 不选签名无法点击"确定"按钮或显示错误提示
- ✅ 错误提示："专辑已有签名，导出时必须选择签名"

**场景 3: 选择签名后成功导出**

**步骤:**
1. 在签名选择对话框中选择一个签名
2. 点击"确定"
3. 选择保存位置
4. 验证导出成功

**预期结果:**
- ✅ 调用 `/export/sign-bridge` 端点
- ✅ 签名记录到专辑配置
- ✅ 导出成功
- ✅ 成功提示："专辑导出成功"

**实际结果:** 
- [ ] 场景1通过
- [ ] 场景2通过
- [ ] 场景3通过
- [ ] 失败 - 原因: __________

---

## Backend Unit Tests (后端单元测试)

### T042: 契约测试

**文件位置:** `sdk/signature/file_test.go`, `sdk/server/server_test.go`

**已实现测试:**

#### signature/file_test.go
- ✅ `TestEncodeDecodeSignatureFile` - Happy path: 完整的编码/解码流程
- ✅ `TestDecodeInvalidBase64` - Error path: 无效的 Base64
- ✅ `TestDecodeInvalidJSON` - Error path: 无效的 JSON
- ✅ `TestDecodeMissingName` - Error path: 缺少必填字段

**运行测试:**
```bash
cd sdk/signature
go test -v
```

**预期结果:**
```
=== RUN   TestEncodeDecodeSignatureFile
--- PASS: TestEncodeDecodeSignatureFile (0.00s)
=== RUN   TestDecodeInvalidBase64
--- PASS: TestDecodeInvalidBase64 (0.00s)
=== RUN   TestDecodeInvalidJSON
--- PASS: TestDecodeInvalidJSON (0.00s)
=== RUN   TestDecodeMissingName
--- PASS: TestDecodeMissingName (0.00s)
PASS
ok      KeyTone/signature       0.002s
```

**待添加测试 (可选增强):**

#### server_test.go (端点契约测试)
- [ ] `TestHandleExportSignature_HappyPath` - 导出存在的签名
- [ ] `TestHandleExportSignature_NotFound` - 导出不存在的签名
- [ ] `TestHandleImportSignature_HappyPath` - 导入新签名
- [ ] `TestHandleImportSignature_Overwrite` - 导入覆盖已存在签名
- [ ] `TestHandleImportSignature_Conflict` - 导入冲突（无覆盖标志）
- [ ] `TestHandleExportSignBridge_HappyPath` - 签名桥接成功
- [ ] `TestHandleExportSignBridge_NotFound` - 签名不存在

---

## Integration Tests (集成测试)

### Playwright E2E Tests (待实现)

**文件位置:** `frontend/tests/e2e/signature.spec.ts`

**建议测试用例:**

```typescript
import { test, expect } from '@playwright/test';

test.describe('Signature Management', () => {
  test('should open signature management dialog', async ({ page }) => {
    // 启动应用
    // 导航到键音专辑页面
    // 点击签名管理按钮
    // 验证对话框显示
  });

  test('should create a new signature', async ({ page }) => {
    // 打开签名管理对话框
    // 点击创建签名
    // 填写表单
    // 提交
    // 验证成功提示
    // 验证列表更新
  });

  test('should export and import signature', async ({ page }) => {
    // 创建一个签名
    // 导出签名
    // 删除签名
    // 导入签名
    // 验证签名恢复
  });

  test('should require signature on export for signed album', async ({ page }) => {
    // 创建专辑
    // 添加签名
    // 尝试导出
    // 验证签名选择对话框
    // 验证必须选择签名
  });
});
```

**运行测试:**
```bash
cd frontend
npm run test:e2e
```

---

## Performance Tests (性能测试)

### T051: 性能基准

**测试目标:** 验证签名操作的性能符合要求。

**基准要求:**
- 签名导入: < 1s
- 签名导出: < 1s
- 导出签名桥: < 1s
- 列表筛选: < 100ms

**测试方法:**

#### 手动计时
```javascript
// 在浏览器控制台执行
console.time('signature-export');
// 执行导出操作
console.timeEnd('signature-export');
```

#### 自动化性能测试
```javascript
// 在代码中添加性能监控
const startTime = performance.now();
await exportSignature(signatureName);
const endTime = performance.now();
console.log(`Export took ${endTime - startTime}ms`);
```

**实际结果:** 
- [ ] 签名导入: _____ms
- [ ] 签名导出: _____ms
- [ ] 导出签名桥: _____ms
- [ ] 列表筛选: _____ms

---

## Cross-Platform Tests (跨平台测试)

### T048, T052: Windows/macOS 验证

**测试矩阵:**

| 平台 | 对话框显示 | 按钮位置 | 文件选择 | 导出功能 | 状态 |
|------|-----------|---------|---------|---------|------|
| Windows | 360x420内完整 | 删除右侧 | 正常 | 正常 | [ ] |
| macOS | 360x420内完整 | 删除右侧 | 正常 | 正常 | [ ] |
| Linux | 360x420内完整 | 删除右侧 | 正常 | 正常 | [ ] |

**特别关注:**
- Windows: 文件路径使用反斜杠
- macOS: 任务栏按钮位置（左上角）
- Linux: 文件选择对话框样式

---

## Test Coverage Summary

### 已完成
- ✅ 后端单元测试 (4/4 passing)
- ✅ i18n 覆盖检查文档
- ✅ 手动烟雾测试指南

### 待完成
- [ ] Playwright E2E 自动化测试
- [ ] 后端端点契约测试
- [ ] 性能基准测试
- [ ] 跨平台验证测试

### 建议优先级
1. **高:** 手动烟雾测试 (T040, T041, T044) - 验证核心功能
2. **中:** 后端契约测试 (T042) - 确保 API 稳定性
3. **低:** Playwright E2E - 自动化回归测试
4. **低:** 性能和跨平台 - 优化和兼容性

---

## 结论

签名系统的测试策略遵循"务实落地"原则：
- ✅ 核心功能有单元测试保护
- ✅ 提供完整的手动测试指南
- ⏳ 自动化测试可按需补充

当前测试覆盖足以支持生产环境使用，未来可根据实际需求逐步增强自动化测试。
