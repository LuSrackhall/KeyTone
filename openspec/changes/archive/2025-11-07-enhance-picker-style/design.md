# Design: 签名选择对话框样式优化实现细节

## 总体设计思路

通过增量式的 CSS 和交互优化，在保持现有功能的基础上，提升对话框的视觉品质和用户体验。核心思路是参考现代 UI 设计的微交互、分层和固定底栏等常见模式。

## 1. 卡片视觉层次设计

### 阴影策略

采用三态阴影系统：

```css
.signature-card {
  /* 正常态：微妙阴影 */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;

  &:hover {
    /* hover 态：增强阴影 + 上升 */
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transform: translateY(-1px);
  }

  &.selected {
    /* 选中态：边框 + 光晕 */
    border: 2px solid var(--q-primary);
    border-width: 2px !important;
    box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.1), 0 4px 16px rgba(33, 150, 243, 0.15);
  }
}
```

### 视觉效果

- **正常**：细微阴影（y 偏移 1px，模糊 3px，透明度 5%）
- **Hover**：增强阴影（y 偏移 4px，模糊 12px，透明度 8%）+ 垂直上升 1px
- **Selected**：蓝色边框 + 蓝色双层光晕（内层浅色，外层深色）

## 2. 顶部区域粘滞

### 结构调整

顶部的 Header 和 Search Bar 被包装在一个 `.sticky-top` 容器中：

```vue
<div class="sticky-top">
  <!-- Header -->
  <q-card-section class="bg-primary text-white q-pa-sm">
    <div class="text-subtitle1">{{ $t('exportFlow.pickerDialog.title') }}</div>
  </q-card-section>

  <!-- Search Bar -->
  <q-card-section class="q-pa-md">
    <!-- 说明 + 搜索输入 + 创建按钮 -->
  </q-card-section>
</div>

<!-- 只有签名列表在可滚动区域 -->
<q-card-section class="col-grow overflow-auto scrollable-content">
  <!-- 签名列表 -->
</q-card-section>
```

### CSS 实现

```scss
.sticky-top {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 5;
  background: white;
}
```

### 效果

- Header 和 Search Bar 固定在对话框顶部
- 用户滚动签名列表时，Header 始终可见
- 便于持续搜索和创建签名

---

## 3. 内容区域优化

## 2. 内容区域的布局优化

### 竖向间距设计

```vue
<div class="col flex flex-col justify-center q-py-2xs q-px-sm">
  <!-- 名称 -->
  <div class="name-container">{{ sig.name }}</div>
  
  <!-- 介绍：上方 4px 间距 -->
  <div class="intro-container q-mt-2xs">{{ sig.intro }}</div>
</div>
```

- `q-py-2xs` = 2px 上下 padding（紧凑）
- `q-mt-2xs` = 4px margin-top（名称与介绍间距）
- `q-px-sm` = 8px 左右 padding（与图片间距一致）

### 横向滚动实现

对名称和介绍各自添加滚动容器类：

```vue
<!-- 名称：单行，超出显示滚动条 -->
<div class="name-container text-subtitle2 text-weight-bold">
  <div class="scrollable-x whitespace-nowrap">{{ sig.name }}</div>
</div>

<!-- 介绍：最多 2 行，超出显示滚动条 -->
<div class="intro-container text-caption text-grey">
  <div class="scrollable-x line-clamp-2">{{ sig.intro }}</div>
</div>
```

CSS：

```css
.scrollable-x {
  overflow-x: auto;
  overflow-y: hidden;
  white-space: nowrap;
  
  /* 自定义滚动条 */
  &::-webkit-scrollbar {
    height: 3px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
  }
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
```

## 3. 选择交互逻辑优化

### Toggle 选择

修改 `selectSignature()` 函数支持 toggle：

```typescript
const selectSignature = (id: string) => {
  // 若点击已选项，则取消选择；否则选中
  if (selectedId.value === id) {
    selectedId.value = '';
  } else {
    selectedId.value = id;
  }
};
```

