# 技术设计：签名管理系统

## 架构概览

签名管理系统采用客户端-服务端分离架构，前端负责UI交互和数据展示，后端负责数据持久化和加密处理。

```text
┌─────────────────────────────────────────────────────────┐
│                    前端 (Vue 3 + Quasar)                  │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ 签名管理页面  │  │ 签名对话框    │  │ 专辑导出页面  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                          │                               │
│  ┌──────────────────────▼────────────────────────────┐ │
│  │          签名服务层 (signature-service.ts)         │ │
│  │  - 签名CRUD  - 加密/解密  - 文件处理  - 数据验证  │ │
│  └───────────────────────┬───────────────────────────┘ │
└────────────────────────┼─────────────────────────────┘
                         │ HTTP/SSE
┌────────────────────────▼─────────────────────────────┐
│                  后端 (Go + Gin)                       │
├─────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  /signature/*│  │  /signature/*│  │    /stream   │ │
│  │  (CRUD/导入导出/解密/图片/排序)   │  │    (SSE)     │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                          │                               │
│  ┌──────────────────────▼────────────────────────────┐ │
│  │         配置管理层 (config/config.go)              │ │
│  │  - GetValue  - SetValue  - 加密处理  - 文件IO    │ │
│  └───────────────────────┬───────────────────────────┘ │
└────────────────────────┼─────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────┐
│              配置文件 (KeyToneSetting.json)            │
│  {                                                     │
│    "signature": {                                     │
│      "<加密ID>": { "value": "<加密签名JSON>", "sort": { "time": 1730000000 } }  │
│    }                                                   │
│  }                                                     │
└─────────────────────────────────────────────────────┘
```

## 数据模型

### 前端数据结构

```typescript
/**
 * 签名接口（前端使用）
 */
interface Signature {
  id: string;              // 唯一标识（nanoid，21字符）
  name: string;            // 签名名称（必填，1-50字符）
  intro: string;           // 个人介绍（选填，0-500字符）
  cardImage: string;       // 名片图片路径（选填）
}

/**
 * 签名管理器（前端使用）
 */
interface SignatureManager {
  [key: string]: Signature;  // key = 签名ID
}

/**
 * .ktsign 文件格式
 */
interface SignatureFile {
  version: string;           // 文件格式版本（"1.0.0"）
  signature: {
    id: string;
    name: string;
    intro: string;
    cardImage: string;       // Base64 编码的图片数据
  };
  checksum: string;          // SHA-256 校验和
}
```

### 后端存储结构（实际实现）

```json
{
  "signature": {
    "<encryptedId>": {
      "value": "<hex-encoded AES-GCM ciphertext>",
      "sort": { "time": 1730000000 }
    }
  }
}
```

加密说明：

- **Key**：使用 AES-256 加密签名ID
- **Value**：使用 AES-256 加密整个签名对象的JSON字符串
- **加密位置**：所有加密/解密逻辑在 Go SDK 中实现
- **前端交互**：前端通过 HTTP 接口调用后端加密/解密服务，不在前端实现加密逻辑

## 加密方案

### 后端加密（Go）

所有加密/解密逻辑在 `sdk/signature/encryption.go` 中实现（实际方案摘要）：

- 使用 AES-256-GCM 模式
- KeyA（32字节）：用于加密签名ID（配置键）；同时作为 PBKDF2 的 password
- 动态密钥：解密出原始ID后取其后7位作为 salt，与 KeyA 通过 PBKDF2(SHA-256, 10000, 32字节) 派生，用于加/解密 Value
- KeyB（32字节）：用于导出/导入文件的外层整体加/解密
- 加密输出采用十六进制编码（hex）进行持久化

## API 设计

### 签名管理 API

签名管理需要新增专用端点，因为涉及加密/解密逻辑和保护码生成。

#### 创建签名

```http
POST /signature/create
Content-Type: multipart/form-data

字段：
- id: string（未加密ID）
- name: string（必填）
- intro: string（可选）
- cardImage: file（可选）
```

响应：

```json
{ "success": true, "data": { "id": "<encryptedId>" } }
```

**处理流程**：
1. 后端接收表单字段及图片文件
2. 将图片写入 ConfigPath/signature 目录（文件名为基于 id|name|originalName|timestamp 的 SHA-1）
3. 使用 KeyA 加密 ID 作为配置键；使用动态密钥加密签名 JSON 作为 value
4. 存储到配置文件 `signature`
5. **配置文件更新后，现有 SSE 机制自动推送全量配置**

#### 获取所有签名

```http
GET /signature/list
```

响应：

```json
{
  "success": true,
  "data": {
    "<encryptedId>": {
      "value": "<hex-ciphertext>",
      "sort": { "time": 1730000000 }
    }
  }
}
```

如需明文数据，逐项调用：

```http
POST /signature/decrypt
{ "encryptedValue": "<hex-ciphertext>", "encryptedId": "<encryptedId>" }
```

→ 返回：`{ "success": true, "data": "{\"name\":...,\"intro\":...,\"cardImage\":...}" }`

#### 更新签名

```http
POST /signature/update
Content-Type: multipart/form-data

字段：
- encryptedId: string（必填）
- name: string（可选）
- intro: string（可选）
- cardImage: file（可选）
- removeImage: "true"|"false"（可选）
- imageChanged: "true"|"false"（可选；未提供视为 true）
```

响应： `{ "success": true }`

