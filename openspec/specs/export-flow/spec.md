# 专辑导出签名流程规格说明

## Purpose

描述 `useExportSignatureFlow` 及配套对话框在专辑导出前指导用户确认签名、授权和选择签名的交互要求，确保导出流程一致且可扩展。
## Requirements
### Requirement: 导出状态机编排

Normative: The composable `useExportSignatureFlow` SHALL 根据专辑状态按顺序驱动签名确认、授权策略、授权门控及签名选择；取消任何步骤 MUST 终止本次导出并重置状态。

#### Scenario: 无签名且无需签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“无需签名”并继续
- **THEN** 状态机立即进入 `done`，调用方可直接触发真实导出

#### Scenario: 无签名且需要签名

- **GIVEN** 专辑尚未添加签名历史
- **WHEN** 用户在“签名确认”选择“需要签名”
- **THEN** 状态机依次展示“授权策略选择 → 风险确认 → 联系方式填写”，完成后进入签名选择

#### Scenario: 已有签名记录

- **GIVEN** 专辑已存在签名
- **WHEN** 用户点击导出
- **THEN** 状态机跳过“签名确认”，根据策略直接进入授权门控（如需）或签名选择

---

### Requirement: 签名确认对话框

Normative: The confirmation dialog SHALL 仅在专辑无签名时显示，默认选中“需要签名”，并允许用户选择“无需签名”；关闭对话框 MUST 恢复空闲状态。

#### Scenario: 选择需要签名

- **GIVEN** 对话框可见且默认选中“需要签名”
- **WHEN** 用户点击“继续”
- **THEN** 对话框关闭并通知状态机进入授权策略步骤

#### Scenario: 选择无需签名

- **GIVEN** 对话框可见
- **WHEN** 用户切换到“无需签名”并点击“继续”
- **THEN** 对话框关闭并直接完成导出前置流程

---

### Requirement: 授权策略设置对话链

Normative: The system SHALL 提供授权策略对话（默认推荐无需授权）、风险提示对话及联系方式收集对话；当用户选择需要授权时，邮箱输入 MUST 校验格式且为空时禁用继续。

#### Scenario: 保持无需授权

- **GIVEN** 授权策略对话默认选择“无需授权”
- **WHEN** 用户直接点击“继续”
- **THEN** 状态机跳过风险提示与联系方式步骤，进入签名选择

#### Scenario: 需要授权并填写联系方式

- **GIVEN** 用户切换到“需要授权”并确认风险提示
- **WHEN** 用户在联系方式对话中输入有效邮箱（可选附加信息）并点击继续
- **THEN** 状态机记录联系方式并进入签名选择

---

### Requirement: 授权门控对话框

Normative: The authorization gate dialog SHALL 在策略要求授权且会话尚未授权时显示作者联系方式、支持复制操作，并在导入授权文件前禁用“继续”按钮。

#### Scenario: 导入授权后继续

- **GIVEN** 授权门控对话框可见且尚未导入文件
- **WHEN** 用户选择授权文件并确认
- **THEN** “继续”按钮启用，点击后关闭对话框并进入签名选择

#### Scenario: 取消授权门控

- **GIVEN** 授权门控对话框可见
- **WHEN** 用户点击“取消”
- **THEN** 对话框关闭，状态机回到 `idle`，导出流程终止

---

### Requirement: 签名选择对话框

Normative: The picker dialog SHALL 展示由父组件提供的签名列表（名称、简介、缩略图），支持名称搜索、空态提示，并提供“新建签名”入口调用签名创建对话框。

#### Scenario: 从列表中选择

- **GIVEN** 传入的签名数组非空
- **WHEN** 用户点击某个签名卡片并确认
- **THEN** 对话框关闭并将选中签名 ID 返回状态机

#### Scenario: 列表为空

- **GIVEN** 传入签名数组为空
- **WHEN** 对话框显示空态提示
- **THEN** 用户可点击“新建签名”触发父级打开签名创建对话框

---

### Requirement: 导出组件目录约束

Normative: All export-flow UI components SHALL 存放于 `frontend/src/components/export-flow/`，状态机逻辑 MUST 位于 `useExportSignatureFlow.ts` 并保持文件头注释说明用途；后续新增步骤组件 SHALL 复用该目录与状态机对接。

#### Scenario: 新增导出组件

- **GIVEN** 开发者实现新的导出步骤对话框
- **WHEN** 组件被提交
- **THEN** 文件路径位于 `components/export-flow/` 并通过 `useExportSignatureFlow` 集成

