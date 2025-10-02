# Data Model: 键音专辑签名系统

## Entities

### DigitalSignature（本地签名实体）

- name: string (1..64, required, unique, immutable, 唯一标识)
- intro: string? (0..512)
- cardImagePath: string? (相对全局配置同级资源目录，uuid.png/jpg 等；可选)
- protectCode: string (nanoid 生成的 uuid, required, immutable, 前端生成)
- createdAt: ISO8601 string (required)

校验：

- name 在本地唯一，作为签名的唯一标识；导入重复时允许覆盖 intro/cardImagePath，但 name/protectCode 不变。
- protectCode 不对用户暴露编辑入口。
- cardImagePath 使用 nanoid uuid 命名，便于覆盖式替换（无需删除旧图片）。

### SignatureManager (逻辑集合)

- signatures: DigitalSignature[]

### AlbumSignatureRecord（专辑内签名记录）

- signerName: string (引用 DigitalSignature.name，唯一标识)
- intro: string?
- cardImagePath: string?
- signedAt: ISO8601 string[] (该签名在此专辑的导出时间戳历史)
- authorization?: AuthorizationBlock (仅原始作者的签名中包含；第二阶段实现)

AuthorizationBlock（第二阶段占位）：

- authCode: string (签名授权码；默认为固定 sha256，匹配即允许二次导出；不匹配则需授权码校验)
- authorizedList: string[] (授权对象资格码列表；每个对应一个三方签名；具体生成/校验规则待第二阶段定义)

### SignatureFile (.ktsign)

- version: string (e.g., "1")
- payload: DigitalSignature（剔除本地绝对路径敏感信息）
- assets: { cardImage?: string (base64 或相对路径) }
- integrity: string (预留：第一阶段可为空或弱校验)

### ExportSession (临时)

- selectedSignatureName: string | null (选中的签名名称)
- allowReexportPreference: boolean | null

## Relationships

- DigitalSignature 1..* → AlbumSignatureRecord（按专辑存储多条记录）
- SignatureFile 是 DigitalSignature 的可移植副本

## Storage Layout（存储布局）

### 全局配置（明文存储）

位于全局配置文件（如 KeyToneSetting.json）的 signatures 字段：

```json
{
  "signatures": {
    "<encrypt(protectCode)>": {
      "name": "签名名称",
      "intro": "可选简介",
      "cardImagePath": "uuid.png",
      "createdAt": "ISO8601"
    }
  }
}
```

**说明**：

- key: `encrypt(protectCode)` - 对称加密后的保护码
- value: 明文 JSON，包含 name、intro、cardImagePath、createdAt 字段（无 id）
- name 字段作为签名的唯一标识
- cardImagePath 的文件名使用 nanoid 生成的 uuid，便于覆盖式替换

### 专辑配置（对称加密存储）

位于键音专辑配置文件的 signatures 字段：

```json
{
  "signatures": {
    "<encrypt(sha256(decrypt(protectCode) + name))>": "<encrypt(JSON_payload)>"
  }
}
```

**说明**：

- key: `encrypt(sha256(decrypt(protectCode) + name))` - 先解密保护码，与 name 拼接后计算 SHA-256 哈希，最后对哈希值加密
- value: `encrypt(JSON_payload)` - 对称加密后的 JSON 字符串
- JSON_payload（明文结构，加密前）:

```json
{
  "name": "签名名称",
  "intro": "可选简介",
  "cardImagePath": "uuid.png",
  "signedAt": ["ISO8601", ...],
  "authorization": {
    "authCode": "固定sha256或授权码",
    "authorizedList": ["资格码1", "资格码2"]
  }
}
```

**authorization 字段**（第二阶段强依赖）：

- 仅在原始作者签名中包含
- authCode: 签名授权码（默认为固定 sha256 = 允许二次导出；不匹配 = 需授权码校验）
- authorizedList: 资格码列表（存储通过授权校验的三方签名资格码）
- 第一阶段：可省略 authorization 或置空

### 对称加密规范（占位）

- 算法：AES-GCM（推荐）或其他对称加密算法（待 research.md 确定）
- 密钥来源：待定（可能为应用级固定密钥、用户密码派生、或其他策略）
- 编码：Base64（便于 JSON 序列化）

## State & Transitions

- 导入：SignatureFile → DigitalSignature（若 name 已存在 → 覆盖 intro/cardImagePath，但 name/protectCode 不变）
- 导出：DigitalSignature → SignatureFile
- 专辑签名：选择 DigitalSignature → 写入/追加 AlbumSignatureRecord 到专辑配置（key 和 value 均对称加密）
