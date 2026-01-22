# Tasks: 在专辑选择器中展示签名作者信息

## 1. 后端 API 扩展

- [x] 1.1 在 `sdk/audioPackage/list/audioPackageList.go` 中扩展专辑列表获取逻辑
  - [x] 1.1.1 在遍历专辑时读取签名配置信息
  - [x] 1.1.2 提取原始作者和直接导出作者的完整信息（名称、图片、介绍）
  - [x] 1.1.3 返回签名信息摘要（包含 isSameAuthor 标记）
  - [x] 1.1.4 独立创建 Viper 实例，函数结束后释放以避免内存泄漏

- [x] 1.2 在 `sdk/server/server.go` 中添加/扩展 API
  - [x] 1.2.1 修改 `GetAudioPackageList` 返回结构，增加签名信息字段

## 2. 前端类型定义

- [x] 2.1 在 `frontend/src/types/album-selector.ts` 中定义签名摘要类型
  - [x] 2.1.1 创建 `AlbumSignatureSummary` 接口
  - [x] 2.1.2 定义字段：原始作者信息、直接导出作者信息、isSameAuthor 标记

## 3. 前端状态管理

- [x] 3.1 扩展 `frontend/src/stores/main-store.ts`
  - [x] 3.1.1 新增 `keyTonePkgSignatureInfo: Map<string, AlbumSignatureSummary>` 状态
  - [x] 3.1.2 修改 `GetKeyToneAlbumList` 函数，同时获取签名信息
  - [x] 3.1.3 提供 `getSignatureInfoByPath(path: string)` 辅助方法

## 4. 创建签名展示组件

- [x] 4.1 创建 `frontend/src/components/album-selector/AlbumSignatureBadge.vue`（芯片样式）
  - [x] 4.1.1 使用 `inline-flex` 确保宽度由内容决定
  - [x] 4.1.2 圆角胶囊形状（`border-radius: 999px`）
  - [x] 4.1.3 支持 `normal` 和 `small` 两种尺寸
  - [x] 4.1.4 添加详细注释说明组件用途和参数

- [x] 4.2 创建 `frontend/src/components/album-selector/AlbumSignatureHoverCard.vue`
  - [x] 4.2.1 实现毛玻璃效果（`backdrop-filter: blur(8px)`）
  - [x] 4.2.2 根据 `isSameAuthor` 决定展示一个或两个作者区块
  - [x] 4.2.3 实现鼠标移入卡片时保持显示、移出后延迟消失的交互逻辑
  - [x] 4.2.4 实现"点击查看详细信息"可点击 label
  - [x] 4.2.5 添加详细注释说明组件交互逻辑

## 5. 主页面选择器改造

- [x] 5.1 修改 `frontend/src/pages/Main_page.vue` 统一模式选择器
  - [x] 5.1.1 实现 Legend 效果：签名徽章位于选择器边框上
  - [x] 5.1.2 使用 `option` 插槽自定义列表项，添加签名芯片和悬停卡片
  - [x] 5.1.3 添加 Legend 效果的 CSS 样式
  - [x] 5.1.4 添加详细注释

- [x] 5.2 修改 `frontend/src/pages/Main_page.vue` 分离模式选择器
  - [x] 5.2.1 键盘专辑选择器：同上改造
  - [x] 5.2.2 鼠标专辑选择器：同上改造
  - [x] 5.2.3 添加详细注释

## 6. 键音专辑页面选择器改造

- [x] 6.1 修改 `frontend/src/pages/Keytone_album_page_new.vue`
  - [x] 6.1.1 实现 Legend 效果：签名徽章位于选择器边框上
  - [x] 6.1.2 使用 `option` 插槽自定义列表项，添加签名芯片和悬停卡片
  - [x] 6.1.3 集成 `SignatureAuthorsDialog` 组件
  - [x] 6.1.4 添加详细注释说明改造目的和实现方式

## 7. 国际化支持

- [x] 7.1 在 `frontend/src/i18n/zh-CN/index.json` 中添加翻译
  - [x] 7.1.1 添加"原始作者"翻译
  - [x] 7.1.2 添加"直接导出作者"翻译
  - [x] 7.1.3 添加"原始作者/直接导出作者"翻译（同一作者时使用）
  - [x] 7.1.4 添加"点击查看详细信息"翻译

- [x] 7.2 在 `frontend/src/i18n/en-US/index.json` 中添加英文翻译

## 8. 文档更新

- [x] 8.1 更新规格文档
  - [x] 8.1.1 更新 `proposal.md`
  - [x] 8.1.2 更新 `design.md`
  - [x] 8.1.3 更新 `tasks.md`

## 9. 测试与验证

- [ ] 9.1 功能测试
  - [ ] 9.1.1 测试有签名专辑的展示效果
  - [ ] 9.1.2 测试无签名专辑的展示效果
  - [ ] 9.1.3 测试悬停卡片的交互（进入、保持、离开消失）
  - [ ] 9.1.4 测试毛玻璃效果
  - [ ] 9.1.5 测试 Legend 效果（签名徽章在边框上）
  - [ ] 9.1.6 测试点击查看详细信息打开 SignatureAuthorsDialog
