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

### 5. 授权元数据结构与三种导出情况

**情况1：首次导出 - 无需签名**
- 前端直接调用原导出API，不触发签名应用流程
- 专辑配置不包含signature字段

**情况2：首次导出 - 需要签名且需要授权**
```json
{
  "name": "张三",
  "intro": "键音创作者",
  "cardImagePath": "audioFiles/card1.jpg",
  "authorization": {
    "requireAuthorization": true,
    "contactEmail": "zhang@example.com",
    "contactAdditional": "微信: zhangsan123",
    "authorizedList": [],  // 初始为空，后续授权后添加
    "directExportAuthor": "<资格码1>"  // 当前导出者的资格码
  }
}
```

**情况3：首次导出 - 需要签名但无需授权**
```json
{
  "name": "张三",
  "intro": "键音创作者",
  "cardImagePath": "audioFiles/card1.jpg",
  "authorization": {
    "requireAuthorization": false,
    "contactEmail": "zhang@example.com",
    "contactAdditional": "",
    "authorizedList": [],
    "directExportAuthor": "<资格码1>"
  }
}
```

**再次导出 - 贡献者签名**（无authorization字段）：
```json
{
  "name": "李四",
  "intro": "二次创作者",
  "cardImagePath": "audioFiles/card2.jpg"
  // 无authorization字段
}
```

**再次导出 - 更新directExportAuthor**：
- 每次导出时，原始作者签名的`authorization.directExportAuthor`更新为当前导出者的资格码
- 其他签名条目保持不变

**authorizedList更新**：
- 触发时机：原作者导入授权文件后
- 存储内容：被授权者的资格码
- 用途：前端选择签名时使能/失能签名选项

**签名作者角色识别**：
- 原始作者：signature中包含authorization字段的签名（只有一个）
- 历史贡献作者：signature中的所有其他签名条目
- 直接导出作者：authorization.directExportAuthor对应的签名

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

## API设计

### 端点1: 应用签名配置
**路径**: `POST /keytone_pkg/apply_signature_config`
**功能**: 将签名写入专辑配置

**请求体**:
```json
{
  "albumPath": "/path/to/album",
  "signatureId": "<encrypted_id>",
  "requireAuthorization": true,
  "contactEmail": "author@example.com",
  "contactAdditional": "微信: xxx",
  "updateSignatureContent": true
}
```

**参数说明**:
- `updateSignatureContent`: (boolean) 是否更新签名内容。
  - `true`: 使用新签名数据覆盖专辑中的 Name, Intro, CardImage。
  - `false`: 保留专辑中原有的 Name, Intro, CardImage。
  - 无论取值如何，若原签名包含 Authorization，必须保留；且必须更新 DirectExportAuthor。

**响应**:
```json
{
  "message": "ok",
  "success": true,
  "qualificationCode": "<sha256_hash>"
}
```

### 端点2: 获取专辑签名信息
**路径**: `POST /keytone_pkg/get_album_signature_info`
**功能**: 读取并解析专辑签名信息（前端需求2和4）

**请求体**:
```json
{
  "albumPath": "/path/to/album"
}
```

**响应**:
```json
{
  "message": "ok",
  "data": {
    "hasSignature": true,
    "originalAuthor": {
      "qualificationCode": "<code>",
      "name": "张三",
      "intro": "...",
      "cardImagePath": "audioFiles/card.jpg",
      "isOriginalAuthor": true,
      "requireAuthorization": true,
      "authorizedList": ["<code2>"]
    },
    "contributorAuthors": [
      {
        "qualificationCode": "<code2>",
        "name": "李四",
        "intro": "...",
        "isOriginalAuthor": false
      }
    ],
    "directExportAuthor": {
      "qualificationCode": "<code2>",
      "name": "李四",
      "isOriginalAuthor": false
    },
    "allSignatures": { ... }
  }
}
```

### 端点3: 检查签名是否在专辑中
**路径**: `POST /keytone_pkg/check_signature_in_album`
**功能**: 标记已在专辑中的签名（前端需求3）

**请求体**:
```json
{
  "albumPath": "/path/to/album",
  "signatureId": "<encrypted_id>"
}
```

**响应**:
```json
{
  "message": "ok",
  "isInAlbum": true,
  "qualificationCode": "<code>"
}
```

### 端点4: 检查签名授权状态
**路径**: `POST /keytone_pkg/check_signature_authorization`
**功能**: 检查签名是否有导出授权（前端需求3）

