# Design: 优化签名选择对话框实现细节

## 设计概述

本文档说明 `SignaturePickerDialog` 的实现思路，包括 UI 布局重构和真实数据集成的具体方式。

## 1. UI 布局设计

当前布局问题

卡片原布局（上图下文）：

```
┌─────────────────┐
│   图片 80px     │  <- 占用过大
├─────────────────┤
│ 名称             │
│ 介绍（2行）      │
├─────────────────┤
│ 选中指示器       │
└─────────────────┘
```

### 新布局目标（左图右文）

```
┌────────┬────────────────────┬──────┐
│        │ 名称               │      │
│ 图片   │ 介绍（最多2行）    │  ✓   │
│ 60x60  │                    │      │
└────────┴────────────────────┴──────┘
  60px        flex-grow       32px
```

### 实现代码参考

参考签名管理页面 (`Signature_management_page.vue` 第 81-130 行)：

```vuevue
<q-card-section class="q-pa-none" style="display: flex; align-items: center;">
  <!-- 图片区 60x60 -->
  <div style="width: 60px; flex-shrink: 0; height: 60px;">
    <q-img v-if="sig.image" :src="sig.image" style="width: 100%; height: 100%;" />
    <div v-else style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;">
      <!-- 占位符 -->
    </div>
  </div>

  <!-- 信息区 -->
  <div class="col flex flex-col justify-center" style="padding: 0 12px; min-width: 0;">
    <!-- 名称 -->
    <div class="text-subtitle2 text-weight-bold truncate">{{ sig.name }}</div>
    <!-- 介绍 -->
    <div class="text-caption text-grey" style="overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2;">
      {{ sig.intro || '无介绍' }}
    </div>
  </div>

  <!-- 选中指示器 右上角 -->
  <div v-if="selectedId === sig.id" style="margin-left: auto;">
    <q-icon name="check_circle" color="positive" />
  </div>
</q-card-section>
```

### 关键 CSS 类

- `truncate` — 单行截断（名称）
- `-webkit-line-clamp: 2` — 限制 2 行（介绍）
- `flex-shrink: 0` — 图片宽度固定不伸缩
- `col / flex-grow` — 文字区域占用剩余宽度

## 2. 真实数据源集成

### 数据流向

```
签名管理 Store (useSignatureStore)
    ↓
后端 /signature/list API (加密列表)
    ↓
前端 decryptSignatureData() (解密逐项)
    ↓
localSignatures 响应式数组
    ↓
SignaturePickerDialog 渲染
```

### 数据结构

```typescript
interface Signature {
  id: string;              // 加密 ID
  name: string;            // 签名名称 1-50 字符
  intro: string | null;    // 介绍 ≤500 字符
  cardImage: string | null; // 图片文件名，用于 getImageUrl() 获取 Blob URL
}
```

### 函数实现伪代码

```typescripttypescript
async function loadSignaturesRealtime() {
  loading.value = true;
  try {
    // 1. 获取加密列表
    const encryptedSignatures = await getSignaturesList();
    
    // 2. 逐项解密（参考签名管理页面的 decryptAndBuildSignatureMap）
    const signatureMap = new Map<string, Signature>();
    for (const [encryptedId, entry] of Object.entries(encryptedSignatures)) {
      const encryptedValue = typeof entry === 'string' ? entry : entry?.value;
      const decrypted = await decryptSignatureData(encryptedValue, encryptedId);
      const data = JSON.parse(decrypted);
      signatureMap.set(encryptedId, {
        id: encryptedId,
        name: data.name,
        intro: data.intro,
        cardImage: data.cardImage,
      });
    }
    
    // 3. 排序（按 sort.time）
    const sorted = extractSortedSignatures(encryptedSignatures, signatureMap);
    localSignatures.value = sorted;
    
    // 4. 注册 SSE 监听
    registerSseCallback(() => handleSseUpdate());
  } finally {
    loading.value = false;
  }
}
```

## 3. SSE 增量更新逻辑

### 监听时机

