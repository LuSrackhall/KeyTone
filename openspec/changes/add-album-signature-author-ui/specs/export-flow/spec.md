# 专辑导出签名流程规格说明 - Delta

## ADDED Requirements

### Requirement: 专辑选择器签名信息展示

Normative: The album selector component SHALL display direct export author signature information for albums that have signatures; the selected state display MUST be positioned on the top-right border line of the selector (legend effect) and MUST mask the border line using a background that matches the selector control to avoid a visible dark block; the list item display MUST be positioned below the album name and rendered as a chip that sizes to content (not full-width); signature information SHALL only be shown when the album has a signature; the signature display MUST include the author's avatar image (or a badge icon as placeholder when image is not loaded) and the author name; the list item signature chip SHALL support the same hover card as the selected state.

#### Scenario: 选中专辑后展示签名信息

- **GIVEN** 用户已选择一个带有签名的专辑
- **WHEN** 选择器显示选中状态
- **THEN** 在选择器上边框靠右位置展示直接导出作者的头像（若无图片则显示 badge 图标）和名称，背景遮挡边框线且不出现明显黑色矩形

#### Scenario: 选中无签名专辑

- **GIVEN** 用户已选择一个没有签名的专辑
- **WHEN** 选择器显示选中状态
- **THEN** 选择器上边框位置不展示任何签名信息

#### Scenario: 列表项中展示签名信息

- **GIVEN** 专辑选择器下拉列表展开
- **WHEN** 列表中存在带有签名的专辑
- **THEN** 该专辑的列表项在专辑名称下方以芯片样式展示直接导出作者的头像（若无图片则显示 badge 图标）和名称，芯片宽度随内容自适应且支持悬停详情卡片

#### Scenario: 列表项中无签名专辑

- **GIVEN** 专辑选择器下拉列表展开
- **WHEN** 列表中存在没有签名的专辑
- **THEN** 该专辑的列表项仅展示专辑名称，不显示签名相关区域

---

### Requirement: 专辑选择器签名悬停详情卡片

Normative: The album selector SHALL provide a hover card for signature information; the hover card MUST remain visible when the user's mouse moves onto the card; the hover card MUST display detailed signature information for both the original author and the direct export author (image, name, intro); when the original author and direct export author are the same, the hover card MUST show a single combined section; the hover card MUST include a clickable "点击查看详细信息" label at the bottom-right corner; clicking the label SHALL open the `SignatureAuthorsDialog` dialog to display complete album and signature information.

#### Scenario: 鼠标悬停显示详情卡片

- **GIVEN** 用户将鼠标悬停在签名信息区域
- **WHEN** 悬停持续一定时间（如 200ms）
- **THEN** 显示一个详情卡片，包含原始作者与直接导出作者的图片、名称、介绍

#### Scenario: 鼠标移入卡片保持显示

- **GIVEN** 详情卡片已显示
- **WHEN** 用户将鼠标从签名信息区域移动到详情卡片上
- **THEN** 详情卡片保持显示不消失

#### Scenario: 鼠标离开卡片后消失

- **GIVEN** 详情卡片已显示
- **WHEN** 用户将鼠标移出签名信息区域和详情卡片
- **THEN** 详情卡片在短暂延迟后消失（如 100-150ms）

#### Scenario: 点击查看详细信息

- **GIVEN** 详情卡片已显示
- **WHEN** 用户点击"点击查看详细信息" label
- **THEN** 打开 `SignatureAuthorsDialog` 对话框，展示完整的专辑及签名信息

#### Scenario: 列表项中点击查看详细信息

- **GIVEN** 专辑选择器列表已展开且签名悬停卡片可见
- **WHEN** 用户在列表项的悬停卡片中按下（pointerdown.capture）"点击查看详细信息"
- **THEN** 选择器弹层保持打开（不改变选择器状态），卡片根节点阻止 pointerdown/mousedown 冒泡，并在 pointerdown.capture 阶段打开 `SignatureAuthorsDialog` 对话框，避免弹层关闭导致 click 丢失