**请求体**:
```json
{
  "albumPath": "/path/to/album",
  "signatureId": "<encrypted_id>"
}
```

**响应**:
```json
{
  "message": "ok",
  "isAuthorized": true,
  "requireAuthorization": true,
  "qualificationCode": "<code>"
}
```

### 端点5: 获取可用签名列表
**路径**: `POST /keytone_pkg/get_available_signatures`
**功能**: 获取所有可用签名及其状态（前端需求3）

**请求体**:
```json
{
  "albumPath": "/path/to/album"
}
```

**响应**:
```json
{
  "message": "ok",
  "signatures": [
    {
      "encryptedId": "<encrypted>",
      "qualificationCode": "<code>",
      "name": "张三",
      "intro": "...",
      "isInAlbum": true,
      "isAuthorized": true,
      "isOriginalAuthor": true
    },
    {
      "encryptedId": "<encrypted2>",
      "qualificationCode": "<code2>",
      "name": "李四",
      "intro": "...",
      "isInAlbum": false,
      "isAuthorized": true,
      "isOriginalAuthor": false
    }
  ]
}
```

## 前端实现

### 组件架构

所有专辑导出流程相关的组件统一存放在 `frontend/src/components/export-flow/` 目录下，包括：
- 导出流程对话框（Export*Dialog.vue）
- 签名选择和管理组件（Signature*Dialog.vue）
- 流程控制逻辑（useExportSignatureFlow.ts）

#### SignatureAuthorsDialog.vue（需求4）
**位置**: `frontend/src/components/export-flow/SignatureAuthorsDialog.vue`
**功能**: 展示专辑签名作者信息
**使用场景**: 用户查看专辑详情时点击"查看签名信息"

**主要功能**:
- 调用 `GetAlbumSignatureInfo` 获取签名信息
- 分区展示：原始作者、直接导出作者、历史贡献作者
- 显示授权状态徽章
- 处理无签名、加载中、错误状态

**UI布局**:
```
┌─────────────────────────────┐
│ [★] 原始作者                │
│  ┌───────────────────────┐  │
│  │ [图片] 名称            │  │
│  │        介绍            │  │
│  │        [需要授权导出]  │  │
│  └───────────────────────┘  │
│                             │
│ [⬇] 直接导出作者            │
│  ┌───────────────────────┐  │
│  │ [图片] 名称            │  │
│  └───────────────────────┘  │
│                             │
│ [👥] 历史贡献作者 (2)       │
│  ┌───────────────────────┐  │
│  │ [图片] 名称            │  │
│  └───────────────────────┘  │
│  ┌───────────────────────┐  │
│  │ [图片] 名称            │  │
│  └───────────────────────┘  │
└─────────────────────────────┘
```

#### SignatureSelectionDialog.vue（需求3）
**位置**: `frontend/src/components/export-flow/SignatureSelectionDialog.vue`
**功能**: 增强的签名选择界面
**使用场景**: 用户导出专辑时选择签名

**主要功能**:
- 调用 `GetAvailableSignatures` 获取签名列表
- 视觉标记：
  - 已在专辑中：蓝色左边框
  - 未授权：置灰 + 锁图标
  - 原始作者：金色星标
- 筛选功能：仅显示已授权 / 仅显示已在专辑中
- 点击未授权签名时提示需要授权

**UI布局**:
```
┌─────────────────────────────────────────┐
│ □ 仅显示已授权  □ 仅显示已在专辑中       │
├─────────────────────────────────────────┤
│ ┌────────────┐ ┌────────────┐           │
│ │[★原始作者] │ │[✓已在专辑] │           │
│ │ 名称       │ │ 名称       │           │
│ │ 介绍       │ │ 介绍 [✓已选]│           │
│ └────────────┘ └────────────┘           │
│ ┌────────────┐                          │
│ │[🔒需要授权]│ (置灰，不可选)           │
│ │ 名称       │                          │
│ └────────────┘                          │
└─────────────────────────────────────────┘
```

### 数据流

#### 需求1：删除"无需签名+需要授权"分支
```
用户选择导出
    ↓
需要签名？
├─ No → 直接调用ExportAlbum()
└─ Yes → 进入签名流程
```