#### 删除签名

```http
POST /signature/delete
{ "id": "<encryptedId>" }
```

响应：`{ "success": true }`

#### 导出签名

```http
POST /signature/export
{ "encryptedId": "<encryptedId>" }
```

响应： `.ktsign` 二进制流（Content-Type: application/octet-stream）。

**处理流程**：
1. 根据加密ID定位并解密签名数据
2. 读取图片文件并转为十六进制字符串，连同 Name/Intro 等构建内部 JSON（包含 CardImageName）
3. 使用 KeyB 对内部 JSON 进行加密，返回加密字符串的字节作为下载内容
4. 前端将响应保存为 `.ktsign`

#### 导入签名

```http
POST /signature/import
Content-Type: multipart/form-data
- file: .ktsign 文件
```

响应：

- 成功：`{ "success": true, "data": { "encryptedId": "...", "name": "..." } }`
- 冲突：HTTP 409 `{ "success": false, "conflict": true, "data": { "encryptedId": "...", "name": "..." } }`

确认覆盖：

```http
POST /signature/import-confirm
{ "file": "<加密字符串>", "overwrite": true }
```

#### 获取图片

```http
POST /signature/get-image
{ "imagePath": "<absolute path from config>" }
```

响应：图片二进制（application/octet-stream）。前端将二进制转 Blob URL 进行展示。

### 复用现有 API

#### SSE 数据同步机制（重要澄清）

**现有机制说明**：

KeyTone 项目中已经实现了完整的 SSE 数据同步机制：

- 后端通过 `/stream` 端点持续推送**全量配置数据**到前端
- 当配置文件发生任何变化时，后端会通过 SSE 直接发送**整个配置文件的完整数据**
- 前端已有监听逻辑，在接收到 SSE 消息后会处理全量配置数据

**签名管理的适配方式**：

❌ **不需要**：
- 不需要重新设计 SSE 通知机制
- 不需要添加新的 `signature_updated` 标志
- 不需要在后端单独发送签名更新通知
- 不需要前端主动调用 `GET /signature/list` 刷新数据

✅ **只需要**：
- 在前端现有的 SSE 全量数据处理逻辑中，添加对 `signature` 字段的解构和处理
- 当 SSE 推送全量配置数据时，提取其中的 `signature` 数据块
- 更新前端签名状态存储（如 Pinia store 或组件状态）
- 触发 Vue 响应式更新，自动刷新相关 UI 组件

**前端实现示例**：

```typescript
// 在现有的 SSE 监听逻辑中（假设在 app-store.ts 或类似文件中）
eventSource.addEventListener('message', (event) => {
  const fullConfig = JSON.parse(event.data);  // 后端直接推送的全量配置数据

  // 现有逻辑：处理其他配置字段
  if (fullConfig.key_package) {
    // ... 处理键盘配置
  }
  
  if (fullConfig.audio_settings) {
    // ... 处理音频配置
  }

  // 新增逻辑：处理签名管理数据
  if (fullConfig.signature) {
    // 解构获取签名管理数据
  const signatureData = fullConfig.signature;
    
    // 更新前端状态（根据实际使用的状态管理方案）
    // 方案 1：使用 Pinia store
    signatureStore.updateSignatures(signatureData);
    
    // 方案 2：直接更新响应式变量
    // signatures.value = signatureData;
    
    // 方案 3：发出自定义事件
    // eventBus.emit('signatures-updated', signatureData);
  }
});
```

**数据流示意**：

```text
后端配置变更
    ↓
后端通过 /stream 推送全量配置（SSE）
    ↓
前端 SSE 监听器接收全量数据
    ↓
解构提取 signature 字段
    ↓
更新前端签名状态
    ↓
Vue 响应式自动刷新 UI
```

**关键要点**：

- ✅ 现有代码已有 SSE 全量数据推送，无需重新设计
- ✅ 签名数据作为配置文件的一部分，会随全量数据一起推送
- ✅ 前端只需在现有的全量数据处理逻辑中添加签名数据的解构和适配
- ✅ 不需要单独的再次调用 `GET /signature/list` 刷新数据
- ✅ 所有签名 CRUD 操作完成后，后端保存配置文件，自动触发 SSE 推送

## 文件处理

### 名片图片存储

1. **目录结构**（实际）：

   ```text
   ConfigPath/
   ├── KeyToneSetting.json
   └── signature/
     ├── <sha1seed>.png
     └── <sha1seed>.jpg
   ```

   **注意**：不再需要 `exported/` 目录，导出文件由用户选择保存位置。

2. **图片命名**：使用基于 `id|name|originalName|timestamp` 的字符串计算 SHA-1 作为文件名（由后端处理）

