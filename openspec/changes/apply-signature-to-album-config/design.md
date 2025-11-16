# 设计：签名应用到专辑配置

## 架构概览

本设计实现将签名管理系统与专辑配置系统集成，在用户导出专辑时将签名元数据安全写入专辑配置。

### 模块交互
```
前端导出流程
    ↓ HTTP POST /keytone_pkg/apply_signature_config
SDK Server (server.go)
    ↓ 调用签名应用逻辑
签名配置应用模块 (signatureConfig.go) ← 新增
    ↓ 读取签名数据
签名管理模块 (signature/signature.go)
    ↓ 解密签名内容
签名加密模块 (signature/encryption.go)
    ↓ 生成资格码、加密专辑签名
签名专辑加密 (signature/album.go) ← 新增
    ↓ 写入专辑配置
音频包配置模块 (audioPackage/config/audioPackageConfig.go)
```

## 核心组件设计

### 1. 专辑签名加密密钥（signature/album.go）

**常量定义**：
```go
// KeyToneAlbumSignatureEncryptionKey 专辑签名字段专用加密密钥
// 用途：加密专辑配置中的signature字段
// 安全等级：标准（固定密钥，源码可见）
// 长度：32字节（AES-256）
const KeyToneAlbumSignatureEncryptionKey = "KeyTone2024Album_Signature_Field_EncryptionKey_32Bytes"
```

**加密函数**：
```go
// EncryptAlbumSignatureField 加密专辑配置中的签名字段
// 输入：signatureJSON string（整个signature对象的JSON序列化）
// 输出：16进制加密字符串
func EncryptAlbumSignatureField(signatureJSON string) (string, error)

// DecryptAlbumSignatureField 解密专辑配置中的签名字段
// 输入：16进制加密字符串
// 输出：JSON字符串
func DecryptAlbumSignatureField(encryptedHex string) (string, error)
```

**设计理由**：
- 独立密钥：避免与签名管理KeyA/KeyB混淆，职责分离
- 固定密钥：简化实现，专辑配置本身已有UUID派生加密保护
- 复用加密算法：使用现有`signature.EncryptData/DecryptData`（AES-GCM）

### 2. 资格码生成

**函数签名**：
```go
// GenerateQualificationCode 根据原始签名ID生成资格码
// 输入：signatureID string（未加密的UUID）
// 输出：SHA256哈希值（64字符十六进制）
func GenerateQualificationCode(signatureID string) string {
    hash := sha256.Sum256([]byte(signatureID))
    return hex.EncodeToString(hash[:])
}
```

**设计理由**：
- SHA256单向哈希：原始签名ID不可逆推，保护隐私
- 十六进制编码：便于JSON存储和调试
- 确定性：相同签名ID始终生成相同资格码

### 3. 签名数据提取（signatureConfig.go）

**核心流程**：
```go
// ApplySignatureToAlbum 将签名应用到专辑配置
// 参数：
//   - albumPath: 专辑目录路径
//   - encryptedSignatureID: 加密的签名ID（从签名管理系统获取）
//   - requireAuthorization: 是否需要二次导出授权
//   - contactEmail: 联系邮箱（requireAuthorization=true时必需）
//   - contactAdditional: 补充联系信息
// 返回：资格码（SHA256）
func ApplySignatureToAlbum(
    albumPath string,
    encryptedSignatureID string,
    requireAuthorization bool,
    contactEmail string,
    contactAdditional string,
) (string, error)
```

**实现步骤**：
1. 从签名管理配置读取encryptedSignatureID对应的StorageEntry
2. 解密签名ID获取原始UUID
3. 使用动态密钥解密签名Value，获取SignatureData（name, intro, cardImage路径）
4. 生成资格码：`qualificationCode := GenerateQualificationCode(原始UUID)`
5. 处理图片资源（见下节）
6. 构建专辑签名对象
7. 序列化为JSON并加密
8. 写入专辑配置`signature`字段
9. 打印调试日志

### 4. 图片资源处理

**文件复制逻辑**：
```go
// 源路径：config.GetValue("signature").encryptedID.value.cardImage
//        示例：/path/to/ConfigPath/signature/abc123.jpg
// 目标路径：{albumPath}/audioFiles/{newFilename}
//        新文件名：SHA1(qualificationCode + originalFilename + timestamp) + ext
```

