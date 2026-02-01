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
文件说明: dialogs/EveryKeyEffectDialog.vue - 全键声效设置对话框
============================================================================

【文件作用】
本组件是联动声效功能（Step4）中的"全键声效"设置对话框，负责：
1. 选择全键按下时播放的声效（支持音频文件/声音/按键音）
2. 选择全键抬起时播放的声效
3. 筛选可选声效的类型（通过 option-group）
4. 至臻键音锚定功能（选择按键音时自动同步 down/up）

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
    │
    └── steps/StepLinkageEffects.vue (Step 4)
      │
      └── dialogs/EveryKeyEffectDialog.vue  <── 当前文件

【数据流】
  父组件状态                              本组件使用方式
  ─────────────────────────────────────────────────────────────
  ctx.showEveryKeyEffectDialog       ->
v-model 控制对话框显示 ctx.keyDownUnifiedSoundEffectSelect -> 按下声效选择 ctx.keyUpUnifiedSoundEffectSelect ->
抬起声效选择 ctx.keyUnifiedSoundEffectOptions -> 可选声效列表（已排序） ctx.unifiedTypeGroup ->
类型筛选（音频文件/声音/按键音） ctx.isShowUltimatePerfectionKeySoundAnchoring -> 是否显示锚定图标
ctx.isAnchoringUltimatePerfectionKeySound -> 锚定开关状态 ctx.saveUnifiedSoundEffectConfig() -> 保存全键声效配置
【锚定功能说明】 当选择"按键音"类型时，会出现锚定图标： - 锚定开启：选择 down 声效时自动将 up 设为相同值，反之亦然 -
锚定关闭：down 和 up 可独立选择 - 删除声效时的联动在父组件的 watch 中处理 【关联文件】 - ../types.ts : 类型定义 -
../../Keytone_album.vue : 父组件，包含状态和 saveUnifiedSoundEffectConfig 函数 - ../../DependencyWarning.vue :
依赖警告组件

【当前状态】
✅ 本组件已集成：由 `StepLinkageEffects` 渲染，并通过 `ctx.showEveryKeyEffectDialog` 的 v-model 控制显示/隐藏。

============================================================================
-->

