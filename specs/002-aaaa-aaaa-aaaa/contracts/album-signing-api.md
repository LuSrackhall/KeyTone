# Contract: Album Signing

本文档定义了专辑签名功能的 API 契约，包括查询、签名和导出操作。

---

## List Album Signatures

**描述**: 获取专辑的所有签名记录

**端点**: `GET /sdk/albums/{albumId}/signatures`

**响应**: `200 OK`

```json
[
  {
    "signerId": "Zhang San",
    "signedAt": [
      "2025-10-01T10:30:00Z",
      "2025-10-02T14:20:00Z"
    ],
    "allowReexport": true
  },
  {
    "signerId": "Li Si",
    "signedAt": [
      "2025-10-02T15:00:00Z"
    ],
    "allowReexport": false
  }
]
```

**字段说明**:
- `signerId`: 签名者名称
- `signedAt`: 签名时间戳数组（支持多次签名）
- `allowReexport`: 是否允许二次导出

**错误响应**:
- `404 Not Found`: 专辑不存在

---

## Sign Album (First or Additional)

**描述**: 为专辑添加签名（首次或追加时间戳）

**端点**: `POST /sdk/albums/{albumId}/sign`

**请求体**:

```json
{
  "signatureId": "Zhang San",
  "protectCode": "a1b2c3d4e5f6g7h8",
  "allowReexport": true
}
```

**字段说明**:
- `signatureId`: 签名名称（必填）
- `protectCode`: 签名保护码（必填）
- `allowReexport`: 是否允许二次导出（默认 true）

**响应**: `201 Created`

```json
{
  "updated": true,
  "signedAt": "2025-10-02T16:45:00Z"
}
```

**错误响应**:
- `400 Bad Request`: 签名不存在

  ```json
  {
    "error": "signature_not_found",
    "message": "签名不存在"
  }
  ```
- `401 Unauthorized`: 保护码错误

  ```json
  {
    "error": "invalid_protect_code",
    "message": "保护码错误"
  }
  ```
- `404 Not Found`: 专辑不存在

---

## Export Album With Signatures

**描述**: 导出带签名的专辑文件

**端点**: `POST /sdk/albums/{albumId}/export`

**请求体**:

```json
{
  "signatureId": "Zhang San",
  "protectCode": "a1b2c3d4e5f6g7h8"
}
```

**字段说明**:
- `signatureId`: 指定签名（可选，不指定则不附加新签名）
- `protectCode`: 签名保护码（如指定 signatureId 则必填）

**响应**: `200 OK`

```json
{
  "fileNameSuggested": "My Album.ktalbum",
  "blobBase64": "UEsDBBQAAAAIAC..."
}
```

**字段说明**:
- `fileNameSuggested`: 建议的专辑文件名
- `blobBase64`: Base64 编码的专辑文件内容

**错误响应**:
- `401 Unauthorized`: 保护码错误
- `404 Not Found`: 专辑或签名不存在
- `403 Forbidden`: 专辑已签名但未指定签名
  
  ```json
  {
    "error": "signature_required",
    "message": "专辑已签名，导出时必须选择签名"
  }
  ```

---

## Export Sign Bridge（导出签名桥）

**描述**: 在导出流程中，将已选择的签名写入专辑配置的 `album_signatures`（合并/去重/排序 `signedAt`），并返回导出继续所需数据。

**端点**: `POST /export/sign-bridge`

**请求体**:

```json
{
  "albumId": "<album-id>",
  "signatureName": "Zhang San"
}
```

**响应**: `200 OK`

```json
{
  "ok": true,
  "signedAtAppended": "2025-10-02T16:45:00Z"
}
```

**错误响应**:
- `404 Not Found`: 专辑不存在或签名不存在
- `409 Conflict`: 重复签名但未成功合并时间戳（应尽量通过去重合并避免该错误）

**说明**:
- 仅在“导出流程”中使用；签名管理对话框不提供“签名专辑/导出专辑”的入口（符合 FR-017）。
- `album_signatures` 中每个签名维护 `signedAt: string[]`；本端点在写入时进行去重并按时间排序。

