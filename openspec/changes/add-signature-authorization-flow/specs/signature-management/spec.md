# 签名管理功能规格增量

## ADDED Requirements

### Requirement: 签名授权申请

Normative: The system SHALL allow users to generate authorization request files (.ktauthreq) for signatures they wish to use on albums requiring authorization; the request file MUST contain encrypted authorization UUID hash, encrypted original author qualification code, and encrypted requester signature ID.

#### Scenario: 生成授权申请文件

- **GIVEN** 用户在签名选择对话框中，专辑需要授权，用户有未授权的签名
- **WHEN** 用户点击"授权申请"按钮并选择一个未授权签名
- **THEN** 系统调用 `POST /signature/generate-auth-request` 生成授权申请文件，用户可下载 .ktauthreq 文件

#### Scenario: 查看原始作者联系方式

- **GIVEN** 用户在授权申请对话框中
- **WHEN** 用户进入"导出申请"步骤
- **THEN** 显示原始作者的联系方式（邮箱、附加信息），用户可复制联系方式

---

### Requirement: 签名授权受理

Normative: The system SHALL allow original authors to import authorization request files, validate them against local signatures, and generate authorization grant files (.ktauth); the system MUST verify that the importing user owns a signature matching the request.

#### Scenario: 导入授权申请文件

- **GIVEN** 原始作者在签名管理页面
- **WHEN** 用户点击"授权受理"并导入 .ktauthreq 文件
- **THEN** 系统调用 `POST /signature/parse-auth-request` 解析文件，并在审核步骤显示：
  - 申请信息中的“申请方名称”
  - 匹配到的本地签名（以完整签名列表项展示：图片（若有）+名称+介绍；若匹配多个则先选择再展示所选项）

#### Scenario: 生成授权文件

- **GIVEN** 原始作者已导入并审核授权申请
- **WHEN** 用户选择对应签名并确认授权
- **THEN** 系统调用 `POST /signature/generate-auth-grant` 生成授权文件，用户可下载 .ktauth 文件

#### Scenario: 无匹配签名

- **GIVEN** 原始作者导入授权申请文件
- **WHEN** 本地签名列表中没有匹配的签名
- **THEN** 系统提示"未找到匹配的签名，无法授权"，禁用授权按钮

---

### Requirement: 签名授权验证与导入

Normative: The system SHALL allow users to import authorization grant files (.ktauth) and verify them against album authorization data; upon successful verification, the requester's qualification code MUST be added to the album's authorizedList.

#### Scenario: 导入授权文件成功

- **GIVEN** 用户在授权门控对话框中导入 .ktauth 文件
- **WHEN** 系统调用 `POST /signature/verify-import-auth-grant` 验证成功
- **THEN** 系统调用 `POST /signature/add-to-authorized-list` 将资格码添加到授权列表，提示用户授权成功

#### Scenario: 验证授权文件失败

- **GIVEN** 用户导入了无效的 .ktauth 文件
- **WHEN** 系统调用验证 API 失败
- **THEN** 提示"授权文件验证失败，请确保文件正确"，不更新授权列表

---

### Requirement: 授权列表管理

Normative: The system SHALL provide an API to add qualification codes to an album's authorizedList; the API MUST check for duplicates and preserve existing entries.

#### Scenario: 添加资格码到授权列表

- **GIVEN** 授权验证通过
- **WHEN** 系统调用 `POST /signature/add-to-authorized-list`
- **THEN** 资格码被添加到专辑配置的 authorizedList，若已存在则不重复添加

#### Scenario: 原始作者签名缺失

- **GIVEN** 专辑配置损坏或未包含原始作者签名
- **WHEN** 系统尝试添加到授权列表
- **THEN** 返回错误"未找到原始作者签名，无法添加授权"

## MODIFIED Requirements

### Requirement: 签名选择对话框（修改）

Normative: The signature picker dialog SHALL additionally provide an "Authorization Request" button when the album requires authorization; clicking this button SHALL open the AuthRequestDialog for the user to generate a request file.

#### Scenario: 显示授权申请按钮

- **GIVEN** 签名选择对话框打开，专辑需要授权（requireAuthorization=true）
- **WHEN** 对话框渲染
- **THEN** 显示"授权申请"按钮在操作栏中

#### Scenario: 隐藏授权申请按钮

- **GIVEN** 签名选择对话框打开，专辑不需要授权（requireAuthorization=false）
- **WHEN** 对话框渲染
- **THEN** 不显示"授权申请"按钮
