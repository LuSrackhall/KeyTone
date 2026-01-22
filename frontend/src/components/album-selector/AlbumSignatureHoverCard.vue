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
    - 鼠标可移动到卡片上，卡片保持显示不消失
    - 卡片右下角提供"点击查看详细信息"可点击标签
    - 点击标签打开 SignatureAuthorsDialog 对话框

  交互逻辑：
    1. 鼠标进入触发区域 → 延迟显示卡片（避免快速划过时闪烁）
    2. 鼠标离开触发区域但进入卡片 → 卡片保持显示
    3. 鼠标离开卡片且不在触发区域 → 延迟隐藏卡片
    4. 点击"查看详细信息" → 触发 @view-details 事件

  Props:
    - albumPath: 专辑路径
    - authorName: 作者名称
    - authorImage: 作者图片路径
    - authorIntro: 作者介绍（可选）

  Events:
    - view-details: 点击"查看详细信息"时触发

  使用示例：
    <AlbumSignatureHoverCard
      :album-path="albumPath"
      :author-name="signatureInfo.directExportAuthorName"
      :author-image="signatureInfo.directExportAuthorImage"
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

  <!-- 悬停卡片：使用 Teleport 渲染到 body，避免被父元素 overflow 裁剪 -->
  <Teleport to="body">
    <Transition name="hover-card-fade">
      <div
        v-if="showCard"
        ref="cardRef"
        class="signature-hover-card fixed z-[9999] bg-white rounded-lg shadow-lg border border-gray-200"
        :style="cardPosition"
        @mouseenter="onCardEnter"
        @mouseleave="onCardLeave"
      >
        <!-- 卡片内容区域 -->
        <div class="card-content p-3">
          <!-- 作者信息头部 -->
          <div class="flex items-start gap-3">
            <!-- 作者头像 -->
            <div class="avatar-wrapper flex-shrink-0 w-12 h-12 rounded-full overflow-hidden bg-gray-100 shadow-sm">
              <img v-if="imageUrl" :src="imageUrl" class="w-full h-full object-cover" @error="handleImageError" />
              <div v-else class="w-full h-full flex items-center justify-center">
                <q-icon name="badge" color="amber-7" size="24px" />
              </div>
            </div>

            <!-- 作者信息 -->
            <div class="author-info flex-1 min-w-0">
              <!-- 作者名称 -->
              <div class="author-name text-sm font-medium text-gray-900 truncate" :title="authorName">
                {{ authorName }}
              </div>
              <!-- 直接导出作者标签 -->
              <div class="author-role text-[10px] text-amber-600 mt-0.5">
                {{ $t('albumSelector.signatureHoverCard.directExportAuthor') }}
              </div>
              <!-- 作者介绍（如有） -->
              <div v-if="authorIntro" class="author-intro text-xs text-gray-500 mt-1 line-clamp-2" :title="authorIntro">
                {{ authorIntro }}
              </div>
            </div>
          </div>

          <!-- 分隔线 -->
          <div class="divider my-2 border-t border-gray-100"></div>

          <!-- 底部操作区域 -->
          <div class="card-footer flex justify-end">
            <!-- 查看详细信息链接 -->
            <span
              class="view-details-link text-xs text-blue-500 hover:text-blue-700 cursor-pointer hover:underline transition-colors"
              @click="handleViewDetails"
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
 */
import { ref, computed, watch, onUnmounted, nextTick } from 'vue';
import { GetAlbumFile } from 'src/boot/query/keytonePkg-query';

// ============================================================================
// Props 定义
// ============================================================================
const props = withDefaults(
  defineProps<{
    /** 专辑路径，用于获取签名图片 */
    albumPath: string;
    /** 作者名称 */
    authorName: string;
    /** 作者图片路径（相对于专辑目录） */
    authorImage: string;
    /** 作者介绍（可选） */
    authorIntro?: string;
  }>(),
  {
    authorIntro: '',
  }
);

// ============================================================================
// Events 定义
// ============================================================================
const emit = defineEmits<{
  /** 点击"查看详细信息"时触发 */
  (e: 'view-details'): void;
}>();

