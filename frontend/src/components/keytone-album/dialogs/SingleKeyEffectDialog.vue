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
文件说明: dialogs/SingleKeyEffectDialog.vue - 单键声效设置对话框
============================================================================

【文件作用】
本组件是联动声效功能（Step4）中的"单键声效"设置对话框，负责：
1. 显示已配置的单键声效列表（以键盘按键芯片形式展示）
2. 添加新的单键声效配置
3. 编辑/删除已有的单键声效配置

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(Context)
    │
    └── steps/StepLinkageEffects.vue (Step 4)
      │
      └── dialogs/SingleKeyEffectDialog.vue  <── 当前文件
                        │
                        ├── 添加单键声效子对话框 (isShowAddOrSettingSingleKeyEffectDialog)
                        └── 编辑单键声效子对话框 (isShowSingleKeySoundEffectEditDialog)

【子对话框说明】
1. 添加单键声效子对话框：
   - 支持多选按键（通过输入搜索或键盘录制）
   - 为选中的按键配置统一的声效

2. 编辑单键声效子对话框：
   - 编辑单个按键的声效配置
   - 支持删除该按键的声效配置

【数据流】
  父组件状态                              本组件使用方式
  ─────────────────────────────────────────────────────────────
  ctx.showSingleKeyEffectDialog      ->
v-model 控制主对话框显示 ctx.isShowAddOrSettingSingleKeyEffectDialog -> 添加对话框显示
ctx.isShowSingleKeySoundEffectEditDialog -> 编辑对话框显示 ctx.keysWithSoundEffect -> 已配置声效的按键 Map
ctx.selectedSingleKeys -> 选中的按键列表 ctx.isRecordingSingleKeys -> 是否正在录制按键
ctx.saveSingleKeySoundEffectConfig() -> 保存单键声效配置 【关联文件】 - ../types.ts : 类型定义 -
../../Keytone_album.vue : 父组件 - ../../DependencyWarning.vue : 依赖警告组件

【当前状态】
✅ 本组件已集成：由 `StepLinkageEffects` 渲染，并通过 `ctx.showSingleKeyEffectDialog` 的 v-model 控制显示/隐藏。

============================================================================
-->

<template>
  <!-- 单键声效主对话框 -->
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.showSingleKeyEffectDialog.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card>
      <!-- 对话框标题 -->
      <q-card-section class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        {{ ctx.$t('KeyToneAlbum.linkageEffects.singleKeySettings') }}
      </q-card-section>

      <!-- 说明文字 -->
      <q-card-section class="q-pt-none pb-0">
        <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.description') }}
        </div>

        <!-- 添加按钮和子对话框 -->
        <div class="flex flex-row items-center gap-2 mb-2 ml-2">
          <q-btn
            flat
            round
            color="primary"
            icon="add"
            @click="ctx.isShowAddOrSettingSingleKeyEffectDialog.value = true"
          />
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.addSingleKeyEffect') }}

          <!-- 添加单键声效子对话框 -->
          <AddSingleKeyEffectSubDialog />
        </div>
      </q-card-section>

      <!-- 已配置的单键声效列表 -->
      <q-card-section>
        <div v-if="ctx.keysWithSoundEffect.value.size === 0" class="text-[1.06rem]">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.noSingleKeyEffects') }}
        </div>
        <div v-else class="text-[1.06rem] pb-2 font-600 text-gray-700 flex flex-row items-center">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.singleKeyEffects') }}
          <div class="text-[0.88rem] ml-1">({{ ctx.$t('KeyToneAlbum.linkageEffects.single.dialog.clickToView') }})</div>
        </div>

        <!-- 按键芯片列表 -->
        <div class="flex flex-wrap gap-0.8">
          <q-chip
            v-for="item in ctx.keysWithSoundEffect.value"
            :key="item[0]"
            dense
            square
            class="p-t-3.25 p-b-3.25 p-x-2.5 bg-gradient-to-b from-gray-50 to-gray-200 border-2 border-gray-300 rounded-[0.18rem] shadow-[1px_2px_1px_3px_rgba(0,0,0,0.2),inset_1px_1px_1px_rgba(255,255,255,0.6)] inset_1px_1px_1px_rgba(255,255,255,0.6)]"
            clickable
            @click="handleEditKey(item)"
          >
            {{ keyEvent_store.dikCodeToName.get(Number(item[0])) || 'Dik-{' + item[0] + '}' }}
          </q-chip>

          <!-- 编辑单键声效子对话框 -->
          <EditSingleKeyEffectSubDialog />
        </div>
      </q-card-section>

      <!-- 关闭按钮 -->
      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="ctx.$t('KeyToneAlbum.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
/**
 * SingleKeyEffectDialog.vue - 单键声效设置对话框
 *
 * 【组件职责】
 * - 管理单键声效配置的主界面
 * - 展示已配置的按键列表
 * - 协调添加/编辑子对话框的显示
 *
 * 【子组件说明】
 * 由于单键声效功能复杂，本组件进一步拆分为两个子组件：
 * - AddSingleKeyEffectSubDialog: 添加单键声效
 * - EditSingleKeyEffectSubDialog: 编辑单键声效
 */

import { inject } from 'vue';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import { useKeyEventStore } from 'src/stores/keyEvent-store';
import AddSingleKeyEffectSubDialog from './AddSingleKeyEffectSubDialog.vue';
import EditSingleKeyEffectSubDialog from './EditSingleKeyEffectSubDialog.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;
const keyEvent_store = useKeyEventStore();

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 处理点击编辑按键
 * 打开编辑对话框并初始化数据
 */
function handleEditKey(item: [string, any]) {
  // 打开查看声效的对话框
  ctx.isShowSingleKeySoundEffectEditDialog.value = true;

  // 记录旧值用于判断
  ctx.currentEditingKey_old = ctx.currentEditingKey.value;
  ctx.currentEditingKey.value = Number(item[0]);

  // 如果点击的是不同的按键，则重新初始化
  if (ctx.currentEditingKey.value !== ctx.currentEditingKey_old) {
    // 为了防止初始有'key_sounds'时触发的锚定，会影响原数据的初始化
    ctx.singleKeyTypeGroup_edit.value = ['sounds'];

    const down = item[1].down;
    const up = item[1].up;

    ctx.keyDownSingleKeySoundEffectSelect_edit.value = ctx.convertValue(down ? down : '');
    ctx.keyUpSingleKeySoundEffectSelect_edit.value = ctx.convertValue(up ? up : '');

    ctx.keyDownSingleKeySoundEffectSelect_edit_old = ctx.keyDownSingleKeySoundEffectSelect_edit.value;
    ctx.keyUpSingleKeySoundEffectSelect_edit_old = ctx.keyUpSingleKeySoundEffectSelect_edit.value;
  }
}
</script>

<style lang="scss" scoped>
/**
 * SingleKeyEffectDialog 组件样式
 */

// 按钮样式
.q-btn {
  @apply text-xs;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}
</style>
