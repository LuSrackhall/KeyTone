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

---

### Requirement: 隐藏滚动条优化用户体验

Normative: The signature picker dialog's scrollbar on the main card element SHALL be hidden from view while maintaining scrolling functionality to improve visual aesthetics and reduce visual clutter; however, the scrollbar within signature items (for horizontal scrolling of long names/descriptions) SHALL remain visible with refined styling.

#### Scenario: 隐藏主对话框滚动条

- **GIVEN** 用户打开签名选择对话框
- **WHEN** 对话框内容超过可视区域，需要垂直滚动
- **THEN** 主对话框的垂直滚动条被隐藏（使用 `[&::-webkit-scrollbar]:hidden`），但用户仍可通过鼠标滚轮或触控板滚动内容，不影响功能

#### Scenario: 保留卡片内横向滚动条样式

- **GIVEN** 签名卡片中的名称或介绍文本过长超出可视区域
- **WHEN** 卡片内的文本需要水平滚动显示
- **THEN** 卡片内横向滚动区域的滚动条保持可见且美观（4px 宽，半透明深灰色）

---

### Requirement: 改善布局紧凑度与间距美观性

Normative: The dialog layout SHALL optimize spacing between search/create area and signature list to improve visual aesthetics, reduce unnecessary whitespace, and maintain adequate readability and visual balance.

#### Scenario: 优化搜索区域内部间距

- **GIVEN** 签名选择对话框打开
- **WHEN** 用户查看搜索框、描述文本和"新建签名"按钮所占的区域
- **THEN** 搜索区域 padding 从 `q-pa-md`（16px）优化为 `q-pa-sm`（8px）；描述文本、搜索输入框下方 margin 从 `q-mb-md` 优化为 `q-mb-sm`，使搜索区域更紧凑

#### Scenario: 减少签名卡片之间的垂直间隔

- **GIVEN** 签名选择对话框列表包含多个签名
- **WHEN** 对话框渲染签名卡片列表
- **THEN** 卡片之间的垂直间隔从 `q-mb-sm`（8px）优化为 `q-mb-xs`（4px），提高列表紧凑度

#### Scenario: 移除过度底部 padding

- **GIVEN** 用户滚动签名列表到底部
- **WHEN** 观察列表底部空白区域
- **THEN** 移除原先主内容区域硬编码的 `padding-bottom: 100px`，改为自适应布局，仅保留必要的下方操作栏空间
