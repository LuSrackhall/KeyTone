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
   - 原始作者签名：包含`authorization`对象（requireAuthorization、contactEmail、contactAdditional、authorizedList、directExportAuthor）
   - 非原始作者签名：不包含authorization字段，但可能出现在其他签名的authorizedList中
   - directExportAuthor：记录直接导出作者的资格码，每次导出时更新为当前签名者

### 三种导出情况

**情况1：专辑无签名**
- 视为首次导出流程。
- 用户可以选择"无需签名"（直接导出）或"需要签名"（进入签名流程）。

**情况2：专辑有签名 + 无需授权**
- 强制要求签名（不可选"无需签名"）。
- 再次导出流程：
  - 弹出提示框确认是否继续签名。
  - 直接进入签名选择页面。
  - 允许选择已存在签名，并询问是否更新签名内容。
  - 始终更新 `directExportAuthor`。

**情况3：专辑有签名 + 需要授权**
- 强制要求签名（不可选"无需签名"）。
- 再次导出流程：
  - 弹出提示框确认是否继续签名。
  - 强制先导入授权文件，验证通过后才可进入签名选择。
  - 仅允许选择已授权的签名。
  - 允许选择已存在签名，并询问是否更新签名内容。
  - 始终更新 `directExportAuthor`。

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
        "authorizedList": ["<资格码1>", "<资格码2>", "<资格码3>"],  // 已授权的资格码列表（包含原始作者自身+已授权的第三方）
        "directExportAuthor": "<资格码4>",  // 直接导出作者的资格码（每次导出更新）
        "authorizationUUID": "<nanoid生成的UUID>"  // 授权标识UUID（首次导出时生成）
      }
    }
  }
}
```

### authorizedList 字段说明

**初始化逻辑**：
- 首次导出且选择"需要授权"（`requireAuthorization=true`）时：
  - `authorizedList` 初始化为 `["<原始作者资格码>"]`
  - 原始作者作为授权创建者，天然拥有导出授权
- 首次导出且选择"无需授权"（`requireAuthorization=false`）时：
  - `authorizedList` 初始化为空数组 `[]`

**更新逻辑**：
- 再次导出不会修改 `authorizedList`，保持原有值不变
- 仅通过未来的"授权导入"功能添加新的第三方资格码

### authorizationUUID 字段说明

**生成时机**：
- 首次导出选择"需要签名"时由前端 `nanoid` 生成
- 无论选择"需要授权"还是"无需授权"都会生成此UUID
- 再次导出时沿用已存储的UUID，不重新生成

**未来用途 - 签名授权导出/导入功能**（本次变更仅存储，不实现具体逻辑）：

授权是"签名+专辑"的特定授权，而非通用签名授权：

1. **授权申请文件生成**（从专辑导出流程发起）：
   - 包含字段1：签名解密后原始key的后11位 + authorizationUUID全部字符的组合码的SHA256值
   - 包含字段2：专辑原始作者UUID的SHA256值的后7位的SHA256值

2. **授权文件生成**（原始作者导入申请文件后）：
   - 原始作者选择对应原始签名完成授权
   - 授权文件中：删除原作者UUID的SHA256后7位的SHA256，改为前11位的SHA256值

3. **授权验证**：
   - 通过authorizationUUID参与的组合哈希校验授权文件的有效性
   - 确保授权仅对特定专辑+签名组合生效

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

**端点1**: `POST /keytone_pkg/apply_signature_config` - 应用签名配置
- **输入**：`albumPath`, `signatureId`, `requireAuthorization`, `contactEmail`, `contactAdditional`, `updateSignatureContent`, `authorizationUUID`
  - `authorizationUUID`: 首次导出时由前端nanoid生成，再次导出时传空字符串（SDK沿用已存储的UUID）
- **输出**：`{ message: "ok", qualificationCode: "<sha256>" }`
- **副作用**：专辑配置文件写入signature字段、图片文件复制到专辑目录

**端点2**: `POST /keytone_pkg/get_album_signature_info` - 获取专辑签名信息
- **输入**：`albumPath`
- **输出**：包含originalAuthor、contributorAuthors、directExportAuthor的完整签名信息
- **用途**：前端需求2（再次导出时的签名识别）和需求4（签名作者信息展示）

**端点3**: `POST /keytone_pkg/check_signature_in_album` - 检查签名是否在专辑中
- **输入**：`albumPath`, `signatureId`
- **输出**：`{ isInAlbum: boolean, qualificationCode: string }`
- **用途**：前端需求3（标记已在专辑中的签名）

**端点4**: `POST /keytone_pkg/check_signature_authorization` - 检查签名授权状态
- **输入**：`albumPath`, `signatureId`
- **输出**：`{ isAuthorized: boolean, requireAuthorization: boolean, qualificationCode: string }`
- **用途**：前端需求3（使能/失能签名选项）

**端点5**: `POST /keytone_pkg/get_available_signatures` - 获取可用签名列表
- **输入**：`albumPath`
- **输出**：签名列表，包含每个签名的isInAlbum、isAuthorized、isOriginalAuthor状态
- **用途**：前端需求3（签名选择页面增强）

## 不在范围内

- 授权文件生成与导入机制
- 专辑导出打包流程（仅负责配置写入）
- 前端签名作者信息展示对话框（在前端变更中实现）

## 前端配合需求（已实现）

### 1. 导出流程适配 ✅
- **实现**: 删除"无需签名+需要授权"分支逻辑
- **代码**: "无需签名"直接调用原导出API，不触发签名相关API
- **位置**: 待集成到导出流程中

### 2. 再次导出时的签名识别 ✅
- **API**: `GetAlbumSignatureInfo(albumPath)`
- **实现**: 
  - 读取并解密专辑配置的 signature 字段
  - 返回原始作者、历史贡献作者、直接导出作者信息
  - 根据 requireAuthorization 决定是否需要授权验证
- **类型**: `AlbumSignatureInfo`（frontend/src/types/export-flow.ts）
- **代码**: frontend/src/boot/query/keytonePkg-query.ts
- **流程集成**: ✅ 已集成到useExportSignatureFlow.start()
  - 自动调用GetAlbumSignatureInfo获取签名状态
  - 三种情况自动识别并跳转到对应对话框
  - 错误处理：失败时按首次导出处理

### 3. 签名选择页面增强 ✅
- **API**: `GetAvailableSignatures(albumPath)`
- **组件**: `SignatureSelectionDialog.vue`
- **实现功能**:
  - ✅ 标记已在专辑中的签名（蓝色左边框）
  - ✅ 根据 authorizedList 使能/失能签名（未授权置灰+锁图标）
  - ✅ 显示签名的授权状态徽章
  - ✅ 筛选功能（仅显示已授权/已在专辑中）
  - ✅ 原始作者标记（金色星标）
- **位置**: frontend/src/components/export-flow/SignatureSelectionDialog.vue

### 4. 签名作者信息展示 ✅
- **API**: `GetAlbumSignatureInfo(albumPath)`
- **组件**: `SignatureAuthorsDialog.vue`
- **实现功能**:
  - ✅ 原始作者：signature 字段中包含 authorization 的签名
  - ✅ 历史贡献作者：signature 字段中的所有其他签名
  - ✅ 直接导出作者：authorization.directExportAuthor 对应的签名
  - ✅ 分区展示，带徽章和图标
  - ✅ 处理无签名、加载中、错误状态
- **位置**: frontend/src/components/export-flow/SignatureAuthorsDialog.vue
- **集成状态**: ✅ 已集成到专辑页面，通过"查看签名信息"按钮调用

## Bug修复记录

### Bug #1: 无需签名时仍进入授权对话框 ✅
- **问题**: 用户选择"无需签名"后仍显示授权要求对话框，违反需求1
- **修复**: 在useExportSignatureFlow.ts的handleConfirmSignatureSubmit中添加条件判断
- **结果**: 选择"无需签名"时直接完成导出流程，不进入授权对话框

### Bug #2: SignatureAuthorsDialog尺寸过大 ✅
- **问题**: 对话框宽度600-800px在固定窗口尺寸应用中溢出
- **修复**: 
  - 对话框尺寸调整为90vw，最大480px，最大高度85vh
  - 图片尺寸从100px缩小为70px
  - 字体从text-h6调整为text-subtitle2
  - 添加滚动支持，优化间距
- **结果**: 对话框完美适配固定窗口尺寸，内容清晰可读

## 待确认问题

- **已明确**：使用固定密钥+AES-GCM加密整个signature对象
- **已明确**：资格码使用SHA256(原始签名ID未加密版本)
- **已明确**：图片复制到`audioFiles`目录，与音频文件共用
- **已明确**：authorizedList存储资格码而非加密ID
- **已明确**：directExportAuthor记录直接导出作者，每次导出更新
- **已明确**：删除"无需签名+需要授权"分支，简化为三种情况
