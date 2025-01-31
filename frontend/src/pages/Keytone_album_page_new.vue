<template>
  <div class="w-[379px] h-[458.5px] overflow-hidden">
    <!-- 展开按钮 -->
    <div v-if="isCollapsed" class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50">
      <q-btn flat class="custom-expand-btn" @click="isCollapsed = false">
        <div class="chevron-down"></div>
      </q-btn>
    </div>

    <!-- 选择器容器使用绝对定位 -->
    <div class="relative">
      <transition name="slide">
        <div
          v-show="!isCollapsed"
          class="selector-container absolute w-[88%] ml-[6.2%] mr-[5.8%] pt-[5%]"
          style="top: 0; left: 0; right: 0; z-index: 1"
        >
          <q-select
            v-model="setting_store.mainHome.selectedKeyTonePkg"
            :options="main_store.keyTonePkgOptions"
            :option-label="(item: any) => main_store.keyTonePkgOptionsName.get(item)"
            :label="$t('mainHome.selectedKeySoundAlbum')"
            :virtual-scroll-slice-size="999999"
            outlined
            dense
            emit-value
            map-options
            ref="selectedKeyTonePkgRef"
            @popup-hide="blur()"
          >
            <template v-if="setting_store.mainHome.selectedKeyTonePkg" v-slot:append>
              <q-icon
                name="cancel"
                @click.stop.prevent="setting_store.mainHome.selectedKeyTonePkg = ''"
                class="cursor-pointer text-lg"
              />
            </template>
          </q-select>

          <q-btn
            flat
            round
            color="grey"
            icon="expand_less"
            class="collapse-btn absolute -bottom-6 left-1/2 transform -translate-x-1/2"
            @click="isCollapsed = true"
          />
        </div>
      </transition>

      <!-- 内容区域添加过渡padding -->
      <div
        class="content-wrapper"
        :style="{
          paddingTop: !isCollapsed ? '60px' : '0',
          transition: `padding-top ${!isCollapsed ? '0.8s' : '1.2s'} ease`,
        }"
      >
        <div :class="{ 'hide-scrollbar': isAtTop }" class="keytone-album-container">
          <KeytoneAlbum
            v-if="setting_store.mainHome.selectedKeyTonePkg"
            :key="setting_store.mainHome.selectedKeyTonePkg"
            ref="keytoneAlbumRef"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { QSelect } from 'quasar';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { useTemplateRef } from 'vue';
import KeytoneAlbum from 'src/components/Keytone_album.vue';
import { ref, onMounted, onUnmounted } from 'vue';

const main_store = useMainStore();
const setting_store = useSettingStore();
const selectedKeyTonePkgRef = useTemplateRef<QSelect>('selectedKeyTonePkgRef');
const keytoneAlbumRef = ref<InstanceType<typeof KeytoneAlbum> | null>(null);
const isCollapsed = ref(false);
let lastScrollTop = 0;
const isAtTop = ref(true);

const blur = () => {
  setTimeout(() => {
    selectedKeyTonePkgRef?.value?.blur();
    // TIPS: 这里需要延迟后再blur, 以确保blur的正确触发(太早触发blur会不起作用, 经验证, 本人电脑延迟10ms后, 可以正确触发blur使焦点丧失, 为确保适配更多的低性能设备, 这里保险起见设置为66ms)
  }, 66);
};

console.log('main_store.keyTonePkgOptions', main_store.keyTonePkgOptions);

// 监听 KeytoneAlbum 内部的滚动
const handleAlbumScroll = (event: Event) => {
  const scrollableElement = (event.target as HTMLElement).closest('.q-scrollarea__container');
  if (!scrollableElement) return;

  const currentScroll = scrollableElement.scrollTop;
  const maxScroll = scrollableElement.scrollHeight - scrollableElement.clientHeight;

  // 更新是否在顶部的状态
  isAtTop.value = currentScroll === 0;

  // 在顶部继续向上滚动时展开
  if (currentScroll === 0 && event instanceof WheelEvent && event.deltaY < 0 && isCollapsed.value) {
    isCollapsed.value = false;
    return;
  }

  // 向下滚动时立即收起
  if (
    (currentScroll > lastScrollTop ||
      (currentScroll >= maxScroll && event instanceof WheelEvent && event.deltaY > 0)) &&
    !isCollapsed.value
  ) {
    isCollapsed.value = true;
  }

  lastScrollTop = currentScroll;
};

onMounted(() => {
  const scrollContainer = keytoneAlbumRef.value?.$el.querySelector('.q-scrollarea__container');
  if (scrollContainer) {
    scrollContainer.addEventListener('scroll', handleAlbumScroll, { passive: true });
    scrollContainer.addEventListener('wheel', handleAlbumScroll, { passive: true });
  }
});

onUnmounted(() => {
  const scrollContainer = keytoneAlbumRef.value?.$el.querySelector('.q-scrollarea__container');
  if (scrollContainer) {
    scrollContainer.removeEventListener('scroll', handleAlbumScroll);
    scrollContainer.removeEventListener('wheel', handleAlbumScroll);
  }
});
</script>

<style lang="scss" scoped>
.selector-container {
  transform-origin: top;
}

// 滑动动画
.slide-enter-active,
.slide-leave-active {
  transition: all 0.8s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
  transform: translateY(0);
  opacity: 1;
}

.collapse-btn {
  position: absolute;
  z-index: 1;
  min-height: 24px;
  backdrop-filter: blur(4px);
  border-radius: 12px;
}

.custom-expand-btn {
  min-height: 24px;
  width: 64px;
  padding: 0;
  background: rgba(100, 100, 100, 0.3);
  backdrop-filter: blur(4px);
  border-radius: 12px;
  margin-top: 8px;

  &:hover {
    background: rgba(100, 100, 100, 0.3);
  }
}

.chevron-down {
  width: 10px;
  height: 10px;
  border-right: 2px solid white;
  border-bottom: 2px solid white;
  transform: rotate(45deg);
  margin-top: -4px;
}

.content-wrapper {
  will-change: padding-top;
  position: relative;
}

.keytone-album-container {
  &.hide-scrollbar :deep(.q-scrollarea__thumb) {
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  :deep(.q-scrollarea__thumb) {
    transition: opacity 0.3s ease;
  }
}
</style>
