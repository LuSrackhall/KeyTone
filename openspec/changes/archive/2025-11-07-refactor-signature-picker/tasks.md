# Tasks: 优化签名选择对话框

## 任务列表

### Task 1: 更新 SignaturePickerDialog 布局


**目标**：将卡片布局从竖向（图片顶部）改为横向（图片左侧）

**工作内容**：
- 修改 `.signature-card` 的 flex 方向为 `row`
- 左侧图片区固定宽度 60px，`flex-shrink: 0`
- 右侧信息区使用 `col` 和 `flex-grow`
- 调整 padding 和边距使卡片紧凑
- 修改选中指示器位置（可改为右侧 top-right 角）

**验证**：
- 视觉检查卡片高度约 70px（含 padding）
- 图片在左侧 60x60 展示
- 名称和介绍在右侧，单行名称 + 最多 2 行介绍

**依赖**：无

---

### Task 2: 导入并集成真实签名数据源


**目标**：从签名管理 store 获取真实签名数据而非 mock

**工作内容**：
- 导入 `useSignatureStore`
- 导入 `getSignaturesList`, `decryptSignatureData`, `getImageUrl` 从 `boot/query/signature-query`
- 在 `setup()` 中添加 `loadSignaturesRealtime()` 函数
- 函数逻辑：
  - 调用 `getSignaturesList()` 获取加密列表
  - 逐个解密并构建 Signature 对象
  - 为每个签名获取图片 URL（异步）
  - 按 `sort.time` 排序
  - 赋值给响应式变量 `localSignatures`

**验证**：
- console 打印签名加载过程无异常
- 对话框显示的签名名称与签名管理页面一致
- 图片 URL 正确加载

**依赖**：Task 1（布局准备好后才能验证内容展示）

---

**依赖**：Task 2

---

### Task 3: 在对话框打开时触发数据加载


**目标**：监听 `visible` prop，打开时加载真实数据，关闭时清理

**工作内容**：
- 添加 `watch(() => props.visible)` 监听器
- 当 `visible` 变为 true 时：
  - 若 `props.signatures` 为空或未提供，调用 `loadSignaturesRealtime()`
  - 显示 loading spinner（可选）
- 当 `visible` 变为 false 时：
  - 清理图片 Blob URL
  - 注销 SSE 监听（如果已监听）
- 支持降级：若传入 `signatures` prop，直接使用（用于 mock 测试）

**验证**：
- 打开对话框，等待数据加载（不应感到卡顿）
- 关闭对话框，不会内存泄漏（Blob URL 清理）
- 重复打开/关闭，数据正确

**依赖**：Task 3

---

### Task 4: SSE 增量更新支持


**目标**：对话框打开时，监听后端配置变更，自动更新签名列表

**工作内容**：
- 在 `loadSignaturesRealtime()` 完成后，注册 SSE 回调
- 回调函数 `handleSseUpdate()`：
  - 重新获取和解密签名列表
  - 更新本地 `localSignatures` 数据
  - 保持当前搜索状态和选中项（如果签名仍存在）
- 对话框关闭时注销此回调

**验证**：
- 打开对话框后，在签名管理页新建签名，观察选择对话框自动刷新并显示新签名
- 删除签名后，选择对话框自动移除该项
- 修改排序后，对话框中的顺序同步更新

**依赖**：Task 3

---

### Task 5: 功能集成和端到端测试

**目标**：整合所有修改，验证导出流程端到端可用

**工作内容**：
- 移除 `Keytone_album_page_new.vue` 中的 mock signatures 数组
- 更新对话框调用处，仍传入空数组或不传 `signatures` prop
- 端到端测试：
  1. 打开签名管理页，创建或查看几个签名
  2. 打开导出流程，到达签名选择对话框
  3. 确认显示的签名与签名管理页一致
  4. 搜索功能正常（按名称筛选）
  5. 选择一个签名，继续导出流程
  6. 在签名管理页修改签名顺序，再打开导出流程，验证选择对话框中的顺序同步

**验证**：
- 所有 5 个测试点通过
- 无 console 错误
- 导出流程其他对话框（授权策略、授权门控）不受影响

**依赖**：Task 1-4 都应完成

---

## 优先级与并行性

- **顺序依赖**：Task 1 → Task 2 → Task 3 → Task 4 → Task 5
  - Task 1 是 UI 基础
  - Task 2 依赖 Task 1 的布局才能正确展示
  - Task 3 依赖 Task 2 的数据加载函数
  - Task 4 依赖 Task 3 的监听机制
  - Task 5 是整体验收

## 完成标准

- [x] Task 1: 布局改为横向，卡片高度 ~70px，图片左侧 60x60
- [x] Task 2: 真实签名数据加载，解密逻辑正确，console 无错误
- [x] Task 3: 打开/关闭对话框，数据和资源正确管理
- [x] Task 4: SSE 回调注册，后端更新自动同步，排序保持一致
- [x] Task 5: 端到端测试通过，移除 mock，所有流程正常

**全部完成日期**：2025-11-06

## 回滚方案

若某个任务出现问题：
- 重置为 `git checkout -- frontend/src/components/export-flow/SignaturePickerDialog.vue`
- 或针对性注释问题代码，恢复到上一个稳定 Task
