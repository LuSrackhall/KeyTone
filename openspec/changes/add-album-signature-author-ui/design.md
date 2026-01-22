# Design: 在专辑选择器中展示签名作者信息

## Context

KeyTone 已实现键音专辑的签名功能，包括：
- 原始作者签名
- 直接导出作者（最近一次导出的作者）
- 签名授权机制

但目前在专辑选择器（主页面、键音专辑页面）中无法展示签名信息，用户无法在选择专辑时识别作者身份。

**约束：**
- 需要与现有 Quasar `q-select` 组件兼容
- 签名信息获取需要在专辑列表遍历过程中完成，复用 Viper 实例
- 遍历后需及时释放 Viper 实例以释放内存

## Goals / Non-Goals

**Goals:**
- 在选择器列表项中展示直接导出作者信息
- 在选择器选中状态（上边框右侧）展示直接导出作者信息
- 提供悬停卡片显示详细信息
- 支持点击查看完整签名对话框

**Non-Goals:**
- 不改变现有签名的数据结构
- 不修改 SignatureAuthorsDialog 组件

## Decisions

### Decision 1: 签名信息获取策略

**选择**：在获取专辑列表时同步获取签名摘要信息

**原因**：
- 避免列表加载后再异步请求造成的闪烁
- 复用现有的专辑遍历逻辑和 Viper 实例
- 签名摘要信息量小，不会显著增加响应大小

**数据结构**：
```typescript
// 签名摘要（轻量级，用于列表展示）
interface AlbumSignatureSummary {
  hasSignature: boolean;           // 是否有签名
  directExportAuthorName: string;  // 直接导出作者名称
  directExportAuthorImage: string; // 直接导出作者图片路径（相对于专辑目录）
}
```

### Decision 2: 选择器自定义策略

**选择**：优先使用 Quasar `q-select` 插槽机制

**使用的插槽**：
- `option`: 自定义列表项内容
- `before-options`: 无（不需要）
- 通过 CSS 定位实现选中状态的签名展示

**备选方案**：若插槽方案无法满足需求，创建自定义的 `AlbumSelector.vue` 组件

### Decision 3: 悬停卡片交互实现

**选择**：自定义组件 + `@mouseenter`/`@mouseleave` 事件

**原因**：
- Quasar 的 `q-tooltip` 不支持鼠标移入保持显示
- 需要自定义卡片内容和"点击查看详细信息"交互

**实现要点**：
```typescript
// 悬停状态管理
const isHoverOnTrigger = ref(false);
const isHoverOnCard = ref(false);
const showCard = computed(() => isHoverOnTrigger.value || isHoverOnCard.value);

// 延迟隐藏，确保鼠标移动到卡片上时不会闪烁
const hideDelay = 150; // ms
```

### Decision 4: 图片加载策略

**选择**：使用 `GetAlbumFile` API 获取签名图片

**原因**：
- 签名图片存储在专辑目录内
- 需要通过后端 API 读取（与 SignatureAuthorsDialog 一致）

**缓存策略**：
- 使用 `URL.createObjectURL()` 创建 Blob URL
- 组件卸载时调用 `URL.revokeObjectURL()` 释放

## Component Architecture

```
album-selector/
├── AlbumSignatureBadge.vue      # 签名徽章（图片+名称）
├── AlbumSignatureHoverCard.vue  # 悬停详情卡片
└── index.ts                     # 导出入口
```

### AlbumSignatureBadge.vue

```vue
<template>
  <!-- 签名作者展示徽章 -->
  <div class="album-signature-badge">
    <!-- 作者头像/签名图标 -->
    <img v-if="imageUrl" :src="imageUrl" />
    <q-icon v-else name="badge" />
    <!-- 作者名称 -->
    <span>{{ authorName }}</span>
  </div>
</template>

<script setup lang="ts">
// Props: albumPath, size ('small' | 'normal')
// 负责：加载图片、展示基本信息
</script>
```

### AlbumSignatureHoverCard.vue

```vue
<template>
  <!-- 触发区域（包裹 AlbumSignatureBadge） -->
  <div @mouseenter="onTriggerEnter" @mouseleave="onTriggerLeave">
    <slot />
  </div>
  
  <!-- 悬停卡片（Teleport 到 body） -->
  <Teleport to="body">
    <div v-if="showCard" class="signature-hover-card">
      <!-- 详细信息展示 -->
      <!-- "点击查看详细信息" label -->
    </div>
  </Teleport>
</template>
```

## API Changes

### Backend: GetAudioPackageList 响应扩展

```json
{
  "list": ["path1", "path2"],
  "signatureInfo": {
    "path1": {
      "hasSignature": true,
      "directExportAuthorName": "作者名称",
      "directExportAuthorImage": "signature/card_xxxx.jpg"
    },
    "path2": {
      "hasSignature": false
    }
  }
}
```

### Frontend: main-store 扩展

```typescript
// 新增状态
const keyTonePkgSignatureInfo = ref<Map<string, AlbumSignatureSummary>>(new Map());

// 修改 GetKeyToneAlbumList
function GetKeyToneAlbumList() {
  GetAudioPackageList().then((res) => {
    if (res.list) {
      keyTonePkgOptions.value = res.list;
      // 处理签名信息
      keyTonePkgSignatureInfo.value.clear();
      if (res.signatureInfo) {
        Object.entries(res.signatureInfo).forEach(([path, info]) => {
          keyTonePkgSignatureInfo.value.set(path, info);
        });
      }
    }
  });
}
```

## Risks / Trade-offs

| 风险                             | 影响 | 缓解措施                             |
| -------------------------------- | ---- | ------------------------------------ |
| 签名信息获取增加列表加载时间     | 中   | 签名摘要数据量小；可考虑后期增加缓存 |
| q-select 插槽自定义受限          | 中   | 预留自定义组件方案作为备选           |
| 悬停卡片定位在不同布局下可能异常 | 低   | 使用 Teleport + 动态计算位置         |

## Open Questions

1. 是否需要为没有签名的专辑显示占位符（保持列表项高度一致）？
   - 建议：不显示占位符，无签名时不展示签名区域

2. 主页面 dialog 模式的选择器是否需要特殊处理？
   - 建议：复用相同的列表项组件，dialog 模式下自动适配