- **打开时**：`watch(() => props.visible, (v) => { if (v) loadSignaturesRealtime() })`
- **关闭时**：清理资源和 SSE 监听

### 更新策略

```typescript
async function handleSseUpdate() {
  const newList = await loadSignaturesRealtime(); // 重新加载
  
  // 保留当前搜索和选中状态
  if (selectedId.value && !newList.find(s => s.id === selectedId.value)) {
    // 若选中的签名被删除，清空选中
    selectedId.value = '';
  }
  
  // 搜索词保持不变，重新过滤即可
}
```

## 4. 降级方案：Mock 兼容

若父组件传入 `signatures` prop（非空数组），优先使用：

```typescripttypescript
const finalSignatures = computed(() => {
  if (props.signatures && props.signatures.length > 0) {
    return props.signatures; // 使用 prop（mock 模式）
  }
  return localSignatures.value; // 使用真实数据
});
```

这样既支持测试（传入 mock），也支持生产（不传，自动加载真实）。

## 5. 图片 URL 管理

### Blob URL 的生命周期

```typescript
const imageUrlCache = new Map<string, string>();

function getImageUrl(cardImage: string): string {
  if (!imageUrlCache.has(cardImage)) {
    // 调用 boot/query/signature-query 的 getImageUrl()
    // 返回 data:image/... 或 blob:... URL
    const url = getImageUrlFromApi(cardImage);
    imageUrlCache.set(cardImage, url);
  }
  return imageUrlCache.get(cardImage)!;
}

// 清理
onUnmounted(() => {
  imageUrlCache.forEach(url => {
    if (url.startsWith('blob:')) {
      URL.revokeObjectURL(url);
    }
  });
  imageUrlCache.clear();
});
```

## 6. 搜索逻辑保持不变

```typescript
const filteredSignatures = computed(() => {
  const list = finalSignatures.value;
  if (!searchQuery.value) return list;
  const q = searchQuery.value.toLowerCase();
  return list.filter(sig => sig.name.toLowerCase().includes(q));
});
```

## 7. 错误处理

### 解密失败

若某个签名解密失败（密钥不匹配或数据损坏）：
- console.warn 记录
- 该签名从列表中跳过
- 用户看不到这个项（与签名管理页面一致）

### 网络错误

若 `getSignaturesList()` 失败：
- 显示 error 状态或空列表
- 用户可重试（关闭对话框重新打开）

### SSE 推送异常

若 SSE 回调出错：
- console.error 记录
- 不影响现有列表展示
- 下次手动打开对话框时重新同步

## 8. 性能优化

### 异步加载图片

```typescript
// 不阻塞列表渲染，图片 URL 异步获取
if (sig.cardImage) {
  getImageUrl(sig.cardImage).then(url => {
    // 更新对应签名的图片 URL
  });
}
```

### 搜索防抖（可选）

若性能需求高，可对 `searchQuery` 添加防抖：

```typescript
const debouncedSearch = useDebounceFn((query: string) => {
  searchQuery.value = query;
}, 300);
```

## 9. 与签名管理页面的一致性

| 方面     | 签名管理页                | 选择对话框 |
| -------- | ------------------------- | ---------- |
| 数据源   | `useSignatureStore` + SSE | 同         |
| 解密逻辑 | `decryptSignatureData()`  | 同         |
| 排序方式 | `sort.time` 升序          | 同         |
| 图片获取 | `getImageUrl()`           | 同         |
| 布局样式 | 左图 60x60 + 右文         | 同         |

## 10. 测试覆盖

### 单元测试（可选）

- `loadSignaturesRealtime()` 返回正确的数据结构
- `filteredSignatures` 搜索过滤正确

### 集成测试（人工）

- 对话框打开时数据正确加载
- SSE 推送时自动更新
- 排序改动后对话框同步
- 删除签名后对话框移除项
- 搜索功能正常

---

**总体目标**：通过参考签名管理页面的成熟设计，快速优化选择对话框的 UI 和数据源，保证前后端一致且用户体验一致。
