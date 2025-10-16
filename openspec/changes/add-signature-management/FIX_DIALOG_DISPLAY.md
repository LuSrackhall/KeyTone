# 签名对话框显示问题修复

## 问题描述

使用 `q-scroll-area` 配合 Flexbox 的方式导致对话框内容无法正常显示。

## 根本原因

Quasar 的 `q-scroll-area` 组件在 Flexbox 容器中的高度计算存在问题：
- 当 `q-scroll-area` 直接作为 flex 子元素时，其 `flex: 1` 和内部的高度无法正确协商
- 即使设置 `height: 100%`，由于父容器高度不确定，也无法工作

## 解决方案

使用中间容器包装 `q-scroll-area`：

```vue
<!-- ✅ 正确做法 -->
<div style="flex: 1; overflow: hidden">
  <q-scroll-area style="height: 100%">
    <!-- 内容 -->
  </q-scroll-area>
</div>

<!-- ❌ 错误做法（之前的方案）-->
<q-scroll-area style="flex: 1; ...">
  <!-- 内容 -->
</q-scroll-area>
```

## 修改内容

### `SignatureFormDialog.vue` 模板结构

**关键变更**：

1. **q-card 设置 Flexbox**：
   ```vue
   <q-card class="w-96" style="display: flex; flex-direction: column; max-height: 90vh">
   ```

2. **中间容器包装 q-scroll-area**：
   ```vue
   <div style="flex: 1; overflow: hidden">
     <q-scroll-area style="height: 100%">
       <!-- 内容 -->
     </q-scroll-area>
   </div>
   ```

3. **底部按钮固定**：
   ```vue
   <q-card-actions align="right" style="flex-shrink: 0; border-top: 1px solid rgba(0, 0, 0, 0.12); background: white">
   ```

4. **CSS 简化**：
   - 移除了不必要的 `.dialog-card` 和 `.scrollable-content` CSS 类
   - 改为使用 inline style 直接控制布局

## 效果验证

✅ **对话框内容正常显示**
- 标题和编辑提示固定在顶部
- 表单内容区域可滚动
- 底部按钮始终固定在对话框底部
- 无论内容多少，按钮都可见

✅ **编译通过**
- No errors found

## 技术细节

**为什么需要中间容器**：

Flexbox 布局中，直接子元素的高度计算基于其内容和 flex 属性。但 `q-scroll-area` 是一个复杂的组件，它的内部结构与标准 HTML 元素不同。通过添加一个中间的 `<div>` 容器：

1. 中间容器获得 `flex: 1`，占用所有剩余空间
2. 中间容器设置 `overflow: hidden` 防止溢出
3. `q-scroll-area` 设置 `height: 100%` 填充中间容器
4. 现在 `q-scroll-area` 有了明确的高度值，可以正确工作

这是处理 Quasar 组件与 Flexbox 交互的标准模式。

