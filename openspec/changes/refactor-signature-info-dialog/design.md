# 专辑签名信息对话框重构 - 技术设计

## 组件结构

### SignatureAuthorsDialog.vue

**文件位置**: `frontend/src/components/export-flow/SignatureAuthorsDialog.vue`

#### Props

```typescript
interface Props {
  albumPath: string;  // 专辑路径
}
```

#### 状态管理

```typescript
const dialogVisible = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);
const signatureInfo = ref<AlbumSignatureInfo | null>(null);

/** 图片URL缓存 Map<cardImagePath, blobUrl> */
const imageUrlCache = ref<Map<string, string>>(new Map());

/** 资格码指纹缓存 Map<qualificationCode, fingerprint> */
const fingerprintCache = ref<Map<string, string>>(new Map());
```

#### 计算属性

```typescript
// 从 allSignatures 获取原始作者完整签名条目
const originalAuthorEntry = computed<AlbumSignatureEntry | null>(() => {
  if (!signatureInfo.value?.originalAuthor || !signatureInfo.value.allSignatures) {
    return null;
  }
  const qualCode = signatureInfo.value.originalAuthor.qualificationCode;
  return signatureInfo.value.allSignatures[qualCode] || null;
});

// 检查直接导出作者是否与原始作者相同
const isDirectExportAuthorSameAsOriginal = computed(() => {
  // 相同时不重复显示直接导出作者区块
});
```

## 新增 SDK 端点

### GET /keytone_pkg/get_album_file

读取专辑目录中的文件（如签名图片）。

**请求参数**：

```typescript
{
  albumPath: string;      // 专辑路径
  relativePath: string;   // 相对路径，如 "audioFiles/xxx.jpg"
}
```

**响应**：二进制文件内容，Content-Type 自动检测。

**安全措施**：验证请求路径在专辑目录内，防止路径遍历攻击。

## 资格码指纹计算

### 背景

- **资格码**：签名原始ID的SHA256哈希结果（64字符十六进制字符串）
- **资格码指纹**：用于UI展示，保护原始资格码不泄漏

### SDK端实现

资格码指纹在SDK端计算，前端不接触计算逻辑，直接使用返回的 `qualificationFingerprint` 字段。

**SDK代码位置**：`sdk/signature/album.go`

```go
// GenerateQualificationFingerprint 根据资格码生成资格码指纹
//
// TIPS: 资格码指纹用于在保护原始资格码不泄漏的前提下，保证签名的可追溯性。
// 计算方式：将资格码去除第2位（索引1）和第11位（索引10）字符后，计算SHA256哈希。
func GenerateQualificationFingerprint(qualificationCode string) string {
    if len(qualificationCode) < 12 {
        return qualificationCode // 无效输入时返回原值
    }

    // TIPS: 去除第2位（索引1）和第11位（索引10）字符
    modified := qualificationCode[0:1] + qualificationCode[2:10] + qualificationCode[11:]

    // 计算SHA256哈希
    hash := sha256.Sum256([]byte(modified))
    return hex.EncodeToString(hash[:])
}
```

### 数据结构变化

`SignatureAuthorInfo` 结构体新增 `qualificationFingerprint` 字段：

```go
type SignatureAuthorInfo struct {
    QualificationCode        string `json:"qualificationCode"`        // 用于内部数据关联
    QualificationFingerprint string `json:"qualificationFingerprint"` // 用于前端展示
    // ...其他字段
}
```

### 安全性

- 前端不接触指纹计算逻辑
- 资格码仅用于内部数据关联（如查找 `allSignatures`）
- 展示给用户的只有指纹

## UI 区块设计

### 1. 原始作者区块

- 背景色：`bg-amber-1`
- 图标：`star`（琥珀色）
- 内容：
  - 签名卡片（图片 64×64 + 名称 + 介绍）
  - 资格码指纹（monospace 字体 + 复制按钮）
  - 分隔线
  - 联系方式区域
    - 邮箱（email 图标）
    - 其他联系方式（chat 图标）
  - 分隔线
  - 授权状态区域
    - 状态徽章（需要授权/无需授权）
    - 已授权数量徽章
    - 授权UUID（vpn_key 图标）
    - 直接导出作者资格码（file_download 图标）
    - 已授权列表（展开/折叠，checklist 图标）

### 2. 直接导出作者区块

