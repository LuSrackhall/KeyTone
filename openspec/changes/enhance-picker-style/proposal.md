# Proposal: 签名选择对话框样式与交互优化

## 概述

在完成数据集成后，进一步优化 `SignaturePickerDialog` 的视觉设计与交互体验，使其更加精致、高效、符合现代设计审美。

## 需求总结

1. **视觉层次感**：增强卡片设计的层次感，提升专业度
2. **内容区域优化**：名称/介绍支持横向滚动，间距合理，视觉紧凑
3. **交互灵活性**：支持点击已选项来取消选择；选中图标具有渐显透明效果
4. **底部操作栏固定**：固定粘滞在底部，毛玻璃背景，无需滚动即可操作
5. **滚动条优化**：自定义滚动条，更细致、更美观

## 范围

- **仅修改** `components/export-flow/SignaturePickerDialog.vue` 的样式（`<style scoped>`）和卡片交互逻辑
- **不修改** 数据加载、SSE 同步、搜索功能等核心逻辑
- **保持兼容** 所有事件和 props 接口不变

## 设计目标

### 1. 卡片视觉层次

当前：扁平的边框卡片
目标：具有微妙阴影、hover 效果、selected 高亮的分层卡片

```css
/* 正常态 */
box-shadow: 0 1px 3px rgba(0,0,0,0.05);

/* hover */
box-shadow: 0 4px 12px rgba(0,0,0,0.08);
transform: translateY(-1px);

/* selected */
border: 2px solid var(--q-primary);
box-shadow: 0 0 0 3px rgba(33,150,243,0.1), 0 4px 16px rgba(33,150,243,0.15);
```

### 2. 内容区域优化

参考签名管理页面的做法：
- 名称：单行，支持横向滚动
- 介绍：最多 2 行，支持横向滚动
- 间距：合理的上下间距（目标 4-8px）
- 字体大小统一（0.9rem 名称，0.75rem 介绍）
- 图片区间距：左右各 8px

```vue
<!-- 目标布局 -->
<div class="flex flex-col justify-center q-py-2xs q-px-sm" style="min-width: 0">
  <!-- 名称：支持横向滚动 -->
  <div class="scrollable-x">{{ sig.name }}</div>
  <!-- 介绍：支持横向滚动，最多 2 行 -->
  <div class="scrollable-x line-clamp-2">{{ sig.intro }}</div>
</div>
```

### 3. 选择交互

当前：选中后显示勾选图标，无法取消
目标：
- 再次点击已选项可取消选择
- 选中图标周围有渐显透明的光晕效果（radial-gradient，动画周期 1.5s）
- 光晕采用脉冲动画（scale: 1 → 1.1 → 1）

```vue
<!-- 选中指示器 -->
<div v-if="selectedId === sig.id" class="selection-indicator-wrapper">
  <div class="selection-glow"></div>
  <q-icon name="check_circle" class="selection-icon" />
</div>
```

### 4. 底部操作栏固定

当前：操作按钮随内容滚动
目标：
- 使用 `position: sticky; bottom: 0` 固定在底部
- 背景使用毛玻璃效果（backdrop-filter: blur(10px); rgba(255,255,255,0.1)）
- 上方边界线以区分内容区
- z-index 确保始终在内容之上

```css
.sticky-bottom {
  position: sticky;
  bottom: 0;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.1);
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  z-index: 10;
}
```

### 5. 主对话框滚动条隐藏

当前：显示原生滚动条，占用空间
目标：
- 使用 `[&::-webkit-scrollbar]:hidden` 隐藏主滚动条
- 保留滚动功能（鼠标滚轮、触控板）
- 内容卡片内部的横向滚动条仍保持可见和美观

```css
.signature-picker-dialog {
  /* 隐藏主滚动条 */
  [&::-webkit-scrollbar]: hidden;
}
```

### 6. 内容区域间距优化

当前：间距过大（padding-bottom: 100px）
目标：
- 搜索区域 padding 从 q-pa-md 改为 q-pa-sm
- 描述和搜索下方 margin 从 q-mb-md 改为 q-mb-sm
- 卡片间距从 q-mb-sm 改为 q-mb-xs
- 移除过度 padding-bottom，改为自适应

## 实现步骤

### Step 1：增强卡片阴影与 hover 效果

- 更新 `.signature-card` 的 `box-shadow` 和 `transform` 属性
- 调整 selected 状态的视觉反馈

### Step 2：优化内容区域布局

- 调整搜索区域 padding：`q-pa-md` → `q-pa-sm`
- 调整描述和搜索下方 margin：`q-mb-md` → `q-mb-sm`
- 调整卡片间距：`q-mb-sm` → `q-mb-xs`
- 移除主内容区域硬编码的 `padding-bottom: 100px`

### Step 3：实现选择取消逻辑

- 修改 `selectSignature()` 函数，支持 toggle（重复点击取消）
- 为选中图标添加渐显光晕 CSS 动画

### Step 4：固定底部操作栏

- 将 `<q-card-actions>` 改为 sticky 定位
- 添加毛玻璃背景和边界线

### Step 5：隐藏主对话框滚动条

- 应用 `[&::-webkit-scrollbar]:hidden` 到对话框卡片
- 保证滚动功能正常

## 测试计划

1. **卡片视觉**：hover 和 selected 状态的阴影、圆角、颜色
2. **内容滚动**：名称/介绍过长时，是否支持横向滚动
3. **选择交互**：点击已选项是否取消，图标光晕是否正常显示
4. **底部固定**：滚动内容时，底部按钮是否保持可见和可点击
5. **主滚动条**：主对话框滚动条是否隐藏，滚动是否正常
6. **间距美观**：整体布局是否更紧凑、更美观
7. **端到端**：完整的导出流程是否依然正常

## 相关 Spec

- `spec/album-export` - 签名选择对话框交互（将新增选择取消需求）

## 风险与缓解

| 风险                          | 缓解方案                           |
| ----------------------------- | ---------------------------------- |
| 毛玻璃效果浏览器兼容性        | 降级：使用半透明白色背景           |
| 滚动条自定义在 Firefox 不支持 | 接受 Firefox 使用原生滚动条        |
| 固定底部导致内容被遮挡        | 内容区添加 padding-bottom 预留空间 |

## 合并条件

1. ✅ 所有视觉效果通过视觉审查
2. ✅ 交互流程验证无误
3. ✅ 没有破坏导出流程的其他功能
4. ✅ 不同分辨率和浏览器兼容性测试通过

## 后续迭代方向

- 考虑签名预览卡片的更多交互（如拖拽排序、快速操作菜单）
- 对话框自适应宽度（当屏幕宽度较小时）
- 暗色主题支持