#### Scenario: 维护状态机

- **GIVEN** 开发者阅读 `useExportSignatureFlow.ts`
- **WHEN** 打开文件
- **THEN** 能看到文件头注释说明状态机用途，从而理解其为生产级编排逻辑

---

### Requirement: 可访问性与反馈

Normative: All export flow dialogs SHALL 支持键盘导航（Tab/Shift+Tab、Enter、Esc），关键操作 MUST 提供禁用或加载态，输入校验错误 MUST 就地提示。

#### Scenario: 键盘操作

- **GIVEN** 任意导出流程对话框处于焦点状态
- **WHEN** 用户使用键盘在控件之间切换并按 Enter/Esc
- **THEN** 主按钮与取消行为与点击操作一致

#### Scenario: 校验提示

- **GIVEN** 用户在联系方式对话中勾选需要授权但邮箱为空
- **WHEN** 系统执行校验
- **THEN** “继续”按钮保持禁用并显示邮箱必填提示

### Requirement: 对话框顶部区域的粘滞定位

Normative: The header and search bar area at the top of the dialog SHALL be fixed to the top using sticky positioning, allowing users to access search and create signature features while scrolling through the signature list without losing the header context.

#### Scenario: 顶部区域粘滞在对话框顶部

- **GIVEN** 用户在签名列表中滚动内容
- **WHEN** 内容向下滚动
- **THEN** Header 和 Search Bar 保持在对话框顶部可见，不跟随滚动

---

### Requirement: 签名选择卡片的视觉设计与交互

Normative: The signature picker SHALL display signature items with enhanced visual hierarchy including micro-interactions (hover shadows, selection glow); users MAY toggle selection by re-clicking an already-selected item; the interface SHALL support horizontal scrolling for long names and descriptions within constrained space.

#### Scenario: 卡片悬停时显示视觉反馈

- **GIVEN** 用户将鼠标悬停在一张签名卡片上
- **WHEN** 鼠标进入卡片区域
- **THEN** 卡片显示增强的阴影（box-shadow: 0 4px 12px rgba(0,0,0,0.08)）并轻微上升（transform: translateY(-1px)），创建交互感

#### Scenario: 卡片被选中时显示光晕效果

- **GIVEN** 用户点击一张签名卡片选中它
- **WHEN** 卡片进入 selected 状态
- **THEN** 卡片显示 2px 蓝色边框（border: 2px solid var(--q-primary)）并带有柔和光晕（box-shadow: 0 0 0 3px rgba(33,150,243,0.1), 0 4px 16px rgba(33,150,243,0.15)），图标周围呈现柔和发光效果

#### Scenario: 重复点击取消选择

- **GIVEN** 用户已选中一张签名卡片
- **WHEN** 用户再次点击同一张卡片
- **THEN** 卡片的选中状态被取消，光晕消退，selectedId 被清空

#### Scenario: 长名称和介绍支持水平滚动

- **GIVEN** 签名的名称或介绍文本过长，超出容器宽度
- **WHEN** 卡片渲染时文字溢出
- **THEN** 名称和介绍分别支持水平滚动而不换行（名称单行，介绍最多 2 行），滚动条宽度为 3px，颜色为半透明黑色（rgba(0,0,0,0.12)），hover 时加深至 rgba(0,0,0,0.2)

---

### Requirement: 对话框底部操作栏的固定定位与视觉设计

Normative: The entire action bar (q-card-actions container) SHALL be fixed to the bottom of the dialog using sticky positioning with a frosted glass background effect (backdrop-filter blur + semi-transparent overlay), allowing users to clearly see the signature list items blurred behind the action bar; the action bar SHALL have a subtle top border for visual separation.

#### Scenario: 按钮所在整个区域具有毛玻璃背景

- **GIVEN** 用户滚动签名列表内容
- **WHEN** 签名列表滚动到底部，与操作栏重叠
- **THEN** 操作栏整个区域显示毛玻璃背景（backdrop-filter: blur(10px); background: rgba(255,255,255,0.1)），用户能隐约看到后方列表项的模糊内容，同时具有上边界线（border-top: 1px solid rgba(0,0,0,0.08)）用于视觉分离

---

### Requirement: 主对话框滚动条隐藏优化

Normative: The scrollbar on the dialog's primary card element SHALL be hidden from view using `[&::-webkit-scrollbar]:hidden` to improve visual aesthetics and reduce visual clutter, while maintaining full scrolling functionality through mouse wheel and trackpad.

#### Scenario: 对话框主滚动条隐藏