3. **图片处理流程（分阶段详解）**：

   #### 阶段 1：创建/修改阶段 - 图片选择与上传
   
   在这个阶段,图片数据以**文件对象(File)**或**Base64字符串**的形式存在于前端内存中:
   
   - **前端选择器**：用户通过 `<q-file>` 组件选择图片文件
   - **文件对象**：选择后,图片以 `File` 对象存在于组件的 `v-model` 中
   - **预览渲染**：通过 `FileReader.readAsDataURL()` 将 `File` 对象转为 Base64 用于预览显示
   - **表单数据**：在表单提交前,图片仍然是 `File` 对象,**不是路径字符串**
   
   ```typescript
   // 前端表单数据结构（创建/编辑阶段）
  interface SignatureFormData {
     name: string;
     intro: string;
     cardImage: File | null;  // 注意：这里是 File 对象,不是字符串路径
   }
   ```
   
  #### 阶段 2：提交存储阶段 - multipart 传输与后端持久化（实际）
   
   当用户点击"创建"或"更新"按钮时:
   
  - **HTTP 传输**：以 multipart/form-data 直接上传 `cardImage` 文件
  - **后端处理**：
    1. 接收文件与表单字段
    2. 基于 `id|name|originalName|timestamp` 生成 SHA-1 文件名
    3. 保存图片文件到 `ConfigPath/signature/` 目录
    4. 在配置文件 `signature` 中存储图片绝对路径
   
  （改为 multipart 表单，不再传 Base64 图片字符串）
   
  ```json
  // 后端配置文件存储（示例）
  {
    "signature": {
      "<encrypted_id>": {
        "value": "<hex-ciphertext>",
        "sort": { "time": 1730000000 }
      }
    }
  }
  ```
   
#### 阶段 3：列表渲染阶段 - 路径字符串到图片资源
   
   当前端需要显示签名列表时:
   
- **获取路径**：通过 API 获取签名数据，其中 `cardImage` 为图片绝对路径（由后端存储）
- **渲染**：通过 `POST /signature/get-image` 读取二进制并转为 Blob URL 渲染
   
  ```typescript
   // 前端签名数据结构（列表显示阶段）
  interface Signature {
    name: string;
    intro: string;
    cardImage: string;  // 绝对路径
  }
  ```
   
  ```vue
  <!-- 在模板中使用（示意：通过 Blob URL 显示） -->
  <q-img :src="imageBlobUrl" />
  ```
   
  #### 阶段 4：导出/导入阶段 - 文件封装与加解密（实际）
   
  **导出流程**：
  - 前端调用 `POST /signature/export`（payload: { encryptedId }）
  - 后端读取签名数据，读取图片并转为十六进制字符串，构建内部 JSON（包含 CardImageName）
  - 用 KeyB 对内部 JSON 加密，返回加密字符串的字节作为下载内容（octet-stream）
  - 前端保存为 `.ktsign` 文件
   
  **导入流程**：
  - 前端以 multipart 方式上传 `.ktsign` 至 `POST /signature/import`
  - 后端用 KeyB 解密并解析内部 JSON；若冲突返回 409 与 conflict 标志
  - 覆盖导入：前端调用 `POST /signature/import-confirm`（{ file, overwrite: true }）
   
#### 图片数据形态总结
   
   | 阶段       | 图片数据形态                        | 说明                                  |
   | ---------- | ----------------------------------- | ------------------------------------- |
   | 选择器选择 | `File` 对象                         | 前端内存中的文件对象                  |
   | 表单编辑   | `File` 对象 / 预览用 Base64         | 未提交前仍是 File 对象,不是路径字符串 |
   | HTTP 传输  | multipart 文件                      | 直接上传文件至后端                    |
   | 后端存储   | 二进制文件 + 绝对路径（配置文件中） | 图片存为文件，配置中存绝对路径        |
   | 列表显示   | 绝对路径 → get-image → Blob URL     | 通过后端接口读取二进制再渲染          |
   | 导出文件   | KeyB 加密的内部 JSON（二进制流）    | 内部 JSON 的图片为十六进制字符串      |
   | 导入文件   | 二进制流 → KeyB 解密 → 文件写入     | 恢复为文件系统中的图片                |
   
**关键要点**：
- ✅ 创建/编辑表单中，图片是 `File` 对象，非路径字符串
- ✅ HTTP 传输使用 multipart 文件上传
- ✅ 配置文件存储图片绝对路径，图片本体存于 ConfigPath/signature
- ✅ 列表渲染通过后端 get-image 接口读取二进制再显示
- ✅ 导出/导入采用 KeyB 对内部 JSON 进行整体加/解密

### .ktsign 文件格式（实际）

`.ktsign` 为二进制文件，内容为使用 KeyB 加密后的内部 JSON 字符串。内部 JSON 示例如下：

```json
{
  "key": "<encryptedId>",
  "name": "张三",
  "intro": "资深音效设计师",
  "cardImage": "<hex-encoded-bytes>",
  "cardImageName": "avatar.png"
}
```

导出：前端调用 `POST /signature/export`，将返回的二进制流保存为 `.ktsign`。

导入：前端以 multipart 上传 `.ktsign` 到 `POST /signature/import`；如遇冲突，按接口返回进行覆盖确认。

导入流程（实际）：

1. 前端以 multipart 上传 `.ktsign` 到 `POST /signature/import`
2. 后端用 KeyB 解密并解析内部 JSON（包含十六进制图片数据与文件名）
3. 若目标加密ID已存在，返回 409 冲突与 `conflict: true`
4. 用户确认覆盖后，调用 `POST /signature/import-confirm` 完成导入
5. 后端写入图片文件至 ConfigPath/signature，并更新配置 `signature`
6. 保存后由 SSE 自动推送全量配置

## 前端服务层设计

### signature-service.ts

前端服务层简化为 HTTP 调用，不包含加密逻辑：

