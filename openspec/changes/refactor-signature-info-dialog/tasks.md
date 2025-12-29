# 专辑签名信息对话框重构 - 任务清单

## 概述

重构 `SignatureAuthorsDialog` 组件，完善信息展示并优化 UI 结构。

## 任务列表

### Phase 1: UI 重构

- [x] 1.1 重构对话框头部
  - [x] 添加 badge 图标
  - [x] 使用 `bg-deep-purple-1` 背景色

- [x] 1.2 重构原始作者区块
  - [x] 使用卡片式布局，`bg-amber-1` 头部
  - [x] 显示签名卡片（图片 + 名称 + 介绍）
  - [x] 显示资格码指纹（带复制按钮）
  - [x] 显示联系方式（邮箱 + 其他联系方式）
  - [x] 显示授权状态（需要授权/无需授权徽章）
  - [x] 显示授权标识UUID
  - [x] 显示直接导出作者资格码
  - [x] 显示已授权签名列表（展开/折叠）
  - [x] 显示已授权数量统计

- [x] 1.3 重构直接导出作者区块
  - [x] 使用卡片式布局，`bg-blue-1` 头部
  - [x] 仅在与原始作者不同时显示
  - [x] 显示签名卡片 + 资格码指纹

- [x] 1.4 重构历史贡献作者区块
  - [x] 使用卡片式布局，`bg-green-1` 头部
  - [x] 使用 q-list 列表形式
  - [x] 每项带复制资格码指纹按钮

- [x] 1.5 添加签名统计摘要
  - [x] 使用 q-chip 展示统计信息
  - [x] 显示总签名数、原始作者数、贡献者数

### Phase 2: 数据逻辑

- [x] 2.1 添加 originalAuthorEntry 计算属性
  - [x] 从 allSignatures 读取完整签名条目
  - [x] 获取授权元数据（contactEmail, contactAdditional, authorizationUUID 等）

- [x] 2.2 添加 isDirectExportAuthorSameAsOriginal 计算属性
  - [x] 比较资格码判断是否相同
  - [x] 相同时隐藏直接导出作者区块

### Phase 3: 交互功能

- [x] 3.1 实现复制功能
  - [x] copyToClipboard 方法
  - [x] 复制成功/失败通知

- [x] 3.2 实现横向滚动
  - [x] 复用签名列表滚动条样式
  - [x] 应用于名称和介绍字段

- [x] 3.3 优化错误处理
  - [x] 添加重试按钮
  - [x] 改进错误提示 UI

- [x] 3.4 优化空态提示
  - [x] 改进无签名时的提示文案
  - [x] 添加引导说明

### Phase 4: 图片处理

- [x] 4.1 修复图片加载问题 ✅
  - [x] 新增 SDK 端点 `get_album_file` 读取专辑内文件
  - [x] 新增前端 API `GetAlbumFile`
  - [x] 使用 Blob URL 缓存机制
  - [x] 对话框关闭时自动释放资源
  - [x] 添加 error slot 处理加载失败
  - [x] 统一占位图标样式

### Phase 5: 样式优化

- [x] 5.1 统一区块样式
  - [x] 圆角 8px
  - [x] 区块间距 12px

### Phase 6: 资格码指纹实现 ✅

- [x] 6.1 实现资格码指纹计算（SDK端）
  - [x] 在 `sdk/signature/album.go` 添加 `GenerateQualificationFingerprint` 函数
  - [x] TIPS: 去除第2位（索引1）和第11位（索引10）字符后计算SHA256
  - [x] 在 `SignatureAuthorInfo` 结构体添加 `qualificationFingerprint` 字段
  - [x] 在 `AuthorizationMetadata` 结构体添加 `directExportAuthorFingerprint` 字段
  - [x] 构建签名信息时自动计算指纹

- [x] 6.2 更新前端
  - [x] 更新类型定义，添加 `qualificationFingerprint` 和 `directExportAuthorFingerprint` 字段
  - [x] 移除前端指纹计算逻辑
  - [x] 直接使用SDK返回的指纹字段展示
  - [x] "直接导出作者资格码"改为"最近导出者资格码指纹"

### Phase 7: i18n 国际化 ✅

- [x] 7.1 添加中文翻译
  - [x] 在 `zh-CN/index.json` 添加 `exportFlow.signatureInfoDialog.*` 配置

- [x] 7.2 添加英文翻译
  - [x] 在 `en-US/index.json` 添加 `exportFlow.signatureInfoDialog.*` 配置

- [x] 7.3 更新组件
  - [x] 引入 `useI18n`
  - [x] 替换所有硬编码文本为 `t()` 调用

## 已完成

所有任务已完成，组件重构成功：

1. ✅ 完整展示原始作者的所有信息（联系方式、授权UUID、已授权列表等）
2. ✅ 分区卡片式布局，颜色区分清晰
3. ✅ 所有关键信息支持一键复制
4. ✅ 横向滚动条样式与签名列表保持一致
5. ✅ 已授权列表可展开/折叠
6. ✅ 签名统计摘要
7. ✅ 错误状态带重试按钮
8. ✅ 资格码指纹（SDK端计算，保护原始资格码）
9. ✅ 最近导出者资格码指纹展示（修复原资格码泄漏问题）
10. ✅ 完整的 i18n 国际化支持

## 测试要点

- [ ] 验证有签名专辑的信息展示完整性
- [ ] 验证无签名专辑的空态提示
- [ ] 验证复制功能正常
- [ ] 验证图片加载失败时的占位显示
- [ ] 验证已授权列表展开/折叠功能
- [ ] 验证长文本横向滚动
- [ ] 验证资格码指纹展示正确
- [ ] 验证中英文切换正常
