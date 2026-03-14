# 签名管理功能规格说明 - 专辑签名应用扩展

## ADDED Requirements

### Requirement: 专辑签名数据加密

Normative: The system SHALL provide dedicated encryption functions for album signature fields using a fixed symmetric key `KeyToneAlbumSignatureEncryptionKey`, encrypting the entire signature object as JSON and storing it as hexadecimal ciphertext in the album configuration.

#### Scenario: 加密专辑签名数据

- **GIVEN** 签名数据已从签名管理系统解密并准备写入专辑配置
- **WHEN** 系统调用 `EncryptAlbumSignatureField` 并传入签名JSON字符串
- **THEN** 函数使用固定密钥和AES-256-GCM算法加密数据，返回16进制编码的密文

#### Scenario: 解密专辑签名数据

- **GIVEN** 专辑配置中包含加密的签名字段
- **WHEN** 系统需要读取签名信息时调用 `DecryptAlbumSignatureField`
- **THEN** 函数解密16进制密文并返回JSON字符串，可解析为签名对象

---

### Requirement: 资格码生成

Normative: The system SHALL generate qualification codes by computing SHA256 hash of the original (unencrypted) signature ID, ensuring one-way mapping that protects the original ID while maintaining deterministic identification.

#### Scenario: 生成签名资格码

- **GIVEN** 原始签名ID为未加密的UUID字符串
- **WHEN** 系统调用 `GenerateQualificationCode` 传入原始ID
- **THEN** 函数返回64字符的SHA256十六进制哈希值，相同ID总是生成相同资格码

#### Scenario: 资格码用于专辑配置

- **GIVEN** 资格码已生成
- **WHEN** 签名写入专辑配置时使用资格码作为key
- **THEN** 专辑配置中的signature对象以资格码索引，无法反推原始签名ID

---

### Requirement: 签名数据提取与转换

Normative: The system SHALL extract signature data from the encrypted signature management configuration, decrypt using dynamic keys, and transform into album signature format including name, intro, cardImagePath, and optional authorization metadata.

#### Scenario: 提取原始作者签名数据

- **GIVEN** 用户选择自己创建的签名应用到专辑且要求授权
- **WHEN** 系统从签名配置解密签名Value获取SignatureData
- **THEN** 生成的专辑签名对象包含authorization字段，requireAuthorization设为true，contactEmail和contactAdditional填入用户提供的值，authorizedList初始化为空数组

#### Scenario: 提取非原始作者签名数据

- **GIVEN** 用户选择导入的第三方签名应用到专辑
- **WHEN** 系统解密签名数据
- **THEN** 生成的专辑签名对象不包含authorization字段，仅包含name、intro和cardImagePath

---

### Requirement: 签名图片资源管理

Normative: The system SHALL copy signature card images from the signature storage directory to the album's audioFiles directory, generating unique filenames using SHA1(qualificationCode + originalFilename + timestamp) and storing relative paths in the album signature configuration.

#### Scenario: 复制签名名片图片

- **GIVEN** 签名数据包含cardImage路径指向有效图片文件
- **WHEN** 系统应用签名到专辑时
- **THEN** 图片文件被复制到`{albumPath}/audioFiles/`，新文件名为哈希值+扩展名，cardImagePath更新为相对路径`audioFiles/{newFilename}`

#### Scenario: 签名图片缺失处理

- **GIVEN** 签名的cardImage路径指向的文件不存在
- **WHEN** 系统尝试复制图片
- **THEN** 复制操作跳过并记录警告日志，cardImagePath字段设为空字符串，不影响其他签名数据写入

---

### Requirement: 专辑配置签名字段写入

Normative: The system SHALL write the encrypted signature object to the album configuration under the `signature` key, preserving existing album data and supporting multiple signatures indexed by qualification codes.

#### Scenario: 首次为专辑添加签名

- **GIVEN** 专辑配置中不存在signature字段
- **WHEN** 系统应用签名并写入配置
- **THEN** 创建新的signature对象，包含单个资格码key及其对应的签名数据，整体加密后存储

#### Scenario: 为已有签名的专辑添加新签名

- **GIVEN** 专辑配置已包含加密的signature字段
- **WHEN** 系统应用第二个签名
- **THEN** 解密现有signature对象，添加新资格码及其数据，重新加密后覆盖原字段

---

### Requirement: 调试信息输出

Normative: The system SHALL output unencrypted signature data to the terminal after successfully applying a signature to an album configuration, prefixed with `[专辑签名调试]`, to facilitate development and troubleshooting.

#### Scenario: 输出签名应用调试日志

- **GIVEN** 签名已成功写入专辑配置
- **WHEN** 操作完成时
- **THEN** 终端输出包含标题"[专辑签名调试] 签名已成功应用到专辑配置 - 未加密内容："，后跟格式化的JSON对象（包含资格码、name、intro等字段）

#### Scenario: 调试日志不影响生产数据

- **GIVEN** 调试日志输出到标准输出
- **WHEN** 日志打印时
- **THEN** 仅用于开发者观察，不写入配置文件或持久化存储，不影响专辑数据安全性
