# Contract: POST /export/sign-bridge

## 描述
导出流程中的签名桥接端点。将前端导出流程中选择的签名写入专辑配置的 `album_signatures`，并返回导出继续所需数据。

## 端点
`POST /export/sign-bridge`

## 请求体
```json
{
  "albumId": "album-uuid-123",
  "signatureName": "Zhang San"
}
```

**字段说明:**
- `albumId` (string, required): 专辑UUID或标识
- `signatureName` (string, required): 签名名称

## 成功响应

### 200 OK
```json
{
  "ok": true,
  "signedAtAppended": "2025-10-02T16:45:00Z"
}
```

**字段说明:**
- `ok` (boolean): 操作是否成功
- `signedAtAppended` (string): 追加的签名时间戳（ISO8601 格式）

## 错误响应

### 400 Bad Request - 请求数据无效
```json
{
  "error": "invalid_request",
  "message": "Invalid request body: ..."
}
```

### 404 Not Found - 签名不存在
```json
{
  "error": "not_found",
  "message": "Signature not found"
}
```

### 404 Not Found - 专辑不存在
```json
{
  "error": "not_found",
  "message": "Album not found"
}
```

### 500 Internal Server Error - 内部错误
```json
{
  "error": "internal_error",
  "message": "Invalid signature_manager format"
}
```

## 处理流程

1. 接收并验证请求体（albumId, signatureName）
2. 从 `/store/get` 读取全局 `signature_manager`
3. 在签名列表中查找匹配的签名（按 name 匹配）
4. 如果签名不存在，返回 404 错误
5. 获取当前时间戳（ISO8601 格式）
6. （TODO Stage 1）通过 `/keytone_pkg/get` 获取专辑的 `album_signatures`
7. （TODO Stage 1）合并/去重/排序签名时间戳到 `signedAt` 数组
8. （TODO Stage 1）通过 `/keytone_pkg/set` 保存更新后的专辑配置
9. 返回成功响应和时间戳

## 说明

### Stage 1 实现范围
- 验证签名存在性
- 生成并返回时间戳
- 基础响应结构

### 后续增强（T026）
- 实际写入专辑配置的 `album_signatures`
- 实现时间戳数组的合并、去重、排序
- 支持签名历史记录

### `album_signatures` 字段结构
```json
{
  "album_signatures": {
    "<encrypt(sha256(decrypt(protectCode) + name))>": "<encrypt(JSON_payload)>"
  }
}
```

**JSON_payload (明文结构，加密前):**
```json
{
  "name": "Zhang San",
  "intro": "My signature",
  "cardImagePath": "uuid.png",
  "signedAt": ["2025-10-01T10:00:00Z", "2025-10-02T16:45:00Z"],
  "authorization": {
    "authCode": "sha256...",
    "authorizedList": ["code1", "code2"]
  }
}
```

**字段说明:**
- `signedAt` (string[]): 每次导出的时间戳数组，自动去重并按时间排序
- `authorization` (object, optional): 仅在原始作者签名中包含，Stage 2 实现

### 合并规则（T026）
1. 从专辑配置读取现有 `album_signatures`
2. 解密对应签名的 key 和 value
3. 获取现有 `signedAt` 数组
4. 追加新时间戳
5. 去重（相同时间戳只保留一个）
6. 按时间排序（从早到晚）
7. 加密并保存回专辑配置

## 使用场景

### 场景1: 首次签名
- 专辑无任何签名记录
- 导出时用户选择签名
- 创建新的签名记录，`signedAt` 包含一个时间戳

### 场景2: 追加签名
- 专辑已有该签名的记录
- 导出时用户再次选择同一签名
- 追加新时间戳到 `signedAt` 数组

### 场景3: 多人签名
- 专辑已有其他人的签名
- 当前用户选择自己的签名
- 创建新的签名记录（不同的 key）

## 测试用例

### Happy Path - 首次签名
```bash
POST /export/sign-bridge
Body: { albumId: "album-123", signatureName: "Zhang San" }
Response: 200 OK with timestamp
Verify: album_signatures contains new entry with signedAt: [timestamp]
```

### Happy Path - 追加签名
```bash
POST /export/sign-bridge
Body: { albumId: "album-123", signatureName: "Zhang San" }
Response: 200 OK with new timestamp
Verify: album_signatures entry has signedAt: [old_timestamp, new_timestamp]
```

### Error Path - 签名不存在
```bash
POST /export/sign-bridge
Body: { albumId: "album-123", signatureName: "NonExistent" }
Response: 404 Not Found
```

### Edge Case - 重复时间戳
- 同一签名在极短时间内导出两次
- 应该去重，`signedAt` 数组不包含重复时间戳

### Edge Case - 乱序时间戳
- 手动修改配置导致时间戳乱序
- 合并后应该自动排序