<template>
  <!-- 全键声效设置对话框 -->
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.showEveryKeyEffectDialog.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="[{ 'mr-0': isMac }]">
      <!-- 对话框标题（sticky 置顶） -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ ctx.$t('KeyToneAlbum.linkageEffects.globalSettings') }}
      </q-card-section>

      <!-- 说明文字 -->
      <q-card-section class="q-pt-none">
        <div class="text-subtitle1 q-mb-md">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.global.description') }}
        </div>
      </q-card-section>

      <!-- 声效选择区域 -->
      <q-card-section>
        <!-- 选择全键按下/抬起声效的选项, 仅支持单选 -->
        <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
          <div class="flex flex-col space-y-4 w-7/8">
            <!-- 按下声效选择 -->
            <q-select
              outlined
              stack-label
              :virtual-scroll-slice-size="999999"
              popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
              v-model="ctx.keyDownUnifiedSoundEffectSelect.value"
              :options="ctx.keyUnifiedSoundEffectOptions.value"
              :option-label="ctx.album_options_select_label"
              :option-value="getOptionValue"
              :label="ctx.$t('KeyToneAlbum.linkageEffects.global.setKeyDownSound')"
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
              outlined
              stack-label
              :virtual-scroll-slice-size="999999"
              popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
              v-model="ctx.keyUpUnifiedSoundEffectSelect.value"
              :options="ctx.keyUnifiedSoundEffectOptions.value"
              :option-label="ctx.album_options_select_label"
              :option-value="getOptionValue"
              :label="ctx.$t('KeyToneAlbum.linkageEffects.global.setKeyUpSound')"
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

          <!-- 锚定图标：仅在选择了按键音类型时显示 -->
          <div class="flex justify-end -m-l-2">
            <q-icon
              @click="
                ctx.isAnchoringUltimatePerfectionKeySound.value = !ctx.isAnchoringUltimatePerfectionKeySound.value
              "
              size="2.75rem"
              v-if="ctx.isShowUltimatePerfectionKeySoundAnchoring.value"
            >
              <template v-if="ctx.isAnchoringUltimatePerfectionKeySound.value">
                <!-- 锚定 -->
                <q-icon name="svguse:icons.svg#锚定"></q-icon>
              </template>
              <template v-else>
                <!-- 锚定解除 -->
                <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
              </template>
              <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                <span class="text-sm">{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title') }}</span>
                <span class="text-sm" v-if="ctx.isAnchoringUltimatePerfectionKeySound.value">
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

        <!-- 类型筛选选项组 -->
        <div class="h-16 m-l-5.8">
          <q-option-group dense v-model="ctx.unifiedTypeGroup.value" :options="ctx.options" type="checkbox">
            <template #label-0="props">
              <q-item-label>
                {{ ctx.$t(props.label) }}
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                    <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.audioFile') }}</div>
                  </q-tooltip>
                </q-icon>
              </q-item-label>
            </template>
            <template v-slot:label-1="props">
              <q-item-label :class="[setting_store.languageDefault === 'es' ? 'text-nowrap' : '']">
                {{ ctx.$t(props.label) }}
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                    <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.soundList') }}</div>
                  </q-tooltip>
                </q-icon>
              </q-item-label>
            </template>
            <template v-slot:label-2="props">
              <q-item-label>
                {{ ctx.$t(props.label) }}
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                    <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.keySounds') }}</div>
                  </q-tooltip>
                </q-icon>
              </q-item-label>
            </template>
          </q-option-group>
        </div>
      </q-card-section>

      <!-- 操作按钮（sticky 置底） -->
      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn
          dense
          :label="ctx.$t('KeyToneAlbum.linkageEffects.confirm')"
          color="primary"
          v-close-popup
          @click="handleConfirm"
        />
        <q-btn flat dense color="primary" :label="ctx.$t('KeyToneAlbum.close')" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * EveryKeyEffectDialog.vue - 全键声效设置对话框
 *
 * 【组件职责】
 * - 提供全键声效（按下/抬起）的选择界面
 * - 支持类型筛选（音频文件/声音/按键音）
 * - 实现至臻键音锚定功能
 *
 * 【锚定功能工作原理】
 * 1. 当用户选择 down 声效且类型为 key_sounds 时，自动将 up 设为相同值
 * 2. 反之亦然
 * 3. 删除声效时的联动逻辑在父组件的 watch 中处理，避免循环依赖
 */

import { inject, computed } from 'vue';
import { useQuasar, Platform } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';
import { useSettingStore } from 'src/stores/setting-store';

const q = useQuasar();
const setting_store = useSettingStore();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

const isMac = computed(() => Platform.is.mac === true);

// ============================================================================
// 工具函数
// ============================================================================

/**
 * 获取选项的唯一值
 * 用于 q-select 组件的 option-value 属性
 */
function getOptionValue(item: any) {
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
 * 获取选项的 itemId（用于依赖警告）
 */
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

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 处理按下声效选择弹窗关闭
 * 实现锚定功能：当选择按键音时，自动同步到抬起声效
 *
 * 【重要】
 * 为避免循环依赖，此处仅处理"选择声效时"的锚定逻辑
 * "删除声效时"的联动逻辑在父组件的 watch 中处理
 */
function handleDownPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring.value &&
    ctx.isAnchoringUltimatePerfectionKeySound.value &&
    // 这里的 ? 是防止在勾选至臻键音的条件下，仅打开选项菜单且未做任何选择就关闭时，
    // null 值内没有 type 字段引起报错
    ctx.keyDownUnifiedSoundEffectSelect.value?.type === 'key_sounds'
  ) {
    ctx.keyUpUnifiedSoundEffectSelect.value = ctx.keyDownUnifiedSoundEffectSelect.value;
  }
}

/**
 * 处理抬起声效选择弹窗关闭
 * 实现锚定功能：当选择按键音时，自动同步到按下声效
 */
function handleUpPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring.value &&
    ctx.isAnchoringUltimatePerfectionKeySound.value &&
    ctx.keyUpUnifiedSoundEffectSelect.value?.type === 'key_sounds'
  ) {
    ctx.keyDownUnifiedSoundEffectSelect.value = ctx.keyUpUnifiedSoundEffectSelect.value;
  }
}

/**
 * 处理确认按钮点击
 * 保存全键声效配置并显示通知
 */
function handleConfirm() {
  ctx.saveUnifiedSoundEffectConfig(
    {
      down: ctx.keyDownUnifiedSoundEffectSelect.value,
      up: ctx.keyUpUnifiedSoundEffectSelect.value,
    },
    () => {
      q.notify({
        type: 'positive',
        position: 'top',
        message: ctx.$t('KeyToneAlbum.notify.configSuccess'),
        timeout: 2000,
      });
    }
  );
}
</script>

<style lang="scss" scoped>
/**
 * EveryKeyEffectDialog 组件样式
 *
 * 【样式说明】
 * 本组件的样式与其他 Dialog 组件保持一致
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

// 多选芯片选择框专用样式
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
