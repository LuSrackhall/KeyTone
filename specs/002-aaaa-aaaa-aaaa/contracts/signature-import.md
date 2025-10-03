# Contract: POST /sdk/signatures/import

## 描述
从 .ktsign 文件导入签名到签名管理器。

## 端点
`POST /sdk/signatures/import`

## 请求体
```json
{
  "fileName": "My Signature.ktsign",
  "blobBase64": "SGVsbG8gV29ybGQ...",
  "overwrite": true
}
```

**字段说明:**
- `fileName` (string, required): 文件名
- `blobBase64` (string, required): Base64 编码的文件内容
- `overwrite` (boolean, optional): 是否覆盖已存在的签名，默认 false

## 成功响应

### 201 Created
```json
{
  "name": "My Signature",
  "overwritten": false
}
```

**字段说明:**
- `name` (string): 导入的签名名称
- `overwritten` (boolean): 是否覆盖了已存在的签名

## 错误响应

### 400 Bad Request - 请求数据无效
```json
{
  "error": "invalid_request",
  "message": "Invalid request body: ..."
}
```

### 400 Bad Request - 文件格式无效（Base64解码失败）
```json
{
  "error": "invalid_format",
  "message": "invalid_format: failed to decode base64"
}
```

### 400 Bad Request - 文件格式无效（JSON解析失败）
```json
{
  "error": "invalid_format",
  "message": "invalid_format: failed to parse JSON"
}
```

### 400 Bad Request - 缺少必填字段
```json
{
  "error": "invalid_format",
  "message": "Missing or invalid name field"
}
```

### 409 Conflict - 签名已存在且未指定覆盖
```json
{
  "error": "exists_without_overwrite",
  "message": "Signature already exists"
}
```

## 处理流程

1. 接收并验证请求体（fileName, blobBase64）
2. 解码 Base64 数据
3. 解密并解析 JSON（获取 key/value 和可选 assets）
4. 验证 JSON 结构和必填字段（name）
5. 从 `/store/get` 读取现有 `signature_manager`
6. 检查签名是否已存在（按 name 匹配）
7. 如果存在且未指定 overwrite，返回 409 错误
8. 如果存在且指定 overwrite，删除旧条目
9. 添加新签名到签名管理器
10. 通过 `/store/set` 保存更新后的签名管理器
11. 触发 SSE 刷新（自动通过 `/store/set` 实现）
12. 返回导入结果

## 说明

- 签名名称(name)作为唯一标识
- 保护码(protectCode)由前端在创建时生成，导入时保持不变
- 如果签名已存在但 intro 或 cardImagePath 不同，需要用户确认是否覆盖
- 导入成功后会自动触发 SSE 刷新，前端签名列表自动更新
- 资源文件（如名片图片）的处理在 Stage 1 简化实现

## 测试用例

### Happy Path - 新建签名
```bash
POST /sdk/signatures/import
Body: { fileName: "test.ktsign", blobBase64: "...", overwrite: false }
Response: 201 Created with name and overwritten: false
```

### Happy Path - 覆盖签名
```bash
POST /sdk/signatures/import
Body: { fileName: "test.ktsign", blobBase64: "...", overwrite: true }
Response: 201 Created with name and overwritten: true
```

### Error Path - 签名已存在且未指定覆盖
```bash
POST /sdk/signatures/import
Body: { fileName: "existing.ktsign", blobBase64: "...", overwrite: false }
Response: 409 Conflict
```

### Error Path - 无效的 Base64
```bash
POST /sdk/signatures/import
Body: { fileName: "test.ktsign", blobBase64: "invalid!@#$", overwrite: false }
Response: 400 Bad Request with invalid_format error
```