#### 需求2：再次导出时的签名识别（已实现）
```
进入导出流程（exportAlbum）
    ↓
调用 exportFlow.start({ albumPath })
    ↓
GetAlbumSignatureInfo(albumPath)
    ↓
专辑已有签名？
├─ No → state.step = 'confirm-signature' (首次导出流程)
└─ Yes → 检查授权
    ↓
requireAuthorization？
├─ No → state.step = 'picker' (直接进入签名选择)
└─ Yes → state.step = 'auth-gate' (授权门控)
    ↓
TODO: 检查当前用户是否有授权
├─ Yes → 进入签名选择
└─ No → 提示导入授权文件
```

**实现细节**:
- `useExportSignatureFlow.start()` 现在接收 `albumPath` 参数
- 自动调用 `GetAlbumSignatureInfo(albumPath)` 获取签名状态
- 根据返回的 `hasSignature` 和 `requireAuthorization` 决定流程
- 错误时默认按首次导出处理
- 向后兼容旧的测试参数

**代码实现**:

```typescript
// useExportSignatureFlow.ts
const start = async (options: ExportSignatureFlowOptions) => {
  const { albumPath } = options;
  
  // 获取专辑签名信息
  const signatureInfo = await GetAlbumSignatureInfo(albumPath);
  
  // 情况1：专辑无签名 → 首次导出流程
  if (!signatureInfo.hasSignature) {
    state.value.step = 'confirm-signature';
    confirmSignatureDialogVisible.value = true;
    return;
  }
  
  // 情况2：专辑有签名且需要授权
  if (signatureInfo.originalAuthor?.requireAuthorization) {
    state.value.step = 'auth-gate';
    authGateDialogVisible.value = true;
    return;
  }
  
  // 情况3：专辑有签名但不需要授权 → 直接进入签名选择
  state.value.step = 'picker';
  pickerDialogVisible.value = true;
};
```

```typescript
// Keytone_album_page_new.vue
const exportAlbum = async () => {
  const albumPath = setting_store.mainHome.selectedKeyTonePkg;
  
  if (!albumPath) {
    q.notify({ type: 'warning', message: '请先选择一个专辑' });
    return;
  }
  
  // 自动识别三种情况
  await exportFlow.start({ albumPath });
};
```

#### 需求3：签名选择页面增强
```
打开SignatureSelectionDialog
    ↓
GetAvailableSignatures(albumPath)
    ↓
渲染签名卡片
├─ isInAlbum=true → 蓝色边框
├─ isAuthorized=false → 置灰 + 锁图标
├─ isOriginalAuthor=true → 星标
└─ 用户点击签名
    ↓
isAuthorized？
├─ No → 提示需要授权
└─ Yes → 选中签名，允许确认
```

### 类型系统

所有类型定义位于 `frontend/src/types/export-flow.ts`:
- `SignatureAuthorInfo` - 签名作者基本信息
- `AlbumSignatureEntry` - 专辑配置中的签名条目（对应SDK）
- `AlbumSignatureInfo` - 完整的专辑签名信息（API返回）
- `AvailableSignature` - 可选签名信息（包含状态标记）

### API调用封装

所有API函数位于 `frontend/src/boot/query/keytonePkg-query.ts`:
- `GetAlbumSignatureInfo(albumPath)` - 需求2、4使用
- `CheckSignatureInAlbum(albumPath, signatureId)` - 需求3辅助
- `CheckSignatureAuthorization(albumPath, signatureId)` - 需求3辅助
- `GetAvailableSignatures(albumPath)` - 需求3主要使用

### 页面集成

#### Keytone_album_page_new.vue
**集成内容**:
1. 添加"查看签名信息"按钮（badge图标，amber颜色）
2. 导入并使用SignatureAuthorsDialog组件
3. 实现showAlbumSignatureInfo方法，检查专辑选中状态

**代码示例**:
```vue
<!-- 按钮 -->
<q-btn
  icon="badge"
  color="amber"
  @click="showAlbumSignatureInfo"
/>

<!-- 对话框 -->
<SignatureAuthorsDialog
  ref="signatureAuthorsDialogRef"
  :album-path="setting_store.mainHome.selectedKeyTonePkg || ''"
/>
```

**方法实现**:
```typescript
const signatureAuthorsDialogRef = ref<InstanceType<typeof SignatureAuthorsDialog> | null>(null);
const showAlbumSignatureInfo = () => {
  if (!setting_store.mainHome.selectedKeyTonePkg) {
    q.notify({ type: 'warning', message: '请先选择一个专辑' });
    return;
  }
  signatureAuthorsDialogRef.value?.open();
};
```

## Bug修复记录

### Bug #1: 无需签名时仍进入授权对话框

