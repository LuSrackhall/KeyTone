# 签名选择对话框优化

## MODIFIED Requirements

### Requirement: 签名选择对话框布局与数据源

Normative: The signature picker dialog SHALL display signatures in a horizontal row layout (left image 60x60px, right text area) and automatically load real signature data from the signature store instead of relying on mock prop data; when visible, the dialog SHALL fetch decrypted signatures and listen for SSE updates.

#### Scenario: 优化后的紧凑布局

- **GIVEN** 签名选择对话框打开且有多个签名
- **WHEN** 对话框渲染签名卡片列表
- **THEN** 每张卡片采用行布局：左侧 60x60 图片、右侧名称和介绍（2行内）、右上角选中指示器，整体高度约 70px

#### Scenario: 加载真实签名数据

- **GIVEN** 对话框的 `visible` prop 变为 true
- **WHEN** 对话框未传入 `signatures` prop（或传入空数组）
- **THEN** 自动调用 `getSignaturesList()` 获取加密签名、逐项解密、按 `sort.time` 排序，显示在列表中

#### Scenario: SSE 增量更新

- **GIVEN** 对话框已打开且显示了真实签名列表
- **WHEN** 后端签名数据变更（新增/修改/删除/排序），触发 SSE 推送
- **THEN** 对话框自动重新加载签名数据，保持现有的搜索和选中状态，用户无需手动刷新

#### Scenario: 降级到 Mock 模式

- **GIVEN** 父组件传入非空 `signatures` prop
- **WHEN** 对话框打开
- **THEN** 使用传入的 mock 数据而非加载真实数据（用于测试场景）

---

### Requirement: 图片 URL 管理与清理

Normative: The dialog SHALL manage Blob URLs for signature images, fetching them asynchronously during list load and revoking them when the dialog closes to prevent memory leaks.

#### Scenario: 异步加载图片 URL

- **GIVEN** 签名列表包含带图片的项
- **WHEN** 列表渲染
- **THEN** 图片 URL 异步获取不阻塞列表显示，图片逐项加载完成后渲染

#### Scenario: 关闭对话框时清理资源

- **GIVEN** 对话框已显示并加载了多张图片
- **WHEN** 用户关闭对话框（`visible` 变为 false）
- **THEN** 所有 Blob URL 被 revoke，SSE 监听被注销，无内存泄漏

---

### Requirement: 保持向后兼容

Normative: The dialog SHALL preserve the `signatures` prop and event signatures (`select`, `createNew`, `cancel`) for backward compatibility with existing callers.

#### Scenario: 事件接口不变

- **GIVEN** 父组件调用 SignaturePickerDialog
- **WHEN** 用户选择签名或点击新建
- **THEN** 触发的 `select` 和 `createNew` 事件与原实现一致

#### Scenario: Props 接口兼容

- **GIVEN** 父组件仍然传入 `signatures` 和 `visible` props
- **WHEN** 对话框接收
- **THEN** 接口有效，行为如 Scenario 中所述（优先使用 prop 或自动加载）