**设计考虑**：
- 复用`audioFiles`目录：与音频文件共用，简化目录结构
- 文件名哈希：避免冲突，保持唯一性
- 相对路径存储：`cardImagePath: "audioFiles/xyz789.jpg"`，便于专辑迁移

### 5. 授权元数据结构

**原始作者签名**（第一次导出）：
```json
{
  "name": "张三",
  "intro": "键音创作者",
  "cardImagePath": "audioFiles/card1.jpg",
  "authorization": {
    "requireAuthorization": true,
    "contactEmail": "zhang@example.com",
    "contactAdditional": "微信: zhangsan123",
    "authorizedList": []  // 初始为空，后续授权后添加
  }
}
```

**授权第三方签名**（后续导出）：
```json
{
  "name": "李四",
  "intro": "二次创作者",
  "cardImagePath": "audioFiles/card2.jpg"
  // 无authorization字段
}
```

**authorizedList更新**：
- 触发时机：原作者导入授权文件后
- 存储内容：被授权者的资格码
- 用途：前端选择签名时提示"已授权"状态

## 安全性分析

### 威胁模型
1. **原始签名ID泄露**：通过资格码反推签名创建者
   - 缓解措施：SHA256单向哈希
2. **签名内容窥探**：未授权用户查看专辑配置中的签名信息
   - 缓解措施：整个signature字段AES-GCM加密
3. **图片文件篡改**：替换签名名片图片
   - 当前限制：无文件完整性校验（未来可考虑哈希验证）

### 加密层次
```
专辑配置文件 (package.json / core)
  ↓ 专辑配置加密（enc.DeriveKey based on albumUUID）
[解密后的配置内容]
  ↓ signature字段二次加密（KeyToneAlbumSignatureEncryptionKey）
[签名元数据明文]
```

## 错误处理

### 常见错误场景
1. **签名不存在**：encryptedSignatureID在配置中找不到
   - 返回：`error: "签名不存在或已被删除"`
2. **解密失败**：签名数据损坏或密钥不匹配
   - 返回：`error: "签名数据解密失败"`
3. **图片文件缺失**：cardImage路径指向的文件不存在
   - 行为：跳过图片复制，cardImagePath设为空字符串
4. **专辑配置加载失败**：albumPath无效或配置损坏
   - 返回：`error: "专辑配置加载失败"`
5. **授权信息不完整**：requireAuthorization=true但contactEmail为空
   - 返回：`error: "需要授权时必须提供联系邮箱"`

## 调试与可观测性

### 日志输出
```go
// 成功案例
logger.Info("签名成功应用到专辑配置",
    "album", albumPath,
    "qualificationCode", qualCode,
    "requireAuthorization", requireAuth,
)

// 调试输出（包含未加密数据）
fmt.Printf(`
[专辑签名调试] 签名已成功应用到专辑配置 - 未加密内容：
%s
`, unencryptedSignatureJSON)
```

### 配置验证
- 写入后立即读取signature字段并解密
- 验证JSON格式和必需字段
- 输出调试信息供开发者检查

## 向后兼容性

### 配置迁移
- **无签名专辑**：signature字段不存在，保持原样
- **旧版签名（如有）**：覆盖写入新格式

### 版本标识
在signature对象中添加元数据（可选）：
```json
{
  "_version": "1.0",  // 签名配置格式版本
  "<资格码>": { ... }
}
```

## 性能考虑

### 优化点
1. **图片复制**：异步操作，不阻塞主流程
2. **加密缓存**：如果同一专辑多次应用签名，复用已复制的图片
3. **大文件处理**：图片大小限制（沿用签名创建时的5MB限制）

### 资源消耗
- **内存**：单次操作<10MB（签名数据+图片）
- **磁盘IO**：读签名配置、写专辑配置、复制图片（3次操作）
- **CPU**：加解密操作（AES-GCM + SHA256），<100ms

## 未来扩展

### 可能的增强
1. **签名链验证**：校验authorizedList中的资格码有效性
2. **图片完整性**：存储cardImage的SHA256哈希用于校验
3. **批量签名**：支持一次应用多个签名（联名专辑）
4. **签名撤销**：从authorizedList移除已授权的资格码
5. **签名历史**：记录每次签名操作的时间戳和操作者
