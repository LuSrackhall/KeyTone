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
文件说明: dialogs/CreateKeySoundDialog.vue - 创建按键音对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于创建新的按键音定义。
功能包括：
1. 设置按键音名称
2. 配置按下(down)时的声音（模式 + 声音列表）
3. 配置抬起(up)时的声音（模式 + 声音列表）
4. 支持选择音频文件、声音、按键音三种类型
5. 支持单次、随机、循环三种播放模式

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepCraftKeySounds.vue
        │         │
        │         └── dialogs/CreateKeySoundDialog.vue  <── 当前文件

【数据流】
  父组件状态                          本组件使用方式
  ───────────────────────────────────────────────────────────
  ctx.createNewKeySound          ->
v-model 控制对话框显示 ctx.keySoundName -> 按键音名称 ctx.configureDownSound -> 控制按下声音配置子对话框
ctx.configureUpSound -> 控制抬起声音配置子对话框 ctx.selectedSoundsForDown -> 按下时选中的声音列表
ctx.selectedSoundsForUp -> 抬起时选中的声音列表 ctx.playModeForDown/Up -> 播放模式 ctx.downSoundList/upSoundList ->
可选择的声音列表 ctx.saveKeySoundConfig() -> 保存按键音
============================================================================
-->

<template>
  <!-- 主对话框：创建新按键音 -->
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.createNewKeySound.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="['min-w-[90%]', { 'mr-0': isMac }]">
      <!-- 对话框标题 -->
      <q-card-section class="row items-center q-pb-none text-h6">
        {{ ctx.$t('KeyToneAlbum.craftKeySounds.newKeySound') }}
      </q-card-section>

      <!-- 按键音名称输入 -->
      <q-card-section :class="['p-b-1']">
        <q-input
          outlined
          stack-label
          dense
          v-model="ctx.keySoundName.value"
          :label="ctx.$t('KeyToneAlbum.craftKeySounds.keySoundName')"
          :placeholder="ctx.$t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder')"
        />

        <!-- 配置按钮组 -->
        <div class="flex flex-col mt-3">
          <!-- 配置按下声音按钮 -->
          <q-btn
            :class="['bg-zinc-300 my-7 w-88% self-center']"
            :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureDownSound')"
            @click="ctx.configureDownSound.value = true"
          />

          <!-- 配置按下声音子对话框 -->
          <q-dialog
            :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
            v-model="ctx.configureDownSound.value"
            backdrop-filter="invert(70%)"
            @mouseup="ctx.preventDefaultMouseWhenRecording"
          >
            <q-card :class="['min-w-[80%]', { 'mr-0': isMac }]">
              <q-card-section class="row items-center q-pb-none text-h6">
                {{ ctx.$t('KeyToneAlbum.craftKeySounds.configureDownSound') }}
              </q-card-section>
              <q-card-section>
                <!-- 使用选择框选择模式 -->
                <q-select
                  outlined
                  popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                  :virtual-scroll-slice-size="999999"
                  stack-label
                  v-model="ctx.playModeForDown.value"
                  :options="ctx.playModeOptions"
                  :option-label="(item: any) => ctx.$t(ctx.playModeLabels.get(item) || '')"
                  :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                  dense
                />
              </q-card-section>
              <q-card-section>
                <!-- 选择声音的选项，支持多选 -->
                <q-select
                  outlined
                  stack-label
                  :virtual-scroll-slice-size="999999"
                  v-model="ctx.selectedSoundsForDown.value"
                  popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                  :options="ctx.downSoundList.value"
                  :option-label="ctx.album_options_select_label"
                  :option-value="getOptionValue"
                  :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                  multiple
                  use-chips
                  :class="['zl-ll']"
                  dense
                  :max-values="ctx.maxSelectionForDown.value"
                  counter
                  :error-message="ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                  :error="ctx.playModeForDown.value === 'single' ? ctx.selectedSoundsForDown.value.length > 1 : false"
                  ref="downSoundSelectDom"
                  @update:model-value="downSoundSelectDom?.hidePopup()"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section>
                        <q-item-label>{{ ctx.album_options_select_label(scope.opt) }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <DependencyWarning
                          v-if="scope.opt.type === 'sounds'"
                          :issues="ctx.dependencyIssues.value"
                          item-type="sounds"
                          :item-id="scope.opt.value.soundKey"
                          :show-details="false"
                        />
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>
                <div class="h-10">
                  <q-option-group
                    dense
                    v-model="ctx.downTypeGroup.value"
                    :options="ctx.options"
                    type="checkbox"
                    class="absolute left-8"
                  >
                    <template #label-0="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                    <template v-slot:label-1="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                    <template v-slot:label-2="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                            <div>⬇</div>
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                            <div>⬇</div>
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                  </q-option-group>
                </div>
              </q-card-section>
              <q-card-actions align="right">
                <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
              </q-card-actions>
            </q-card>
          </q-dialog>

          <!-- 配置抬起声音按钮 -->
          <q-btn
            :class="['bg-zinc-300 m-b-7 w-88% self-center']"
            :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
            @click="ctx.configureUpSound.value = true"
          />

          <!-- 配置抬起声音子对话框 -->
          <q-dialog
            :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
            v-model="ctx.configureUpSound.value"
            backdrop-filter="invert(70%)"
            @mouseup="ctx.preventDefaultMouseWhenRecording"
          >
            <q-card :class="['min-w-[80%]', { 'mr-0': isMac }]">
              <q-card-section class="row items-center q-pb-none text-h6">
                {{ ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound') }}
              </q-card-section>
              <q-card-section>
                <!-- 使用选择框选择模式 -->
                <q-select
                  outlined
                  popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                  :virtual-scroll-slice-size="999999"
                  stack-label
                  v-model="ctx.playModeForUp.value"
                  :options="ctx.playModeOptions"
                  :option-label="(item: any) => ctx.$t(ctx.playModeLabels.get(item) || '')"
                  :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                  dense
                />
              </q-card-section>
              <q-card-section>
                <!-- 选择声音的选项，支持多选 -->
                <q-select
                  outlined
                  stack-label
                  :virtual-scroll-slice-size="999999"
                  v-model="ctx.selectedSoundsForUp.value"
                  popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                  :options="ctx.upSoundList.value"
                  :option-label="ctx.album_options_select_label"
                  :option-value="getOptionValue"
                  :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                  multiple
                  use-chips
                  :class="['zl-ll']"
                  dense
                  :max-values="ctx.maxSelectionForUp.value"
                  counter
                  :error-message="ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                  :error="ctx.playModeForUp.value === 'single' ? ctx.selectedSoundsForUp.value.length > 1 : false"
                  ref="upSoundSelectDom"
                  @update:model-value="upSoundSelectDom?.hidePopup()"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section>
                        <q-item-label>{{ ctx.album_options_select_label(scope.opt) }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <DependencyWarning
                          v-if="scope.opt.type === 'sounds'"
                          :issues="ctx.dependencyIssues.value"
                          item-type="sounds"
                          :item-id="scope.opt.value.soundKey"
                          :show-details="false"
                        />
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>
                <div class="h-10">
                  <q-option-group
                    dense
                    v-model="ctx.upTypeGroup.value"
                    :options="ctx.options"
                    type="checkbox"
                    class="absolute left-8"
                  >
                    <template #label-0="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                    <template v-slot:label-1="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                    <template v-slot:label-2="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                            <div>⬇</div>
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                            <div>⬇</div>
                            <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                  </q-option-group>
                </div>
              </q-card-section>
              <q-card-actions align="right">
                <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
              </q-card-actions>
            </q-card>
          </q-dialog>
        </div>
      </q-card-section>

      <!-- 操作按钮 -->
      <q-card-actions align="right">
        <q-btn color="primary" :label="ctx.$t('KeyToneAlbum.craftKeySounds.confirmAdd')" @click="handleSave" />
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * CreateKeySoundDialog.vue - 创建按键音对话框
 *
 * 【组件职责】
 * 提供创建新按键音的完整界面，包括：
 * - 按键音名称设置
 * - 按下声音配置（内嵌子对话框）
 * - 抬起声音配置（内嵌子对话框）
 * - 保存功能
 */

import { inject, ref, computed } from 'vue';
import { Platform } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

const isMac = computed(() => Platform.is.mac === true);

// ============================================================================
// DOM 引用
// ============================================================================
const downSoundSelectDom = ref<any>(null);
const upSoundSelectDom = ref<any>(null);

// ============================================================================
// 工具函数
// ============================================================================

/**
 * 获取选项的唯一值
 * 用于 q-select 组件的 option-value 属性
 */
function getOptionValue(item: any) {
  if (item.type === 'audio_files') {
    return item.value.sha256 + item.value.name_id;
  }
  if (item.type === 'sounds') {
    return item.value.soundKey;
  }
  if (item.type === 'key_sounds') {
    return item.value.keySoundKey;
  }
}

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 保存按键音配置
 */
function handleSave() {
  ctx.saveKeySoundConfig(
    {
      key: '',
      name: ctx.keySoundName.value,
      down: {
        mode: ctx.playModeForDown.value,
        value: ctx.selectedSoundsForDown.value,
      },
      up: {
        mode: ctx.playModeForUp.value,
        value: ctx.selectedSoundsForUp.value,
      },
    },
    () => {
      // 关闭对话框
      ctx.createNewKeySound.value = false;

      // 重置表单变量
      ctx.keySoundName.value = ctx.$t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder');
      ctx.selectedSoundsForDown.value = [];
      ctx.playModeForDown.value = 'random';
      ctx.selectedSoundsForUp.value = [];
      ctx.playModeForUp.value = 'random';
    }
  );
}
</script>

<style lang="scss" scoped>
/**
 * CreateKeySoundDialog 组件样式
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

// 按键音选择器专用样式 - 用于多选芯片选择框
.zl-ll {
  :deep(.q-field__native) {
    @apply h-auto;
  }
  :deep(.q-field__messages) {
    @apply text-nowrap;
  }
}

// 椭圆省略样式 - 溢出处理
:deep(.ellipsis) {
  @apply max-w-full overflow-auto whitespace-nowrap text-clip;
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}
</style>
