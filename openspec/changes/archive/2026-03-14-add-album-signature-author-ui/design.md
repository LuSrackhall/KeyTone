# Design: 在专辑选择器中展示签名作者信息

## Context

KeyTone 已实现键音专辑的签名功能，包括：
- 原始作者签名
- 直接导出作者（最近一次导出的作者）
- 签名授权机制

但目前在专辑选择器（主页面、键音专辑页面）中无法展示签名信息，用户无法在选择专辑时识别作者身份。

**约束：**
- 需要与现有 Quasar `q-select` 组件兼容
- 签名信息获取需要在专辑列表遍历过程中完成
- 独立创建 Viper 实例，函数结束后释放以避免内存泄漏

## Goals / Non-Goals

**Goals:**
- 在选择器列表项中展示直接导出作者信息（芯片样式）
- 在选择器选中状态（上边框右侧，Legend 效果）展示直接导出作者信息
- 提供悬停卡片显示详细信息（原始作者 + 直接导出作者）
- 支持点击查看完整签名对话框
- 悬停卡片具有毛玻璃效果

**Non-Goals:**
- 不改变现有签名的数据结构
- 不修改 SignatureAuthorsDialog 组件

## Decisions

### Decision 1: 签名信息获取策略

**选择**：在获取专辑列表时同步获取签名摘要信息

**原因**：
- 避免列表加载后再异步请求造成的闪烁
- 签名摘要信息量小，不会显著增加响应大小

**数据结构**：

```typescript
// 签名摘要（用于列表展示和悬停卡片）
interface AlbumSignatureSummary {
  hasSignature: boolean;              // 是否有签名
  // 原始作者信息
  originalAuthorName: string;         // 原始作者名称
  originalAuthorImage: string;        // 原始作者图片路径
  originalAuthorIntro: string;        // 原始作者介绍
  // 直接导出作者信息
  directExportAuthorName: string;     // 直接导出作者名称
  directExportAuthorImage: string;    // 直接导出作者图片路径
  directExportAuthorIntro: string;    // 直接导出作者介绍
  // 是否为同一作者
  isSameAuthor: boolean;              // 原始作者与直接导出作者是否为同一人
}
```

### Decision 2: Legend 效果实现

**选择**：使用绝对定位 + 背景色遮挡边框

**实现方式**：
 
```css
/* 带有 legend 效果的选择器容器 */
.selector-with-legend-container {
  position: relative;
}

/* 签名徽章包装器：定位在边框上 */
.signature-legend-wrapper {
  position: absolute;
  top: -9px;           /* 垂直定位：与边框顶部对齐 */
  right: 12px;         /* 水平定位：靠右 */
  z-index: 10;
  /* 使用与选择器控件一致的半透明背景，避免黑色矩形 */
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  padding: 0 6px;      /* 左右 padding 确保边框被完全遮挡 */
  border-radius: 999px; /* 轻微圆角，避免背景块显得生硬 */
}
```

**效果示意**：
 
```text
      ┌─────────────────────────[签名徽章]─┐
      │                                    │
      │     专辑名称                        │
      │                                    │
      └────────────────────────────────────┘
```

### Decision 3: 芯片样式徽章

**选择**：签名徽章采用芯片（Chip）样式

**设计特点**：
- 使用 `inline-flex` 确保宽度由内容决定，不占整行
- 圆角胶囊形状（`border-radius: 999px`）
- 琥珀色调背景与签名主题一致

### Decision 4: 悬停卡片内容展示

**选择**：根据 `isSameAuthor` 标记决定展示方式

**展示逻辑**：
- `isSameAuthor=true`：只展示一个作者区块（标签为"原始作者/直接导出作者"）
- `isSameAuthor=false`：展示两个作者区块（原始作者 + 直接导出作者）

**卡片特效**：
- 毛玻璃效果：`backdrop-filter: blur(8px)`
- 半透明白色背景：`background: rgba(255, 255, 255, 0.85)`

### Decision 5: 悬停卡片消失逻辑

**选择**：使用独立的状态标志 + 延迟检查

**实现逻辑**：

