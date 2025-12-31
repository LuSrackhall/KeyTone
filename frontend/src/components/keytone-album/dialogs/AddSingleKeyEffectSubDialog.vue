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
文件说明: dialogs/AddSingleKeyEffectSubDialog.vue - 添加单键声效子对话框
============================================================================

【文件作用】
本组件是单键声效对话框的子组件，负责：
1. 选择要配置的按键（支持多选、搜索、键盘录制）
2. 为选中的按键配置按下/抬起声效
3. 支持类型筛选（音频文件/声音/按键音）

【按键选择功能】
- 支持在输入框中搜索按键名称
- 支持键盘录制模式（点击键盘图标启用）
- 录制时直接按下按键即可添加到选择列表

【数据流】
  父组件状态                              本组件使用方式
  ─────────────────────────────────────────────────────────────
  ctx.isShowAddOrSettingSingleKeyEffectDialog -->
v-model ctx.selectedSingleKeys --> 选中的按键列表 ctx.isRecordingSingleKeys --> 是否正在录制
ctx.keyDownSingleKeySoundEffectSelect --> 按下声效选择 ctx.keyUpSingleKeySoundEffectSelect --> 抬起声效选择
ctx.saveSingleKeySoundEffectConfig() --> 保存配置 【关联文件】 - ./SingleKeyEffectDialog.vue : 父对话框 - ../types.ts :
类型定义 - ../../DependencyWarning.vue : 依赖警告组件 【当前状态】 ✅ 本组件已创建，待集成。
============================================================================ -->