```typescript
import { api } from 'boot/axios';

export class SignatureService {
  /**
   * 获取所有签名
   */
  async getAllSignatures(): Promise<SignatureManager> {
    const response = await api.get('/signature/list');
    
    if (response.data.message === 'ok') {
      return response.data.signatures || {};
    }
    
    return {};
  }

  /**
   * 创建签名
   */
  async createSignature(data: {
    name: string;
    intro?: string;
    cardImage?: File;
  }): Promise<Signature> {
    let cardImageBase64 = '';
    
    // 将图片转为 Base64
    if (data.cardImage) {
      cardImageBase64 = await this.fileToBase64(data.cardImage);
    }

    const response = await api.post('/signature/create', {
      name: data.name,
      intro: data.intro || '',
      cardImage: cardImageBase64,
    });

    if (response.data.message === 'ok') {
      return response.data.signature;
    }

    throw new Error('Failed to create signature');
  }

  /**
   * 更新签名
   */
  async updateSignature(data: {
    id: string;
    name: string;
    intro?: string;
    cardImage?: File;
  }): Promise<Signature> {
    let cardImageBase64 = '';
    
    // 如果有新图片，转为 Base64
    if (data.cardImage) {
      cardImageBase64 = await this.fileToBase64(data.cardImage);
    }

    const response = await api.put('/signature/update', {
      id: data.id,
      name: data.name,
      intro: data.intro || '',
      cardImage: cardImageBase64,
    });

    if (response.data.message === 'ok') {
      return response.data.signature;
    }

    throw new Error('Failed to update signature');
  }

  /**
   * 删除签名
   */
  async deleteSignature(id: string): Promise<void> {
    const response = await api.delete(`/signature/delete/${id}`);
    
    if (response.data.message !== 'ok') {
      throw new Error('Failed to delete signature');
    }
  }

  /**
   * 导出签名到文件
   */
  async exportSignature(id: string): Promise<void> {
    // 调用后端 API 获取导出数据（包含 Base64 图片）
    const response = await api.get(`/signature/export/${id}`);
    
    if (response.data.version && response.data.signature) {
      const jsonStr = JSON.stringify(response.data, null, 2);
      const blob = new Blob([jsonStr], { type: 'application/json' });
      
      // 使用 Web API 保存文件
      try {
        const handle = await window.showSaveFilePicker({
          suggestedName: `${response.data.signature.name}.ktsign`,
          types: [{
            description: 'KeyTone Signature File',
            accept: { 'application/json': ['.ktsign'] },
          }],
        });
        
        const writable = await handle.createWritable();
        await writable.write(blob);
        await writable.close();
      } catch (err) {
        // 用户取消了保存操作，或浏览器不支持 File System Access API
        if (err.name === 'AbortError') {
          return;
        }
        
        // 降级方案：使用 <a download>
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${response.data.signature.name}.ktsign`;
        a.click();
        URL.revokeObjectURL(url);
      }
    } else {
      throw new Error('Invalid export data');
    }
  }

  /**
   * 从文件导入签名
   */
  async importSignature(file: File, overwrite: boolean = false): Promise<{ 
    success: boolean; 
    exists?: boolean; 
    signature?: Signature 
  }> {
    const text = await file.text();
    const signatureFile = JSON.parse(text);

    // 调用后端 API 导入
    const response = await api.post('/signature/import', {
      ...signatureFile,
      overwrite,
    });

    if (response.data.message === 'ok') {
      if (response.data.exists && !overwrite) {
        // 签名已存在，需要用户确认
        return { success: false, exists: true };
      }
      
      return { success: true, signature: response.data.signature };
    }

    throw new Error('Failed to import signature');
  }

  /**
   * 获取图片 URL（用于显示）
   */
  getImageUrl(cardImagePath: string): string {
    if (!cardImagePath) return '';
    
    const filename = cardImagePath.split('/').pop();
    return `${api.defaults.baseURL}/signature/image/${filename}`;
  }

  /**
   * 将 File 转为 Base64
   */
  private async fileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      
      reader.onload = () => {
        const result = reader.result as string;
        resolve(result); // 格式：data:image/png;base64,...
      };
      
      reader.onerror = () => {
        reject(new Error('Failed to read file'));
      };
      
      reader.readAsDataURL(file);
    });
  }
}

