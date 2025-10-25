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
  <!-- 背景遮罩 - 用于点击外部关闭菜单 -->
  <div
    v-if="isVisible"
    class="fixed inset-0 z-40"
    style="background-color: transparent"
    @click="close"
    @contextmenu.prevent="close"
  />

  <!-- 自定义菜单容器 -->
  <div
    v-if="isVisible"
    :style="menuStyle"
    class="fixed bg-white rounded-lg z-50 overflow-hidden"
    style="box-shadow: 0 10px 40px rgba(0, 0, 0, 0.16), 0 0 0 1px rgba(0, 0, 0, 0.08); backdrop-filter: blur(4px)"
    @click.stop
    @contextmenu.stop.prevent
  >
    <!-- 菜单内容区域 -->
    <div class="min-w-max py-1">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';

interface Props {
  modelValue: boolean;
  x?: number;
  y?: number;
  maxWidth?: number;
}

// ========== 常量定义 ==========
const DEFAULT_MENU_WIDTH = 150;
const DEFAULT_MENU_HEIGHT = 120;
const MENU_OFFSET = 4; // 距离边界的最小距离

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  x: 0,
  y: 0,
  maxWidth: DEFAULT_MENU_WIDTH,
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  show: [];
  hide: [];
}>();

// ========== 状态管理 ==========
const isVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val);
    if (val) {
      emit('show');
    } else {
      emit('hide');
    }
  },
});

// 菜单计算后的实际位置
const menuX = ref(props.x);
const menuY = ref(props.y);

// 事件监听状态标记
let resizeListenerActive = false;
let keydownListenerActive = false;

// ========== 计算属性 ==========
const menuStyle = computed(() => {
  return {
    left: menuX.value + 'px',
    top: menuY.value + 'px',
    maxWidth: props.maxWidth + 'px',
    minWidth: props.maxWidth + 'px',
  };
});

// ========== 方法 ==========

/** 关闭菜单 */
function close() {
  isVisible.value = false;
}

/** 更新菜单位置，确保不超出视口边界 */
const updateMenuPosition = () => {
  if (!isVisible.value) return;

  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;
  const menuWidth = props.maxWidth;
  const menuHeight = DEFAULT_MENU_HEIGHT;

  let x = props.x;
  let y = props.y;

  // 检查是否超出右边界
  if (x + menuWidth > viewportWidth) {
    x = Math.max(0, viewportWidth - menuWidth - MENU_OFFSET);
  }

  // 检查是否超出下边界
  if (y + menuHeight > viewportHeight) {
    y = Math.max(0, viewportHeight - menuHeight - MENU_OFFSET);
  }

  menuX.value = x;
  menuY.value = y;
};

/** 处理菜单可见性变化 */
const handleVisibilityChange = () => {
  if (isVisible.value && !resizeListenerActive) {
    // 菜单打开：添加 resize 监听
    setTimeout(updateMenuPosition, 0);
    window.addEventListener('resize', updateMenuPosition);
    resizeListenerActive = true;
  } else if (!isVisible.value && resizeListenerActive) {
    // 菜单关闭：移除 resize 监听
    window.removeEventListener('resize', updateMenuPosition);
    resizeListenerActive = false;
  }
};

/** 处理键盘事件 */
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isVisible.value) {
    close();
  }
};

// ========== 生命周期 ==========

onMounted(() => {
  // 初始化
  handleVisibilityChange();

  // 添加键盘事件监听（只添加一次）
  if (!keydownListenerActive) {
    document.addEventListener('keydown', handleKeyDown);
    keydownListenerActive = true;
  }
});

onUnmounted(() => {
  // 清理 resize 监听
  if (resizeListenerActive) {
    window.removeEventListener('resize', updateMenuPosition);
    resizeListenerActive = false;
  }

  // 清理键盘事件监听
  if (keydownListenerActive) {
    document.removeEventListener('keydown', handleKeyDown);
    keydownListenerActive = false;
  }
});

// 监听 modelValue 变化
watch(isVisible, handleVisibilityChange);

// 监听坐标变化
watch(
  () => ({ x: props.x, y: props.y }),
  () => {
    menuX.value = props.x;
    menuY.value = props.y;
    updateMenuPosition();
  }
);
</script>

<style scoped>
/* 菜单项的基础样式由父组件定义 */
</style>
