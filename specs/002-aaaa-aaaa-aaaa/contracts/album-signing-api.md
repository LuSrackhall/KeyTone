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

