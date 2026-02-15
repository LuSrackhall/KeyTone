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
文件说明: dialogs/EditKeySoundDialog.vue - 编辑按键音对话框
============================================================================

【文件作用】
本组件是一个可复用的对话框，用于编辑已有的按键音定义。
功能包括：
1. 从按键音列表中选择要编辑的按键音
2. 修改按键音名称
3. 重新配置按下(down)时的声音（模式 + 声音列表）
4. 重新配置抬起(up)时的声音（模式 + 声音列表）
5. 保存或删除按键音
6. 显示依赖警告（当引用的声音被删除时）

【特殊说明】
- 在选项列表中显示依赖警告（DependencyWarning）
- 选中按键音后展示详细编辑卡片
- 包含嵌套的子对话框用于配置 down/up 声音

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
        │
        ├── steps/StepCraftKeySounds.vue
        │         │
        │         └── dialogs/EditKeySoundDialog.vue  <── 当前文件
        │                 │
        │                 ├── (内嵌) EditConfigureDownSoundDialog
        │                 └── (内嵌) EditConfigureUpSoundDialog

【数据流】
  父组件状态                              本组件使用方式
  ─────────────────────────────────────────────────────────────
  ctx.editExistingKeySound           -> v-model 控制对话框显示
  ctx.keySoundList                   -> 可选择的按键音列表
  ctx.selectedKeySound               -> 当前选中的按键音
  ctx.edit_configureDownSound        -> 控制按下声音编辑子对话框
  ctx.edit_configureUpSound          -> 控制抬起声音编辑子对话框
  ctx.edit_downSoundList/upSoundList -> 可选择的声音列表
  ctx.edit_downTypeGroup/upTypeGroup -> 类型筛选（音频文件/声音/按键音）
  ctx.saveKeySoundConfig()           -> 保存按键音
  ctx.deleteKeySound()               -> 删除按键音
  ctx.dependencyIssues               -> 依赖问题列表

【关联文件】
- ../types.ts : 类型定义，包含 KEYTONE_ALBUM_CONTEXT_KEY
- ../steps/StepCraftKeySounds.vue : 使用此对话框的 Step 组件
- ../../DependencyWarning.vue : 依赖警告组件

【当前状态】
✅ 本组件已集成到父组件中。

