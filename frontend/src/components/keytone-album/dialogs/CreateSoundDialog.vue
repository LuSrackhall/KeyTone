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
    class="create-sound-dialog-单独影响global"
    @before-hide="onDialogBeforeHide"
  >
    <!--
      重要：KeyTone 窗口有固定宽度（约 379~389px）。
      - 对话框不能超出窗口宽度，但需要尽可能利用可视空间。
      - 这里直接使用视口宽度减去极小边距（8px），避免“看起来还能更宽”但被人为留白。
    -->
    <q-card
      style="
        /*
          目标：
          - 在固定窗口宽度下，尽可能放大对话框宽度（不超出窗口）。
          - 为波形外侧刻度留空间，但不压缩波形本体宽度。
        */
        width: calc(100vw - 8px);
        max-width: calc(100vw - 8px);
      "
      :class="['p-l-2 p-r-5',
        { 'mr-0': isMac } // Mac 平台下, 右侧不留额外空隙, 因为阴影用的是原生的
      ]"
    >
      <!-- 对话框标题（sticky 置顶） -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
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
          ref="waveformRef"
          v-show="ctx.createNewSound.value"
          :sha256="ctx.sourceFileForSound.value.sha256"
          :file-type="ctx.sourceFileForSound.value.type"
          v-model:volume="ctx.soundVolume.value"
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
          v-model.number="soundVolumeDb"
          :label="ctx.$t('KeyToneAlbum.defineSounds.volumeDb')"
          type="number"
          :step="0.1"
          :hint="isSoundVolumeDbOutOfRange ? ctx.$t('KeyToneAlbum.defineSounds.volumeOutOfRange') : ''"
          :hide-bottom-space="!isSoundVolumeDbOutOfRange"
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

      <!-- 操作按钮（sticky 置底） -->
      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm', { 'mr-0': isMac }]">
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

import { inject, computed, ref, watch } from 'vue';
import { Platform } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import WaveformTrimmer from '../components/WaveformTrimmer.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

const waveformRef = ref<InstanceType<typeof WaveformTrimmer> | null>(null);

/**
 * 监听对话框关闭（before-hide）：
 * - 在动画开始前立即停止音频播放，避免“关闭后仍有声音”的不符合直觉现象。
 * - 使用 v-show 保持组件存在，避免销毁重计算导致关闭动画丢失。
 */
function onDialogBeforeHide() {
  waveformRef.value?.stopPlayback?.();
}

// 兜底：当 dialog 通过任意方式关闭（包括外部强制设置 v-model=false）时，立即停止播放
watch(
  () => ctx.createNewSound.value,
  (val) => {
    if (!val) waveformRef.value?.stopPlayback?.();
  }
);

// 使用 Quasar 提供的前端平台检测，仅依赖前端环境
const isMac = computed(() => Platform.is.mac === true);

// ============================================================================
// dB <-> cut.volume 换算（Base=1.6）
// - SDK：gain = 1.6 ^ volume
// - dB = 20 * log10(gain) = 20 * volume * log10(1.6)
// - UI 以 dB 显示，内部仍使用 cut.volume
// ==========================================================================
const dbPerVolume = 20 * Math.log10(1.6);
const volumeToDb = (volume: number) => volume * dbPerVolume;
const dbToVolume = (db: number) => db / dbPerVolume;

const soundVolumeDb = computed({
  get: () => Number(volumeToDb(ctx.soundVolume.value || 0).toFixed(1)),
  set: (db: number) => {
    ctx.soundVolume.value = dbToVolume(Number(db));
  },
});

const isSoundVolumeDbOutOfRange = computed(() => Math.abs(soundVolumeDb.value) > 18);

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
  // 恢复默认大小：避免“按钮被放大”的非预期视觉变化。
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

// TIPS: 对话框实际宽度调整, 只能通过覆盖全局样式实现(因为 q-dialog 实际是基于当前组件外部的全局组件实现的)
// :global(.q-dialog__inner--minimized){ // TIPS: 我们可以通过添加类名的方式, 只修改特定对话框的样式, 具体见下方的操作。
//   @apply p-4;
// }

// TIPS: 虽然对于全局样式的覆盖, 只能通过 :global 实现, 想要进修改单个组件的样式(不影响其他用到此组件的业务)
//       > 可以在 :global 内部继续使用组件作用域的类名选择器继承的方式, 以避免影响其他组件的同名类选择器
:global(.create-sound-dialog-单独影响global .q-dialog__inner--minimized) {
  @apply p-x-2;
}
</style>
