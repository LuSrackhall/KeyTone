<template>
  <!-- 导入一个已经制作好的音频包, 比如其它用户制作了一个音频包分享给你使用。
       * 为防止用户随便的将一个音频文件作为音频包来导入引起的错误, 这里需要加校验, 以确保导入的音频包是可以正常工作的。
       * 请放心, 使用keytone软件制作的音频包, 是一定能够成功导入的。
  -->
  <!-- <div>导入音频包</div> -->

  <!-- 编辑一个已经制作好的音频包, 比如其它用户制作了一个音频包分享给你使用, 但有些地方你像要进一步的进行定制。
       * 用户先选择好要编辑的音频包, 然后点击编辑进入编辑页面, 最后进行编辑(编辑时所作出的更改是实时保存的, 无主动保存按钮)。
  -->
  <!-- <div>编辑音频包</div> -->

  <!-- 制作一个新的音频包
       * 做出最符合自身需求的音频包后, 您还可以将这个音频包进行导出分享, 以供其它用户使用你制作的音频包。
       * 用户先选择好要编辑的音频包, 然后点击编辑进入编辑页面, 最后进行编辑(编辑时所作出的更改是实时保存的, 无主动保存按钮)。
  -->
  <!-- <div>制作音频包</div> -->

  <!-- 导出已有KeyTone音频包(仅由KeyTone制作或可用KeyTone进行编辑的音频包支持导出功能) -->
  <!-- 不将其作为单独的按钮项, 而是集成到制作音频包和编辑音频包的页面中去 -->
  <!-- <div>导出音频包</div> -->

  <q-page>
    <q-scroll-area class="w-[379px] h-[458.5px]">
      <div
        :class="[
          '',
          'w-full flex flex-col gap-10 items-center  mt-25 ',
          // 使页面文本无法选择, 免得影响界面体验。
          'select-none',
        ]"
      >
        <q-btn no-caps :label="$t('KeyTonePackage.new.new')" icon="flight_takeoff" />
        <q-btn no-caps :label="$t('KeyTonePackage.edit.edit')" icon="flight_takeoff" />
        <q-btn no-caps :label="$t('KeyTonePackage.import.import')" icon="flight_takeoff" />
      </div>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped></style>
