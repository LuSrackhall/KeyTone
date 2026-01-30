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
文件说明: dialogs/EditSoundDialog.vue - 编辑声音对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于编辑已有的声音定义。
功能包括：
1. 从声音列表中选择要编辑的声音
2. 修改声音名称
3. 更换源文件
4. 调整裁剪时间范围
5. 调整音量
6. 预览声音效果
7. 保存或删除声音

【特殊功能】
- 在选项列表中显示依赖警告（DependencyWarning）
- 选中声音后展示详细编辑卡片
- 支持实时预览修改效果

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepDefineSounds.vue
        │         │
        │         └── dialogs/EditSoundDialog.vue  <── 当前文件
        │
        └── (其他 Step 也可以调用此对话框)

【数据流】
  父组件状态                      本组件使用方式
  ─────────────────────────────────────────────────────
  ctx.showEditSoundDialog    ->
v-model 控制对话框显示 ctx.soundList -> 可选择的声音列表 ctx.selectedSound -> 当前选中的声音 ctx.soundFileList ->
源文件列表（用于更换源文件） ctx.dependencyIssues -> 依赖问题列表 ctx.saveSoundConfig() -> 保存声音 ctx.deleteSound()
-> 删除声音 ctx.previewSound() -> 预览声音 【关联文件】 - ../types.ts : 类型定义 - ../steps/StepDefineSounds.vue :
使用此对话框的 Step 组件 - ../../DependencyWarning.vue : 依赖警告组件
============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.showEditSoundDialog.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card class="min-w-[106%]">
      <!-- 对话框标题（sticky 置顶） -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ ctx.$t('KeyToneAlbum.defineSounds.editExistingSound') }}
      </q-card-section>

      <!-- 声音选择下拉框 -->
      <q-card-section>
        <q-select
          outlined
          stack-label
          :virtual-scroll-slice-size="999999"
          clearable
          v-model="ctx.selectedSound.value"
          popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
          :options="ctx.soundList.value"
          :option-label="getSoundLabel"
          :option-value="(item: any) => item.soundKey"
          :label="ctx.$t('KeyToneAlbum.defineSounds.selectSoundToManage')"
          dense
        >
          <!-- 自定义选项渲染，包含依赖警告 -->
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ getSoundLabel(scope.opt) }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <DependencyWarning
                  :issues="ctx.dependencyIssues.value"
                  item-type="sounds"
                  :item-id="scope.opt.soundKey"
                  :show-details="false"
                />
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </q-card-section>

      <!-- 选中声音的编辑卡片 -->
      <q-card-section
        :class="['flex flex-col m-t-3']"
        v-if="ctx.selectedSound.value?.soundKey !== '' && ctx.selectedSound.value !== undefined"
      >
        <q-card :class="['flex flex-col pb-3 w-[100%]']" v-if="ctx.selectedSound.value">
          <!-- 声音名称输入 -->
          <q-card-section :class="['p-b-1 mt-3']">
            <q-input
              outlined
              stack-label
              dense
              v-model="ctx.selectedSound.value.soundValue.name"
              :label="ctx.$t('KeyToneAlbum.defineSounds.soundName')"
              :placeholder="selectedSoundPlaceholder"
              :input-style="{ textOverflow: 'ellipsis' }"
              :input-class="'text-truncate'"
            >
              <template v-slot:append>
                <q-icon name="info" color="primary">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                    {{
                      ctx.$t('KeyToneAlbum.defineSounds.tooltip.soundName') +
                      ' : \n' +
                      (ctx.selectedSound.value.soundValue.name === ''
                        ? selectedSoundPlaceholder
                        : ctx.selectedSound.value.soundValue.name)
                    }}
                  </q-tooltip>
                </q-icon>
              </template>
            </q-input>
          </q-card-section>

          <!-- 源文件选择 -->
          <q-card-section :class="['p-b-1 w-68']">
            <q-select
              outlined
              stack-label
              :virtual-scroll-slice-size="999999"
              popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
              v-model="ctx.selectedSound.value.soundValue.source_file_for_sound"
              :options="ctx.soundFileList.value"
              :option-label="getSourceFileLabel"
              :option-value="(item: any) => item.sha256 + item.name_id"
              :label="ctx.$t('KeyToneAlbum.defineSounds.sourceFile')"
              dense
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
              v-if="ctx.showEditSoundDialog.value"
              :sha256="ctx.selectedSound.value.soundValue.source_file_for_sound.sha256"
              :file-type="ctx.selectedSound.value.soundValue.source_file_for_sound.type"
                v-model:volume="ctx.selectedSound.value.soundValue.cut.volume"
              v-model:startMs="ctx.selectedSound.value.soundValue.cut.start_time"
              v-model:endMs="ctx.selectedSound.value.soundValue.cut.end_time"
            />

            <!--
              TIPS: 注意 number 类型使用时需要使用 v-model.number
              这样可以自动处理 01、00.55 这种输入
            -->
            <div class="flex flex-row">
              <q-input
                :class="['w-1/2 p-r-1']"
                outlined
                stack-label
                dense
                v-model.number="ctx.selectedSound.value.soundValue.cut.start_time"
                :label="ctx.$t('KeyToneAlbum.defineSounds.startTime')"
                type="number"
                :error-message="ctx.$t('KeyToneAlbum.defineSounds.error.negativeTime')"
                :error="ctx.selectedSound.value.soundValue.cut.start_time < 0"
              />
              <q-input
                :class="['w-1/2 p-l-1']"
                outlined
                stack-label
                dense
                v-model.number="ctx.selectedSound.value.soundValue.cut.end_time"
                :label="ctx.$t('KeyToneAlbum.defineSounds.endTime')"
                type="number"
                :error-message="ctx.$t('KeyToneAlbum.defineSounds.error.negativeTime')"
                :error="ctx.selectedSound.value.soundValue.cut.end_time < 0"
              />
            </div>
          </q-card-section>

          <!-- 音量设置 -->
          <q-card-section :class="['p-y-0']">
            <q-input
              outlined
              stack-label
              dense
              v-model.number="ctx.selectedSound.value.soundValue.cut.volume"
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

          <!-- 操作按钮组 -->
          <q-card-section :class="['flex justify-center gap-4']">
            <!-- 预览按钮 -->
            <q-btn
              class="pr-2.3"
              dense
              color="secondary"
              icon="play_arrow"
              :label="ctx.$t('KeyToneAlbum.defineSounds.previewSound')"
              @click="handlePreview"
            >
              <q-tooltip
                :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']"
                :delay="600"
              >
                {{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.previewSound') }}
              </q-tooltip>
            </q-btn>

            <!-- 保存按钮 -->
            <q-btn
              class="pr-2.3"
              dense
              color="primary"
              icon="save"
              :label="ctx.$t('KeyToneAlbum.confirmEdit')"
              @click="handleSave"
            />

            <!-- 删除按钮 -->
            <q-btn
              class="pr-2.3"
              dense
              color="negative"
              icon="delete"
              :label="ctx.$t('KeyToneAlbum.delete')"
              @click="handleDelete"
            />
          </q-card-section>
        </q-card>
      </q-card-section>

      <!-- 关闭按钮（sticky 置底） -->
      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * EditSoundDialog.vue - 编辑声音对话框
 *
 * 【组件职责】
 * 提供编辑已有声音的完整界面，包括：
 * - 从列表选择声音
 * - 修改声音属性
 * - 预览、保存、删除功能
 *
 * 【依赖警告】
 * 在选项列表中显示 DependencyWarning 组件，
 * 提示用户该声音是否被其他按键音引用。
 */

import { inject, computed } from 'vue';
import { useQuasar } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';
import WaveformTrimmer from '../components/WaveformTrimmer.vue';

const q = useQuasar();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 辅助函数
// ============================================================================

/**
 * 获取声音的显示标签
 * 如果声音有自定义名称则显示名称，否则显示 "源文件名 - [开始时间 ~ 结束时间]"
 */
function getSoundLabel(item: any): string {
  if (item.soundValue.name !== '' && item.soundValue.name !== undefined) {
    return item.soundValue.name;
  } else {
    const soundFile = ctx.soundFileList.value.find(
      (sf) =>
        sf.sha256 === item.soundValue.source_file_for_sound.sha256 &&
        sf.name_id === item.soundValue.source_file_for_sound.name_id
    );
    return `${soundFile?.name}     - [${item.soundValue.cut.start_time} ~ ${item.soundValue.cut.end_time}]`;
  }
}

/**
 * 获取源文件的显示标签
 * 通过 find 查找确保能享受 name 变化时的实时更新
 */
function getSourceFileLabel(item: any): string {
  const soundFile = ctx.soundFileList.value.find((sf) => sf.sha256 === item.sha256 && sf.name_id === item.name_id);
  return soundFile ? soundFile.name + soundFile.type : '';
}

// ============================================================================
// 计算属性
// ============================================================================

/**
 * 选中声音的默认名称占位符
 */
const selectedSoundPlaceholder = computed(() => {
  if (!ctx.selectedSound.value) return '';

  const soundFile = ctx.soundFileList.value.find(
    (sf) =>
      sf.sha256 === ctx.selectedSound.value?.soundValue.source_file_for_sound.sha256 &&
      sf.name_id === ctx.selectedSound.value?.soundValue.source_file_for_sound.name_id
  );

  return `${soundFile?.name}     - [${ctx.selectedSound.value.soundValue.cut.start_time} ~ ${ctx.selectedSound.value.soundValue.cut.end_time}]`;
});

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 预览声音
 */
function handlePreview() {
  if (!ctx.selectedSound.value) return;

  ctx.previewSound({
    source_file_for_sound: ctx.selectedSound.value.soundValue.source_file_for_sound,
    cut: ctx.selectedSound.value.soundValue.cut,
  });
}

/**
 * 保存声音修改
 * 手动选择字段传递，避免传递不需要的属性
 */
function handleSave() {
  if (!ctx.selectedSound.value) return;

  ctx.saveSoundConfig({
    soundKey: ctx.selectedSound.value.soundKey,
    source_file_for_sound: {
      sha256: ctx.selectedSound.value.soundValue.source_file_for_sound.sha256,
      name_id: ctx.selectedSound.value.soundValue.source_file_for_sound.name_id,
      type: ctx.selectedSound.value.soundValue.source_file_for_sound.type,
    },
    name: ctx.selectedSound.value.soundValue.name,
    cut: ctx.selectedSound.value.soundValue.cut,
    onSuccess: () => {
      // 保存成功后的回调（通知已在 saveSoundConfig 中处理）
    },
  });
}

/**
 * 删除声音
 * 删除成功后清空选中状态并显示通知
 */
function handleDelete() {
  if (!ctx.selectedSound.value) return;

  ctx.deleteSound({
    soundKey: ctx.selectedSound.value.soundKey,
    onSuccess: () => {
      ctx.selectedSound.value = undefined;
      q.notify({
        type: 'positive',
        position: 'top',
        message: ctx.$t('KeyToneAlbum.notify.deleteSuccess'),
        timeout: 5,
      });
    },
  });
}
</script>

<style lang="scss" scoped>
/**
 * EditSoundDialog 组件样式
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

// 选项列表项样式 - 溢出处理
:deep(.q-item__section) {
  @apply max-w-full overflow-auto whitespace-nowrap;
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

// 椭圆省略样式 - 溢出处理
:deep(.ellipsis) {
  @apply max-w-full overflow-auto whitespace-nowrap text-clip;
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}
</style>
