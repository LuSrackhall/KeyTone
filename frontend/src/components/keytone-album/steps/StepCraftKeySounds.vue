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
文件说明: steps/StepCraftKeySounds.vue - 步骤3：制作按键音
============================================================================

【文件作用】
本组件是键音专辑编辑器的第三个步骤（Step 3），负责：
1. 创建新的按键音（组合声音、源文件或其他按键音）
2. 配置按下(down)和抬起(up)时的声音
3. 设置播放模式（单次、随机、循环）
4. 编辑已有的按键音配置
5. 显示依赖警告（当引用的声音被删除时）

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)
        │
        ▼
  ┌─────────────────────────────────────────────────────┐
  │  q-stepper                                          │
  │  ├── StepLoadAudioFiles.vue       (Step 1)          │
  │  ├── StepDefineSounds.vue         (Step 2)          │
  │  ├── StepCraftKeySounds.vue  <── 当前文件 (Step 3)  │
  │  └── StepLinkageEffects.vue       (Step 4) [待创建] │
  └─────────────────────────────────────────────────────┘

【数据流】
本组件不持有任何状态，所有状态都通过 inject 从父组件获取：

  父组件状态                          本组件使用方式
  ─────────────────────────────────────────────────────────
  ctx.step                       -->
控制当前步骤 ctx.keySoundList --> 已定义的按键音列表 ctx.createNewKeySound --> 控制"创建按键音"对话框
ctx.editExistingKeySound --> 控制"编辑按键音"对话框 ctx.selectedSoundsForDown/Up --> 按下/抬起时选中的声音
ctx.playModeForDown/Up --> 播放模式 ctx.downSoundList/upSoundList --> 可选择的声音列表 ctx.saveKeySoundConfig() -->
保存按键音配置 ctx.deleteKeySound() --> 删除按键音 ctx.dependencyIssues --> 依赖问题列表 【关联文件】 - ../types.ts :
类型定义 - ../../Keytone_album.vue : 父组件 - ../../DependencyWarning.vue : 依赖警告组件 【当前状态】 ⚠️
注意：本组件已创建但尚未集成到父组件中！ ============================================================================
-->

<template>
  <q-step
    :name="3"
    :title="ctx.$t('KeyToneAlbum.craftKeySounds.title')"
    icon="piano"
    :done="ctx.keySoundList.value.length !== 0"
    :disable="ctx.step.value === 99 && ctx.keySoundList.value.length === 0"
    :header-nav="false"
    @click="handleStepClick"
  >
    <!-- 步骤说明 -->
    <div :class="['mb-3', ctx.step_introduce_fontSize.value]">
      {{ ctx.$t('KeyToneAlbum.craftKeySounds.description') }}
      <q-icon name="info" color="primary">
        <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
          <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.keySoundExplain') }}</div>
          <div>{{ ctx.$t('KeyToneAlbum.craftKeySounds.tooltip.downUpExplain') }}</div>
        </q-tooltip>
      </q-icon>
    </div>

    <!-- 制作按键音的业务逻辑 -->
    <div>
      <!-- 制作新的按键音 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.craftKeySounds.createNewKeySound')"
          @click="ctx.createNewKeySound.value = !ctx.createNewKeySound.value"
        />

        <!-- TODO: 创建按键音对话框 - 待抽离为独立组件 dialogs/CreateKeySoundDialog.vue -->
        <!-- 当前保持原有实现，后续迁移 -->
      </div>

      <div :class="['p-2 text-zinc-600']">{{ ctx.$t('KeyToneAlbum.or') }}</div>

      <!-- 编辑已有按键音 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.craftKeySounds.editExistingKeySound')"
          @click="handleEditKeySound"
        />

        <!-- TODO: 编辑按键音对话框 - 待抽离为独立组件 dialogs/EditKeySoundDialog.vue -->
        <!-- 当前保持原有实现，后续迁移 -->
      </div>
    </div>

    <!-- 导航按钮 -->
    <q-stepper-navigation>
      <q-btn @click="ctx.step.value = 4" color="primary" :label="ctx.$t('KeyToneAlbum.continue')" />
      <q-btn flat @click="ctx.step.value = 2" color="primary" :label="ctx.$t('KeyToneAlbum.back')" class="q-ml-sm" />
    </q-stepper-navigation>
  </q-step>
</template>

<script setup lang="ts">
/**
 * StepCraftKeySounds.vue - 步骤3：制作按键音
 *
 * 【组件职责】
 * - 只负责 UI 渲染和用户交互
 * - 所有状态和业务逻辑通过 inject 从父组件获取
 * - 不持有任何本地状态（除了纯 UI 状态）
 *
 * 【按键音概念】
 * 按键音(KeySound)是由多个声音(Sound)或音频文件(AudioFile)组合而成的：
 * - down: 按键按下时播放的声音配置
 * - up: 按键抬起时播放的声音配置
 * - 每个配置包含：模式(single/random/loop)和声音列表
 *
 * 【依赖警告】
 * 当按键音引用的声音或音频文件被删除时，会显示依赖警告。
 * 使用 DependencyWarning 组件来展示这些警告。
 */

import { inject } from 'vue';
import { useQuasar } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';

const q = useQuasar();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 处理 step header 点击事件
 */
function handleStepClick(event: MouseEvent) {
  const header = (event.target as HTMLElement).closest('.q-stepper__tab');
  if (header) {
    ctx.step.value = ctx.step.value === 3 ? 99 : 3;
  }
}

/**
 * 处理"编辑已有按键音"按钮点击
 */
function handleEditKeySound() {
  if (ctx.keySoundList.value.length === 0) {
    q.notify({
      type: 'warning',
      message: ctx.$t('KeyToneAlbum.notify.noKeySoundsToEdit'),
      position: 'top',
    });
    return;
  }
  ctx.editExistingKeySound.value = true;
}
</script>
