# Contract: Signature Operations

## List Signatures

GET /sdk/signatures
Response: 200 [{ name, intro?, cardImagePath?, createdAt }]

## Create Signature

POST /sdk/signatures
Body: { name, intro?, cardImageBase64? }
Response: 201 { name }
Errors: 400 (duplicate name)

## Delete Signature

DELETE /sdk/signatures/{name}
Response: 204

## Export Signature (.ktsign)

POST /sdk/signatures/{name}/export
Body: { }
Response: 200 { fileNameSuggested: string, blobBase64: string }
Notes:
- 导出为 .ktsign（二进制），其内部为加密后的原始 JSON（等价于全局配置单一签名项）。

## Import Signature (.ktsign)

POST /sdk/signatures/import
Body: { fileName: string, blobBase64: string, overwrite?: boolean }
Response: 201 { name, overwritten: boolean }
Process:
1) 解密 blobBase64 → 原始 JSON（key/value 与可选 assets）
2) 写入全局配置（/store/set），并落盘资源（如名片图片）
3) 触发 SSE 刷新
Errors:
- 400 invalid_format（解密或 JSON 结构不合法）
- 409 exists_without_overwrite（已存在且未指定 overwrite）
