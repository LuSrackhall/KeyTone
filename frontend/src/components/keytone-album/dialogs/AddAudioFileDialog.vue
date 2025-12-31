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
文件说明: dialogs/AddAudioFileDialog.vue - 添加音频文件对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于添加新的音频源文件到键音专辑。
功能包括：
1. 提供文件选择器，支持拖拽上传
2. 支持多文件同时上传
3. 限制文件类型为 .wav, .mp3, .ogg
4. 上传成功/失败时显示通知

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepLoadAudioFiles.vue
        │         │
        │         └── dialogs/AddAudioFileDialog.vue  <── 当前文件
        │
        └── (其他 Step 也可以调用此对话框)

【为什么对话框要独立成组件】
1. 可复用性：同一个对话框可能被多个 Step 调用
2. 代码清晰：将对话框 UI 与 Step 逻辑分离
3. 便于维护：修改对话框不影响 Step 组件

【数据流】
  父组件状态                  本组件使用方式
  ─────────────────────────────────────────────────
  ctx.addNewSoundFile    --> v-model 控制对话框显示/隐藏
  ctx.files              --> v-model 绑定文件列表
  ctx.$t()               --> i18n 翻译函数
  ctx.i18n_fontSize      --> 字体大小样式变量

【关联文件】
- ../types.ts                        : 类型定义
- ../steps/StepLoadAudioFiles.vue    : 使用此对话框的 Step 组件
- src/boot/query/keytonePkg-query.ts : SendFileToServer API

【当前状态】
⚠️ 注意：本组件已创建但尚未集成到父组件中！

============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize }"
    v-model="ctx.addNewSoundFile.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card>
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ ctx.$t('KeyToneAlbum.loadAudioFile.addNewFile_1') }}
      </q-card-section>

      <q-card-section>
        <div class="text-gray-600 text-xs">{{ ctx.$t('KeyToneAlbum.loadAudioFile.dragAndDrop') }}</div>
        <q-file
          :class="['w-56', 'zl-ll']"
          dense
          v-model="ctx.files.value"
          :label="ctx.$t('KeyToneAlbum.loadAudioFile.audioFile')"
          outlined
          use-chips
          multiple
          append
          accept=".wav,.mp3,.ogg"
          excludeAcceptAllOption
          style="max-width: 300px"
          :hint="ctx.$t('KeyToneAlbum.loadAudioFile.supportedFormats')"
        />
      </q-card-section>

      <q-card-section>
        <div>{{ ctx.$t('KeyToneAlbum.loadAudioFile.addAsNeeded') }}</div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          flat
          @click="handleConfirmAdd"
          color="primary"
          :label="ctx.$t('KeyToneAlbum.loadAudioFile.confirmAdd')"
        />
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { inject, nextTick } from 'vue';
import { useQuasar } from 'quasar';
import { SendFileToServer } from 'src/boot/query/keytonePkg-query';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';

const q = useQuasar();

// 注入父组件提供的上下文
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// 确认添加文件
async function handleConfirmAdd() {
  if (!ctx.files.value || ctx.files.value.length === 0) {
    console.warn('No files selected for upload');
    return;
  }

  // 使用 slice() 创建浅拷贝，避免遍历过程中修改原数组导致问题
  for (const file of ctx.files.value.slice()) {
    try {
      const re = await SendFileToServer(file);
      if (re === true) {
        console.info(`File ${file.name} uploaded successfully`);
        // 上传成功后从列表中移除
        const index = ctx.files.value.indexOf(file);
        if (index > -1) {
          ctx.files.value.splice(index, 1);
        }
      } else {
        console.error(`File ${file.name} uploading error`);
        q.notify({
          type: 'negative',
          position: 'top',
          message: `${ctx.$t('KeyToneAlbum.notify.addFailed')} '${file.name}'`,
          timeout: 5,
        });
        return;
      }
    } catch (error) {
      console.error(`Error uploading file ${file.name}:`, error);
      q.notify({
        type: 'negative',
        position: 'top',
        message: `${ctx.$t('KeyToneAlbum.notify.addFailed')} '${file.name}'`,
        timeout: 5,
      });
      return;
    }
  }

  nextTick(() => {
    q.notify({
      type: 'positive',
      position: 'top',
      message: ctx.$t('KeyToneAlbum.notify.addSuccess'),
      timeout: 5,
    });
  });
}
</script>

<style lang="scss" scoped>
/**
 * AddAudioFileDialog 组件样式
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
</style>
