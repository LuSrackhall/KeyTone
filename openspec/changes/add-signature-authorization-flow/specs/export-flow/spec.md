# 专辑导出签名流程规格增量

## ADDED Requirements

### Requirement: 授权申请对话框数据加载（新增）

Normative: The AuthRequestDialog component SHALL automatically load signature data with authorization status when opened; the component MUST filter and display only unauthorized signatures to allow users to request authorization for them.

#### Scenario: 打开授权申请对话框加载数据

- **GIVEN** 用户点击"申请授权"按钮
- **WHEN** AuthRequestDialog 打开（visible=true）
- **THEN** 系统自动调用 GetAvailableSignaturesForExport 和 getSignaturesList 获取签名及授权状态，筛选未授权签名显示在列表中

#### Scenario: 有未授权签名

- **GIVEN** 用户有多个签名，其中部分未被授权
- **WHEN** 授权申请对话框打开
- **THEN** 仅显示未授权的签名供用户选择

#### Scenario: 所有签名均已授权

- **GIVEN** 用户的所有签名均已被授权
- **WHEN** 授权申请对话框打开
- **THEN** 显示提示"所有签名均已授权，您可以使用现有签名或创建新签名后重新申请授权"

#### Scenario: 无任何签名

- **GIVEN** 用户没有任何签名
- **WHEN** 授权申请对话框打开
- **THEN** 显示提示"暂无可用签名，请先创建签名"，并提供"创建签名"按钮

#### Scenario: 创建签名后实时刷新列表

- **GIVEN** 用户在授权申请对话框中点击"创建签名"
- **WHEN** 用户在签名创建对话框中成功创建签名
- **THEN** 授权申请对话框的签名列表自动刷新，显示新创建的签名（无需关闭再重新打开对话框）

---

### Requirement: 授权申请对话框

Normative: The system SHALL provide an AuthRequestDialog component that guides users through a 3-step wizard to generate authorization request files; the dialog MUST display only unauthorized signatures, show the selected signature as a complete card (with image and intro), and provide original author contact information split into email and additional contact.

#### Scenario: 授权申请向导步骤1 - 选择签名

- **GIVEN** 用户从签名选择对话框点击"授权申请"
- **WHEN** AuthRequestDialog 打开
- **THEN** 显示步骤1（带图标的步骤条），列出所有未授权的签名供选择，选中后启用"下一步"按钮

#### Scenario: 授权申请向导步骤2 - 导出申请

- **GIVEN** 用户已选择一个签名
- **WHEN** 用户进入步骤2
- **THEN** 显示已选签名的完整卡片（图片+名称+介绍），分开展示原始作者的邮箱和备用联系方式（各带复制按钮），显示带有"建议先沟通"提示的操作说明，用户点击"导出授权申请"按钮生成 .ktauthreq 文件

#### Scenario: 授权申请文件保存需先选择路径

- **GIVEN** 用户在步骤2 点击"导出授权申请"
- **WHEN** 浏览器支持 File System Access API
- **THEN** 弹出保存对话框（showSaveFilePicker），用户确认保存路径并写入完成后才提示导出成功；若用户取消，则不进入成功步骤

#### Scenario: 授权申请文件保存回退

- **GIVEN** 浏览器不支持 File System Access API
- **WHEN** 用户点击"导出授权申请"
- **THEN** 回退为浏览器下载链接方案，触发下载后提示导出成功

#### Scenario: 授权申请向导步骤3 - 完成

- **GIVEN** 用户已导出授权申请文件
- **WHEN** 导出成功
- **THEN** 显示完成提示，指导用户将文件发送给原始作者

#### Scenario: 步骤条 UI

- **GIVEN** 用户在授权申请对话框的任意步骤
- **WHEN** 查看步骤条
- **THEN** 显示带有步骤编号/勾选图标的圆形指示器，每个步骤下方显示步骤标题，当前步骤高亮显示，已完成步骤显示勾选图标和绿色，步骤之间用连接线连接

#### Scenario: 已选签名展示

- **GIVEN** 用户在步骤2
- **WHEN** 查看已选签名区域
- **THEN** 显示完整的签名卡片，包含：48x48像素的图片（无图片时显示人物图标占位符）、签名名称（加粗）、签名介绍（最多2行，超出省略）

#### Scenario: 联系方式分开展示

- **GIVEN** 用户在步骤2
- **WHEN** 查看原始作者联系方式区域
- **THEN** 分两个区块展示：邮箱区块（带邮件图标和复制按钮）、备用联系方式区块（带联系人图标和复制按钮，仅在有备用联系方式时显示）

#### Scenario: 操作说明提示

- **GIVEN** 用户在步骤2
- **WHEN** 查看操作说明区域
- **THEN** 显示醒目的提示条（建议先与原始作者沟通确认授权意向），然后显示3个步骤的说明文字

---

### Requirement: 授权受理对话框

