# 签名管理功能规格说明（变更：签名图片异步删除改为重命名标记 + 启动按文件名清理）

## MODIFIED Requirements

### Requirement: 删除签名

Normative: The system SHALL 在删除签名前尝试将该签名的名片图片重命名为 `Delete___<original_filename>` 作为删除标记（即在原文件名前增加固定前缀 `Delete___`）；仅当重命名成功（或签名无图片）时才删除配置中的签名条目并返回成功；若重命名失败则 MUST 返回失败且 MUST NOT 删除配置条目。

#### Scenario: 删除带图片的签名（成功）

- **GIVEN** 签名条目存在且其 `CardImage` 指向 `ConfigPath/signature` 下的现有文件
- **WHEN** 前端调用 `POST /signature/delete`（负载 `{ id }`）
- **THEN** 后端将图片文件重命名为 `Delete___<original_filename>`
- **AND** 后端删除配置中的该签名条目并返回成功响应
- **AND** 前端提示“删除成功”，签名在 SSE 推送后从列表移除

#### Scenario: 删除带图片的签名（重命名失败）

- **GIVEN** 签名条目存在且其 `CardImage` 指向的文件无法被重命名或无法被定位（权限/占用/路径非法/跨平台路径不可解析等）
- **WHEN** 前端调用 `POST /signature/delete`
- **THEN** 后端返回失败响应
- **AND** 后端 MUST NOT 删除配置中的该签名条目
- **AND** 前端提示“删除失败”并保留原列表项

#### Scenario: 删除不带图片的签名

- **GIVEN** 签名条目存在且 `CardImage` 为空
- **WHEN** 前端调用 `POST /signature/delete`
- **THEN** 后端删除配置中的该签名条目并返回成功响应

---

### Requirement: 签名图片路径初始化

Normative: The backend SHALL 在 `ConfigPath/signature` 下创建并使用签名图片目录；启动后的清理任务 MUST 删除所有文件名以 `Delete___` 开头的图片文件，且 MUST NOT 通过解密签名配置来推断需要删除的图片。

Non-normative note:
- 经检查后端实现，签名图片文件名不是基于图片内容去重，而是由 `id|name|originalImageName|timestamp`（导入为 `id|name|importFileName|timestamp`）参与生成，因此正常流程不会出现多个签名共享同一个图片文件。

#### Scenario: 初始化图片目录

- **GIVEN** 签名模块首次保存或更新图片
- **WHEN** 后端写入文件
- **THEN** 系统确保 `ConfigPath/signature` 目录存在

#### Scenario: 启动清理仅删除 Delete___* 文件

- **GIVEN** `ConfigPath/signature` 目录中同时存在 `Delete___*` 文件与普通图片文件
- **WHEN** 启动后执行签名图片清理任务
- **THEN** 所有 `Delete___*` 文件被删除
- **AND** 普通图片文件不会因为“签名配置解密失败/缺失”而被删除