export const signatureService = new SignatureService();
```

## UI 组件设计

### 签名列表项布局设计

**列表项结构（新增设计）**：

列表项采用**左图右文**的布局设计：

```text
┌─────────────────────────────────────────────────────────┐
│  ┌────────┐                                    ┌─────┐  │
│  │        │  签名名称（粗体大字）     [编辑] [删除] │
│  │ 图片   │  个人介绍（灰色小字）...             │  └─ [导出]
│  │ 缩略图 │  最多2行，超出省略号                │  │
│  │(80x80) │                                    └──────┤
│  │        │                                           │
│  └────────┘                                           │
└─────────────────────────────────────────────────────────┘
```

**布局特点**：

1. **图片区域**（左侧）：
   - 尺寸：80x80px（小尺寸缩略图）
   - 位置：列表项左侧
   - **功能**：仅用于展示，点击时打开预览对话框查看大图
   - **样式**：圆角 8px，可选阴影

2. **信息区域**（中间主要内容）：
   - **签名名称**：大字体（如 16px），粗体
   - **个人介绍**：小字体（如 12px），灰色，最多 2 行，超出显示省略号
   - **高度**：约 80px，与图片高度对齐

3. **操作区域**（右侧）：
   - 按钮：编辑、删除、导出
   - 使用图标按钮节省空间

4. **交互**：
   - 点击**图片缩略图** → 打开图片预览对话框
   - 点击**签名名称或介绍区域** → 打开编辑对话框
   - **操作按钮** → 相应操作（编辑、删除、导出）

**实现示例**：

```vue
<template>
  <q-card class="signature-list-item" @click="handleInfoClick">
    <q-card-section horizontal class="q-pa-none">
      <!-- 左侧图片 -->
      <div class="signature-image-container">
        <q-img
          v-if="signature.cardImage"
          :src="getImageUrl(signature.cardImage)"
          class="signature-thumbnail"
          @click.stop="showImagePreview = true"
        />
        <div v-else class="signature-image-placeholder">
          <q-icon name="image" size="lg" />
        </div>
      </div>

      <!-- 中间信息 -->
      <q-card-section class="signature-info-container">
        <div class="signature-name text-h6 text-weight-bold">
          {{ signature.name }}
        </div>
        <div class="signature-intro text-caption text-grey">
          {{ truncateText(signature.intro, 2) }}
        </div>
      </q-card-section>

      <!-- 右侧操作按钮 -->
      <q-card-actions vertical class="signature-actions">
        <q-btn
          flat
          dense
          round
          icon="edit"
          @click.stop="handleEdit"
          size="sm"
        />
        <q-btn
          flat
          dense
          round
          icon="delete"
          @click.stop="handleDelete"
          size="sm"
          color="negative"
        />
        <q-btn
          flat
          dense
          round
          icon="download"
          @click.stop="handleExport"
          size="sm"
        />
      </q-card-actions>
    </q-card-section>
  </q-card>

  <!-- 图片预览对话框 -->
  <q-dialog
    v-model="showImagePreview"
    backdrop-filter="invert(70%)"
  >
    <q-card class="image-preview-card">
      <q-card-section class="q-pa-none">
        <q-img
          :src="getImageUrl(signature.cardImage)"
          fit="contain"
          style="max-width: 90vw; max-height: 90vh"
        />
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="关闭" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Signature } from 'types/signature';

const props = defineProps<{
  signature: Signature;
}>();

const emit = defineEmits<{
  edit: [id: string];
  delete: [id: string];
  export: [id: string];
}>();

const showImagePreview = ref(false);

function handleInfoClick() {
  emit('edit', props.signature.id);
}

function handleEdit() {
  emit('edit', props.signature.id);
}

function handleDelete() {
  emit('delete', props.signature.id);
}

function handleExport() {
  emit('export', props.signature.id);
}

function truncateText(text: string, lines: number): string {
  if (!text) return '';
  const lineArray = text.split('\n');
  return lineArray.slice(0, lines).join('\n');
}

function getImageUrl(imagePath: string): string {
  // 假设 imagePath 是 "signatures/card_images/xyz789.png"
  return `/signature/image/${imagePath.split('/').pop()}`;
}
</script>

<style scoped>
.signature-list-item {
  min-height: 100px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
}

.signature-image-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  flex-shrink: 0;
  padding: 8px;
}

.signature-thumbnail {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  object-fit: cover;
  cursor: pointer;

  &:hover {
    opacity: 0.9;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
}

.signature-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  border-radius: 8px;
  color: #ccc;
}

.signature-info-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 12px 16px;
}

.signature-name {
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.signature-intro {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  line-height: 1.4;
}

.signature-actions {
  flex-shrink: 0;
  padding: 8px;
}

.image-preview-card {
  background: transparent;
}
</style>
```

### 图片选择器和预览位置（创建/编辑对话框）

**表单布局结构**：

```text
┌────────────────────────────────────────┐
│ 创建/编辑签名                           │
├────────────────────────────────────────┤
│                                        │
│ 签名名称：[________]                   │
│ 个人介绍：[________________]            │
│                                        │
│ 名片图片：[选择图片] [拖拽上传]         │
│                                        │
│ ┌────────────────────────────────────┐ │
│ │   图片快速预览（可点击打开大图）    │ │
│ │        [点击放大预览]              │ │
│ └────────────────────────────────────┘ │
│                                        │
│ [创建/更新]  [取消]                    │
└────────────────────────────────────────┘
```

**关键设计点**：

1. **表单字段顺序**：
   - 签名名称（必填）
   - 个人介绍（选填）
   - 图片选择器

2. **图片预览位置**：
   - **在图片选择器下方**（而不是上方）
   - 尺寸：约 120x120px
   - 可点击打开大图预览

3. **图片选择器**：
   - 支持点击选择文件
   - 支持拖拽上传
   - 文件类型不受严格限制，仅支持 webview 支持的图片格式

**实现示例**：

```vue
<template>
  <q-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    backdrop-filter="invert(70%)"
  >
    <q-card style="min-width: 400px; max-width: 600px">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">
          {{ isEdit ? '编辑签名' : '创建签名' }}
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section>
        <q-form @submit="handleSubmit" class="q-gutter-md">
          <!-- 签名名称 -->
          <q-input
            v-model="formData.name"
            label="签名名称"
            :readonly="isEdit"
            :disable="isEdit"
            filled
            :rules="[
              val => val && val.length > 0 || '请输入签名名称',
              val => val && val.length <= 50 || '不超过50字符'
            ]"
            hint="必填，1-50字符"
          />

          <!-- 个人介绍 -->
          <q-input
            v-model="formData.intro"
            label="个人介绍"
            type="textarea"
            filled
            :rules="[
              val => !val || val.length <= 500 || '不超过500字符'
            ]"
            hint="选填，最多500字符"
            rows="3"
          />

          <!-- 图片选择器 -->
          <div class="q-mt-md">
            <div class="text-subtitle2 q-mb-sm">名片图片</div>
            <q-file
              v-model="formData.cardImageFile"
              label="选择图片或拖拽上传"
              filled
              @update:model-value="handleImageSelect"
              accept="image/*"
              max-file-size="5242880"
              hint="选填，支持 webview 支持的所有图片格式，最大 5MB"
            >
              <template v-slot:prepend>
                <q-icon name="attach_file" />
              </template>
            </q-file>
          </div>

          <!-- 图片快速预览（在选择器下方） -->
          <div v-if="imagePreviewUrl" class="q-mt-md">
            <div class="text-subtitle2 q-mb-sm">预览</div>
            <div class="image-preview-container">
              <q-img
                :src="imagePreviewUrl"
                class="image-preview"
                @click="showFullPreview = true"
              >
                <div class="absolute-full flex flex-center">
                  <div class="text-white text-center">
                    <div class="text-caption">点击放大</div>
                  </div>
                </div>
              </q-img>
            </div>
          </div>

          <!-- 提交按钮 -->
          <div class="q-mt-lg row q-gutter-md">
            <q-btn
              type="submit"
              color="primary"
              :label="isEdit ? '更新' : '创建'"
              class="col"
            />
            <q-btn
              type="button"
              flat
              label="取消"
              color="grey"
              class="col"
              v-close-popup
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>

    <!-- 全屏图片预览对话框 -->
    <q-dialog
      v-model="showFullPreview"
      backdrop-filter="invert(70%)"
    >
      <q-card class="image-full-preview-card">
        <q-card-section class="q-pa-none">
          <q-img
            :src="imagePreviewUrl"
            fit="contain"
            style="max-width: 90vw; max-height: 90vh"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="关闭" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { Signature } from 'types/signature';

