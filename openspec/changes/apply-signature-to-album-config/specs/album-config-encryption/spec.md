# 专辑配置加密功能规格说明 - 签名字段扩展

## ADDED Requirements

### Requirement: 签名配置字段结构

Normative: The album configuration SHALL support a `signature` field containing a map of qualification codes to signature objects, where each signature object includes name, intro, cardImagePath, and optionally an authorization object for original author signatures.

#### Scenario: 专辑配置包含单个签名

- **GIVEN** 用户首次为专辑应用签名
- **WHEN** 系统写入signature字段
- **THEN** 配置包含一个资格码key，对应的value对象包含name、intro、cardImagePath，原始作者签名额外包含authorization对象

#### Scenario: 专辑配置包含多个签名

- **GIVEN** 专辑已有一个签名，用户应用第二个签名
- **WHEN** 系统更新signature字段
- **THEN** 配置包含两个资格码key，每个对应独立的签名对象，支持原始作者和第三方签名混合存在

---

### Requirement: 授权元数据管理

Normative: Original author signatures SHALL include an `authorization` object with requireAuthorization boolean, contactEmail string, contactAdditional optional string, and authorizedList array of qualification codes; non-author signatures SHALL NOT include this object.

#### Scenario: 原始作者签名包含授权信息

- **GIVEN** 用户导出自己创建的专辑并选择"需要授权"
- **WHEN** 系统写入签名数据
- **THEN** authorization对象包含requireAuthorization=true，有效的contactEmail，可选的contactAdditional，authorizedList初始化为空数组[]

#### Scenario: 授权列表更新

- **GIVEN** 原作者导入授权文件批准第三方导出
- **WHEN** 系统处理授权
- **THEN** 原作者签名的authorizedList数组新增被授权者的资格码，保持其他字段不变

#### Scenario: 非原始作者签名不含授权字段

- **GIVEN** 用户应用从他人导入的签名到专辑
- **WHEN** 系统写入签名数据
- **THEN** 签名对象仅包含name、intro、cardImagePath，不存在authorization字段

---

### Requirement: API端点完善

Normative: The `/keytone_pkg/apply_signature_config` endpoint SHALL accept albumPath, signatureId, requireAuthorization, contactEmail, and contactAdditional parameters, validate authorization completeness, apply the signature, and return the generated qualification code.

#### Scenario: API请求成功应用签名

- **GIVEN** 前端发送包含有效参数的POST请求
- **WHEN** SDK接收请求并调用ApplySignatureToAlbum函数
- **THEN** 返回200状态码和`{ message: "ok", qualificationCode: "<sha256>" }`响应体

#### Scenario: API请求参数不完整

- **GIVEN** 前端发送requireAuthorization=true但contactEmail为空的请求
- **WHEN** SDK接收并验证参数
- **THEN** 返回400错误和`{ message: "error: 需要授权时必须提供联系邮箱" }`

#### Scenario: API请求签名不存在

- **GIVEN** 前端发送的signatureId在签名配置中找不到
- **WHEN** SDK尝试读取签名数据
- **THEN** 返回404错误和`{ message: "error: 签名不存在或已被删除" }`

---

### Requirement: 配置字段加密层次

Normative: The album's `signature` field SHALL be encrypted twice: first using the album-specific encryption key derived from albumUUID, then the signature object itself encrypted with KeyToneAlbumSignatureEncryptionKey before being stored in the outer configuration.

#### Scenario: 双重加密保护签名数据

- **GIVEN** 签名数据已准备写入专辑配置
- **WHEN** 系统执行写入操作
- **THEN** 签名对象先用KeyToneAlbumSignatureEncryptionKey加密为十六进制字符串，该字符串作为signature字段值存入配置JSON，最终整个配置用albumUUID派生密钥加密存储为core文件

#### Scenario: 解密签名数据进行读取

- **GIVEN** 专辑配置已加密存储
- **WHEN** 系统需要读取签名信息
- **THEN** 先用albumUUID派生密钥解密core文件获取配置JSON，再从JSON中提取signature字段值，最后用KeyToneAlbumSignatureEncryptionKey解密获得签名对象

---

### Requirement: 错误处理与降级

Normative: The signature application process SHALL handle missing image files gracefully by skipping the copy and setting cardImagePath to empty string, log errors for debugging, and ensure partial failures do not corrupt the album configuration.

#### Scenario: 图片文件缺失降级处理

- **GIVEN** 签名cardImage路径指向已删除的文件
- **WHEN** 系统尝试复制图片
- **THEN** 记录警告日志但继续执行，cardImagePath设为""，签名的其他字段正常写入

#### Scenario: 专辑配置写入失败回滚

- **GIVEN** 签名数据和图片已准备就绪
- **WHEN** 写入配置时发生IO错误
- **THEN** 返回错误信息，不修改原配置文件，已复制的图片文件保留（不影响下次重试）

---

### Requirement: 配置兼容性

Normative: Albums without signatures SHALL remain valid with no `signature` field; adding a signature SHALL not break existing configuration parsing; multiple signature applications SHALL merge into the same signature map without overwriting unrelated keys.

#### Scenario: 无签名专辑保持兼容

- **GIVEN** 专辑配置不包含signature字段
- **WHEN** 系统加载配置
- **THEN** 正常解析其他字段（audio_pkg_uuid、audio_files等），无错误或警告

#### Scenario: 添加签名不影响现有配置

- **GIVEN** 专辑已有完整的音频和按键映射配置
- **WHEN** 用户应用签名到专辑
- **THEN** 仅新增signature字段，audio_files、keymap等字段保持原样
