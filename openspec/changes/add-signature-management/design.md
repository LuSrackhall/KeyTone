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
│  │  /store/get  │  │  /store/set  │  │    /stream   │ │
│  │  (GET)       │  │  (POST)      │  │    (SSE)     │ │
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
│    "signature_manager": {                             │
│      "加密的保护码1": "加密的签名数据1",                │
│      "加密的保护码2": "加密的签名数据2"                 │
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

### 后端存储结构

```json
{
  "signature_manager": {
    "encrypted_id_1": "encrypted_signature_data_1",
    "encrypted_id_2": "encrypted_signature_data_2"
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

所有加密/解密逻辑在 `sdk/signature/encryption.go` 中实现：

```go
package signature

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
    "github.com/jaevor/go-nanoid"
)

const encryptionKey = "KeyTone2024SecretKey_SignatureProtection"

// GenerateProtectCode 生成21位保护码（使用nanoid算法）
func GenerateProtectCode() (string, error) {
    canonicGenerator, err := nanoid.Standard(21)
    if err != nil {
        return "", err
    }
    return canonicGenerator(), nil
}

// EncryptSignature 加密签名数据
func EncryptSignature(data string) (string, error) {
    block, err := aes.NewCipher([]byte(encryptionKey))
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSignature 解密签名数据
func DecryptSignature(encryptedData string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher([]byte(encryptionKey))
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

## API 设计

### 签名管理 API

签名管理需要新增专用端点，因为涉及加密/解密逻辑和保护码生成。

#### 创建签名

```http
POST /signature/create
Content-Type: application/json

{
  "name": "张三",
  "intro": "资深音效设计师",
  "cardImage": "data:image/png;base64,..." // Base64 编码的图片（可选）
}
```

响应：

```json
{
  "message": "ok",
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师",
    "cardImage": "signatures/card_images/xyz789.png"
  }
}
```

**处理流程**：
1. 后端接收前端发来的签名数据
2. 生成签名 ID（nanoid，21位）
3. 生成保护码（nanoid，21位，用于加密，UI 中不可见）
4. 如果有 cardImage（Base64），解码并保存到本地文件系统
5. 加密签名 ID 和签名数据
6. 存储到配置文件
7. **配置文件更新后，现有的 SSE 机制会自动推送全量配置数据**（无需单独实现通知）
8. 返回未加密的签名数据（供前端显示）

#### 获取所有签名

```http
GET /signature/list
```

响应：

```json
{
  "message": "ok",
  "signatures": {
    "abc123...": {
      "id": "abc123...",
      "name": "张三",
      "intro": "资深音效设计师",
      "cardImage": "signatures/card_images/xyz789.png"
    },
    "def456...": {
      "id": "def456...",
      "name": "李四",
      "intro": "独立开发者",
      "cardImage": ""
    }
  }
}
```

**处理流程**：
1. 从配置文件读取加密的签名管理器
2. 解密所有签名数据
3. 返回解密后的签名列表

#### 更新签名

```http
PUT /signature/update
Content-Type: application/json

{
  "id": "abc123...",
  "name": "张三",
  "intro": "资深音效设计师 + 键盘爱好者",
  "cardImage": "data:image/png;base64,..." // 新图片（可选）
}
```

响应：

```json
{
  "message": "ok",
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师 + 键盘爱好者",
    "cardImage": "signatures/card_images/new_hash.png"
  }
}
```

#### 删除签名

```http
DELETE /signature/delete/:id
```

响应：

```json
{
  "message": "ok"
}
```

#### 导出签名

```http
GET /signature/export/:id
```

响应：

```json
{
  "version": "1.0.0",
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师",
    "cardImage": "data:image/png;base64,..." // 图片转为 Base64
  },
  "checksum": "sha256_hash"
}
```

**处理流程**：
1. 根据 ID 查找签名
2. 如果有图片，读取图片文件并转为 Base64
3. 计算签名数据的 SHA-256 校验和
4. 返回 .ktsign 格式的 JSON
5. 前端接收 JSON，使用 `window.showSaveFilePicker()` 保存文件

#### 导入签名

```http
POST /signature/import
Content-Type: application/json

{
  "version": "1.0.0",
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师",
    "cardImage": "data:image/png;base64,..."
  },
  "checksum": "sha256_hash"
}
```

响应：

```json
{
  "message": "ok",
  "exists": false,  // 如果签名已存在则为 true
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师",
    "cardImage": "signatures/card_images/restored_hash.png"
  }
}
```

**处理流程**：
1. 验证文件格式和版本
2. 验证校验和
3. 检查签名 ID 是否已存在
4. 如果存在，返回 `exists: true`，前端提示用户是否覆盖
5. 如果不存在或用户确认覆盖：
   - 将 Base64 图片解码并保存到文件系统
   - 生成新的保护码（因为不保存原保护码）
   - 加密并存储签名数据
6. 返回导入结果

#### 获取图片

```http
GET /signature/image/:filename
```

响应：图片文件的二进制数据，Content-Type 为对应的图片类型（image/png, image/jpeg 等）

**处理流程**：
1. 验证文件名合法性（防止路径遍历攻击）
2. 从 `signatures/card_images/` 目录读取图片
3. 设置正确的 Content-Type
4. 返回图片二进制数据

**前端使用示例**：

```typescript
<img :src="`http://localhost:port/signature/image/${signature.cardImage.split('/').pop()}`" />
```

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
- 在前端现有的 SSE 全量数据处理逻辑中，添加对 `signature_manager` 字段的解构和处理
- 当 SSE 推送全量配置数据时，提取其中的 `signature_manager` 数据块
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
  if (fullConfig.signature_manager) {
    // 解构获取签名管理数据
    const signatureData = fullConfig.signature_manager;
    
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
解构提取 signature_manager 字段
    ↓
更新前端签名状态
    ↓
Vue 响应式自动刷新 UI
```

**关键要点**：

- ✅ 现有代码已有 SSE 全量数据推送，无需重新设计
- ✅ 签名数据作为配置文件的一部分，会随全量数据一起推送
- ✅ 前端只需在现有的全量数据处理逻辑中添加签名数据的解构和适配
- ✅ 不需要单独的 `GET /signature/list` API 调用来刷新数据
- ✅ 所有签名 CRUD 操作完成后，后端保存配置文件，自动触发 SSE 推送

## 文件处理

### 名片图片存储

1. **目录结构**：

   ```text
   配置目录/
   ├── KeyToneSetting.json
   └── signatures/
       └── card_images/
           ├── abc123...def.png
           └── xyz789...uvw.jpg
   ```

   **注意**：不再需要 `exported/` 目录，导出文件由用户选择保存位置。

2. **图片命名**：使用 SHA-256 哈希值作为文件名（由后端处理）

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
   
   #### 阶段 2：提交存储阶段 - Base64 传输与后端持久化
   
   当用户点击"创建"或"更新"按钮时:
   
   - **Base64 转换**：前端将 `File` 对象转为 Base64 字符串（通过 `fileToBase64()` 方法）
   - **HTTP 传输**：Base64 字符串通过 HTTP POST/PUT 发送到后端
   - **后端处理**：
     1. 接收 Base64 字符串
     2. 解码为二进制图片数据
     3. 计算 SHA-256 哈希值作为文件名
     4. 保存图片文件到 `signatures/card_images/` 目录
     5. 在配置文件中只存储**文件路径字符串**（如 `signatures/card_images/xyz789.png`）
   
   ```typescript
   // HTTP 请求体（Base64 传输）
   {
     name: "张三",
     intro: "资深音效设计师",
     cardImage: "data:image/png;base64,iVBORw0KGg..."  // Base64 字符串
   }
   ```
   
   ```json
   // 后端配置文件存储（加密后的路径字符串）
   {
     "signature_manager": {
       "encrypted_id": "encrypted_data_contains_path_string"
     }
   }
   ```
   
   #### 阶段 3：列表渲染阶段 - 路径字符串到图片资源
   
   当前端需要显示签名列表时:
   
   - **获取路径**：通过 API 获取签名数据,其中 `cardImage` 字段是**路径字符串**（如 `signatures/card_images/xyz789.png`）
   - **路径转换**：前端将路径字符串转为 HTTP URL
   - **图片渲染**：通过 HTTP 请求获取图片二进制数据并渲染
   
   ```typescript
   // 前端签名数据结构（列表显示阶段）
   interface Signature {
     id: string;
     name: string;
     intro: string;
     cardImage: string;  // 注意：这里是路径字符串,如 "signatures/card_images/xyz789.png"
   }
   
   // 前端渲染时转换路径为 URL
   function getImageUrl(cardImagePath: string): string {
     const filename = cardImagePath.split('/').pop();  // 提取文件名
     return `${baseURL}/signature/image/${filename}`;  // 转为 HTTP URL
   }
   ```
   
   ```vue
   <!-- 在模板中使用 -->
   <q-img :src="getImageUrl(signature.cardImage)" />
   ```
   
   #### 阶段 4：导出/导入阶段 - 文件系统与 Base64 互转
   
   **导出流程**：
   - 前端调用 `/signature/export/:id` API
   - 后端读取配置文件中的路径字符串
   - 后端根据路径读取图片文件,转为 Base64
   - 后端返回包含 Base64 图片的 JSON 数据
   - 前端保存为 .ktsign 文件
   
   **导入流程**：
   - 前端读取 .ktsign 文件,获取 Base64 图片
   - 前端调用 `/signature/import` API,传递 Base64 数据
   - 后端解码 Base64 为二进制数据
   - 后端计算哈希值,保存为新图片文件
   - 后端在配置文件中存储新的路径字符串
   
   #### 图片数据形态总结
   
   | 阶段       | 图片数据形态                          | 说明                                  |
   | ---------- | ------------------------------------- | ------------------------------------- |
   | 选择器选择 | `File` 对象                           | 前端内存中的文件对象                  |
   | 表单编辑   | `File` 对象 / 预览用 Base64           | 未提交前仍是 File 对象,不是路径字符串 |
   | HTTP 传输  | Base64 字符串                         | 通过 HTTP 传输到后端                  |
   | 后端存储   | 二进制文件 + 路径字符串（配置文件中） | 图片存为文件,配置中只存路径           |
   | 列表显示   | 路径字符串 → HTTP URL → 渲染图片      | 通过 HTTP 接口访问图片资源            |
   | 导出文件   | Base64 字符串（嵌入 JSON）            | 便于文件独立传输                      |
   | 导入文件   | Base64 字符串 → 二进制文件            | 恢复为文件系统中的图片                |
   
   **关键要点**：
   - ✅ 创建/编辑表单中,图片是 `File` 对象,**不是路径字符串**
   - ✅ HTTP 传输使用 Base64 字符串
   - ✅ 配置文件中只存储路径字符串,图片本体存在文件系统
   - ✅ 列表渲染时,路径字符串通过 HTTP 接口转为真实图片
   - ✅ 导出/导入使用 Base64 嵌入 JSON 文件

### .ktsign 文件格式

导出的签名文件包含 Base64 编码的图片数据：

```json
{
  "version": "1.0.0",
  "signature": {
    "id": "abc123...",
    "name": "张三",
    "intro": "资深音效设计师",
    "cardImage": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  },
  "checksum": "sha256_hash_of_signature_data"
}
```

导出流程：

1. 前端调用 `GET /signature/export/:id`
2. 后端读取签名数据
3. 如果有图片，读取图片文件并转为 Base64（格式：`data:image/png;base64,...`）
4. 计算签名数据的 SHA-256 校验和
5. 返回 JSON 格式数据
6. 前端使用 `window.showSaveFilePicker()` API 让用户选择保存路径
7. 将 JSON 数据写入用户选择的文件

导入流程：

1. 前端读取用户选择的 `.ktsign` 文件内容
2. 调用 `POST /signature/import` 传递文件内容
3. 后端验证文件格式和版本
4. 验证 checksum
5. 解析签名数据
6. 如果有 Base64 图片，解码并保存到 `signatures/card_images/`
7. 检查签名 ID 是否已存在
8. 如已存在，返回 `exists: true`，前端提示用户是否覆盖
9. 保存到配置文件（加密存储）
10. **配置文件更新后，现有的 SSE 机制会自动推送全量配置数据**（无需单独实现通知）

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

### 图片预览功能

**图片预览组件要求**：

1. **缩略图预览**：在创建/编辑表单中显示缩略图（如 120x120px）
2. **点击放大**：点击缩略图打开全尺寸预览对话框
3. **预览对话框特性**：
   - 显示原始尺寸图片
   - 支持图片缩放（鼠标滚轮）
   - 支持图片平移（鼠标拖拽）
   - 背景使用 `backdrop-filter="invert(70%)"`
   - 提供关闭按钮

**实现示例**：

```vue
<template>
  <!-- 缩略图 -->
  <q-img
    v-if="cardImageUrl"
    :src="cardImageUrl"
    class="card-image-thumbnail"
    @click="showPreview = true"
  >
    <div class="absolute-bottom text-subtitle2 text-center">
      点击放大
    </div>
  </q-img>
  
  <!-- 放大预览对话框 -->
  <q-dialog
    v-model="showPreview"
    backdrop-filter="invert(70%)"
  >
    <q-card class="image-preview-card">
      <q-card-section class="q-pa-none">
        <q-img
          :src="cardImageUrl"
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
import { ref, computed } from 'vue';
import { signatureService } from 'services/signature-service';

const props = defineProps<{
  cardImagePath?: string;
}>();

const showPreview = ref(false);

const cardImageUrl = computed(() => {
  return props.cardImagePath 
    ? signatureService.getImageUrl(props.cardImagePath)
    : '';
});
</script>

<style scoped>
.card-image-thumbnail {
  width: 120px;
  height: 120px;
  cursor: pointer;
  border-radius: 8px;
}

.image-preview-card {
  background: transparent;
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

### 专辑文件格式扩展

在现有的 `.ktalbum` 文件中添加签名信息：

```json
{
  "magicNumber": "KTAF",
  "version": "1.0.0",
  "exportTime": "2025-10-15T10:30:00.000Z",
  "albumUUID": "abc123...",
  "albumName": "My Album",
  "signatures": [
    {
      "signatureId": "xyz789...",
      "signatureName": "张三",
      "signedAt": "2025-10-15T10:30:00.000Z",
      "protectCode": "encrypted_protect_code"
    }
  ],
  "config": {...},
  "sounds": {...}
}
```

### 签名选择流程

1. 用户点击"导出专辑"
2. 显示签名选择对话框
3. 用户选择签名（或选择"不签名"）
4. 如果选择签名，嵌入签名信息到专辑文件
5. 继续原有的导出流程

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
  if (fullConfig.signature_manager) {
    // 直接更新签名状态，无需调用 API 重新获取
    const decryptedSignatures = decryptSignatureManager(fullConfig.signature_manager);
    
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

- 旧版本不包含 `signature_manager` 字段，读取时返回空对象
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