============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.editExistingKeySound.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="['min-w-[100%]', { 'mr-0': isMac }]">
      <!-- 对话框标题（sticky 置顶） -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ ctx.$t('KeyToneAlbum.craftKeySounds.editExistingKeySound') }}
      </q-card-section>

      <!-- 按键音选择下拉框 -->
      <q-card-section>
        <q-select
          outlined
          stack-label
          clearable
          :virtual-scroll-slice-size="999999"
          popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
          v-model="ctx.selectedKeySound.value"
          :options="ctx.keySoundList.value"
          :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectKeySoundToEdit')"
          :option-label="(item: any) => item.keySoundValue.name"
          :option-value="(item: any) => item.keySoundKey"
          dense
        >
          <!-- 自定义选项渲染，包含依赖警告 -->
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ scope.opt.keySoundValue.name }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <DependencyWarning
                  :issues="ctx.dependencyIssues.value"
                  item-type="key_sounds"
                  :item-id="scope.opt.keySoundKey"
                  :show-details="false"
                />
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </q-card-section>

      <!-- 选中按键音的编辑卡片 -->
      <q-card-section
        :class="['flex flex-col -m-t-2']"
        v-if="ctx.selectedKeySound.value?.keySoundKey !== '' && ctx.selectedKeySound.value !== undefined"
      >
        <q-card :class="['flex flex-col pb-3', { 'mr-0': isMac }]" v-if="ctx.selectedKeySound.value">
          <!-- 按键音名称输入 -->
          <q-card-section :class="['p-b-1 mt-3']">
            <q-input
              outlined
              stack-label
              dense
              v-model="ctx.selectedKeySound.value.keySoundValue.name"
              :label="ctx.$t('KeyToneAlbum.craftKeySounds.keySoundName')"
              :placeholder="ctx.$t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder')"
            />

            <!-- 配置按钮组 -->
            <div class="flex flex-col mt-1">
              <!-- 配置按下声音按钮 -->
              <q-btn
                :class="['bg-zinc-300 my-7 w-88% self-center']"
                :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureDownSound')"
                @click="ctx.edit_configureDownSound.value = true"
              />

              <!-- 配置按下声音编辑子对话框 -->
              <q-dialog
                :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
                v-model="ctx.edit_configureDownSound.value"
                backdrop-filter="invert(70%)"
                @mouseup="ctx.preventDefaultMouseWhenRecording"
              >
                <q-card :class="['min-w-[80%]', { 'mr-0': isMac }]">
                    <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
                      {{ ctx.$t('KeyToneAlbum.craftKeySounds.configureDownSound') }}
                    </q-card-section>
                  <q-card-section>
                    <q-select
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.selectedKeySound.value.keySoundValue.down.mode"
                      :options="ctx.playModeOptions"
                      :option-label="(item: any) => ctx.$t(ctx.playModeLabels.get(item) || '')"
                      :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                      dense
                    />
                  </q-card-section>
                  <q-card-section class="pb-8">
                    <q-select
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.selectedKeySound.value.keySoundValue.down.value"
                      :options="ctx.edit_downSoundList.value"
                      :option-label="ctx.album_options_select_label"
                      :option-value="getOptionValue"
                      :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                      multiple
                      use-chips
                      :class="['zl-ll']"
                      dense
                      :max-values="getPlayMode(ctx.selectedKeySound.value.keySoundValue.down.mode) === 'single' ? 1 : Infinity"
                      counter
                      :error-message="ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                      :error="
                        getPlayMode(ctx.selectedKeySound.value.keySoundValue.down.mode) === 'single' &&
                        ctx.selectedKeySound.value.keySoundValue.down.value.length > 1
                      "
                      ref="edit_downSoundSelectDom"
                      @update:model-value="edit_downSoundSelectDom?.hidePopup()"
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
                              :item-id="
                                scope.opt.type === 'audio_files'
                                  ? `${scope.opt.value?.sha256}:${scope.opt.value?.name_id}`
                                  : scope.opt.type === 'sounds'
                                  ? scope.opt.value?.soundKey
                                  : scope.opt.value?.keySoundKey
                              "
                              :show-details="false"
                            />
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>
                    <div class="h-10">
                      <q-option-group
                        dense
                        v-model="ctx.edit_downTypeGroup.value"
                        :options="ctx.options"
                        type="checkbox"
                        class="absolute left-8"
                      >
                        <template #label-0="props">
                          <q-item-label>
                            {{ ctx.$t(props.label) }}
                            <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                              <q-tooltip
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                  <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                    <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                  </q-card-actions>
                </q-card>
              </q-dialog>

              <!-- 配置抬起声音按钮 -->
              <q-btn
                :class="['bg-zinc-300 m-b-7 w-88% self-center']"
                :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
                @click="ctx.edit_configureUpSound.value = true"
              />

              <!-- 配置抬起声音编辑子对话框 -->
              <q-dialog
                :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
                v-model="ctx.edit_configureUpSound.value"
                backdrop-filter="invert(70%)"
                @mouseup="ctx.preventDefaultMouseWhenRecording"
              >
                <q-card :class="['min-w-[80%]', { 'mr-0': isMac }]">
                  <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
                    {{ ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound') }}
                  </q-card-section>
                  <q-card-section>
                    <q-select
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.selectedKeySound.value.keySoundValue.up.mode"
                      :options="ctx.playModeOptions"
                      :option-label="(item: any) => ctx.$t(ctx.playModeLabels.get(item) || '')"
                      :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                      dense
                    />
                  </q-card-section>
                  <q-card-section class="pb-8">
                    <q-select
                      outlined
                      stack-label
                      :virtual-scroll-slice-size="999999"
                      popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      v-model="ctx.selectedKeySound.value.keySoundValue.up.value"
                      :options="ctx.edit_upSoundList.value"
                      :option-label="ctx.album_options_select_label"
                      :option-value="getOptionValue"
                      :label="ctx.$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                      multiple
                      use-chips
                      :class="['zl-ll']"
                      dense
                      :max-values="getPlayMode(ctx.selectedKeySound.value.keySoundValue.up.mode) === 'single' ? 1 : Infinity"
                      counter
                      :error-message="ctx.$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                      :error="
                        getPlayMode(ctx.selectedKeySound.value.keySoundValue.up.mode) === 'single' &&
                        ctx.selectedKeySound.value.keySoundValue.up.value.length > 1
                      "
                      ref="edit_upSoundSelectDom"
                      @update:model-value="edit_upSoundSelectDom?.hidePopup()"
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
                              :item-id="
                                scope.opt.type === 'audio_files'
                                  ? `${scope.opt.value?.sha256}:${scope.opt.value?.name_id}`
                                  : scope.opt.type === 'sounds'
                                  ? scope.opt.value?.soundKey
                                  : scope.opt.value?.keySoundKey
                              "
                              :show-details="false"
                            />
                          </q-item-section>
                        </q-item>
                      </template>
                    </q-select>
                    <div class="h-10">
                      <q-option-group
                        dense
                        v-model="ctx.edit_upTypeGroup.value"
                        :options="ctx.options"
                        type="checkbox"
                        class="absolute left-8"
                      >
                        <template #label-0="props">
                          <q-item-label>
                            {{ ctx.$t(props.label) }}
                            <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                              <q-tooltip
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                                :class="[
                                  'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                ]"
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
                  <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                    <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                  </q-card-actions>
                </q-card>
              </q-dialog>
            </div>
          </q-card-section>

          <!-- 操作按钮组 -->
          <q-card-section :class="['flex justify-center gap-4 -mt-3']">
            <!-- 保存按钮 -->
            <q-btn
              dense
              class="pr-2.3"
              color="primary"
              icon="save"
              :label="ctx.$t('KeyToneAlbum.confirmEdit')"
              @click="handleSave"
            />

            <!-- 删除按钮 -->
            <q-btn
              dense
              class="pr-2.3"
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
 * EditKeySoundDialog.vue - 编辑按键音对话框
 *
 * 【组件职责】
 * 提供编辑已有按键音的完整界面，包括：
 * - 从列表选择按键音
 * - 修改按键音属性
 * - 重新配置 down/up 声音（内嵌子对话框）
 * - 保存和删除功能
 *
 * 【注意事项】
 * selectedKeySound 的 down/up value 结构在父组件的 watch 中会被转换，
 * 以适配选择输入框组件的使用需求。
 * - value 中的 uuid 被转换为完整的对象引用
 * - mode 可能来自历史结构或 UI 中间态，既可能是字符串，也可能是 { mode: string }
 *   因此这里统一通过解析函数读取，避免保存时写回 undefined。
 * 详见父组件中的 watch(selectedKeySound) 实现。
 */

import { inject, ref, computed } from 'vue';
import { useQuasar, Platform } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';

const q = useQuasar();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY);
if (!ctx) {
  // 这里必须存在父组件的 provide；若不存在说明组件被错误地单独挂载
  throw new Error('EditKeySoundDialog requires KEYTONE_ALBUM_CONTEXT_KEY');
}

// ============================================================================
// DOM 引用
// ============================================================================
const edit_downSoundSelectDom = ref<{ hidePopup?: () => void } | null>(null);
const edit_upSoundSelectDom = ref<{ hidePopup?: () => void } | null>(null);

const isMac = computed(() => Platform.is.mac === true);

// ============================================================================
// 工具函数
// ============================================================================

/**
 * 获取选项的唯一值
 * 用于 q-select 组件的 option-value 属性
 *
 * 【说明】
 * 虽然 json 中的存储格式分别是：
 * - {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
 * - {key:'sounds', value:string} // 此处 value 是 soundKey
 * - {key:'key_sounds', value:string} // 此处 value 是 keySoundKey
 *
 * 但是，我们通过 watch 对当前组件的 model 做了变更，
 * 使其类型提前由 uuid 转换成了相关对象。
 * 因此，此处仍按照对应对象处理即可。
 */
type KeySoundOption = {
  type: 'audio_files' | 'sounds' | 'key_sounds';
  value?: {
    sha256?: string;
    name_id?: string;
    soundKey?: string;
    keySoundKey?: string;
  };
};

function getOptionValue(item: KeySoundOption) {
  if (item.type === 'audio_files') {
    return item.value?.sha256 + item.value?.name_id;
  }
  if (item.type === 'sounds') {
    return item.value?.soundKey;
  }
  if (item.type === 'key_sounds') {
    return item.value?.keySoundKey;
  }
}

/**
 * 解析播放模式值，统一为 string。
 *
 * 【兼容性说明】
 * - 旧数据或中间态可能使用 { mode: string }
 * - 现有 UI 期望使用 string
 * 这里做兼容解析，避免保存时写回 undefined。
 */
function getPlayMode(mode: unknown, fallback = 'random') {
  if (typeof mode === 'string') {
    return mode;
  }
  if (mode && typeof mode === 'object' && 'mode' in mode && typeof (mode as { mode?: unknown }).mode === 'string') {
    return (mode as { mode: string }).mode;
  }
  return fallback;
}

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 保存按键音修改
 *
 * 【重要】
 * selectedKeySound.keySoundValue.down/up 的结构已在父组件的 watch 中被转换，
 * 其 mode 可能是字符串或对象，因此这里必须使用解析函数读取。
 */
function handleSave() {
  if (!ctx.selectedKeySound.value) return;

  ctx.saveKeySoundConfig({
    key: ctx.selectedKeySound.value.keySoundKey,
    name: ctx.selectedKeySound.value.keySoundValue.name,
    down: {
      mode: getPlayMode(ctx.selectedKeySound.value.keySoundValue.down.mode),
      value: ctx.selectedKeySound.value.keySoundValue.down.value,
    },
    up: {
      mode: getPlayMode(ctx.selectedKeySound.value.keySoundValue.up.mode),
      value: ctx.selectedKeySound.value.keySoundValue.up.value,
    },
  });
}

/**
 * 删除按键音
 * 删除成功后清空选中状态并显示通知
 */
function handleDelete() {
  if (!ctx.selectedKeySound.value) return;

  ctx.deleteKeySound({
    keySoundKey: ctx.selectedKeySound.value.keySoundKey,
    onSuccess: () => {
      ctx.selectedKeySound.value = undefined;
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
 * EditKeySoundDialog 组件样式
 *
 * 【样式说明】
 * 本组件使用的样式与 CreateKeySoundDialog 基本一致，
 * 主要用于处理 Quasar 组件的溢出和滚动显示问题。
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