interface FormData {
  name: string;
  intro: string;
  cardImageFile: File | null;
}

const props = defineProps<{
  modelValue: boolean;
  signature?: Signature;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  submit: [data: Omit<Signature, 'id' | 'cardImage'> & { cardImageFile: File | null }];
}>();

const isEdit = computed(() => !!props.signature);

const formData = ref<FormData>({
  name: props.signature?.name || '',
  intro: props.signature?.intro || '',
  cardImageFile: null,
});

const imagePreviewUrl = ref<string>('');
const showFullPreview = ref(false);

function handleImageSelect(file: File | null) {
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      imagePreviewUrl.value = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  } else {
    imagePreviewUrl.value = '';
  }
}

function handleSubmit() {
  emit('submit', {
    name: formData.value.name,
    intro: formData.value.intro,
    cardImageFile: formData.value.cardImageFile,
  });
}
</script>

<style scoped>
.image-preview-container {
  display: flex;
  justify-content: center;
}

.image-preview {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    opacity: 0.9;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  }
}

.image-full-preview-card {
  background: transparent;
}
</style>
```

### 图片格式支持

**格式支持原则**：

- 不严格限制图片格式
- 允许用户上传 webview（Chromium）支持的所有图片格式
- 常见支持格式：PNG、JPG/JPEG、WebP、GIF、BMP、SVG 等
- 后端负责验证图片的实际有效性

**前端文件选择器配置**：

```vue
<!-- 使用宽泛的 MIME 类型，允许用户选择所有图片 -->
<q-file
  v-model="imageFile"
  accept="image/*"  <!-- 允许所有图片格式 -->
  max-file-size="5242880"  <!-- 5MB 限制 -->
/>
```

### 尺寸适配说明

**应用固定尺寸限制**：

- KeyTone 应用窗口尺寸为固定值（需要查看项目配置确认具体尺寸）
- 所有组件必须同时适配：
  - **页面模式**：作为独立页面显示在主窗口中
  - **对话框模式**：作为对话框弹出显示

**设计原则**：

1. **响应式布局**：使用 Quasar 的栅格系统（`q-col`, `q-row`）
2. **固定高度**：对话框模式下设置最大高度，内容可滚动
3. **组件复用**：核心组件设计为既能独立使用，也能嵌入对话框
4. **适配方案**：

   ```vue
   <template>
     <div :class="isDialog ? 'signature-dialog-mode' : 'signature-page-mode'">
       <!-- 共享的组件内容 -->
     </div>
   </template>

   <style scoped>
   .signature-page-mode {
     padding: 24px;
     height: 100%;
   }

   .signature-dialog-mode {
     padding: 16px;
     max-height: 80vh;
     overflow-y: auto;
   }
   </style>
   ```

### 保护码处理说明

**保护码生成时机**：

- 保护码在用户点击"创建"或"更新"按钮后，由后端自动生成
- 前端不生成、不显示、不存储保护码
- UI 中完全不可见，用户无需关心保护码

**前端表单简化**：

```vue
<template>
  <q-form @submit="handleSubmit">
    <!-- 签名名称（必填） -->
    <q-input
      v-model="formData.name"
      label="签名名称"
      :rules="[val => val && val.length > 0 || '请输入签名名称']"
      maxlength="50"
    />
    
    <!-- 个人介绍（选填） -->
    <q-input
      v-model="formData.intro"
      label="个人介绍"
      type="textarea"
      maxlength="500"
      hint="选填，最多500字符"
    />
    
    <!-- 名片图片（选填） -->
    <q-file
      v-model="formData.cardImage"
      label="名片图片"
      accept="image/png,image/jpeg,image/jpg,image/gif"
      max-file-size="5242880"
      hint="选填，支持 PNG/JPG/GIF，最大5MB"
    >
      <template v-slot:prepend>
        <q-icon name="attach_file" />
      </template>
    </q-file>
    
    <!-- 图片预览（如果已有图片） -->
    <ImagePreview v-if="previewUrl" :image-url="previewUrl" />
    
    <!-- 提交按钮 -->
    <q-btn
      type="submit"
      color="primary"
      :label="isEdit ? '更新' : '创建'"
    />
  </q-form>
