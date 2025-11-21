# 前端API使用指南

本文档说明如何使用SDK提供的签名相关API实现前端的四个需求。

## API基础信息

**基础URL**: `http://localhost:8088` (或配置的SDK服务地址)

**通用错误响应**:
```json
{
  "message": "error: 错误描述"
}
```

## 需求1: 删除"无需签名+需要授权"分支

### 实现要点
- 前端导出流程中，如果用户选择"无需签名"，直接调用原导出API
- 不再需要授权相关的UI和逻辑
- 简化为三种情况：无需签名、需要签名+需要授权、需要签名+无需授权

### 代码示例
```typescript
// 导出专辑逻辑
async function exportAlbum(albumPath: string, needSignature: boolean) {
  if (!needSignature) {
    // 情况1: 无需签名，直接导出
    return await callOriginalExportAPI(albumPath);
  }
  
  // 情况2和3: 需要签名，进入签名流程
  return await showSignatureSelectionDialog(albumPath);
}
```

---

## 需求2: 再次导出时的签名识别

### API端点
`POST /keytone_pkg/get_album_signature_info`

### 请求示例
```typescript
interface GetAlbumSignatureInfoRequest {
  albumPath: string;
}

interface AlbumSignatureInfo {
  hasSignature: boolean;
  originalAuthor?: SignatureAuthorInfo;
  contributorAuthors: SignatureAuthorInfo[];
  directExportAuthor?: SignatureAuthorInfo;
  allSignatures: Record<string, AlbumSignatureEntry>;
}

interface SignatureAuthorInfo {
  qualificationCode: string;
  name: string;
  intro: string;
  cardImagePath: string;
  isOriginalAuthor: boolean;
  requireAuthorization?: boolean;
  authorizedList?: string[];
}

async function getAlbumSignatureInfo(albumPath: string): Promise<AlbumSignatureInfo> {
  const response = await fetch('http://localhost:8088/keytone_pkg/get_album_signature_info', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ albumPath })
  });
  
  const data = await response.json();
  if (data.message !== 'ok') {
    throw new Error(data.message);
  }
  
  return data.data;
}
```

### 使用场景
```typescript
// 用户再次导出专辑时，检查是否需要授权
async function handleReExport(albumPath: string) {
  const signatureInfo = await getAlbumSignatureInfo(albumPath);
  
  if (!signatureInfo.hasSignature) {
    // 首次导出，按正常流程处理
    return await handleFirstExport(albumPath);
  }
  
  // 判断是否需要授权
  if (signatureInfo.originalAuthor?.requireAuthorization) {
    // 需要授权验证
    const currentUserQualCode = await getCurrentUserQualificationCode();
    
    // 检查当前用户是否有权限
    if (currentUserQualCode === signatureInfo.originalAuthor.qualificationCode) {
      // 原始作者，可以导出
      return await showSignatureSelectionDialog(albumPath);
    }
    
    // 检查是否在授权列表中
    const isAuthorized = signatureInfo.originalAuthor.authorizedList?.includes(currentUserQualCode);
    if (!isAuthorized) {
      // 未授权，提示导入授权文件
      return await showAuthorizationRequiredDialog();
    }
  }
  
  // 不需要授权或已授权，进入签名选择
  return await showSignatureSelectionDialog(albumPath);
}
```

---

## 需求3: 签名选择页面增强

### API端点
`POST /keytone_pkg/get_available_signatures`

### 请求示例
```typescript
interface AvailableSignature {
  encryptedId: string;
  qualificationCode: string;
  name: string;
  intro: string;
  isInAlbum: boolean;        // 是否已在专辑中（用于标记）
  isAuthorized: boolean;     // 是否有导出授权（用于使能/失能）
  isOriginalAuthor: boolean; // 是否为原始作者
}

async function getAvailableSignatures(albumPath: string): Promise<AvailableSignature[]> {
  const response = await fetch('http://localhost:8088/keytone_pkg/get_available_signatures', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ albumPath })
  });
  
  const data = await response.json();
  if (data.message !== 'ok') {
    throw new Error(data.message);
  }
  
  return data.signatures;
}
```

