# 前端集成使用示例

本文档展示如何在现有的导出流程中集成新实现的签名功能（需求1.2.3）。

## 快速开始

### 1. 在专辑详情页添加"查看签名信息"按钮

```vue
<template>
  <q-page class="album-detail-page">
    <!-- 专辑信息 -->
    <div class="album-info">
      <h2>{{ albumName }}</h2>
      <q-btn
        label="查看签名信息"
        icon="badge"
        color="primary"
        @click="showSignatureInfo"
      />
    </div>

    <!-- 签名信息对话框 -->
    <SignatureAuthorsDialog ref="signatureDialog" :album-path="albumPath" />
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import SignatureAuthorsDialog from 'src/components/export-flow/SignatureAuthorsDialog.vue';

const albumPath = ref('/path/to/album');
const signatureDialog = ref();

function showSignatureInfo() {
  signatureDialog.value?.open();
}
</script>
```

### 2. 在导出流程中集成签名选择

#### 步骤1：检查专辑签名状态（需求2）

```typescript
import { GetAlbumSignatureInfo } from 'src/boot/query/keytonePkg-query';

async function startExportFlow(albumPath: string) {
  // 1. 用户选择是否需要签名
  const needSignature = await showNeedSignatureDialog();
  
  if (!needSignature) {
    // 需求1：无需签名，直接导出
    return await exportAlbumDirectly(albumPath);
  }

  // 2. 获取专辑签名信息（需求2）
  const signatureInfo = await GetAlbumSignatureInfo(albumPath);
  
  // 3. 判断是否需要授权验证
  if (signatureInfo.hasSignature && 
      signatureInfo.originalAuthor?.requireAuthorization) {
    // 需要授权，检查当前用户是否有权限
    const hasAuth = await checkCurrentUserAuthorization(signatureInfo);
    
    if (!hasAuth) {
      // 未授权，提示导入授权文件
      showAuthorizationRequiredDialog();
      return;
    }
  }

  // 4. 进入签名选择
  const selectedSignatureId = await showSignatureSelection(albumPath);
  
  // 5. 应用签名并导出
  await applySignatureAndExport(albumPath, selectedSignatureId);
}
```

#### 步骤2：显示签名选择对话框（需求3）

```vue
<template>
  <div>
    <q-btn label="选择签名" @click="selectSignature" />
    
    <!-- 签名选择对话框 -->
    <SignatureSelectionDialog
      ref="signatureSelectionDialog"
      :album-path="albumPath"
      @confirm="onSignatureSelected"
      @cancel="onSelectionCancelled"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import SignatureSelectionDialog from 'src/components/export-flow/SignatureSelectionDialog.vue';

const albumPath = ref('/path/to/album');
const signatureSelectionDialog = ref();

async function selectSignature() {
  signatureSelectionDialog.value?.open();
}

function onSignatureSelected(signatureId: string) {
  console.log('用户选择的签名ID:', signatureId);
  // 继续导出流程
  proceedWithExport(signatureId);
}

function onSelectionCancelled() {
  console.log('用户取消签名选择');
}
</script>
```

### 3. 完整的导出流程实现