- **GIVEN** 用户在签名列表中滚动内容
- **WHEN** 内容超过对话框高度，需要垂直滚动
- **THEN** 对话框边缘的垂直滚动条被隐藏，但用户仍可正常滚动，不影响任何功能

---

### Requirement: 内容区域的布局和间距优化

Normative: The name and description fields within each signature item SHALL have optimized vertical spacing, support horizontal scrolling for long text, and maintain visual consistency with the signature management page.

#### Scenario: 合理的竖向间距

- **GIVEN** 签名卡片渲染
- **WHEN** 名称和介绍显示在右侧信息区
- **THEN** 名称和介绍之间的竖向间距为 4-8px（使用 q-mt-2xs），内容区与图片区的间距为 8px（margin: 0 8px）

#### Scenario: 名称支持单行展示

- **GIVEN** 签名名称较长
- **WHEN** 超出容器宽度
- **THEN** 名称单行显示，支持水平滚动（overflow-x: auto），不会换行或被截断

#### Scenario: 介绍支持最多 2 行展示

- **GIVEN** 签名介绍较长
- **WHEN** 超出容器高度（2 行）或宽度
- **THEN** 介绍最多显示 2 行（-webkit-line-clamp: 2），若超出宽度则支持水平滚动
```

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

---

### Requirement: 专辑签名信息对话框

Normative: The `SignatureAuthorsDialog` SHALL 完整展示专辑中所有签名的详细信息，包括原始作者、直接导出作者、历史贡献作者；对话框 MUST 显示"资格码指纹"而非原始资格码，以保护资格码不泄漏的前提下保证签名的可追溯性；原始作者区块 MUST 显示联系方式（邮箱、其他联系方式）、授权状态、授权标识UUID、最近导出者资格码指纹及已授权签名列表；所有关键信息 MUST 支持一键复制；签名图片 MUST 通过 `GetAlbumFile` API 从专辑目录读取；所有用户可见文本 SHALL 使用 i18n 翻译键（`exportFlow.signatureInfoDialog.*`）。

#### Scenario: 查看有签名专辑信息

- **GIVEN** 用户在专辑页面点击"查看签名信息"按钮
- **WHEN** 对话框加载专辑签名数据成功
- **THEN** 对话框以分区卡片形式展示：
  - 原始作者区块（琥珀色）：签名卡片、资格码指纹、联系方式、授权状态、授权UUID、最近导出者资格码指纹、已授权列表（可展开）
  - 直接导出作者区块（蓝色）：仅在与原始作者不同时显示
  - 历史贡献作者区块（绿色）：列表形式展示
  - 签名统计摘要

#### Scenario: 查看无签名专辑

- **GIVEN** 用户在专辑页面点击"查看签名信息"按钮
- **WHEN** 专辑不包含任何签名
- **THEN** 对话框显示空态提示："此专辑尚未包含任何签名"及引导说明

#### Scenario: 复制签名信息

- **GIVEN** 对话框展示签名详情
- **WHEN** 用户点击任意复制按钮（资格码指纹/邮箱/UUID等）
- **THEN** 对应信息被复制到剪贴板，并显示成功通知

#### Scenario: 加载失败处理

- **GIVEN** 对话框尝试加载签名信息
- **WHEN** API 请求失败
- **THEN** 对话框显示错误信息并提供"重试"按钮

#### Scenario: 横向滚动条样式一致性

- **GIVEN** 签名名称或介绍文本过长
- **WHEN** 需要横向滚动查看完整内容
- **THEN** 横向滚动条样式与签名管理页签名列表保持一致（细滚动条样式）

#### Scenario: 签名图片加载

- **GIVEN** 对话框加载专辑签名数据
- **WHEN** 签名包含名片图片路径
- **THEN** 通过 `GetAlbumFile` API 读取专辑目录中的图片文件并显示

---

### Requirement: 资格码指纹计算

Normative: 资格码指纹 SHALL 在SDK端计算并返回给前端，前端不应接触原始资格码的计算逻辑；计算方式为：将资格码去除第2位（索引1）和第11位（索引10）字符后，计算 SHA256 哈希。此设计用于保护原始资格码不泄漏，同时保证签名的可追溯性。

#### Scenario: 资格码指纹生成

- **GIVEN** SDK 需要向前端返回签名信息
- **WHEN** 构建 `SignatureAuthorInfo` 结构
- **THEN** SDK 计算资格码指纹并填充 `qualificationFingerprint` 字段，前端直接使用该值展示

