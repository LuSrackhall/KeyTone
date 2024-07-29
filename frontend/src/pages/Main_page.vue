<template>
  <!-- <div>Silent Typist's Friend</div> -->
  <!--
   q-layout 有个 style="min-height: 803.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->
  <q-page style="min-height: 0px">
    <div
      :class="[
        '',
        'w-auto h-60 flex flex-col items-center',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <q-avatar
        :class="[
          '',
          'm-t-10',
          // 设置整体头像的尺寸大小。 至于内部相关文字的尺寸大小(如果有的话  可通过设置text-[8rem] leading-[10rem]实现)
          'size-56',
          '',
        ]"
      >
        <!-- draggable="false"使得图片无法拖动, 免得影响界面体验 -->
        <img :src="logoUrl" draggable="false" />
      </q-avatar>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import logoUrl from 'assets/img/KeyTone.png?url';
</script>

<style lang="scss" scoped></style>