```typescript
// 文件：src/composables/useAlbumExport.ts

import { ref } from 'vue';
import { Dialog, Notify } from 'quasar';
import {
  GetAlbumSignatureInfo,
  ApplySignatureConfig,
  ExportAlbum,
} from 'src/boot/query/keytonePkg-query';
import type { AlbumSignatureInfo } from 'src/types/export-flow';

export function useAlbumExport(albumPath: string) {
  const isExporting = ref(false);

  /**
   * 开始导出流程
   */
  async function startExport() {
    if (isExporting.value) return;
    isExporting.value = true;

    try {
      // 步骤1：询问是否需要签名
      const needSignature = await askNeedSignature();
      
      if (!needSignature) {
        // 需求1：无需签名，直接导出
        await exportAlbumDirectly();
        return;
      }

      // 步骤2：检查签名状态（需求2）
      const signatureInfo = await GetAlbumSignatureInfo(albumPath);
      
      // 步骤3：验证授权
      const canExport = await validateAuthorization(signatureInfo);
      if (!canExport) {
        return;
      }

      // 步骤4：选择签名（需求3）
      const signatureId = await selectSignatureWithDialog();
      if (!signatureId) {
        return; // 用户取消
      }

      // 步骤5：收集授权信息
      const authInfo = await getAuthorizationInfo(signatureInfo);

      // 步骤6：应用签名配置
      const success = await ApplySignatureConfig({
        albumPath,
        needSignature: true,
        signatureId,
        ...authInfo,
      });

      if (!success) {
        throw new Error('应用签名配置失败');
      }

      // 步骤7：执行导出
      await exportAlbumDirectly();

    } catch (error: any) {
      console.error('导出失败:', error);
      Notify.create({
        type: 'negative',
        message: '导出失败: ' + (error.message || '未知错误'),
      });
    } finally {
      isExporting.value = false;
    }
  }

  /**
   * 询问是否需要签名
   */
  async function askNeedSignature(): Promise<boolean> {
    return new Promise((resolve) => {
      Dialog.create({
        title: '导出设置',
        message: '是否为此专辑添加签名？',
        options: {
          type: 'radio',
          model: 'no',
          items: [
            { label: '无需签名', value: 'no' },
            { label: '需要签名', value: 'yes' },
          ],
        },
        cancel: true,
      }).onOk((value) => {
        resolve(value === 'yes');
      }).onCancel(() => {
        resolve(false);
      });
    });
  }

  /**
   * 验证授权（需求2）
   */
  async function validateAuthorization(
    signatureInfo: AlbumSignatureInfo
  ): Promise<boolean> {
    // 首次导出，无需验证
    if (!signatureInfo.hasSignature) {
      return true;
    }

    // 不需要授权，直接通过
    if (!signatureInfo.originalAuthor?.requireAuthorization) {
      return true;
    }

    // 获取当前用户的签名资格码
    const currentUserQualCode = await getCurrentUserQualificationCode();

    // 是原始作者本人
    if (currentUserQualCode === signatureInfo.originalAuthor.qualificationCode) {
      return true;
    }

    // 检查是否在授权列表中
    const isAuthorized = signatureInfo.originalAuthor.authorizedList?.includes(
      currentUserQualCode
    );

    if (!isAuthorized) {
      // 未授权，提示用户
      Notify.create({
        type: 'warning',
        message: '此专辑需要原始作者授权后才能导出',
        caption: '请联系原始作者获取授权文件，或导入已有的授权文件',
        actions: [
          { label: '导入授权文件', color: 'white', handler: () => importAuthFile() },
          { label: '取消', color: 'white' },
        ],
        timeout: 0,
      });
      return false;
    }

    return true;
  }

  /**
   * 使用对话框选择签名（需求3）
   */
  function selectSignatureWithDialog(): Promise<string | null> {
    return new Promise((resolve) => {
      // 这里需要集成 SignatureSelectionDialog 组件
      // 实际使用时需要通过事件或回调获取结果
      resolve('selected-signature-id');
    });
  }

  /**
   * 获取授权信息
   */
  async function getAuthorizationInfo(signatureInfo: AlbumSignatureInfo) {
    // 如果是首次导出，询问是否需要授权
    if (!signatureInfo.hasSignature) {
      return await Dialog.create({
        title: '授权设置',
        message: '是否需要对二次创作进行授权控制？',
        options: {
          type: 'radio',
          model: 'no',
          items: [
            { label: '不需要授权', value: 'no' },
            { label: '需要授权', value: 'yes' },
          ],
        },
        prompt: {
          model: '',
          type: 'text',
          label: '联系邮箱',
          hint: '需要授权时必填',
        },
        cancel: true,
      }).onOk((data) => ({
        requireAuthorization: data.value === 'yes',
        contactEmail: data.value === 'yes' ? data.prompt : undefined,
      }));
    }

    return {
      requireAuthorization: signatureInfo.originalAuthor?.requireAuthorization || false,
      contactEmail: signatureInfo.originalAuthor?.contactEmail,
    };
  }

  /**
   * 直接导出专辑
   */
  async function exportAlbumDirectly() {
    Notify.create({
      type: 'info',
      message: '正在导出专辑...',
    });

    const blob = await ExportAlbum(albumPath);
    
    // 下载文件
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `album-${Date.now()}.ktpkg`;
    a.click();
    URL.revokeObjectURL(url);

    Notify.create({
      type: 'positive',
      message: '专辑导出成功！',
    });
  }

  /**
   * 获取当前用户的资格码（示例）
   */
  async function getCurrentUserQualificationCode(): Promise<string> {
    // 这里需要实际实现获取当前用户签名的逻辑
    // 可以从用户配置中读取默认签名，然后生成资格码
    return 'current-user-qualification-code';
  }

  /**
   * 导入授权文件（示例）
   */
  function importAuthFile() {
    // 这里需要实现授权文件导入逻辑
    console.log('导入授权文件');
  }

  return {
    isExporting,
    startExport,
  };
}
```

### 4. 在页面中使用

