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
  AlbumSignatureBadge.vue - 专辑签名徽章组件（芯片样式）

  功能说明：
    - 在专辑选择器中展示签名作者信息
    - 显示直接导出作者的头像（图片或签名图标占位符）和名称
    - 采用芯片（Chip）样式，宽度根据内容自适应，不占整行
    - 支持两种尺寸模式：normal（正常）和 small（紧凑）

  设计特点：
    - 内联显示（inline-flex），宽度由内容决定
    - 圆角胶囊形状，视觉上类似芯片/标签
    - 背景采用琥珀色调，与签名主题一致

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
  <!-- ============================================================================
       签名徽章容器（芯片样式）
       - 使用 inline-flex 确保宽度由内容决定
       - 圆角胶囊形状
       ============================================================================ -->
  <div
    class="album-signature-badge"
    :class="[
      size === 'small' ? 'badge-small' : 'badge-normal',
    ]"
  >
    <!-- 作者头像 / 签名图标占位符 -->
    <div class="avatar-container" :class="[size === 'small' ? 'avatar-small' : 'avatar-normal']">
      <!-- 有图片时显示图片 -->
      <img
        v-if="imageUrl"
        :src="imageUrl"
        class="avatar-img"
        @error="handleImageError"
      />
      <!-- 无图片时显示签名图标 -->
      <div v-else class="avatar-placeholder">
        <q-icon name="badge" color="amber-7" :size="size === 'small' ? '10px' : '12px'" />
      </div>
    </div>

    <!-- 作者名称 -->
    <span
      class="author-name"
      :class="[size === 'small' ? 'name-small' : 'name-normal']"
      :title="authorName"
    >
      {{ authorName }}
    </span>
  </div>
</template>

<script setup lang="ts">
/**
 * AlbumSignatureBadge - 专辑签名徽章组件（芯片样式）
 *
 * 用于在专辑选择器的列表项和选中状态中展示签名作者信息
 * 采用芯片样式，宽度自适应，视觉效果更紧凑美观
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
/* ============================================================================
   签名徽章基础样式（芯片/Chip 样式）
   - inline-flex: 宽度由内容决定，不占整行
   - 圆角胶囊形状
   - 琥珀色背景
   ============================================================================ */
.album-signature-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(251, 191, 36, 0.12);
  border: 1px solid rgba(251, 191, 36, 0.25);
  border-radius: 999px; /* 胶囊形状 */
  user-select: none;
  white-space: nowrap;
  max-width: 100%;
}

/* ============================================================================
   正常尺寸模式
   ============================================================================ */
.badge-normal {
  padding: 2px 8px 2px 3px;
}

/* ============================================================================
   紧凑尺寸模式（用于列表项）
   ============================================================================ */
.badge-small {
  padding: 1px 6px 1px 2px;
}

/* ============================================================================
   头像容器
   ============================================================================ */
.avatar-container {
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  background: rgba(251, 191, 36, 0.15);
}

.avatar-normal {
  width: 18px;
  height: 18px;
}

.avatar-small {
  width: 14px;
  height: 14px;
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
}

/* ============================================================================
   作者名称
   ============================================================================ */
.author-name {
  color: #92400e; /* amber-800 */
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
}

.name-normal {
  font-size: 11px;
  max-width: 100px;
}

.name-small {
  font-size: 10px;
  max-width: 80px;
}
</style>
