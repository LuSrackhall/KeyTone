# Tasks: 在专辑选择器中展示签名作者信息

## 1. 后端 API 扩展

- [x] 1.1 在 `sdk/audioPackage/list/audioPackageList.go` 中扩展专辑列表获取逻辑
  - [x] 1.1.1 在遍历专辑时读取签名配置信息
  - [x] 1.1.2 提取直接导出作者（directExportAuthor）的基本信息
  - [x] 1.1.3 返回签名信息摘要（名称、图片路径、是否有签名）
  - [x] 1.1.4 确保 Viper 实例遍历后及时释放内存

- [x] 1.2 在 `sdk/server/server.go` 中添加/扩展 API
  - [x] 1.2.1 修改 `GetAudioPackageList` 返回结构，增加签名信息字段
  - [ ] 1.2.2 添加 `GetAlbumSignatureSummary` 接口用于获取单个专辑的签名摘要

## 2. 前端类型定义

- [x] 2.1 在 `frontend/src/types/` 中定义签名摘要类型
  - [x] 2.1.1 创建 `AlbumSignatureSummary` 接口
  - [x] 2.1.2 定义字段：`hasSignature`, `directExportAuthorName`, `directExportAuthorImage`, `albumPath`

## 3. 前端状态管理

- [x] 3.1 扩展 `frontend/src/stores/main-store.ts`
  - [x] 3.1.1 新增 `keyTonePkgSignatureInfo: Map<string, AlbumSignatureSummary>` 状态
  - [x] 3.1.2 修改 `GetKeyToneAlbumList` 函数，同时获取签名信息
  - [x] 3.1.3 提供 `getSignatureInfoByPath(path: string)` 辅助方法

## 4. 前端 API 调用

- [x] 4.1 在 `frontend/src/boot/query/keytonePkg-query.ts` 中
  - [x] 4.1.1 修改或新增获取专辑列表的 API 调用以支持签名信息
  - [ ] 4.1.2 添加获取专辑签名摘要的单独 API 调用（如需要）

## 5. 创建签名展示组件

- [x] 5.1 创建 `frontend/src/components/album-selector/AlbumSignatureBadge.vue`
  - [x] 5.1.1 实现签名作者头像展示（图片或签名图标）
  - [x] 5.1.2 实现签名作者名称展示
  - [x] 5.1.3 支持不同尺寸（列表项/选中状态）
  - [x] 5.1.4 添加详细注释说明组件用途和参数

- [x] 5.2 创建 `frontend/src/components/album-selector/AlbumSignatureHoverCard.vue`
  - [x] 5.2.1 实现悬停详情卡片布局
  - [x] 5.2.2 实现鼠标移入卡片时保持显示的交互逻辑
  - [x] 5.2.3 实现"点击查看详细信息"可点击 label
  - [x] 5.2.4 实现点击后打开 SignatureAuthorsDialog 的功能
  - [x] 5.2.5 添加详细注释说明组件交互逻辑

## 6. 键音专辑页面选择器改造

- [ ] 6.1 修改 `frontend/src/pages/Keytone_album_page_new.vue`
  - [ ] 6.1.1 使用 `option` 插槽自定义列表项，在专辑名称下方展示签名信息
  - [ ] 6.1.2 在选择器上边框右侧位置添加签名信息展示区域
  - [ ] 6.1.3 集成 `AlbumSignatureBadge` 组件
  - [ ] 6.1.4 集成 `AlbumSignatureHoverCard` 组件
  - [ ] 6.1.5 添加详细注释说明改造目的和实现方式

## 7. 主页面选择器改造

- [x] 7.1 修改 `frontend/src/pages/Main_page.vue` 统一模式选择器
  - [x] 7.1.1 使用 `option` 插槽自定义列表项展示签名信息
  - [x] 7.1.2 在选择器上边框右侧位置添加签名信息展示
  - [x] 7.1.3 集成签名相关组件
  - [x] 7.1.4 添加详细注释

- [x] 7.2 修改 `frontend/src/pages/Main_page.vue` 分离模式选择器
  - [x] 7.2.1 键盘专辑选择器：同上改造
  - [x] 7.2.2 鼠标专辑选择器：同上改造
  - [x] 7.2.3 添加详细注释

## 8. 国际化支持

- [x] 8.1 在 `frontend/src/i18n/` 中添加相关翻译键
  - [x] 8.1.1 添加"点击查看详细信息"等文本翻译
  - [x] 8.1.2 添加悬停卡片中的文本翻译
  - [ ] 8.1.3 添加无签名状态的提示文本

## 9. 测试与验证

- [ ] 9.1 功能测试
  - [ ] 9.1.1 测试有签名专辑的展示效果
  - [ ] 9.1.2 测试无签名专辑的展示效果
  - [ ] 9.1.3 测试悬停卡片的交互
  - [ ] 9.1.4 测试点击查看详细信息的功能
  - [ ] 9.1.5 测试 SignatureAuthorsDialog 正确打开

## 10. 文档更新

- [x] 10.1 更新规格文档
  - [x] 10.1.1 更新 `export-flow` spec 的相关需求
