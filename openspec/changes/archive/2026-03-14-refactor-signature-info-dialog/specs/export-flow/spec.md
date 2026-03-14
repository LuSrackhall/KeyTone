# 专辑签名信息对话框规格（增量）

## Purpose

本增量规格描述 `SignatureAuthorsDialog` 组件重构后的展示要求，补充 `openspec/specs/export-flow/spec.md` 中新增的 **专辑签名信息对话框** 要求。

---

### Requirement: 完整展示原始作者授权元数据

Normative: The `SignatureAuthorsDialog` SHALL 从 `allSignatures` 中读取原始作者签名条目的完整授权元数据（`authorization` 字段），并分区展示联系方式、授权状态、授权标识UUID、直接导出作者资格码及已授权签名列表。

#### Scenario: 展示联系方式

- **GIVEN** 对话框加载了有签名的专辑
- **WHEN** 原始作者签名条目包含 `authorization.contactEmail`
- **THEN** 在原始作者区块显示邮箱信息，带复制按钮

#### Scenario: 展示其他联系方式

- **GIVEN** 原始作者签名条目包含 `authorization.contactAdditional`
- **WHEN** 对话框渲染
- **THEN** 显示其他联系方式（支持多行），带复制按钮

#### Scenario: 展示授权标识UUID

- **GIVEN** 原始作者签名条目包含 `authorization.authorizationUUID`
- **WHEN** 对话框渲染
- **THEN** 以等宽字体显示 UUID，带复制按钮

#### Scenario: 展示最近导出者资格码指纹

- **GIVEN** 原始作者签名条目包含 `authorization.directExportAuthor`
- **WHEN** 对话框渲染
- **THEN** 显示最近导出者的资格码指纹（由SDK计算的 `directExportAuthorFingerprint`）

#### Scenario: 展示已授权签名列表

- **GIVEN** 原始作者签名条目的 `requireAuthorization` 为 true
- **AND** `authorization.authorizedFingerprintList` 非空（不含原始作者自己）
- **WHEN** 用户展开"已授权签名列表"
- **THEN** 以列表形式显示所有已授权签名的资格码指纹（非原始资格码）

#### Scenario: 已授权信息仅在需要授权时显示

- **GIVEN** 原始作者签名条目的 `requireAuthorization` 为 false
- **WHEN** 对话框渲染
- **THEN** 不显示"已授权 N 个签名"标签和"已授权签名列表"展开项

#### Scenario: 已授权列表过滤原始作者

- **GIVEN** `authorization.authorizedList` 中包含原始作者自己的资格码
- **WHEN** SDK 计算 `authorizedFingerprintList`
- **THEN** 自动过滤掉原始作者自己，仅返回其他被授权者的指纹

---

### Requirement: 直接导出作者区块条件显示

Normative: The dialog SHALL 仅在 `directExportAuthor` 与 `originalAuthor` 的资格码不同时显示独立的"直接导出作者"区块，避免重复展示。

#### Scenario: 原始作者与直接导出作者相同

- **GIVEN** 专辑的原始作者即为直接导出作者
- **WHEN** 对话框渲染
- **THEN** 不显示独立的"直接导出作者"区块

#### Scenario: 原始作者与直接导出作者不同

- **GIVEN** 专辑的直接导出作者是被授权的贡献者
- **WHEN** 对话框渲染
- **THEN** 显示独立的"直接导出作者"区块（蓝色背景）

---

### Requirement: 复制功能

Normative: All key identifiers displayed in the dialog SHALL support one-click copy to clipboard, including qualification code fingerprints, email addresses, authorization UUIDs, and contact information.

#### Scenario: 复制资格码指纹

- **GIVEN** 对话框展示原始作者或贡献者的资格码指纹
- **WHEN** 用户点击复制按钮
- **THEN** 资格码指纹被复制到剪贴板，显示成功通知

---

### Requirement: 签名统计摘要

Normative: The dialog SHALL 在底部显示签名统计摘要，包括总签名数、原始作者数、贡献者数。

#### Scenario: 展示统计信息

- **GIVEN** 对话框加载了有签名的专辑
- **WHEN** 对话框渲染完成
- **THEN** 底部显示 q-chip 形式的统计信息

---

### Requirement: 签名图片加载

Normative: 签名图片 SHALL 通过 `GetAlbumFile` API 从专辑目录读取，使用 Blob URL 缓存机制，对话框关闭时 MUST 释放 Blob URL 资源。

#### Scenario: 加载签名图片

- **GIVEN** 对话框打开并加载签名数据
- **WHEN** 签名包含 `cardImagePath` 字段
- **THEN** 调用 `GetAlbumFile` API 读取图片，创建 Blob URL 并显示

#### Scenario: 图片加载失败

- **GIVEN** 签名图片路径无效或文件不存在
- **WHEN** API 请求失败
- **THEN** 显示默认占位图标（person 图标）

#### Scenario: 资源释放

- **GIVEN** 对话框已加载签名图片
- **WHEN** 用户关闭对话框
- **THEN** 所有 Blob URL 被释放

---

### Requirement: 资格码指纹

Normative: 对话框 SHALL 显示"资格码指纹"而非原始资格码；资格码指纹 MUST 由SDK端计算并返回，前端不应接触指纹计算逻辑。

#### Scenario: 资格码指纹展示

- **GIVEN** SDK 返回签名信息
- **WHEN** 签名信息包含 `qualificationFingerprint` 字段
- **THEN** 前端直接使用该字段值展示，无需本地计算

#### Scenario: 最近导出者资格码指纹

- **GIVEN** SDK 返回签名信息
- **WHEN** `AuthorizationMetadata` 包含 `directExportAuthor` 字段
- **THEN** SDK 自动计算 `directExportAuthorFingerprint` 并返回
- **AND** 前端标签显示为"最近导出者资格码指纹"

#### Scenario: 资格码指纹计算（SDK端）

- **GIVEN** SDK 需要构建 `SignatureAuthorInfo`
- **WHEN** 设置签名的资格码
- **THEN** 调用 `GenerateQualificationFingerprint` 计算指纹：
  - TIPS: 去除第2位（索引1）和第11位（索引10）字符
  - 对处理后的字符串计算 SHA256 哈希
  - 结果为64字符十六进制字符串

---

### Requirement: i18n 国际化

Normative: 对话框所有用户可见文本 SHALL 使用 i18n 翻译键，支持中英文切换。

#### Scenario: 文本国际化

- **GIVEN** 用户切换语言设置
- **WHEN** 打开签名信息对话框
- **THEN** 所有标签、按钮、提示文本均显示对应语言

#### 翻译键路径

`exportFlow.signatureInfoDialog.*`
