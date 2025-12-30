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

【复杂度说明】
本对话框包含两个嵌套的子对话框（配置按下声音、配置抬起声音），
结构较为复杂。在未来的迭代中，可以考虑将子对话框进一步拆分。

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepCraftKeySounds.vue
        │         │
        │         └── dialogs/CreateKeySoundDialog.vue  <── 当前文件
        │                 │
        │                 ├── (内嵌) ConfigureDownSoundDialog
        │                 └── (内嵌) ConfigureUpSoundDialog

【数据流】
  父组件状态                          本组件使用方式
  ───────────────────────────────────────────────────────────
  ctx.createNewKeySound          -->
v-model 控制对话框显示 ctx.keySoundName --> 按键音名称 ctx.configureDownSound --> 控制按下声音配置子对话框
ctx.configureUpSound --> 控制抬起声音配置子对话框 ctx.selectedSoundsForDown --> 按下时选中的声音列表
ctx.selectedSoundsForUp --> 抬起时选中的声音列表 ctx.playModeForDown/Up --> 播放模式 ctx.downSoundList/upSoundList -->
可选择的声音列表 ctx.saveKeySoundConfig() --> 保存按键音 【关联文件】 - ../types.ts : 类型定义 -
../steps/StepCraftKeySounds.vue : 使用此对话框的 Step 组件 - ../../DependencyWarning.vue : 依赖警告组件
============================================================================ -->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.createNewKeySound.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="['min-w-[90%]']">
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
          <ConfigureDownSoundSubDialog />

          <!-- 配置抬起声音按钮 -->
          <q-btn
            :class="['bg-zinc-300 m-b-7 w-88% self-center']"
            :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
            @click="ctx.configureUpSound.value = true"
          />

          <!-- 配置抬起声音子对话框 -->
          <ConfigureUpSoundSubDialog />
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
 * - 按下声音配置
 * - 抬起声音配置
 * - 保存功能
 *
 * 【子对话框】
 * 配置按下/抬起声音的子对话框内嵌在此组件中，
 * 未来可以考虑进一步拆分为独立组件。
 */

import { inject, h, defineComponent } from 'vue';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 配置按下声音子对话框（内嵌组件）
// ============================================================================
const ConfigureDownSoundSubDialog = defineComponent({
  name: 'ConfigureDownSoundSubDialog',
  setup() {
    return () =>
      h(
        'q-dialog',
        {
          style: { '--i18n_fontSize': ctx.i18n_fontSize.value },
          modelValue: ctx.configureDownSound.value,
          'onUpdate:modelValue': (val: boolean) => {
            ctx.configureDownSound.value = val;
          },
          'backdrop-filter': 'invert(70%)',
          onMouseup: ctx.preventDefaultMouseWhenRecording,
        },
        [
          h('q-card', { class: 'min-w-[80%]' }, [
            // 标题
            h(
              'q-card-section',
              { class: 'row items-center q-pb-none text-h6' },
              ctx.$t('KeyToneAlbum.craftKeySounds.configureDownSound')
            ),
            // 播放模式选择
            h('q-card-section', [
              h('q-select', {
                outlined: true,
                'stack-label': true,
                'virtual-scroll-slice-size': 999999,
                'popup-content-class':
                  'w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50',
                modelValue: ctx.playModeForDown.value,
                'onUpdate:modelValue': (val: string) => {
                  ctx.playModeForDown.value = val;
                },
                options: ctx.playModeOptions,
                'option-label': (item: string) => ctx.$t(ctx.playModeLabels.get(item) || ''),
                label: ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode'),
                dense: true,
              }),
            ]),
            // 声音选择
            h('q-card-section', [
              h('q-select', {
                outlined: true,
                'stack-label': true,
                'virtual-scroll-slice-size': 999999,
                'popup-content-class':
                  'w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50',
                modelValue: ctx.selectedSoundsForDown.value,
                'onUpdate:modelValue': (val: any[]) => {
                  ctx.selectedSoundsForDown.value = val;
                },
                options: ctx.downSoundList.value,
                'option-label': ctx.album_options_select_label,
                'option-value': getOptionValue,
                label: ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds'),
                multiple: true,
                'use-chips': true,
                class: 'zl-ll',
                dense: true,
                'max-values': ctx.maxSelectionForDown.value,
                counter: true,
                'error-message': ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode'),
                error: ctx.playModeForDown.value === 'single' ? ctx.selectedSoundsForDown.value.length > 1 : false,
              }),
              // 类型选择组
              h('div', { class: 'h-10' }, [
                h('q-option-group', {
                  dense: true,
                  modelValue: ctx.downTypeGroup.value,
                  'onUpdate:modelValue': (val: string[]) => {
                    ctx.downTypeGroup.value = val;
                  },
                  options: ctx.options,
                  type: 'checkbox',
                  class: 'absolute left-8',
                }),
              ]),
            ]),
            // 关闭按钮
            h('q-card-actions', { align: 'right' }, [
              h(
                'q-btn',
                {
                  flat: true,
                  label: ctx.$t('KeyToneAlbum.close'),
                  color: 'primary',
                  onClick: () => {
                    ctx.configureDownSound.value = false;
                  },
                },
                null
              ),
            ]),
          ]),
        ]
      );
  },
});