- 背景色：`bg-blue-1`
- 图标：`file_download`（蓝色）
- 条件：仅在 `directExportAuthor` 与 `originalAuthor` 不同时显示
- 内容：
  - 签名卡片（图片 56×56 + 名称 + 介绍）
  - 资格码指纹

### 3. 历史贡献作者区块

- 背景色：`bg-green-1`
- 图标：`group`（绿色）
- 内容：
  - q-list 列表形式
  - 每项：头像 40×40 + 名称 + 介绍 + 复制资格码指纹按钮

### 4. 签名统计摘要

- 背景色：`bg-grey-1`
- 内容：q-chip 芯片组展示统计

## 图片加载机制

### 加载流程

```typescript
async function loadAllImages() {
  // 1. 收集所有需要加载的图片路径
  // 2. 去重
  // 3. 并行调用 GetAlbumFile API
  // 4. 创建 Blob URL 并存入缓存
}

function getImageUrl(cardImagePath: string): string {
  return imageUrlCache.value.get(cardImagePath) || '';
}
```

### 资源清理

```typescript
// 对话框关闭时释放 Blob URL
watch(dialogVisible, (visible) => {
  if (!visible) {
    for (const url of imageUrlCache.value.values()) {
      URL.revokeObjectURL(url);
    }
    imageUrlCache.value.clear();
  }
});
```

## 交互设计

### 复制功能

```typescript
function copyToClipboard(text: string, label: string) {
  navigator.clipboard.writeText(text).then(
    () => {
      $q.notify({
        type: 'positive',
        message: t('exportFlow.signatureInfoDialog.copySuccess', { label }),
        position: 'top',
        timeout: 1500,
      });
    },
    () => {
      $q.notify({
        type: 'negative',
        message: t('exportFlow.signatureInfoDialog.copyFailed'),
        position: 'top',
        timeout: 1500,
      });
    }
  );
}
```

### 横向滚动样式

复用签名列表的统一滚动条样式：

```typescript
const scrollableTextClasses = {
  name: [
    'max-w-full !overflow-x-auto whitespace-nowrap !text-clip',
    'h-5.5 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50 [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
  ],
  intro: [
    'max-w-full !overflow-x-auto whitespace-nowrap',
    'h-4.4 [&::-webkit-scrollbar]:h-0.3 ...',
  ],
};
```

### 展开/折叠

使用 Quasar 的 `q-expansion-item` 组件展示已授权签名列表。

## i18n 国际化

组件使用 `vue-i18n` 实现国际化支持：

```typescript
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
```

**翻译键路径**：`exportFlow.signatureInfoDialog.*`

**主要翻译键**：
| 键                          | 中文                  | 英文                        |
| --------------------------- | --------------------- | --------------------------- |
| `title`                     | 专辑签名信息          | Album Signature Info        |
| `originalAuthor`            | 原始作者              | Original Author             |
| `qualificationFingerprint`  | 资格码指纹            | Qualification Fingerprint   |
| `latestExporterFingerprint` | 最近导出者资格码指纹  | Latest Exporter Fingerprint |
| `contributorAuthors`        | 历史贡献作者          | Contributor Authors         |
| `copySuccess`               | {label}已复制到剪贴板 | {label} copied to clipboard |

## 数据流

```text
用户点击"查看签名信息"
        ↓
open() 方法被调用
        ↓
调用 GetAlbumSignatureInfo(albumPath)
        ↓
signatureInfo 赋值（SDK已计算好所有指纹）
        ↓
调用 loadAllImages() 加载图片
        ↓
模板根据数据渲染各区块
        ↓
originalAuthorEntry 计算属性从 allSignatures 提取完整授权元数据
```

## 图片加载

通过 `GetAlbumFile` API 从专辑目录读取图片：

```typescript
async function loadImage(cardImagePath: string) {
  const blob = await GetAlbumFile(props.albumPath, cardImagePath);
  if (blob) {
    const url = URL.createObjectURL(blob);
    imageUrlCache.value.set(cardImagePath, url);
  }
}
```

配合 `q-img` 的 error slot 处理加载失败情况。

## 错误处理

- 加载失败时显示错误图标和消息
- 提供"重试"按钮，点击重新调用 `open()` 方法

## 样式

```scss
.signature-section {
  border-radius: 8px;
  overflow: hidden;

  & + .signature-section {
    margin-top: 12px;
  }
}
```
