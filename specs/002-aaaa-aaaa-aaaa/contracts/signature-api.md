# Contract: Signature Operations

## Export Signature (.ktsign)

**描述**: 导出指定签名为 .ktsign 文件格式

**端点**: `POST /sdk/signatures/{name}/export`

**路径参数**:
- `name`: 签名名称（必填）

**请求体**: `{}` (空对象)

**响应**: `200 OK`

```json
{
  "fileNameSuggested": "My Signature.ktsign",
  "blobBase64": "SGVsbG8gV29ybGQ..."
}
```

**字段说明**:
- `fileNameSuggested`: 建议的文件名
- `blobBase64`: Base64 编码的加密签名数据

**错误响应**:
- `400 Bad Request`: 签名名称缺失
  ```json
  {
    "error": "invalid_request",
    "message": "Signature name is required"
  }
  ```
- `404 Not Found`: 签名不存在
  ```json
  {
    "error": "not_found",
    "message": "Signature not found"
  }
  ```

**说明**:
- 导出为 .ktsign（二进制），其内部为加密后的原始 JSON（等价于全局配置单一签名项）
- 前端接收到 base64 数据后需解码为二进制并触发下载

---

## Import Signature (.ktsign)

**描述**: 从 .ktsign 文件导入签名

**端点**: `POST /sdk/signatures/import`

**请求体**:

```json
{
  "fileName": "My Signature.ktsign",
  "blobBase64": "SGVsbG8gV29ybGQ...",
  "overwrite": true
}
```

**字段说明**:
- `fileName`: 文件名（必填）
- `blobBase64`: Base64 编码的文件内容（必填）
- `overwrite`: 是否覆盖已存在的签名（可选，默认 false）

**响应**: `201 Created`

```json
{
  "name": "My Signature",
  "overwritten": false
}
```

**字段说明**:
- `name`: 导入的签名名称
- `overwritten`: 是否覆盖了已存在的签名

**错误响应**:
- `400 Bad Request`: 请求数据无效
  ```json
  {
    "error": "invalid_request",
    "message": "Invalid request body: ..."
  }
  ```
- `400 Bad Request`: 文件格式无效
  ```json
  {
    "error": "invalid_format",
    "message": "invalid_format: failed to decode base64"
  }
  ```
- `409 Conflict`: 签名已存在且未指定覆盖
  ```json
  {
    "error": "exists_without_overwrite",
    "message": "Signature already exists"
  }
  ```

**处理流程**:
1. 解密 blobBase64 → 原始 JSON（key/value 与可选 assets）
2. 写入全局配置（/store/set），并落盘资源（如名片图片）
3. 触发 SSE 刷新

**说明**:
- 签名名称(name)作为唯一标识
- 保护码(protectCode)由前端在创建时生成，导入时保持不变
- 如果签名已存在但 intro 或 cardImagePath 不同，需要用户确认是否覆盖

---

## 签名管理（通过 /store 接口）

签名的列表、创建、删除操作通过现有的 `/store/get` 和 `/store/set` 端点实现：

### 获取签名列表

**端点**: `GET /store/get?key=signature_manager`

**响应**: `200 OK`

```json
{
  "message": "ok",
  "key": "signature_manager",
  "value": {
    "key_pc_123": {
      "name": "Zhang San",
      "intro": "My personal signature",
      "cardImagePath": "uuid123.png",
      "createdAt": "2025-10-01T10:00:00Z"
    },
    "key_pc_456": {
      "name": "Li Si",
      "intro": "",
      "createdAt": "2025-10-02T15:00:00Z"
    }
  }
}
```

### 创建/更新签名

**端点**: `POST /store/set`

**请求体**:

```json
{
  "key": "signature_manager",
  "value": {
    "key_pc_123": {
      "name": "Zhang San",
      "intro": "My personal signature",
      "cardImagePath": "uuid123.png",
      "createdAt": "2025-10-01T10:00:00Z"
    }
  }
}
```

**响应**: `200 OK`

```json
{
  "message": "ok"
}
```

**说明**:
- `key_pc_*` 是加密后的保护码，作为签名的存储键
- 前端在创建签名时自动生成保护码，格式：`pc_` + timestamp + random
- 保护码不在 UI 中展示