// ============================================================================
// 配置抬起声音子对话框（内嵌组件）
// ============================================================================
const ConfigureUpSoundSubDialog = defineComponent({
  name: 'ConfigureUpSoundSubDialog',
  setup() {
    return () =>
      h(
        'q-dialog',
        {
          style: { '--i18n_fontSize': ctx.i18n_fontSize.value },
          modelValue: ctx.configureUpSound.value,
          'onUpdate:modelValue': (val: boolean) => {
            ctx.configureUpSound.value = val;
          },
          'backdrop-filter': 'invert(70%)',
          onMouseup: ctx.preventDefaultMouseWhenRecording,
        },
        [
          h('q-card', { class: 'min-w-[80%]' }, [
            // 标题
            h(
              'q-card-section',
              { class: 'row items-center q-pb-none text-h6' },
              ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound')
            ),
            // 播放模式选择
            h('q-card-section', [
              h('q-select', {
                outlined: true,
                'stack-label': true,
                'virtual-scroll-slice-size': 999999,
                'popup-content-class':
                  'w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50',
                modelValue: ctx.playModeForUp.value,
                'onUpdate:modelValue': (val: string) => {
                  ctx.playModeForUp.value = val;
                },
                options: ctx.playModeOptions,
                'option-label': (item: string) => ctx.$t(ctx.playModeLabels.get(item) || ''),
                label: ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode'),
                dense: true,
              }),
            ]),
            // 声音选择
            h('q-card-section', [
              h('q-select', {
                outlined: true,
                'stack-label': true,
                'virtual-scroll-slice-size': 999999,
                'popup-content-class':
                  'w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50',
                modelValue: ctx.selectedSoundsForUp.value,
                'onUpdate:modelValue': (val: any[]) => {
                  ctx.selectedSoundsForUp.value = val;
                },
                options: ctx.upSoundList.value,
                'option-label': ctx.album_options_select_label,
                'option-value': getOptionValue,
                label: ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds'),
                multiple: true,
                'use-chips': true,
                class: 'zl-ll',
                dense: true,
                'max-values': ctx.maxSelectionForUp.value,
                counter: true,
                'error-message': ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode'),
                error: ctx.playModeForUp.value === 'single' ? ctx.selectedSoundsForUp.value.length > 1 : false,
              }),
              // 类型选择组
              h('div', { class: 'h-10' }, [
                h('q-option-group', {
                  dense: true,
                  modelValue: ctx.upTypeGroup.value,
                  'onUpdate:modelValue': (val: string[]) => {
                    ctx.upTypeGroup.value = val;
                  },
                  options: ctx.options,
                  type: 'checkbox',
                  class: 'absolute left-8',
                }),
              ]),
            ]),
            // 关闭按钮
            h('q-card-actions', { align: 'right' }, [
              h(
                'q-btn',
                {
                  flat: true,
                  label: ctx.$t('KeyToneAlbum.close'),
                  color: 'primary',
                  onClick: () => {
                    ctx.configureUpSound.value = false;
                  },
                },
                null
              ),
            ]),
          ]),
        ]
      );
  },
});

// ============================================================================
// 辅助函数
// ============================================================================

/**
 * 获取选项的唯一标识值
 */
function getOptionValue(item: any): string {
  if (item.type === 'audio_files') {
    return item.value.sha256 + item.value.name_id;
  }
  if (item.type === 'sounds') {
    return item.value.soundKey;
  }
  if (item.type === 'key_sounds') {
    return item.value.keySoundKey;
  }
  return '';
}

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 保存按键音配置
 * 保存成功后关闭对话框并重置表单
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
