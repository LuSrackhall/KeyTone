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
文件说明: dialogs/EditSingleKeyEffectSubDialog.vue - 编辑单键声效子对话框
============================================================================

【文件作用】
本组件是单键声效对话框的子组件，负责：
1. 编辑单个按键的声效配置
2. 支持修改按下/抬起声效
3. 支持删除该按键的声效配置

【数据流】
  父组件状态                              本组件使用方式
  ─────────────────────────────────────────────────────────────
  ctx.isShowSingleKeySoundEffectEditDialog ->
v-model ctx.currentEditingKey -> 当前编辑的按键 ctx.currentEditingKeyOfName -> 当前按键名称
ctx.keyDownSingleKeySoundEffectSelect_edit -> 按下声效选择 ctx.keyUpSingleKeySoundEffectSelect_edit -> 抬起声效选择
ctx.saveSingleKeySoundEffectConfig() -> 保存配置 【关联文件】 - ./SingleKeyEffectDialog.vue : 父对话框 - ../types.ts :
类型定义 - ../../DependencyWarning.vue : 依赖警告组件

【当前状态】
✅ 本组件已集成：由 `SingleKeyEffectDialog` 拉起，通过 `ctx.isShowSingleKeySoundEffectEditDialog` 控制显示/隐藏。

============================================================================ -->

<template>
  <!-- 编辑单键声效子对话框 -->
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.isShowSingleKeySoundEffectEditDialog.value"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card style="min-width: 350px" :class="[{ 'mr-0': isMac }]">
      <!-- 对话框标题 -->
      <q-card-section class="sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        <div class="text-base flex flex-row items-center">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.editSingleKey') }} -
          <div class="text-sm font-bold">[ {{ ctx.currentEditingKeyOfName.value }} ]</div>
          - {{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.soundEffect') }}
        </div>
      </q-card-section>

      <q-card-section class="q-pt-none pb-1">
        <!-- 声效编辑区域 -->
        <div class="w-full">
          <q-card-section>
            <div class="flex flex-row flex-nowrap items-center mb-3">
              <div class="flex flex-col space-y-4 w-full">
                <!-- 按下声效选择 -->
                <q-select
                  outlined
                  stack-label
                  :virtual-scroll-slice-size="999999"
                  popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                  v-model="ctx.keyDownSingleKeySoundEffectSelect_edit.value"
                  :options="ctx.keySingleKeySoundEffectOptions_edit.value"
                  :option-label="ctx.album_options_select_label"
                  :option-value="getOptionValue"
                  :label="`${ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.editSingleKey')} -[ ${
                    ctx.currentEditingKeyOfName.value
                  } ]- ${ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.soundEffect-down')} `"
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
                  v-model="ctx.keyUpSingleKeySoundEffectSelect_edit.value"
                  :options="ctx.keySingleKeySoundEffectOptions_edit.value"
                  :option-label="ctx.album_options_select_label"
                  :option-value="getOptionValue"
                  :label="`${ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.editSingleKey')} -[ ${
                    ctx.currentEditingKeyOfName.value
                  } ]- ${ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.soundEffect-up')} `"
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
              <div class="flex justify-end -m-l-2">
                <q-icon
                  @click="
                    ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value =
                      !ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value
                  "
                  size="2.75rem"
                  v-if="ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value"
                >
                  <template v-if="ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value">
                    <q-icon name="svguse:icons.svg#锚定"></q-icon>
                  </template>
                  <template v-else>
                    <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
                  </template>
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                    <span class="text-sm">{{
                      ctx.$t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title')
                    }}</span>
                    <span class="text-sm" v-if="ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value">
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
            <div class="h-16 m-l-9">
              <q-option-group dense v-model="ctx.singleKeyTypeGroup_edit.value" :options="ctx.options" type="checkbox">
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
                  <q-item-label>
                    {{ ctx.$t(props.label) }}
                    <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
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

            <!-- 操作按钮 -->
            <div class="flex justify-center gap-8 -m-l-3 m-t-5">
              <q-btn
                class="p-r-2"
                dense
                :label="ctx.$t('KeyToneAlbum.confirmEdit')"
                color="primary"
                icon="save"
                @click="handleSave"
              >
                <q-tooltip>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.save') }}</q-tooltip>
              </q-btn>
              <q-btn
                class="p-r-2"
                dense
                :label="ctx.$t('KeyToneAlbum.delete')"
                color="negative"
                icon="delete"
                v-close-popup
                @click="handleDelete"
              >
                <q-tooltip>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.delete') }}</q-tooltip>
              </q-btn>
            </div>
          </q-card-section>
        </div>
      </q-card-section>

      <!-- 关闭按钮 -->
      <q-card-actions class="pt-0" align="right">
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * EditSingleKeyEffectSubDialog.vue - 编辑单键声效子对话框
 *
 * 【组件职责】
 * - 编辑单个按键的声效配置
 * - 支持保存修改和删除配置
 */

import { inject, computed } from 'vue';
import { useQuasar, Platform } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';

const q = useQuasar();
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

const isMac = computed(() => Platform.is.mac === true);

// ============================================================================
// 工具函数
// ============================================================================

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

function getUuid(item: any) {
  if (item?.type === 'audio_files') {
    return item?.value.sha256 + item?.value.name_id;
  }
  if (item?.type === 'sounds') {
    return item?.value.soundKey;
  }
  if (item?.type === 'key_sounds') {
    return item?.value.keySoundKey;
  }
}

// ============================================================================
// 事件处理函数
// ============================================================================

function handleDownPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
    ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value &&
    ctx.keyDownSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds'
  ) {
    ctx.keyUpSingleKeySoundEffectSelect_edit.value = ctx.keyDownSingleKeySoundEffectSelect_edit.value;
  }
}