<template>
  <!-- 添加单键声效子对话框 -->
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    :no-esc-dismiss="ctx.isRecordingSingleKeys.value && ctx.isGetsFocused.value"
    v-model="ctx.isShowAddOrSettingSingleKeyEffectDialog.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card>
      <!-- 对话框标题 -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ ctx.$t('KeyToneAlbum.linkageEffects.single.addSingleKeyEffect') }}
      </q-card-section>

      <q-card-section class="q-pt-none">
        <!-- 说明文字 -->
        <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.selectKeyAndEffect') }}
          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
              <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.fastDelete') }}</div>
            </q-tooltip>
          </q-icon>
        </div>

        <div class="flex flex-col gap-4">
          <div class="flex flex-row items-center gap-2 w-full">
            <!-- 按键选择器 -->
            <q-select
              :label="ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey')"
              ref="singleKeysSelectRef"
              popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
              v-model="ctx.selectedSingleKeys.value"
              :options="ctx.filterOptions.value"
              :virtual-scroll-slice-size="9999"
              dense
              filled
              hide-dropdown-icon
              multiple
              outlined
              stack-label
              :placeholder="
                ctx.isRecordingSingleKeys.value
                  ? ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey-placeholder_record')
                  : ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey-placeholder_search')
              "
              use-input
              use-chips
              :class="['zl-ll']"
              class="flex-1"
              @focus="ctx.isGetsFocused.value = true"
              @blur="ctx.isGetsFocused.value = false"
              :maxlength="ctx.isRecordingSingleKeys.value ? 0 : Infinity"
              @keydown="ctx.preventDefaultKeyBehaviorWhenRecording"
              :option-label="(option: number) => keyEvent_store.dikCodeToName.get(option) || 'Dik-{' + option + '}'"
              @filter="handleFilter"
              @remove="handleRemove"
            >
              <template v-slot:append>
                <q-btn
                  dense
                  flat
                  :color="ctx.isRecordingSingleKeys.value ? 'primary' : ''"
                  icon="keyboard"
                  @click="ctx.isRecordingSingleKeys.value = !ctx.isRecordingSingleKeys.value"
                >
                  <q-tooltip>
                    {{
                      ctx.isRecordingSingleKeys.value
                        ? ctx.$t('KeyToneAlbum.linkageEffects.tooltips.stopRecording')
                        : ctx.$t('KeyToneAlbum.linkageEffects.tooltips.startRecording')
                    }}
                  </q-tooltip>
                </q-btn>
              </template>
            </q-select>

            <!-- 声效选择区域 -->
            <div class="flex flex-row items-center justify-center gap-x-9 gap-y-2 w-[95%]">
              <q-checkbox
                dense
                :class="getCheckboxClass('down')"
                v-model="ctx.isDownSoundEffectSelectEnabled.value"
                :label="ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.downSoundEffect')"
              />
              <q-checkbox
                dense
                :class="getCheckboxClass('up')"
                v-model="ctx.isUpSoundEffectSelectEnabled.value"
                :label="ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.upSoundEffect')"
              />
            </div>

            <!-- 声效选择下拉框 -->
            <div class="w-full">
              <q-card-section>
                <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                  <div class="flex flex-col space-y-4 w-[223px]">
                    <!-- 按下声效选择 -->
                    <q-select
                      v-show="ctx.isDownSoundEffectSelectEnabled.value"
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.keyDownSingleKeySoundEffectSelect.value"
                      :options="ctx.keySingleKeySoundEffectOptions.value"
                      :option-label="ctx.album_options_select_label"
                      :option-value="getOptionValue"
                      :label="ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.setDownSoundEffect')"
                      use-chips
                      :class="['zl-ll']"
                      dense
                      @popup-hide="handleDownPopupHide"
                      class="max-w-full"
                    >
                      <template v-slot:option="scope">
                        <q-item v-bind="scope.itemProps">
                          <q-item-section>
                            <q-item-label>{{ ctx.album_options_select_label(scope.opt) }}</q-item-label>
                          </q-item-section>
                          <q-item-section side>
                            <DependencyWarning
                              :issues="ctx.dependencyIssues.value"
                              :item-type="scope.opt.type"
                              :item-id="getItemId(scope.opt)"
                              :show-details="false"
                            />
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>

                    <!-- 抬起声效选择 -->
                    <q-select
                      v-show="ctx.isUpSoundEffectSelectEnabled.value"
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.keyUpSingleKeySoundEffectSelect.value"
                      :options="ctx.keySingleKeySoundEffectOptions.value"
                      :option-label="ctx.album_options_select_label"
                      :option-value="getOptionValue"
                      :label="ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.setUpSoundEffect')"
                      use-chips
                      :class="['zl-ll']"
                      dense
                      @popup-hide="handleUpPopupHide"
                      class="max-w-full"
                    >
                      <template v-slot:option="scope">
                        <q-item v-bind="scope.itemProps">
                          <q-item-section>
                            <q-item-label>{{ ctx.album_options_select_label(scope.opt) }}</q-item-label>
                          </q-item-section>
                          <q-item-section side>
                            <DependencyWarning
                              :issues="ctx.dependencyIssues.value"
                              :item-type="scope.opt.type"
                              :item-id="getItemId(scope.opt)"
                              :show-details="false"
                            />
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>
                  </div>

                  <!-- 锚定图标 -->
                  <div
                    v-show="ctx.isDownSoundEffectSelectEnabled.value && ctx.isUpSoundEffectSelectEnabled.value"
                    :class="['absolute -right-2']"
                  >
                    <q-icon
                      @click="
                        ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value =
                          !ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value
                      "
                      size="2.75rem"
                      v-if="ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey.value"
                    >
                      <template v-if="ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value">
                        <q-icon name="svguse:icons.svg#锚定"></q-icon>
                      </template>
                      <template v-else>
                        <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
                      </template>
                      <q-tooltip
                        :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                      >
                        <span class="text-sm">{{
                          ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title')
                        }}</span>
                        <span class="text-sm" v-if="ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value">
                          {{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.anchored') }}<br />
                        </span>
                        <span class="text-sm" v-else>
                          {{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.unanchored') }}<br />
                        </span>
                        <span>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.tooltip') }}</span>
                      </q-tooltip>
                    </q-icon>
                  </div>
                </div>

                <!-- 类型筛选 -->
                <div
                  class="h-16 m-l-9"
                  v-show="ctx.isDownSoundEffectSelectEnabled.value || ctx.isUpSoundEffectSelectEnabled.value"
                >
                  <q-option-group dense v-model="ctx.singleKeyTypeGroup.value" :options="ctx.options" type="checkbox">
                    <template #label-0="props">
                      <q-item-label>
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.audioFile') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                    <template v-slot:label-1="props">
                      <q-item-label :class="getTypeGroupLabelClass()">
                        {{ ctx.$t(props.label) }}
                        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                          <q-tooltip
                            :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                          >
                            <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.soundList') }}</div>
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
                            <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.keySounds') }}</div>
                          </q-tooltip>
                        </q-icon>
                      </q-item-label>
                    </template>
                  </q-option-group>
                </div>
              </q-card-section>
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- 操作按钮 -->
      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="ctx.$t('KeyToneAlbum.linkageEffects.confirm')" color="primary" @click="handleConfirm" />
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * AddSingleKeyEffectSubDialog.vue - 添加单键声效子对话框
 *
 * 【组件职责】
 * - 提供按键选择界面（支持搜索和键盘录制）
 * - 为选中的多个按键配置统一的声效
 */

import { inject, ref } from 'vue';
import { useQuasar, QSelect } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import { useKeyEventStore } from 'src/stores/keyEvent-store';
import { useSettingStore } from 'src/stores/setting-store';
import DependencyWarning from '../../DependencyWarning.vue';

const q = useQuasar();
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;
const keyEvent_store = useKeyEventStore();
const setting_store = useSettingStore();

// 选择器引用
const singleKeysSelectRef = ref<QSelect>();

// ============================================================================
// 工具函数
// ============================================================================

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

function getItemId(opt: any) {
  if (opt.type === 'audio_files') {
    return `${opt.value?.sha256}:${opt.value?.name_id}`;
  }
  if (opt.type === 'sounds') {
    return opt.value?.soundKey;
  }
  if (opt.type === 'key_sounds') {
    return opt.value?.keySoundKey;
  }
  return '';
}

function getCheckboxClass(type: 'down' | 'up') {
  const lang = setting_store.languageDefault;
  if (['ru', 'pl'].includes(lang)) return 'w-[200px] text-nowrap';
  if (['fr', 'pt', 'pt-BR', 'vi'].includes(lang)) return 'w-[180px] text-nowrap';
  if (type === 'down' && ['it', 'tr'].includes(lang)) return 'w-[102px] text-nowrap';
  if (type === 'up' && lang === 'it') return 'w-[102px] text-nowrap';
  return '';
}

function getTypeGroupLabelClass() {
  const lang = setting_store.languageDefault;
  return ['fr', 'it', 'pt', 'pt-BR'].includes(lang) ? 'text-nowrap' : '';
}

// ============================================================================
// 事件处理函数
// ============================================================================

function handleFilter(inputValue: string, doneFn: (fn: () => void) => void) {
  if (inputValue === '') {
    doneFn(() => {
      ctx.filterOptions.value = ctx.keyOptions.value;
    });
    return;
  }

  doneFn(() => {
    const inputValueLowerCase = inputValue.toLowerCase();
    ctx.filterOptions.value = ctx.keyOptions.value.filter((item) => {
      const ifre = keyEvent_store.dikCodeToName.get(item)?.toLowerCase()?.indexOf(inputValueLowerCase);
      return ifre !== undefined && ifre > -1;
    });
  });
}

function handleRemove() {
  // 与旧实现保持一致：避免录制模式切换/移除芯片瞬间误记录鼠标行为
  ctx.setSingleKeyRecordingClearFlag();
  // 处理移除时保持焦点
  singleKeysSelectRef.value?.focus();
}

function handleDownPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
    ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value &&
    ctx.keyDownSingleKeySoundEffectSelect.value?.type === 'key_sounds'
  ) {
    ctx.keyUpSingleKeySoundEffectSelect.value = ctx.keyDownSingleKeySoundEffectSelect.value;
  }
}

function handleUpPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
    ctx.isAnchoringUltimatePerfectionKeySound_singleKey.value &&
    ctx.keyUpSingleKeySoundEffectSelect.value?.type === 'key_sounds'
  ) {
    ctx.keyDownSingleKeySoundEffectSelect.value = ctx.keyUpSingleKeySoundEffectSelect.value;
  }
}

