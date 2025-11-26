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

Normative: Original author signatures SHALL include an `authorization` object with requireAuthorization boolean, contactEmail string, contactAdditional optional string, authorizedList array of qualification codes, and directExportAuthor string storing the current exporter's qualification code; non-author signatures SHALL NOT include this object.

#### Scenario: 原始作者签名包含授权信息

- **GIVEN** 用户导出自己创建的专辑并选择"需要授权"
- **WHEN** 系统写入签名数据
- **THEN** authorization对象包含requireAuthorization=true，有效的contactEmail，可选的contactAdditional，authorizedList初始化为空数组[]，directExportAuthor设置为当前签名者的资格码

#### Scenario: 授权列表更新

- **GIVEN** 原作者导入授权文件批准第三方导出
- **WHEN** 系统处理授权
- **THEN** 原作者签名的authorizedList数组新增被授权者的资格码，保持其他字段不变

#### Scenario: 非原始作者签名不含授权字段

- **GIVEN** 用户应用从他人导入的签名到专辑
- **WHEN** 系统写入签名数据
- **THEN** 签名对象仅包含name、intro、cardImagePath，不存在authorization字段

#### Scenario: 直接导出作者记录更新

- **GIVEN** 用户再次导出已签名的专辑
- **WHEN** 用户选择签名并执行导出
- **THEN** 原始作者签名的authorization.directExportAuthor更新为当前导出者的资格码

---

### Requirement: API端点 - 应用签名配置

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

### Requirement: API端点 - 获取专辑签名信息

Normative: The `/keytone_pkg/get_album_signature_info` endpoint SHALL read, decrypt, and parse the album's signature field, identify original author, contributors, and direct export author, and return structured signature information for frontend display.

#### Scenario: 获取包含签名的专辑信息

- **GIVEN** 专辑配置包含加密的signature字段
- **WHEN** 前端请求获取签名信息
- **THEN** 返回包含originalAuthor、contributorAuthors、directExportAuthor和allSignatures的完整信息

#### Scenario: 获取无签名的专辑信息

- **GIVEN** 专辑配置不包含signature字段
- **WHEN** 前端请求获取签名信息
- **THEN** 返回hasSignature=false和空的签名列表

---

### Requirement: API端点 - 检查签名状态

Normative: The system SHALL provide endpoints to check if a signature exists in the album (`/check_signature_in_album`) and whether it has export authorization (`/check_signature_authorization`), enabling frontend to mark signatures and control UI state.

#### Scenario: 检查签名是否在专辑中

- **GIVEN** 用户在签名选择页面
- **WHEN** 前端调用check_signature_in_album检查签名
- **THEN** 返回isInAlbum布尔值和qualificationCode，用于UI标记

#### Scenario: 检查签名是否有授权

- **GIVEN** 专辑需要授权导出
- **WHEN** 前端调用check_signature_authorization检查签名
- **THEN** 返回isAuthorized（是否可导出）、requireAuthorization（是否需要授权）和qualificationCode

#### Scenario: 首次导出时所有签名都有授权

- **GIVEN** 专辑不包含signature字段
- **WHEN** 检查任意签名的授权状态
- **THEN** 返回isAuthorized=true, requireAuthorization=false

---

### Requirement: API端点 - 获取可用签名列表

Normative: The `/keytone_pkg/get_available_signatures` endpoint SHALL retrieve all signatures from user configuration, enrich each with album-specific metadata (isInAlbum, isAuthorized, isOriginalAuthor), and return a complete list for signature selection UI.

#### Scenario: 获取首次导出的可用签名

- **GIVEN** 专辑无签名字段，用户配置有3个签名
- **WHEN** 前端请求可用签名列表
- **THEN** 返回3个签名，全部isAuthorized=true, isInAlbum=false

#### Scenario: 获取需要授权专辑的可用签名

- **GIVEN** 专辑需要授权，authorizedList包含1个资格码，用户配置有3个签名
- **WHEN** 前端请求可用签名列表
- **THEN** 返回3个签名，其中原始作者和授权签名isAuthorized=true，其他isAuthorized=false

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

### Requirement: 导出场景与流程控制

Normative: The system SHALL support three distinct export scenarios for both initial and subsequent exports: unsigned albums, signed with authorization, and signed without authorization. The system SHALL enforce strict flow control for re-exports based on the original author's configuration.

#### Scenario: 首次导出无需签名

- **GIVEN** 专辑配置不包含signature字段，用户选择"无需签名"
- **WHEN** 前端触发导出流程
- **THEN** 跳过签名应用API，直接调用原导出API，不创建signature字段

#### Scenario: 首次导出需要签名

- **GIVEN** 专辑配置不包含signature字段，用户选择"需要签名"
- **WHEN** 用户选择是否需要授权
- **THEN** 系统应用签名到专辑，创建signature字段，包含原始作者签名及authorization对象（根据选择设置requireAuthorization），directExportAuthor设为当前签名者资格码