</template>

<script setup lang="ts">
// 注意：没有 protectCode 字段
// 没有 createdAt 字段
// 后端会自动处理这些

    const signatureFile: SignatureFile = {
      version: '1.0.0',
      signature,
      checksum: this.calculateChecksum(signature),
    };

    const json = JSON.stringify(signatureFile, null, 2);
    return new Blob([json], { type: 'application/json' });
  }

  /**
   * 从文件导入签名
   */
  async importSignature(file: File): Promise<Signature> {
    const text = await file.text();
    const signatureFile: SignatureFile = JSON.parse(text);

    // 验证文件格式
    if (!signatureFile.version || !signatureFile.signature) {
      throw new Error('Invalid signature file format');
    }

    // 验证校验和
    const expectedChecksum = this.calculateChecksum(signatureFile.signature);
    if (signatureFile.checksum !== expectedChecksum) {
      throw new Error('Signature file checksum mismatch');
    }

    // 检查是否已存在
    const manager = await this.getAllSignatures();
    if (manager[signatureFile.signature.id]) {
      throw new Error('SIGNATURE_EXISTS');
    }

    await this.saveSignature(signatureFile.signature);
    return signatureFile.signature;
  }

  /**
   * 保存签名到配置文件
   */
  private async saveSignature(signature: Signature): Promise<void> {
    const manager = await this.getAllSignatures();
    manager[signature.id] = signature;
    await this.saveSignatureManager(manager);
  }

  /**
   * 保存签名管理器
   */
  private async saveSignatureManager(manager: SignatureManager): Promise<void> {
    const encrypted = this.encryptSignatureManager(manager);
    await api.post('/store/set', {
      key: 'signature_manager',
      value: encrypted,
    });
  }

  /**
   * 加密签名管理器
   */
  private encryptSignatureManager(manager: SignatureManager): Record<string, string> {
    const encrypted: Record<string, string> = {};

    for (const [id, signature] of Object.entries(manager)) {
      const encryptedId = this.encrypt(id);
      const encryptedData = this.encrypt(JSON.stringify(signature));
      encrypted[encryptedId] = encryptedData;
    }

    return encrypted;
  }

  /**
   * 解密签名管理器
   */
  private decryptSignatureManager(encrypted: Record<string, string>): SignatureManager {
    const manager: SignatureManager = {};

    for (const [encryptedId, encryptedData] of Object.entries(encrypted)) {
      const id = this.decrypt(encryptedId);
      const signature = JSON.parse(this.decrypt(encryptedData));
      manager[id] = signature;
    }

    return manager;
  }

  /**
   * 加密字符串
   */
  private encrypt(data: string): string {
    return CryptoJS.AES.encrypt(data, ENCRYPTION_KEY).toString();
  }

  /**
   * 解密字符串
   */
  private decrypt(encrypted: string): string {
    const bytes = CryptoJS.AES.decrypt(encrypted, ENCRYPTION_KEY);
    return bytes.toString(CryptoJS.enc.Utf8);
  }

  /**
   * 计算签名校验和
   */
  private calculateChecksum(signature: Signature): string {
    const json = JSON.stringify(signature);
    return CryptoJS.SHA256(json).toString();
  }

  /**
   * 将 File 转为 Base64（已在 signature-service.ts 中实现）
   * 注意：图片不需要单独上传，而是在创建/更新签名时作为 Base64 字段传递
   */
  private async fileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      
      reader.onload = () => {
        const result = reader.result as string;
        resolve(result); // 格式：data:image/png;base64,...
      };
      
      reader.onerror = () => {
        reject(new Error('Failed to read file'));
      };
      
      reader.readAsDataURL(file);
    });
  }
}

export const signatureService = new SignatureService();
```

## 专辑签名集成

此部分尚未在当前代码中落地实现，后续如需支持将通过独立 Proposal 与变更集补充详细设计与接口。

## 性能优化

### 缓存策略

```typescript
class SignatureCache {
  private cache: SignatureManager | null = null;
  private cacheTime: number = 0;
  private cacheTTL: number = 5000; // 5秒

  async get(): Promise<SignatureManager> {
    const now = Date.now();

    if (this.cache && (now - this.cacheTime) < this.cacheTTL) {
      return this.cache;
    }

    this.cache = await signatureService.getAllSignatures();
    this.cacheTime = now;
    return this.cache;
  }

