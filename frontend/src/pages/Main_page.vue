<template>
  <!-- <div>Silent Typist's Friend</div> -->
  <q-page>
    <div :class="['', 'w-auto h-60 flex flex-col items-center']">
      <q-avatar
        :class="[
          '',
          'm-t-10',
          // 设置整体头像的尺寸大小。 至于内部相关文字的尺寸大小(如果有的话  可通过设置text-[8rem] leading-[10rem]实现)
          'size-56',
          '',
        ]"
      >
        <img :src="logoUrl" />
      </q-avatar>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import logoUrl from 'assets/img/KeyTone.png?url';
</script>

<style lang="scss" scoped></style>