#### Scenario: 再次导出无签名专辑

- **GIVEN** 专辑配置不包含signature字段（即原始作者未签名）
- **WHEN** 用户再次导出
- **THEN** 视为首次导出流程，用户可以选择"无需签名"或"需要签名"。

#### Scenario: 再次导出有签名专辑

- **GIVEN** 专辑已有签名
- **WHEN** 用户再次导出
- **THEN** 系统强制要求签名，不提供"无需签名"选项。
- **AND** 弹出提示对话框："该键音专辑原始作者明确了该键音包的二次导出必须实施签名..."。

#### Scenario: 再次导出需要签名且无需授权

- **GIVEN** 专辑已有签名且requireAuthorization=false，用户确认提示对话框
- **WHEN** 进入签名选择流程
- **THEN** 直接进入签名选择界面。用户可以选择新签名或已存在于专辑中的签名。
- **AND** 每次进入界面时，必须重新加载签名列表及专辑内的签名状态，确保标签显示最新。
- **AND** 签名选择列表中，原始作者签名标记为"原始作者"，其他已存在于专辑中的签名标记为"贡献作者"。
- **AND** 如果选择已存在的签名，弹出二次确认对话框："所选签名与专辑中已存在的该签名内容有变更。是否确认更新该签名的信息（介绍、图片）？"。
- **AND** 如果确认更新，且该签名是原始作者签名，仅更新name、intro、cardImagePath，保持authorization字段不变。
- **AND** 无论选择何种签名，导出时directExportAuthor更新为当前签名者的资格码。

#### Scenario: 再次导出需要签名且需要授权

- **GIVEN** 专辑已有签名且requireAuthorization=true，用户确认提示对话框
- **WHEN** 进入签名选择流程
- **THEN** 首先弹出授权文件导入对话框。
- **AND** 授权验证通过后，进入签名选择界面。
- **AND** 每次进入界面时，必须重新加载签名列表及专辑内的签名状态。
- **AND** 界面仅允许选择已获得授权的签名（在authorizedList中）或原始作者签名（如果是原作者本人）。
- **AND** 签名选择列表中，原始作者签名标记为"原始作者"，其他已存在于专辑中的签名标记为"贡献作者"。
- **AND** 如果选择已存在的签名，弹出二次确认对话框（同上）。
- **AND** 如果确认更新，且该签名是原始作者签名，仅更新基本信息，保持authorization字段不变。
- **AND** 导出时directExportAuthor更新为当前签名者的资格码。

### Requirement: 签名更新与直接导出作者

Normative: When re-exporting, the `directExportAuthor` field in the original author's authorization object SHALL always be updated to the current exporter's qualification code. When updating an existing signature, the system SHALL allow updating content (name, intro, image) while preserving authorization metadata for the original author. The API SHALL accept an `updateSignatureContent` flag to control this behavior.

#### Scenario: 智能更新检测

- **GIVEN** 用户选择了一个已存在于专辑中的签名
- **WHEN** 系统检查签名状态
- **THEN** 系统对比本地签名配置与专辑内签名的 Name, Intro 和 CardImage（通过SHA256哈希比对）
- **AND** 仅当检测到实际内容变更时，才弹出"更新确认"对话框
- **AND** 如果内容完全一致，系统自动选择"不更新"模式，跳过确认对话框，直接进行后续流程（仅更新DirectExportAuthor）

#### Scenario: 更新已存在的签名内容

- **GIVEN** 用户选择了一个已存在于专辑中的签名，并确认更新 (updateSignatureContent=true)
- **WHEN** 系统应用签名
- **THEN** 使用当前本地签名配置中的name、intro、cardImagePath覆盖专辑中的对应字段
- **AND** 如果是原始作者签名，**必须**保留原有的authorization对象（包括requireAuthorization, authorizedList等），但更新directExportAuthor
- **AND** 如果签名包含新的图片，系统将复制新图片到专辑目录，并**删除**旧的图片文件（如果存在且路径不同），避免垃圾文件堆积。

#### Scenario: 选择已存在签名但不更新

- **GIVEN** 用户选择了一个已存在于专辑中的签名，但在确认对话框中选择"不更新" (updateSignatureContent=false)
- **WHEN** 系统应用签名
- **THEN** 保持专辑中该签名的name、intro、cardImagePath不变
- **AND** 系统**跳过**图片复制操作，不产生新的图片文件。
- **BUT** 仍然更新原始作者签名中的authorization.directExportAuthor为当前导出者的资格码

#### Scenario: 始终更新直接导出作者

- **GIVEN** 任何再次导出操作（只要选择了签名）
- **WHEN** 导出完成前
- **THEN** 原始作者签名的authorization.directExportAuthor字段必须被更新为当前所选签名的资格码

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
