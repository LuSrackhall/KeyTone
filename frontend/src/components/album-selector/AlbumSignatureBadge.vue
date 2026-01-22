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
  AlbumSignatureBadge.vue - 专辑签名徽章组件

  功能说明：
    - 在专辑选择器中展示签名作者信息
    - 显示直接导出作者的头像（图片或签名图标占位符）和名称
    - 支持两种尺寸模式：normal（正常）和 small（紧凑）

  Props:
    - albumPath: 专辑路径，用于获取签名图片
    - authorName: 作者名称
    - authorImage: 作者图片路径（相对于专辑目录）
    - size: 尺寸模式，'normal' | 'small'

  使用示例：
    <AlbumSignatureBadge
      :album-path="albumPath"
      :author-name="signatureInfo.directExportAuthorName"
      :author-image="signatureInfo.directExportAuthorImage"
      size="small"
    />
  ============================================================================
-->

<template>
  <!-- 签名徽章容器 -->
  <div
    class="album-signature-badge"
    :class="[size === 'small' ? 'badge-small' : 'badge-normal', 'flex items-center gap-1.5', 'select-none']"
  >
    <!-- 作者头像 / 签名图标占位符 -->
    <div
      class="avatar-container flex-shrink-0 rounded-full overflow-hidden bg-gray-200/50"
      :class="[size === 'small' ? 'w-4 h-4' : 'w-5 h-5']"
    >
      <!-- 有图片时显示图片 -->
      <img v-if="imageUrl" :src="imageUrl" class="w-full h-full object-cover" @error="handleImageError" />
      <!-- 无图片时显示签名图标 -->
      <div
        v-else
        class="w-full h-full flex items-center justify-center"
        :class="[size === 'small' ? 'text-[10px]' : 'text-xs']"
      >
        <q-icon name="badge" color="amber-7" />
      </div>
    </div>

    <!-- 作者名称 -->
    <span
      class="author-name text-gray-600 truncate"
      :class="[size === 'small' ? 'text-[10px] max-w-[80px]' : 'text-xs max-w-[100px]']"
      :title="authorName"
    >
      {{ authorName }}
    </span>
  </div>
</template>

<script setup lang="ts">
/**
 * AlbumSignatureBadge - 专辑签名徽章组件
 *
 * 用于在专辑选择器的列表项和选中状态中展示签名作者信息
 */
import { ref, watch, onUnmounted } from 'vue';
import { GetAlbumFile } from 'src/boot/query/keytonePkg-query';

// ============================================================================
// Props 定义
// ============================================================================
const props = withDefaults(
  defineProps<{
    /** 专辑路径，用于通过 GetAlbumFile API 获取签名图片 */
    albumPath: string;
    /** 作者名称 */
    authorName: string;
    /** 作者图片路径（相对于专辑目录），如 "signature/card_xxx.jpg" */
    authorImage: string;
    /** 尺寸模式：'normal' 正常模式，'small' 紧凑模式（用于列表项） */
    size?: 'normal' | 'small';
  }>(),
  {
    size: 'normal',
  }
);

// ============================================================================
// 图片 URL 管理
// 通过 GetAlbumFile API 获取图片 Blob，创建 Object URL 用于显示
// 组件卸载时释放 Object URL 以防止内存泄漏
// ============================================================================
const imageUrl = ref<string | null>(null);
const imageLoadFailed = ref(false);

/**
 * 加载签名图片
 * 通过 GetAlbumFile API 从专辑目录读取图片文件
 */
async function loadImage() {
  // 重置状态
  imageUrl.value = null;
  imageLoadFailed.value = false;

  // 如果没有图片路径，不加载
  if (!props.authorImage || !props.albumPath) {
    return;
  }

  try {
    // 调用 API 获取图片 Blob
    const blob = await GetAlbumFile(props.albumPath, props.authorImage);
    if (blob) {
      // 释放旧的 Object URL（如果存在）
      if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value);
      }
      // 创建新的 Object URL
      imageUrl.value = URL.createObjectURL(blob);
    }
  } catch (error) {
    console.warn('AlbumSignatureBadge: 加载签名图片失败', error);
    imageLoadFailed.value = true;
  }
}

/**
 * 处理图片加载错误
 * 当 <img> 标签加载失败时调用
 */
function handleImageError() {
  imageLoadFailed.value = true;
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

// 组件卸载时释放 Object URL
onUnmounted(() => {
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
    imageUrl.value = null;
  }
});
</script>

<style scoped>
/* 签名徽章基础样式 */
.album-signature-badge {
  /* 使用 CSS 变量便于主题定制 */
  --badge-bg: rgba(251, 191, 36, 0.1); /* amber-400/10 */
  --badge-border: rgba(251, 191, 36, 0.3); /* amber-400/30 */
}

/* 正常尺寸模式 */
.badge-normal {
  padding: 2px 6px;
  background: var(--badge-bg);
  border: 1px solid var(--badge-border);
  border-radius: 4px;
}

/* 紧凑尺寸模式 */
.badge-small {
  padding: 1px 4px;
  background: var(--badge-bg);
  border: 1px solid var(--badge-border);
  border-radius: 3px;
}

/* 头像容器阴影效果 */
.avatar-container {
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 作者名称文本样式 */
.author-name {
  line-height: 1.2;
}
</style>
