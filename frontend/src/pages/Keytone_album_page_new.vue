<template>
  <div class="w-[379px] h-[458.5px] overflow-hidden">
    <!-- 展开按钮 -->
    <div v-if="isCollapsed" class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50">
      <q-btn flat class="custom-expand-btn" @click="isCollapsed = false">
        <div class="chevron-down"></div>
      </q-btn>
    </div>

    <!-- 使用 Vue 的过渡组件 -->
    <transition name="slide">
      <div v-show="!isCollapsed" class="selector-container w-[88%] ml-[6.2%] mr-[5.8%] pt-[5%] relative">
        <q-select
          v-model="setting_store.mainHome.selectedKeyTonePkg"
          :options="main_store.keyTonePkgOptions"
          :option-label="(item: any) => {
            return main_store.keyTonePkgOptionsName.get(item)
          }"
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
            <!-- 由于直接使用默认的clearable, 会使得mode=null, 而我希望点击清楚按钮时mode=""即空字符串。因此使用插槽来实现。 -->
            <q-icon
              name="cancel"
              @click.stop.prevent="setting_store.mainHome.selectedKeyTonePkg = ''"
              class="cursor-pointer text-lg"
            />
          </template>
        </q-select>

        <!-- 收起按钮 -->
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

    <div>
      <!-- 引入键音包编辑组件 (v-if保证了, 当键音包未被选中时, 不会引入键音包编辑组件。 :key保证了, 当键音包被切换时, 会重新引入键音包编辑组件。)-->
      <KeytoneAlbum v-if="setting_store.mainHome.selectedKeyTonePkg" :key="setting_store.mainHome.selectedKeyTonePkg" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { QSelect } from 'quasar';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { useTemplateRef } from 'vue';
import KeytoneAlbum from 'src/components/Keytone_album.vue';
import { ref } from 'vue';

const main_store = useMainStore();
const setting_store = useSettingStore();

const selectedKeyTonePkgRef = useTemplateRef<QSelect>('selectedKeyTonePkgRef');
const blur = () => {
  setTimeout(() => {
    selectedKeyTonePkgRef?.value?.blur();
    // TIPS: 这里需要延迟后再blur, 以确保blur的正确触发(太早触发blur会不起作用, 经验证, 本人电脑延迟10ms后, 可以正确触发blur使焦点丧失, 为确保适配更多的低性能设备, 这里保险起见设置为66ms)
  }, 66);
};

console.log('main_store.keyTonePkgOptions', main_store.keyTonePkgOptions);

// 确保 isCollapsed 的初始值为 false
const isCollapsed = ref(false);
</script>

<style lang="scss" scoped>
.selector-container {
  transform-origin: top;
}

// 滑动动画
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
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
</style>
