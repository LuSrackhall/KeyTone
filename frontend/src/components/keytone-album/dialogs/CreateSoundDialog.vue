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
文件说明: dialogs/CreateSoundDialog.vue - 创建声音对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于从音频源文件创建新的声音定义。
功能包括：
1. 选择音频源文件
2. 设置声音名称（可选，默认使用源文件名+时间范围）
3. 裁剪时间范围（开始时间、结束时间）
4. 调整音量
5. 预览声音效果
6. 保存声音配置

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepDefineSounds.vue
        │         │
        │         └── dialogs/CreateSoundDialog.vue  <── 当前文件
        │
        └── (其他 Step 也可以调用此对话框)

【数据流】
  父组件状态                      本组件使用方式
  ─────────────────────────────────────────────────────
  ctx.createNewSound         ->
v-model 控制对话框显示 ctx.soundName -> 声音名称输入 ctx.sourceFileForSound -> 选中的源文件 ctx.soundFileList ->
可选择的源文件列表 ctx.soundStartTime -> 裁剪开始时间 ctx.soundEndTime -> 裁剪结束时间 ctx.soundVolume -> 音量调整
ctx.saveSoundConfig() -> 保存声音 ctx.previewSound() -> 预览声音 【关联文件】 - ../types.ts : 类型定义 -
../steps/StepDefineSounds.vue : 使用此对话框的 Step 组件
============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.createNewSound.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card>
      <!-- 对话框标题 -->
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ ctx.$t('KeyToneAlbum.defineSounds.createNewSound') }}
      </q-card-section>

      <!-- 声音名称输入 -->
      <q-card-section :class="['p-b-1']">
        <q-input
          outlined
          stack-label
          dense
          v-model="ctx.soundName.value"
          :label="ctx.$t('KeyToneAlbum.defineSounds.soundName')"
          :placeholder="defaultSoundName"
          :input-style="{ textOverflow: 'ellipsis' }"
          :input-class="'text-truncate'"
        >
          <template v-slot:append>
            <q-icon name="info" color="primary">
              <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                {{
                  ctx.$t('KeyToneAlbum.defineSounds.tooltip.soundName') +
                  ' : \n' +
                  (ctx.soundName.value === '' ? defaultSoundName : ctx.soundName.value)
                }}
              </q-tooltip>
            </q-icon>
          </template>
        </q-input>
      </q-card-section>

      <!-- 选择源文件 -->
      <q-card-section :class="['p-b-1']">
        <q-select
          outlined
          stack-label
          :virtual-scroll-slice-size="999999"
          v-model="ctx.sourceFileForSound.value"
          :options="ctx.soundFileList.value"
          :option-label="(item: any) => item.name + item.type"
          :label="ctx.$t('KeyToneAlbum.defineSounds.sourceFile')"
          dense
          popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
        />
      </q-card-section>

      <!-- 裁剪时间设置 -->
      <q-card-section :class="['p-b-1']">
        <div class="text-[13.5px] text-gray-600 p-b-2">
          {{ ctx.$t('KeyToneAlbum.defineSounds.cropSound') }}
          <q-icon name="info" color="primary">
            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
              {{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.soundDuration') }}
            </q-tooltip>
          </q-icon>
        </div>

        <!-- 波形裁剪（可视化选区） -->
        <WaveformTrimmer
          v-if="ctx.createNewSound.value"
          :sha256="ctx.sourceFileForSound.value.sha256"
          :file-type="ctx.sourceFileForSound.value.type"
          :volume="ctx.soundVolume.value"
          v-model:startMs="ctx.soundStartTime.value"
          v-model:endMs="ctx.soundEndTime.value"
        />

        <!--
          TIPS: 注意 number 类型使用时需要使用 v-model.number
          这样可以自动处理 01、00.55 这种输入，将其自动变更为 1、0.55
        -->
        <div class="flex flex-row">
          <q-input
            :class="['w-1/2 p-r-1']"
            outlined
            stack-label
            dense
            v-model.number="ctx.soundStartTime.value"
            :label="ctx.$t('KeyToneAlbum.defineSounds.startTime')"
            type="number"
            :error-message="ctx.$t('KeyToneAlbum.defineSounds.error.negativeTime')"
            :error="ctx.soundStartTime.value < 0"
          />
          <q-input
            :class="['w-1/2 p-l-1']"
            outlined
            stack-label
            dense
            v-model.number="ctx.soundEndTime.value"
            :label="ctx.$t('KeyToneAlbum.defineSounds.endTime')"
            type="number"
            :error-message="ctx.$t('KeyToneAlbum.defineSounds.error.negativeTime')"
            :error="ctx.soundEndTime.value < 0"
          />
        </div>
      </q-card-section>

      <!-- 音量设置 -->
      <q-card-section :class="['p-y-0']">
        <q-input
          outlined
          stack-label
          dense
          v-model.number="ctx.soundVolume.value"
          :label="ctx.$t('KeyToneAlbum.defineSounds.volume')"
          type="number"
          :step="0.1"
        >
          <template v-slot:append>
            <q-icon name="info" color="primary">
              <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                {{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.volume') }}
                {{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.volume_1') }}
              </q-tooltip>
            </q-icon>
          </template>
        </q-input>
      </q-card-section>

      <!-- 操作按钮 -->
      <q-card-actions align="right">
        <!-- 预览按钮 -->
        <q-btn
          class="mt-2"
          dense
          @click="handlePreview"
          :label="ctx.$t('KeyToneAlbum.defineSounds.previewSound')"
          color="secondary"
        >
          <q-tooltip
            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']"
            :delay="600"
          >
            {{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.previewSound') }}
          </q-tooltip>
        </q-btn>

        <!-- 确认添加按钮 -->
        <q-btn
          class="mt-2"
          @click="handleSave"
          :label="ctx.$t('KeyToneAlbum.defineSounds.confirmAdd')"
          color="primary"
        />

        <!-- 关闭按钮 -->
        <q-btn class="mt-2" flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * CreateSoundDialog.vue - 创建声音对话框
 *
 * 【组件职责】
 * 提供创建新声音的完整表单界面，包括：
 * - 声音名称设置（可选）
 * - 源文件选择
 * - 时间范围裁剪
 * - 音量调整
 * - 预览和保存功能
 *
 * 【重要说明】
 * 保存时需要手动选择字段传递给 saveSoundConfig，
 * 避免将 sourceFileForSound 中不需要的 name 字段传入。
 * 这是因为 JS/TS 中对象是引用传递，解构也会复制所有属性。
 */

import { inject, computed } from 'vue';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import WaveformTrimmer from '../components/WaveformTrimmer.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 计算属性
// ============================================================================

/**
 * 默认声音名称
 * 当用户未输入声音名称时，使用 "源文件名 - [开始时间 ~ 结束时间]" 格式
 */
const defaultSoundName = computed(() => {
  return `${ctx.sourceFileForSound.value.name}     - [${ctx.soundStartTime.value} ~ ${ctx.soundEndTime.value}]`;
});

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 预览声音
 * 使用当前配置播放声音预览
 */
function handlePreview() {
  ctx.previewSound({
    source_file_for_sound: ctx.sourceFileForSound.value,
    cut: {
      start_time: ctx.soundStartTime.value,
      end_time: ctx.soundEndTime.value,
      volume: ctx.soundVolume.value,
    },
  });
}

/**
 * 保存声音配置
 *
 * 【重要】手动选择字段传递
 * 通过手动选择字段，确保只传递所需的字段（sha256, name_id, type），
 * 而不会包含不需要的 name 字段。
 *
 * 保存成功后重置表单状态。
 */
function handleSave() {
  ctx.saveSoundConfig({
    // 手动选择字段：确保只传递所需的字段，不包含任何不需要的字段
    source_file_for_sound: {
      sha256: ctx.sourceFileForSound.value.sha256,
      name_id: ctx.sourceFileForSound.value.name_id,
      type: ctx.sourceFileForSound.value.type,
    },
    name: ctx.soundName.value,
    cut: {
      start_time: ctx.soundStartTime.value,
      end_time: ctx.soundEndTime.value,
      volume: ctx.soundVolume.value,
    },
    onSuccess: () => {
      // 重置表单状态
      ctx.soundName.value = '';
      ctx.sourceFileForSound.value = {
        sha256: '',
        name_id: '',
        name: '',
        type: '',
      };
      ctx.soundStartTime.value = 0;
      ctx.soundEndTime.value = 0;
      ctx.soundVolume.value = 0.0;
    },
  });
}
</script>

<style lang="scss" scoped>
/**
 * CreateSoundDialog 组件样式
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