---

### Requirement: 专辑列表签名信息获取

Normative: The system SHALL fetch signature summary information for each album when loading the album list; the signature information MUST be obtained during the existing album traversal process (using an isolated Viper instance per album); the Viper instances MUST be released promptly after traversal to free memory; the signature summary SHALL include: `hasSignature`, original author fields (`originalAuthorName`, `originalAuthorImage`, `originalAuthorIntro`), direct export author fields (`directExportAuthorName`, `directExportAuthorImage`, `directExportAuthorIntro`), and `isSameAuthor`.

#### Scenario: 获取专辑列表时同步获取签名摘要

- **GIVEN** 前端调用获取专辑列表 API
- **WHEN** 后端遍历专辑目录
- **THEN** 对每个专辑同时读取签名配置，返回签名摘要信息（是否有签名、直接导出作者名称、图片路径）

#### Scenario: 签名摘要数据结构

- **GIVEN** 专辑包含签名
- **WHEN** API 返回签名摘要
- **THEN** 摘要包含 `hasSignature: true`、原始作者字段、直接导出作者字段，以及 `isSameAuthor` 标记

#### Scenario: 无签名专辑的摘要

- **GIVEN** 专辑不包含签名
- **WHEN** API 返回签名摘要
- **THEN** 摘要包含 `hasSignature: false`，其他字段为空字符串或 false

#### Scenario: 及时释放 Viper 实例

- **GIVEN** 后端遍历专辑列表并读取签名配置
- **WHEN** 单个专辑的配置读取完成
- **THEN** 立即释放该专辑对应的 Viper 实例以释放内存

---

### Requirement: 专辑选择器签名组件目录约束

Normative: All album selector signature display components SHALL be stored in `frontend/src/components/album-selector/`; components MUST include detailed comments explaining their purpose, parameters, and interaction logic; component naming SHALL follow PascalCase convention.

#### Scenario: 创建签名徽章组件

- **GIVEN** 开发者创建签名徽章展示组件
- **WHEN** 组件被提交
- **THEN** 组件位于 `components/album-selector/AlbumSignatureBadge.vue`，包含详细注释

#### Scenario: 创建悬停卡片组件

- **GIVEN** 开发者创建悬停详情卡片组件
- **WHEN** 组件被提交
- **THEN** 组件位于 `components/album-selector/AlbumSignatureHoverCard.vue`，包含详细的交互逻辑注释

---

### Requirement: 主页面与专辑页面选择器兼容

Normative: The signature display functionality SHALL be implemented in both Main_page.vue (unified and split mode selectors) and Keytone_album_page_new.vue selectors; each page MAY have different styling based on its specific context while maintaining consistent core functionality; the implementation SHALL prioritize using Quasar `q-select` slot mechanisms; if slot limitations are encountered, a custom album selector component SHALL be created as an alternative.

#### Scenario: 主页面统一模式选择器

- **GIVEN** 用户在主页面使用统一模式（单选择器）
- **WHEN** 选择器展开或选中专辑
- **THEN** 签名信息正确展示在列表项和选中状态位置

#### Scenario: 主页面分离模式选择器

- **GIVEN** 用户在主页面使用分离模式（键盘/鼠标双选择器）
- **WHEN** 任一选择器展开或选中专辑
- **THEN** 对应选择器的签名信息正确展示

#### Scenario: 键音专辑页面选择器

- **GIVEN** 用户在键音专辑页面
- **WHEN** 选择器展开或选中专辑
- **THEN** 签名信息正确展示，与主页面功能一致

#### Scenario: 插槽方案遇到瓶颈

- **GIVEN** 使用 Quasar `q-select` 插槽无法实现所需效果
- **WHEN** 开发者评估后确认需要自定义
- **THEN** 创建自定义的 `AlbumSelector.vue` 组件替代原有选择器