// ============================================================================
// 悬停状态管理
// 使用两个独立的状态变量跟踪鼠标是否在触发区域或卡片上
// ============================================================================
const isHoverOnTrigger = ref(false);
const isHoverOnCard = ref(false);

/** 控制卡片显示的计算属性 */
const showCard = computed(() => isHoverOnTrigger.value || isHoverOnCard.value);

// 延迟定时器引用
let showDelayTimer: ReturnType<typeof setTimeout> | null = null;
let hideDelayTimer: ReturnType<typeof setTimeout> | null = null;

// 延迟时间配置（毫秒）
const SHOW_DELAY = 200; // 显示延迟，避免快速划过时闪烁
const HIDE_DELAY = 150; // 隐藏延迟，确保鼠标能移动到卡片上

/**
 * 清除所有延迟定时器
 */
function clearTimers() {
  if (showDelayTimer) {
    clearTimeout(showDelayTimer);
    showDelayTimer = null;
  }
  if (hideDelayTimer) {
    clearTimeout(hideDelayTimer);
    hideDelayTimer = null;
  }
}

/**
 * 鼠标进入触发区域
 */
function onTriggerEnter() {
  clearTimers();
  // 延迟显示卡片
  showDelayTimer = setTimeout(() => {
    isHoverOnTrigger.value = true;
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
  clearTimers();
  // 延迟隐藏，给鼠标移动到卡片的时间
  hideDelayTimer = setTimeout(() => {
    isHoverOnTrigger.value = false;
  }, HIDE_DELAY);
}

/**
 * 鼠标进入卡片
 */
function onCardEnter() {
  clearTimers();
  isHoverOnCard.value = true;
}

/**
 * 鼠标离开卡片
 */
function onCardLeave() {
  clearTimers();
  // 延迟隐藏
  hideDelayTimer = setTimeout(() => {
    isHoverOnCard.value = false;
  }, HIDE_DELAY);
}

/**
 * 点击"查看详细信息"
 */
function handleViewDetails() {
  // 立即隐藏卡片
  isHoverOnTrigger.value = false;
  isHoverOnCard.value = false;
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
 * 将卡片定位在触发区域的正下方
 */
function updateCardPosition() {
  if (!triggerRef.value) return;

  const triggerRect = triggerRef.value.getBoundingClientRect();
  const cardWidth = 240; // 卡片预估宽度
  const cardHeight = 150; // 卡片预估高度
  const offset = 8; // 与触发区域的间距

  // 计算初始位置（触发区域下方居中）
  let top = triggerRect.bottom + offset;
  let left = triggerRect.left + (triggerRect.width - cardWidth) / 2;

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

  cardPosition.value = {
    top: `${top}px`,
    left: `${left}px`,
  };
}

// ============================================================================
// 图片 URL 管理
// ============================================================================
const imageUrl = ref<string | null>(null);

/**
 * 加载签名图片
 */
async function loadImage() {
  imageUrl.value = null;

  if (!props.authorImage || !props.albumPath) {
    return;
  }

  try {
    const blob = await GetAlbumFile(props.albumPath, props.authorImage);
    if (blob) {
      if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value);
      }
      imageUrl.value = URL.createObjectURL(blob);
    }
  } catch (error) {
    console.warn('AlbumSignatureHoverCard: 加载签名图片失败', error);
  }
}

/**
 * 处理图片加载错误
 */
function handleImageError() {
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
    imageUrl.value = null;
  }
}

// 监听 props 变化，重新加载图片
watch(
  () => [props.albumPath, props.authorImage],
  () => {
    loadImage();
  },
  { immediate: true }
);

// 组件卸载时清理
onUnmounted(() => {
  clearTimers();
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
  }
});
</script>

<style scoped>
/* 悬停卡片样式 */
.signature-hover-card {
  width: 240px;
  max-width: calc(100vw - 20px);
}

/* 卡片淡入淡出动画 */
.hover-card-fade-enter-active,
.hover-card-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.hover-card-fade-enter-from,
.hover-card-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* 作者介绍文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 头像包装器阴影 */
.avatar-wrapper {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
