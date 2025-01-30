<template>
  <div class="w-[379px] h-[458.5px] overflow-hidden">
    <div :class="['w-[88%] flex flex-col ml-[6.2%] mr-[5.8%] pt-[5%]']">
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
    </div>

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
</script>

<style lang="scss" scoped></style>