  invalidate(): void {
    this.cache = null;
  }
}
```

### SSE 数据同步与缓存管理

**重要说明**：根据现有代码架构，SSE 直接推送全量配置数据，因此不需要单独的监听逻辑和缓存失效机制。

```typescript
// 在前端现有的 SSE 全量数据处理逻辑中（如 app-store.ts）
eventSource.addEventListener('message', (event) => {
  const fullConfig = JSON.parse(event.data);  // 全量配置数据

  // 处理签名管理数据（新增逻辑）
  if (fullConfig.signature) {
    // 直接更新签名状态，无需调用 API 重新获取
  const decryptedSignatures = fullConfig.signature;
    
    // 更新 Pinia store 或响应式状态
    signatureStore.setSignatures(decryptedSignatures);
    
    // Vue 响应式系统会自动更新所有使用该数据的组件
  }
  
  // 处理其他配置字段...
});
```

**缓存策略（可选优化）**：

虽然 SSE 会推送全量数据，但如果需要在 SSE 断开期间临时访问签名数据，可以实现简单的内存缓存：

```typescript
class SignatureCache {
  private cache: SignatureManager | null = null;

  // 从 SSE 接收的数据更新缓存
  update(signatures: SignatureManager): void {
    this.cache = signatures;
  }

  // 获取缓存（仅在 SSE 暂时不可用时使用）
  get(): SignatureManager | null {
    return this.cache;
  }

  // 清空缓存
  clear(): void {
    this.cache = null;
  }
}
```

**关键要点**：

- ✅ SSE 推送全量配置数据，签名数据包含在其中
- ✅ 前端在 SSE 监听器中直接解构和使用签名数据
- ✅ 不需要单独的 `emit('signatures-updated')` 事件
- ✅ 不需要调用 API 重新获取数据（SSE 已提供最新数据）
- ✅ Vue 响应式系统会自动处理 UI 更新
- ✅ 缓存仅作为 SSE 断开时的临时备用方案

## 安全考虑

1. **XSS 防护**：所有用户输入都需要进行 HTML 转义
2. **文件验证**：
   - 验证文件扩展名（.ktsign）
   - 验证文件大小（< 1MB）
   - 验证 JSON 格式
   - 验证校验和
3. **加密强度**：使用 AES-256-GCM 模式
4. **密钥管理**：密钥硬编码在代码中（开源项目限制）
5. **图片安全**：
   - 限制图片大小（< 5MB）
   - 限制图片格式（PNG, JPG, JPEG, GIF）
   - 图片扫描（防止恶意文件）

## 错误处理

### 错误码定义

```typescript
enum SignatureErrorCode {
  INVALID_FILE_FORMAT = 'INVALID_FILE_FORMAT',
  CHECKSUM_MISMATCH = 'CHECKSUM_MISMATCH',
  SIGNATURE_EXISTS = 'SIGNATURE_EXISTS',
  SIGNATURE_NOT_FOUND = 'SIGNATURE_NOT_FOUND',
  ENCRYPTION_FAILED = 'ENCRYPTION_FAILED',
  DECRYPTION_FAILED = 'DECRYPTION_FAILED',
  IMAGE_TOO_LARGE = 'IMAGE_TOO_LARGE',
  INVALID_IMAGE_FORMAT = 'INVALID_IMAGE_FORMAT',
}
```

### 用户友好的错误消息

```typescript
const errorMessages = {
  [SignatureErrorCode.INVALID_FILE_FORMAT]: '签名文件格式无效',
  [SignatureErrorCode.CHECKSUM_MISMATCH]: '签名文件已损坏或被篡改',
  [SignatureErrorCode.SIGNATURE_EXISTS]: '签名已存在，是否覆盖？',
  [SignatureErrorCode.SIGNATURE_NOT_FOUND]: '签名不存在',
  [SignatureErrorCode.ENCRYPTION_FAILED]: '加密失败',
  [SignatureErrorCode.DECRYPTION_FAILED]: '解密失败',
  [SignatureErrorCode.IMAGE_TOO_LARGE]: '图片文件过大（最大5MB）',
  [SignatureErrorCode.INVALID_IMAGE_FORMAT]: '不支持的图片格式',
};
```

## 测试策略

### 单元测试

- [ ] 加密/解密函数
- [ ] 校验和计算
- [ ] 签名验证逻辑
- [ ] 文件格式解析

### 集成测试

- [ ] 签名 CRUD 操作
- [ ] 导入/导出流程
- [ ] SSE 数据同步
- [ ] 专辑签名流程

### E2E 测试

- [ ] 完整的用户操作流程
- [ ] 跨页面导航
- [ ] 错误场景处理

## 迁移与兼容性

### 向后兼容

- 旧版本不包含 `signature` 字段，读取时返回空对象
- 旧版本专辑文件不包含 `signatures` 字段，视为未签名

### 数据迁移

不需要数据迁移，这是全新功能。


## 监控与日志

### 关键指标

- 签名创建成功率
- 签名导入/导出成功率
- 平均响应时间
- 错误率

### 日志记录

```go
// 后端日志
logger.Info("Signature created", "id", signature.ID, "name", signature.Name)
logger.Error("Signature decryption failed", "error", err)
```

```typescript
// 前端日志
console.log('[Signature] Created:', signature.id);
console.error('[Signature] Import failed:', error);
```

## 未来扩展

1. **在线验证**：提供在线签名验证服务
2. **签名链**：支持多重签名和签名历史追溯
3. **吊销机制**：允许用户吊销已发布的签名
4. **公钥加密**：升级为非对称加密方案
5. **区块链集成**：使用区块链技术实现不可篡改的签名记录