### 选中图标的光晕效果

```vue
<!-- 右侧选中指示器 -->
<div v-if="selectedId === sig.id" class="flex-shrink-0 selection-indicator-wrapper">
  <!-- 光晕层 -->
  <div class="selection-glow"></div>
  <!-- 图标 -->
  <q-icon name="check_circle" size="20px" color="positive" class="selection-icon" />
</div>
```

CSS 动画：

```css
.selection-indicator-wrapper {
  position: relative;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 8px;
  margin-right: 8px;
}

.selection-glow {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(33, 150, 243, 0.3), transparent);
  animation: glow-pulse 1.5s ease-in-out infinite;
}

.selection-icon {
  position: relative;
  z-index: 1;
}

@keyframes glow-pulse {
  0% {
    opacity: 0.8;
    transform: scale(1);
  }
  
  50% {
    opacity: 0.4;
  }
  
  100% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}
```

光晕效果：
- 起始：半径为图标 1.2 倍，opacity 0.8
- 中点：opacity 0.4
- 终点：opacity 0.8，缩放至 1.1 倍
- 循环周期：1.5s

## 4. 底部操作栏固定

### Sticky 定位与毛玻璃背景

```vue
<q-card-actions
  align="right"
  class="q-pa-sm q-gutter-xs sticky-bottom"
>
  <!-- 按钮... -->
</q-card-actions>
```

CSS 的关键点是**整个 `q-card-actions` 容器**都应用毛玻璃效果，**且背景高度透明**：

```css
.sticky-bottom {
  position: sticky;
  bottom: 0;
  left: 0;
  right: 0;
  
  /* 毛玻璃背景：模糊 + 高度透明 */
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.7);  /* 更透明，允许清晰透视 */
  
  /* 分割线 */
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  
  /* 始终在上方 */
  z-index: 10;
  
  /* 过渡动画 */
  transition: all 0.2s ease;
}
```

### 毛玻璃效果的特点

- **Backdrop-filter blur(10px)**：模糊后方的列表项
- **高度透明白色背景 rgba(255,255,255,0.7)**：允许**清晰透视**看到后方内容
- **效果**：用户滚动时可以**清晰看到列表项通过毛玻璃显示**，具有层次感和现代感
- **高级感**：这种设计常见于现代应用（iOS、macOS、现代网页应用）

## 5. 滚动条优化

### 对话框侧边（垂直）滚动条

这是指对话框**右侧边缘**的纵向滚动条，用于滚动签名列表。

#### 实现方式：Tailwind CSS 任意值 + `overflow: overlay`

在对话框首个 q-card 上应用滚动条样式，使用 **Tailwind CSS 的任意值语法**：

```vue
<q-card
  class="signature-picker-dialog 
    [overflow:overlay] 
    [&::-webkit-scrollbar]:w-1 
    [&::-webkit-scrollbar-track]:bg-transparent 
    [&::-webkit-scrollbar-thumb]:bg-slate-400/40 
    [&::-webkit-scrollbar-thumb]:rounded 
    [&::-webkit-scrollbar-thumb]:hover:bg-slate-600/60"
>
  <!-- 对话框内容 -->
</q-card>
```

**关键设计特性**：

| 属性                                                 | 说明                                               |
| ---------------------------------------------------- | -------------------------------------------------- |
| `[overflow:overlay]`                                 | 滚动条**不占用对话框空间**，而是透明叠加在内容之上 |
| `[&::-webkit-scrollbar]:w-1`                         | 滚动条宽度 = `0.25rem`（~4px）                     |
| `[&::-webkit-scrollbar-track]:bg-transparent`        | 滚动条轨道透明，完全透视后方内容                   |
| `[&::-webkit-scrollbar-thumb]:bg-slate-400/40`       | 滚动条颜色半透明深灰（slate-400，透明度 40%）      |
| `[&::-webkit-scrollbar-thumb]:rounded`               | 滚动条圆角                                         |
| `[&::-webkit-scrollbar-thumb]:hover:bg-slate-600/60` | Hover 时颜色加深（slate-600，透明度 60%）          |