function handleUpPopupHide() {
  if (
    ctx.isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
    ctx.isAnchoringUltimatePerfectionKeySound_singleKey_edit.value &&
    ctx.keyUpSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds'
  ) {
    ctx.keyDownSingleKeySoundEffectSelect_edit.value = ctx.keyUpSingleKeySoundEffectSelect_edit.value;
  }
}

function handleSave() {
  const downUuid = getUuid(ctx.keyDownSingleKeySoundEffectSelect_edit.value);
  const upUuid = getUuid(ctx.keyUpSingleKeySoundEffectSelect_edit.value);
  const downOldUuid = getUuid(ctx.keyDownSingleKeySoundEffectSelect_edit_old);
  const upOldUuid = getUuid(ctx.keyUpSingleKeySoundEffectSelect_edit_old);

  if (downUuid !== downOldUuid || upUuid !== upOldUuid) {
    // 检查是否为删除操作（两个都为空）
    if (!ctx.keyDownSingleKeySoundEffectSelect_edit.value && !ctx.keyUpSingleKeySoundEffectSelect_edit.value) {
      // 显示确认对话框
      q.dialog({
        title: ctx.$t('KeyToneAlbum.notify.confirmDeleteSingleKeyEffect'),
        message: ctx.$t('KeyToneAlbum.notify.confirmDeleteSingleKeyEffectMessage'),
        ok: {
          label: ctx.$t('KeyToneAlbum.cancel'),
          color: 'primary',
          flat: true,
        },
        cancel: {
          label: ctx.$t('KeyToneAlbum.confirm'),
          color: 'primary',
          flat: true,
        },
        persistent: true,
        focus: 'cancel',
      }).onCancel(() => {
        doSave();
      });
    } else {
      doSave();
    }
  } else {
    q.notify({
      type: 'warning',
      position: 'top',
      message: ctx.$t('KeyToneAlbum.notify.noChangesDetected'),
      timeout: 2000,
    });
  }
}

function doSave() {
  ctx.saveSingleKeySoundEffectConfig(
    {
      singleKeys: ctx.currentEditingKey.value ? [ctx.currentEditingKey.value] : [],
      down: ctx.keyDownSingleKeySoundEffectSelect_edit.value,
      up: ctx.keyUpSingleKeySoundEffectSelect_edit.value,
    },
    () => {
      if (!ctx.keyDownSingleKeySoundEffectSelect_edit.value && !ctx.keyUpSingleKeySoundEffectSelect_edit.value) {
        q.notify({
          type: 'positive',
          position: 'top',
          message: ctx.$t('KeyToneAlbum.notify.deleteSuccess'),
          timeout: 2000,
        });
      } else {
        q.notify({
          type: 'positive',
          position: 'top',
          message: ctx.$t('KeyToneAlbum.notify.saveSuccess'),
          timeout: 2000,
        });
      }
      // 更新旧值记录
      ctx.keyDownSingleKeySoundEffectSelect_edit_old = ctx.keyDownSingleKeySoundEffectSelect_edit.value;
      ctx.keyUpSingleKeySoundEffectSelect_edit_old = ctx.keyUpSingleKeySoundEffectSelect_edit.value;
      // 关闭对话框
      ctx.isShowSingleKeySoundEffectEditDialog.value = false;
    }
  );
}

function handleDelete() {
  ctx.saveSingleKeySoundEffectConfig(
    {
      singleKeys: ctx.currentEditingKey.value ? [ctx.currentEditingKey.value] : [],
      down: null,
      up: null,
    },
    () => {
      q.notify({
        type: 'positive',
        position: 'top',
        message: ctx.$t('KeyToneAlbum.notify.deleteSuccess'),
        timeout: 2000,
      });
    }
  );
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
