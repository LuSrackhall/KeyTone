# 专辑签名信息对话框重构提案

## 概述

重构键音专辑页面的 `SignatureAuthorsDialog` 组件，修复信息展示不完整的问题，并采用更清晰的分区结构显示签名的所有相关信息。

## 问题分析

当前 `SignatureAuthorsDialog` 存在以下问题：

1. **信息展示不完整**
   - 缺少联系方式（contactEmail, contactAdditional）
   - 缺少授权标识UUID（authorizationUUID）
   - 缺少直接导出作者资格码
   - 缺少已授权签名列表详情

2. **UI 结构不够清晰**
   - 缺乏信息分组和层级
   - 没有复制功能支持
   - 长文本溢出处理不当

3. **视觉设计待优化**
   - 缺少区块颜色区分
   - 错误状态无重试功能
   - 统计信息缺失

4. **图片无法正常展示** ✅ 已修复
   - 使用 `file://` 协议在 Electron 中有安全限制
   - 需要通过 SDK API 读取专辑目录中的图片文件

5. **资格码隐私保护** ✅ 已实现
   - 直接展示资格码存在泄漏风险
   - 改为展示"资格码指纹"以保护隐私

## 解决方案

### 1. 重构信息展示结构

采用分区卡片式布局，每个区块有独立的颜色标识：

- **原始作者区块**（琥珀色背景）
  - 签名卡片（图片 + 名称 + 介绍）
  - 资格码指纹（带复制按钮）
  - 联系方式（邮箱 + 其他联系方式）
  - 授权状态（需要授权/无需授权徽章）
  - 授权标识UUID
  - 直接导出作者资格码
  - 已授权签名列表（可展开/折叠）

- **直接导出作者区块**（蓝色背景）
  - 仅在与原始作者不同时显示
  - 签名卡片 + 资格码指纹

- **历史贡献作者区块**（绿色背景）
  - 列表形式展示
  - 每项带复制资格码指纹按钮

- **签名统计摘要**（灰色背景）
  - 总签名数
  - 原始作者/贡献者计数

### 2. 数据源优化

从 `allSignatures` 中读取完整的 `AlbumSignatureEntry`，包含：
- `authorization.contactEmail` - 联系邮箱
- `authorization.contactAdditional` - 其他联系方式
- `authorization.authorizationUUID` - 授权标识UUID
- `authorization.directExportAuthor` - 直接导出作者资格码
- `authorization.authorizedList` - 已授权签名列表

### 3. 交互增强

- 所有关键信息支持一键复制（资格码指纹、邮箱、UUID等）
- 已授权列表使用展开/折叠组件
- 错误状态提供重试按钮
- 横向滚动条样式与签名列表保持一致

### 4. 图片加载修复

- 新增 SDK 端点 `GET /keytone_pkg/get_album_file` 用于读取专辑内文件
- 前端通过 `GetAlbumFile` API 加载签名图片
- 使用 Blob URL 缓存机制，对话框关闭时自动释放

### 5. 资格码指纹实现

**背景说明**：
- 资格码是签名原始ID的SHA256哈希结果（64字符十六进制字符串）
- 直接展示资格码可能导致隐私泄漏

**SDK端计算（安全性考虑）**：
- 指纹计算在SDK端进行，前端不接触计算逻辑
- SDK 新增 `GenerateQualificationFingerprint` 函数
- `SignatureAuthorInfo` 结构体新增 `qualificationFingerprint` 字段
- `AuthorizationMetadata` 结构体新增 `directExportAuthorFingerprint` 字段（动态计算，不存储）
- 前端直接使用SDK返回的指纹值展示

**资格码指纹计算方式**：
- TIPS: 将资格码去除第2位（索引1）和第11位（索引10）字符
- 对处理后的字符串计算SHA256哈希
- 结果为64字符十六进制字符串

### 6. i18n 国际化支持

**新增 i18n 配置**：
- 中文：`exportFlow.signatureInfoDialog.*`
- 英文：`exportFlow.signatureInfoDialog.*`

**主要翻译键**：
- `title` - 对话框标题
- `originalAuthor` - 原始作者
- `qualificationFingerprint` - 资格码指纹
- `latestExporterFingerprint` - 最近导出者资格码指纹
- `contributorAuthors` - 历史贡献作者
- 等...

## 涉及文件

- `frontend/src/components/export-flow/SignatureAuthorsDialog.vue` - 完整重构
- `frontend/src/boot/query/keytonePkg-query.ts` - 新增 GetAlbumFile API
- `frontend/src/types/export-flow.ts` - 新增 qualificationFingerprint 和 directExportAuthorFingerprint 字段
- `frontend/src/i18n/zh-CN/index.json` - 新增 signatureInfoDialog i18n 配置
- `frontend/src/i18n/en-US/index.json` - 新增 signatureInfoDialog i18n 配置
- `sdk/server/server.go` - 新增 get_album_file 端点
- `sdk/signature/album.go` - 新增 GenerateQualificationFingerprint 函数
- `sdk/audioPackage/config/signatureHelper.go` - SignatureAuthorInfo 新增指纹字段
- `sdk/audioPackage/config/signatureConfig.go` - AuthorizationMetadata 新增 directExportAuthorFingerprint 字段

## UI 布局示意

```text
┌──────────────────────────────────────────────┐
│ [badge] 专辑签名信息                     [×] │
├──────────────────────────────────────────────┤
│ ┌────────────────────────────────────────┐   │
│ │ [★] 原始作者              [需授权导出] │   │
│ ├────────────────────────────────────────┤   │
│ │ [图片] 名称                            │   │
│ │        介绍                            │   │
│ │                                        │   │
│ │ [指纹] 资格码指纹: xxx...xxx    [复制] │   │
│ │ ───────────────────────────────────    │   │
│ │ 联系方式                               │   │
│ │ [邮箱] author@example.com       [复制] │   │
│ │ [聊天] 微信: xxx               [复制] │   │
│ │ ───────────────────────────────────    │   │
│ │ 授权状态                               │   │
│ │ [需要授权] [已授权 3 个签名]           │   │
│ │ [钥匙] 授权UUID: xxx...xxx     [复制] │   │
│ │ [下载] 直接导出作者: xxx...xxx         │   │
│ │ [列表] 已授权签名列表 (3)        [▼]   │   │
│ └────────────────────────────────────────┘   │
│                                              │
│ ┌────────────────────────────────────────┐   │
│ │ [⬇] 直接导出作者          [最近导出者] │   │
│ ├────────────────────────────────────────┤   │
│ │ [图片] 名称                            │   │
│ │        介绍                            │   │
│ │ [指纹] 资格码: xxx...xxx               │   │
│ └────────────────────────────────────────┘   │
│                                              │
│ ┌────────────────────────────────────────┐   │
│ │ [👥] 历史贡献作者 (2)                   │   │
│ ├────────────────────────────────────────┤   │
│ │ [头像] 名称 介绍                [复制] │   │
│ │ [头像] 名称 介绍                [复制] │   │
│ └────────────────────────────────────────┘   │
│                                              │
│ ┌────────────────────────────────────────┐   │
│ │ 签名统计                               │   │
│ │ [共 3 个签名] [1 原始作者] [2 贡献者]  │   │
│ └────────────────────────────────────────┘   │
├──────────────────────────────────────────────┤
│                              [关闭]          │
└──────────────────────────────────────────────┘
```

## 兼容性

- 无 SDK/API 变更
- 无数据结构变更
- 纯前端 UI 重构