```typescript
const isHoverOnTrigger = ref(false);  // 鼠标是否在触发区域
const isHoverOnCard = ref(false);      // 鼠标是否在卡片上
const isCardVisible = ref(false);      // 卡片是否可见

// 显示延迟：避免快速划过时闪烁
const SHOW_DELAY = 200;
// 隐藏延迟：确保鼠标能移动到卡片上
const HIDE_DELAY = 100;

// 检查是否应该隐藏
function checkAndHideCard() {
  if (!isHoverOnTrigger.value && !isHoverOnCard.value) {
    isCardVisible.value = false;
  }
}
```

**列表项点击打开对话框**：
- 在选择器列表中点击“查看详细信息”时，先关闭 QSelect 弹层，再延迟打开 `SignatureAuthorsDialog`，避免弹层关闭流程打断对话框打开。

## Component Architecture

 
```text
album-selector/
├── AlbumSignatureBadge.vue      # 签名徽章（芯片样式）
├── AlbumSignatureHoverCard.vue  # 悬停详情卡片（毛玻璃效果）
└── index.ts                     # 导出入口
```

### AlbumSignatureBadge.vue

```vue
<template>
  <!-- 签名徽章（芯片样式） -->
  <div class="album-signature-badge" :class="size">
    <div class="avatar-container">
      <img v-if="imageUrl" :src="imageUrl" />
      <q-icon v-else name="badge" />
    </div>
    <span class="author-name">{{ authorName }}</span>
  </div>
</template>
```

**样式特点**：
- `display: inline-flex`：宽度自适应
- `border-radius: 999px`：胶囊形状
- 支持 `normal` 和 `small` 两种尺寸

### AlbumSignatureHoverCard.vue

```vue
<template>
  <!-- 触发区域 -->
  <div @mouseenter @mouseleave>
    <slot />
  </div>
  
  <!-- 悬停卡片（Teleport + 毛玻璃） -->
  <Teleport to="body">
    <div v-if="isCardVisible" class="signature-hover-card">
      <!-- isSameAuthor=true: 一个作者区块 -->
      <!-- isSameAuthor=false: 两个作者区块 -->
      <!-- "点击查看详细信息" -->
    </div>
  </Teleport>
</template>
```

**Props**：
- `albumPath`: 专辑路径
- `signatureInfo`: `AlbumSignatureSummary` 类型的签名摘要

## API Changes

### Backend: GetAudioPackageList 响应扩展

```json
{
  "list": ["path1", "path2"],
  "signatureInfo": {
    "path1": {
      "hasSignature": true,
      "originalAuthorName": "原始作者",
      "originalAuthorImage": "signature/card_orig.jpg",
      "originalAuthorIntro": "原始作者介绍",
      "directExportAuthorName": "导出作者",
      "directExportAuthorImage": "signature/card_export.jpg",
      "directExportAuthorIntro": "导出作者介绍",
      "isSameAuthor": false
    }
  }
}
```

### Backend: AlbumSignatureSummary 结构体

```go
type AlbumSignatureSummary struct {
    HasSignature            bool   `json:"hasSignature"`
    // 原始作者信息
    OriginalAuthorName      string `json:"originalAuthorName"`
    OriginalAuthorImage     string `json:"originalAuthorImage"`
    OriginalAuthorIntro     string `json:"originalAuthorIntro"`
    // 直接导出作者信息
    DirectExportAuthorName  string `json:"directExportAuthorName"`
    DirectExportAuthorImage string `json:"directExportAuthorImage"`
    DirectExportAuthorIntro string `json:"directExportAuthorIntro"`
    // 是否为同一作者
    IsSameAuthor            bool   `json:"isSameAuthor"`
}
```

## Risks / Trade-offs

| 风险                             | 影响 | 缓解措施                             |
| -------------------------------- | ---- | ------------------------------------ |
| 签名信息获取增加列表加载时间     | 中   | 签名摘要数据量小；可考虑后期增加缓存 |
| Legend 背景色硬编码              | 低   | 需要与页面背景保持一致               |
| 毛玻璃效果在某些浏览器可能不支持 | 低   | 使用回退的半透明背景                 |

## Open Questions

1. ~~是否需要为没有签名的专辑显示占位符？~~ - **已决定**：不显示占位符
2. ~~悬停卡片展示内容？~~ - **已决定**：展示原始作者和直接导出作者，相同时只展示一个
3. ~~悬停卡片消失逻辑？~~ - **已决定**：鼠标离开触发区域和卡片后延迟消失