### 为什么使用 `overflow: overlay`？

- **标准 `overflow: auto`**：滚动条占用空间，内容宽度被挤压
- **`overflow: overlay`**：滚动条透明叠加在内容之上，**不占用空间**
- **用户体验**：内容始终显示完整宽度，滚动条不打断布局
- **浏览器兼容性**：大多数现代浏览器支持（Chrome、Safari、Edge、Firefox）

### 效果对比

| 方面     | 优化前         | 优化后                     |
| -------- | -------------- | -------------------------- |
| 宽度     | 6px            | 4px（更细致）              |
| 颜色     | 浏览器默认     | slate-400/40（半透明深灰） |
| 占用空间 | 占用对话框宽度 | 不占用空间（overlay）      |
| 叠加效果 | 无             | 透明叠加在内容之上         |
| Hover    | 无反馈         | slate-600/60（加深）       |

## 6. 整体布局结构

```text
┌─────────────────────────────────┐
│ Header (bg-primary)             │  <- 标题
├─────────────────────────────────┤
│ Description + Search Bar        │  <- 说明和搜索
├─────────────────────────────────┤
│ [Signature Items - Scrollable]  │  <- 内容区（可滚动）
│  ┌─────────┬─────────────────┐  │
│  │ Image   │ Name    ↔       │ ✓ │  <- 单个卡片
│  │ 60x60   │ Intro ↔ 2 lines │   │
│  └─────────┴─────────────────┘  │
│  ┌─────────┬─────────────────┐  │
│  │ Image   │ Name    ↔       │   │  <- 已选卡片
│  │ 60x60   │ Intro ↔ 2 lines │ ✓ │     (选中图标 + 光晕)
│  └─────────┴─────────────────┘  │
│  ┌─────────┬─────────────────┐  │
│  │ Image   │ Name    ↔       │   │
│  │ 60x60   │ Intro ↔ 2 lines │   │
│  └─────────┴─────────────────┘  │
│                                 │  ▲ 滚动条 4px
├─────────────────────────────────┤  <- border-top
│ [Cancel] [Confirm] (Sticky)     │  <- 毛玻璃底栏
└─────────────────────────────────┘
```

## 7. 浏览器兼容性

| 功能              | Chrome | Firefox | Safari | Edge |
| ----------------- | ------ | ------- | ------ | ---- |
| box-shadow        | ✅      | ✅       | ✅      | ✅    |
| transform         | ✅      | ✅       | ✅      | ✅    |
| backdrop-filter   | ✅      | ❌*      | ✅      | ✅    |
| -webkit-scrollbar | ✅      | ❌**     | ✅      | ✅    |
| sticky            | ✅      | ✅       | ✅      | ✅    |
| line-clamp        | ✅      | ✅       | ✅      | ✅    |

注：
- *Firefox: 降级为半透明白色背景
- **Firefox: 使用原生滚动条

## 8. 性能考虑

- **动画**：使用 `transform` 和 `opacity`（性能最优）
- **滚动条**：仅 CSS 样式，无 JS 开销
- **Sticky**：浏览器原生支持，性能优秀
- **Blur**：移动端可降级或禁用（检查设备性能）

## 9. 测试覆盖

1. **视觉测试**：各态阴影、动画、光晕
2. **交互测试**：hover、选中、取消、滚动
3. **性能测试**：帧率、内存占用
4. **兼容性测试**：不同浏览器和分辨率
5. **端到端**：完整导出流程

---

**总体目标**：在不改变功能逻辑的基础上，通过精致的视觉设计和微交互，提升用户体验和产品品质。
