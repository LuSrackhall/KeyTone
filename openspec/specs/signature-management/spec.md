# 签名管理功能规格说明

## Purpose

本文档定义桌面端签名管理模块在前端与后端之间的交互契约，确保创建、维护、导入导出以及排序流程一致且可追溯。

## Requirements

### Requirement: 签名列表加载与增量同步

Normative: The frontend SHALL retrieve the encrypted signature map via `GET /signature/list`, decrypt entries via `POST /signature/decrypt`（携带 `encryptedId`），并按 `sort.time` 升序维护顺序；SSE 回调 MUST 触发增量刷新而不是清空列表。

#### Scenario: 首次加载成功
- **GIVEN** 后端配置中存在一个或多个签名条目
- **WHEN** 前端调用 `GET /signature/list` 并针对每个条目调用 `POST /signature/decrypt`
- **THEN** 前端构建签名列表，按 `sort.time` 排序，并对有图片的条目调用 `POST /signature/get-image` 生成预览

#### Scenario: SSE 推送增量更新
- **GIVEN** 前端已注册 `useSignatureStore.sseSync` 作为配置变更回调
- **WHEN** 后端配置更新触发 SSE 消息
- **THEN** 前端重新获取并解密签名数据，合并新增或更新的条目，同时保持既有项的顺序

---

### Requirement: 创建签名

Normative: The system SHALL allow creating signatures with必填名称（1-50 字符）、可选介绍（≤500 字符）和可选名片图片（PNG/JPG/JPEG/GIF，≤5MB）；客户端 MUST 以 `multipart/form-data` 调用 `POST /signature/create`，后端 MUST 保存图片至 `ConfigPath/signature` 并写入加密配置。

#### Scenario: 成功创建签名
- **GIVEN** 用户填写合法名称、介绍和合规图片
- **WHEN** 前端提交包含 `id`、`name`、`intro` 及可选 `cardImage` 的请求到 `POST /signature/create`
- **THEN** 后端生成加密 ID、保存图片并更新配置，前端提示创建成功并在 SSE 推送后看到新签名

#### Scenario: 本地校验失败
- **GIVEN** 用户选择超过 5MB 或不支持格式的图片
- **WHEN** 前端执行表单校验
- **THEN** 阻止请求并提示错误，后端不会收到调用

---

### Requirement: 编辑签名

Normative: The system SHALL 打开编辑对话框时按 `encryptedId` 重新加载数据，保持名称字段只读，并通过 `POST /signature/update` 提交变更；后端 MUST 保留原始 `sort.time` 并根据 `removeImage`/`imageChanged` 标识处理图片。

#### Scenario: 对话框数据保持同步
- **GIVEN** 用户从签名列表打开编辑对话框
- **WHEN** 再次打开同一签名
- **THEN** 对话框使用当前签名数据初始化，名称字段只读，图片预览可用

#### Scenario: 更新签名成功
- **GIVEN** 用户修改介绍或更换图片
- **WHEN** 前端提交 `POST /signature/update`（必要时附带 `removeImage` 或 `imageChanged`）
- **THEN** 后端更新加密内容并保留排序时间戳，SSE 推送后列表展示最新数据

---

### Requirement: 删除签名

Normative: The system SHALL 在删除前显示确认对话框；客户端 MUST 调用 `POST /signature/delete`（负载 `{ id }`）并根据后端响应展示成功或失败提示。

#### Scenario: 删除成功
- **GIVEN** 用户在确认对话框中选择“删除”
- **WHEN** 前端调用 `POST /signature/delete` 并接收成功响应
- **THEN** 前端提示“删除成功”，签名在 SSE 推送后从列表移除

#### Scenario: 删除失败
- **GIVEN** 后端返回 4xx/5xx 错误
- **WHEN** 前端接收失败响应
- **THEN** 前端提示“删除失败”并保留原列表项

---

### Requirement: 导出签名

Normative: The client SHALL 通过 `POST /signature/export` 请求导出；后端 MUST 返回经 KeyB 加密并以十六进制编码的 `.ktsign` 内容，客户端 MUST 在用户确认保存后才提示成功。

#### Scenario: 导出成功
- **GIVEN** 用户在上下文菜单选择导出并确认风险提示
- **WHEN** 前端调用 `POST /signature/export` 并使用 `showSaveFilePicker`（或回退下载）保存 `{name}.ktsign`
- **THEN** 文件保存成功且内容可被解密为包含 `key`、`name`、`intro`、`cardImage`(十六进制) 的 JSON

#### Scenario: 用户取消保存
- **GIVEN** 保存对话框已打开
- **WHEN** 用户取消或关闭对话框
- **THEN** 前端不提示成功并终止导出流程

---

### Requirement: 导入签名

Normative: The client SHALL 上传 `.ktsign` 文件到 `POST /signature/import`；后端 MUST 使用 KeyB 解密并校验字段，当签名已存在时返回 `409` 且 `conflict: true`；覆盖流程 SHALL 通过 `POST /signature/import-confirm` 携带原始加密字符串和 `overwrite` 标识完成导入。

#### Scenario: 导入成功
- **GIVEN** 用户选择合法 `.ktsign` 文件且不存在同名签名
- **WHEN** 前端调用 `POST /signature/import`
- **THEN** 后端写入签名并返回成功，前端提示“导入成功”并等待 SSE 更新

#### Scenario: 导入冲突并确认覆盖
- **GIVEN** 目标签名已存在
- **WHEN** 初次导入返回 `409 conflict`
- **THEN** 前端提示用户选择；若用户确认覆盖，则调用 `POST /signature/import-confirm` 并在成功后展示提示

---

### Requirement: 签名图片路径初始化

Normative: The backend SHALL 在 `ConfigPath/signature` 下创建并使用签名图片目录，且 `CleanupOrphanCardImages` MUST 删除配置中未引用的文件。

#### Scenario: 初始化图片目录
- **GIVEN** 签名模块首次保存或更新图片
- **WHEN** 后端写入文件
- **THEN** 系统确保 `ConfigPath/signature` 目录存在，且项目根目录不会产生多余的 `signatures/`

#### Scenario: 清理孤立图片
- **GIVEN** 某些签名已被删除但图片仍存在磁盘
- **WHEN** 执行 `CleanupOrphanCardImages`
- **THEN** 所有未被引用的图片文件被移除

---

### Requirement: 签名排序持久化

Normative: The frontend SHALL 支持拖拽排序并通过 `POST /signature/update-sort` 提交新顺序；后端 MUST 更新每个条目的 `sort.time`。

#### Scenario: 排序成功
- **GIVEN** 用户拖动签名改变顺序
- **WHEN** 前端计算新的时间戳并调用 `POST /signature/update-sort`
- **THEN** 后端持久化排序并返回成功，SSE 推送后列表保持新顺序

#### Scenario: 排序失败回滚
- **GIVEN** 后端返回错误
- **WHEN** 排序更新请求失败
- **THEN** 前端提示失败并重新加载签名列表恢复原顺序