function handleConfirm() {
  if (ctx.selectedSingleKeys.value.length !== 0) {
    // 保存之前先记录下keysWithSoundEffect的值
    const keysWithSoundEffect_old = [...ctx.keysWithSoundEffect.value];

    ctx.saveSingleKeySoundEffectConfig(
      {
        singleKeys: ctx.selectedSingleKeys.value,
        down: ctx.keyDownSingleKeySoundEffectSelect.value,
        up: ctx.keyUpSingleKeySoundEffectSelect.value,
      },
      () => {
        handleSaveCallback(keysWithSoundEffect_old);
      }
    );
  } else {
    q.notify({
      type: 'warning',
      position: 'top',
      message: ctx.$t('KeyToneAlbum.notify.selectKey'),
      timeout: 2000,
    });
  }
}

function handleSaveCallback(keysWithSoundEffect_old: [string, any][]) {
  if (!ctx.keyDownSingleKeySoundEffectSelect.value && !ctx.keyUpSingleKeySoundEffectSelect.value) {
    const isDeletes: number[] = [];
    const isNotSoundEffectSelect: number[] = [];

    ctx.selectedSingleKeys.value.forEach((key) => {
      if (keysWithSoundEffect_old.some((item) => item[0] === String(key))) {
        isDeletes.push(key);
      } else {
        isNotSoundEffectSelect.push(key);
      }
    });

    if (isNotSoundEffectSelect.length !== 0) {
      if (isDeletes.length === 0) {
        q.notify({
          type: 'warning',
          position: 'top',
          message: ctx.$t('KeyToneAlbum.notify.selectSoundEffect'),
          timeout: 2000,
        });
        return;
      } else {
        // 同时存在删除和未选择声效的情况
        let deleteString = '';
        let notSoundEffectSelectString = '';

        isDeletes.forEach((key) => {
          deleteString += '-[' + (keyEvent_store.dikCodeToName.get(key) || 'Dik-{' + key + '}') + ']';
        });
        isNotSoundEffectSelect.forEach((key) => {
          notSoundEffectSelectString += '-[' + (keyEvent_store.dikCodeToName.get(key) || 'Dik-{' + key + '}') + ']';
        });

        q.notify({
          type: 'warning',
          position: 'top',
          message: ctx.$t('KeyToneAlbum.notify.selectSoundEffectForKeys', { keys: notSoundEffectSelectString }),
          timeout: 8000,
        });
        q.notify({
          type: 'positive',
          position: 'top',
          message: ctx.$t('KeyToneAlbum.notify.deleteSuccessForKeys', { keys: deleteString }),
          timeout: 3000,
        });
        ctx.selectedSingleKeys.value = isNotSoundEffectSelect;
        return;
      }
    } else {
      // 纯删除操作
      q.notify({
        type: 'positive',
        position: 'top',
        message: ctx.$t('KeyToneAlbum.notify.deleteSuccess'),
        timeout: 2000,
      });
      ctx.selectedSingleKeys.value = [];
      ctx.keyDownSingleKeySoundEffectSelect.value = null;
      ctx.keyUpSingleKeySoundEffectSelect.value = null;
      ctx.isShowAddOrSettingSingleKeyEffectDialog.value = false;
      return;
    }
  }

  q.notify({
    type: 'positive',
    position: 'top',
    message: ctx.$t('KeyToneAlbum.notify.configSuccess'),
    timeout: 2000,
  });
  ctx.selectedSingleKeys.value = [];
  ctx.keyDownSingleKeySoundEffectSelect.value = null;
  ctx.keyUpSingleKeySoundEffectSelect.value = null;
  ctx.isShowAddOrSettingSingleKeyEffectDialog.value = false;
}
</script>

<style lang="scss" scoped>
.q-btn {
  @apply text-xs;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}

:deep(.q-field__native) {
  @apply max-w-full overflow-auto whitespace-nowrap;
  @apply h-5.8 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50 [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400;
}

:deep(.q-field__label) {
  @apply overflow-visible -ml-1.5 text-[0.8rem];
}

:deep(.q-placeholder) {
  @apply h-auto;
}

.zl-ll {
  :deep(.q-field__native) {
    @apply h-auto;
  }
  :deep(.q-field__messages) {
    @apply text-nowrap;
  }
}

:deep(.q-item__section) {
  @apply max-w-full overflow-auto whitespace-nowrap;
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

:deep(.ellipsis) {
  @apply max-w-full overflow-auto whitespace-nowrap text-clip;
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}
</style>
