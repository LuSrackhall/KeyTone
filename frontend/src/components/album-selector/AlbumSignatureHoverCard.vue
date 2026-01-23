<!--
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<!--
  ============================================================================
  AlbumSignatureHoverCard.vue - 专辑签名悬停详情卡片组件

  功能说明：
    - 包裹签名徽章组件，提供悬停时的详情卡片展示
    - 鼠标悬停在触发区域时显示详情卡片
    - 鼠标可移动到卡片上，卡片保持显示
    - 鼠标离开触发区域和卡片后，延迟隐藏卡片
    - 卡片右下角提供"点击查看详细信息"可点击标签
    - 点击标签打开 SignatureAuthorsDialog 对话框
    - 卡片具有毛玻璃效果（backdrop-filter: blur）

  展示内容：
    - 当 isSameAuthor=true 时，只展示一个作者区块（原始作者=直接导出作者）
    - 当 isSameAuthor=false 时，展示两个作者区块（原始作者 + 直接导出作者）

  交互逻辑：
    1. 鼠标进入触发区域 → 延迟显示卡片（避免快速划过时闪烁）
    2. 鼠标离开触发区域但进入卡片 → 卡片保持显示
    3. 鼠标离开卡片 → 检查是否在触发区域，若不在则延迟隐藏
    4. 点击"查看详细信息" → 触发 @view-details 事件

  Props:
    - albumPath: 专辑路径
    - signatureInfo: 签名摘要信息（包含原始作者和直接导出作者）

  Events:
    - view-details: 点击"查看详细信息"时触发

  使用示例：
    <AlbumSignatureHoverCard
      :album-path="albumPath"
      :signature-info="signatureInfo"
      @view-details="openSignatureDialog"
    >
      <AlbumSignatureBadge ... />
    </AlbumSignatureHoverCard>
  ============================================================================
-->

