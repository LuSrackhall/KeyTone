# Data Model: 键音专辑签名系统

## Entities

### DigitalSignature
- id: string (uuid)
- name: string (唯一且不可变)
- intro: string? (可选)
- cardImagePath: string? (相对全局配置同级资源目录，可选)
- protectCodeHash: string (只读，创建时生成)
- createdAt: ISO8601 string

### SignatureManager (逻辑集合)
- signatures: DigitalSignature[]

### AlbumSignatureRecord
- signerId: string (引用 DigitalSignature.id)
- signedAt: ISO8601 string[] (该签名在此专辑的导出时间历史)
- allowReexport: boolean (首次作者选择；第一阶段仅记录)

### SignatureFile (.ktsign)
- payload: DigitalSignature（含必要校验字段）
- integrity: string (校验信息)
- version: string

### ExportSession (临时)
- selectedSignatureId: string | null
- allowReexportPreference: boolean | null
