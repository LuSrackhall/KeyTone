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
3. 重新配置按下(down)时的声音
4. 重新配置抬起(up)时的声音
5. 保存或删除按键音

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
  ctx.editExistingKeySound           -->
v-model 控制对话框显示 ctx.keySoundList --> 可选择的按键音列表 ctx.selectedKeySound --> 当前选中的按键音
ctx.edit_configureDownSound --> 控制按下声音编辑子对话框 ctx.edit_configureUpSound --> 控制抬起声音编辑子对话框
ctx.edit_downSoundList/upSoundList --> 可选择的声音列表 ctx.edit_downTypeGroup/upTypeGroup --> 类型筛选
ctx.saveKeySoundConfig() --> 保存按键音 ctx.deleteKeySound() --> 删除按键音 ctx.dependencyIssues --> 依赖问题列表
【关联文件】 - ../types.ts : 类型定义 - ../steps/StepCraftKeySounds.vue : 使用此对话框的 Step 组件 -
../../DependencyWarning.vue : 依赖警告组件 ============================================================================
-->

<template>
  <q-dialog
    :style="{ '--i18n_fontSize': ctx.i18n_fontSize.value }"
    v-model="ctx.editExistingKeySound.value"
    backdrop-filter="invert(70%)"
    @mouseup="ctx.preventDefaultMouseWhenRecording"
  >
    <q-card :class="['min-w-[100%]']">
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
        <q-card :class="['flex flex-col pb-3']" v-if="ctx.selectedKeySound.value">
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

              <!-- TODO: 配置按下声音编辑子对话框 -->
              <!-- 由于结构复杂，暂保持在父组件中，后续迁移 -->

              <!-- 配置抬起声音按钮 -->
              <q-btn
                :class="['bg-zinc-300 m-b-7 w-88% self-center']"
                :label="ctx.$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
                @click="ctx.edit_configureUpSound.value = true"
              />

              <!-- TODO: 配置抬起声音编辑子对话框 -->
              <!-- 由于结构复杂，暂保持在父组件中，后续迁移 -->
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
 * - 重新配置 down/up 声音
 * - 保存和删除功能
 *
 * 【注意事项】
 * selectedKeySound 的 down/up value 结构在 watch 中会被转换，
 * 以适配选择输入框组件的使用需求。
 * 详见父组件中的 watch(selectedKeySound) 实现。
 */

import { inject } from 'vue';
import { useQuasar } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import DependencyWarning from '../../DependencyWarning.vue';

const q = useQuasar();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 保存按键音修改
 *
 * 【重要】
 * selectedKeySound.keySoundValue.down/up 的结构已在 watch 中被转换，
 * 其 mode 变成了 { mode: string } 对象形式。
 */
function handleSave() {
  if (!ctx.selectedKeySound.value) return;

  ctx.saveKeySoundConfig({
    key: ctx.selectedKeySound.value.keySoundKey,
    name: ctx.selectedKeySound.value.keySoundValue.name,
    down: {
      mode: ctx.selectedKeySound.value.keySoundValue.down.mode.mode,
      value: ctx.selectedKeySound.value.keySoundValue.down.value,
    },
    up: {
      mode: ctx.selectedKeySound.value.keySoundValue.up.mode.mode,
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
