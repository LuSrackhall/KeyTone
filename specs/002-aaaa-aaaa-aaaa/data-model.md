# Data Model: 键音专辑签名系统

## Entities

### DigitalSignature
- id: string (uuid, required, immutable)
- name: string (1..64, required, unique, immutable)
- intro: string? (0..512)
- cardImagePath: string? (相对全局配置同级资源目录，可选)
- protectCodeHash: string (sha256, required, immutable)
- createdAt: ISO8601 string (required)

校验：
- name 在本地唯一；导入重复时允许覆盖 intro/cardImagePath，但 id/name/protectCodeHash 不变。
- protectCodeHash 不对用户暴露编辑入口。

 
### SignatureManager (逻辑集合)
- signatures: DigitalSignature[]

 
### AlbumSignatureRecord
- signerId: string (引用 DigitalSignature.id)
- signerName: string
- protectCodeHash: string
- signedAt: ISO8601 string[] (该签名在此专辑的导出时间历史)
- allowReexport: boolean (首次作者选择；第一阶段仅记录)

 
### SignatureFile (.ktsign)
- version: string (e.g., "1")
- payload: DigitalSignature（剔除本地绝对路径敏感信息）
- assets: { cardImage?: string (base64 或相对路径) }
- integrity: string (预留：第一阶段可为空或弱校验)

 
### ExportSession (临时)
- selectedSignatureId: string | null
- allowReexportPreference: boolean | null

 
## Relationships
- DigitalSignature 1..* → AlbumSignatureRecord（按专辑存储多条记录）
- SignatureFile 是 DigitalSignature 的可移植副本

 
## State & Transitions
- 导入：SignatureFile → DigitalSignature（若 id/name 已存在 → 覆盖 intro/cardImagePath）
- 导出：DigitalSignature → SignatureFile
- 专辑导出：选择 DigitalSignature → 写入/追加 AlbumSignatureRecord