<template>
  <!-- 触发区域：包裹默认插槽内容 -->
  <div
    ref="triggerRef"
    class="hover-card-trigger inline-flex"
    @mouseenter="onTriggerEnter"
    @mouseleave="onTriggerLeave"
  >
    <slot />
  </div>

  <!-- ============================================================================
       悬停卡片：使用 Teleport 渲染到 body，避免被父元素 overflow 裁剪
       具有毛玻璃效果（backdrop-filter: blur(8px)）
       ============================================================================ -->
  <Teleport to="body">
    <Transition name="hover-card-fade">
      <div
        v-if="isCardVisible"
        ref="cardRef"
        class="signature-hover-card"
        :style="cardPosition"
        @mouseenter="onCardEnter"
        @mouseleave="onCardLeave"
        @pointerdown.stop
        @mousedown.stop
      >
        <!-- 卡片内容区域 -->
        <div class="card-content">
          <!-- ============================================================================
               作者信息展示
               - isSameAuthor=true: 只展示一个作者区块
               - isSameAuthor=false: 展示原始作者 + 直接导出作者
               ============================================================================ -->

          <!-- 情况1: 原始作者与直接导出作者是同一人 -->
          <template v-if="signatureInfo.isSameAuthor">
            <div class="author-section">
              <!-- 作者标签 -->
              <div class="author-label text-amber-600">
                {{ $t('albumSelector.signatureHoverCard.sameAuthor') }}
              </div>
              <!-- 作者信息 -->
              <div class="author-row">
                <div class="avatar-wrapper">
                  <img
                    v-if="directExportAuthorImageUrl"
                    :src="directExportAuthorImageUrl"
                    class="avatar-img"
                    @error="handleDirectExportImageError"
                  />
                  <div v-else class="avatar-placeholder">
                    <q-icon name="badge" color="amber-7" size="20px" />
                  </div>
                </div>
                <div class="author-info">
                  <div class="author-name" :title="signatureInfo.directExportAuthorName">
                    {{ signatureInfo.directExportAuthorName }}
                  </div>
                  <div
                    v-if="signatureInfo.directExportAuthorIntro"
                    class="author-intro"
                    :title="signatureInfo.directExportAuthorIntro"
                  >
                    {{ signatureInfo.directExportAuthorIntro }}
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- 情况2: 原始作者与直接导出作者不同 -->
          <template v-else>
            <!-- 原始作者区块 -->
            <div class="author-section">
              <div class="author-label text-amber-600">
                {{ $t('albumSelector.signatureHoverCard.originalAuthor') }}
              </div>
              <div class="author-row">
                <div class="avatar-wrapper">
                  <img
                    v-if="originalAuthorImageUrl"
                    :src="originalAuthorImageUrl"
                    class="avatar-img"
                    @error="handleOriginalImageError"
                  />
                  <div v-else class="avatar-placeholder">
                    <q-icon name="badge" color="amber-7" size="20px" />
                  </div>
                </div>
                <div class="author-info">
                  <div class="author-name" :title="signatureInfo.originalAuthorName">
                    {{ signatureInfo.originalAuthorName }}
                  </div>
                  <div
                    v-if="signatureInfo.originalAuthorIntro"
                    class="author-intro"
                    :title="signatureInfo.originalAuthorIntro"
                  >
                    {{ signatureInfo.originalAuthorIntro }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 分隔线 -->
            <div class="section-divider"></div>

            <!-- 直接导出作者区块 -->
            <div class="author-section">
              <div class="author-label text-blue-600">
                {{ $t('albumSelector.signatureHoverCard.directExportAuthor') }}
              </div>
              <div class="author-row">
                <div class="avatar-wrapper">
                  <img
                    v-if="directExportAuthorImageUrl"
                    :src="directExportAuthorImageUrl"
                    class="avatar-img"
                    @error="handleDirectExportImageError"
                  />
                  <div v-else class="avatar-placeholder">
                    <q-icon name="badge" color="blue-6" size="20px" />
                  </div>
                </div>
                <div class="author-info">
                  <div class="author-name" :title="signatureInfo.directExportAuthorName">
                    {{ signatureInfo.directExportAuthorName }}
                  </div>
                  <div
                    v-if="signatureInfo.directExportAuthorIntro"
                    class="author-intro"
                    :title="signatureInfo.directExportAuthorIntro"
                  >
                    {{ signatureInfo.directExportAuthorIntro }}
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- 分隔线 -->
          <div class="footer-divider"></div>

          <!-- 底部操作区域 -->
          <div class="card-footer">
            <span
              class="view-details-link"
              @pointerdown.capture.stop.prevent="handleViewDetails"
            >
              {{ $t('albumSelector.signatureHoverCard.viewDetails') }}
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
/**
 * AlbumSignatureHoverCard - 专辑签名悬停详情卡片组件
 *
 * 提供悬停时的签名详情预览，支持点击查看完整签名信息
 * 展示原始作者和直接导出作者的信息
 */
import { ref, watch, onUnmounted, nextTick } from 'vue';
import { GetAlbumFile } from 'src/boot/query/keytonePkg-query';
import type { AlbumSignatureSummary } from 'src/types/album-selector';

// ============================================================================
// Props 定义
// ============================================================================
const props = defineProps<{
  /** 专辑路径，用于通过 GetAlbumFile API 获取签名图片 */
  albumPath: string;
  /** 签名摘要信息（包含原始作者和直接导出作者的完整信息） */
  signatureInfo: AlbumSignatureSummary;
}>();

// ============================================================================
// Events 定义
// ============================================================================
const emit = defineEmits<{
  /** 点击"查看详细信息"时触发 */
  (e: 'view-details'): void;
}>();

// ============================================================================
// 悬停状态管理
// 使用标志位跟踪鼠标是否在触发区域或卡片上
// ============================================================================
const isHoverOnTrigger = ref(false);
const isHoverOnCard = ref(false);
const isCardVisible = ref(false);

// 延迟定时器引用
let showDelayTimer: ReturnType<typeof setTimeout> | null = null;
let hideDelayTimer: ReturnType<typeof setTimeout> | null = null;

// 延迟时间配置（毫秒）
const SHOW_DELAY = 200; // 显示延迟，避免快速划过时闪烁
const HIDE_DELAY = 100; // 隐藏延迟，确保鼠标能移动到卡片上

/**
 * 清除显示定时器
 */
function clearShowTimer() {
  if (showDelayTimer) {
    clearTimeout(showDelayTimer);
    showDelayTimer = null;
  }
}

/**
 * 清除隐藏定时器
 */
function clearHideTimer() {
  if (hideDelayTimer) {
    clearTimeout(hideDelayTimer);
    hideDelayTimer = null;
  }
}

/**
 * 检查是否应该隐藏卡片
 * 只有当鼠标既不在触发区域也不在卡片上时才隐藏
 */
function checkAndHideCard() {
  if (!isHoverOnTrigger.value && !isHoverOnCard.value) {
    isCardVisible.value = false;
  }
}

/**
 * 鼠标进入触发区域
 */
function onTriggerEnter() {
  isHoverOnTrigger.value = true;
  clearHideTimer(); // 取消隐藏计时

  // 延迟显示卡片
  clearShowTimer();
  showDelayTimer = setTimeout(() => {
    isCardVisible.value = true;
    // 显示后更新卡片位置
    nextTick(() => {
      updateCardPosition();
    });
  }, SHOW_DELAY);
}

/**
 * 鼠标离开触发区域
 */
function onTriggerLeave() {
  isHoverOnTrigger.value = false;
  clearShowTimer(); // 如果还在显示延迟中，取消显示

  // 延迟检查是否隐藏
  clearHideTimer();
  hideDelayTimer = setTimeout(() => {
    checkAndHideCard();
  }, HIDE_DELAY);
}

/**
 * 鼠标进入卡片
 */
function onCardEnter() {
  isHoverOnCard.value = true;
  clearHideTimer(); // 取消隐藏计时
}

/**
 * 鼠标离开卡片
 */
function onCardLeave() {
  isHoverOnCard.value = false;

  // 延迟检查是否隐藏
  clearHideTimer();
  hideDelayTimer = setTimeout(() => {
    checkAndHideCard();
  }, HIDE_DELAY);
}

/**
 * 点击"查看详细信息"
 *
 * 说明：
 * - 使用 pointerdown.capture 触发，确保在 QSelect 弹层监听到外部点击前就打开对话框
 * - 卡片根节点也会 stop pointerdown/mousedown，避免弹层把该点击视为"外部点击"
 * - 避免 click 阶段因列表项卸载而导致事件丢失
 */
function handleViewDetails() {
  // 立即隐藏卡片
  isHoverOnTrigger.value = false;
  isHoverOnCard.value = false;
  isCardVisible.value = false;
  clearShowTimer();
  clearHideTimer();
  // 触发事件
  emit('view-details');
}

// ============================================================================
// 卡片位置计算
// 根据触发区域的位置动态计算卡片的显示位置
// ============================================================================
const triggerRef = ref<HTMLElement | null>(null);
const cardRef = ref<HTMLElement | null>(null);
const cardPosition = ref<{ top: string; left: string }>({ top: '0px', left: '0px' });

/**
 * 更新卡片位置
 * 将卡片定位在触发区域的正下方或正上方（根据空间情况）
 */
function updateCardPosition() {
  if (!triggerRef.value) return;

  const triggerRect = triggerRef.value.getBoundingClientRect();
  const cardWidth = 280; // 卡片宽度
  const cardHeight = props.signatureInfo.isSameAuthor ? 140 : 220; // 根据内容估算高度
  const offset = 8; // 与触发区域的间距

  // 计算初始位置（触发区域下方，左对齐）
  let top = triggerRect.bottom + offset;
  let left = triggerRect.left;

  // 边界检查：防止卡片超出视口
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  // 右边界检查
  if (left + cardWidth > viewportWidth - 10) {
    left = viewportWidth - cardWidth - 10;
  }
  // 左边界检查
  if (left < 10) {
    left = 10;
  }
  // 下边界检查：如果下方空间不足，显示在上方
  if (top + cardHeight > viewportHeight - 10) {
    top = triggerRect.top - cardHeight - offset;
  }
  // 上边界检查
  if (top < 10) {
    top = 10;
  }

  cardPosition.value = {
    top: `${top}px`,
    left: `${left}px`,
  };
}

// ============================================================================
// 图片 URL 管理
// 分别管理原始作者和直接导出作者的图片
// ============================================================================
const originalAuthorImageUrl = ref<string | null>(null);
const directExportAuthorImageUrl = ref<string | null>(null);

/**
 * 加载原始作者图片
 */
async function loadOriginalAuthorImage() {
  if (!props.signatureInfo.originalAuthorImage || !props.albumPath) {
    originalAuthorImageUrl.value = null;
    return;
  }

  try {
    const blob = await GetAlbumFile(props.albumPath, props.signatureInfo.originalAuthorImage);
    if (blob) {
      if (originalAuthorImageUrl.value) {
        URL.revokeObjectURL(originalAuthorImageUrl.value);
      }
      originalAuthorImageUrl.value = URL.createObjectURL(blob);
    }
  } catch (error) {
    console.warn('AlbumSignatureHoverCard: 加载原始作者图片失败', error);
    originalAuthorImageUrl.value = null;
  }
}

/**
 * 加载直接导出作者图片
 */
async function loadDirectExportAuthorImage() {
  if (!props.signatureInfo.directExportAuthorImage || !props.albumPath) {
    directExportAuthorImageUrl.value = null;
    return;
  }

  try {
    const blob = await GetAlbumFile(props.albumPath, props.signatureInfo.directExportAuthorImage);
    if (blob) {
      if (directExportAuthorImageUrl.value) {
        URL.revokeObjectURL(directExportAuthorImageUrl.value);
      }
      directExportAuthorImageUrl.value = URL.createObjectURL(blob);
    }
  } catch (error) {
    console.warn('AlbumSignatureHoverCard: 加载直接导出作者图片失败', error);
    directExportAuthorImageUrl.value = null;
  }
}

/**
 * 处理原始作者图片加载错误
 */
function handleOriginalImageError() {
  if (originalAuthorImageUrl.value) {
    URL.revokeObjectURL(originalAuthorImageUrl.value);
    originalAuthorImageUrl.value = null;
  }
}

/**
 * 处理直接导出作者图片加载错误
 */
function handleDirectExportImageError() {
  if (directExportAuthorImageUrl.value) {
    URL.revokeObjectURL(directExportAuthorImageUrl.value);
    directExportAuthorImageUrl.value = null;
  }
}

// 监听 props 变化，重新加载图片
watch(
  () => [props.albumPath, props.signatureInfo],
  () => {
    loadOriginalAuthorImage();
    loadDirectExportAuthorImage();
  },
  { immediate: true, deep: true }
);

// 组件卸载时清理
onUnmounted(() => {
  clearShowTimer();
  clearHideTimer();
  if (originalAuthorImageUrl.value) {
    URL.revokeObjectURL(originalAuthorImageUrl.value);
  }
  if (directExportAuthorImageUrl.value) {
    URL.revokeObjectURL(directExportAuthorImageUrl.value);
  }
});
</script>

<style scoped>
/* ============================================================================
   悬停卡片主样式
   - 毛玻璃效果：backdrop-filter: blur(8px)
   - 半透明白色背景
   ============================================================================ */
.signature-hover-card {
  position: fixed;
  z-index: 9999;
  width: 280px;
  max-width: calc(100vw - 20px);
  border-radius: 12px;
  overflow: hidden;

  /* 毛玻璃效果：降低不透明度，确保下方内容可见且模糊 */
  background: rgba(255, 255, 255, 0.35);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);

  /* 轻微叠加层，增强可读性又不遮挡背景 */
  background-image: linear-gradient(
    180deg,
    rgba(255, 255, 255, 0.25) 0%,
    rgba(255, 255, 255, 0.15) 100%
  );

  /* 阴影和边框 */
  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.12),
    0 0 1px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

