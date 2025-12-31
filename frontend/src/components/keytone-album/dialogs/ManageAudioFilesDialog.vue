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
文件说明: dialogs/ManageAudioFilesDialog.vue - 管理音频文件对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于管理（重命名/删除）已有的音频源文件。
功能包括：
1. 下拉选择已有的音频文件
2. 显示选中文件的详情卡片
3. 支持重命名文件
4. 支持删除文件

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepLoadAudioFiles.vue
        │         │
        │         └── dialogs/ManageAudioFilesDialog.vue  <── 当前文件
        │
        └── (其他 Step 也可以调用此对话框)

【数据流】
  父组件状态                  本组件使用方式
  ─────────────────────────────────────────────────
  ctx.editSoundFile      --> v-model 控制对话框显示/隐藏
  ctx.soundFileList      --> 下拉选项列表
  ctx.selectedSoundFile  --> 当前选中的文件
  ctx.$t()               --> i18n 翻译函数

【关联文件】
- ../types.ts                        : 类型定义
- ../steps/StepLoadAudioFiles.vue    : 使用此对话框的 Step 组件
- src/boot/query/keytonePkg-query.ts : SoundFileDelete API

【当前状态】
⚠️ 注意：本组件已创建但尚未集成到父组件中！

============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize }"
    v-model="ctx.editSoundFile.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="['p-x-3 w-[96%]']">
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ ctx.$t('KeyToneAlbum.loadAudioFile.manageExistingFiles') }}
      </q-card-section>

      <q-card-section>
        <q-select
          outlined
          :virtual-scroll-slice-size="999999"
          stack-label
          v-model="ctx.selectedSoundFile.value"
          :options="ctx.soundFileList.value"
          :option-label="(item: any) => item.name + item.type"
          :label="ctx.$t('KeyToneAlbum.loadAudioFile.selectFileToManage')"
          dense
          popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
        >
          <!-- 清除按钮 -->
          <template
            v-if="ctx.selectedSoundFile.value.sha256 !== '' && ctx.selectedSoundFile.value.name_id !== ''"
            v-slot:append
          >
            <q-icon name="cancel" @click.stop.prevent="clearSelection" class="cursor-pointer text-lg" />
          </template>
        </q-select>
      </q-card-section>

      <!-- 分割线 -->
      <q-separator v-if="ctx.selectedSoundFile.value.sha256 !== '' && ctx.selectedSoundFile.value.name_id !== ''" />

      <!-- 选中的音频源文件卡片 -->
      <q-card-section
        v-if="ctx.selectedSoundFile.value.sha256 !== '' && ctx.selectedSoundFile.value.name_id !== ''"
        :class="['flex flex-col m-t-3']"
      >
        <q-card :class="['flex flex-col']">
          <q-badge
            transparent
            color="orange"
            :label="ctx.selectedSoundFile.value.type"
            :class="['absolute overflow-visible right-0']"
          />
          <q-card-section
            v-if="ctx.selectedSoundFile.value.sha256 !== '' && ctx.selectedSoundFile.value.name_id !== ''"
            :class="['flex flex-col m-t-3']"
          >
            <!-- 重命名输入框 -->
            <q-input
              outlined
              stack-label
              dense
              :error-message="ctx.$t('KeyToneAlbum.notify.emptyFileName')"
              :error="
                ctx.selectedSoundFile.value.name === '' ||
                ctx.selectedSoundFile.value.name === undefined ||
                ctx.selectedSoundFile.value.name === null
              "
              v-model="ctx.selectedSoundFile.value.name"
              :label="ctx.$t('KeyToneAlbum.loadAudioFile.renameFile')"
            />

            <!-- 删除按钮 -->
            <q-btn
              :class="['w-20 self-center bg-pink-700 text-zinc-50']"
              dense
              no-caps
              :label="ctx.$t('KeyToneAlbum.delete')"
              icon="flight_takeoff"
              @click="handleDelete"
            />
          </q-card-section>
        </q-card>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { inject } from 'vue';
import { useQuasar } from 'quasar';
import { SoundFileDelete } from 'src/boot/query/keytonePkg-query';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';

const q = useQuasar();

// 注入父组件提供的上下文
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// 清除选择
function clearSelection() {
  ctx.selectedSoundFile.value = {
    sha256: '',
    name_id: '',
    name: '',
    type: '',
  };
}

// 删除文件
async function handleDelete() {
  const re = await SoundFileDelete(
    ctx.selectedSoundFile.value.sha256,
    ctx.selectedSoundFile.value.name_id,
    ctx.selectedSoundFile.value.type
  );

  if (re) {
    q.notify({
      type: 'positive',
      position: 'top',
      message: ctx.$t('KeyToneAlbum.notify.deleteSuccess'),
      timeout: 5,
    });
    // 清除前端结构体对象
    clearSelection();
  } else {
    q.notify({
      type: 'negative',
      position: 'top',
      message: ctx.$t('KeyToneAlbum.notify.deleteFailed'),
      timeout: 5,
    });
  }
}
</script>

<style lang="scss" scoped>
/**
 * ManageAudioFilesDialog 组件样式
 */

// 按钮样式 - 统一按钮外观
.q-btn {
  @apply text-xs;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}

// 选择器样式 - 处理溢出
:deep(.q-field__native) {
  @apply max-w-full overflow-auto whitespace-nowrap;
  @apply h-5.8 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50 [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400;
}

// 输入框标签样式
:deep(.q-field__label) {
  @apply overflow-visible -ml-1.5 text-[0.8rem];
}

// 输入框 placeholder 高度修复
:deep(.q-placeholder) {
  @apply h-auto;
}
</style>
