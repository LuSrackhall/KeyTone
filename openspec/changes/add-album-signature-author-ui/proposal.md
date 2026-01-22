# Change: 在专辑选择器中展示签名作者信息

## Why

目前键音专辑已支持签名功能，但在主页面和键音专辑页面的选择器中，用户无法直观看到专辑的签名作者信息。需要为选择器添加签名作者信息展示功能，让用户在选择专辑时能快速识别专辑来源和作者身份。

## What Changes

### 1. 后端 API 扩展

- 在获取专辑列表时，同时获取每个专辑的签名信息（包含原始作者和直接导出作者）
- 在现有的专辑遍历解析过程中适配签名信息获取（独立创建 Viper 实例，函数结束后释放）
- 扩展 `/get_audio_package_list` API 返回签名信息
- 签名摘要包含：原始作者（名称、图片、介绍）、直接导出作者（名称、图片、介绍）、是否为同一作者标记

### 2. 前端状态管理

- 扩展 `main-store.ts` 中的 `keyTonePkgOptions` 数据结构，增加签名信息字段
- 新增 `keyTonePkgSignatureInfo` Map 用于存储每个专辑的签名摘要
- 新增 `getSignatureInfoByPath` 辅助方法

### 3. 选择器 UI 组件

#### 3.1 选中状态展示（选择器上边框右侧位置 - Legend 效果）

- 位置：选择器上边框靠右侧位置，**直接位于边框线上**（类似 HTML `<fieldset>` 的 `<legend>` 效果）
- 效果：使用背景色遮挡边框线，创造"打断"效果
- 内容：直接导出作者的图片（无图片则显示签名图标）+ 名称
- 条件：仅当专辑有签名时显示

#### 3.2 列表项展示（芯片样式）

- 位置：专辑名称下方
- 样式：**芯片（Chip）样式**，宽度根据内容自适应，不占整行
- 内容：直接导出作者的图片（无图片则显示签名图标）+ 名称
- 条件：仅当专辑有签名时显示
- 交互：支持悬停显示详情卡片

#### 3.3 悬停详情卡片

- 触发：鼠标悬停在签名信息区域
- 交互：鼠标可移动到卡片上，卡片保持显示；鼠标离开后延迟消失
- 内容：
  - **当原始作者=直接导出作者时**：只展示一个作者区块（避免重复）
  - **当原始作者≠直接导出作者时**：分两个区块展示（原始作者 + 直接导出作者）
  - 每个作者展示：图片、名称、介绍
- 视觉效果：**毛玻璃效果**（backdrop-filter: blur(8px)），增加高级感
- 操作：右下角显示"点击查看详细信息"可点击 label
- 点击行为：打开 `SignatureAuthorsDialog` 对话框

### 4. 涉及页面

- **Main_page.vue**：主页面（统一/分离模式两个选择器）
- **Keytone_album_page_new.vue**：键音专辑页面（一个选择器）- 待实现

### 5. 实现策略

- 使用 Quasar `q-select` 的插槽机制自定义：
  - `option` 插槽：自定义列表项，添加签名芯片和悬停卡片
- 使用绝对定位实现 legend 效果：
  - 签名徽章包装器定位在边框线上（top: -9px）
  - 使用与页面背景相同的渐变色遮挡边框

## Impact

- Affected specs: `export-flow`
- Affected code:
  - `sdk/audioPackage/list/audioPackageList.go` - 签名摘要获取逻辑
  - `sdk/server/server.go` - API 扩展
  - `frontend/src/types/album-selector.ts` - 类型定义（新增）
  - `frontend/src/stores/main-store.ts` - 状态管理
  - `frontend/src/pages/Main_page.vue` - 主页面选择器
  - `frontend/src/pages/Keytone_album_page_new.vue` - 键音专辑页面选择器（待实现）
  - `frontend/src/components/album-selector/` - 专辑选择器相关组件（新增）
    - `AlbumSignatureBadge.vue` - 签名徽章组件（芯片样式）
    - `AlbumSignatureHoverCard.vue` - 悬停详情卡片组件
    - `index.ts` - 组件导出入口
  - `frontend/src/i18n/zh-CN/index.json` - 中文翻译
  - `frontend/src/i18n/en-US/index.json` - 英文翻译