**问题描述**：
用户在导出流程中选择"无需签名"后，仍然会进入"二次创作是否需要授权"的对话框，违反了需求1（删除"无需签名+需要授权"分支）。

**根本原因**：
`useExportSignatureFlow.ts`中的`handleConfirmSignatureSubmit`函数没有检查`needSignature`标志，无论用户选择什么，都会进入授权流程：

```typescript
// 错误代码
const handleConfirmSignatureSubmit = (payload: { needSignature: boolean }) => {
  state.value.flowData = { ...(state.value.flowData ?? {}), needSignature: payload.needSignature };
  confirmSignatureDialogVisible.value = false;
  // 无论最终是否需要签名，都要做二次创作授权判断 ❌
  state.value.step = 'auth-requirement';
  authRequirementDialogVisible.value = true;
};
```

**修复方案**：
在`handleConfirmSignatureSubmit`中添加条件判断，如果用户选择"无需签名"，直接将流程状态设为`done`，不进入授权流程：

```typescript
// 修复后代码
const handleConfirmSignatureSubmit = (payload: { needSignature: boolean }) => {
  state.value.flowData = { ...(state.value.flowData ?? {}), needSignature: payload.needSignature };
  confirmSignatureDialogVisible.value = false;
  
  // 如果选择"无需签名"，直接完成，不进入授权流程 ✅
  if (!payload.needSignature) {
    state.value.step = 'done';
    return;
  }

  // 选择"需要签名"，进入授权判断
  state.value.step = 'auth-requirement';
  authRequirementDialogVisible.value = true;
};
```

**影响范围**：
- 文件：`frontend/src/components/export-flow/useExportSignatureFlow.ts`
- 影响流程：首次导出专辑的签名确认流程
- 验证方法：选择"无需签名"后应直接完成，不显示任何授权相关对话框

### Bug #2: SignatureAuthorsDialog尺寸过大导致溢出

**问题描述**：
`SignatureAuthorsDialog`对话框的尺寸（min-width: 600px, max-width: 800px）在固定窗口尺寸的应用中导致界面溢出，无法正常查看内容。

**根本原因**：
对话框设计时未考虑固定窗口尺寸的约束，使用了较大的固定宽度和较大的字体/图片尺寸。

**修复方案**：

1. **对话框尺寸调整**：
```vue
<!-- 修复前 -->
<q-card style="min-width: 600px; max-width: 800px">

<!-- 修复后 -->
<q-card style="width: 90vw; max-width: 480px; max-height: 85vh">
```

2. **添加滚动支持**：
```vue
<q-card-section 
  v-else-if="signatureInfo" 
  style="max-height: calc(85vh - 100px); overflow-y: auto"
>
```

3. **图片尺寸缩小**：
```vue
<!-- 修复前 -->
style="width: 100px; height: 100px"

<!-- 修复后 -->
style="width: 70px; height: 70px"
```

4. **字体大小调整**：
```vue
<!-- 修复前 -->
<q-icon size="24px" />
<span class="text-h6">
<div class="text-h6">

<!-- 修复后 -->
<q-icon size="20px" />
<span class="text-subtitle1">
<div class="text-subtitle2">
```

5. **间距优化**：
```vue
<!-- 内边距从 q-pa-md 改为 q-pa-sm -->
<q-card-section class="col q-pa-sm">

<!-- 边距从 q-mt-md 改为 q-mt-sm -->
class="author-card q-mt-sm"

<!-- 区块间距从 q-mb-lg 改为 q-mb-md -->
class="author-section q-mb-md"
```

**影响范围**：
- 文件：`frontend/src/components/SignatureAuthorsDialog.vue`
- 改进效果：
  - 对话框宽度适配小屏幕（最大480px）
  - 内容区域可滚动，避免溢出
  - 更紧凑的布局，信息密度更高
  - 保持可读性的同时节省空间

**验证方法**：
1. 在固定窗口尺寸下打开对话框
2. 确认对话框不溢出屏幕
3. 确认所有内容可通过滚动查看
4. 确认文字清晰可读，图片不失真

## 未来扩展

### 可能的增强
1. **签名链验证**：校验authorizedList中的资格码有效性
2. **图片完整性**：存储cardImage的SHA256哈希用于校验
3. **批量签名**：支持一次应用多个签名（联名专辑）
4. **签名撤销**：从authorizedList移除已授权的资格码
5. **签名历史**：记录每次签名操作的时间戳和操作者
6. **签名预览**：在选择前预览签名效果
7. **离线授权**：支持导入离线授权文件