### UI实现示例
```vue
<template>
  <div class="signature-selection">
    <h2>选择签名</h2>
    <div v-for="sig in signatures" :key="sig.qualificationCode" class="signature-item">
      <!-- 签名卡片 -->
      <div 
        class="signature-card" 
        :class="{ 
          disabled: !sig.isAuthorized,
          'in-album': sig.isInAlbum 
        }"
        @click="selectSignature(sig)"
      >
        <img :src="sig.cardImagePath" alt="签名名片" />
        <div class="info">
          <h3>{{ sig.name }}</h3>
          <p>{{ sig.intro }}</p>
        </div>
        
        <!-- 标记 -->
        <div class="badges">
          <span v-if="sig.isOriginalAuthor" class="badge original">原始作者</span>
          <span v-if="sig.isInAlbum" class="badge in-album">已在专辑中</span>
          <span v-if="!sig.isAuthorized" class="badge unauthorized">需要授权</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const props = defineProps<{ albumPath: string }>();
const signatures = ref<AvailableSignature[]>([]);

onMounted(async () => {
  signatures.value = await getAvailableSignatures(props.albumPath);
});

function selectSignature(sig: AvailableSignature) {
  if (!sig.isAuthorized) {
    // 提示需要授权
    showMessage('此签名需要原始作者授权后才能导出');
    return;
  }
  
  // 继续导出流程
  emit('signature-selected', sig.encryptedId);
}
</script>

<style scoped>
.signature-card.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.signature-card.in-album {
  border: 2px solid #ffd700; /* 金色边框标记 */
}

.badge.original {
  background: #4CAF50;
  color: white;
}

.badge.in-album {
  background: #2196F3;
  color: white;
}

.badge.unauthorized {
  background: #f44336;
  color: white;
}
</style>
```

### 辅助API: 单独检查签名状态

如果只需要检查单个签名的状态，可以使用以下API：

#### 检查签名是否在专辑中
```typescript
async function checkSignatureInAlbum(albumPath: string, signatureId: string) {
  const response = await fetch('http://localhost:8088/keytone_pkg/check_signature_in_album', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ albumPath, signatureId })
  });
  
  const data = await response.json();
  return {
    isInAlbum: data.isInAlbum,
    qualificationCode: data.qualificationCode
  };
}
```

#### 检查签名授权状态
```typescript
async function checkSignatureAuthorization(albumPath: string, signatureId: string) {
  const response = await fetch('http://localhost:8088/keytone_pkg/check_signature_authorization', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ albumPath, signatureId })
  });
  
  const data = await response.json();
  return {
    isAuthorized: data.isAuthorized,
    requireAuthorization: data.requireAuthorization,
    qualificationCode: data.qualificationCode
  };
}
```

---

## 需求4: 签名作者信息展示

### UI组件示例
```vue
<template>
  <q-dialog v-model="showDialog">
    <q-card class="signature-authors-info">
      <q-card-section>
        <div class="text-h6">专辑签名信息</div>
      </q-card-section>

      <!-- 原始作者 -->
      <q-card-section v-if="signatureInfo?.originalAuthor">
        <div class="author-section">
          <h3>原始作者</h3>
          <div class="author-card">
            <img :src="signatureInfo.originalAuthor.cardImagePath" />
            <div>
              <h4>{{ signatureInfo.originalAuthor.name }}</h4>
              <p>{{ signatureInfo.originalAuthor.intro }}</p>
              <q-badge v-if="signatureInfo.originalAuthor.requireAuthorization" color="orange">
                需要授权导出
              </q-badge>
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- 直接导出作者 -->
      <q-card-section v-if="signatureInfo?.directExportAuthor">
        <div class="author-section">
          <h3>直接导出作者</h3>
          <div class="author-card">
            <img :src="signatureInfo.directExportAuthor.cardImagePath" />
            <div>
              <h4>{{ signatureInfo.directExportAuthor.name }}</h4>
              <p>{{ signatureInfo.directExportAuthor.intro }}</p>
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- 历史贡献作者 -->
      <q-card-section v-if="signatureInfo?.contributorAuthors.length">
        <div class="author-section">
          <h3>历史贡献作者</h3>
          <div 
            v-for="contributor in signatureInfo.contributorAuthors" 
            :key="contributor.qualificationCode"
            class="author-card"
          >
            <img :src="contributor.cardImagePath" />
            <div>
              <h4>{{ contributor.name }}</h4>
              <p>{{ contributor.intro }}</p>
            </div>
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="关闭" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  albumPath: string;
}

const props = defineProps<Props>();
const showDialog = ref(false);
const signatureInfo = ref<AlbumSignatureInfo | null>(null);

async function open() {
  signatureInfo.value = await getAlbumSignatureInfo(props.albumPath);
  showDialog.value = true;
}

defineExpose({ open });
</script>

<style scoped>
.author-section {
  margin-bottom: 20px;
}

.author-card {
  display: flex;
  gap: 16px;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  margin-bottom: 12px;
}

.author-card img {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
}
</style>
```

