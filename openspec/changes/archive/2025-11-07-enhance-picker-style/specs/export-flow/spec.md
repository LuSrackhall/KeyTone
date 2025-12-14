# 签名选择对话框样式优化

---

## ADDED Requirements

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