Normative: The system SHALL provide an AuthGrantDialog component on the signature management page that allows original authors to import authorization requests and generate grant files; the dialog MUST validate requests against local signatures.

#### Scenario: 授权受理向导步骤1 - 导入申请

- **GIVEN** 原始作者点击"授权受理"按钮
- **WHEN** AuthGrantDialog 打开
- **THEN** 显示步骤1，提供文件选择器导入 .ktauthreq 文件

#### Scenario: 授权受理向导步骤2 - 审核授权

- **GIVEN** 用户已导入授权申请文件
- **WHEN** 文件解析成功且找到匹配签名
- **THEN** 显示申请详情和匹配的签名，用户可选择签名并点击"授权"
- **NOTE** 警告提示应明确：授权后，申请方将可以使用**其自己的签名**导出该专辑，而非原始作者的签名

#### Scenario: 授权受理向导步骤3 - 完成

- **GIVEN** 用户点击授权并导出成功
- **WHEN** .ktauth 文件生成完成
- **THEN** 显示完成提示，指导用户将授权文件发送给请求方

---

### Requirement: 导出签名流程状态机扩展

Normative: The useExportSignatureFlow composable SHALL extend its state management to include authorization request flow states; the state machine MUST handle transitions between signature picker and auth request dialog.

#### Scenario: 从签名选择进入授权申请

- **GIVEN** 签名选择对话框打开，用户点击"授权申请"
- **WHEN** 状态机收到 authRequest 事件
- **THEN** 关闭签名选择对话框，打开授权申请对话框

#### Scenario: 授权申请完成后返回

- **GIVEN** 用户完成授权申请流程
- **WHEN** AuthRequestDialog 发出 done 事件
- **THEN** 关闭授权申请对话框，返回签名选择对话框

#### Scenario: 授权申请取消

- **GIVEN** 用户在授权申请对话框点击取消
- **WHEN** AuthRequestDialog 发出 cancel 事件
- **THEN** 关闭授权申请对话框，返回签名选择对话框

## MODIFIED Requirements

### Requirement: 签名选择对话框按钮布局（修改）

Normative: The signature picker dialog's action buttons in the search area SHALL be arranged in a single row with consistent spacing; when authorization is required, the "Authorization Request" button SHALL appear on the left, "Import Authorization" in the middle, and "Create Signature" on the right; when authorization is not required, only the "Create Signature" button SHALL be displayed aligned to the right.

#### Scenario: 需要授权时的按钮布局

- **GIVEN** 签名选择对话框打开，专辑需要授权（requireAuthorization=true）
- **WHEN** 对话框渲染
- **THEN** 三个按钮在同一行显示：左侧"申请授权"、中间"导入授权"、右侧"创建签名"

#### Scenario: 无需授权时的按钮布局

- **GIVEN** 签名选择对话框打开，专辑不需要授权（requireAuthorization=false）
- **WHEN** 对话框渲染
- **THEN** 仅显示"创建签名"按钮，右对齐

---

### Requirement: 授权门控对话框增强（修改）

Normative: The ExportAuthorizationGateDialog SHALL additionally support importing .ktauth authorization files and verifying them via API; the dialog MUST display author contact information split into email and additional contact (each with copy button); upon successful verification, the dialog MUST update the album's authorized list with the requester's qualification code and allow the user to continue.

#### Scenario: 授权门控联系方式展示

- **GIVEN** 用户在授权门控对话框
- **WHEN** 查看作者联系方式区域
- **THEN** 分两个区块展示：邮箱区块（带邮件图标和复制按钮）、备用联系方式区块（带联系人图标和复制按钮，仅在有备用联系方式时显示）

#### Scenario: 导入 .ktauth 授权文件

- **GIVEN** 用户在授权门控对话框选择 .ktauth 文件
- **WHEN** 用户点击"继续"
- **THEN** 系统验证授权文件
  - 若提供了请求方签名ID，则验证该特定签名
  - 若未提供请求方签名ID（如直接导入），系统自动遍历本地所有签名进行匹配验证
- **AND** 验证成功后获取请求方签名的资格码，并将其添加到专辑的授权列表中，最后关闭对话框

#### Scenario: 授权验证失败提示

- **GIVEN** 用户导入无效的授权文件
- **WHEN** 验证失败
- **THEN** 显示错误提示"授权文件验证失败，请确保文件正确"，不关闭对话框

#### Scenario: 授权文件保存需先选择路径

- **GIVEN** 用户在授权门控对话框点击"同意并导出"
- **WHEN** 浏览器支持 File System Access API
- **THEN** 弹出保存对话框（showSaveFilePicker），用户确认保存路径并写入完成后才提示导出成功；若用户取消，则不进入成功步骤

#### Scenario: 授权文件保存回退

- **GIVEN** 浏览器不支持 File System Access API
- **WHEN** 用户点击"同意并导出"
- **THEN** 回退为浏览器下载链接方案，触发下载后提示导出成功