### 使用示例
```vue
<template>
  <q-page>
    <!-- 专辑详情页 -->
    <q-btn 
      label="查看签名信息" 
      @click="showSignatureInfo"
      icon="badge"
    />
    
    <!-- 签名信息对话框 -->
    <SignatureAuthorsDialog ref="signatureDialog" :album-path="albumPath" />
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import SignatureAuthorsDialog from './SignatureAuthorsDialog.vue';

const albumPath = ref('/path/to/album');
const signatureDialog = ref();

function showSignatureInfo() {
  signatureDialog.value?.open();
}
</script>
```

---

## 完整导出流程示例

```typescript
async function handleAlbumExport(albumPath: string) {
  // 步骤1: 用户选择是否需要签名
  const needSignature = await showSignatureChoiceDialog();
  
  if (!needSignature) {
    // 情况1: 无需签名，直接导出
    return await exportAlbumDirectly(albumPath);
  }
  
  // 步骤2: 获取专辑签名信息
  const signatureInfo = await getAlbumSignatureInfo(albumPath);
  
  // 步骤3: 判断是否需要授权验证
  if (signatureInfo.hasSignature && 
      signatureInfo.originalAuthor?.requireAuthorization) {
    // 情况2: 需要授权验证
    const currentUserQualCode = await getCurrentUserQualificationCode();
    const isAuthorized = 
      currentUserQualCode === signatureInfo.originalAuthor.qualificationCode ||
      signatureInfo.originalAuthor.authorizedList?.includes(currentUserQualCode);
    
    if (!isAuthorized) {
      // 未授权，提示导入授权文件
      return await showAuthorizationRequiredDialog();
    }
  }
  
  // 步骤4: 显示签名选择页面
  const signatures = await getAvailableSignatures(albumPath);
  const selectedSignatureId = await showSignatureSelectionDialog(signatures);
  
  // 步骤5: 应用签名并导出
  const { requireAuthorization, contactEmail, contactAdditional } = 
    await getAuthorizationSettings();
  
  const result = await applySignatureConfig(
    albumPath,
    selectedSignatureId,
    requireAuthorization,
    contactEmail,
    contactAdditional
  );
  
  // 步骤6: 执行导出
  return await exportAlbumWithSignature(albumPath);
}
```

---

## 错误处理

### 常见错误及处理方式

1. **签名不存在**
```typescript
try {
  await applySignatureConfig(albumPath, signatureId, ...);
} catch (error) {
  if (error.message.includes('签名不存在')) {
    showMessage('签名已被删除，请选择其他签名');
  }
}
```

2. **专辑配置损坏**
```typescript
try {
  const info = await getAlbumSignatureInfo(albumPath);
} catch (error) {
  if (error.message.includes('解密失败')) {
    showMessage('专辑配置文件损坏，无法读取签名信息');
  }
}
```

3. **授权验证失败**
```typescript
const authStatus = await checkSignatureAuthorization(albumPath, signatureId);
if (!authStatus.isAuthorized && authStatus.requireAuthorization) {
  showMessage('此专辑需要原始作者授权后才能导出\n请联系原始作者获取授权文件');
}
```

---

## 调试技巧

### 查看未加密的签名内容
SDK在应用签名时会在终端输出调试日志：
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

### 验证资格码生成
资格码是签名ID的SHA256哈希，可以在前端验证：
```typescript
import { sha256 } from 'crypto-js';

function generateQualificationCode(originalSignatureID: string): string {
  return sha256(originalSignatureID).toString();
}
```

---

## 总结

通过以上5个API端点，前端可以完整实现：

1. ✅ **删除"无需签名+需要授权"分支** - 简化导出流程
2. ✅ **再次导出时的签名识别** - `get_album_signature_info`
3. ✅ **签名选择页面增强** - `get_available_signatures`
4. ✅ **签名作者信息展示** - `get_album_signature_info`

所有API都已实现并通过编译验证，可以开始前端集成开发！
