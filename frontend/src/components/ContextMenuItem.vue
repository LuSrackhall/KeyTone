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

<template>
  <!-- 菜单项容器 -->
  <div class="menu-item-wrapper">
    <!-- 菜单项内容 -->
    <div
      class="flex items-center gap-3 px-3 py-2 transition-colors duration-100 relative"
      :class="[
        disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
        isNegative ? 'text-negative hover:bg-red-50' : 'hover:bg-blue-50',
        loading ? 'opacity-75' : '',
      ]"
      :style="{ userSelect: 'none' }"
      @click.stop="handleClick"
      @mouseenter="isHovered = true"
      @mouseleave="isHovered = false"
    >
      <!-- 左侧指示条 - 悬停时显示 -->
      <div
        class="absolute left-0 top-0 bottom-0 w-1 rounded-r transition-colors duration-100"
        :class="[isHovered ? (isNegative ? 'bg-red-400' : 'bg-blue-400') : 'bg-transparent']"
      />

      <!-- 图标或加载转圈 -->
      <q-icon v-if="!loading" :name="icon" size="sm" />
      <q-spinner v-else size="20px" color="primary" />

      <!-- 文本内容 -->
      <span class="text-sm whitespace-nowrap">
        <slot />
      </span>
    </div>

    <!-- 菜单项分隔线 - 仅当不是最后一个菜单项时显示 -->
    <div v-if="showDivider" class="h-px bg-gradient-to-r from-gray-100 via-gray-150 to-gray-100" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  icon?: string;
  isNegative?: boolean;
  disabled?: boolean;
  loading?: boolean;
  showDivider?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'edit',
  isNegative: false,
  disabled: false,
  loading: false,
  showDivider: true,
});

const emit = defineEmits<{
  click: [];
}>();

// 悬停状态
const isHovered = ref(false);

function handleClick() {
  // 禁用或加载中时不响应点击
  if (props.disabled || props.loading) {
    return;
  }
  emit('click');
}
</script>

<style scoped>
/* 菜单项活跃反馈 */
.menu-item-wrapper > div:first-child:active:not(.opacity-50) {
  transform: scale(0.98);
}

/* 菜单项的左侧指示条添加发光效果 */
.menu-item-wrapper > div:first-child:hover .absolute {
  box-shadow: 0 0 8px rgba(59, 130, 246, 0.3);
}
</style>
