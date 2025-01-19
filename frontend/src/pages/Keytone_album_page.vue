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
        <q-btn
          no-caps
          :label="$t('KeyToneAlbum.new.new')"
          icon="flight_takeoff"
          @click="() => router.push('/keytone_album-create')"
        />
        <q-btn
          no-caps
          :label="$t('KeyToneAlbum.edit.edit')"
          icon="flight_takeoff"
          @click="
            () => {
              q.notify({
                type: 'info',
                position: 'top',
                message: '键音包编辑功能正在积极开发中...',
                timeout: 2000,
              });
              // q.notify({
              //   type: 'info',
              //   position: 'top',
              //   message: '赞助作者或是前往steam购买本软件来支持该项目, 以进一步增加作者的开发动力。',
              //   timeout: 2000,
              // });
            }
          "
        />
        <q-btn
          no-caps
          :label="$t('KeyToneAlbum.import.import')"
          icon="flight_takeoff"
          @click="
            () => {
              q.notify({
                type: 'info',
                position: 'top',
                message: '键音包导入功能正在积极开发中...',
                timeout: 2000,
              });
              // q.notify({
              //   type: 'info',
              //   position: 'top',
              //   message: '赞助作者或是前往steam购买本软件来支持该项目, 以进一步增加作者的开发动力。',
              //   timeout: 2000,
              // });
            }
          "
        />
      </div>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';

const router = useRouter();
const q = useQuasar();
</script>

<style lang="scss" scoped></style>