```vue
<template>
  <q-page class="album-page">
    <div class="album-header">
      <h2>{{ albumName }}</h2>
      <div class="actions">
        <!-- 查看签名信息（需求4） -->
        <q-btn
          outline
          label="查看签名信息"
          icon="badge"
          color="primary"
          @click="showSignatureInfo"
        />
        
        <!-- 导出按钮 -->
        <q-btn
          unelevated
          label="导出专辑"
          icon="file_download"
          color="primary"
          :loading="isExporting"
          @click="startExport"
        />
      </div>
    </div>

    <!-- 对话框组件 -->
    <SignatureAuthorsDialog ref="signatureDialog" :album-path="albumPath" />
    <SignatureSelectionDialog
      ref="signatureSelectionDialog"
      :album-path="albumPath"
      @confirm="onSignatureSelected"
      @cancel="onSelectionCancelled"
    />
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAlbumExport } from 'src/composables/useAlbumExport';
import SignatureAuthorsDialog from 'src/components/export-flow/SignatureAuthorsDialog.vue';
import SignatureSelectionDialog from 'src/components/export-flow/SignatureSelectionDialog.vue';

const props = defineProps<{
  albumPath: string;
  albumName: string;
}>();

const signatureDialog = ref();
const signatureSelectionDialog = ref();

const { isExporting, startExport } = useAlbumExport(props.albumPath);

function showSignatureInfo() {
  signatureDialog.value?.open();
}

function onSignatureSelected(signatureId: string) {
  console.log('选择的签名:', signatureId);
  // 在实际的 useAlbumExport 实现中处理
}

function onSelectionCancelled() {
  console.log('取消选择签名');
}
</script>

<style scoped lang="scss">
.album-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
}

.actions {
  display: flex;
  gap: 12px;
}
</style>
```

## API参考

### GetAlbumSignatureInfo

获取专辑签名信息。

```typescript
async function GetAlbumSignatureInfo(
  albumPath: string
): Promise<AlbumSignatureInfo>
```

**返回值**:
```typescript
{
  hasSignature: boolean;
  originalAuthor?: SignatureAuthorInfo;
  contributorAuthors: SignatureAuthorInfo[];
  directExportAuthor?: SignatureAuthorInfo;
  allSignatures: Record<string, AlbumSignatureEntry>;
}
```

### GetAvailableSignatures

获取可用于导出的签名列表。

```typescript
async function GetAvailableSignatures(
  albumPath: string
): Promise<AvailableSignature[]>
```

**返回值**:
```typescript
[
  {
    encryptedId: string;
    qualificationCode: string;
    name: string;
    intro: string;
    isInAlbum: boolean;        // 是否已在专辑中
    isAuthorized: boolean;     // 是否有导出授权
    isOriginalAuthor: boolean; // 是否为原始作者
  },
  // ...
]
```

### CheckSignatureInAlbum

检查签名是否在专辑中。

```typescript
async function CheckSignatureInAlbum(
  albumPath: string,
  signatureId: string
): Promise<{ isInAlbum: boolean; qualificationCode: string }>
```

### CheckSignatureAuthorization

检查签名授权状态。

```typescript
async function CheckSignatureAuthorization(
  albumPath: string,
  signatureId: string
): Promise<{
  isAuthorized: boolean;
  requireAuthorization: boolean;
  qualificationCode: string;
}>
```

## 常见问题

### Q: 如何判断专辑是否已有签名？

A: 调用 `GetAlbumSignatureInfo` 并检查返回的 `hasSignature` 字段。

```typescript
const signatureInfo = await GetAlbumSignatureInfo(albumPath);
if (!signatureInfo.hasSignature) {
  console.log('这是首次导出');
}
```

### Q: 如何显示未授权的提示？

A: 检查 `isAuthorized` 字段，并显示相应的UI。

```typescript
const signatures = await GetAvailableSignatures(albumPath);
const unauthorizedSigs = signatures.filter(sig => !sig.isAuthorized);

if (unauthorizedSigs.length > 0) {
  showUnauthorizedWarning();
}
```

### Q: 如何标记已在专辑中的签名？

A: 使用 `isInAlbum` 字段添加视觉标记。

```vue
<q-card :class="{ 'in-album': signature.isInAlbum }">
  <!-- 签名内容 -->
</q-card>

<style>
.in-album {
  border-left: 4px solid #2196F3;
}
</style>
```

## 下一步

- [ ] 将组件集成到实际的导出流程中
- [ ] 实现授权文件导入功能
- [ ] 添加签名验证和错误处理
- [ ] 完善用户体验（加载状态、错误提示等）
