# 变更：将签名应用到专辑配置

## 为什么需要

当用户在导出专辑时选择"需要签名"，系统需要将所选签名的信息正式写入专辑配置文件。这是专辑导出签名流程的核心环节，确保：
1. 专辑配置包含完整的签名元数据（名称、介绍、图片等）
2. 签名内容使用专用对称密钥加密，保护敏感信息
3. 使用资格码（签名ID的SHA256哈希值）作为标识，保护原始签名ID不被泄露
4. 支持原始作者的导出授权机制和授权记录

## 变更内容

本变更实现SDK侧的签名配置应用功能，具体包括：

### 核心功能
1. **定义专用签名加密密钥**：为专辑配置中的签名字段创建独立的固定对称加密密钥（区别于现有的签名管理KeyA/KeyB）
2. **生成资格码**：将原始签名ID（未加密的UUID）进行SHA256哈希，作为专辑配置中的签名标识
3. **签名数据提取与加密**：从签名配置中解密获取签名详细信息，按专辑签名字段格式重新加密并写入
4. **图片资源复制**：将签名引用的名片图片复制到专辑目录，更新路径为专辑内相对路径
5. **授权元数据处理**：
   - 原始作者签名：包含`authorization`对象（requireAuthorization、contactEmail、contactAdditional、authorizedList）
   - 非原始作者签名：不包含authorization字段，但可能出现在其他签名的authorizedList中

### 数据格式
专辑配置中的签名字段结构（加密存储）：
```json
{
  "signature": {
    "<资格码1>": {
      "name": "签名名称",
      "intro": "个人介绍",
      "cardImagePath": "audioFiles/相对路径.jpg",
      "authorization": {  // 仅原始作者签名包含此字段
        "requireAuthorization": true,
        "contactEmail": "author@example.com",
        "contactAdditional": "补充联系方式",
        "authorizedList": ["<资格码2>", "<资格码3>"]  // 已授权的第三方签名资格码列表
      }
    }
  }
}
```

### 加密方案
- **签名加密密钥**：固定常量`KeyToneAlbumSignatureEncryptionKey`（32字节）
- **加密算法**：AES-256-GCM（复用现有`signature.EncryptData/DecryptData`函数）
- **存储方式**：整个signature对象JSON序列化后加密，存储为16进制字符串

### 终端调试日志
签名成功写入专辑配置后，在终端输出调试日志：
```
[专辑签名调试] 签名已成功应用到专辑配置 - 未加密内容：
{
  "<资格码>": {
    "name": "...",
    "intro": "...",
    ...
  }
}
```

## 影响范围

### 新增能力
- 扩展现有能力：`signature-management`（签名数据读取）、`album-config-encryption`（专辑配置写入）
- 跨模块集成：签名管理 + 专辑配置加密

### 受影响代码
- **新增**：`sdk/audioPackage/config/signatureConfig.go`（签名应用逻辑）
- **新增**：`sdk/signature/album.go`（专辑签名专用加密函数）
- **修改**：`sdk/server/server.go`（完善`/keytone_pkg/apply_signature_config` API实现）
- **依赖**：现有签名管理模块（解密签名数据）、音频包配置模块（写入配置）

### API接口
完善现有的`POST /keytone_pkg/apply_signature_config`端点：
- **输入**：`albumPath`, `signatureId`, `requireAuthorization`, `contactEmail`, `contactAdditional`
- **输出**：`{ message: "ok", qualificationCode: "<sha256>" }`
- **副作用**：专辑配置文件写入signature字段、图片文件复制到专辑目录

## 不在范围内

- 前端导出流程UI实现（已在其他变更中定义）
- 签名验证逻辑（导入或播放时的校验）
- 授权文件生成与导入机制
- 专辑导出打包流程（仅负责配置写入）

## 待确认问题

- **已明确**：使用固定密钥+AES-GCM加密整个signature对象
- **已明确**：资格码使用SHA256(原始签名ID未加密版本)
- **已明确**：图片复制到`audioFiles`目录，与音频文件共用
- **已明确**：authorizedList存储资格码而非加密ID
