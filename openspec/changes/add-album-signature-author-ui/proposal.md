# Change: 在专辑选择器中展示签名作者信息

## Why

目前键音专辑已支持签名功能，但在主页面和键音专辑页面的选择器中，用户无法直观看到专辑的签名作者信息。需要为选择器添加签名作者信息展示功能，让用户在选择专辑时能快速识别专辑来源和作者身份。

## What Changes

### 1. 后端 API 扩展

- 在获取专辑列表时，同时获取每个专辑的签名信息
- 在现有的专辑遍历解析过程中适配签名信息获取（复用 Viper 实例，遍历后及时释放）
- 新增 `GetAudioPackageListWithSignature` API 或扩展现有 API 返回签名信息

### 2. 前端状态管理

- 扩展 `main-store.ts` 中的 `keyTonePkgOptions` 数据结构，增加签名信息字段
- 新增 `keyTonePkgSignatureInfo` Map 用于存储每个专辑的签名信息

### 3. 选择器 UI 组件

#### 3.1 选中状态展示（选择器上边框右侧位置）

- 位置：选择器上边框靠右侧位置（利用空白区域）
- 内容：直接导出作者的图片（无图片则显示签名图标）+ 名称
- 条件：仅当专辑有签名时显示

#### 3.2 列表项展示

- 位置：专辑名称下方
- 内容：直接导出作者的图片（无图片则显示签名图标）+ 名称
- 条件：仅当专辑有签名时显示
- 增加列表项的信息展示密度

#### 3.3 悬停详情卡片

- 触发：鼠标悬停在签名信息区域
- 交互：鼠标可移动到卡片上，卡片保持显示
- 内容：显示更详细的签名信息
- 操作：右下角显示"点击查看详细信息"可点击 label
- 点击行为：打开 `SignatureAuthorsDialog` 对话框

### 4. 涉及页面

- **Main_page.vue**：主页面（统一/分离模式两个选择器）
- **Keytone_album_page_new.vue**：键音专辑页面（一个选择器）

### 5. 实现策略

- 优先使用 Quasar `q-select` 的插槽机制自定义：
  - `option` 插槽：自定义列表项
  - `selected-item` 插槽：自定义选中后的显示（上边框位置）
- 若插槽方案遇到瓶颈，则重新实现专门的专辑选择器组件

## Impact

- Affected specs: `export-flow`
- Affected code:
  - `sdk/server/server.go` - API 扩展
  - `sdk/audioPackage/list/list.go` - 专辑列表获取逻辑
  - `frontend/src/stores/main-store.ts` - 状态管理
  - `frontend/src/boot/query/keytonePkg-query.ts` - API 调用
  - `frontend/src/pages/Main_page.vue` - 主页面选择器
  - `frontend/src/pages/Keytone_album_page_new.vue` - 键音专辑页面选择器
  - 新增 `frontend/src/components/album-selector/` - 专辑选择器相关组件