/* 卡片内容区域 */
.card-content {
  padding: 12px;
}

/* 作者区块 */
.author-section {
  margin-bottom: 8px;
}

/* 作者标签（原始作者/直接导出作者） */
.author-label {
  font-size: 10px;
  font-weight: 500;
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 作者信息行 */
.author-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

/* 头像包装器 */
.avatar-wrapper {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.05);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 头像图片 */
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 头像占位符 */
.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(251, 191, 36, 0.1);
}

/* 作者信息容器 */
.author-info {
  flex: 1;
  min-width: 0;
}

/* 作者名称 */
.author-name {
  font-size: 13px;
  font-weight: 500;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 作者介绍 */
.author-intro {
  font-size: 11px;
  color: #6b7280;
  margin-top: 2px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

/* 区块分隔线 */
.section-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.06);
  margin: 8px 0;
}

/* 底部分隔线 */
.footer-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.06);
  margin: 8px 0 6px 0;
}

/* 底部操作区域 */
.card-footer {
  display: flex;
  justify-content: flex-end;
}

/* 查看详细信息链接 */
.view-details-link {
  font-size: 11px;
  color: #3b82f6;
  cursor: pointer;
  transition: color 0.15s ease;
}

.view-details-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

/* ============================================================================
   卡片淡入淡出动画
   ============================================================================ */
.hover-card-fade-enter-active,
.hover-card-fade-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}

.hover-card-fade-enter-from,
.hover-card-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
