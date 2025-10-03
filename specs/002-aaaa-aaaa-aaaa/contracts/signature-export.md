# Contract: POST /sdk/signatures/:name/export

## 描述
导出指定签名为 .ktsign 文件格式。

## 端点
`POST /sdk/signatures/:name/export`

## 路径参数
- `name` (string, required): 签名名称

## 请求体
```json
{}
```
空对象

## 成功响应

### 200 OK
```json
{
  "fileNameSuggested": "My Signature.ktsign",
  "blobBase64": "SGVsbG8gV29ybGQ..."
}
```

**字段说明:**
- `fileNameSuggested` (string): 建议的文件名
- `blobBase64` (string): Base64 编码的加密签名数据

## 错误响应

### 400 Bad Request - 签名名称缺失
```json
{
  "error": "invalid_request",
  "message": "Signature name is required"
}
```

### 404 Not Found - 签名不存在
```json
{
  "error": "not_found",
  "message": "Signature not found"
}
```

### 500 Internal Server Error - 编码失败
```json
{
  "error": "encode_error",
  "message": "Failed to encode signature file: ..."
}
```

## 处理流程

1. 从路径参数获取签名名称
2. 从 `/store/get` 读取 `signature_manager`
3. 在签名列表中查找匹配的签名（按 name 字段匹配）
4. 构建签名文件载荷（SignatureFilePayload）
5. 使用 XOR 加密和 Base64 编码生成 .ktsign 数据
6. 返回文件名建议和 Base64 数据

## 说明

- 导出的 .ktsign 文件是二进制格式，内部为加密后的 JSON
- JSON 结构等价于全局配置中的单一签名项（key/value 对）
- 前端接收到 base64 数据后需解码为二进制并触发下载
- 保护码(protectCode)在导出文件中加密存储，不可见

## 测试用例

### Happy Path
```bash
POST /sdk/signatures/Zhang%20San/export
Response: 200 OK with valid base64 data
```

### Error Path - 签名不存在
```bash
POST /sdk/signatures/NonExistent/export
Response: 404 Not Found
```
